package testutil

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/config"
	"github.com/sasalatart/quizory/db/migrations"
	models "github.com/sasalatart/quizory/db/model"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
)

var dbModule = fx.Module(
	"db-test",
	fx.Provide(newPostgresContainer),
	fx.Provide(newDB),
)

// newPostgresContainer creates a new testcontainer running a postgres instance, waits for it to be
// ready, runs migrations on it, and makes sure it is terminated when the fx app is stopped.
func newPostgresContainer(lc fx.Lifecycle, dbCfg config.DBConfig) *postgres.PostgresContainer {
	ctx := context.Background()
	container, err := postgres.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		log.Fatalf("failed to start postgres testcontainer: %s", err)
	}

	connString := container.MustConnectionString(ctx, "sslmode=disable")
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return migrations.Up(connString, dbCfg.MigrationsDir)
		},
		OnStop: func(context.Context) error {
			return container.Terminate(ctx)
		},
	})

	return container
}

// newDB creates a new sql.DB instance connected to the test postgres container, and makes sure it
// is closed when the fx app is stopped.
func newDB(lc fx.Lifecycle, container *postgres.PostgresContainer) *sql.DB {
	ctx := context.Background()
	connString := container.MustConnectionString(ctx, "sslmode=disable")
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatalf("failed to open connection to postgres testcontainer: %s", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return db.Close()
		},
	})

	return db
}

// DeleteData deletes all entities from the test database.
func DeleteData(ctx context.Context, db *sql.DB) error {
	if _, err := models.Questions().DeleteAll(ctx, db); err != nil {
		return errors.Wrap(err, "deleting all questions")
	}
	return nil
}
