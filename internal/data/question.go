package data

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type QuestionID uuid.UUID

func (id QuestionID) String() string {
	return uuid.UUID(id).String()
}

// Question is an aggregate root representing a question based on a Dota 2 match.
type Question struct {
	ID QuestionID `db:"question_id"`

	MatchID  MatchID `db:"match_id"`
	PlayerID SteamID `db:"player_steam_id"`

	Options []Option `db:"-"`

	TelegramFileID string    `db:"telegram_file_id"`
	CreatedAt      time.Time `db:"created_at"`
}

func (q *Question) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", q.ID.String()),
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

type UserAnswerID uuid.UUID

func (id UserAnswerID) String() string {
	return uuid.UUID(id).String()
}

type UserAnswer struct {
	Option     `db:"-"`
	ID         UserAnswerID `db:"user_answer_id"`
	AnsweredAt time.Time    `db:"answered_at"`
}

func (a *UserAnswer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", a.ID.String()),
		slog.Any("option", a.Option),
	)
}
