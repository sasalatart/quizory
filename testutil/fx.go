package testutil

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/http/server"
	"github.com/sasalatart/quizory/infra"
	"github.com/sasalatart/quizory/infra/otel"
	"github.com/sasalatart/quizory/llm"
	"go.uber.org/fx"
)

// Module defines a reusable module so that we do not need to manually provide all the dependencies
// in every test suite. It also provides test-specific defaults.
// Module DOES NOT include server.TestModule. Use ModuleWithAPI for that instead.
var Module = fx.Module(
	"testutil",

	fx.Provide(config.NewTestConfig),
	db.TestModule,
	llm.TestModule,

	answer.Module,
	question.Module,
	otel.TestModule,

	// Repositories are injected privately in the modules above, so we provide them here to make
	// them available for tests (e.g. for seeding the database with test data).
	fx.Provide(answer.NewRepository),
	fx.Provide(question.NewRepository),
)

// ModuleWithAPI injects the API module in addition to the dependencies provided by Module.
// It is intended to be used in test suites that require the API module, as it also turns on the API
// server, manages its lifecycle, and waits for it to be ready before running tests.
var ModuleWithAPI = fx.Module(
	"testutil-with-api",

	Module,
	server.TestModule,
	fx.Invoke(serverLC),
)

// serverLC starts the server and a test client, and waits for the server to be ready before
// returning. It is intended to be used in test suites to ensure that the server is ready before
// running tests.
func serverLC(lc fx.Lifecycle, server *server.Server, clientFactory server.TestClientFactory) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go server.Start()

			client, err := clientFactory(uuid.New())
			if err != nil {
				return errors.Wrap(err, "creating test client")
			}

			check := func() bool {
				resp, err := client.HealthCheck(ctx)
				return err == nil && resp.StatusCode == http.StatusNoContent
			}
			if err := infra.WaitFor(check, 5, 1*time.Second); err != nil {
				return errors.Wrap(err, "waiting for server to start")
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
