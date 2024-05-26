package question_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sasalatart.com/quizory/domain/answer"
	"github.com/sasalatart.com/quizory/domain/question"
	"github.com/sasalatart.com/quizory/llm"
	"github.com/sasalatart.com/quizory/testutil"
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

	q1 := question.Mock(nil)
	err := s.QuestionRepo.Insert(ctx, q1)
	s.Require().NoError(err)

	q2 := question.Mock(nil)
	err = s.QuestionRepo.Insert(ctx, q2)
	s.Require().NoError(err)

	got, err := s.Service.NextFor(ctx, userID)
	s.Require().NoError(err)
	s.Equal(q1.ID, got.ID)

	a1 := answer.New(userID, q1.Choices[0].ID)
	err = s.AnswerRepo.Insert(ctx, *a1)
	s.Require().NoError(err)

	got, err = s.Service.NextFor(ctx, userID)
	s.Require().NoError(err)
	s.Equal(q2.ID, got.ID)

	a2 := answer.New(userID, q2.Choices[0].ID)
	err = s.AnswerRepo.Insert(ctx, *a2)
	s.Require().NoError(err)

	_, err = s.Service.NextFor(ctx, userID)
	s.Require().ErrorIs(err, question.ErrNoQuestionsLeft)
}
