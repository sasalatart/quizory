package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db/migrations"
	"github.com/sasalatart/quizory/infra"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"db",
	fx.Provide(newDB),
)

// newDB creates a database connection, and adds an fx hook to close it when the application stops.
func newDB(lc fx.Lifecycle, cfg config.DBConfig) *sql.DB {
	rootDB := mustOpen(cfg.URL())
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return rootDB.Close()
		},
	})
	return rootDB
}

// mustOpen opens a database connection with the given URL, and waits for it to be ready.
// It retries up to 5 times with an exponential backoff, starting at 1 second.
// It panics if the connection cannot be established after the retries.
func mustOpen(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(errors.Wrapf(err, "opening database with URL %s", dbURL))
	}

	check := func() bool {
		err := db.Ping()
		return err == nil
	}
	if err := infra.WaitFor(check, 5, 1*time.Second); err != nil {
		panic(errors.Wrapf(err, "timeout pinging database with URL %s", dbURL))
	}

	return db
}

var TestModule = fx.Module(
	"db-test",
	fx.Provide(newTempDB),
)

// newTempDB creates a temporary database for testing purposes, and adds fx hooks to clean it up.
func newTempDB(lc fx.Lifecycle, cfg config.DBConfig) *sql.DB {
	rootDB := mustOpen(cfg.URL())
	db, dbName := mustCreateTempDB(rootDB)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			db.Close()
			if _, err := rootDB.Exec(fmt.Sprintf("DROP DATABASE %s", dbName)); err != nil {
				return errors.Wrapf(err, "dropping DB %s", dbName)
			}
			return rootDB.Close()
		},
	})
	return db
}

// mustCreateTempDB creates a temporary database with a random name.
func mustCreateTempDB(rootDB *sql.DB) (*sql.DB, string) {
	dbName := "tmp_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	cfg := config.NewDBConfig(dbName)

	if _, err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		panic(errors.Wrapf(err, "creating db %s", dbName))
	}
	if err := migrations.Up(cfg); err != nil {
		panic(errors.Wrapf(err, "running migrations for db %s", dbName))
	}
	return mustOpen(cfg.URL()), dbName
}
