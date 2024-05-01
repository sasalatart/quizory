package answer

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/question"
)

// Service represents the service that manages answers.
type Service struct {
	repo            *Repository
	questionService *question.Service
}

// NewService creates a new instance of answer.Service.
func NewService(repo *Repository, questionService *question.Service) *Service {
	return &Service{
		repo:            repo,
		questionService: questionService,
	}
}

type submissionResponse struct {
	CorrectChoiceID uuid.UUID
	MoreInfo        string
}

// Submit registers the choice made by a user for a specific question, and returns the correct
// choices for it, plus some more info for the user to know how they did.
func (s Service) Submit(
	ctx context.Context,
	userID, choiceID uuid.UUID,
) (*submissionResponse, error) {
	q, err := s.questionService.FromChoice(ctx, choiceID)
	if err != nil {
		return nil, errors.Wrapf(err, "getting question for choice %s", choiceID)
	}

	a := New(userID, choiceID)
	if err := s.repo.Insert(ctx, *a); err != nil {
		return nil, errors.Wrapf(err, "inserting answer %+v", a)
	}

	correctChoice, err := q.CorrectChoice()
	if err != nil {
		return nil, errors.Wrapf(err, "getting correct choice for question %s", q.ID)
	}

	return &submissionResponse{
		CorrectChoiceID: correctChoice.ID,
		MoreInfo:        q.MoreInfo,
	}, nil
}
