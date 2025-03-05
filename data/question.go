package data

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID      uuid.UUID `db:"question_id"`
	MatchID int64     `db:"match_id"`

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

type Option struct {
	Hero           `db:"-"`
	IsCorrect      bool   `db:"is_correct"`
	TelegramFileID string `db:"telegram_file_id"`
}

type UserOption struct {
	ID uuid.UUID `db:"user_question_id"`
	Option
	AnsweredAt time.Time `db:"answered_at"`
}
