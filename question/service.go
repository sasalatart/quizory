package question

import (
	"context"
	_ "embed"
	"log/slog"
	"time"

	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/llm"
	"github.com/sasalatart.com/quizory/question/enums"
	"github.com/sasalatart.com/quizory/question/internal/ai"
)

// Service represents the service that manages questions.
type Service struct {
	repo       *Repository
	llmService llm.ChatCompletioner
}

// NewService creates a new instance of question.Service.
func NewService(repo *Repository, llmService llm.ChatCompletioner) *Service {
	return &Service{
		repo:       repo,
		llmService: llmService,
	}
}

// StartGeneration generates questions about random topics at a given frequency.
func (s Service) StartGeneration(ctx context.Context, freq time.Duration, batchSize int) {
	slog.Info("Starting generation loop", slog.Duration("freq", freq))

	ticker := time.NewTicker(freq)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			topic := enums.RandomTopic()
			slog.Info(
				"Generating questions",
				slog.String("topic", topic.String()),
				slog.Int("amount", batchSize),
			)
			if err := s.handleGeneration(ctx, topic, batchSize); err != nil {
				if !errors.Is(err, context.Canceled) {
					slog.Error("Error generating question set", err)
				}
				return
			}
		}
	}
}

// handleGeneration generates and stores a set of questions about a given topic.
func (s Service) handleGeneration(ctx context.Context, topic enums.Topic, amount int) error {
	results := make(chan ai.Result)
	defer close(results)

	recentlyGenerated, err := s.getRecentlyGenerated(ctx, topic, 100)
	if err != nil {
		return errors.Wrapf(err, "getting recently generated questions about %s", topic)
	}

	go ai.Generate(ctx, s.llmService, topic, recentlyGenerated, amount, results)

	select {
	case <-ctx.Done():
		return nil
	case result := <-results:
		if result.Err != nil {
			return result.Err
		}
		for _, aiQuestion := range result.Questions {
			q, err := parseAIQuestion(aiQuestion, topic)
			if err != nil {
				return errors.Wrap(err, "parsing AI question")
			}
			slog.Info("Inserting question", slog.String("q", q.Question))
			if err := s.repo.Insert(ctx, *q); err != nil {
				return errors.Wrap(err, "inserting question")
			}
		}
	}
	return nil
}

// getRecentlyGenerated returns the most recent questions generated about a given topic.
func (s Service) getRecentlyGenerated(
	ctx context.Context,
	topic enums.Topic,
	amount int,
) ([]string, error) {
	questions, err := s.repo.GetMany(
		ctx,
		WhereTopicIs(topic),
		OrderByCreatedAtDesc(),
		Limit(amount),
	)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, q := range questions {
		result = append(result, q.Question)
	}
	return result, nil
}

// parseAIQuestion converts an ai.Question to a Question.
func parseAIQuestion(aiQuestion ai.Question, topic enums.Topic) (*Question, error) {
	difficulty, err := enums.DifficultyString(aiQuestion.Difficulty)
	if err != nil {
		return nil, errors.Wrap(err, "parsing difficulty")
	}

	q := New(aiQuestion.Question, aiQuestion.Hint, aiQuestion.MoreInfo).
		WithTopic(topic).
		WithDifficulty(difficulty)
	for _, c := range aiQuestion.Choices {
		q.WithChoice(c.Text, c.IsCorrect)
	}
	return q, nil
}
