package config

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"go.uber.org/fx"
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
		Server: NewServerConfig("localhost", 8080),
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
