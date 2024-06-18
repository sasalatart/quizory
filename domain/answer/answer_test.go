package answer_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/sasalatart/quizory/domain/answer"
	"github.com/stretchr/testify/assert"
)

func TestAnswer_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		answer  *answer.Answer
		wantErr bool
	}{
		{
			name:    "valid answer",
			answer:  answer.New(uuid.New(), uuid.New()),
			wantErr: false,
		},
		{
			name:    "missing user ID",
			answer:  answer.New(uuid.Nil, uuid.New()),
			wantErr: true,
		},
		{
			name:    "missing choice ID",
			answer:  answer.New(uuid.New(), uuid.Nil),
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.answer.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, answer.ErrInvalidRecord)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
