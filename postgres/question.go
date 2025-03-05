package postgres

import (
	"context"

	"github.com/akionka/akionkabot/data"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRepository struct {
	db *pgxpool.Pool
}

func NewQuestionRepository(db *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{
		db: db,
	}
}

func (r *QuestionRepository) GetQuestion(ctx context.Context, id uuid.UUID) (*data.Question, error) {
	var question *data.Question

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return question, err
	}

	if err = transaction(ctx, tx, "GetQuestion", func() error {
		question, err = r.getQuestionTx(ctx, tx, id)
		if err != nil {
			return err
		}
		question, err = r.enrichQuestionTx(ctx, tx, question)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return question, err
	}

	return question, nil
}

func (r *QuestionRepository) getQuestionTx(ctx context.Context, tx pgx.Tx, id uuid.UUID) (*data.Question, error) {
	const sql = `
	SELECT q.question_id, q.match_id, q.player_id, q.player_name, q.player_is_pro, q.player_pos, q.player_mmr, q.is_won, q.created_at, q.telegram_file_id
	FROM questions q
	WHERE q.question_id = $1
	LIMIT 1`

	rows, err := tx.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.Question])
}

func (r *QuestionRepository) GetQuestionAvailableForUser(ctx context.Context, id uuid.UUID) (*data.Question, error) {
	var question *data.Question

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return question, err
	}

	if err = transaction(ctx, tx, "GetQuestionAvailableForUser", func() error {
		question, err = r.getQuestionAvailableForUserTx(ctx, tx, id)
		if err != nil {
			return err
		}
		question, err = r.enrichQuestionTx(ctx, tx, question)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return question, err
	}

	return question, nil
}

func (r *QuestionRepository) getQuestionAvailableForUserTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID) (*data.Question, error) {
	const sql = `
	SELECT q.question_id, q.match_id, q.player_id, q.player_name, q.player_is_pro, q.player_pos, q.player_mmr, q.is_won, q.created_at, q.telegram_file_id
	FROM questions q
	LEFT JOIN user_questions uq ON q.question_id = uq.question_id AND uq.user_id = $1
	WHERE uq.question_id IS NULL
	ORDER BY q.created_at DESC
	LIMIT 1;`

	rows, err := tx.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.Question])
}

func (r *QuestionRepository) enrichQuestionTx(ctx context.Context, tx pgx.Tx, q *data.Question) (*data.Question, error) {
	const heroSQL = `
	SELECT qo.is_correct, h.hero_id, h.display_name, h.short_name, qo.telegram_file_id
	FROM question_options qo
	INNER JOIN heroes h ON qo.hero_id = h.hero_id
	WHERE qo.question_id = $1
	ORDER BY qo.hero_id`

	const itemSQL = `
	SELECT i.item_id, i.display_name, i.short_name
	FROM question_items qi
	INNER JOIN items i ON qi.item_id = i.item_id
	WHERE qi.question_id = $1`

	heroRows, err := tx.Query(ctx, heroSQL, q.ID)
	if err != nil {
		return q, err
	}

	q.Options, err = pgx.CollectRows(heroRows, pgx.RowToStructByName[data.Option])
	if err != nil {
		return q, err
	}

	itemRows, err := tx.Query(ctx, itemSQL, q.ID)
	if err != nil {
		return q, err
	}

	q.Items, err = pgx.CollectRows(itemRows, pgx.RowToStructByName[data.Item])
	return q, err
}

func (r *QuestionRepository) SaveQuestion(ctx context.Context, q *data.Question) (*data.Question, error) {
	var question *data.Question

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return question, err
	}

	if err = transaction(ctx, tx, "SaveQuestion", func() error {
		id, err := r.saveQuestionTx(ctx, tx, q)
		if err != nil {
			return err
		}

		question, err = r.getQuestionTx(ctx, tx, id)
		if err != nil {
			return err
		}
		question, err = r.enrichQuestionTx(ctx, tx, question)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return question, err
	}

	return question, nil
}

func (r *QuestionRepository) saveQuestionTx(ctx context.Context, tx pgx.Tx, q *data.Question) (uuid.UUID, error) {
	const sql = `
	INSERT INTO questions (question_id, match_id, player_id, player_name, player_is_pro, player_pos, player_mmr, is_won, created_at, telegram_file_id) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING question_id`

	qID := uuid.Nil

	rows, err := tx.Query(ctx, sql,
		q.ID, q.MatchID, q.PlayerID, q.PlayerName, q.PlayerIsPro, q.PlayerPos, q.PlayerMMR, q.IsWon, q.CreatedAt, q.TelegramFileID,
	)
	if err != nil {
		return qID, err
	}
	qID, err = pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return qID, err
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"question_items"}, []string{"question_id", "item_id"}, pgx.CopyFromSlice(len(q.Items), func(i int) ([]any, error) {
		return []any{qID, q.Items[i].ID}, nil
	}))
	if err != nil {
		return qID, err
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"question_options"}, []string{"question_id", "hero_id", "is_correct", "telegram_file_id"}, pgx.CopyFromSlice(len(q.Options), func(i int) ([]any, error) {
		return []any{qID, q.Options[i].Hero.ID, q.Options[i].IsCorrect, q.Options[i].TelegramFileID}, nil
	}))
	if err != nil {
		return qID, err
	}

	return qID, nil
}

func (r *QuestionRepository) AnswerQuestion(ctx context.Context, u *data.User, q *data.Question, answer *data.UserOption) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err = transaction(ctx, tx, "AnswerQuestion", func() error {
		return r.answerQuestionTx(ctx, tx, u, q, answer)
	}); err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) answerQuestionTx(ctx context.Context, tx pgx.Tx, u *data.User, q *data.Question, answer *data.UserOption) error {
	const sql = `INSERT INTO user_questions (user_question_id, user_id, question_id, hero_id, answered_at) VALUES
	($1, $2, $3, $4, $5)`

	_, err := tx.Exec(ctx, sql, answer.ID, u.ID, q.ID, answer.Hero.ID, answer.AnsweredAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) UpdateQuestionImage(ctx context.Context, q *data.Question, fileID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err = transaction(ctx, tx, "UpdateQuestionImage", func() error {
		return r.updateQuestionImageTx(ctx, tx, q, fileID)
	}); err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) updateQuestionImageTx(ctx context.Context, tx pgx.Tx, q *data.Question, fileID string) error {
	const sql = `UPDATE questions SET telegram_file_id = $1 WHERE question_id = $2`

	_, err := tx.Exec(ctx, sql,fileID, q.ID)
	if err != nil {
		return err
	}

	return nil
}
