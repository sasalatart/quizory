package otel

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/sdk/log"
)

func newLoggerProvider() (*log.LoggerProvider, error) {
	ctx := context.Background()

	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, err
	}

	lp := log.NewLoggerProvider(log.WithProcessor(log.NewBatchProcessor(logExporter)))
	return lp, nil
}
