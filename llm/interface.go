package llm

import "context"

// Chater is an interface that defines a method to complete an LLM chat message.
type Chater interface {
	Chat(ctx context.Context, systemContent, userContent string) (string, error)
}
