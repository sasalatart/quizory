package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/config"
	"github.com/sasalatart.com/quizory/db/migrations"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"db",
	fx.Provide(newDB),
)

func newDB(lc fx.Lifecycle, cfg config.Config) *sql.DB {
	rootDB := mustOpen(cfg.Database.URL)
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return rootDB.Close()
		},
	})
	return rootDB
}

func mustOpen(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(errors.Wrapf(err, "opening database with URL %s", dbURL))
	}
	return db
}

var TestModule = fx.Module(
	"db-test",
	fx.Provide(newTempDB),
)

func newTempDB(lc fx.Lifecycle, cfg config.Config) *sql.DB {
	rootDB := mustOpen(cfg.Database.URL)
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

func mustCreateTempDB(rootDB *sql.DB) (*sql.DB, string) {
	dbName := "tmp_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	cfg := config.NewDBConfig(dbName)

	if _, err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		panic(errors.Wrapf(err, "creating db %s", dbName))
	}
	if err := migrations.Up(cfg); err != nil {
		panic(errors.Wrapf(err, "running migrations for db %s", dbName))
	}
	return mustOpen(cfg.URL), dbName
}
