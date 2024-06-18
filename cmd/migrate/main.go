package main

import (
	"log"
	"log/slog"

	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db/migrations"
)

func main() {
	slog.Info("Running migrations...")
	defer slog.Info("Migrations complete.")

	dbCfg := config.NewConfig().DB
	if err := migrations.Up(dbCfg); err != nil {
		log.Fatal(err)
	}
}
