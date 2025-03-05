package data

type Hero struct {
	ID          int    `db:"hero_id"`
	DisplayName string `db:"display_name"`
	ShortName   string `db:"short_name"`
	Image       Image  `db:"-"`
}
