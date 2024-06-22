package migrations

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/sasalatart/quizory/config"
)

// Up runs the migrations on the database specified by dbUrl.
func Up(dbCfg config.DBConfig) error {
	m, err := migrate.New(fmt.Sprintf("file://%s", dbCfg.MigrationsDir), dbCfg.URL())
	defer func() {
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
