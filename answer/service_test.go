package answer_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart.com/quizory/answer"
	"github.com/sasalatart.com/quizory/question"
	"github.com/sasalatart.com/quizory/testutil"
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
		answer.Module,
		question.Module,
		fx.Provide(answer.NewRepository),
		fx.Provide(question.NewRepository),
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

	testCases := []struct {
		name string
		req  answer.LogRequest
		want []answer.LogItem
	}{
		{
			name: "First Page Full",
			req:  answer.LogRequest{UserID: userID, Page: 0, PerPage: 2},
			want: []answer.LogItem{
				{
					ChoiceID:     q3.Choices[0].ID,
					Question:     q3,
					IsSuccessful: false,
				},
				{
					ChoiceID:     q2.Choices[1].ID,
					Question:     q2,
					IsSuccessful: true,
				},
			},
		},
		{
			name: "Second Page With Some Results",
			req:  answer.LogRequest{UserID: userID, Page: 1, PerPage: 2},
			want: []answer.LogItem{
				{
					ChoiceID:     q1.Choices[0].ID,
					Question:     q1,
					IsSuccessful: false,
				},
			},
		},
		{
			name: "Third Page Empty",
			req:  answer.LogRequest{UserID: userID, Page: 2, PerPage: 2},
			want: []answer.LogItem{},
		},
		{
			name: "Another User's Log",
			req:  answer.LogRequest{UserID: anotherUserID, Page: 0, PerPage: 2},
			want: []answer.LogItem{
				{
					ChoiceID:     q1.Choices[0].ID,
					Question:     q1,
					IsSuccessful: false,
				},
			},
		},
		{
			name: "Without Any Answers",
			req:  answer.LogRequest{UserID: uuid.New(), Page: 0, PerPage: 2},
			want: []answer.LogItem{},
		},
	}
	for _, tt := range testCases {
		s.Run(tt.name, func() {
			got, err := s.AnswerService.LogFor(ctx, tt.req)
			s.Require().NoError(err)
			s.Require().Len(got, len(tt.want))
			for i, want := range tt.want {
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
				s.Equal(
					want.IsSuccessful,
					got[i].IsSuccessful,
					"IsSuccessful at index %d", i,
				)
			}
		})
	}
}
