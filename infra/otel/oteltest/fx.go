package oteltest

import (
	"github.com/sasalatart/quizory/infra/otel"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"otel-test",
	fx.Provide(fx.Annotate(newMockMeter, fx.As(new(otel.Meter)))),
)
