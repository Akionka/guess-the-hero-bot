package data

import "log/slog"

type ItemID int32

// Item is an entity representing an item in the Dota 2.
type Item struct {
	ID          ItemID `db:"item_id"`
	DisplayName string `db:"display_name"`
	ShortName   string `db:"short_name"`
	Image       Image  `db:"-"`
}

func (i Item) String() string {
	return i.DisplayName
}

func (i Item) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("short_name", i.ShortName),
	)
}
