package otel

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"otel",
	fx.Provide(newProvider),
	fx.Provide(fx.Annotate(newMeter, fx.As(new(Meter)))),
	fx.Invoke(providerLC),
)

type Provider struct {
	fx.Out

	LoggerProvider *log.LoggerProvider
	MeterProvider  *metric.MeterProvider
	TracerProvider *trace.TracerProvider
}

func newProvider() (Provider, error) {
	ctx := context.Background()
	provider := Provider{}

	res, err := newResource(ctx)
	if err != nil {
		return provider, errors.Wrap(err, "creating resource")
	}

	lp, err := initLoggerProvider(ctx, res)
	if err != nil {
		return provider, errors.Wrap(err, "creating logger provider")
	}
	provider.LoggerProvider = lp

	mp, err := initMeterProvider(ctx, res)
	if err != nil {
		return provider, errors.Wrap(err, "creating meter provider")
	}
	provider.MeterProvider = mp

	tp, err := initTracerProvider(ctx, res)
	if err != nil {
		return provider, errors.Wrap(err, "creating tracer provider")
	}
	provider.TracerProvider = tp

	return provider, nil
}

func providerLC(
	lc fx.Lifecycle,
	lp *log.LoggerProvider,
	mp *metric.MeterProvider,
	tp *trace.TracerProvider,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.SetDefault(newDefaultLogger())
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := lp.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "shutting down logger provider")
			}
			if err := mp.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "shutting down meter provider")
			}
			if err := tp.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "shutting down tracer provider")
			}
			return nil
		},
	})
}
