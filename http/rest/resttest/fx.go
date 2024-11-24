package resttest

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/http/rest"
	"github.com/sasalatart/quizory/infra"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

var Module = fx.Module(
	"test-server",
	fx.Provide(rest.NewServer),
	fx.Provide(newClientFactory),
	fx.Invoke(serverLC),
)

// serverLC starts the server and a test client, and waits for the server to be ready before
// returning. It is intended to be used in test suites to ensure that the server is ready before
// running tests.
func serverLC(lc fx.Lifecycle, server *rest.Server, clientFactory ClientFactory) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start(); err != nil {
					slog.Error("Failed to start server", slog.Any("error", err))
				}
			}()

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
