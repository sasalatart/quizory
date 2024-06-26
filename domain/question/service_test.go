package question_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/llm"
	"github.com/sasalatart/quizory/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type questionServiceTestSuiteParams struct {
	fx.In
	DB           *sql.DB
	LLMService   llm.ChatCompletioner
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

func (s *QuestionServiceTestSuite) TestStartGeneration() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	llmCallsDone := 0
	s.LLMService.(*llm.MockService).ChatCompletionFn = func(_, _ string) (string, error) {
		llmCallsDone++
		return fmt.Sprintf(`
			[
				{
					"question": "Test question %d",
					"hint": "Test hint",
					"choices": [
						{"text": "Choice A", "isCorrect": false},
						{"text": "Choice B", "isCorrect": true},
						{"text": "Choice C", "isCorrect": false},
						{"text": "Choice D", "isCorrect": false}
					],
					"moreInfo": ["Test more info", "Test fun fact"],
					"difficulty": "novice historian"
				}
			]`, llmCallsDone), nil
	}

	freq := 500 * time.Millisecond
	batchSize := 1
	go s.Service.StartGeneration(ctx, freq, batchSize)

	s.EventuallyWithT(func(c *assert.CollectT) {
		questions, err := s.QuestionRepo.GetMany(ctx, question.OrderByCreatedAtDesc())
		assert.NoError(c, err)
		assert.Len(c, questions, 2)

		for i, q := range questions {
			assert.Equal(c, fmt.Sprintf("Test question %d", len(questions)-i), q.Question)
		}
	}, 5*time.Second, 500*time.Millisecond)
}

func (s *QuestionServiceTestSuite) TestNextFor() {
	ctx := context.Background()
	userID := uuid.New()

	q1 := s.mustSeedQuestion(nil)
	q2 := s.mustSeedQuestion(nil)

	got, err := s.Service.NextFor(ctx, userID)
	s.Require().NoError(err)
	s.Equal(q1.ID, got.ID)

	s.mustSeedAnswer(userID, q1.Choices[0].ID)

	got, err = s.Service.NextFor(ctx, userID)
	s.Require().NoError(err)
	s.Equal(q2.ID, got.ID)

	s.mustSeedAnswer(userID, q2.Choices[0].ID)

	_, err = s.Service.NextFor(ctx, userID)
	s.Require().ErrorIs(err, question.ErrNoQuestionsLeft)
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
