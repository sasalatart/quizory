package llm

import (
	"context"
	"time"

	"github.com/sasalatart/quizory/config"
	"github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	openaiClient *openai.Client
	tracer       trace.Tracer
}

func NewService(cfg config.LLMConfig) *Service {
	return &Service{
		openaiClient: openai.NewClient(cfg.OpenAIKey),
		tracer:       otel.Tracer("llm.Service"),
	}
}

func (s *Service) ChatCompletion(
	ctx context.Context,
	systemContent, userContent string,
) (string, error) {
	ctx, span := s.tracer.Start(ctx, "ChatCompletion")
	defer span.End()

	var seed int = time.Now().Nanosecond()
	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemContent,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userContent,
				},
			},
			Seed:      &seed,
			MaxTokens: 4096,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

var _ ChatCompletioner = &Service{}
