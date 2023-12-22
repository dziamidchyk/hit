package hit

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
}

func (c *Client) Do(r *http.Request, n int) *Result {
	t := time.Now()
	sum := c.do(r, n)
	return sum.Finalize(time.Since(t))
}

func (c *Client) do(r *http.Request, n int) *Result {
	var sum Result
	for ; n > 0; n-- {
		sum.Merge(Send(r))
	}
	return &sum
}

type SendFunc func(*http.Request) *Result

func Send(r *http.Request) *Result {
	t := time.Now()

	fmt.Printf("request: %s\n", r.URL)
	time.Sleep(100 * time.Millisecond)

	return &Result{
		Duration: time.Since(t),
		Bytes:    10,
		Status:   http.StatusOK,
	}
}
