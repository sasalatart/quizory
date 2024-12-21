package llmtest

import (
	"context"

	"github.com/sasalatart/quizory/llm"
)

type MockService struct {
	ChatFn func(systemContent, userContent string) (string, error)
}

func newMockService() *MockService {
	return &MockService{}
}

func (s *MockService) Chat(ctx context.Context, systemContent, userContent string) (string, error) {
	return s.ChatFn(systemContent, userContent)
}

var _ llm.Chater = &MockService{}
