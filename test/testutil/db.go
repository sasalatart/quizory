package testutil

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
	models "github.com/sasalatart/quizory/db/model"
)

// DeleteData deletes all entities from the test database.
func DeleteData(ctx context.Context, db *sql.DB) error {
	if _, err := models.Questions().DeleteAll(ctx, db); err != nil {
		return errors.Wrap(err, "deleting all questions")
	}
	return nil
}
