package main

import (
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

type parseFunc func(string) error

func (f *flags) parse() (err error) {
	parsers := map[string]parseFunc{
		"url": f.urlVar(&f.url),
		"n":   f.intVar(&f.n),
		"c":   f.intVar(&f.c),
	}
	for _, arg := range os.Args[1:] {
		n, v, ok := strings.Cut(arg, "=")
		if !ok {
			continue
		}
		parse, ok := parsers[strings.TrimPrefix(n, "-")]
		if !ok {
			continue
		}
		if err := parse(v); err != nil {
			err = fmt.Errorf("invalid value %q for flag %s: %w", v, n, err)
			break
		}
	}
	return err
}

func (f *flags) urlVar(p *string) parseFunc {
	return func(s string) error {
		_, err := url.Parse(s)
		*p = s
		return err
	}
}

func (f *flags) intVar(p *int) parseFunc {
	return func(s string) (err error) {
		*p, err = strconv.Atoi(s)
		return err
	}
}
