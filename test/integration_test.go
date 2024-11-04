package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart/quizory/domain/pagination"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/generator"
	"github.com/sasalatart/quizory/http/grpc/proto"
	"github.com/sasalatart/quizory/http/rest/oapi"
	"github.com/sasalatart/quizory/http/rest/resttest"
	"github.com/sasalatart/quizory/llm"
	"github.com/sasalatart/quizory/llm/llmtest"
	"github.com/sasalatart/quizory/test/testutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type testSuiteParams struct {
	fx.In
	DB                *sql.DB
	GRPCClient        proto.QuizoryServiceClient
	RESTClientFactory resttest.ClientFactory
	LLMService        llm.ChatCompletioner
	GeneratorService  *generator.Service
	QuestionRepo      *question.Repository
}

type TestSuite struct {
	suite.Suite
	testSuiteParams
	app *fx.App

	llmCalls int
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	ctx := context.Background()

	s.app = fx.New(
		fx.NopLogger,
		testutil.ModuleWithHTTP,
		fx.Populate(&s.testSuiteParams),
	)
	err := s.app.Start(ctx)
	s.Require().NoError(err)
}

func (s *TestSuite) TearDownSuite() {
	_ = s.app.Stop(context.Background())
}

func (s *TestSuite) TearDownTest() {
	_ = testutil.DeleteData(context.Background(), s.DB)
}

func (s *TestSuite) TestIntegration() {
	ctx := context.Background()

	ancientGreeceQuestions := s.assertGenerateBatchOfQuestions(ctx, 2, enums.TopicAncientGreece)
	ancientGreeceQ1 := ancientGreeceQuestions[0]
	ancientGreeceQ2 := ancientGreeceQuestions[1]

	ancientRomeQuestions := s.assertGenerateBatchOfQuestions(ctx, 1, enums.TopicAncientRome)
	ancientRomeQ1 := ancientRomeQuestions[0]

	userID := uuid.New()
	restClient, err := s.RESTClientFactory(userID)
	s.Require().NoError(err)

	s.assertRemainingTopicsAre(ctx, restClient, []oapi.RemainingTopic{
		{Topic: enums.TopicAncientGreece.String(), AmountOfQuestions: 2},
		{Topic: enums.TopicAncientRome.String(), AmountOfQuestions: 1},
	})
	s.assertNextQuestionFor(ctx, restClient, enums.TopicAncientGreece, ancientGreeceQ1)
	s.assertNextQuestionFor(ctx, restClient, enums.TopicAncientRome, ancientRomeQ1)
	s.assertNoNextQuestionFor(ctx, restClient, enums.TopicNapoleonicWars)

	s.submitAnswer(ctx, restClient, ancientGreeceQ1, ancientGreeceQ1.Choices[0].ID)

	s.assertRemainingTopicsAre(ctx, restClient, []oapi.RemainingTopic{
		{Topic: enums.TopicAncientGreece.String(), AmountOfQuestions: 1},
		{Topic: enums.TopicAncientRome.String(), AmountOfQuestions: 1},
	})
	s.assertNextQuestionFor(ctx, restClient, enums.TopicAncientGreece, ancientGreeceQ2)
	s.assertNextQuestionFor(ctx, restClient, enums.TopicAncientRome, ancientRomeQ1)
	s.assertNoNextQuestionFor(ctx, restClient, enums.TopicNapoleonicWars)

	s.submitAnswer(ctx, restClient, ancientGreeceQ2, ancientGreeceQ2.Choices[1].ID)

	s.assertRemainingTopicsAre(ctx, restClient, []oapi.RemainingTopic{
		{Topic: enums.TopicAncientRome.String(), AmountOfQuestions: 1},
	})
	s.assertNoNextQuestionFor(ctx, restClient, enums.TopicAncientGreece)
	s.assertNextQuestionFor(ctx, restClient, enums.TopicAncientRome, ancientRomeQ1)
	s.assertNoNextQuestionFor(ctx, restClient, enums.TopicNapoleonicWars)

	s.submitAnswer(ctx, restClient, ancientRomeQ1, ancientRomeQ1.Choices[0].ID)

	s.assertRemainingTopicsAre(ctx, restClient, []oapi.RemainingTopic{})
	s.assertNoNextQuestionFor(ctx, restClient, enums.TopicAncientGreece)
	s.assertNoNextQuestionFor(ctx, restClient, enums.TopicAncientRome)
	s.assertNoNextQuestionFor(ctx, restClient, enums.TopicNapoleonicWars)

	s.assertPaginatedLog(
		ctx,
		restClient,
		userID,
		pagination.Pagination{Page: 0, PageSize: 2},
		[]wantLogItem{
			{QuestionID: ancientRomeQ1.ID, ChoiceID: ancientRomeQ1.Choices[0].ID},
			{QuestionID: ancientGreeceQ2.ID, ChoiceID: ancientGreeceQ2.Choices[1].ID},
		},
	)
	s.assertPaginatedLog(
		ctx,
		restClient,
		userID,
		pagination.Pagination{Page: 1, PageSize: 2},
		[]wantLogItem{
			{QuestionID: ancientGreeceQ1.ID, ChoiceID: ancientGreeceQ1.Choices[0].ID},
		},
	)
	s.assertPaginatedLog(
		ctx,
		restClient,
		userID,
		pagination.Pagination{Page: 2, PageSize: 2},
		nil,
	)
}

// assertGenerateBatchOfQuestions checks that the GeneratorService.GenerateBatch feature produces
// the expected questions E2E. This function also returns the newly generated questions.
func (s *TestSuite) assertGenerateBatchOfQuestions(
	ctx context.Context,
	amount int,
	topic enums.Topic,
) []question.Question {
	getCurrentQuestionsIDs := func() []uuid.UUID {
		questions, err := s.QuestionRepo.GetMany(
			ctx,
			question.WhereTopicEq(topic),
			question.OrderByCreatedAtAsc(),
		)
		s.Require().NoError(err)
		var questionsIDs []uuid.UUID
		for _, q := range questions {
			questionsIDs = append(questionsIDs, q.ID)
		}
		return questionsIDs
	}

	previousQuestionsIDs := getCurrentQuestionsIDs()

	s.LLMService.(*llmtest.MockService).ChatCompletionFn = func(_, _ string) (string, error) {
		s.llmCalls++
		return newMockLLMResult(amount, s.llmCalls), nil
	}
	err := s.GeneratorService.GenerateBatch(ctx, amount, topic)
	s.Require().NoError(err)

	allQuestionsIDs := getCurrentQuestionsIDs()
	s.Require().Equal(len(previousQuestionsIDs)+amount, len(allQuestionsIDs))

	newQuestionsIDs := allQuestionsIDs[len(previousQuestionsIDs):]
	newQuestions, err := s.QuestionRepo.GetMany(
		ctx,
		question.WhereIDIn(newQuestionsIDs...),
		question.OrderByCreatedAtAsc(),
	)
	s.Require().NoError(err)

	for i, q := range newQuestions {
		qIndex := i + s.llmCalls
		s.Equal(fmt.Sprintf("Test question %d", qIndex), q.Question)
		s.Equal(fmt.Sprintf("Test hint %d", qIndex), q.Hint)
		s.Equal(topic.String(), q.Topic.String())
		s.Equal(fmt.Sprintf("Test more info %d\nTest fun fact", qIndex), q.MoreInfo)
		s.Equal("novice historian", q.Difficulty.String())
	}
	return newQuestions
}

func (s *TestSuite) assertRemainingTopicsAre(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	wantRemainingTopics []oapi.RemainingTopic,
) {
	res, err := client.GetRemainingTopics(ctx)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	got := parseResponse[[]oapi.RemainingTopic](s.T(), res)
	s.ElementsMatch(wantRemainingTopics, got)
}

func (s *TestSuite) assertNoNextQuestionFor(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	topic enums.Topic,
) {
	res, err := client.GetNextQuestion(ctx, &oapi.GetNextQuestionParams{
		Topic: topic.String(),
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, res.StatusCode)
}

func (s *TestSuite) assertNextQuestionFor(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	topic enums.Topic,
	wantQuestion question.Question,
) {
	res, err := client.GetNextQuestion(ctx, &oapi.GetNextQuestionParams{
		Topic: topic.String(),
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	got := parseResponse[oapi.UnansweredQuestion](s.T(), res)
	s.Equal(wantQuestion.ID.String(), got.Id.String())
	s.Len(got.Choices, len(wantQuestion.Choices))
}

func (s *TestSuite) submitAnswer(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	q question.Question,
	selectedChoice uuid.UUID,
) {
	correctChoice, err := q.CorrectChoice()
	s.Require().NoError(err)

	res, err := client.SubmitAnswer(ctx, oapi.SubmitAnswerJSONRequestBody{
		ChoiceId: selectedChoice,
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, res.StatusCode)
	got := parseResponse[oapi.SubmitAnswerResult](s.T(), res)
	s.Equal(correctChoice.ID.String(), got.CorrectChoiceId.String())
	s.Equal(q.MoreInfo, got.MoreInfo)
}

type wantLogItem struct {
	QuestionID uuid.UUID
	ChoiceID   uuid.UUID
}

func (s *TestSuite) assertPaginatedLog(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	userID uuid.UUID,
	p pagination.Pagination,
	wantLog []wantLogItem,
) {
	res, err := client.GetAnswersLog(ctx, userID, &oapi.GetAnswersLogParams{
		Page:     &p.Page,
		PageSize: &p.PageSize,
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)

	gotLog := parseResponse[[]oapi.AnswersLogItem](s.T(), res)
	s.Require().Len(gotLog, len(wantLog))

	for i, want := range wantLog {
		baseMsg := fmt.Sprintf("Log item %d for page %d", i, p.Page)
		got := gotLog[i]
		s.Equal(want.QuestionID.String(), got.Question.Id.String(), baseMsg)
		s.Equal(want.ChoiceID.String(), got.ChoiceId.String(), baseMsg)
	}
}

func parseResponse[T any](t *testing.T, res *http.Response) T {
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var response T
	require.NoError(t, json.Unmarshal(body, &response))
	return response
}

// newMockLLMResult generates a JSON string with a list of questions that can be used as a mock
// response from the LLM.
func newMockLLMResult(amount, startIndex int) string {
	var questions []map[string]interface{}
	for i := range amount {
		questionIndex := i + startIndex
		questions = append(questions, map[string]interface{}{
			"question": fmt.Sprintf("Test question %d", questionIndex),
			"hint":     fmt.Sprintf("Test hint %d", questionIndex),
			"choices": []map[string]interface{}{
				{"text": fmt.Sprintf("Choice A%d", questionIndex), "isCorrect": false},
				{"text": fmt.Sprintf("Choice B%d", questionIndex), "isCorrect": true},
				{"text": fmt.Sprintf("Choice C%d", questionIndex), "isCorrect": false},
				{"text": fmt.Sprintf("Choice D%d", questionIndex), "isCorrect": false},
			},
			"moreInfo":   []string{fmt.Sprintf("Test more info %d", questionIndex), "Test fun fact"},
			"difficulty": "novice historian",
		})
	}
	result, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		return ""
	}
	return string(result)
}
