package otel

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"otel",
	fx.Provide(newLoggerProvider),
	fx.Provide(newMetricsProvider),
	fx.Invoke(loggerLC),
	fx.Invoke(metricsLC),
)

func loggerLC(lc fx.Lifecycle, lp *log.LoggerProvider) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.SetDefault(otelslog.NewLogger("quizory"))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return lp.Shutdown(ctx)
		},
	})
}

func metricsLC(lc fx.Lifecycle, lp *metric.MeterProvider) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return lp.Shutdown(ctx)
		},
	})
}
