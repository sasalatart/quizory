package migrations

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
)

// Up runs the migrations on the database specified by dbUrl.
func Up(url, migrationsDir string) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", migrationsDir), url)
	defer func() {
		if m == nil {
			return
		}
		if _, err := m.Close(); err != nil {
			slog.Error("Failed to close migrations instance", slog.Any("error", err))
		}
	}()
	if err != nil {
		return errors.Wrap(err, "creating migrations instance")
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		slog.Info("No migrations to apply")
		return nil
	}
	return err
}
