package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver

	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db/migrations"
	"github.com/sasalatart/quizory/infra"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"db",
	fx.Provide(newDB),
	fx.Invoke(migrationsLC),
)

// newDB creates a database connection, and adds an fx hook to close it when the application stops.
func newDB(lc fx.Lifecycle, cfg config.DBConfig) *sql.DB {
	rootDB := mustOpen(cfg.URL)
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

func migrationsLC(lc fx.Lifecycle, dbCfg config.DBConfig) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return migrations.Up(dbCfg.URL, dbCfg.MigrationsDir())
		},
	})
}
