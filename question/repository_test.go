package question_test

import (
	"context"
	"testing"

	"github.com/sasalatart.com/quizory/answer"
	"github.com/sasalatart.com/quizory/db/testutil"
	"github.com/sasalatart.com/quizory/question"
	"github.com/stretchr/testify/suite"
)

type QuestionRepoTestSuite struct {
	suite.Suite

	testDB *testutil.TestDB
	r      *question.Repository
}

func TestQuestionRepoTestSuite(t *testing.T) {
	suite.Run(t, new(QuestionRepoTestSuite))
}

func (s *QuestionRepoTestSuite) SetupSuite() {
	ctx := context.Background()
	testDB, err := testutil.NewTestDB(ctx)
	s.Require().NoError(err)

	s.testDB = testDB
	s.r = question.NewRepository(testDB.DB())
}

func (s *QuestionRepoTestSuite) TearDownSuite() {
	_ = s.testDB.Teardown()
}

func (s *QuestionRepoTestSuite) TearDownTest() {
	_ = s.testDB.DeleteData(context.Background())
}

func (s *QuestionRepoTestSuite) TestGetMany() {
	ctx := context.Background()

	q1 := question.Mock(func(q *question.Question) {
		q.Question = "Test Question 1"
		q.Hint = "Test Hint 1"
		q.MoreInfo = "Test More Info 1"
	})
	err := s.r.Insert(ctx, q1)
	s.Require().NoError(err)

	q2 := question.Mock(func(q *question.Question) {
		q.Question = "Test Question 2"
		q.Hint = "Test Hint 2"
		q.MoreInfo = "Test More Info 2"
	})
	err = s.r.Insert(ctx, q2)
	s.Require().NoError(err)

	got, err := s.r.GetMany(ctx, answer.OrderByCreatedAtDesc())
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

func (s *QuestionRepoTestSuite) TestInsert() {
	ctx := context.Background()
	q := question.Mock(nil)

	err := s.r.Insert(ctx, q)
	s.Require().NoError(err)

	questions, err := s.r.GetMany(ctx)
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
