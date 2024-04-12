package question_test

import (
	"testing"

	"github.com/sasalatart.com/quizory/question"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	q := question.
		New("Test Question", "Test Hint").
		WithChoice("Choice 1", false).
		WithChoice("Choice 2", true)

	assert.Equal(t, "Test Question", q.Question)
	assert.Equal(t, "Test Hint", q.Hint)

	assert.Len(t, q.Choices, 2)

	assert.Equal(t, "Choice 1", q.Choices[0].Choice)
	assert.False(t, q.Choices[0].IsCorrect)

	assert.Equal(t, "Choice 2", q.Choices[1].Choice)
	assert.True(t, q.Choices[1].IsCorrect)
}

func TestQuestion_Validate(t *testing.T) {
	validQuestion := *question.
		New("Test Question", "Test Hint").
		WithChoice("Choice 1", false).
		WithChoice("Choice 2", true)

	testCases := []struct {
		name    string
		factory func() question.Question
		wantErr bool
	}{
		{
			name: "Valid",
			factory: func() question.Question {
				return validQuestion
			},
			wantErr: false,
		},
		{
			name: "Without Question",
			factory: func() question.Question {
				q := validQuestion
				q.Question = ""
				return q
			},
			wantErr: true,
		},
		{
			name: "Without Hint",
			factory: func() question.Question {
				q := validQuestion
				q.Hint = ""
				return q
			},
			wantErr: true,
		},
		{
			name: "Without Choices",
			factory: func() question.Question {
				q := validQuestion
				q.Choices = []question.Choice{}
				return q
			},
			wantErr: true,
		},
		{
			name: "Without Correct Choices",
			factory: func() question.Question {
				q := validQuestion
				q.Choices = []question.Choice{}
				q.WithChoice("Choice 1", false).WithChoice("Choice 2", false)
				return q
			},
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			q := tt.factory()
			err := q.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
