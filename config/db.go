package config

import (
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/pkg/errors"
)

// DBCOnfig represents the configuration of the database.
type DBConfig struct {
	URL           string
	MigrationsDir string
	User          string
	Password      string
	Host          string
	Port          string
	Name          string
}

// NewDBConfig returns a new DBConfig instance with values loaded from environment variables.
// The dbName argument is used to specify the name of the database to connect to.
func NewDBConfig(dbName string) DBConfig {
	dbURL := os.Getenv("DB_URL")

	u, err := url.Parse(dbURL)
	if err != nil {
		panic(errors.Wrap(err, "parsing DB_URL"))
	}

	psqlUser := u.User.Username()
	psqlPassword, _ := u.User.Password()
	psqlHost, psqlPort, err := net.SplitHostPort(u.Host)
	if err != nil {
		panic(errors.Wrap(err, "splitting host and port"))
	}

	return DBConfig{
		URL: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			psqlUser, psqlPassword, psqlHost, psqlPort, dbName,
		),
		MigrationsDir: mustGetModuleRoot() + "/db/migrations",
		User:          psqlUser,
		Password:      psqlPassword,
		Host:          psqlHost,
		Port:          psqlPort,
		Name:          dbName,
	}
}
