package hit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClientDo(t *testing.T) {
	t.Parallel()

	const wantHits, wantErrors = 10, 0
	var gotHits atomic.Int64

	handler := func(_ http.ResponseWriter, _ *http.Request) {
		gotHits.Add(1)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()
	request, err := http.NewRequest(http.MethodGet, server.URL, http.NoBody)
	if err != nil {
		t.Fatalf("NewRequest err=%q; want nil", err)
	}

	c := &Client{
		C: 1,
	}
	sum := c.Do(context.Background(), request, wantHits)
	if got := gotHits.Load(); got != wantHits {
		t.Errorf("hits=%d; want %d", got, wantHits)
	}
	if got := sum.Requests; got != wantHits {
		t.Errorf("Requests=%d; want %d", got, wantHits)
	}
	if got := sum.Errors; got != wantErrors {
		t.Errorf("Errors=%d; want %d", got, wantErrors)
	}
}
