package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gaus57/http-multiplexer/server/middleware"
)

// TestLimitRequests checks the limit on the number of requests
func TestLimitRequests(t *testing.T) {
	wait := make(chan struct{})
	handler := middleware.LimitRequests(2, func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusOK)
		<-wait
	})

	r1 := httptest.NewRecorder()
	r2 := httptest.NewRecorder()
	r3 := httptest.NewRecorder()
	r4 := httptest.NewRecorder()

	go handler.ServeHTTP(r1, nil)
	go handler.ServeHTTP(r2, nil)

	time.Sleep(1 * time.Millisecond)
	go handler.ServeHTTP(r3, nil)

	time.Sleep(1 * time.Millisecond)
	close(wait)

	wait = make(chan struct{})
	go handler.ServeHTTP(r4, nil)

	close(wait)

	if r1.Code != http.StatusOK {
		t.Errorf("unexpected response status for request 1: %d", r1.Code)
	}
	if r2.Code != http.StatusOK {
		t.Errorf("unexpected response status for request 2: %d", r2.Code)
	}
	if r3.Code != http.StatusTooManyRequests {
		t.Errorf("unexpected response status for request 3: %d", r3.Code)
	}
	if r4.Code != http.StatusOK {
		t.Errorf("unexpected response status for request 4: %d", r4.Code)
	}
}
