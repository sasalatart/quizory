package testutil

import (
	"context"
	"database/sql"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/db/migrations"
	models "github.com/sasalatart.com/quizory/db/model"
)

const dbPort = 5433
const dbConnString = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

// NewDB creates a new embedded PostgreSQL database for testing purposes and returns a connection to
// it, plus a teardown function that should be called after the test finishes.
func NewDB() (*sql.DB, func(), error) {
	config := embeddedpostgres.
		DefaultConfig().
		Port(dbPort).
		Logger(nil)
	database := embeddedpostgres.NewDatabase(config)
	if err := database.Start(); err != nil {
		return nil, nil, errors.Wrap(err, "starting embedded database")
	}

	teardown := func() {
		database.Stop()
	}

	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		teardown()
		return nil, nil, errors.Wrap(err, "opening connection to embedded database")
	}

	if err := migrations.Up(dbConnString); err != nil {
		teardown()
		return nil, nil, errors.Wrap(err, "running migrations")
	}

	return db, teardown, nil
}

func WipeDB(ctx context.Context, db *sql.DB) error {
	if _, err := models.Questions().DeleteAll(ctx, db); err != nil {
		return errors.Wrap(err, "deleting all questions")
	}
	return nil
}
