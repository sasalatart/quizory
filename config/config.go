package config

import (
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Config struct {
	fx.Out

	DB     DBConfig     `mapstructure:"db"`
	LLM    LLMConfig    `mapstructure:"llm"`
	Server ServerConfig `mapstructure:"server"`
}

type DBConfig struct {
	URL string `mapstructure:"url"`
}

func (c DBConfig) MigrationsDir() string {
	return mustFindAbsoluteFilePath("db/migrations")
}

type LLMConfig struct {
	OpenAIKey string `mapstructure:"openai_key"`

	Questions struct {
		BatchSize int           `mapstructure:"batch_size"`
		Frequency time.Duration `mapstructure:"frequency"`
	} `mapstructure:"questions"`
}

type ServerConfig struct {
	Host         string        `mapstructure:"host"`
	RESTPort     int           `mapstructure:"rest_port"`
	GRPCPort     int           `mapstructure:"grpc_port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	JWTSecret    string        `mapstructure:"jwt_secret"`
}

func (c ServerConfig) RESTAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.RESTPort)
}

func (c ServerConfig) OAPISchemaDir() string {
	return mustFindAbsoluteFilePath("http/rest/oapi/schema.yaml")
}

func (c ServerConfig) GRPCAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.GRPCPort)
}

// NewConfig returns a new Config instance with values loaded from environment variables.
func NewConfig() Config {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	configPath := mustFindAbsoluteFilePath(path.Join("config", "config.yaml"))
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		log.Fatal("error reading config file: ", err)
	}

	var config Config
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     &config,
		TagName:    "mapstructure",
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		log.Fatal("error creating mapstructure decoder: ", err)
	}
	if err := decoder.Decode(v.AllSettings()); err != nil {
		log.Fatal("error decoding config: ", err)
	}
	return config
}

// NewTestConfig returns a Config instance intended for testing.
func NewTestConfig() Config {
	cfg := NewConfig()
	cfg.Server.Host = "localhost" // Test can run outside Docker
	cfg.Server.RESTPort = 8081    // Avoid conflicts with the main server in case it's running in dev mode
	cfg.Server.GRPCPort = 8086    // Avoid conflicts with the main server in case it's running in dev mode
	cfg.LLM.OpenAIKey = "test"
	return cfg
}
