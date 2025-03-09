package data

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID uuid.UUID `db:"question_id"`

	MatchID        int64     `db:"match_id"`
	MatchStartedAt time.Time `db:"match_started_at"`

	PlayerID    int64    `db:"player_id"`
	PlayerName  string   `db:"player_name"`
	PlayerIsPro bool     `db:"player_is_pro"`
	PlayerPos   Position `db:"player_pos"`
	PlayerMMR   int      `db:"player_mmr"`
	IsWon       bool     `db:"is_won"`

	Items   []Item   `db:"-"`
	Options []Option `db:"-"`

	TelegramFileID string    `db:"telegram_file_id"`
	CreatedAt      time.Time `db:"created_at"`
}

func (q *Question) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("uuid", q.ID.String()),
	)
}

type Option struct {
	Hero           `db:"-"`
	IsCorrect      bool   `db:"is_correct"`
	TelegramFileID string `db:"telegram_file_id"`
}

func (o *Option) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Bool("is_correct", o.IsCorrect),
		slog.Any("hero", o.Hero),
	)
}

type UserAnswer struct {
	Option
	ID         uuid.UUID `db:"user_answer_id"`
	AnsweredAt time.Time `db:"answered_at"`
}

func (a *UserAnswer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("uuid", a.ID.String()),
		slog.Any("option", a.Option),
	)
}
