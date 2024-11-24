package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/generator"
	grpclient "github.com/sasalatart/quizory/http/grpc/client"
	"github.com/sasalatart/quizory/llm"
	"go.uber.org/fx"

	"github.com/sasalatart/quizory/infra/otel"
)

func main() {
	ctx := context.Background()

	app := fx.New(
		fx.Provide(config.NewConfig),
		otel.Module,
		llm.Module,
		grpclient.Module,
		generator.Module,
		fx.Invoke(generatorLC),
	)

	if err := app.Start(ctx); err != nil {
		slog.Error("Error starting app", slog.Any("error", err))
		os.Exit(1)
	}
	if err := app.Stop(ctx); err != nil {
		slog.Error("Error stopping app", slog.Any("error", err))
		os.Exit(1)
	}
}

func generatorLC(lc fx.Lifecycle, s *generator.Service, cfg config.LLMConfig) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
			defer cancel()

			err := s.GenerateBatch(ctx, cfg.Questions.BatchSize, enums.RandomTopic())
			if err != nil {
				slog.Error("Failed to generate questions", slog.Any("error", err))
				return err
			}

			slog.Info("Questions generated successfully")
			return nil
		},
	})
}
