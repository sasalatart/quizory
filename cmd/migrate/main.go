package main

import (
	"log"
	"log/slog"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sasalatart.com/quizory/db/migrations"
)

func main() {
	slog.Info("Running migrations...")
	defer slog.Info("Migrations complete.")

	if err := migrations.Up(os.Getenv("DB_URL")); err != nil {
		log.Fatal(err)
	}
}
