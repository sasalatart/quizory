package repo_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sasalatart.com/quizory/db/testutil"
	"github.com/sasalatart.com/quizory/question"
	"github.com/sasalatart.com/quizory/question/internal/repo"
	"github.com/stretchr/testify/suite"
)

type QuestionRepoTestSuite struct {
	suite.Suite

	db *sql.DB
	r  *repo.QuestionRepo
}

func (s *QuestionRepoTestSuite) SetupSuite() {
	db, teardown, err := testutil.NewDB()
	s.Require().NoError(err)
	s.T().Cleanup(teardown)

	s.db = db
	s.r = repo.New(db)
}

func (s *QuestionRepoTestSuite) TearDownTest() {
	testutil.WipeDB(context.Background(), s.db)
}

func TestQuestionRepoTestSuite(t *testing.T) {
	suite.Run(t, new(QuestionRepoTestSuite))
}

func (s *QuestionRepoTestSuite) TestGetMany() {
	ctx := context.Background()

	q1 := question.New("Test Question 1", "Test Hint 1")
	err := s.r.Insert(ctx, q1)
	s.Require().NoError(err)

	q2 := question.New("Test Question 2", "Test Hint 2")
	err = s.r.Insert(ctx, q2)
	s.Require().NoError(err)

	got, err := s.r.GetMany(ctx, repo.OrderByCreatedAtDesc())
	s.Require().NoError(err)
	s.Require().Len(got, 2)

	want := []question.Question{q2, q1}
	for i, w := range want {
		s.Equal(w.ID.String(), got[i].ID.String())
		s.Equal(w.Question, got[i].Question)
		s.Equal(w.Hint, got[i].Hint)
	}
}

func (s *QuestionRepoTestSuite) TestInsert() {
	ctx := context.Background()
	q := question.New("Test Question", "Test Hint")

	err := s.r.Insert(ctx, q)
	s.Require().NoError(err)

	questions, err := s.r.GetMany(ctx)
	s.Require().NoError(err)
	s.Require().Len(questions, 1)
	got := questions[0]

	s.Equal(q.ID.String(), got.ID.String())
	s.Equal(q.Question, got.Question)
	s.Equal(q.Hint, got.Hint)
}
