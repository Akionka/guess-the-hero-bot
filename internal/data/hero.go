package data

import "log/slog"

type HeroID int32

// Hero is an entity representing a hero in the Dota 2.
type Hero struct {
	ID          HeroID `db:"hero_id"`
	DisplayName string `db:"display_name"`
	ShortName   string `db:"short_name"`
	Image       Image  `db:"-"`
}

func (h Hero) String() string {
	return h.DisplayName
}

func (h Hero) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("short_name", h.ShortName),
	)
}
