package postgres

import (
	"context"
	"log/slog"

	"github.com/akionka/akionkabot/internal/data"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewQuestionRepository(db *pgxpool.Pool, logger *slog.Logger) *QuestionRepository {
	return &QuestionRepository{
		db:     db,
		logger: logger,
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
	r.logger.DebugContext(ctx, "getting question by id", slog.String("uuid", id.String()))

	rows, err := tx.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.Question])
}

func (r *QuestionRepository) GetQuestionAvailableForUser(ctx context.Context, userID uuid.UUID, isWon bool) (*data.Question, error) {
	var question *data.Question

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return question, err
	}

	if err = transaction(ctx, tx, "GetQuestionAvailableForUser", func() error {
		question, err = r.getQuestionAvailableForUserTx(ctx, tx, userID, isWon)
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
	LEFT JOIN user_answers ua ON q.question_id = ua.question_id AND ua.user_id = $1
	WHERE q.is_won = $2 AND ua.question_id IS NULL
	ORDER BY q.created_at DESC
	LIMIT 1;`
	r.logger.DebugContext(ctx, "getting question available for user", slog.String("user_uuid", userID.String()), slog.Bool("is_won", isWon))

	rows, err := tx.Query(ctx, sql, userID, isWon)
	if err != nil {
		return nil, err
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.Question])
}

func (r *QuestionRepository) enrichQuestionTx(ctx context.Context, tx pgx.Tx, question *data.Question) (*data.Question, error) {
	const (
		heroSQL = `
		SELECT qo.is_correct, h.hero_id, h.display_name, h.short_name, qo.telegram_file_id
		FROM question_options qo
		INNER JOIN heroes h ON qo.hero_id = h.hero_id
		WHERE qo.question_id = $1
		ORDER BY qo.order`

		itemSQL = `
		SELECT i.item_id, i.display_name, i.short_name
		FROM question_items qi
		INNER JOIN items i ON qi.item_id = i.item_id
		WHERE qi.question_id = $1`
	)
	r.logger.DebugContext(ctx, "enriching question", slog.Any("question", question))

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
	slog.DebugContext(ctx, "saving question", slog.Any("question", question))

	var questionID uuid.UUID

	if err := tx.QueryRow(ctx, sql,
		question.ID, question.MatchID, question.MatchStartedAt, question.PlayerID, question.PlayerName, question.PlayerIsPro, question.PlayerPos, question.PlayerMMR, question.IsWon, question.CreatedAt, question.TelegramFileID,
	).Scan(&questionID); err != nil {
		return questionID, err
	}

	if _, err := tx.CopyFrom(ctx, pgx.Identifier{"question_items"}, []string{"question_id", "item_id"}, pgx.CopyFromSlice(len(question.Items), func(i int) ([]any, error) {
		return []any{questionID, question.Items[i].ID}, nil
	})); err != nil {
		return questionID, err
	}

	if _, err := tx.CopyFrom(ctx, pgx.Identifier{"question_options"}, []string{"question_id", "hero_id", "is_correct", "telegram_file_id", "order"}, pgx.CopyFromSlice(len(question.Options), func(i int) ([]any, error) {
		return []any{questionID, question.Options[i].Hero.ID, question.Options[i].IsCorrect, question.Options[i].TelegramFileID, i}, nil
	})); err != nil {
		return questionID, err
	}

	return questionID, nil
}

func (r *QuestionRepository) AnswerQuestion(ctx context.Context, userID uuid.UUID, question *data.Question, answer data.UserAnswer) (data.UserAnswer, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return answer, err
	}

	if err = transaction(ctx, tx, "AnswerQuestion", func() error {
		if _, err := r.answerQuestionTx(ctx, tx, userID, question, answer); err != nil {
			return err
		}
		answer, err = r.getUserAnswerTx(ctx, tx, question.ID, userID)
		return err
	}); err != nil {
		return answer, err
	}

	return answer, nil
}

func (r *QuestionRepository) answerQuestionTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID, question *data.Question, answer data.UserAnswer) (uuid.UUID, error) {
	const sql = `INSERT INTO user_answers (user_answer_id, user_id, question_id, hero_id, answered_at) VALUES
	($1, $2, $3, $4, $5) RETURNING user_answer_id`
	slog.DebugContext(ctx, "saving user's answer", slog.String("user_uuid", userID.String()), slog.Any("question", question), slog.Any("answer", answer))

	var userQuestionID uuid.UUID
	if err := tx.QueryRow(ctx, sql, answer.ID, userID, question.ID, answer.Hero.ID, answer.AnsweredAt).Scan(&userQuestionID); err != nil {
		return userQuestionID, err
	}

	return userQuestionID, nil
}

func (r *QuestionRepository) GetUserAnswer(ctx context.Context, id uuid.UUID, userID uuid.UUID) (data.UserAnswer, error) {
	var answer data.UserAnswer

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

func (r *QuestionRepository) getUserAnswerTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, userID uuid.UUID) (data.UserAnswer, error) {
	const sql = `
	SELECT ua.user_answer_id, ua.answered_at, qo.is_correct, qo.telegram_file_id, h.hero_id, h.display_name, h.short_name
	FROM user_answers ua
	INNER JOIN public.question_options qo on ua.hero_id = qo.hero_id AND ua.question_id = qo.question_id
	INNER JOIN public.heroes h on ua.hero_id = h.hero_id
	WHERE ua.question_id = $1 AND ua.user_id = $2`
	slog.DebugContext(ctx, "getting user's answer to question", slog.String("question_uuid", id.String()), slog.String("user_id", id.String()))

	var answer data.UserAnswer

	rows, err := tx.Query(ctx, sql, id, userID)
	if err != nil {
		return answer, err
	}

	return pgx.CollectOneRow(rows, pgx.RowToStructByName[data.UserAnswer])
}

func (r *QuestionRepository) UpdateQuestionImage(ctx context.Context, id uuid.UUID, fileID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err = transaction(ctx, tx, "UpdateQuestionImage", func() error {
		return r.updateQuestionImageTx(ctx, tx, id, fileID)
	}); err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) updateQuestionImageTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, fileID string) error {
	const sql = `UPDATE questions SET telegram_file_id = $1 WHERE question_id = $2`
	slog.DebugContext(ctx, "updating question image", slog.String("question_uuid", id.String()), slog.String("file_id", fileID))

	_, err := tx.Exec(ctx, sql, fileID, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) UpdateOptionImage(ctx context.Context, id uuid.UUID, option data.Option, fileID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	if err = transaction(ctx, tx, "UpdateOptionImage", func() error {
		return r.updateOptionImageTx(ctx, tx, id, option, fileID)
	}); err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) updateOptionImageTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, option data.Option, fileID string) error {
	const sql = `UPDATE question_options SET telegram_file_id = $1 WHERE question_id = $2 AND hero_id = $3`
	slog.DebugContext(ctx, "updating option image", slog.String("question_uuid", id.String()), slog.String("file_id", fileID), slog.Any("option", option))

	_, err := tx.Exec(ctx, sql, fileID, id, option.Hero.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *QuestionRepository) GetQuestionStats(ctx context.Context, questionID uuid.UUID) (map[int]int, error) {
	const sql = `
	SELECT qo.hero_id, COUNT(ua.user_id) as answer_count
	FROM question_options qo
	LEFT JOIN user_answers ua ON qo.question_id = ua.question_id AND qo.hero_id = ua.hero_id
	WHERE qo.question_id = $1
	GROUP BY qo.hero_id`
	slog.DebugContext(ctx, "getting question stats", slog.String("question_uuid", questionID.String()))

	rows, err := r.db.Query(ctx, sql, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[int]int)
	for rows.Next() {
		var heroID int
		var count int
		if err := rows.Scan(&heroID, &count); err != nil {
			return nil, err
		}
		stats[heroID] = count
	}

	return stats, nil
}
