package llm

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"llm",
	fx.Provide(fx.Annotate(NewService, fx.As(new(ChatCompletioner)))),
)

var TestModule = fx.Module(
	"llm-test",
	fx.Provide(fx.Annotate(newMockService, fx.As(new(ChatCompletioner)))),
)
