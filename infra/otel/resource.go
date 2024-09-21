package otel

import (
	"context"
	"runtime"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// newResource creates a new OTEL resource with the given service name and version, such that it may
// be associated with signals (e.g. traces, metrics, logs) emitted by the service.
func newResource(
	ctx context.Context,
	serviceName, serviceVersion string,
) (*resource.Resource, error) {
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String(serviceVersion),
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
