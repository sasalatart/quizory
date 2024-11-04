package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter is a wrapper around an http.ResponseWriter that captures the status code written
// to it. This is necessary because the built-in http.ResponseWriter does not expose it.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// WithLogger decorates HTTP requests such that two logs are emitted for each request: one when the
// request is just received to record details such as the HTTP method and path, and another when the
// request is completed to record its duration and resulting HTTP status.
func WithLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		method := slog.String("method", r.Method)
		path := slog.String("path", r.URL.Path)
		slog.Info(
			"Handling request",
			slog.Group("http", method, path),
		)

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		status := slog.Int("status", rw.status)
		duration := slog.String("duration", time.Since(start).String())
		slog.Info(
			"Completed request",
			slog.Group("http", method, path, status, duration),
		)
	})
}
