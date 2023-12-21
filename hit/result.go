package hit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Result struct {
	RPS      float64       // requests per second
	Requests int           //number of requesys made
	Errors   int           // number of errors occured
	Bytes    int64         //number of bytesdownloaded
	Duration time.Duration //single or all requests' duration
	Fastest  time.Duration //fastest request result duration among other
	Slowest  time.Duration //slowest request result duration among other
	Status   int           // request HTTP status code
	Error    error         //error is not nill if the request is failed
}

func (r *Result) Merge(o *Result) {
	r.Requests++
	r.Bytes += o.Bytes

	if r.Fastest == 0 || o.Duration < r.Fastest {
		r.Fastest = o.Duration
	}
	if o.Duration > r.Slowest {
		r.Slowest = o.Duration
	}

	switch {
	case o.Error != nil:
		fallthrough
	case o.Status >= http.StatusBadRequest:
		r.Errors++
	}
}

func (r *Result) Finalize(total time.Duration) *Result {
	r.Duration = total
	r.RPS = float64(r.Requests) / total.Seconds()
	return r
}

func (r *Result) Fprint(out io.Writer) {
	p := func(format string, args ...any) {
		fmt.Fprintf(out, format, args...)
	}
	p("\nSummary:\n")
	p("\tSuccess    : %.0f%%\n", r.success())
	p("\tRPS        : %.1f\n", r.RPS)
	p("\tRequests   : %d\n", r.Requests)
	p("\tErrors     : %d\n", r.Errors)
	p("\tBytes      : %d\n", r.Bytes)
	p("\tDuration   : %s\n", round(r.Duration))
	if r.Requests > 1 {
		p("\tFastet     : %s\n", round(r.Fastest))
		p("\tSlowest    : %s\n", round(r.Slowest))
	}
}

func (r *Result) success() float64 {
	rr, e := float64(r.Requests), float64(r.Errors)
	return (rr - e) / rr * 100
}

func (r *Result) String() string {
	var s strings.Builder
	r.Fprint(&s)
	return s.String()
}

func round(t time.Duration) time.Duration {
	return t.Round(time.Microsecond)
}
