package postgres

import (
	"errors"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// pgErrToDomain converts a pgx error to a domain error. It returns a domain error if the pgx error is a known error, otherwise it returns the original error.
func pgErrToDomain(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return data.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return data.ErrAlreadyExists
		}
	}

	return err
}
