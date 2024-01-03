package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
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
		fmt.Fprintln(os.Stderr, "error occurred:", err)
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
	if f.rps > 0 {
		fmt.Fprintf(out, "(RPS: %d)\n", f.rps)
	}

	request, err := http.NewRequest(http.MethodGet, f.url, http.NoBody)
	if err != nil {
		return err
	}
	c := &hit.Client{
		C:   f.c,
		RPS: f.rps,
	}

	ctx, cancel := context.WithTimeout(context.Background(), f.t)
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	defer stop()
	sum := c.Do(ctx, request, f.n)
	sum.Fprint(out)

	if err := ctx.Err(); errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("timed out in %s", f.t)
	}
	return nil
}
