package servertest

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/http/server"
	"github.com/sasalatart/quizory/infra"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"test-server",
	fx.Provide(server.NewServer),
	fx.Provide(newClientFactory),
	fx.Invoke(serverLC),
)

// serverLC starts the server and a test client, and waits for the server to be ready before
// returning. It is intended to be used in test suites to ensure that the server is ready before
// running tests.
func serverLC(lc fx.Lifecycle, server *server.Server, clientFactory ClientFactory) {
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
