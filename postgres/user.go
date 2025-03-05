package postgres

import (
	"context"

	"github.com/akionka/akionkabot/data"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*data.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	var user *data.User

	err = transaction(ctx, tx, "GetUser", func() error {
		user, err = r.getUserTx(ctx, tx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}

func (r *UserRepository) getUserTx(ctx context.Context, tx pgx.Tx, userID uuid.UUID) (*data.User, error) {
	const sql = `SELECT user_id, telegram_id, username, first_name, last_name, created_at FROM users u WHERE u.user_id = $1`

	rows, err := tx.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[data.User])
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByTelegramID(ctx context.Context, userID int64) (*data.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	var user *data.User

	err = transaction(ctx, tx, "GetUserByTelegramID", func() error {
		user, err = r.getUserByTelegramIDTx(ctx, tx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	return user, err
}

func (r *UserRepository) getUserByTelegramIDTx(ctx context.Context, tx pgx.Tx, userID int64) (*data.User, error) {
	const sql = `SELECT user_id, telegram_id, username, first_name, last_name, created_at FROM users u WHERE u.telegram_id = $1`

	rows, err := tx.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[data.User])
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) SaveUser(ctx context.Context, user *data.User) (*data.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	var savedUser *data.User

	err = transaction(ctx, tx, "SaveUser", func() error {
		id, err := r.saveUserTx(ctx, tx, user)
		if err != nil {
			return err
		}

		savedUser, err = r.getUserTx(ctx, tx, id)
		return err
	})

	return savedUser, err
}

func (r *UserRepository) saveUserTx(ctx context.Context, tx pgx.Tx, user *data.User) (uuid.UUID, error) {
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

	rows, err := tx.Query(
		ctx,
		sql,
		user.ID,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
}
