package question_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sasalatart.com/quizory/domain/answer"
	"github.com/sasalatart.com/quizory/domain/question"
	"github.com/sasalatart.com/quizory/testutil"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type questionRepoTestSuiteParams struct {
	fx.In
	DB   *sql.DB
	Repo *question.Repository
}

type QuestionRepoTestSuite struct {
	suite.Suite
	questionRepoTestSuiteParams
	app *fx.App
}

func TestQuestionRepoTestSuite(t *testing.T) {
	suite.Run(t, new(QuestionRepoTestSuite))
}

func (s *QuestionRepoTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.NopLogger,
		testutil.Module,
		fx.Populate(&s.questionRepoTestSuiteParams),
	)
	err := s.app.Start(context.Background())
	s.Require().NoError(err)
}

func (s *QuestionRepoTestSuite) TearDownSuite() {
	_ = s.app.Stop(context.Background())
}

func (s *QuestionRepoTestSuite) TearDownTest() {
	_ = testutil.DeleteData(context.Background(), s.DB)
}

func (s *QuestionRepoTestSuite) TestGetMany() {
	ctx := context.Background()

	q1 := question.Mock(func(q *question.Question) {
		q.Question = "Test Question 1"
		q.Hint = "Test Hint 1"
		q.MoreInfo = "Test More Info 1"
	})
	err := s.Repo.Insert(ctx, q1)
	s.Require().NoError(err)

	q2 := question.Mock(func(q *question.Question) {
		q.Question = "Test Question 2"
		q.Hint = "Test Hint 2"
		q.MoreInfo = "Test More Info 2"
	})
	err = s.Repo.Insert(ctx, q2)
	s.Require().NoError(err)

	got, err := s.Repo.GetMany(ctx, answer.OrderByCreatedAtDesc())
	s.Require().NoError(err)

	want := []question.Question{q2, q1}
	s.Require().Len(got, len(want))
	for i, w := range want {
		g := got[i]
		s.Equal(w.ID.String(), g.ID.String())
		s.Equal(w.Topic, g.Topic)
		s.Equal(w.Question, g.Question)
		s.Equal(w.Hint, g.Hint)
		s.Equal(w.MoreInfo, g.MoreInfo)
		s.Equal(w.Difficulty, g.Difficulty)
		s.Require().Len(g.Choices, len(w.Choices))
		for j, c := range w.Choices {
			s.Equal(c.ID.String(), g.Choices[j].ID.String())
			s.Equal(c.Choice, g.Choices[j].Choice)
			s.Equal(c.IsCorrect, g.Choices[j].IsCorrect)
		}
	}
}

func (s *QuestionRepoTestSuite) TestGetOne() {
	ctx := context.Background()

	q1 := question.Mock(nil)
	err := s.Repo.Insert(ctx, q1)
	s.Require().NoError(err)

	q2 := question.Mock(nil)
	err = s.Repo.Insert(ctx, q2)
	s.Require().NoError(err)

	got, err := s.Repo.GetOne(ctx, question.WhereChoiceIDIn(q1.Choices[0].ID))
	s.Require().NoError(err)
	s.Equal(q1.ID.String(), got.ID.String())
}

func (s *QuestionRepoTestSuite) TestInsert() {
	ctx := context.Background()
	q := question.Mock(nil)

	err := s.Repo.Insert(ctx, q)
	s.Require().NoError(err)

	questions, err := s.Repo.GetMany(ctx)
	s.Require().NoError(err)
	s.Require().Len(questions, 1)
	got := questions[0]

	s.Equal(q.ID.String(), got.ID.String())
	s.Equal(q.Topic, got.Topic)
	s.Equal(q.Question, got.Question)
	s.Equal(q.Hint, got.Hint)
	s.Equal(q.MoreInfo, got.MoreInfo)
	s.Equal(q.Difficulty, got.Difficulty)
	for _, c := range q.Choices {
		s.Contains(got.Choices, c)
	}
}
