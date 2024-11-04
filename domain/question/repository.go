package question

import (
	"context"
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	models "github.com/sasalatart/quizory/db/model"
	"github.com/sasalatart/quizory/domain/question/enums"
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

func (r *Repository) GetOne(ctx context.Context, qms ...qm.QueryMod) (*Question, error) {
	q, err := models.Questions(r.withChoices(qms)...).One(ctx, r.db)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving question")
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

// GetRemainingTopics returns a map such that each key is a topic for which the user still has
// unanswered questions, and each value is the amount of remaining questions for that topic.
// This might look like it does not belong in the repository, but if we do this at the service level
// by loading all unanswered questions for a user and grouping via code, we risk loading ALL
// questions in the database in the worst case (e.g. if the user has not answered any question).
func (r *Repository) GetRemainingTopics(
	ctx context.Context,
	userID uuid.UUID,
) (map[enums.Topic]uint, error) {
	query := `
		SELECT topic, COUNT(*) AS count
		FROM questions
		WHERE NOT EXISTS (
			SELECT 1
			FROM answers a
			JOIN choices c ON c.id = a.choice_id
			WHERE c.question_id = questions.id AND a.user_id = $1
		)
		GROUP BY topic
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topicCounts := make(map[enums.Topic]uint)
	for rows.Next() {
		var tc struct {
			Topic string `boil:"topic"`
			Count uint   `boil:"count"`
		}
		if err := rows.Scan(&tc.Topic, &tc.Count); err != nil {
			return nil, err
		}
		topic, err := enums.TopicString(tc.Topic)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing topic %s", tc.Topic)
		}
		topicCounts[topic] = tc.Count
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return topicCounts, nil
}

// withChoices eager loads the choices of the questions.
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

func WhereIDIn(ids ...uuid.UUID) qm.QueryMod {
	var strIDs []string
	for _, id := range ids {
		strIDs = append(strIDs, id.String())
	}
	return models.QuestionWhere.ID.IN(strIDs)
}

func WhereChoiceIDIn(ids ...uuid.UUID) qm.QueryMod {
	strIDs := make([]interface{}, len(ids))
	placeholders := make([]string, len(ids))

	for i, id := range ids {
		strIDs[i] = id.String()
		placeholders[i] = "?"
	}

	return qm.Where(`
		questions.id IN (
			SELECT DISTINCT(c.question_id)
			FROM choices c
			WHERE c.id IN (`+strings.Join(placeholders, ", ")+`)
		)
	`, strIDs...)
}

func WhereTopicEq(topic enums.Topic) qm.QueryMod {
	return models.QuestionWhere.Topic.EQ(topic.String())
}

func WhereNotAnsweredBy(userID uuid.UUID) qm.QueryMod {
	return qm.Where(`
		NOT EXISTS (
			SELECT 1
			FROM answers a
			JOIN choices c ON a.choice_id = c.id
			WHERE c.question_id = questions.id AND a.user_id = ?
		)
		`, userID.String(),
	)
}

func OrderByCreatedAtAsc() qm.QueryMod {
	return qm.OrderBy(models.QuestionColumns.CreatedAt + " ASC")
}

func OrderByCreatedAtDesc() qm.QueryMod {
	return qm.OrderBy(models.QuestionColumns.CreatedAt + " DESC")
}

func Limit(n int) qm.QueryMod {
	return qm.Limit(n)
}
