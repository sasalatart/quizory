package otel

import (
	"context"
	"log/slog"
	"os"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"otel",
	fx.Provide(newProvider),
	fx.Invoke(providerLC),
	fx.Provide(fx.Annotate(newMeter, fx.As(new(Meter)))),
)

type Provider struct {
	fx.Out

	LoggerProvider *log.LoggerProvider
	MeterProvider  *metric.MeterProvider
}

func newProvider() (Provider, error) {
	ctx := context.Background()
	provider := Provider{}

	res, err := newResource(ctx, getServiceName(), getServiceVersion())
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
			slog.SetDefault(otelslog.NewLogger(getServiceName()))
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

func getServiceName() string {
	return os.Getenv("OTEL_SERVICE_NAME")
}

func getServiceVersion() string {
	return "0.1.0" // TODO: make this dynamic (e.g. via env var or similar)
}
