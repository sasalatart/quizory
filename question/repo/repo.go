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
	queryMods := append(
		[]qm.QueryMod{qm.Load(models.QuestionRels.Choices)},
		qms...,
	)
	questions, err := models.Questions(queryMods...).All(ctx, r.db)
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
	dbQuestion, dbChoices, err := r.toDB(q)
	if err != nil {
		return errors.Wrap(err, "mapping question to DB")
	}

	txn, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "beginning transaction")
	}
	if err := dbQuestion.Insert(ctx, txn, boil.Infer()); err != nil {
		txn.Rollback()
		return errors.Wrapf(err, "inserting question %v", q)
	}
	if err := dbQuestion.AddChoices(ctx, txn, true, dbChoices...); err != nil {
		txn.Rollback()
		return errors.Wrapf(err, "adding choices to question %v", q)

	}
	return txn.Commit()
}

func (r *QuestionRepo) fromDB(q *models.Question) (*question.Question, error) {
	id, err := uuid.Parse(q.ID)
	if err != nil {
		return nil, errors.Wrap(err, "parsing question ID")
	}

	result := question.Question{
		ID:        id,
		Question:  q.Question,
		Hint:      q.Hint,
		CreatedAt: q.CreatedAt,
	}
	for _, c := range q.R.Choices {
		choiceID, err := uuid.Parse(c.ID)
		if err != nil {
			return nil, errors.Wrap(err, "parsing choice ID")
		}
		result.Choices = append(result.Choices, question.Choice{
			ID:        choiceID,
			Choice:    c.Choice,
			IsCorrect: c.IsCorrect,
		})
	}
	return &result, nil
}

func (r *QuestionRepo) toDB(q question.Question) (*models.Question, models.ChoiceSlice, error) {
	if err := q.Validate(); err != nil {
		return nil, nil, errors.Wrap(err, "validating question")
	}

	choices := make(models.ChoiceSlice, len(q.Choices))
	for i, c := range q.Choices {
		choices[i] = &models.Choice{
			ID:         c.ID.String(),
			QuestionID: q.ID.String(),
			Choice:     c.Choice,
			IsCorrect:  c.IsCorrect,
		}
	}

	return &models.Question{
		ID:        q.ID.String(),
		Question:  q.Question,
		Hint:      q.Hint,
		CreatedAt: q.CreatedAt,
	}, choices, nil
}
