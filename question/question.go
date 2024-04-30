package question

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sasalatart.com/quizory/question/enums"
)

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
	if !q.Topic.IsATopic() {
		return errors.New("invalid topic")
	}
	if q.Question == "" {
		return errors.New("question is required")
	}
	if q.Hint == "" {
		return errors.New("hint is required")
	}
	if len(q.Choices) < 2 {
		return errors.New("at least two choices are required")
	}
	if q.MoreInfo == "" {
		return errors.New("more info field is required")
	}
	if !q.Difficulty.IsADifficulty() {
		return errors.New("invalid difficulty")
	}

	seenChoices := make(map[string]struct{})
	hasAtLeastOneCorrectChoice := false
	for _, c := range q.Choices {
		seenChoices[c.Choice] = struct{}{}
		if c.IsCorrect {
			hasAtLeastOneCorrectChoice = true
		}
	}
	if len(seenChoices) != len(q.Choices) {
		return errors.New("choices must be unique")
	}
	if !hasAtLeastOneCorrectChoice {
		return errors.New("at least one choice must be correct")
	}
	return nil
}

// CorrectChoices returns the choices that are correct.
func (q *Question) CorrectChoices() []Choice {
	var choices []Choice
	for _, c := range q.Choices {
		if c.IsCorrect {
			choices = append(choices, c)
		}
	}
	return choices
}
