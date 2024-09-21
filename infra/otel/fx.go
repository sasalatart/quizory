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
	LoggerProvider *log.LoggerProvider
	MeterProvider  *metric.MeterProvider
}

func newProvider() (*Provider, error) {
	ctx := context.Background()

	res, err := newResource(ctx, getServiceName(), getServiceVersion())
	if err != nil {
		return nil, errors.Wrap(err, "creating resource")
	}

	lp, err := newLoggerProvider(ctx, res)
	if err != nil {
		return nil, errors.Wrap(err, "creating logger provider")
	}

	mp, err := newMeterProvider(ctx, res)
	if err != nil {
		return nil, errors.Wrap(err, "creating meter provider")
	}

	return &Provider{
		LoggerProvider: lp,
		MeterProvider:  mp,
	}, nil
}

func providerLC(lc fx.Lifecycle, op *Provider) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slog.SetDefault(otelslog.NewLogger(getServiceName()))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if err := op.LoggerProvider.Shutdown(ctx); err != nil {
				return errors.Wrap(err, "shutting down logger provider")
			}
			if err := op.MeterProvider.Shutdown(ctx); err != nil {
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
