package llm

import "context"

// ChatCompletioner is an interface that defines a method to complete an LLM chat message.
type ChatCompletioner interface {
	ChatCompletion(ctx context.Context, systemContent, userContent string) (string, error)
}
