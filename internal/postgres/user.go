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

func (r *UserRepository) GetUser(ctx context.Context, id data.UserID) (*data.User, error) {
	const sql = `
		SELECT
			u.user_id, u.telegram_id, u.username, u.first_name, u.last_name, u.created_at, u.player_steam_id
		FROM users u
		WHERE u.user_id = $1
	`
	var user data.User

	r.logger.DebugContext(ctx, "getting user by id", slog.String("id", id.String()))
	if err := r.db.QueryRow(ctx, sql, id).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.SteamAccountID,
	); err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", pgErrToDomain(err))
	}

	return &user, nil
}

func (r *UserRepository) GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error) {
	const sql = `
		SELECT
			u.user_id, u.telegram_id, u.username, u.first_name, u.last_name, u.created_at, u.player_steam_id
		FROM users u
		WHERE u.telegram_id = $1
	`
	var user data.User

	r.logger.DebugContext(ctx, "getting user by telegram id", slog.Int64("id", id))
	if err := r.db.QueryRow(ctx, sql, id).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.SteamAccountID,
	); err != nil {
		return nil, fmt.Errorf("error getting user by telegram id: %w", pgErrToDomain(err))
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *data.User) (data.UserID, error) {
	const sql = `
		INSERT INTO users
		(user_id, telegram_id, username, first_name, last_name, created_at, player_steam_id) VALUES
		($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (telegram_id)
			DO UPDATE SET
				username = EXCLUDED.username,
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name,
				player_steam_id = EXCLUDED.player_steam_id
		RETURNING user_id
	`

	r.logger.DebugContext(ctx, "saving user", slog.Any("user", user))
	rows, err := r.db.Query(
		ctx, sql,
		user.ID,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.SteamAccountID,
	)
	if err != nil {
		return data.UserID(uuid.Nil), fmt.Errorf("error saving user: %w", err)
	}

	userID, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return data.UserID(uuid.Nil), fmt.Errorf("error collecting user id: %w", err)
	}

	return data.UserID(userID), nil
}

func (r *UserRepository) UpdateByID(ctx context.Context, id data.UserID, updateFn func(user *data.User) (bool, error)) error {
	return runInTx(ctx, r.db, func(tx pgx.Tx) error {
		const (
			selectSQL = `
				SELECT
					u.user_id, u.telegram_id, u.username, u.first_name, u.last_name, u.created_at, u.player_steam_id
				FROM users u
				WHERE u.user_id = $1
				FOR UPDATE
			`
			updateSQL = `
				UPDATE users SET username = $1, first_name = $2, last_name = $3, created_at = $4, player_steam_id = $5 WHERE user_id = $6
			`
		)

		var user data.User

		logger := r.logger.With(slog.String("id", id.String()))

		logger.DebugContext(ctx, "getting user by id")
		if err := r.db.QueryRow(ctx, selectSQL, id).Scan(
			&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt, &user.SteamAccountID,
		); err != nil {
			return fmt.Errorf("error getting user by id: %w", pgErrToDomain(err))
		}

		logger.DebugContext(ctx, "running updatefn")
		updated, err := updateFn(&user)
		if err != nil {
			return err
		}

		if !updated {
			return nil
		}

		// var steamAccID *int64
		// if user.SteamAccountID != nil {
		// 	steamAccID = &acc.ID
		// 	logger.DebugContext(ctx, "running upsert for steam account", slog.Int64("steam_id", acc.ID))
		// 	if _, err = tx.Exec(ctx, insertPlayerSQL, acc.ID, acc.Name, acc.IsPro, acc.ProName); err != nil {
		// 		return pgErrToDomain(err)
		// 	}
		// }

		logger.DebugContext(ctx, "running update for user")
		if _, err := tx.Exec(ctx, updateSQL, user.Username, user.FirstName, user.LastName, user.CreatedAt, user.SteamAccountID, id); err != nil {
			return pgErrToDomain(err)
		}

		return nil
	})
}
