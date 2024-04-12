package question

import (
	"github.com/sasalatart.com/quizory/question"
	"github.com/sasalatart.com/quizory/question/internal/repo"
)

type Service struct {
	repo repo.QuestionRepo
}

// TODO: implement with AI
func (s Service) GenerateMany() ([]question.Question, error) {
	return nil, nil
}
