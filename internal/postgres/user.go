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
	const sql = `
		SELECT
			u.user_id, u.telegram_id, u.username, u.first_name, u.last_name, u.created_at,
			pa.player_steam_id, pa.name, pa.is_pro, pa.pro_name
		FROM users u
		LEFT JOIN player_accounts pa ON pa.player_steam_id = u.player_steam_id
		WHERE u.user_id = $1
	`
	var (
		user    data.User
		steamID *int64
		name    *string
		isPro   *bool
		proName *string
	)

	r.logger.DebugContext(ctx, "getting user by id", slog.String("uuid", id.String()))
	if err := r.db.QueryRow(ctx, sql, id).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt,
		&steamID, &name, &isPro, &proName,
	); err != nil {
		return nil, fmt.Errorf("error getting user by id: %w", pgErrToDomain(err))
	}

	if steamID != nil {
		user.SteamAcc = &data.SteamAccount{
			SteamID: *steamID,
			Name:    *name,
			IsPro:   *isPro,
			ProName: *proName,
		}
	}

	return &user, nil
}

func (r *UserRepository) GetUserByTelegramID(ctx context.Context, id int64) (*data.User, error) {
	const sql = `
		SELECT
			u.user_id, u.telegram_id, u.username, u.first_name, u.last_name, u.created_at,
			pa.player_steam_id, pa.name, pa.is_pro, pa.pro_name
		FROM users u
		LEFT JOIN player_accounts pa ON pa.player_steam_id = u.player_steam_id
		WHERE u.telegram_id = $1
	`
	var (
		user    data.User
		steamID *int64
		name    *string
		isPro   *bool
		proName *string
	)

	r.logger.DebugContext(ctx, "getting user by telegram id", slog.Int64("id", id))
	if err := r.db.QueryRow(ctx, sql, id).Scan(
		&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt,
		&steamID, &name, &isPro, &proName,
	); err != nil {
		return nil, fmt.Errorf("error getting user by telegram id: %w", pgErrToDomain(err))
	}

	if steamID != nil {
		user.SteamAcc = &data.SteamAccount{
			SteamID: *steamID,
			Name:    *name,
			IsPro:   *isPro,
			ProName: *proName,
		}
	}

	return &user, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, user *data.User) (uuid.UUID, error) {
	const sql = `
		INSERT INTO users
		(user_id, telegram_id, username, first_name, last_name, created_at) VALUES
		($1, $2, $3, $4, $5, $6)
		ON CONFLICT (telegram_id)
		DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name
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

func (r *UserRepository) UpdateByID(ctx context.Context, id uuid.UUID, updateFn func(user *data.User) (bool, error)) error {
	return runInTx(ctx, r.db, func(tx pgx.Tx) error {
		const (
			selectSQL = `
				SELECT
					u.user_id, u.telegram_id, u.username, u.first_name, u.last_name, u.created_at,
					pa.player_steam_id, pa.name, pa.is_pro, pa.pro_name
				FROM users u
				LEFT JOIN player_accounts pa ON pa.player_steam_id = u.player_steam_id
				WHERE u.user_id = $1
				FOR UPDATE OF u
			`
			updateSQL = `
				UPDATE users SET username = $1, first_name = $2, last_name = $3, created_at = $4, player_steam_id = $5 WHERE user_id = $6
			`
			insertPlayerSQL = `
				INSERT INTO player_accounts (player_steam_id, name, is_pro, pro_name) VALUES
				($1, $2, $3, $4)
				ON CONFLICT (player_steam_id)
				DO UPDATE SET
					name = EXCLUDED.name,
					is_pro = EXCLUDED.is_pro,
					pro_name = EXCLUDED.pro_name
			`
		)

		var (
			user    data.User
			steamID *int64
			name    *string
			isPro   *bool
			proName *string
		)

		logger := r.logger.With(slog.String("uuid", id.String()))

		logger.DebugContext(ctx, "getting user by id")
		if err := r.db.QueryRow(ctx, selectSQL, id).Scan(
			&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt,
			&steamID, &name, &isPro, &proName,
		); err != nil {
			return fmt.Errorf("error getting user by id: %w", pgErrToDomain(err))
		}

		if steamID != nil {
			user.SteamAcc = &data.SteamAccount{
				SteamID: *steamID,
				Name:    *name,
				IsPro:   *isPro,
				ProName: *proName,
			}
		}

		logger.DebugContext(ctx, "running updatefn")
		updated, err := updateFn(&user)
		if err != nil {
			return err
		}

		if !updated {
			return nil
		}

		var steamAccID *int64
		if user.SteamAcc != nil {
			acc := user.SteamAcc
			steamAccID = &acc.SteamID
			logger.DebugContext(ctx, "running upsert for steam account", slog.Int64("steam_id", acc.SteamID))
			if _, err = tx.Exec(ctx, insertPlayerSQL, acc.SteamID, acc.Name, acc.IsPro, acc.ProName); err != nil {
				return pgErrToDomain(err)
			}
		}

		logger.DebugContext(ctx, "running update for user")
		if _, err := tx.Exec(ctx, updateSQL, user.Username, user.FirstName, user.LastName, user.CreatedAt, steamAccID, id); err != nil {
			return pgErrToDomain(err)
		}

		return nil
	})
}
