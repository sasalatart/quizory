package otel

import (
	"context"

	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/embedded"
)

type MockMeter struct{}

func newMockMeter() *MockMeter {
	return &MockMeter{}
}

var _ Meter = MockMeter{}

func (m MockMeter) Int64Counter(
	name string,
	opts ...otelmetric.Int64CounterOption,
) (otelmetric.Int64Counter, error) {
	return mockAdder{}, nil
}

func (m MockMeter) Int64Histogram(
	name string,
	opts ...otelmetric.Int64HistogramOption,
) (otelmetric.Int64Histogram, error) {
	return mockRecorder{}, nil
}

type mockAdder struct {
	embedded.Int64Counter
}

func (mockAdder) Add(ctx context.Context, incr int64, options ...otelmetric.AddOption) {}

type mockRecorder struct {
	embedded.Int64Histogram
}

func (mockRecorder) Record(ctx context.Context, incr int64, options ...otelmetric.RecordOption) {}
