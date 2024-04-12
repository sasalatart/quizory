package repo

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	models "github.com/sasalatart.com/quizory/db/model"
	"github.com/sasalatart.com/quizory/question"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type QuestionRepo struct {
	db *sql.DB
}

func New(db *sql.DB) *QuestionRepo {
	return &QuestionRepo{db: db}
}

func (r *QuestionRepo) GetMany(ctx context.Context, qms ...qm.QueryMod) ([]question.Question, error) {
	questions, err := models.Questions(qms...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving questions")
	}

	var result []question.Question
	for _, q := range questions {
		question, err := r.fromDB(q)
		if err != nil {
			return nil, errors.Wrapf(err, "mapping question %s", q.ID)
		}
		result = append(result, *question)
	}
	return result, nil
}

func (r *QuestionRepo) Insert(ctx context.Context, q question.Question) error {
	dbQuestion, err := r.toDB(q)
	if err != nil {
		return errors.Wrap(err, "mapping question to DB")
	}

	if err := dbQuestion.Insert(ctx, r.db, boil.Infer()); err != nil {
		return errors.Wrapf(err, "inserting question %v", q)
	}
	return nil
}

func (r *QuestionRepo) fromDB(q *models.Question) (*question.Question, error) {
	id, err := uuid.Parse(q.ID)
	if err != nil {
		return nil, errors.Wrap(err, "parsing question ID")
	}

	return &question.Question{
		ID:        id,
		Question:  q.Question,
		Hint:      q.Hint,
		CreatedAt: q.CreatedAt,
	}, nil
}

func (r *QuestionRepo) toDB(q question.Question) (*models.Question, error) {
	if err := q.Validate(); err != nil {
		return nil, errors.Wrap(err, "validating question")
	}

	return &models.Question{
		ID:        q.ID.String(),
		Question:  q.Question,
		Hint:      q.Hint,
		CreatedAt: q.CreatedAt,
	}, nil
}
