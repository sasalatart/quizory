package main

import (
	"log/slog"

	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db/migrations"
)

func main() {
	slog.Info("Running migrations...")
	defer slog.Info("Migrations complete.")

	dbCfg := config.NewConfig().Database
	if err := migrations.Up(dbCfg); err != nil {
		slog.Error("Error running migrations", err)
	}
}
