package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db"
	"github.com/sasalatart.com/quizory/domain/question"
	"github.com/sasalatart.com/quizory/llm"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(config.NewConfig),
		db.Module,
		llm.Module,
		question.Module,
		fx.Invoke(func(lc fx.Lifecycle, llmCfg config.LLMConfig, service *question.Service) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					service.StartGeneration(ctx, llmCfg.Frequency, llmCfg.BatchSize)
					return nil
				},
			})
		}),
	)

	ctx := context.Background()

	go func() {
		if err := app.Start(ctx); err != nil {
			panic(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	if err := app.Stop(ctx); err != nil {
		panic(err)
	}
}
