package answer

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrInvalidRecord = errors.New("invalid answer record")

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
		return errors.Wrap(ErrInvalidRecord, "user ID is required")
	}
	if a.ChoiceID == uuid.Nil {
		return errors.Wrap(ErrInvalidRecord, "choice ID is required")
	}
	return nil
}
