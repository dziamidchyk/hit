package main

import (
	"flag"
)

type flags struct {
	url  string
	n, c int
}

func (f *flags) parse() error {
	flag.StringVar(&f.url, "url", "", "HTTP server `URL` to make requests (required)")
	flag.IntVar(&f.n, "n", f.n, "Number of requests to make")
	flag.IntVar(&f.c, "c", f.c, "Concurrency level")
	flag.Parse()
	return nil
}
