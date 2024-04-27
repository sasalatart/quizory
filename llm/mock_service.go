package llm

import "context"

type MockService struct {
	ChatCompletionFn func(systemContent, userContent string) (string, error)
}

func newMockService() *MockService {
	return &MockService{}
}

func (s *MockService) ChatCompletion(
	ctx context.Context,
	systemContent, userContent string,
) (string, error) {
	return s.ChatCompletionFn(systemContent, userContent)
}

var _ ChatCompletioner = &MockService{}
