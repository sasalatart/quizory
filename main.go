package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db"
	"github.com/sasalatart.com/quizory/llm"
	"github.com/sasalatart.com/quizory/question"
	"go.uber.org/fx"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	app := fx.New(
		fx.Provide(config.NewConfig),
		db.Module,
		llm.Module,
		question.Module,
		fx.Invoke(func(lc fx.Lifecycle, cfg config.Config, service *question.Service) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					freq := cfg.QuestionGeneration.Frequency
					batchSize := cfg.QuestionGeneration.BatchSize
					service.StartGeneration(ctx, freq, batchSize)
					return nil
				},
				OnStop: func(context.Context) error {
					cancel()
					return nil
				},
			})
		}),
	)

	go func() {
		if err := app.Start(ctx); err != nil {
			panic(err)
		}
	}()

	handleStop := func() {
		if err := app.Stop(ctx); err != nil {
			panic(err)
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sig:
		handleStop()
	case <-time.After(4 * time.Minute):
		handleStop()
	}
}
