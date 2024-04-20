package question_test

import (
	"testing"

	"github.com/sasalatart.com/quizory/question"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	q := question.
		New("Test Question", "Test Hint", "Test More Info").
		WithTopic(question.TopicAncientRome).
		WithDifficulty(question.DifficultyAvidHistorian).
		WithChoice("Choice 1", false).
		WithChoice("Choice 2", true)

	assert.Equal(t, question.TopicAncientRome, q.Topic)
	assert.Equal(t, "Test Question", q.Question)
	assert.Equal(t, "Test Hint", q.Hint)
	assert.Equal(t, "Test More Info", q.MoreInfo)
	assert.Equal(t, question.DifficultyAvidHistorian, q.Difficulty)

	assert.Len(t, q.Choices, 2)

	assert.Equal(t, "Choice 1", q.Choices[0].Choice)
	assert.False(t, q.Choices[0].IsCorrect)

	assert.Equal(t, "Choice 2", q.Choices[1].Choice)
	assert.True(t, q.Choices[1].IsCorrect)
}

func TestQuestion_Validate(t *testing.T) {
	validQuestion := question.Mock(nil)

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
			name: "With Invalid Topic",
			factory: func() question.Question {
				q := validQuestion
				q.WithTopic(-1)
				return q
			},
			wantErr: true,
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
			name: "Without MoreInfo",
			factory: func() question.Question {
				q := validQuestion
				q.MoreInfo = ""
				return q
			},
			wantErr: true,
		},
		{
			name: "With Invalid Difficulty",
			factory: func() question.Question {
				q := validQuestion
				q.WithDifficulty(-1)
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
			name: "With Duplicate Choices",
			factory: func() question.Question {
				q := validQuestion
				q.Choices = []question.Choice{}
				q.WithChoice("Choice 1", true).
					WithChoice("Choice 1", false).
					WithChoice("Choice 2", false)
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
