package data

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"user_id"`
	TelegramID int64     `db:"telegram_id"`
	Username   string    `db:"username"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	CreatedAt  time.Time `db:"created_at"`
}
