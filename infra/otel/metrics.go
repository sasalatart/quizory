package otel

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

const (
	metricsExportInterval   = 10 * time.Second
	minReadMemStatsInterval = 5 * time.Second
)

func newMeterProvider(ctx context.Context, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating OTLP metrics exporter")
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithInterval(metricsExportInterval),
			),
		),
	)

	// Auto-instrument runtime metrics (e.g., memory, CPU, GC stats)
	if err := runtime.Start(
		runtime.WithMeterProvider(provider),
		runtime.WithMinimumReadMemStatsInterval(minReadMemStatsInterval),
	); err != nil {
		return nil, errors.Wrap(err, "starting runtime metrics")
	}

	return provider, nil
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
	return provider.Meter("quizory")
}
