package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/question"
	"github.com/sashabaranov/go-openai"

	_ "embed"
)

// StartGeneration starts a loop that generates questions about random topics at a given frequency.
func (s Service) StartGeneration(
	ctx context.Context,
	freq time.Duration,
	amountPerBatch int,
	cancel <-chan struct{},
) error {
	slog.Info("Starting generation loop", slog.Duration("freq", freq))
	ticker := time.NewTicker(freq)
	for {
		select {
		case <-cancel:
			return nil
		case <-ticker.C:
			topic := topics[rand.Intn(len(topics))]
			slog.Info(
				"Generating questions",
				slog.String("topic", topic),
				slog.Int("amount", amountPerBatch),
			)
			if err := s.generateQuestionSet(ctx, topic, amountPerBatch); err != nil {
				return errors.Wrap(err, "generating question set")
			}
		}
	}
}

//go:embed prompt.txt
var aiSystemPrompt string

// generateQuestionSet generates and stores a set of questions about a given topic.
func (s Service) generateQuestionSet(ctx context.Context, topic string, amount int) error {
	var seed int = time.Now().Nanosecond()
	resp, err := s.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: aiSystemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Generate %d questions about '%s'", amount, topic),
				},
			},
			Seed:      &seed,
			MaxTokens: 2000,
		},
	)
	if err != nil {
		return errors.Wrap(err, "generating question set")
	}

	var questions []aiQuestion
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &questions); err != nil {
		return errors.Wrap(err, "unmarshalling questions")
	}
	for _, q := range questions {
		slog.Info("Inserting question", slog.String("question", q.Question))
		if err := s.repo.Insert(ctx, q.toQuestion()); err != nil {
			return errors.Wrap(err, "creating question")
		}
	}
	return nil
}

// aiChoice represents a choice generated by the AI.
type aiChoice struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

// aiQuestion represents a question generated by the AI.
type aiQuestion struct {
	Question   string     `json:"question"`
	Hint       string     `json:"hint"`
	Choices    []aiChoice `json:"choices"`
	MoreInfo   string     `json:"moreInfo"`
	Difficulty string     `json:"difficulty"`
}

// aiQuestion.toQuestion converts an aiQuestion to a question.Question.
func (a aiQuestion) toQuestion() question.Question {
	q := question.New(a.Question, a.Hint)
	for _, c := range a.Choices {
		q.WithChoice(c.Text, c.IsCorrect)
	}
	return *q
}
