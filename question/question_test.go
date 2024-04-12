package question_test

import (
	"testing"

	"github.com/sasalatart.com/quizory/question"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	q := question.New("Test Question", "Test Hint")
	assert.Equal(t, "Test Question", q.Question)
	assert.Equal(t, "Test Hint", q.Hint)
}

func TestQuestion_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		q       question.Question
		wantErr bool
	}{
		{
			name:    "Valid",
			q:       question.New("Test Question", "Test Hint"),
			wantErr: false,
		},
		{
			name:    "Without Question",
			q:       question.New("", "Test Hint"),
			wantErr: true,
		},
		{
			name:    "Without Hint",
			q:       question.New("Test Question", ""),
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.q.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
