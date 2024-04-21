package question

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/sashabaranov/go-openai"
)

type Service struct {
	repo         *Repository
	openaiClient *openai.Client
}

// NewService creates a new instance of question.Service.
func NewService(repo *Repository, openaiClient *openai.Client) Service {
	return Service{repo: repo, openaiClient: openaiClient}
}

// StartGeneration generates questions about random topics at a given frequency.
func (s Service) StartGeneration(
	ctx context.Context,
	freq time.Duration,
	amountPerBatch int,
	cancel <-chan struct{},
) {
	slog.Info("Starting generation loop", slog.Duration("freq", freq))
	ticker := time.NewTicker(freq)
	for {
		select {
		case <-cancel:
			return
		case <-ticker.C:
			topic := randomTopic()
			slog.Info(
				"Generating questions",
				slog.String("topic", topic.String()),
				slog.Int("amount", amountPerBatch),
			)
			if err := s.generateQuestionSet(ctx, topic, amountPerBatch); err != nil {
				slog.Error("Error generating question set", err)
				return
			}
		}
	}
}
