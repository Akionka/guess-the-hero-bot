package data

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID uuid.UUID `db:"question_id"`

	Match  Match        `db:"-"`
	Player *MatchPlayer `db:"-"`

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
	Hero           Hero   `db:"-"`
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
	Option     `db:"-"`
	ID         uuid.UUID `db:"user_answer_id"`
	AnsweredAt time.Time `db:"answered_at"`
}

func (a *UserAnswer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("uuid", a.ID.String()),
		slog.Any("option", a.Option),
	)
}
