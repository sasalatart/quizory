package answer_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/pagination"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/testutil"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type answerServiceTestSuiteParams struct {
	fx.In
	DB            *sql.DB
	AnswerService *answer.Service
	AnswerRepo    *answer.Repository
	QuestionRepo  *question.Repository
}

type AnswerServiceTestSuite struct {
	suite.Suite
	answerServiceTestSuiteParams
	app *fx.App
}

func TestAnswerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AnswerServiceTestSuite))
}

func (s *AnswerServiceTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.NopLogger,
		testutil.Module,
		fx.Populate(&s.answerServiceTestSuiteParams),
	)
	err := s.app.Start(context.Background())
	s.Require().NoError(err)
}

func (s *AnswerServiceTestSuite) TearDownSuite() {
	_ = s.app.Stop(context.Background())
}

func (s *AnswerServiceTestSuite) TearDownTest() {
	_ = testutil.DeleteData(context.Background(), s.DB)
}

func (s *AnswerServiceTestSuite) TestSubmit() {
	ctx := context.Background()

	q := question.Mock(func(q *question.Question) {
		q.Choices = nil
		q.WithChoice("Choice 1", false).WithChoice("Choice 2", true)
		q.MoreInfo = "Test Some More Info"
	})
	err := s.QuestionRepo.Insert(ctx, q)
	s.Require().NoError(err)

	userID := uuid.New()
	choiceID := q.Choices[0].ID

	got, err := s.AnswerService.Submit(ctx, userID, choiceID)
	s.Require().NoError(err)
	s.Equal(q.Choices[1].ID.String(), got.CorrectChoiceID.String())
	s.Equal(q.MoreInfo, got.MoreInfo)

	a, err := s.AnswerRepo.GetOne(ctx, answer.WhereIDEq(got.ID))
	s.Require().NoError(err)
	s.Equal(userID, a.UserID)
	s.Equal(choiceID, a.ChoiceID)
}

func (s *AnswerServiceTestSuite) TestLogFor() {
	ctx := context.Background()
	userID := uuid.New()
	anotherUserID := uuid.New()

	mustCreateQuestion := func(qNumber int) question.Question {
		q := question.Mock(func(q *question.Question) {
			q.Choices = nil
			q.
				WithChoice(fmt.Sprintf("Choice %da", qNumber), false).
				WithChoice(fmt.Sprintf("Choice %db", qNumber), true)
		})
		err := s.QuestionRepo.Insert(ctx, q)
		s.Require().NoError(err)
		return q
	}

	mustSubmitAnswer := func(userID uuid.UUID, choiceID uuid.UUID) {
		_, err := s.AnswerService.Submit(ctx, userID, choiceID)
		s.Require().NoError(err)
	}

	q1 := mustCreateQuestion(1)
	q2 := mustCreateQuestion(2)
	q3 := mustCreateQuestion(3)
	_ = mustCreateQuestion(4)

	mustSubmitAnswer(userID, q1.Choices[0].ID)
	mustSubmitAnswer(anotherUserID, q1.Choices[0].ID)
	mustSubmitAnswer(userID, q2.Choices[1].ID)
	mustSubmitAnswer(userID, q3.Choices[0].ID)

	toPtr := func(i int) *int { return &i }

	testCases := []struct {
		name string
		req  answer.LogRequest
		want []answer.LogItem
	}{
		{
			name: "First Page Full",
			req: answer.LogRequest{
				UserID:     userID,
				Pagination: pagination.New(toPtr(0), toPtr(2)),
			},
			want: []answer.LogItem{
				{
					ChoiceID: q3.Choices[0].ID,
					Question: q3,
				},
				{
					ChoiceID: q2.Choices[1].ID,
					Question: q2,
				},
			},
		},
		{
			name: "Second Page With Some Results",
			req: answer.LogRequest{
				UserID:     userID,
				Pagination: pagination.New(toPtr(1), toPtr(2)),
			},
			want: []answer.LogItem{
				{
					ChoiceID: q1.Choices[0].ID,
					Question: q1,
				},
			},
		},
		{
			name: "Third Page Empty",
			req: answer.LogRequest{
				UserID:     userID,
				Pagination: pagination.New(toPtr(2), toPtr(2)),
			},
			want: []answer.LogItem{},
		},
		{
			name: "Another User's Log",
			req: answer.LogRequest{
				UserID:     anotherUserID,
				Pagination: pagination.New(toPtr(0), toPtr(2)),
			},
			want: []answer.LogItem{
				{
					ChoiceID: q1.Choices[0].ID,
					Question: q1,
				},
			},
		},
		{
			name: "Without Any Answers",
			req: answer.LogRequest{
				UserID:     uuid.New(),
				Pagination: pagination.New(toPtr(0), toPtr(2)),
			},
			want: []answer.LogItem{},
		},
	}
	for _, tt := range testCases {
		s.Run(tt.name, func() {
			got, err := s.AnswerService.LogFor(ctx, tt.req)
			s.Require().NoError(err)
			s.Require().Len(got, len(tt.want))

			for i, want := range tt.want {
				s.NotEqual(
					uuid.Nil.String(),
					got[i].ID.String(),
					"ID at index %d", i,
				)
				s.Equalf(
					want.ChoiceID.String(),
					got[i].ChoiceID.String(),
					"ChoiceID at index %d", i,
				)
				s.Equal(
					want.Question.ID.String(),
					got[i].Question.ID.String(),
					"Question ID at index %d", i,
				)
			}
		})
	}
}
