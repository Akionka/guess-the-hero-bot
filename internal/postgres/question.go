package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"slices"
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

func (r *QuestionRepository) GetQuestion(ctx context.Context, id uuid.UUID) (*data.Question, error) {
	const sql = `
		SELECT
			-- Question fields
			q.question_id AS q_id, q.created_at AS q_created_at, q.telegram_file_id AS q_telegram_file_id, q.player_steam_id AS q_p_id,

			-- Question.Match fields
			m.match_id AS m_id, m.radiant_won AS m_radiant_won, m.actual_rank AS m_actual_rank, m.started_at AS m_started_at, m.avg_mmr AS m_avg_mmr,

			-- Question.Match.MatchPlayer fields
			mp.player_steam_id AS mp_id, mp.hero_id AS mp_hero_id, mp.is_radiant mp_id_radiant, mp.position mp_position,

			-- Question.Match.MatchPlayer.Player fields
			pa.is_pro AS p_is_pro, pa.name AS p_name, pa.pro_name AS p_pro_name,
			h.display_name AS h_display_name, h.short_name AS h_short_name,

			-- Question.Match.MatchPlayer.Items fields
			i.item_id AS i_id, i.display_name AS i_display_name, i.short_name AS i_short_name, mpi."order" AS i_order,

			-- Question.Options fields
			qo.hero_id AS o_hero_id, qo."order" AS o_order, qo.telegram_file_id AS o_telegram_file_id,
			oh.display_name AS o_h_display_name, oh.short_name AS o_h_short_name

		FROM        questions               q                                           	    -- Question
		LEFT JOIN   matches                 m   ON  q.match_id = m.match_id               		-- Question.Match
		LEFT JOIN   match_players           mp  ON  m.match_id = mp.match_id                	-- Question.Match.Players
		LEFT JOIN   player_accounts         pa  ON  mp.player_steam_id = pa.player_steam_id     -- Question.Match.Players.Player
		LEFT JOIN   heroes                  h   ON 	mp.hero_id = h.hero_id                      -- Question.Match.Players.Hero
		LEFT JOIN   match_player_items      mpi ON  mp.player_steam_id = mpi.player_steam_id    -- Question.Match.Players.Items
												AND mp.match_id = mpi.match_id
		LEFT JOIN   items                   i   ON  mpi.item_id = i.item_id                     -- Question.Match.Players.Items
		LEFT JOIN   question_options        qo  ON  q.question_id = qo.question_id              -- Question.Options
		LEFT JOIN   heroes                  oh  ON  qo.hero_id = oh.hero_id                     -- Question.Options.Hero

		WHERE q.question_id = $1`

	r.logger.DebugContext(ctx, "getting question by id", slog.String("uuid", id.String()))
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

func (r *QuestionRepository) GetQuestionAvailableForUser(ctx context.Context, userID uuid.UUID, isWon bool) (*data.Question, error) {
	const sql = `
		SELECT
			-- Question fields
			q.question_id AS q_id, q.created_at AS q_created_at, q.telegram_file_id AS q_telegram_file_id, q.player_steam_id AS q_p_id,

			-- Question.Match fields
			m.match_id AS m_id, m.radiant_won AS m_radiant_won, m.actual_rank AS m_actual_rank, m.started_at AS m_started_at, m.avg_mmr AS m_avg_mmr,

			-- Question.Match.MatchPlayer fields
			mp.player_steam_id AS mp_id, mp.hero_id AS mp_hero_id, mp.is_radiant mp_id_radiant, mp.position mp_position,

			-- Question.Match.MatchPlayer.Player fields
			pa.is_pro AS p_is_pro, pa.name AS p_name, pa.pro_name AS p_pro_name,
			h.display_name AS h_display_name, h.short_name AS h_short_name,

			-- Question.Match.MatchPlayer.Items fields
			i.item_id AS i_id, i.display_name AS i_display_name, i.short_name AS i_short_name, mpi."order" AS i_order,

			-- Question.Options fields
			qo.hero_id AS o_hero_id, qo."order" AS o_order, qo.telegram_file_id AS o_telegram_file_id,
			oh.display_name AS o_h_display_name, oh.short_name AS o_h_short_name

		FROM (
			SELECT q.*
            FROM 		questions               q                                           	    -- Question
            LEFT JOIN 	user_answers			ua	ON	q.question_id = ua.question_id
                                                AND ua.user_id = $1
            LEFT JOIN   matches                 m   ON  q.match_id = m.match_id               		-- Question.Match
            LEFT JOIN   match_players           mp  ON  m.match_id = mp.match_id                	-- Question.Match.Players
                                                    AND (mp.is_radiant = m.radiant_won) = $2
            WHERE ua.question_id IS NULL
            GROUP BY q.question_id, q.created_at
            ORDER BY q.created_at DESC
            LIMIT 1
		) 									q
		LEFT JOIN   matches                 m   ON  q.match_id = m.match_id               		-- Question.Match
		LEFT JOIN   match_players           mp  ON  m.match_id = mp.match_id                	-- Question.Match.Players
		LEFT JOIN   player_accounts         pa  ON  mp.player_steam_id = pa.player_steam_id     -- Question.Match.Players.Player
		LEFT JOIN   heroes                  h   ON 	mp.hero_id = h.hero_id                      -- Question.Match.Players.Hero
		LEFT JOIN   match_player_items      mpi ON  mp.player_steam_id = mpi.player_steam_id    -- Question.Match.Players.Items
												AND mp.match_id = mpi.match_id
		LEFT JOIN   items                   i   ON  mpi.item_id = i.item_id                     -- Question.Match.Players.Items
		LEFT JOIN   question_options        qo  ON  q.question_id = qo.question_id              -- Question.Options
		LEFT JOIN   heroes                  oh  ON  qo.hero_id = oh.hero_id                     -- Question.Options.Hero
	`

	r.logger.DebugContext(ctx, "getting question available for user", slog.String("user_uuid", userID.String()), slog.Bool("is_won", isWon))
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
		QID             uuid.UUID `db:"q_id"`
		QCreatedAt      time.Time `db:"q_created_at"`
		QTelegramFileID string    `db:"q_telegram_file_id"`

		// Question Player (one of match players)
		QPlayerSteamID int64 `db:"q_p_id"`

		MID         int64     `db:"m_id"`
		MRadiantWon bool      `db:"m_radiant_won"`
		MStartedAt  time.Time `db:"m_started_at"`
		MAvgMMR     *int      `db:"m_avg_mmr"`
		MActualRank int       `db:"m_actual_rank"`

		// Match Player
		PPlayerSteamID *int64         `db:"mp_id"`
		PHeroID        *int           `db:"mp_hero_id"`
		PIsRadiant     *bool          `db:"mp_id_radiant"`
		PPosition      *data.Position `db:"mp_position"`

		PIsPro   *bool   `db:"p_is_pro"`
		PName    *string `db:"p_name"`
		PProName *string `db:"p_pro_name"`

		// Player Hero
		PHDisplayName *string `db:"h_display_name"`
		PHShortName   *string `db:"h_short_name"`

		// Match Player Items
		PItemID       *int    `db:"i_id"`
		PIDisplayName *string `db:"i_display_name"`
		PIShortName   *string `db:"i_short_name"`
		PIOrder       *int    `db:"i_order"`

		// Question Options
		OHeroID         *int    `db:"o_hero_id"`
		OOrder          *int    `db:"o_order"`
		OTelegramFileID *string `db:"o_telegram_file_id"`

		// Option Hero
		OHDisplayName *string `db:"o_h_display_name"`
		OHShortName   *string `db:"o_h_short_name"`
	}

	const (
		itemCount   = 6
		optionCount = 4
	)

	var question *data.Question
	playerMap := make(map[int64]*data.MatchPlayer)     // Key: Player.SteamID
	itemMap := make(map[string]*data.Item)             // Key: Player.SteamID-Order
	options := make(map[int]*data.Option, optionCount) // Key: Order

	qRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[questionRow])
	if err != nil {
		return nil, fmt.Errorf("error collecting question rows: %w", err)
	}
	if len(qRows) == 0 {
		return nil, data.ErrNotFound
	}

	for _, q := range qRows {
		if question == nil {
			question = &data.Question{
				ID: q.QID,
				Match: data.Match{
					ID:         q.MID,
					RadiantWon: q.MRadiantWon,
					ActualRank: q.MActualRank,
					StartedAt:  q.MStartedAt,
					AvgMMR:     q.MAvgMMR,
					Players:    []data.MatchPlayer{},
				},
				Player:         nil,
				Options:        []data.Option{},
				TelegramFileID: q.QTelegramFileID,
				CreatedAt:      q.QCreatedAt,
			}
		}

		if q.OOrder != nil && options[*q.OOrder] == nil {
			option := data.Option{
				Hero: data.Hero{
					ID:          *q.OHeroID,
					DisplayName: *q.OHDisplayName,
					ShortName:   *q.OHShortName,
					Image:       data.Image{Image: nil},
				},
				IsCorrect:      false,
				TelegramFileID: *q.OTelegramFileID,
			}
			options[*q.OOrder] = &option
		}

		if q.PPlayerSteamID != nil && playerMap[*q.PPlayerSteamID] == nil {
			player := data.MatchPlayer{
				Player: data.Player{
					SteamID: *q.PPlayerSteamID,
					Name:    *q.PName,
					IsPro:   *q.PIsPro,
					ProName: *q.PProName,
				},
				Hero: data.Hero{
					ID:          *q.PHeroID,
					DisplayName: *q.PHDisplayName,
					ShortName:   *q.PHShortName,
					Image:       data.Image{Image: nil},
				},
				IsRadiant: *q.PIsRadiant,
				Position:  *q.PPosition,
				Items:     make([]data.Item, itemCount),
			}
			playerMap[*q.PPlayerSteamID] = &player
			question.Match.Players = append(question.Match.Players, player)

			if player.Player.SteamID == q.QPlayerSteamID {
				question.Player = &player
			}
		}

		if q.PPlayerSteamID != nil && q.PIOrder != nil {
			key := fmt.Sprintf("%d-%d", *q.PPlayerSteamID, *q.PIOrder)
			if itemMap[key] == nil {
				item := data.Item{
					ID:          *q.PItemID,
					DisplayName: *q.PIDisplayName,
					ShortName:   *q.PIShortName,
				}
				if _, ok := itemMap[key]; !ok {
					itemMap[key] = &item
				}
				playerMap[*q.PPlayerSteamID].Items[*q.PIOrder] = item
			}
		}
	}

	optionOrders := make([]int, 0, optionCount)
	for k := range maps.Keys(options) {
		optionOrders = append(optionOrders, k)
	}
	slices.Sort(optionOrders)

	for _, order := range optionOrders {
		question.Options = append(question.Options, *options[order])
	}

	for i := range question.Options {
		question.Options[i].IsCorrect = question.Options[i].Hero.ID == question.Player.Hero.ID
	}

	return question, nil
}

func (r *QuestionRepository) SaveQuestion(ctx context.Context, q *data.Question) (questionID uuid.UUID, err error) {
	const insertMatch = `INSERT INTO matches (match_id, radiant_won, actual_rank, started_at, avg_mmr) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	const insertMatchPlayer = `INSERT INTO match_players (match_id, player_steam_id, hero_id, is_radiant, position) VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`
	const insertPlayer = `INSERT INTO player_accounts (player_steam_id, is_pro, name, pro_name) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
	const insertItem = `INSERT INTO match_player_items (player_steam_id, match_id, item_id, "order") VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING`
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

	if _, err = tx.Exec(ctx, insertMatch, q.Match.ID, q.Match.RadiantWon, q.Match.ActualRank, q.Match.StartedAt, q.Match.AvgMMR); err != nil {
		return questionID, fmt.Errorf("error inserting match: %w", err)
	}

	for _, player := range q.Match.Players {
		if _, err = tx.Exec(ctx, insertMatchPlayer, q.Match.ID, player.Player.SteamID, player.Hero.ID, player.IsRadiant, player.Position); err != nil {
			return questionID, fmt.Errorf("error inserting match player: %w", err)
		}
		if _, err = tx.Exec(ctx, insertPlayer, player.Player.SteamID, player.Player.IsPro, player.Player.Name, player.Player.ProName); err != nil {
			return questionID, fmt.Errorf("error inserting player: %w", err)
		}
		for order, item := range player.Items {
			if _, err = tx.Exec(ctx, insertItem, player.Player.SteamID, q.Match.ID, item.ID, order); err != nil {
				return questionID, fmt.Errorf("error inserting item: %w", err)
			}
		}
	}

	if err = tx.QueryRow(ctx, insertQuestion, q.ID, q.Match.ID, q.Player.Player.SteamID, q.CreatedAt, q.TelegramFileID).Scan(&questionID); err != nil {
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

func (r *QuestionRepository) AnswerQuestion(ctx context.Context, userID uuid.UUID, question *data.Question, answer data.UserAnswer) error {
	const sql = `INSERT INTO user_answers (user_answer_id, user_id, question_id, hero_id, answered_at) VALUES
		($1, $2, $3, $4, $5) RETURNING user_answer_id`
	slog.DebugContext(ctx, "saving user's answer", slog.String("user_uuid", userID.String()), slog.Any("question", question), slog.Any("answer", answer))

	var userQuestionID uuid.UUID
	if err := r.db.QueryRow(ctx, sql, answer.ID, userID, question.ID, answer.Hero.ID, answer.AnsweredAt).Scan(&userQuestionID); err != nil {
		return fmt.Errorf("error inserting user's answer: %w", err)
	}

	return nil
}

func (r *QuestionRepository) GetUserAnswer(ctx context.Context, id uuid.UUID, userID uuid.UUID) (data.UserAnswer, error) {
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

	slog.DebugContext(ctx, "getting user's answer to question", slog.String("question_uuid", id.String()), slog.String("user_id", id.String()))
	rows, err := r.db.Query(ctx, sql, id, userID)
	if err != nil {
		return answer, fmt.Errorf("error getting user's answer: %w", err)
	}

	answer, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[data.UserAnswer])
	if err != nil {
		err = pgErrToDomain(err)
		return answer, fmt.Errorf("error collecting user's answer: %w", err)
	}

	return answer, nil
}

func (r *QuestionRepository) UpdateQuestionImage(ctx context.Context, id uuid.UUID, fileID string) error {
	const sql = `UPDATE questions SET telegram_file_id = $1 WHERE question_id = $2`

	slog.DebugContext(ctx, "updating question image", slog.String("question_uuid", id.String()), slog.String("file_id", fileID))
	_, err := r.db.Exec(ctx, sql, fileID, id)
	if err != nil {
		return fmt.Errorf("error updating question image: %w", err)
	}

	return nil
}

func (r *QuestionRepository) UpdateOptionImage(ctx context.Context, id uuid.UUID, option data.Option, fileID string) error {
	const sql = `UPDATE question_options SET telegram_file_id = $1 WHERE question_id = $2 AND hero_id = $3`

	slog.DebugContext(ctx, "updating option image", slog.String("question_uuid", id.String()), slog.String("file_id", fileID), slog.Any("option", option))
	_, err := r.db.Exec(ctx, sql, fileID, id, option.Hero.ID)
	if err != nil {
		return fmt.Errorf("error updating option image: %w", err)
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
		return nil, fmt.Errorf("error getting question stats: %w", err)
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
