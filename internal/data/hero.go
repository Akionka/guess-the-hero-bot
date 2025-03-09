package data

import "log/slog"

type Hero struct {
	ID          int    `db:"hero_id"`
	DisplayName string `db:"display_name"`
	ShortName   string `db:"short_name"`
	Image       Image  `db:"-"`
}

func (h Hero) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("short_name", h.ShortName),
	)
}
