package config

import (
	"time"
)

// LLMConfig represents the configuration of the Language Model.
type LLMConfig struct {
	// OpenAIKey is the API key used to interact with OpenAI's API.
	OpenAIKey string

	// Frequency is the time interval between each generation of questions.
	Frequency time.Duration

	// BatchSize is the number of questions to generate in each batch.
	BatchSize int
}

// NewLLMConfig returns a new LLMConfig instance with values loaded from environment variables.
func NewLLMConfig(openAIKey string) LLMConfig {
	return LLMConfig{
		OpenAIKey: openAIKey,
		Frequency: 12 * time.Hour,
		BatchSize: 8,
	}
}
