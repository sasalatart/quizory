package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/generator"
	grpclient "github.com/sasalatart/quizory/http/grpc/client"
	"github.com/sasalatart/quizory/infra"
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
		generator.Module,
		grpclient.Module,
		fx.Invoke(generatorLC),
	)

	infra.RunFXApp(ctx, app)
}

func generatorLC(lc fx.Lifecycle, s *generator.Service, cfg config.LLMConfig) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := s.GenerateBatch(ctx, cfg.Questions.BatchSize, enums.RandomTopic())
				if err != nil {
					slog.Error("Failed to generate questions", slog.Any("error", err))
					os.Exit(1)
				}
				os.Exit(0)
			}()
			return nil
		},
	})
}
