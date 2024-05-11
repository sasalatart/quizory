package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db"
	"github.com/sasalatart.com/quizory/domain/answer"
	"github.com/sasalatart.com/quizory/domain/question"
	"github.com/sasalatart.com/quizory/http/server"
	"github.com/sasalatart.com/quizory/llm"
	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()

	app := fx.New(
		fx.Provide(config.NewConfig),
		db.Module,
		llm.Module,
		answer.Module,
		question.Module,
		server.Module,
		fx.Invoke(func(lc fx.Lifecycle, s *server.Server) {
			lc.Append(fx.Hook{
				OnStart: func(context.Context) error {
					s.Start()
					return nil
				},
				OnStop: func(context.Context) error {
					return s.Shutdown()
				},
			})
		}),
	)

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
