package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func transaction(ctx context.Context, tx pgx.Tx, label string, f func() error) error {
	if err := f(); err != nil {
		_ = tx.Rollback(ctx)

		return fmt.Errorf("failed to perform transaction %s: %w", label, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction %s: %w", label, err)
	}

	return nil
}
