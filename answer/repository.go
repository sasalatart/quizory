package answer

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	models "github.com/sasalatart.com/quizory/db/model"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMany(ctx context.Context, qms ...qm.QueryMod) ([]Answer, error) {
	answers, err := models.Answers(qms...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving answers")
	}

	var result []Answer
	for _, a := range answers {
		answer, err := r.fromDB(a)
		if err != nil {
			return nil, errors.Wrapf(err, "mapping answer %s", a.ID)
		}
		result = append(result, *answer)
	}
	return result, nil
}

func (r *Repository) Insert(ctx context.Context, a Answer) error {
	dbAnswer, err := r.toDB(a)
	if err != nil {
		return errors.Wrap(err, "mapping answer to DB")
	}
	if err := dbAnswer.Insert(ctx, r.db, boil.Infer()); err != nil {
		return errors.Wrapf(err, "inserting answer %v", a)
	}
	return nil
}

func (r *Repository) fromDB(a *models.Answer) (*Answer, error) {
	id, err := uuid.Parse(a.ID)
	if err != nil {
		return nil, errors.Wrap(err, "parsing answer ID")
	}
	userID, err := uuid.Parse(a.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "parsing user ID")
	}
	choiceID, err := uuid.Parse(a.ChoiceID)
	if err != nil {
		return nil, errors.Wrap(err, "parsing choice ID")
	}

	return &Answer{
		ID:        id,
		UserID:    userID,
		ChoiceID:  choiceID,
		CreatedAt: a.CreatedAt,
	}, nil
}

func (r *Repository) toDB(a Answer) (*models.Answer, error) {
	if err := a.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating answer")
	}

	return &models.Answer{
		ID:        a.ID.String(),
		UserID:    a.UserID.String(),
		ChoiceID:  a.ChoiceID.String(),
		CreatedAt: a.CreatedAt,
	}, nil
}

func OrderByCreatedAtDesc() qm.QueryMod {
	return qm.OrderBy(models.AnswerColumns.CreatedAt + " DESC")
}

func WhereUserID(userID uuid.UUID) qm.QueryMod {
	return models.AnswerWhere.UserID.EQ(userID.String())
}
