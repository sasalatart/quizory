package repo

import (
	"github.com/google/uuid"
	models "github.com/sasalatart.com/quizory/db/model"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func OrderByCreatedAtDesc() qm.QueryMod {
	return qm.OrderBy(models.AnswerColumns.CreatedAt + " DESC")
}

func WhereUserID(userID uuid.UUID) qm.QueryMod {
	return models.AnswerWhere.UserID.EQ(userID.String())
}
