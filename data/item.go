package data

type Item struct {
	ID          int    `db:"item_id"`
	DisplayName string `db:"display_name"`
	ShortName   string `db:"short_name"`
	Image       Image  `db:"-"` // Might be nil if item somehow does not have an image
}
