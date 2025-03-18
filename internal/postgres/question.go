package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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

func (r *QuestionRepository) GetQuestion(ctx context.Context, id data.QuestionID) (*data.Question, error) {
	const sql = `
		SELECT
			q.question_id, q.match_id, q.player_steam_id, q.created_at, q.telegram_file_id AS q_telegram_file_id,
			mp.hero_id = qo.hero_id AS is_correct, qo.telegram_file_id AS qo_telegram_file_id, qo."order",
			h.hero_id, h.display_name, h.short_name

		FROM		questions			q
		LEFT JOIN	question_options	qo	ON qo.question_id = q.question_id
		INNER JOIN	match_players		mp	ON mp.player_steam_id = q.player_steam_id
		INNER JOIN	heroes				h	ON qo.hero_id = h.hero_id

		WHERE q.question_id = $1
		ORDER BY qo."order"
	`

	r.logger.DebugContext(ctx, "getting question by id", slog.String("id", id.String()))
	rows, err := r.db.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("error getting question by id: %w", err)
	}
	defer rows.Close()

	question, err := r.scanQuestionFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning question: %w", err)
	}

	return question, nil
}

func (r *QuestionRepository) GetQuestionAvailableForUser(ctx context.Context, userID data.UserID, isWon bool) (*data.Question, error) {
	const sql = `
		SELECT
			q.question_id, q.match_id, q.player_steam_id, q.created_at, q.telegram_file_id AS q_telegram_file_id,
			mp.hero_id = qo.hero_id AS is_correct, qo.telegram_file_id AS qo_telegram_file_id, qo."order",
			h.hero_id, h.display_name, h.short_name

		FROM (
			SELECT q.*
			FROM 		questions		q
			INNER JOIN	matches			m	ON	q.match_id = m.match_id
			INNER JOIN	match_players	mp	ON	m.match_id = mp.match_id AND q.player_steam_id = mp.player_steam_id
			LEFT JOIN	user_answers	ua	ON	q.question_id = ua.question_id
											AND	ua.user_id = $1
			WHERE	ua.user_answer_id IS NULL
			AND		(m.winning_team = mp.team) = $2
			ORDER BY q.created_at DESC
			LIMIT 1) q

		LEFT JOIN	question_options	qo	ON qo.question_id = q.question_id
		INNER JOIN match_players		mp	ON mp.player_steam_id = q.player_steam_id
		INNER JOIN	heroes				h	ON qo.hero_id = h.hero_id

		ORDER BY qo."order"
	`

	r.logger.DebugContext(ctx, "getting question available for user", slog.String("user_id", userID.String()), slog.Bool("is_won", isWon))
	rows, err := r.db.Query(ctx, sql, userID, isWon)
	if err != nil {
		return nil, fmt.Errorf("error getting question available for user: %w", err)
	}

	question, err := r.scanQuestionFromRows(rows)
	if err != nil {
		return nil, fmt.Errorf("error scanning question: %w", err)
	}

	return question, nil
}

func (r *QuestionRepository) scanQuestionFromRows(rows pgx.Rows) (*data.Question, error) {
	type questionRow struct {
		QuestionID     uuid.UUID `db:"question_id"`
		MatchID        int64     `db:"match_id"`
		PlayerID       int64     `db:"player_steam_id"`
		CreatedAt      time.Time `db:"created_at"`
		TelegramFileID string    `db:"q_telegram_file_id"`

		OptionIsCorrect      *bool   `db:"is_correct"`
		OptionTelegramFileID *string `db:"qo_telegram_file_id"`
		OptionOrder          *int    `db:"order"`

		HeroID          *int    `db:"hero_id"`
		HeroShortName   *string `db:"short_name"`
		HeroDisplayName *string `db:"display_name"`
	}

	var question *data.Question
	options := make(map[int]*data.Option, maxOptionCount) // Key: Order

	qRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[questionRow])
	if err != nil {
		return nil, fmt.Errorf("error collecting question rows: %w", err)
	}
	if len(qRows) == 0 {
		return nil, data.ErrNotFound
	}

	for _, row := range qRows {
		if question == nil {
			question = &data.Question{
				ID:             data.QuestionID(row.QuestionID),
				MatchID:        data.MatchID(row.MatchID),
				PlayerID:       data.SteamID(row.PlayerID),
				TelegramFileID: row.TelegramFileID,
				CreatedAt:      row.CreatedAt,
				Options:        make([]data.Option, 0, maxOptionCount),
			}
		}

		if row.OptionOrder != nil && options[*row.OptionOrder] == nil {
			OptionTelegramFileID := ""
			if row.OptionTelegramFileID != nil {
				OptionTelegramFileID = *row.OptionTelegramFileID
			}

			option := data.Option{
				Hero: data.Hero{
					ID:          data.HeroID(*row.HeroID),
					DisplayName: *row.HeroDisplayName,
					ShortName:   *row.HeroShortName,
				},
				IsCorrect:      *row.OptionIsCorrect,
				TelegramFileID: OptionTelegramFileID,
			}
			options[*row.OptionOrder] = &option
			question.Options = append(question.Options, option)
		}
	}

	return question, nil
}

func (r *QuestionRepository) SaveQuestion(ctx context.Context, q *data.Question) (questionID data.QuestionID, err error) {
	const insertQuestion = `INSERT INTO questions (question_id, match_id, player_steam_id, created_at, telegram_file_id) VALUES ($1, $2, $3, $4, $5) RETURNING question_id`
	const insertOption = `INSERT INTO question_options (question_id, hero_id, "order", telegram_file_id) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`

	r.logger.DebugContext(ctx, "saving question", slog.Any("question", q))
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return questionID, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				err = fmt.Errorf("error rolling back transaction: %w", rbErr)
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			err = fmt.Errorf("error committing transaction: %w", err)
		}
	}()

	if err = tx.QueryRow(ctx, insertQuestion, q.ID, q.MatchID, q.PlayerID, q.CreatedAt, q.TelegramFileID).Scan(&questionID); err != nil {
		err = pgErrToDomain(err)
		return questionID, fmt.Errorf("error inserting question: %w", err)
	}

	for i, option := range q.Options {
		if _, err = tx.Exec(ctx, insertOption, questionID, option.Hero.ID, i, option.TelegramFileID); err != nil {
			return questionID, fmt.Errorf("error inserting option: %w", err)
		}
	}

	return questionID, nil
}

func (r *QuestionRepository) AnswerQuestion(ctx context.Context, id data.QuestionID, userID data.UserID, answer data.UserAnswer) (data.UserAnswerID, error) {
	const sql = `INSERT INTO user_answers (user_answer_id, user_id, question_id, hero_id, answered_at) VALUES
		($1, $2, $3, $4, $5) RETURNING user_answer_id
	`
	var userAnswerID data.UserAnswerID

	slog.DebugContext(ctx, "saving user's answer", slog.String("user_id", userID.String()), slog.String("question_id", id.String()), slog.Any("answer", answer))
	if err := r.db.QueryRow(ctx, sql, answer.ID, userID, id, answer.Hero.ID, answer.AnsweredAt).Scan(&userAnswerID); err != nil {
		err = pgErrToDomain(err)
		return data.UserAnswerID(uuid.Nil), fmt.Errorf("error inserting user's answer: %w", err)
	}

	return userAnswerID, nil
}

func (r *QuestionRepository) GetUserAnswer(ctx context.Context, id data.QuestionID, userID data.UserID) (*data.UserAnswer, error) {
	const sql = `
		SELECT ua.user_answer_id, ua.answered_at, qo.hero_id = mp.hero_id AS is_correct, qo.telegram_file_id, h.hero_id, h.display_name, h.short_name

		FROM		questions			q
		INNER JOIN  match_players		mp	ON	q.match_id = mp.match_id
											AND	q.player_steam_id = mp.player_steam_id
		INNER JOIN	question_options	qo	ON	q.question_id = qo.question_id
		INNER JOIN 	user_answers		ua	ON	qo.question_id = ua.question_id AND qo.hero_id = ua.hero_id
		INNER JOIN	heroes				h	ON	qo.hero_id = h.hero_id

		WHERE	q.question_id = $1
		AND		ua.user_id = $2
		ORDER BY qo."order"
	`
	var answer data.UserAnswer

	slog.DebugContext(ctx, "getting user's answer to question", slog.String("question_id", id.String()), slog.String("user_id", id.String()))
	rows, err := r.db.Query(ctx, sql, id, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user's answer: %w", err)
	}

	answer, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[data.UserAnswer])
	if err != nil {
		err = pgErrToDomain(err)
		return nil, fmt.Errorf("error collecting user's answer: %w", err)
	}

	return &answer, nil
}

func (r *QuestionRepository) UpdateQuestionImage(ctx context.Context, id data.QuestionID, fileID string) error {
	const sql = `UPDATE questions SET telegram_file_id = $1 WHERE question_id = $2`

	slog.DebugContext(ctx, "updating question image", slog.String("question_id", id.String()), slog.String("file_id", fileID))
	_, err := r.db.Exec(ctx, sql, fileID, id)
	if err != nil {
		return fmt.Errorf("error updating question image: %w", err)
	}

	return nil
}

func (r *QuestionRepository) UpdateOptionImage(ctx context.Context, id data.QuestionID, option data.Option, fileID string) error {
	const sql = `UPDATE question_options SET telegram_file_id = $1 WHERE question_id = $2 AND hero_id = $3`

	slog.DebugContext(ctx, "updating option image", slog.String("question_id", id.String()), slog.String("file_id", fileID), slog.Any("option", option))
	_, err := r.db.Exec(ctx, sql, fileID, id, option.Hero.ID)
	if err != nil {
		return fmt.Errorf("error updating option image: %w", err)
	}

	return nil
}

func (r *QuestionRepository) GetQuestionStats(ctx context.Context, id data.QuestionID) (map[data.HeroID]int, error) {
	const sql = `
		SELECT qo.hero_id, COUNT(ua.user_id) as answer_count
		FROM question_options qo
		LEFT JOIN user_answers ua ON qo.question_id = ua.question_id AND qo.hero_id = ua.hero_id
		WHERE qo.question_id = $1
		GROUP BY qo.hero_id
	`

	slog.DebugContext(ctx, "getting question stats", slog.String("question_id", id.String()))
	rows, err := r.db.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("error getting question stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[data.HeroID]int)
	for rows.Next() {
		var heroID data.HeroID
		var count int
		if err := rows.Scan(&heroID, &count); err != nil {
			return nil, err
		}
		stats[heroID] = count
	}

	return stats, nil
}
