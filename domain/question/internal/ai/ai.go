package ai

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/domain/question/enums"
	"github.com/sasalatart/quizory/llm"
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
	recentlyGenerated []string,
	amount int,
	results chan Result,
) {
	send := func(res Result) {
		// Avoid sending results if the context has been cancelled
		if ctx.Err() != nil {
			return
		}
		results <- res
	}

	resp, err := llmService.ChatCompletion(
		ctx,
		prompt,
		newUserContent(topic, recentlyGenerated, amount),
	)
	if err != nil {
		send(Result{Err: errors.Wrap(err, "generating AI questions")})
		return
	}

	var questions []Question
	if err := json.Unmarshal([]byte(resp), &questions); err != nil {
		send(Result{Err: errors.Wrap(err, "unmarshalling AI questions")})
		return
	}
	send(Result{Questions: questions})
}

// newUserContent returns a message to be sent to the LLM model, requesting the generation of new
// questions about a given topic, excluding the ones that have been recently generated.
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

//go:embed prompt.txt
var prompt string
