package data

type Player struct {
	SteamID int64  `db:"player_steam_id"`
	Name    string `db:"name"`
	IsPro   bool   `db:"is_pro"`
	ProName string `db:"pro_name"`
}

type MatchPlayer struct {
	Player    Player   `db:"-"`
	Hero      Hero     `db:"-"`
	IsRadiant bool     `db:"is_radiant"`
	Position  Position `db:"position"`
	Items     []Item   `db:"-"`
}
