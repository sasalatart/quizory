package generator

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"generator",
	fx.Provide(NewService),
)
