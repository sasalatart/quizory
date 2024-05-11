package question

import "go.uber.org/fx"

var Module = fx.Module(
	"question",
	fx.Provide(fx.Private, NewRepository),
	fx.Provide(NewService),
)
