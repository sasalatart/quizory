package main

import (
	"log"
	"log/slog"

	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db/migrations"
)

func main() {
	slog.Info("Running migrations...")
	defer slog.Info("Migrations complete.")

	dbCfg := config.NewConfig().DB
	if err := migrations.Up(dbCfg); err != nil {
		log.Fatal(err)
	}
}
