package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

const envFileName = ".env.quizory"

func init() {
	mustLoadEnvVars()
}

// mustLoadEnvVars loads environment variables from the root .env file if it exists.
func mustLoadEnvVars() {
	envFileDir, err := findFilePath(envFileName)
	if err != nil {
		panic(errors.Wrapf(err, "finding %s file", envFileName))
	}
	if err := godotenv.Load(envFileDir); err != nil && !os.IsNotExist(err) {
		panic(errors.Wrapf(err, "loading %s file", envFileName))
	}
}

// findFilePath searches for a file with the given name starting from the current working directory
// and going up to the root directory. It returns the path of the file if found, or an error if it
// does not exist or if there was an error while searching.
func findFilePath(fileName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "getting current working directory")
	}

	for dir := cwd; dir != "/"; dir = filepath.Dir(dir) {
		possibleFilePath := filepath.Join(dir, fileName)
		if _, err := os.Stat(possibleFilePath); err == nil {
			return possibleFilePath, nil
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}

	return "", fmt.Errorf("no file found called %s", fileName)
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
