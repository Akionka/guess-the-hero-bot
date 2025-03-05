package postgres

import (
	"cmp"
	"context"
	"slices"

	"github.com/akionka/akionkabot/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HeroRepository struct {
	db *pgxpool.Pool
}

func NewHeroRepository(db *pgxpool.Pool) *HeroRepository {
	return &HeroRepository{
		db: db,
	}
}

func (r *HeroRepository) GetHeroByID(ctx context.Context, heroID int) (data.Hero, error) {
	var hero data.Hero

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return hero, err
	}

	if err = transaction(ctx, tx, "GetHeroByID", func() error {
		hero, err = r.getHeroByIDTx(ctx, tx, heroID)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return hero, err
	}

	return hero, nil
}

func (r *HeroRepository) getHeroByIDTx(ctx context.Context, tx pgx.Tx, heroID int) (data.Hero, error) {
	const sql = `SELECT h.hero_id, h.display_name, h.short_name FROM heroes h WHERE h.hero_id = $1`

	var hero data.Hero
	rows, err := tx.Query(ctx, sql, heroID)
	if err != nil {
		return hero, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[data.Hero])
}

func (r *HeroRepository) GetHeroesByIDs(ctx context.Context, heroIDs []int) ([]data.Hero, error) {
	heroes := []data.Hero{}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return heroes, err
	}

	if err = transaction(ctx, tx, "GetHeroesByIDs", func() error {
		heroes, err = r.getHeroesByIDsTx(ctx, tx, heroIDs)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return heroes, err
	}

	idIndex := make(map[int]int)
	for i, id := range heroIDs {
		idIndex[id] = i
	}

	slices.SortStableFunc(heroes, func(h1 data.Hero, h2 data.Hero) int {
		return cmp.Compare(idIndex[h1.ID], idIndex[h2.ID])
	})

	return heroes, nil
}

func (r *HeroRepository) getHeroesByIDsTx(ctx context.Context, tx pgx.Tx, heroIDs []int) ([]data.Hero, error) {
	const sql = `SELECT h.hero_id, h.display_name, h.short_name FROM heroes h WHERE h.hero_id = ANY($1)`

	heroes := []data.Hero{}
	rows, err := tx.Query(ctx, sql, heroIDs)
	if err != nil {
		return heroes, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[data.Hero])
}
