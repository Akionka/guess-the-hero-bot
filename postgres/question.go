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
	SELECT q.question_id, q.match_id, q.match_started_at, q.player_id, q.player_name, q.player_is_pro, q.player_pos, q.player_mmr, q.is_won, q.created_at, q.telegram_file_id
	FROM questions q
	WHERE q.question_id = $1
	LIMIT 1`

	rows, err := tx.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.Question])
}

func (r *QuestionRepository) GetQuestionAvailableForUser(ctx context.Context, id uuid.UUID, isWon bool) (*data.Question, error) {
	var question *data.Question

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return question, err
	}

	if err = transaction(ctx, tx, "GetQuestionAvailableForUser", func() error {
		question, err = r.getQuestionAvailableForUserTx(ctx, tx, id, isWon)
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

func (r *QuestionRepository) getQuestionAvailableForUserTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID, isWon bool) (*data.Question, error) {
	const sql = `
	SELECT q.question_id, q.match_id, q.match_started_at, q.player_id, q.player_name, q.player_is_pro, q.player_pos, q.player_mmr, q.is_won, q.created_at, q.telegram_file_id
	FROM questions q
	LEFT JOIN user_questions uq ON q.question_id = uq.question_id AND uq.user_id = $1
	WHERE q.is_won = $2 AND uq.question_id IS NULL
	ORDER BY q.created_at DESC
	LIMIT 1;`

	rows, err := tx.Query(ctx, sql, userID, isWon)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.Question])
}

func (r *QuestionRepository) enrichQuestionTx(ctx context.Context, tx pgx.Tx, question *data.Question) (*data.Question, error) {
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

	heroRows, err := tx.Query(ctx, heroSQL, question.ID)
	if err != nil {
		return question, err
	}

	question.Options, err = pgx.CollectRows(heroRows, pgx.RowToStructByName[data.Option])
	if err != nil {
		return question, err
	}

	itemRows, err := tx.Query(ctx, itemSQL, question.ID)
	if err != nil {
		return question, err
	}

	question.Items, err = pgx.CollectRows(itemRows, pgx.RowToStructByName[data.Item])
	return question, err
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

func (r *QuestionRepository) saveQuestionTx(ctx context.Context, tx pgx.Tx, question *data.Question) (uuid.UUID, error) {
	const sql = `
	INSERT INTO questions (question_id, match_id, match_started_at, player_id, player_name, player_is_pro, player_pos, player_mmr, is_won, created_at, telegram_file_id) VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING question_id`

	qID := uuid.Nil

	rows, err := tx.Query(ctx, sql,
		question.ID, question.MatchID, question.MatchStartedAt, question.PlayerID, question.PlayerName, question.PlayerIsPro, question.PlayerPos, question.PlayerMMR, question.IsWon, question.CreatedAt, question.TelegramFileID,
	)
	if err != nil {
		return qID, err
	}
	qID, err = pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return qID, err
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"question_items"}, []string{"question_id", "item_id"}, pgx.CopyFromSlice(len(question.Items), func(i int) ([]any, error) {
		return []any{qID, question.Items[i].ID}, nil
	}))
	if err != nil {
		return qID, err
	}

	_, err = tx.CopyFrom(ctx, pgx.Identifier{"question_options"}, []string{"question_id", "hero_id", "is_correct", "telegram_file_id"}, pgx.CopyFromSlice(len(question.Options), func(i int) ([]any, error) {
		return []any{qID, question.Options[i].Hero.ID, question.Options[i].IsCorrect, question.Options[i].TelegramFileID}, nil
	}))
	if err != nil {
		return qID, err
	}

	return qID, nil
}

func (r *QuestionRepository) AnswerQuestion(ctx context.Context, user *data.User, question *data.Question, answer *data.UserOption) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err = transaction(ctx, tx, "AnswerQuestion", func() error {
		return r.answerQuestionTx(ctx, tx, user, question, answer)
	}); err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) answerQuestionTx(ctx context.Context, tx pgx.Tx, user *data.User, question *data.Question, answer *data.UserOption) error {
	const sql = `INSERT INTO user_questions (user_question_id, user_id, question_id, hero_id, answered_at) VALUES
	($1, $2, $3, $4, $5)`

	_, err := tx.Exec(ctx, sql, answer.ID, user.ID, question.ID, answer.Hero.ID, answer.AnsweredAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) GetUserAnswer(ctx context.Context, id uuid.UUID, userID uuid.UUID) (data.UserOption, error) {
	var answer data.UserOption

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return answer, err
	}

	if err = transaction(ctx, tx, "GetUserAnswer", func() error {
		answer, err = r.getUserAnswerTx(ctx, tx, id, userID)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return answer, err
	}

	return answer, nil
}

func (r *QuestionRepository) getUserAnswerTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, userID uuid.UUID) (data.UserOption, error) {
	const sql = `
	SELECT uq.user_question_id, uq.answered_at, qo.is_correct, qo.telegram_file_id, h.hero_id, h.display_name, h.short_name
	FROM user_questions uq
	INNER JOIN public.question_options qo on uq.hero_id = qo.hero_id AND uq.question_id = qo.question_id
	INNER JOIN public.heroes h on uq.hero_id = h.hero_id
	WHERE uq.question_id = $1 AND uq.user_id = $2`

	var answer data.UserOption

	rows, err := tx.Query(ctx, sql, id, userID)
	if err != nil {
		return answer, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[data.UserOption])
}

func (r *QuestionRepository) UpdateQuestionImage(ctx context.Context, question *data.Question, fileID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err = transaction(ctx, tx, "UpdateQuestionImage", func() error {
		return r.updateQuestionImageTx(ctx, tx, question, fileID)
	}); err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) updateQuestionImageTx(ctx context.Context, tx pgx.Tx, question *data.Question, fileID string) error {
	const sql = `UPDATE questions SET telegram_file_id = $1 WHERE question_id = $2`

	_, err := tx.Exec(ctx, sql, fileID, question.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) UpdateOptionImage(ctx context.Context, question *data.Question, option *data.Option, fileID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err = transaction(ctx, tx, "UpdateOptionImage", func() error {
		return r.updateOptionImageTx(ctx, tx, question, option, fileID)
	}); err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) updateOptionImageTx(ctx context.Context, tx pgx.Tx, q *data.Question, o *data.Option, fileID string) error {
	const sql = `UPDATE question_options SET telegram_file_id = $1 WHERE question_id = $2 AND hero_id = $3`

	_, err := tx.Exec(ctx, sql, fileID, q.ID, o.Hero.ID)
	if err != nil {
		return err
	}

	return nil
}
