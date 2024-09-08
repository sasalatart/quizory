package llmtest

import (
	"github.com/sasalatart/quizory/llm"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"llm-test",
	fx.Provide(fx.Annotate(newMockService, fx.As(new(llm.ChatCompletioner)))),
)
