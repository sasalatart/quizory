package infra

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
)

// WaitFor retries a check function up to maxAttempts times with an exponential backoff, starting at
// the given timeout. It returns an error if the check does not return true after all attempts.
func WaitFor(check func() bool, maxAttempts int, timeout time.Duration) error {
	for i := 0; i < maxAttempts; i++ {
		if check() {
			return nil
		}
		if i == maxAttempts-1 {
			return errors.New("max attempts reached")
		}
		time.Sleep(timeout)
		timeout *= 2
	}
	return nil
}

func RunFXApp(ctx context.Context, app *fx.App) {
	go func() {
		if err := app.Start(ctx); err != nil {
			slog.Error("Error starting app", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	if err := app.Stop(ctx); err != nil {
		slog.Error("Error stopping app", slog.Any("error", err))
		os.Exit(1)
	}
}
