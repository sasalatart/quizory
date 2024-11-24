package otel

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	metricsExportInterval   = 5 * time.Second
	minReadMemStatsInterval = 5 * time.Second
)

func initMeterProvider(ctx context.Context, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating OTLP metrics exporter")
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(metricsExportInterval),
			),
		),
	)

	otel.SetMeterProvider(mp)
	if err := autoInstrumentRuntime(mp); err != nil {
		return nil, errors.Wrap(err, "auto-instrumenting runtime")
	}

	return mp, nil
}

func autoInstrumentRuntime(mp *sdkmetric.MeterProvider) error {
	err := runtime.Start(
		runtime.WithMeterProvider(mp),
		runtime.WithMinimumReadMemStatsInterval(minReadMemStatsInterval),
	)
	if err != nil {
		return errors.Wrap(err, "starting runtime metrics")
	}
	return nil
}

type Meter interface {
	Int64Counter(
		name string,
		opts ...metric.Int64CounterOption,
	) (metric.Int64Counter, error)

	Int64Histogram(
		name string,
		opts ...metric.Int64HistogramOption,
	) (metric.Int64Histogram, error)
}

func newMeter(provider *sdkmetric.MeterProvider) Meter {
	return provider.Meter(getServiceName())
}
