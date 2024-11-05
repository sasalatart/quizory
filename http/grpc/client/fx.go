package client

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var Module = fx.Module(
	"grpc-client",
	fx.Provide(new),
	fx.Invoke(clientLC),
)

func clientLC(lc fx.Lifecycle, conn *grpc.ClientConn) {
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return conn.Close()
		},
	})
}
