package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingRepsonseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingRepsonseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingRepsonseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r) // Call the next handler

		duration_ms := float64(time.Since(start).Microseconds()) / 1000
		// Log the request details
		slog.Info("request completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lrw.statusCode,
			"duration_ms", duration_ms,
			"remote", r.RemoteAddr,
		)
	})
}
