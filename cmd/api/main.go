package main

import (
	"context"
	"log/slog"
	"os"

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

	if err := app.Start(ctx); err != nil {
		slog.Error("Error starting app", slog.Any("error", err))
		os.Exit(1)
	}

	// Wait for app to finish (handles SIGINT and SIGTERM)
	sig := <-app.Wait()

	slog.Info("Shutting down app")
	if err := app.Stop(ctx); err != nil {
		slog.Error("Error stopping app", slog.Any("error", err))
		os.Exit(1)
	}

	os.Exit(sig.ExitCode)
}
