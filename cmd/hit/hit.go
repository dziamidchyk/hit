package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/a-demidchik/hit/hit"
)

const bannerText = `
 __  __     __     ______
/\ \_\ \   /\ \   /\__  _\
\ \  __ \  \ \ \  \/_/\ \/
 \ \_\ \_\  \ \_\    \ \_\
  \/_/\/_/   \/_/     \/_/
`

func banner() string { return bannerText[1:] }

func main() {
	if err := run(flag.CommandLine, os.Args[1:], os.Stdout); err != nil {
		os.Exit(1)
	}
}

func run(s *flag.FlagSet, args []string, out io.Writer) error {
	f := &flags{
		n: 100,
		c: runtime.NumCPU(),
		t: time.Minute,
		m: "GET",
	}
	if err := f.parse(s, args); err != nil {
		return err
	}
	fmt.Fprintln(out, banner())
	if len(f.h) > 0 {
		fmt.Fprintf(out, "Headers: %v\n", strings.Join(f.h, ", "))
	}
	fmt.Fprintf(out, "Making %d %s requests to %s with a concurrency level of %d (Timeout=%v).\n", f.n, f.m, f.url, f.c, f.t)

	var sum hit.Result
	sum.Merge(&hit.Result{
		Bytes:    1000,
		Status:   http.StatusOK,
		Duration: time.Second,
	})
	sum.Merge(&hit.Result{
		Bytes:    1000,
		Status:   http.StatusOK,
		Duration: time.Second,
	})
	sum.Merge(&hit.Result{
		Status:   http.StatusTeapot,
		Duration: 2 * time.Second,
	})
	sum.Finalize(2 * time.Second)
	sum.Fprint(out)

	return nil
}
