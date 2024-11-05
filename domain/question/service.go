package question

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/domain/question/enums"
)

// ErrNoQuestionsLeft is returned when there are no questions left to be answered for a user.
var ErrNoQuestionsLeft = errors.New("no questions left")

// Service represents the service that manages questions.
type Service struct {
	repo *Repository
}

// NewService creates a new instance of question.Service.
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Insert inserts a new question.
func (s Service) Insert(ctx context.Context, q *Question) error {
	return s.repo.Insert(ctx, *q)
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
	topic enums.Topic,
) (*Question, error) {
	q, err := s.repo.GetOne(
		ctx,
		WhereTopicEq(topic),
		WhereNotAnsweredBy(userID),
		OrderByCreatedAtAsc(),
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrapf(
			ErrNoQuestionsLeft,
			"getting next question about %s for %s",
			topic, userID,
		)
	}
	return q, err
}

// RemainingTopicsFor returns a map such that each key is a topic for which the user still has
// unanswered questions, and each value is the amount of remaining questions for that topic.
func (s Service) RemainingTopicsFor(
	ctx context.Context,
	userID uuid.UUID,
) (map[enums.Topic]uint, error) {
	// Arguably the following logic could be owned by the service instead of the repository so that
	// we avoid making this a pass-through method. However, if we do this at the service level by
	// loading all unanswered questions for a user and grouping via code, we risk loading ALL
	// questions in the database in the worst case (e.g. if the user has not answered any question).
	// Thus, the tradeoff is to do this in the repository level, where we can leverage the database
	// to avoid this issue.
	return s.repo.GetRemainingTopics(ctx, userID)
}

// Latest returns the most recent questions about a given topic, capped to amount.
func (s Service) Latest(ctx context.Context, topic enums.Topic, amount int) ([]Question, error) {
	return s.repo.GetMany(
		ctx,
		WhereTopicEq(topic),
		OrderByCreatedAtDesc(),
		Limit(amount),
	)
}
