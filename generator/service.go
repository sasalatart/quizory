package generator

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/http/grpc/proto"
	"github.com/sasalatart/quizory/infra/otel"
	"github.com/sasalatart/quizory/llm"
	"go.opentelemetry.io/otel/metric"
)

//go:embed prompt.txt
var generatorPrompt string

// recentlyGeneratedLimit is the amount of questions that are considered recent, and are used as
// part of the LLM's context when generating new questions.
const recentlyGeneratedLimit = 100

// Service represents the service that generates questions via LLMs.
type Service struct {
	durationHistogram metric.Int64Histogram
	llm               llm.ChatCompletioner
	quizoryClient     proto.QuizoryServiceClient
}

// NewService creates a new instance of question.Service.
func NewService(
	quizoryClient proto.QuizoryServiceClient,
	llm llm.ChatCompletioner,
	meter otel.Meter,
) *Service {
	durationHistogram, err := meter.Int64Histogram(
		"questions_generation_duration",
		metric.WithUnit("ms"),
	)
	if err != nil {
		log.Fatal("unable to create questions_generation_duration histogram")
	}

	return &Service{
		durationHistogram: durationHistogram,
		llm:               llm,
		quizoryClient:     quizoryClient,
	}
}

// GenerateBatch generates and persists a batch of questions about the given topic.
func (s Service) GenerateBatch(ctx context.Context, batchSize int, topic enums.Topic) (err error) {
	slog.Info(
		"Generating questions",
		slog.String("topic", topic.String()),
		slog.Int("batchSize", batchSize),
	)

	startTime := time.Now()
	defer func() {
		s.durationHistogram.Record(ctx, time.Since(startTime).Milliseconds())
	}()

	questionsGenResult := make(chan error)
	go func() {
		var questions []question
		questions, err = s.newBatchFromLLM(ctx, topic, batchSize)
		if err != nil {
			questionsGenResult <- errors.Wrap(err, "generating questions")
			return
		}
		if err := s.persistQuestions(ctx, topic, questions); err != nil {
			questionsGenResult <- errors.Wrap(err, "persisting questions")
			return
		}
		questionsGenResult <- nil
	}()

	select {
	case err := <-questionsGenResult:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// newBatchFromLLM generates a set of unpersisted questions about a given topic using an LLM model.
func (s Service) newBatchFromLLM(
	ctx context.Context,
	topic enums.Topic,
	batchSize int,
) ([]question, error) {
	recentlyGenerated, err := s.quizoryClient.GetLatestQuestions(
		ctx,
		&proto.GetLatestQuestionsRequest{
			Topic:  topic.String(),
			Amount: recentlyGeneratedLimit,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "getting recently generated questions about %s", topic)
	}

	llmResp, err := s.llm.ChatCompletion(
		ctx,
		generatorPrompt,
		newUserContent(topic, recentlyGenerated.GetQuestions(), batchSize),
	)
	if err != nil {
		return nil, errors.Wrap(err, "calling LLM")
	}

	var questions []question
	if err := json.Unmarshal([]byte(llmResp), &questions); err != nil {
		return nil, errors.Wrap(err, "unmarshalling LLM questions")
	}
	return questions, nil
}

// newUserContent returns the USER content to be used as context for the LLM when generating new
// questions.
func newUserContent(topic enums.Topic, recentlyGenerated []string, amount int) string {
	baseMsg := fmt.Sprintf("Generate %d new questions about '%s'.", amount, topic)
	if len(recentlyGenerated) == 0 {
		return baseMsg
	}

	return fmt.Sprintf(`
		%s
		You have already generated the following questions, therefore provide new ones:
		- %s
		`, baseMsg, strings.Join(recentlyGenerated, "\n -"),
	)
}

// persistQuestions makes the necessary RPCs to persist the given generated questions.
func (s Service) persistQuestions(
	ctx context.Context,
	topic enums.Topic,
	questions []question,
) error {
	for _, q := range questions {
		var choices []*proto.Choice
		for _, c := range q.Choices {
			choices = append(choices, &proto.Choice{
				Choice:    c.Text,
				IsCorrect: c.IsCorrect,
			})
		}
		resp, err := s.quizoryClient.CreateQuestion(ctx, &proto.CreateQuestionRequest{
			Question:   q.Question,
			Hint:       q.Hint,
			Topic:      topic.String(),
			Difficulty: q.parseDifficulty(),
			MoreInfo:   strings.Join(q.MoreInfo, "\n"),
			Choices:    choices,
		})
		if err != nil {
			return errors.Wrap(err, "creating question")
		}
		slog.Info(
			"Persisted question",
			slog.String("q", q.Question),
			slog.String("id", resp.GetId()),
		)
	}
	return nil
}

type question struct {
	Question   string   `json:"question"`
	Hint       string   `json:"hint"`
	Choices    []choice `json:"choices"`
	MoreInfo   []string `json:"moreInfo"`
	Difficulty string   `json:"difficulty"`
}

func (q question) parseDifficulty() proto.Difficulty {
	switch q.Difficulty {
	case enums.DifficultyNoviceHistorian.String():
		return proto.Difficulty_DIFFICULTY_NOVICE_HISTORIAN
	case enums.DifficultyAvidHistorian.String():
		return proto.Difficulty_DIFFICULTY_AVID_HISTORIAN
	case enums.DifficultyHistoryScholar.String():
		return proto.Difficulty_DIFFICULTY_HISTORY_SCHOLAR
	default:
		return proto.Difficulty_DIFFICULTY_UNSPECIFIED
	}
}

type choice struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}
