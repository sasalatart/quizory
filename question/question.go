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
	CreatedAt time.Time
}

// New creates a new Question with a random ID, a question, and a hint.
func New(question, hint string) Question {
	return Question{
		ID:        uuid.New(),
		Question:  question,
		Hint:      hint,
		CreatedAt: time.Now(),
	}
}

// Validate checks if a question is valid or not, and returns an error if it's not.
func (q Question) Validate() error {
	if q.Question == "" {
		return errors.New("question is required")
	}
	if q.Hint == "" {
		return errors.New("hint is required")
	}
	return nil
}
