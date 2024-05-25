package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db"
	"github.com/sasalatart.com/quizory/db/migrations"
	"github.com/sasalatart.com/quizory/domain/answer"
	"github.com/sasalatart.com/quizory/domain/question"
	"github.com/sasalatart.com/quizory/http/server"
	"github.com/sasalatart.com/quizory/llm"
	"go.uber.org/fx"
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
		fx.Invoke(migrationsLC),
		fx.Invoke(questionsGenLC),
		fx.Invoke(serverLC),
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
		OnStop: func(context.Context) error {
			return s.Shutdown()
		},
	})
}
