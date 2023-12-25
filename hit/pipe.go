package hit

import (
	"net/http"
	"sync"
	"time"
)

func Produce(out chan<- *http.Request, n int, fn func() *http.Request) {
	for ; n > 0; n-- {
		out <- fn()
	}
}

func produce(n int, fn func() *http.Request) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Produce(out, n, fn)
	}()
	return out
}

func Throttle(in <-chan *http.Request, out chan<- *http.Request, delay time.Duration) {
	t := time.NewTicker(delay)
	defer t.Stop()

	for r := range in {
		<-t.C
		out <- r
	}
}

func throttle(in <-chan *http.Request, delay time.Duration) <-chan *http.Request {
	out := make(chan *http.Request)
	go func() {
		defer close(out)
		Throttle(in, out, delay)
	}()
	return out
}

func Split(in <-chan *http.Request, out chan<- *Result, c int, fn SendFunc) {
	send := func() {
		for r := range in {
			out <- fn(r)
		}
	}

	var wg sync.WaitGroup
	wg.Add(c)
	for ; c > 0; c-- {
		go func() {
			defer wg.Done()
			send()
		}()
	}
	wg.Wait()
}

func split(in <-chan *http.Request, c int, fn SendFunc) <-chan *Result {
	out := make(chan *Result)
	go func() {
		defer close(out)
		Split(in, out, c, fn)
	}()
	return out
}
