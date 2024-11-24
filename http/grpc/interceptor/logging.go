package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryLogging logs the method of a gRPC request and its duration for unary handlers.
func UnaryLogging(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now()

	slog.Info("Handling request", slog.String("method", info.FullMethod))
	defer func() {
		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}
		slog.Info(
			"Completed request",
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
			slog.String("code", code.String()),
		)
	}()

	resp, err = handler(ctx, req)
	return resp, err
}
