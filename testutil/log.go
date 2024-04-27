package testutil

import (
	"log/slog"
	"os"
)

func init() {
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	})
	slog.SetDefault(slog.New(slogHandler))
}
