package rest

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"rest-server",
	fx.Provide(NewServer),
	fx.Invoke(serverLC),
)

func serverLC(lc fx.Lifecycle, s *Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go s.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
}
