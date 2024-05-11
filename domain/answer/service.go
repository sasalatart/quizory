package answer

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/domain/pagination"
	"github.com/sasalatart.com/quizory/domain/question"
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

type SubmissionResponse struct {
	ID              uuid.UUID
	CorrectChoiceID uuid.UUID
	MoreInfo        string
}

// Submit registers the choice made by a user for a specific question, and returns the correct
// choice for it, plus some more info for the user to know how they did.
func (s Service) Submit(
	ctx context.Context,
	userID, choiceID uuid.UUID,
) (*SubmissionResponse, error) {
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

	return &SubmissionResponse{
		ID:              a.ID,
		CorrectChoiceID: correctChoice.ID,
		MoreInfo:        q.MoreInfo,
	}, nil
}

// LogItem represents a previous attempt at answering a question.
type LogItem struct {
	ID       uuid.UUID
	Question question.Question
	ChoiceID uuid.UUID
}

// IsSuccessful returns whether the answer was correct or not.
func (l LogItem) IsSuccessful() (bool, error) {
	correctChoice, err := l.Question.CorrectChoice()
	if err != nil {
		return false, errors.Wrapf(err, "getting correct choice for question %s", l.Question.ID)
	}
	return l.ChoiceID == correctChoice.ID, nil
}

// LogRequest has the parameters for retrieving a user's history of answers.
type LogRequest struct {
	UserID     uuid.UUID
	Pagination pagination.Pagination
}

// LogFor returns the paginated history of answers for a user.
func (s Service) LogFor(ctx context.Context, req LogRequest) ([]LogItem, error) {
	answers, err := s.repo.GetMany(
		ctx,
		WhereUserIDEq(req.UserID),
		OrderByCreatedAtDesc(),
		Offset(req.Pagination.Offset()),
		Limit(req.Pagination.PageSize),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "getting answers for %+v", req)
	}

	if len(answers) == 0 {
		return []LogItem{}, nil
	}

	choicesIDs := make([]uuid.UUID, len(answers))
	for _, a := range answers {
		choicesIDs = append(choicesIDs, a.ChoiceID)
	}
	questions, err := s.questionService.FromChoices(ctx, choicesIDs...)
	if err != nil {
		return nil, errors.Wrapf(err, "getting questions for choices %+v", choicesIDs)
	}

	result, err := s.composeLog(answers, questions)
	if err != nil {
		return nil, errors.Wrapf(err, "composing log for %+v", req)
	}
	return result, nil
}

// composeLog returns the log of answers for a user, given a set of answers and their questions.
func (s Service) composeLog(answers []Answer, questions []question.Question) ([]LogItem, error) {
	findQuestion := func(a Answer) (*question.Question, error) {
		for _, q := range questions {
			for _, c := range q.Choices {
				if c.ID == a.ChoiceID {
					return &q, nil
				}
			}
		}
		return nil, errors.New("question not found")
	}

	var result []LogItem
	for _, a := range answers {
		q, err := findQuestion(a)
		if err != nil {
			return nil, errors.Wrapf(err, "finding question for answer %s", a.ID)
		}

		result = append(result, LogItem{
			ID:       a.ID,
			Question: *q,
			ChoiceID: a.ChoiceID,
		})
	}
	return result, nil
}
