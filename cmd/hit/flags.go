package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type flags struct {
	url  string
	n, c int
}

func (f *flags) parse() error {
	flag.StringVar(&f.url, "url", "", "HTTP server `URL` to make requests (required)")
	flag.Var(toNumber(&f.n), "n", "Number of requests to make")
	flag.Var(toNumber(&f.c), "c", "Concurrency level")
	flag.Parse()
	if err := f.validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		return err
	}
	return nil
}

func (f *flags) validate() error {
	if strings.TrimSpace(f.url) == "" {
		return errors.New("-url: required")
	}
	if f.c > f.n {
		return fmt.Errorf("-c=%d: should be less than or equal to -n=%d", f.c, f.n)
	}
	if err := f.validateURL(f.url); err != nil {
		return fmt.Errorf("invalid value %q for flag -url: %w", f.url, err)
	}
	return nil
}

func (f *flags) validateURL(s string) error {
	u, err := url.Parse(s)
	switch {
	case strings.TrimSpace(s) == "":
		err = errors.New("required")
	case err != nil:
		err = errors.New("parse error")
	case u.Scheme != "http":
		err = errors.New("only supported scheme is http")
	case u.Host == "":
		err = errors.New("missing host")
	}
	return err
}

type number int

func toNumber(p *int) *number {
	return (*number)(p)
}

func (n *number) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	switch {
	case err != nil:
		err = errors.New("parse error")
	case v <= 0:
		err = errors.New("should be positive")
	}
	*n = number(v)
	return err
}

func (n *number) String() string {
	return strconv.Itoa(int(*n))
}
