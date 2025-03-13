package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool, logger *slog.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*data.User, error) {
	const sql = `SELECT user_id, telegram_id, username, first_name, last_name, created_at FROM users u WHERE u.user_id = $1`

	r.logger.DebugContext(ctx, "getting user by id", slog.String("uuid", id.String()))
	rows, err := r.db.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.User])
	if err != nil {
		err = pgErrToDomain(err)
		return nil, fmt.Errorf("error collecting user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error) {
	const sql = `SELECT user_id, telegram_id, username, first_name, last_name, created_at FROM users u WHERE u.telegram_id = $1`

	r.logger.DebugContext(ctx, "getting user by telegram id", slog.Int64("telegram_id", id))
	rows, err := r.db.Query(ctx, sql, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user by telegram id: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[data.User])
	if err != nil {
		err = pgErrToDomain(err)
		return nil, fmt.Errorf("error collecting user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *data.User) (uuid.UUID, error) {
	const sql = `
		INSERT INTO users
		(user_id, telegram_id, username, first_name, last_name, created_at) VALUES
		($1, $2, $3, $4, $5, $6)
		ON CONFLICT (telegram_id)
		DO UPDATE SET
		username = EXCLUDED.username,
		first_name = EXCLUDED.first_name,
		last_name = EXCLUDED.last_name
		RETURNING user_id`

	r.logger.DebugContext(ctx, "saving user", slog.Any("user", user))
	rows, err := r.db.Query(
		ctx, sql,
		user.ID,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error saving user: %w", err)
	}

	userID, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return uuid.Nil, fmt.Errorf("error collecting user id: %w", err)
	}

	return userID, nil
}
