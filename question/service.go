package question

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/sasalatart.com/quizory/config"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	cfg          config.Config
	repo         *Repository
	openaiClient *openai.Client
}

// NewService creates a new instance of question.Service.
func NewService(cfg config.Config, repo *Repository) *Service {
	return &Service{cfg: cfg, repo: repo, openaiClient: openai.NewClient(cfg.OpenAIKey)}
}

// StartGeneration generates questions about random topics at a given frequency.
func (s Service) StartGeneration(ctx context.Context) {
	freq := s.cfg.QuestionGeneration.Frequency
	batchSize := s.cfg.QuestionGeneration.BatchSize
	slog.Info("Starting generation loop", slog.Duration("freq", freq))

	ticker := time.NewTicker(freq)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			topic := randomTopic()
			slog.Info(
				"Generating questions",
				slog.String("topic", topic.String()),
				slog.Int("amount", batchSize),
			)
			if err := s.generateQuestionSet(ctx, topic, batchSize); err != nil {
				slog.Error("Error generating question set", err)
				return
			}
		}
	}
}
