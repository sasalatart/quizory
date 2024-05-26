package question

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/domain/question/enums"
)

var ErrInvalidRecord = errors.New("invalid question record")

// Question represents a question about history that users need to answer.
type Question struct {
	ID         uuid.UUID
	Topic      enums.Topic
	Question   string
	Hint       string
	MoreInfo   string
	Difficulty enums.Difficulty
	Choices    []Choice
	CreatedAt  time.Time
}

// Choice represents a possible answer to a question.
type Choice struct {
	ID        uuid.UUID
	Choice    string
	IsCorrect bool
}

// New creates a new Question with a random ID, a question, and a hint.
func New(question, hint, moreInfo string) *Question {
	return &Question{
		ID:        uuid.New(),
		Question:  question,
		Hint:      hint,
		MoreInfo:  moreInfo,
		CreatedAt: time.Now(),
	}
}

// WithTopic sets the topic of a question.
func (q *Question) WithTopic(topic enums.Topic) *Question {
	q.Topic = topic
	return q
}

// WithDifficulty sets the difficulty of a question.
func (q *Question) WithDifficulty(difficulty enums.Difficulty) *Question {
	q.Difficulty = difficulty
	return q
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
	invalidErr := func(msg string) error {
		return errors.Wrap(ErrInvalidRecord, msg)
	}

	if !q.Topic.IsATopic() {
		return invalidErr("invalid topic")
	}
	if q.Question == "" {
		return invalidErr("question is required")
	}
	if q.Hint == "" {
		return invalidErr("hint is required")
	}
	if len(q.Choices) < 2 {
		return invalidErr("at least two choices are required")
	}
	if q.MoreInfo == "" {
		return invalidErr("more info field is required")
	}
	if !q.Difficulty.IsADifficulty() {
		return invalidErr("invalid difficulty")
	}

	seenChoices := make(map[string]struct{})
	correctChoices := 0
	for _, c := range q.Choices {
		seenChoices[c.Choice] = struct{}{}
		if c.IsCorrect {
			correctChoices++
		}
	}
	if len(seenChoices) != len(q.Choices) {
		return invalidErr("choices must be unique")
	}
	if correctChoices != 1 {
		return invalidErr("one choice must be correct")
	}
	return nil
}

// CorrectChoices returns the choices that are correct.
func (q *Question) CorrectChoice() (*Choice, error) {
	for _, c := range q.Choices {
		if c.IsCorrect {
			return &c, nil
		}
	}
	return nil, errors.New("no correct choice found")
}
