package question

import (
	"github.com/sasalatart/quizory/domain/question/internal/metrics"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"question",
	fx.Provide(fx.Private, NewRepository),
	fx.Provide(fx.Private, metrics.NewService),
	fx.Provide(NewService),
)
