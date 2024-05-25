package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

func init() {
	mustLoadEnvVars()
}

// mustLoadEnvVars loads environment variables from the root .env file if it exists.
func mustLoadEnvVars() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic(errors.Wrap(err, "loading .env file"))
	}
}

// Config represents the configuration of the application.
type Config struct {
	fx.Out

	DB     DBConfig
	LLM    LLMConfig
	Server ServerConfig
}

// NewConfig returns a new Config instance with values loaded from environment variables.
func NewConfig() Config {
	openAIKey := os.Getenv("OPENAI_API_KEY")

	return Config{
		DB:     NewDBConfig("postgres"),
		LLM:    NewLLMConfig(openAIKey),
		Server: NewServerConfig("0.0.0.0", 8080),
	}
}

// NewTestConfig returns a Config instance intended for testing.
func NewTestConfig() Config {
	return Config{
		DB:     NewDBConfig("postgres"),
		LLM:    NewLLMConfig("test"),
		Server: NewServerConfig("localhost", 8081),
	}
}
