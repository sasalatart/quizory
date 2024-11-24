package rest

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"rest-server",
	fx.Provide(NewServer),
	fx.Invoke(serverLC),
)

func serverLC(lc fx.Lifecycle, shutdowner fx.Shutdowner, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := s.Start(); err != nil {
					slog.Error("Error starting REST server", slog.Any("error", err))
					_ = shutdowner.Shutdown(fx.ExitCode(1))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
}
