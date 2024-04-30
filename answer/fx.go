package answer

import "go.uber.org/fx"

var Module = fx.Module(
	"answer",
	fx.Provide(fx.Private, NewRepository),
	fx.Provide(NewService),
)
