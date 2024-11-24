package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/infra/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

// WithMonitoring instruments HTTP handlers with OpenTelemetry logs, metrics and traces.
//
// The following metrics are emitted:
//   - http_requests_total: A counter that tracks the total number of HTTP requests.
//   - http_request_duration_ms: A histogram that tracks the duration of HTTP requests in milliseconds.
func WithMonitoring(meter otel.Meter) (func(http.Handler) http.Handler, error) {
	requestCounter, err := meter.Int64Counter("http_requests_total")
	if err != nil {
		return nil, errors.Wrap(err, "creating http_requests_total counter")
	}

	requestHistogram, err := meter.Int64Histogram("http_request_duration_ms")
	if err != nil {
		return nil, errors.Wrap(err, "creating http_request_duration_ms histogram")
	}

	return func(next http.Handler) http.Handler {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			method := slog.String("method", r.Method)
			path := slog.String("path", r.URL.Path)

			slog.Info("Handling request", slog.Group("http", method, path))
			defer func() {
				status := slog.Int("status", rw.status)
				duration := slog.String("duration", time.Since(start).String())

				slog.Info("Completed request", slog.Group("http", method, path, status, duration))

				attributes := metric.WithAttributes(
					attribute.String("method", r.Method),
					attribute.String("path", r.Pattern),
					attribute.Int("status", rw.status),
				)
				requestCounter.Add(ctx, 1, attributes)
				requestHistogram.Record(ctx, int64(time.Since(start).Milliseconds()), attributes)
			}()

			next.ServeHTTP(rw, r)
		})
		return otelhttp.NewHandler(
			handler,
			"http",
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
		)
	}, nil
}

// spanNameFormatter generates a span name based on the HTTP method and request path.
func spanNameFormatter(_ string, r *http.Request) string {
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
}
