package answer_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/testutil"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type answerRepoTestSuiteParams struct {
	fx.In
	DB           *sql.DB
	AnswerRepo   *answer.Repository
	QuestionRepo *question.Repository
}

type AnswerRepoTestSuite struct {
	suite.Suite
	answerRepoTestSuiteParams
	app *fx.App
}

func TestAnswerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AnswerRepoTestSuite))
}

func (s *AnswerRepoTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.NopLogger,
		testutil.Module,
		fx.Populate(&s.answerRepoTestSuiteParams),
	)
	err := s.app.Start(context.Background())
	s.Require().NoError(err)
}

func (s *AnswerRepoTestSuite) TearDownSuite() {
	_ = s.app.Stop(context.Background())
}

func (s *AnswerRepoTestSuite) TearDownTest() {
	_ = testutil.DeleteData(context.Background(), s.DB)
}

func (s *AnswerRepoTestSuite) TestGetMany() {
	ctx := context.Background()

	userID1 := uuid.New()
	userID2 := uuid.New()

	q1 := question.Mock(nil)
	err := s.QuestionRepo.Insert(ctx, q1)
	s.Require().NoError(err)

	a1q1 := answer.New(userID1, q1.Choices[0].ID)
	err = s.AnswerRepo.Insert(ctx, *a1q1)
	s.Require().NoError(err)

	q2 := question.Mock(nil)
	err = s.QuestionRepo.Insert(ctx, q2)
	s.Require().NoError(err)

	a1q2 := answer.New(userID2, q2.Choices[1].ID)
	err = s.AnswerRepo.Insert(ctx, *a1q2)
	s.Require().NoError(err)
	a2q2 := answer.New(userID1, q2.Choices[0].ID)
	err = s.AnswerRepo.Insert(ctx, *a2q2)
	s.Require().NoError(err)

	answers, err := s.AnswerRepo.GetMany(
		ctx,
		answer.WhereUserIDEq(userID1),
		answer.OrderByCreatedAtDesc(),
	)
	s.Require().NoError(err)

	want := []*answer.Answer{a2q2, a1q1}
	s.Require().Len(answers, len(want))
	for i, got := range answers {
		s.Equal(want[i].ID.String(), got.ID.String())
		s.Equal(want[i].UserID.String(), got.UserID.String())
		s.Equal(want[i].ChoiceID.String(), got.ChoiceID.String())
	}
}

func (s *AnswerRepoTestSuite) TestGetOne() {
	ctx := context.Background()

	q := question.Mock(nil)
	err := s.QuestionRepo.Insert(ctx, q)
	s.Require().NoError(err)

	a1 := answer.New(uuid.New(), q.Choices[0].ID)
	err = s.AnswerRepo.Insert(ctx, *a1)
	s.Require().NoError(err)

	a2 := answer.New(uuid.New(), q.Choices[0].ID)
	err = s.AnswerRepo.Insert(ctx, *a2)
	s.Require().NoError(err)

	got, err := s.AnswerRepo.GetOne(ctx, answer.WhereIDEq(a1.ID))
	s.Require().NoError(err)
	s.Equal(a1.ID.String(), got.ID.String())
}

func (s *AnswerRepoTestSuite) TestInsert() {
	ctx := context.Background()

	q := question.Mock(nil)
	err := s.QuestionRepo.Insert(ctx, q)
	s.Require().NoError(err)

	userID := uuid.New()
	a := answer.New(userID, q.Choices[1].ID)
	err = s.AnswerRepo.Insert(ctx, *a)
	s.Require().NoError(err)

	answers, err := s.AnswerRepo.GetMany(ctx)
	s.Require().NoError(err)
	s.Require().Len(answers, 1)
	got := answers[0]

	s.Equal(a.ID.String(), got.ID.String())
	s.Equal(a.UserID.String(), got.UserID.String())
	s.Equal(a.ChoiceID.String(), got.ChoiceID.String())
}
