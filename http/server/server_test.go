package server_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart.com/quizory/http/oapi"
	"github.com/sasalatart.com/quizory/http/server"
	"github.com/sasalatart.com/quizory/question"
	"github.com/sasalatart.com/quizory/testutil"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type serverTestSuiteParams struct {
	fx.In
	DB           *sql.DB
	Client       *oapi.ClientWithResponses
	QuestionRepo *question.Repository
}

type ServerTestSuite struct {
	suite.Suite
	serverTestSuiteParams
	app *fx.App

	q1 question.Question
	q2 question.Question
}

func TestServer(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	s.app = fx.New(
		fx.NopLogger,
		testutil.Module,
		fx.Populate(&s.serverTestSuiteParams),
		fx.Invoke(func(lc fx.Lifecycle, server *server.Server) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					go server.Start()
					return nil
				},
				OnStop: func(context.Context) error {
					return server.Shutdown()
				},
			})
		}),
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

	s.mustNotHaveNextQuestion(ctx)

	s.seedQuestions(ctx)

	s.mustGetNextQuestion(ctx, s.q1)
	s.mustSubmitAnswer(ctx, s.q1, s.q1.Choices[0].ID)

	s.mustGetNextQuestion(ctx, s.q2)
	s.mustSubmitAnswer(ctx, s.q2, s.q2.Choices[1].ID)

	s.mustNotHaveNextQuestion(ctx)
}

func (s *ServerTestSuite) mustNotHaveNextQuestion(ctx context.Context) {
	s.T().Helper()
	res, err := s.Client.GetNextQuestion(ctx)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNoContent, res.StatusCode)
}

func (s *ServerTestSuite) mustGetNextQuestion(ctx context.Context, wantQuestion question.Question) {
	s.T().Helper()
	res, err := s.Client.GetNextQuestion(ctx)
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, res.StatusCode)
	got := parseResponse[oapi.UnansweredQuestion](s.T(), res)
	s.Equal(wantQuestion.ID.String(), got.Id.String())
}

func (s *ServerTestSuite) mustSubmitAnswer(
	ctx context.Context,
	q question.Question,
	selectedChoice uuid.UUID,
) {
	s.T().Helper()

	correctChoice, err := q.CorrectChoice()
	s.Require().NoError(err)

	res, err := s.Client.SubmitAnswer(ctx, oapi.SubmitAnswerJSONRequestBody{
		ChoiceId: selectedChoice,
	})
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, res.StatusCode)
	got := parseResponse[oapi.SubmitAnswerResult](s.T(), res)
	s.Equal(correctChoice.ID.String(), got.CorrectChoiceId.String())
	s.Equal(q.MoreInfo, got.MoreInfo)
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
}

func parseResponse[T any](t *testing.T, res *http.Response) T {
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var response T
	require.NoError(t, json.Unmarshal(body, &response))
	return response
}
