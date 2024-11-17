package otel

import (
	"context"
	"log/slog"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
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
}

func newProvider() (Provider, error) {
	ctx := context.Background()
	provider := Provider{}

	res, err := newResource(ctx)
	if err != nil {
		return provider, errors.Wrap(err, "creating resource")
	}

	lp, err := newLoggerProvider(ctx, res)
	if err != nil {
		return provider, errors.Wrap(err, "creating logger provider")
	}
	provider.LoggerProvider = lp

	mp, err := newMeterProvider(ctx, res)
	if err != nil {
		return provider, errors.Wrap(err, "creating meter provider")
	}
	provider.MeterProvider = mp

	return provider, nil
}

func providerLC(lc fx.Lifecycle, lp *log.LoggerProvider, mp *metric.MeterProvider) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			global.SetLoggerProvider(lp)
			slog.SetDefault(newDefaultLogger())

			otel.SetMeterProvider(mp)
			if err := autoInstrumentRuntime(mp); err != nil {
				return errors.Wrap(err, "auto-instrumenting runtime")
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := lp.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "shutting down logger provider")
			}
			if err := mp.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "shutting down meter provider")
			}
			return nil
		},
	})
}
