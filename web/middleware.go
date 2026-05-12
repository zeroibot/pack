package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zeroibot/pack/clock"
	"github.com/zeroibot/pack/sys"
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

// NewCORSMiddleware creates a new CORS middleware
func NewCORSMiddleware(appEnv sys.Env, allowedOrigins []string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if appEnv == sys.EnvDev {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			} else if len(allowedOrigins) > 0 {
				w.Header().Set("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ", "))
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, User-Agent, Cache-Control")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
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
