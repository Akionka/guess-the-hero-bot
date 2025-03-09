package data

import "log/slog"

type Item struct {
	ID          int    `db:"item_id"`
	DisplayName string `db:"display_name"`
	ShortName   string `db:"short_name"`
	Image       Image  `db:"-"` // Might be nil if item somehow does not have an image
}

func (i Item) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("short_name", i.ShortName),
	)
}
