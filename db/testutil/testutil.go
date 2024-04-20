package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sasalatart.com/quizory/db/migrations"
	models "github.com/sasalatart.com/quizory/db/model"
)

const dbPort = 5433

// TestDB wraps a connection to a test database and provides utility methods to interact with it.
type TestDB struct {
	db   *sql.DB
	name string
}

// NewTestDB creates a new test database and runs migrations on it.
func NewTestDB(ctx context.Context) (*TestDB, error) {
	rootDB := mustOpenDB("postgres")
	defer rootDB.Close()

	dbName := "test_" + strings.ReplaceAll(uuid.New().String(), "-", "")
	testDB := &TestDB{
		name: dbName,
		db:   mustOpenDB(dbName),
	}
	if _, err := rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		return testDB, errors.Wrapf(err, "creating db %s", dbName)
	}
	if err := migrations.Up(dbConnString(dbName)); err != nil {
		return testDB, errors.Wrapf(err, "running migrations for db %s", dbName)
	}
	return testDB, nil
}

// DeleteData deletes all entities from the test database.
func (d *TestDB) DeleteData(ctx context.Context) error {
	if _, err := models.Questions().DeleteAll(ctx, d.db); err != nil {
		return errors.Wrapf(err, "deleting all questions from DB %s", d.name)
	}
	return nil
}

// Teardown closes the connection to the test database and drops it.
func (d *TestDB) Teardown() error {
	rootDB := mustOpenDB("postgres")
	defer rootDB.Close()

	if err := d.db.Close(); err != nil {
		return errors.Wrapf(err, "closing connection to DB %s", d.name)
	}
	if _, err := rootDB.Exec(fmt.Sprintf("DROP DATABASE %s", d.name)); err != nil {
		return errors.Wrapf(err, "dropping DB %s", d.name)
	}
	return nil
}

// DB returns the connection to the test database.
func (d *TestDB) DB() *sql.DB {
	return d.db
}

// mustOpenDB opens a connection to a database. It panics if it fails to do so.
func mustOpenDB(dbName string) *sql.DB {
	db, err := sql.Open("postgres", dbConnString(dbName))
	if err != nil {
		panic(errors.Wrapf(err, "opening connection to db %s", dbName))
	}
	return db
}

// dbConnString returns a connection string for a database with the given name.
func dbConnString(dbName string) string {
	return fmt.Sprintf(
		"postgres://postgres:postgres@localhost:%d/%s?sslmode=disable",
		dbPort, dbName,
	)
}
