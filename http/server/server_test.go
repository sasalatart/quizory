package server_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"testing"

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
	res, err := s.Client.GetQuestionsNext(ctx)
	s.Require().NoError(err)
	s.Equal(http.StatusNoContent, res.StatusCode)

	q1 := question.Mock(func(q *question.Question) {
		q.Question = "When was Napoleon born?"
		q.Choices = nil
		q.WithChoice("1769", true).WithChoice("1770", false)
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, q1))

	q2 := question.Mock(func(q *question.Question) {
		q.Question = "When did Napoleon die?"
		q.Choices = nil
		q.WithChoice("1821", true).WithChoice("1870", false)
	})
	s.Require().NoError(s.QuestionRepo.Insert(ctx, q2))

	res, err = s.Client.GetQuestionsNext(ctx)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, res.StatusCode)
	got := parseResponse[oapi.UnansweredQuestion](s.T(), res)
	s.Equal(q1.ID.String(), got.Id.String())
}

func parseResponse[T any](t *testing.T, res *http.Response) T {
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	var response T
	require.NoError(t, json.Unmarshal(body, &response))
	return response
}
