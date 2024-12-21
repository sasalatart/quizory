package llm

import (
	"context"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sasalatart/quizory/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	openaiClient *openai.Client
	tracer       trace.Tracer
}

func NewService(cfg config.LLMConfig) *Service {
	return &Service{
		openaiClient: openai.NewClient(option.WithAPIKey(cfg.OpenAIKey)),
		tracer:       otel.Tracer("llm.Service"),
	}
}

func (s *Service) Chat(
	ctx context.Context,
	systemContent, userContent string,
) (string, error) {
	ctx, span := s.tracer.Start(ctx, "ChatCompletion")
	defer span.End()

	resp, err := s.openaiClient.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Model: openai.F(openai.ChatModelO1Preview),
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.UserMessage(systemContent), // This model does not support system messages.
				openai.UserMessage(userContent),
			}),
			Seed: openai.Int(time.Now().UnixNano()),
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

var _ Chater = &Service{}
