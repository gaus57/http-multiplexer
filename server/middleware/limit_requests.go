package middleware

import (
	"net/http"
	"sync"
)

// LimitRequests limits the number of simultaneously executed requests
func LimitRequests(limit int64, next http.HandlerFunc) http.HandlerFunc {
	l := &limiter{
		limit: limit,
	}

	return func(res http.ResponseWriter, req *http.Request) {
		if !l.Increase() {
			res.WriteHeader(http.StatusTooManyRequests)
			return
		}
		defer l.Decrease()

		next(res, req)
	}
}

type limiter struct {
	limit int64
	count int64
	mux   sync.Mutex
}

func (l *limiter) Increase() bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.count >= l.limit {
		return false
	}

	l.count++

	return true
}

func (l *limiter) Decrease() {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.count--
}
