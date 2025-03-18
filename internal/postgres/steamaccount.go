package postgres

import (
	"context"
	"log/slog"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SteamAccountRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewSteamAccountRepository(db *pgxpool.Pool, logger *slog.Logger) *SteamAccountRepository {
	return &SteamAccountRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SteamAccountRepository) SaveAccount(ctx context.Context, account *data.SteamAccount) error {
	const sql = `
		INSERT INTO player_accounts (player_steam_id, name, is_pro, pro_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (player_steam_id)
			DO UPDATE SET
				name = EXCLUDED.name,
				is_pro = EXCLUDED.is_pro,
				pro_name = EXCLUDED.pro_name
	`

	r.logger.DebugContext(ctx, "saving steam account", slog.Int64("steam_id", int64(account.ID)))
	if _, err := r.db.Exec(ctx, sql, account.ID, account.Name, account.IsPro, account.ProName); err != nil {
		return pgErrToDomain(err)
	}

	return nil
}
