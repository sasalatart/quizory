package grpc

import (
	"context"
	"log/slog"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc-server",
	fx.Provide(NewServer),
	fx.Invoke(serverLC),
)

func serverLC(lc fx.Lifecycle, shutdowner fx.Shutdowner, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := s.Start(); err != nil {
					slog.Error("Error starting gRPC server", slog.Any("error", err))
					_ = shutdowner.Shutdown(fx.ExitCode(1))
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			s.Shutdown()
			return nil
		},
	})
}
