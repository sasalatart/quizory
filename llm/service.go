package llm

import (
	"context"
	"time"

	"github.com/sasalatart/quizory/config"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	openaiClient *openai.Client
}

func NewService(cfg config.LLMConfig) *Service {
	return &Service{
		openaiClient: openai.NewClient(cfg.OpenAIKey),
	}
}

func (s *Service) ChatCompletion(
	ctx context.Context,
	systemContent, userContent string,
) (string, error) {
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
