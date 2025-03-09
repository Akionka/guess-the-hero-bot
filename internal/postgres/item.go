package postgres

import (
	"cmp"
	"context"
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
	var item data.Item

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return item, err
	}

	if err = transaction(ctx, tx, "GetItemByID", func() error {
		item, err = r.getItemByIDTx(ctx, tx, id)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return item, err
	}

	return item, nil
}

func (r *ItemRepository) getItemByIDTx(ctx context.Context, tx pgx.Tx, id int) (data.Item, error) {
	const sql = `SELECT i.item_id, i.display_name, i.short_name FROM items i WHERE i.item_id = $1`
	r.logger.DebugContext(ctx, "getting item by id", slog.Int("id", id))

	var item data.Item
	rows, err := tx.Query(ctx, sql, id)
	if err != nil {
		return item, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[data.Item])
}

func (r *ItemRepository) GetItemsByIDs(ctx context.Context, ids []int) ([]data.Item, error) {
	var items []data.Item

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return items, err
	}

	if err = transaction(ctx, tx, "GetItemsByIDs", func() error {
		items, err = r.getItemsByIDsTx(ctx, tx, ids)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return items, err
	}

	idIndex := make(map[int]int)
	for i, id := range ids {
		idIndex[id] = i
	}

	slices.SortStableFunc(items, func(h1 data.Item, h2 data.Item) int {
		return cmp.Compare(idIndex[h1.ID], idIndex[h2.ID])
	})

	return items, nil
}

func (r *ItemRepository) getItemsByIDsTx(ctx context.Context, tx pgx.Tx, ids []int) ([]data.Item, error) {
	const sql = `SELECT i.item_id, i.display_name, i.short_name FROM items i WHERE i.item_id = ANY($1)`
	r.logger.DebugContext(ctx, "getting items by ids", slog.Any("ids", ids))

	var items []data.Item
	rows, err := tx.Query(ctx, sql, ids)
	if err != nil {
		return items, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[data.Item])
}
