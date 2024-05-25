package config

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
)

// ServerConfig represents the configuration of the server.
type ServerConfig struct {
	Host      string
	Port      int
	JWTSecret string
	SchemaDir string
}

// NewServerConfig returns a new ServerConfig instance with values loaded from environment variables.
func NewServerConfig(host string, port int) ServerConfig {
	jwtSecret := os.Getenv("JWT_SECRET")

	schemaDir, err := findFilePath(path.Join("http", "oapi", "schema.yml"))
	if err != nil {
		panic(errors.Wrap(err, "finding OAPI schema file"))
	}

	return ServerConfig{
		Host:      host,
		Port:      port,
		JWTSecret: jwtSecret,
		SchemaDir: schemaDir,
	}
}

func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
