package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/http/grpc"
	"github.com/sasalatart/quizory/http/rest"
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
		rest.Module,
	)

	go func() {
		if err := app.Start(ctx); err != nil {
			slog.Error("Error starting app", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	if err := app.Stop(ctx); err != nil {
		slog.Error("Error stopping app", slog.Any("error", err))
		os.Exit(1)
	}
}
