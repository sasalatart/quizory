package otel

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

const logExportInterval = 10 * time.Second

func newLoggerProvider() (*log.LoggerProvider, error) {
	ctx := context.Background()

	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating OTLP logs exporter")
	}

	lp := log.NewLoggerProvider(
		log.WithProcessor(
			log.NewBatchProcessor(logExporter, log.WithExportInterval(logExportInterval)),
		),
	)

	global.SetLoggerProvider(lp)

	return lp, nil
}
