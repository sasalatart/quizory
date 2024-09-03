package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"otel",
	fx.Provide(newLoggerProvider),
	fx.Invoke(loggerLC),
)

func loggerLC(lc fx.Lifecycle, lp *log.LoggerProvider) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			global.SetLoggerProvider(lp)
			slog.SetDefault(otelslog.NewLogger("quizory"))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return lp.Shutdown(ctx)
		},
	})
}
