package question

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	models "github.com/sasalatart.com/quizory/db/model"
	"github.com/sasalatart.com/quizory/question/enums"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetMany(ctx context.Context, qms ...qm.QueryMod) ([]Question, error) {
	questions, err := models.Questions(r.withChoices(qms)...).All(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving questions")
	}

	var result []Question
	for _, q := range questions {
		question, err := r.fromDB(q)
		if err != nil {
			return nil, errors.Wrapf(err, "mapping question %s", q.ID)
		}
		result = append(result, *question)
	}
	return result, nil
}

func (r *Repository) GetByChoiceID(ctx context.Context, choiceID uuid.UUID) (*Question, error) {
	c, err := models.FindChoice(ctx, r.db, choiceID.String())
	if err != nil {
		return nil, errors.Wrap(err, "retrieving choice")
	}
	q, err := c.Question(r.withChoices(nil)...).One(ctx, r.db)
	if err != nil {
		return nil, errors.Wrapf(err, "retrieving question for choice %s", choiceID)
	}
	return r.fromDB(q)
}

func (r *Repository) Insert(ctx context.Context, q Question) error {
	dbQuestion, dbChoices, err := r.toDB(q)
	if err != nil {
		return errors.Wrap(err, "mapping question to DB")
	}

	txn, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "beginning transaction")
	}
	if err := dbQuestion.Insert(ctx, txn, boil.Infer()); err != nil {
		return multierror.Append(
			errors.Wrapf(err, "inserting question %v", q),
			txn.Rollback(),
		)
	}
	if err := dbQuestion.AddChoices(ctx, txn, true, dbChoices...); err != nil {
		return multierror.Append(
			errors.Wrapf(err, "adding choices to question %v", q),
			txn.Rollback(),
		)
	}
	return txn.Commit()
}

func (r *Repository) withChoices(qms []qm.QueryMod) []qm.QueryMod {
	return append(qms, qm.Load(models.QuestionRels.Choices))
}

func (r *Repository) fromDB(q *models.Question) (*Question, error) {
	id, err := uuid.Parse(q.ID)
	if err != nil {
		return nil, errors.Wrap(err, "parsing question ID")
	}
	topic, err := enums.TopicString(q.Topic)
	if err != nil {
		return nil, errors.Wrap(err, "parsing question topic")
	}
	difficulty, err := enums.DifficultyString(q.Difficulty)
	if err != nil {
		return nil, errors.Wrap(err, "parsing question difficulty")
	}

	result := Question{
		ID:         id,
		Topic:      topic,
		Question:   q.Question,
		Hint:       q.Hint,
		MoreInfo:   q.MoreInfo,
		Difficulty: difficulty,
		CreatedAt:  q.CreatedAt,
	}
	for _, c := range q.R.Choices {
		choiceID, err := uuid.Parse(c.ID)
		if err != nil {
			return nil, errors.Wrap(err, "parsing choice ID")
		}
		result.Choices = append(result.Choices, Choice{
			ID:        choiceID,
			Choice:    c.Choice,
			IsCorrect: c.IsCorrect,
		})
	}
	return &result, nil
}

func (r *Repository) toDB(q Question) (*models.Question, models.ChoiceSlice, error) {
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
		ID:         q.ID.String(),
		Topic:      q.Topic.String(),
		Question:   q.Question,
		Hint:       q.Hint,
		MoreInfo:   q.MoreInfo,
		Difficulty: q.Difficulty.String(),
		CreatedAt:  q.CreatedAt,
	}, choices, nil
}

func WhereTopicIs(topic enums.Topic) qm.QueryMod {
	return models.QuestionWhere.Topic.EQ(topic.String())
}

func OrderByCreatedAtDesc() qm.QueryMod {
	return qm.OrderBy(models.QuestionColumns.CreatedAt + " DESC")
}

func Limit(n int) qm.QueryMod {
	return qm.Limit(n)
}
