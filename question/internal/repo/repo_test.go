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

	q1 := question.
		New("Test Question 1", "Test Hint 1").
		WithChoice("Choice 11", false).
		WithChoice("Choice 12", true)
	err := s.r.Insert(ctx, *q1)
	s.Require().NoError(err)

	q2 := question.
		New("Test Question 2", "Test Hint 2").
		WithChoice("Choice 21", true).
		WithChoice("Choice 22", false).
		WithChoice("Choice 23", false)
	err = s.r.Insert(ctx, *q2)
	s.Require().NoError(err)

	got, err := s.r.GetMany(ctx, repo.OrderByCreatedAtDesc())
	s.Require().NoError(err)
	s.Require().Len(got, 2)

	want := []question.Question{*q2, *q1}
	for i, w := range want {
		g := got[i]
		s.Equal(w.ID.String(), g.ID.String())
		s.Equal(w.Question, g.Question)
		s.Equal(w.Hint, g.Hint)
		s.Require().Len(g.Choices, len(w.Choices))
		for j, c := range w.Choices {
			s.Equal(c.ID.String(), g.Choices[j].ID.String())
			s.Equal(c.Choice, g.Choices[j].Choice)
			s.Equal(c.IsCorrect, g.Choices[j].IsCorrect)
		}
	}
}

func (s *QuestionRepoTestSuite) TestInsert() {
	ctx := context.Background()
	q := question.New("Test Question 1", "Test Hint 1").
		WithChoice("Choice 11", false).
		WithChoice("Choice 12", true)

	err := s.r.Insert(ctx, *q)
	s.Require().NoError(err)

	questions, err := s.r.GetMany(ctx)
	s.Require().NoError(err)
	s.Require().Len(questions, 1)
	got := questions[0]

	s.Equal(q.ID.String(), got.ID.String())
	s.Equal(q.Question, got.Question)
	s.Equal(q.Hint, got.Hint)
	for _, c := range q.Choices {
		s.Contains(got.Choices, c)
	}
}
