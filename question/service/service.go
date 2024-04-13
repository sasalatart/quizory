package service

import (
	"github.com/sasalatart.com/quizory/question/repo"
	"github.com/sashabaranov/go-openai"
)

type Service struct {
	repo         *repo.QuestionRepo
	openaiClient *openai.Client
}

// New creates a new instance of Service.
func New(repo *repo.QuestionRepo, openaiClient *openai.Client) Service {
	return Service{repo: repo, openaiClient: openaiClient}
}
