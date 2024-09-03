package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db"
	"github.com/sasalatart/quizory/db/migrations"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/http/server"
	"github.com/sasalatart/quizory/llm"
	"go.uber.org/fx"

	"github.com/sasalatart/quizory/infra/otel"
)

var generateQuestions bool

func init() {
	flag.BoolVar(&generateQuestions, "generate", false, "generate questions")
	flag.Parse()
}

func main() {
	app := fx.New(
		fx.Provide(config.NewConfig),
		db.Module,
		llm.Module,
		answer.Module,
		question.Module,
		server.Module,
		otel.Module,
		fx.Invoke(migrationsLC),
		fx.Invoke(questionsGenLC),
		fx.Invoke(serverLC),
	)

	ctx := context.Background()

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

func migrationsLC(lc fx.Lifecycle, dbCfg config.DBConfig) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return migrations.Up(dbCfg)
		},
	})
}

func questionsGenLC(lc fx.Lifecycle, llmCfg config.LLMConfig, service *question.Service) {
	if !generateQuestions {
		return
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go service.StartGeneration(ctx, llmCfg.Frequency, llmCfg.BatchSize)
			return nil
		},
	})
}

func serverLC(lc fx.Lifecycle, s *server.Server) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			s.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return s.Shutdown(ctx)
		},
	})
}
