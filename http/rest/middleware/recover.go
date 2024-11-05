package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

// WithRecover recovers from panics that occur during the execution of the next http.Handler and
// handles it gracefully by logging the error and responding with a 500 Internal Server Error.
func WithRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error(
					"Recovered from panic",
					slog.Any("error", err),
					slog.String("stack", string(debug.Stack())),
				)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
