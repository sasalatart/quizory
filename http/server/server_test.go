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
	"github.com/sasalatart.com/quizory/domain/pagination"
	"github.com/sasalatart.com/quizory/domain/question"
	"github.com/sasalatart.com/quizory/http/oapi"
	"github.com/sasalatart.com/quizory/http/server"
	"github.com/sasalatart.com/quizory/testutil"
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

	q1 question.Question
	q2 question.Question
	q3 question.Question
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.NopLogger,
		testutil.ModuleWithAPI,
		fx.Populate(&s.serverTestSuiteParams),
	)
	err := s.app.Start(context.Background())
	s.Require().NoError(err)
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

	s.mustNotHaveNextQuestion(ctx, client)

	s.seedQuestions(ctx)

	s.mustGetNextQuestion(ctx, client, s.q1)
	s.mustSubmitAnswer(ctx, client, s.q1, s.q1.Choices[0].ID)

	s.mustGetNextQuestion(ctx, client, s.q2)
	s.mustSubmitAnswer(ctx, client, s.q2, s.q2.Choices[1].ID)

	s.mustGetNextQuestion(ctx, client, s.q3)
	s.mustSubmitAnswer(ctx, client, s.q3, s.q3.Choices[0].ID)

	s.mustNotHaveNextQuestion(ctx, client)

	s.mustGetLog(
		ctx,
		client,
		userID,
		pagination.Pagination{Page: 0, PageSize: 2},
		[]wantLogItem{
			{QuestionID: s.q3.ID, ChoiceID: s.q3.Choices[0].ID},
			{QuestionID: s.q2.ID, ChoiceID: s.q2.Choices[1].ID},
		},
	)
	s.mustGetLog(
		ctx,
		client,
		userID,
		pagination.Pagination{Page: 1, PageSize: 2},
		[]wantLogItem{
			{QuestionID: s.q1.ID, ChoiceID: s.q1.Choices[0].ID},
		},
	)
	s.mustGetLog(
		ctx,
		client,
		userID,
		pagination.Pagination{Page: 2, PageSize: 2},
		nil,
	)
}

func (s *ServerTestSuite) mustNotHaveNextQuestion(
	ctx context.Context,
	client *oapi.ClientWithResponses,
) {
	s.T().Helper()
	res, err := client.GetNextQuestion(ctx)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, res.StatusCode)
}

func (s *ServerTestSuite) mustGetNextQuestion(
	ctx context.Context,
	client *oapi.ClientWithResponses,
	wantQuestion question.Question,
) {
	s.T().Helper()
	res, err := client.GetNextQuestion(ctx)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	got := parseResponse[oapi.UnansweredQuestion](s.T(), res)
	s.Equal(wantQuestion.ID.String(), got.Id.String())
	s.Len(got.Choices, len(wantQuestion.Choices))
}

func (s *ServerTestSuite) mustSubmitAnswer(
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

func (s *ServerTestSuite) mustGetLog(
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
	s.q1 = question.Mock(func(q *question.Question) {
		q.Question = "Question 1"
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, s.q1))

	s.q2 = question.Mock(func(q *question.Question) {
		q.Question = "Question 2"
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, s.q2))

	s.q3 = question.Mock(func(q *question.Question) {
		q.Question = "Question 3"
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, s.q3))
}

func parseResponse[T any](t *testing.T, res *http.Response) T {
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var response T
	require.NoError(t, json.Unmarshal(body, &response))
	return response
}
