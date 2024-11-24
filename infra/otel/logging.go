package otel

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/pkg/errors"
	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
)

const logExportInterval = 5 * time.Second

func initLoggerProvider(ctx context.Context, res *resource.Resource) (*sdklog.LoggerProvider, error) {
	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "creating OTLP logs exporter")
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(
			sdklog.NewBatchProcessor(
				logExporter,
				sdklog.WithExportInterval(logExportInterval),
			),
		),
	)

	global.SetLoggerProvider(lp)

	return lp, nil
}

func newDefaultLogger() *slog.Logger {
	multiHandler := slogmulti.Fanout(
		otelslog.NewHandler(getServiceName()),
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	return slog.New(multiHandler)
}
