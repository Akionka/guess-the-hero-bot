package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func runInTx(ctx context.Context, db *pgxpool.Pool, fn func(pgx.Tx) error) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rbErr := tx.Rollback(ctx)
			if rbErr != nil {
				err = errors.Join(rbErr, err)
			}
			return
		}
		err = tx.Commit(ctx)
	}()

	return fn(tx)
}
