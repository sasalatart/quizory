package repo

import (
	models "github.com/sasalatart.com/quizory/db/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func OrderByCreatedAtDesc() qm.QueryMod {
	return qm.OrderBy(models.QuestionColumns.CreatedAt + " DESC")
}
