package client

import (
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/http/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func new(cfg config.ServerConfig) (proto.QuizoryServiceClient, *grpc.ClientConn, error) {
	var client proto.QuizoryServiceClient

	conn, err := grpc.NewClient(
		cfg.GRPCAddress(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return client, conn, err
	}

	return proto.NewQuizoryServiceClient(conn), conn, nil
}
