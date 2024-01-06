package hit_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/dziamidchyk/hit/hit"
)

func ExampleDo() {
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	res, err := hit.Do(context.Background(), server.URL, 10, hit.WithConcurrency(2), hit.WithTimeout(time.Second*30))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
