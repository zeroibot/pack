package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zeroibot/pack/clock"
)

type Middleware = func(next http.Handler) http.Handler

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggingWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware prints the request timestamp, method, path, and duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		now := clock.StandardFormat(start)
		lw := new(loggingWriter{w, http.StatusOK})
		next.ServeHTTP(lw, r)
		duration := time.Since(start).Round(time.Millisecond)
		fmt.Printf("[%s] (%d) %s %s - %v\n", now, lw.statusCode, r.Method, r.URL.Path, duration)
	})
}

// StackMiddlewares combines multiple middlewares into a single middleware
func StackMiddlewares(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
