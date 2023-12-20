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
			out: "Making 100 requests to http://foo with a concurrency level of " + strconv.Itoa(runtime.NumCPU()),
		},
		"n_c": {
			in:  "-n=20 -c=5 http://foo",
			out: "Making 20 requests to http://foo with a concurrency level of 5",
		},
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
}
