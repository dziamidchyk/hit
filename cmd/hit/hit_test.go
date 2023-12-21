package main

import (
	"bytes"
	"flag"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

type testEnv struct {
	args           string
	stdout, stderr bytes.Buffer
}

func (e *testEnv) run() error {
	s := flag.NewFlagSet("hit", flag.ContinueOnError)
	s.SetOutput(&e.stderr)
	return run(s, strings.Fields(e.args), &e.stdout)
}

func TestRun(t *testing.T) {
	t.Parallel()

	happy := map[string]struct{ in, out string }{
		"url": {
			in:  "http://foo",
			out: "Making 100 GET requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()) + " (Timeout=1m0s).",
		},
		"n_c": {
			in:  "-n=20 -c=5 http://foo",
			out: "Making 20 GET requests to http://foo with a concurrency level of 5 (Timeout=1m0s).",
		},
		"t": {
			in:  "-t=5s http://foo",
			out: "Making 100 GET requests to http://foo with a concurrency level of 8 (Timeout=5s).",
		},
		"m": {
			in:  "-m=POST http://foo",
			out: "Making 100 POST requests to http://foo with a concurrency level of 8 (Timeout=1m0s).",
		},
		// "H": {
		// 	in:  "-H=\"Accept: text/json\" -H=\"User-agent test\" http://foo",
		// 	out: "Making 100 GET requests to http://foo with a concurrency level of 8 (Timeout=1m0s).",
		// },
	}
	sad := map[string]string{
		"url/missing": "",
		"url/err":     "://foo",
		"url/host":    "http://",
		"url/scheme":  "ftp://",
		"c/err":       "-c=x http://foo",
		"n/err":       "-n=x http://foo",
		"c/neg":       "-c=-1 http://foo",
		"n/neg":       "-n=-1 http://foo",
		"c/zero":      "-c=0 http://foo",
		"n/zero":      "-n=0 http://foo",
		"c/greater":   "-n=1 -c=2 http://foo",
		"t/err":       "-t=foo http://foo",
		"m/err":       "-m=DELETE http://foo",
	}
	for name, tt := range happy {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &testEnv{args: tt.in}
			if err := e.run(); err != nil {
				t.Fatalf("got %q;\nwant nil error", err)
			}
			if out := e.stdout.String(); !strings.Contains(out, tt.out) {
				t.Errorf("got:\n%s\nwant %q", out, tt.out)
			}
		})
	}
	for name, in := range sad {
		in := in
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &testEnv{args: in}
			if e.run() == nil {
				t.Fatal("got nil; want err")
			}
			if e.stderr.Len() == 0 {
				t.Fatal("stderr = 0 bytes; want > 0")
			}
		})
	}
}
