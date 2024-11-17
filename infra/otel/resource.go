package otel

import (
	"context"
	"os"
	"runtime"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// newResource creates a new OTEL resource such that it may be associated with signals (e.g. traces,
// metrics, logs) emitted by the service.
func newResource(ctx context.Context) (*resource.Resource, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(getServiceName()),
			semconv.ServiceVersionKey.String(getServiceVersion()),
			attribute.String("go.version", runtime.Version()),
			semconv.OSNameKey.String(runtime.GOOS),
			semconv.HostArchKey.String(runtime.GOARCH),
		),
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating resource")
	}
	return res, nil
}

func getServiceName() string {
	return os.Getenv("OTEL_SERVICE_NAME")
}

func getServiceVersion() string {
	return "0.1.0" // TODO: make this dynamic (e.g. via env var or similar)
}
