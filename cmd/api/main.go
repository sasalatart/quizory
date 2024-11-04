package main

import (
	"context"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/http/grpc"
	"github.com/sasalatart/quizory/http/rest"
	"github.com/sasalatart/quizory/infra"
	"go.uber.org/fx"

	"github.com/sasalatart/quizory/infra/otel"
)

func main() {
	ctx := context.Background()

	app := fx.New(
		fx.Provide(config.NewConfig),
		otel.Module,
		db.Module,
		answer.Module,
		question.Module,
		grpc.Module,
		rest.Module, // TODO: make this module's Invoke non-blocking
	)

	infra.RunFXApp(ctx, app)
}
