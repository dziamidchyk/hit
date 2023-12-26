package main

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type flags struct {
	url, m    string
	n, c, rps int
	t         time.Duration
	h         []string
}

const usageText = `
Usage:
  hit [options] url
Options:`

func (f *flags) parse(s *flag.FlagSet, args []string) error {
	s.Usage = func() {
		fmt.Fprintln(s.Output(), usageText[1:])
		s.PrintDefaults()
	}
	s.Var(toNumber(&f.n), "n", "Number of requests to make")
	s.Var(toNumber(&f.c), "c", "Concurrency level")
	s.Var(toNumber(&f.rps), "rps", "Throttle requests per second")
	s.Var(toMethod(&f.m), "m", "Method")
	s.Var(toHeaders(&f.h), "H", "Headers")
	s.DurationVar(&f.t, "t", f.t, "Timeout")
	if err := s.Parse(args); err != nil {
		return err
	}
	f.url = s.Arg(0)

	if err := f.validate(); err != nil {
		fmt.Fprintln(s.Output(), err)
		s.Usage()
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
		return fmt.Errorf("url: %w", err)
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

type method string

func toMethod(s *string) *method {
	return (*method)(s)
}

func (m *method) Set(s string) error {
	switch s {
	case "GET", "POST", "PUT":
		*m = method(s)
		return nil
	default:
		return fmt.Errorf("incorrect method: %s", s)
	}
}

func (m *method) String() string {
	return string(*m)
}

type headers []string

func toHeaders(s *[]string) *headers {
	return (*headers)(s)
}

func (h *headers) Set(s string) error {
	*h = append(*h, s)
	return nil
}

func (h *headers) String() string {
	return strings.Join(*h, ", ")
}
