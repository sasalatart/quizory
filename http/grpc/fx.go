package grpc

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"grpc-server",
	fx.Provide(NewServer),
	fx.Invoke(serverLC),
)

func serverLC(lc fx.Lifecycle, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.Start()
			return nil
		},
		OnStop: func(context.Context) error {
			s.Shutdown()
			return nil
		},
	})
}
