package hit

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"
)

type Client struct {
	C       int           // concurrency level
	RPS     int           // throttles the requests per second
	Timeout time.Duration // Timeout per request
}

func (c *Client) Do(ctx context.Context, r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(ctx, r, n)
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(ctx context.Context, r *http.Request, n int) *Result {
	p := produce(ctx, n, func() *http.Request {
		return r.Clone(ctx)
	})
	if c.RPS > 0 {
		p = throttle(p, time.Second/time.Duration(c.RPS*c.concurrency()))
	}
	var (
		sum    Result
		client = c.client()
	)
	defer client.CloseIdleConnections()
	for result := range split(p, c.concurrency(), c.send(client)) {
		sum.Merge(result)
	}
	return &sum
}

type SendFunc func(*http.Request) *Result

func Send(c *http.Client, r *http.Request) *Result {
	t := time.Now()

	var (
		code  int
		bytes int64
	)
	response, err := c.Do(r)
	if err == nil {
		code = response.StatusCode
		bytes, err = io.Copy(io.Discard, response.Body)
		_ = response.Body.Close()
	}

	return &Result{
		Duration: time.Since(t),
		Bytes:    bytes,
		Status:   code,
		Error:    err,
	}
}

func (c *Client) send(client *http.Client) SendFunc {
	return func(r *http.Request) *Result {
		return Send(client, r)
	}
}

func (c *Client) client() *http.Client {
	return &http.Client{
		Timeout: c.Timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: c.concurrency(),
		},
	}
}

func (c *Client) concurrency() int {
	if c.C > 0 {
		return c.C
	}
	return runtime.NumCPU()
}

type OptionFunc func(*Client)

func WithConcurrency(n int) OptionFunc {
	return func(c *Client) {
		c.C = n
	}
}

func WithTimeout(d time.Duration) OptionFunc {
	return func(c *Client) {
		c.Timeout = d
	}
}

func Do(ctx context.Context, url string, n int, options ...OptionFunc) (*Result, error) {
	r, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("new http request: %w", err)
	}
	var c Client
	for _, fn := range options {
		fn(&c)
	}
	return c.Do(ctx, r, n), nil
}
