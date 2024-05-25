package question

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/domain/question/enums"
	"github.com/sasalatart.com/quizory/domain/question/internal/ai"
	"github.com/sasalatart.com/quizory/llm"
)

// ErrNoQuestionsLeft is returned when there are no questions left to be answered for a user.
var ErrNoQuestionsLeft = errors.New("no questions left")

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
	slog.Info("Starting questions generation loop", slog.Duration("freq", freq))

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
					slog.Error("Error generating question set", "error", err)
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

	recentlyGenerated, err := s.recentlyGenerated(ctx, topic, 100)
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
			if err := s.repo.Insert(context.WithoutCancel(ctx), *q); err != nil {
				return errors.Wrap(err, "inserting question")
			}
		}
	}
	return nil
}

// recentlyGenerated returns the most recent questions generated about a given topic.
func (s Service) recentlyGenerated(
	ctx context.Context,
	topic enums.Topic,
	amount int,
) ([]string, error) {
	questions, err := s.repo.GetMany(
		ctx,
		WhereTopicEq(topic),
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

// FromChoice returns the question associated with a given choice.
func (s Service) FromChoice(ctx context.Context, choiceID uuid.UUID) (*Question, error) {
	return s.repo.GetOne(ctx, WhereChoiceIDIn(choiceID))
}

// FromChoices returns the questions associated with a given set of choices.
func (s Service) FromChoices(ctx context.Context, ids ...uuid.UUID) ([]Question, error) {
	return s.repo.GetMany(ctx, WhereChoiceIDIn(ids...))
}

// NextFor returns the next question that a user should answer.
func (s Service) NextFor(
	ctx context.Context,
	userID uuid.UUID,
) (*Question, error) {
	q, err := s.repo.GetOne(
		ctx,
		WhereNotAnsweredBy(userID),
		OrderByCreatedAtAsc(),
	)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(ErrNoQuestionsLeft, "getting next question for %s", userID)
	}
	return q, err
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
