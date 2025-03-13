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

type ItemRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func NewItemRepository(db *pgxpool.Pool, logger *slog.Logger) *ItemRepository {
	return &ItemRepository{
		db:     db,
		logger: logger,
	}
}

func (r *ItemRepository) GetItemByID(ctx context.Context, id int) (data.Item, error) {
	const sql = `SELECT i.item_id, i.display_name, i.short_name FROM items i WHERE i.item_id = $1`
	var item data.Item

	r.logger.DebugContext(ctx, "getting item by id", slog.Int("id", id))
	rows, err := r.db.Query(ctx, sql, id)
	if err != nil {
		return item, fmt.Errorf("error getting item by id: %w", err)
	}

	item, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[data.Item])
	if err != nil {
		err = pgErrToDomain(err)
		return item, fmt.Errorf("error collecting item: %w", err)
	}

	return item, nil
}

func (r *ItemRepository) GetItemsByIDs(ctx context.Context, ids []int) ([]data.Item, error) {
	const sql = `SELECT i.item_id, i.display_name, i.short_name FROM items i WHERE i.item_id = ANY($1)`

	r.logger.DebugContext(ctx, "getting items by ids", slog.Any("ids", ids))
	rows, err := r.db.Query(ctx, sql, ids)
	if err != nil {
		return nil, fmt.Errorf("error getting items by ids: %w", err)
	}

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[data.Item])
	if err != nil {
		return nil, fmt.Errorf("error collecting items: %w", err)
	}

	if len(items) != len(ids) {
		return nil, data.ErrNotFound
	}

	idIndex := make(map[int]int)
	for i, id := range ids {
		idIndex[id] = i
	}

	slices.SortStableFunc(items, func(i1 data.Item, i2 data.Item) int {
		return cmp.Compare(idIndex[i1.ID], idIndex[i2.ID])
	})

	return items, nil
}
