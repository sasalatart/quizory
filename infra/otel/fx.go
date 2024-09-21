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
	fx.Provide(newOTELProvider),
	fx.Invoke(otelProviderLC),
	fx.Provide(fx.Annotate(newMeter, fx.As(new(Meter)))),
)

type OTELProvider struct {
	fx.Out

	LoggerProvider *log.LoggerProvider
	MeterProvider  *metric.MeterProvider
}

func newOTELProvider() (OTELProvider, error) {
	ctx := context.Background()
	provider := OTELProvider{}

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

func otelProviderLC(lc fx.Lifecycle, op *OTELProvider) {
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
