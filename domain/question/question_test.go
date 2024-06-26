package question_test

import (
	"testing"

	"github.com/sasalatart/quizory/domain/question"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	q := question.
		New("Test Question", "Test Hint", "Test More Info").
		WithTopic(enums.TopicAncientRome).
		WithDifficulty(enums.DifficultyAvidHistorian).
		WithChoice("Choice 1", false).
		WithChoice("Choice 2", true)

	assert.Equal(t, enums.TopicAncientRome, q.Topic)
	assert.Equal(t, "Test Question", q.Question)
	assert.Equal(t, "Test Hint", q.Hint)
	assert.Equal(t, "Test More Info", q.MoreInfo)
	assert.Equal(t, enums.DifficultyAvidHistorian, q.Difficulty)

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
				q.Choices = nil
				return q
			},
			wantErr: true,
		},
		{
			name: "With Duplicate Choices",
			factory: func() question.Question {
				q := validQuestion
				q.Choices = nil
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
				q.Choices = nil
				q.WithChoice("Choice 1", false).WithChoice("Choice 2", false)
				return q
			},
			wantErr: true,
		},
		{
			name: "With More Than One Correct Choice",
			factory: func() question.Question {
				q := validQuestion
				q.Choices = nil
				q.WithChoice("Choice 1", true).
					WithChoice("Choice 2", true).
					WithChoice("Choice 3", false)
				return q
			},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			q := tc.factory()
			err := q.Validate()
			if tc.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, question.ErrInvalidRecord)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQuestion_CorrectChoice(t *testing.T) {
	testCases := []struct {
		name    string
		q       question.Question
		want    string
		wantErr bool
	}{
		{
			name: "With Two Choices",
			q: question.Mock(func(q *question.Question) {
				q.Choices = nil
				q.WithChoice("Choice 1", false).WithChoice("Choice 2", true)
			}),
			want: "Choice 2",
		},
		{
			name: "With Four Choices",
			q: question.Mock(func(q *question.Question) {
				q.Choices = nil
				q.WithChoice("Choice 1", false).
					WithChoice("Choice 2", false).
					WithChoice("Choice 3", true).
					WithChoice("Choice 4", false)
			}),
			want: "Choice 3",
		},
		{
			name: "Without Correct Choice",
			q: question.Mock(func(q *question.Question) {
				q.Choices = nil
				q.WithChoice("Choice 1", false).WithChoice("Choice 2", false)
			}),
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := tc.q.CorrectChoice()
			if tc.wantErr {
				require.Error(t, err)
				require.Nil(t, c)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, c.Choice)
		})
	}
}
