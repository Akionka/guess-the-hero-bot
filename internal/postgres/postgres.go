package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func transaction(ctx context.Context, tx pgx.Tx, label string, f func() error) error {
	if err := f(); err != nil {
		_ = tx.Rollback(ctx)

		if errors.Is(err, pgx.ErrNoRows) {
			return data.ErrNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return data.ErrAlreadyExists
		}

		return fmt.Errorf("failed to perform transaction %s: %w", label, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction %s: %w", label, err)
	}

	return nil
}
