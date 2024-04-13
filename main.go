package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sasalatart.com/quizory/question/repo"
	"github.com/sasalatart.com/quizory/question/service"
	"github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	questionRepo := repo.New(db)
	questionService := service.New(questionRepo, openai.NewClient(os.Getenv("OPENAI_API_KEY")))

	cancel := make(chan struct{})
	go questionService.StartGeneration(ctx, 5*time.Second, 5, cancel)

	<-time.After(4 * time.Minute)
	close(cancel)
}
