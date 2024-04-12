package question

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Question represents a question about history that users need to answer.
type Question struct {
	ID        uuid.UUID
	Question  string
	Hint      string
	Choices   []Choice
	CreatedAt time.Time
}

// Choice represents a possible answer to a question.
type Choice struct {
	ID        uuid.UUID
	Choice    string
	IsCorrect bool
}

// New creates a new Question with a random ID, a question, and a hint.
func New(question, hint string) *Question {
	return &Question{
		ID:        uuid.New(),
		Question:  question,
		Hint:      hint,
		CreatedAt: time.Now(),
	}
}

// WithChoice adds a choice to a question.
func (q *Question) WithChoice(choice string, isCorrect bool) *Question {
	q.Choices = append(q.Choices, Choice{
		ID:        uuid.New(),
		Choice:    choice,
		IsCorrect: isCorrect,
	})
	return q
}

// Validate checks if a question is valid or not, and returns an error if it's not.
func (q *Question) Validate() error {
	if q.Question == "" {
		return errors.New("question is required")
	}
	if q.Hint == "" {
		return errors.New("hint is required")
	}
	if len(q.Choices) < 2 {
		return errors.New("at least two choices are required")
	}

	hasAtLeastOneCorrectChoice := false
	for _, c := range q.Choices {
		if c.IsCorrect {
			hasAtLeastOneCorrectChoice = true
			break
		}
	}
	if !hasAtLeastOneCorrectChoice {
		return errors.New("at least one choice must be correct")
	}
	return nil
}
