package question

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sasalatart.com/quizory/question/enums"
)

// Mock creates a new random Question for testing purposes with the specified overrides.
func Mock(applyOverrides func(*Question)) Question {
	q := New(
		fmt.Sprintf("%s-%s", "Test Question", uuid.New()),
		"Test Hint",
		"Test More Info",
	).
		WithTopic(enums.TopicAncientRome).
		WithDifficulty(enums.DifficultyAvidHistorian).
		WithChoice("Choice 1", false).
		WithChoice("Choice 2", true).
		WithChoice("Choice 3", false).
		WithChoice("Choice 4", false)

	if applyOverrides != nil {
		applyOverrides(q)
	}

	return *q
}
