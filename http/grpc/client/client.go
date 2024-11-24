package client

import (
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/http/grpc/proto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func new(cfg config.ServerConfig) (proto.QuizoryServiceClient, *grpc.ClientConn, error) {
	var client proto.QuizoryServiceClient

	statsHandler := otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
	conn, err := grpc.NewClient(
		cfg.GRPCAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(statsHandler),
	)
	if err != nil {
		return client, conn, err
	}

	return proto.NewQuizoryServiceClient(conn), conn, nil
}
