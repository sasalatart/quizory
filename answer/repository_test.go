package answer_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart.com/quizory/answer"
	"github.com/sasalatart.com/quizory/db/testutil"
	"github.com/sasalatart.com/quizory/question"
	"github.com/stretchr/testify/suite"
)

type AnswerRepoTestSuite struct {
	suite.Suite

	testDB       *testutil.TestDB
	answerRepo   *answer.Repository
	questionRepo *question.Repository
}

func TestAnswerRepoTestSuite(t *testing.T) {
	suite.Run(t, new(AnswerRepoTestSuite))
}

func (s *AnswerRepoTestSuite) SetupSuite() {
	ctx := context.Background()
	testDB, err := testutil.NewTestDB(ctx)
	s.Require().NoError(err)

	s.testDB = testDB
	s.answerRepo = answer.NewRepository(testDB.DB())
	s.questionRepo = question.NewRepository(testDB.DB())
}

func (s *AnswerRepoTestSuite) TearDownSuite() {
	_ = s.testDB.Teardown()
}

func (s *AnswerRepoTestSuite) TearDownTest() {
	_ = s.testDB.DeleteData(context.Background())
}

func (s *AnswerRepoTestSuite) TestGetMany() {
	ctx := context.Background()

	userID1 := uuid.New()
	userID2 := uuid.New()

	q1 := question.Mock(nil)
	err := s.questionRepo.Insert(ctx, q1)
	s.Require().NoError(err)

	a1q1 := answer.New(userID1, q1.Choices[0].ID)
	err = s.answerRepo.Insert(ctx, *a1q1)
	s.Require().NoError(err)

	q2 := question.Mock(nil)
	err = s.questionRepo.Insert(ctx, q2)
	s.Require().NoError(err)

	a1q2 := answer.New(userID2, q2.Choices[1].ID)
	err = s.answerRepo.Insert(ctx, *a1q2)
	s.Require().NoError(err)
	a2q2 := answer.New(userID1, q2.Choices[0].ID)
	err = s.answerRepo.Insert(ctx, *a2q2)
	s.Require().NoError(err)

	answers, err := s.answerRepo.GetMany(
		ctx,
		answer.WhereUserID(userID1),
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

func (s *AnswerRepoTestSuite) TestInsert() {
	ctx := context.Background()

	q := question.Mock(nil)
	err := s.questionRepo.Insert(ctx, q)
	s.Require().NoError(err)

	userID := uuid.New()
	a := answer.New(userID, q.Choices[1].ID)
	err = s.answerRepo.Insert(ctx, *a)
	s.Require().NoError(err)

	answers, err := s.answerRepo.GetMany(ctx)
	s.Require().NoError(err)
	s.Require().Len(answers, 1)
	got := answers[0]

	s.Equal(a.ID.String(), got.ID.String())
	s.Equal(a.UserID.String(), got.UserID.String())
	s.Equal(a.ChoiceID.String(), got.ChoiceID.String())
}
