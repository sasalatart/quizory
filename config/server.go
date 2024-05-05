package config

import (
	"fmt"
	"os"
)

// ServerConfig represents the configuration of the server.
type ServerConfig struct {
	Host      string
	Port      int
	JWTSecret string
}

// NewServerConfig returns a new ServerConfig instance with values loaded from environment variables.
func NewServerConfig(host string, port int) ServerConfig {
	jwtSecret := os.Getenv("JWT_SECRET")

	return ServerConfig{
		Host:      host,
		Port:      port,
		JWTSecret: jwtSecret,
	}
}

func (c ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
