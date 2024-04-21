package config

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func init() {
	mustLoadEnvVars()
}

// mustLoadEnvVars loads environment variables from the root .env file if it exists.
func mustLoadEnvVars() {
	dotEnvDir := mustGetModuleRoot() + "/.env"
	_, err := os.Stat(dotEnvDir)

	if os.IsNotExist(err) {
		return
	}
	if err := godotenv.Load(dotEnvDir); err != nil {
		panic(errors.Wrap(err, "loading .env file"))
	}
}

// Config represents the configuration of the application.
type Config struct {
	Database           DBConfig
	OpenAIKey          string
	QuestionGeneration questionGenerationConfig
}

type questionGenerationConfig struct {
	Frequency time.Duration
	BatchSize int
}

// NewConfig returns a new Config instance with values loaded from environment variables.
func NewConfig() Config {
	return Config{
		Database:  NewDBConfig("postgres"),
		OpenAIKey: os.Getenv("OPENAI_API_KEY"),
		QuestionGeneration: questionGenerationConfig{
			Frequency: 5 * time.Second,
			BatchSize: 5,
		},
	}
}

// mustGetModuleRoot returns the root directory of the module.
func mustGetModuleRoot() string {
	cmd := exec.Command("go", "list", "-m", "-f", "{{.Dir}}")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		panic(errors.Wrap(err, "getting module root"))
	}
	return strings.TrimSpace(stdout.String())
}
