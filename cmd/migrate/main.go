package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sasalatart.com/quizory/db/migrations"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if err := migrations.Up(os.Getenv("DB_URL")); err != nil {
		log.Fatal(err)
	}
}
