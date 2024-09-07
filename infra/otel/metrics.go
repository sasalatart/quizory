package otel

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	metricsExportInterval   = 10 * time.Second
	minReadMemStatsInterval = 4 * time.Second
)

func newMetricsProvider() (*metric.MeterProvider, error) {
	ctx := context.Background()

	exporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating OTLP metrics exporter")
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(
			metric.NewPeriodicReader(exporter, metric.WithInterval(metricsExportInterval)),
		),
	)

	otel.SetMeterProvider(provider)

	// Auto-instrument runtime metrics (e.g., memory, CPU, GC stats)
	if err := runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(minReadMemStatsInterval),
	); err != nil {
		return nil, errors.Wrap(err, "starting runtime metrics")
	}

	return provider, nil
}

type Meter interface {
	Int64Counter(
		name string,
		opts ...otelmetric.Int64CounterOption,
	) (otelmetric.Int64Counter, error)

	Int64Histogram(
		name string,
		opts ...otelmetric.Int64HistogramOption,
	) (otelmetric.Int64Histogram, error)
}

func newMeter(provider *metric.MeterProvider) Meter {
	return provider.Meter("quizory")
}
