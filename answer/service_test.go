package answer_test

import (
	"context"
	"database/sql"
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

	resp, err := s.AnswerService.Submit(ctx, userID, choiceID)
	s.Require().NoError(err)
	s.Equal(q.Choices[1].ID, resp.CorrectChoiceID)
	s.Equal(q.MoreInfo, resp.MoreInfo)

	answers, err := s.AnswerRepo.GetMany(ctx, answer.WhereUserID(userID))
	s.Require().NoError(err)
	s.Len(answers, 1)
	gotAnswer := answers[0]
	s.Equal(userID, gotAnswer.UserID)
	s.Equal(choiceID, gotAnswer.ChoiceID)
}
