package question_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/llm"
	"github.com/sasalatart/quizory/test/testutil"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type questionServiceTestSuiteParams struct {
	fx.In
	DB           *sql.DB
	LLMService   llm.Chater
	AnswerRepo   *answer.Repository
	QuestionRepo *question.Repository
	Service      *question.Service
}

type QuestionServiceTestSuite struct {
	suite.Suite
	questionServiceTestSuiteParams
	app *fx.App
}

func TestQuestionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(QuestionServiceTestSuite))
}

func (s *QuestionServiceTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.NopLogger,
		testutil.Module,
		fx.Populate(&s.questionServiceTestSuiteParams),
	)
	err := s.app.Start(context.Background())
	s.Require().NoError(err)
}

func (s *QuestionServiceTestSuite) TearDownSuite() {
	_ = s.app.Stop(context.Background())
}

func (s *QuestionServiceTestSuite) TearDownTest() {
	_ = testutil.DeleteData(context.Background(), s.DB)
}

func (s *QuestionServiceTestSuite) TestNextFor() {
	ctx := context.Background()
	userID1 := uuid.New()
	userID2 := uuid.New()

	q1 := s.mustSeedQuestion(func(q *question.Question) {
		q.Topic = enums.TopicNapoleonicWars
	})
	q2 := s.mustSeedQuestion(func(q *question.Question) {
		q.Topic = enums.TopicFrenchRevolution
	})
	q3 := s.mustSeedQuestion(func(q *question.Question) {
		q.Topic = enums.TopicFrenchRevolution
	})

	// Next question is the oldest unanswered question for the given topic
	s.assertNextQuestionIs(enums.TopicFrenchRevolution, userID1, q2.ID)

	// ...and will be the same for the topic unless the user answers it
	s.assertNextQuestionIs(enums.TopicFrenchRevolution, userID1, q2.ID)

	// ...only to change once it has been answered
	s.mustSeedAnswer(userID1, q2.Choices[0].ID)
	s.assertNextQuestionIs(enums.TopicFrenchRevolution, userID1, q3.ID)

	// ...which does not affect the next question for the same topic for other users
	s.assertNextQuestionIs(enums.TopicFrenchRevolution, userID2, q2.ID)

	// ...but the user may change the topic without answering the current question for a given topic
	s.assertNextQuestionIs(enums.TopicNapoleonicWars, userID1, q1.ID)

	// ...and a sentinel error is returned when there are no more questions for the given topic
	_, err := s.Service.NextFor(ctx, userID1, enums.TopicAncientGreece)
	s.Require().ErrorIs(err, question.ErrNoQuestionsLeft)
}

func (s *QuestionServiceTestSuite) assertNextQuestionIs(
	topic enums.Topic,
	userID uuid.UUID,
	wantQuestionID uuid.UUID,
) {
	s.T().Helper()
	got, err := s.Service.NextFor(context.Background(), userID, topic)
	s.Require().NoError(err)
	s.Equal(wantQuestionID.String(), got.ID.String())
}

func (s *QuestionServiceTestSuite) TestRemainingTopicsFor() {
	userID1 := uuid.New()
	userID2 := uuid.New()
	userID3 := uuid.New()

	q1 := s.mustSeedQuestion(func(q *question.Question) {
		q.Topic = enums.TopicNapoleonicWars
	})
	s.mustSeedAnswer(userID1, q1.Choices[0].ID)

	q2 := s.mustSeedQuestion(func(q *question.Question) {
		q.Topic = enums.TopicNapoleonicWars
	})
	s.mustSeedAnswer(userID1, q2.Choices[0].ID)

	q3 := s.mustSeedQuestion(func(q *question.Question) {
		q.Topic = enums.TopicFrenchRevolution
	})
	s.mustSeedAnswer(userID1, q3.Choices[0].ID)
	s.mustSeedAnswer(userID2, q3.Choices[0].ID)

	testCases := []struct {
		name   string
		userID uuid.UUID
		want   map[enums.Topic]uint
	}{
		{
			name:   "A user that has answered all questions",
			userID: userID1,
			want:   map[enums.Topic]uint{},
		},
		{
			name:   "A user that has answered some questions",
			userID: userID2,
			want: map[enums.Topic]uint{
				enums.TopicNapoleonicWars: 2,
			},
		}, {
			name:   "A user that has not answered any question",
			userID: userID3,
			want: map[enums.Topic]uint{
				enums.TopicFrenchRevolution: 1,
				enums.TopicNapoleonicWars:   2,
			},
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			got, err := s.Service.RemainingTopicsFor(context.Background(), tc.userID)
			s.Require().NoError(err)
			s.Equal(tc.want, got)
		})
	}
}

func (s *QuestionServiceTestSuite) mustSeedQuestion(
	overrides func(q *question.Question),
) question.Question {
	q := question.Mock(overrides)
	err := s.QuestionRepo.Insert(context.Background(), q)
	s.Require().NoError(err)
	return q
}

func (s *QuestionServiceTestSuite) mustSeedAnswer(
	userID uuid.UUID,
	choiceID uuid.UUID,
) answer.Answer {
	a := answer.New(userID, choiceID)
	err := s.AnswerRepo.Insert(context.Background(), *a)
	s.Require().NoError(err)
	return *a
}
