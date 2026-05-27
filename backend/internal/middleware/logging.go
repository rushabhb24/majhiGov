package middleware

import (
	"log"
	"net/http"
	"time"
)

// statusResponseWriter wraps http.ResponseWriter to capture the status code
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware logs HTTP method, path, status code, and request duration
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		sw := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(sw, r)

		duration := time.Since(start)
		log.Printf("[%s] %s %d %v", r.Method, r.URL.Path, sw.statusCode, duration)
	})
}
