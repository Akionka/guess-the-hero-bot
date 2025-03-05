package postgres

import (
	"cmp"
	"context"
	"slices"

	"github.com/akionka/akionkabot/data"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository struct {
	db *pgxpool.Pool
}

func NewItemRepository(db *pgxpool.Pool) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) GetItemByID(ctx context.Context, itemID int) (data.Item, error) {
	var item data.Item

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return item, err
	}

	if err = transaction(ctx, tx, "GetItemByID", func() error {
		item, err = r.getItemByIDTx(ctx, tx, itemID)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return item, err
	}

	return item, nil
}

func (r *ItemRepository) getItemByIDTx(ctx context.Context, tx pgx.Tx, itemID int) (data.Item, error) {
	const sql = `SELECT i.item_id, i.display_name, i.short_name FROM items i WHERE i.item_id = $1`

	var item data.Item
	rows, err := tx.Query(ctx, sql, itemID)
	if err != nil {
		return item, err
	}

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[data.Item])
}

func (r *ItemRepository) GetItemsByIDs(ctx context.Context, itemIDs []int) ([]data.Item, error) {
	items := []data.Item{}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return items, err
	}

	if err = transaction(ctx, tx, "GetItemsByIDs", func() error {
		items, err = r.getItemsByIDsTx(ctx, tx, itemIDs)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return items, err
	}

	idIndex := make(map[int]int)
	for i, id := range itemIDs {
		idIndex[id] = i
	}

	slices.SortStableFunc(items, func(h1 data.Item, h2 data.Item) int {
		return cmp.Compare(idIndex[h1.ID], idIndex[h2.ID])
	})

	return items, nil
}

func (r *ItemRepository) getItemsByIDsTx(ctx context.Context, tx pgx.Tx, itemIDs []int) ([]data.Item, error) {
	const sql = `SELECT i.item_id, i.display_name, i.short_name FROM items i WHERE i.item_id = ANY($1)`

	items := []data.Item{}
	rows, err := tx.Query(ctx, sql, itemIDs)
	if err != nil {
		return items, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByName[data.Item])
}
