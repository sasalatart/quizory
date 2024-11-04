package generator

import (
	"github.com/sasalatart/quizory/generator/internal/metrics"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"generator",
	fx.Provide(fx.Private, metrics.NewService),
	fx.Provide(NewService),
)
