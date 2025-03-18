package postgres

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/akionka/akionkabot/internal/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HeroRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewHeroRepository(db *pgxpool.Pool, logger *slog.Logger) *HeroRepository {
	return &HeroRepository{
		db:     db,
		logger: logger,
	}
}

func (r *HeroRepository) GetHero(ctx context.Context, id data.HeroID) (data.Hero, error) {
	const sql = `SELECT h.hero_id, h.display_name, h.short_name FROM heroes h WHERE h.hero_id = $1`
	var hero data.Hero

	r.logger.DebugContext(ctx, "getting hero by id", slog.Int("id", int(id)))
	rows, err := r.db.Query(ctx, sql, id)
	if err != nil {
		return hero, fmt.Errorf("error getting hero by id: %w", err)
	}

	hero, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[data.Hero])
	if err != nil {
		err = pgErrToDomain(err)
		return hero, fmt.Errorf("error collecting hero: %w", err)
	}

	return hero, nil
}

func (r *HeroRepository) GetHeroes(ctx context.Context, ids []data.HeroID) ([]data.Hero, error) {
	const sql = `SELECT h.hero_id, h.display_name, h.short_name FROM heroes h WHERE h.hero_id = ANY($1)`

	r.logger.DebugContext(ctx, "getting heroes by ids", slog.Any("ids", ids))
	rows, err := r.db.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("error getting heroes by ids: %w", err)
	}

	heroes, err := pgx.CollectRows(rows, pgx.RowToStructByName[data.Hero])
	if err != nil {
		return nil, fmt.Errorf("error collecting heroes: %w", err)
	}

	if len(heroes) != len(ids) {
		return nil, data.ErrNotFound
	}

	idIndex := make(map[data.HeroID]int)
	for i, id := range ids {
		idIndex[id] = i
	}

	slices.SortStableFunc(heroes, func(h1 data.Hero, h2 data.Hero) int {
		return cmp.Compare(idIndex[h1.ID], idIndex[h2.ID])
	})

	return heroes, nil
}
