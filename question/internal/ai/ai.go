package ai

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/llm"
	"github.com/sasalatart.com/quizory/question/enums"
)

// Result represents the result of generating questions, intended to be used for communication
// between goroutines.
type Result struct {
	Questions []Question
	Err       error
}

// Generate generates a set of questions about a given topic using an LLM model, and sends the
// result to the input channel.
func Generate(
	ctx context.Context,
	llmService llm.ChatCompletioner,
	topic enums.Topic,
	amount int,
	results chan Result,
) {
	resp, err := llmService.ChatCompletion(
		ctx,
		prompt,
		fmt.Sprintf("Generate %d questions about '%s'", amount, topic),
	)
	if err != nil {
		results <- Result{Err: errors.Wrap(err, "generating AI questions")}
		return
	}

	var questions []Question
	if err := json.Unmarshal([]byte(resp), &questions); err != nil {
		results <- Result{Err: errors.Wrap(err, "unmarshalling AI questions")}
		return
	}
	results <- Result{Questions: questions}
}

//go:embed prompt.txt
var prompt string
