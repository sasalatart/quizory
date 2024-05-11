package answer

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Answer represents the choice selected by a user for a question.
type Answer struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ChoiceID  uuid.UUID
	CreatedAt time.Time
}

// New creates a new Answer with a random ID, and with the specified user ID and choice ID.
func New(userID, choiceID uuid.UUID) *Answer {
	return &Answer{
		ID:        uuid.New(),
		UserID:    userID,
		ChoiceID:  choiceID,
		CreatedAt: time.Now(),
	}
}

// Validate checks if the answer is valid.
func (a Answer) Validate() error {
	if a.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}
	if a.ChoiceID == uuid.Nil {
		return errors.New("choice ID is required")
	}
	return nil
}
