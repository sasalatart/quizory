package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/sasalatart.com/quizory/question"
	"github.com/sashabaranov/go-openai"
)

func main() {
	ctx := context.Background()

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	questionRepo := question.NewRepository(db)
	questionService := question.NewService(questionRepo, openai.NewClient(os.Getenv("OPENAI_API_KEY")))

	cancel := make(chan struct{})
	go questionService.StartGeneration(ctx, 5*time.Second, 5, cancel)

	<-time.After(4 * time.Minute)
	close(cancel)
}
