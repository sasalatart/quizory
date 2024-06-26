package server_test

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
	"github.com/sasalatart/quizory/http/oapi"
	"github.com/sasalatart/quizory/http/server"
	"github.com/sasalatart/quizory/testutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type serverTestSuiteParams struct {
	fx.In
	DB            *sql.DB
	ClientFactory server.TestClientFactory
	QuestionRepo  *question.Repository
}

type ServerTestSuite struct {
	suite.Suite
	serverTestSuiteParams
	app *fx.App

	ancientGreeceQ1 question.Question
	ancientGreeceQ2 question.Question
	ancientRomeQ1   question.Question
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	ctx := context.Background()

	s.app = fx.New(
		fx.NopLogger,
		testutil.ModuleWithAPI,
		fx.Populate(&s.serverTestSuiteParams),
	)
	err := s.app.Start(ctx)
	s.Require().NoError(err)

	s.seedQuestions(ctx)
}

func (s *ServerTestSuite) TearDownSuite() {
	_ = s.app.Stop(context.Background())
}

func (s *ServerTestSuite) TearDownTest() {
	_ = testutil.DeleteData(context.Background(), s.DB)
}

func (s *ServerTestSuite) TestIntegration() {
	ctx := context.Background()

	userID := uuid.New()
	client, err := s.ClientFactory(userID)
	s.Require().NoError(err)

	s.assertRemainingTopicsAre(ctx, client, []oapi.RemainingTopic{
		{Topic: enums.TopicAncientGreece.String(), AmountOfQuestions: 2},
		{Topic: enums.TopicAncientRome.String(), AmountOfQuestions: 1},
	})
	s.assertNextQuestionFor(ctx, client, enums.TopicAncientGreece, s.ancientGreeceQ1)
	s.assertNextQuestionFor(ctx, client, enums.TopicAncientRome, s.ancientRomeQ1)
	s.assertNoNextQuestionFor(ctx, client, enums.TopicNapoleonicWars)

	s.submitAnswer(ctx, client, s.ancientGreeceQ1, s.ancientGreeceQ1.Choices[0].ID)

	s.assertRemainingTopicsAre(ctx, client, []oapi.RemainingTopic{
		{Topic: enums.TopicAncientGreece.String(), AmountOfQuestions: 1},
		{Topic: enums.TopicAncientRome.String(), AmountOfQuestions: 1},
	})
	s.assertNextQuestionFor(ctx, client, enums.TopicAncientGreece, s.ancientGreeceQ2)
	s.assertNextQuestionFor(ctx, client, enums.TopicAncientRome, s.ancientRomeQ1)
	s.assertNoNextQuestionFor(ctx, client, enums.TopicNapoleonicWars)

	s.submitAnswer(ctx, client, s.ancientGreeceQ2, s.ancientGreeceQ2.Choices[1].ID)

	s.assertRemainingTopicsAre(ctx, client, []oapi.RemainingTopic{
		{Topic: enums.TopicAncientRome.String(), AmountOfQuestions: 1},
	})
	s.assertNoNextQuestionFor(ctx, client, enums.TopicAncientGreece)
	s.assertNextQuestionFor(ctx, client, enums.TopicAncientRome, s.ancientRomeQ1)
	s.assertNoNextQuestionFor(ctx, client, enums.TopicNapoleonicWars)

	s.submitAnswer(ctx, client, s.ancientRomeQ1, s.ancientRomeQ1.Choices[0].ID)

	s.assertRemainingTopicsAre(ctx, client, []oapi.RemainingTopic{})
	s.assertNoNextQuestionFor(ctx, client, enums.TopicAncientGreece)
	s.assertNoNextQuestionFor(ctx, client, enums.TopicAncientRome)
	s.assertNoNextQuestionFor(ctx, client, enums.TopicNapoleonicWars)

	s.assertPaginatedLog(
		ctx,
		client,
		userID,
		pagination.Pagination{Page: 0, PageSize: 2},
		[]wantLogItem{
			{QuestionID: s.ancientRomeQ1.ID, ChoiceID: s.ancientRomeQ1.Choices[0].ID},
			{QuestionID: s.ancientGreeceQ2.ID, ChoiceID: s.ancientGreeceQ2.Choices[1].ID},
		},
	)
	s.assertPaginatedLog(
		ctx,
		client,
		userID,
		pagination.Pagination{Page: 1, PageSize: 2},
		[]wantLogItem{
			{QuestionID: s.ancientGreeceQ1.ID, ChoiceID: s.ancientGreeceQ1.Choices[0].ID},
		},
	)
	s.assertPaginatedLog(
		ctx,
		client,
		userID,
		pagination.Pagination{Page: 2, PageSize: 2},
		nil,
	)
}

func (s *ServerTestSuite) assertRemainingTopicsAre(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	wantRemainingTopics []oapi.RemainingTopic,
) {
	s.T().Helper()
	res, err := client.GetRemainingTopics(ctx)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	got := parseResponse[[]oapi.RemainingTopic](s.T(), res)
	s.ElementsMatch(wantRemainingTopics, got)
}

func (s *ServerTestSuite) assertNoNextQuestionFor(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	topic enums.Topic,
) {
	s.T().Helper()
	res, err := client.GetNextQuestion(ctx, &oapi.GetNextQuestionParams{
		Topic: topic.String(),
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, res.StatusCode)
}

func (s *ServerTestSuite) assertNextQuestionFor(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	topic enums.Topic,
	wantQuestion question.Question,
) {
	s.T().Helper()
	res, err := client.GetNextQuestion(ctx, &oapi.GetNextQuestionParams{
		Topic: topic.String(),
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	got := parseResponse[oapi.UnansweredQuestion](s.T(), res)
	s.Equal(wantQuestion.ID.String(), got.Id.String())
	s.Len(got.Choices, len(wantQuestion.Choices))
}

func (s *ServerTestSuite) submitAnswer(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	q question.Question,
	selectedChoice uuid.UUID,
) {
	s.T().Helper()

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

func (s *ServerTestSuite) assertPaginatedLog(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	userID uuid.UUID,
	p pagination.Pagination,
	wantLog []wantLogItem,
) {
	s.T().Helper()

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

func (s *ServerTestSuite) seedQuestions(ctx context.Context) {
	s.ancientGreeceQ1 = question.Mock(func(q *question.Question) {
		q.Question = "Question 1"
		q.Topic = enums.TopicAncientGreece
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, s.ancientGreeceQ1))

	s.ancientGreeceQ2 = question.Mock(func(q *question.Question) {
		q.Question = "Question 2"
		q.Topic = enums.TopicAncientGreece
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, s.ancientGreeceQ2))

	s.ancientRomeQ1 = question.Mock(func(q *question.Question) {
		q.Question = "Question 3"
		q.Topic = enums.TopicAncientRome
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, s.ancientRomeQ1))
}

func parseResponse[T any](t *testing.T, res *http.Response) T {
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var response T
	require.NoError(t, json.Unmarshal(body, &response))
	return response
}
