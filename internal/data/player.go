package data

// Team is a value object representing the two teams in a Dota 2 match.
type Team string

const (
	TeamRadiant Team = "Radiant"
	TeamDire    Team = "Dire"
)

// Position is a value object representing the position of a player in a Dota 2 match.
type Position string

const (
	PositionCarry       Position = "Carry"
	PositionMid         Position = "Mid"
	PositionOfflane     Position = "Offlane"
	PositionSoftSupport Position = "Soft Support"
	PositionHardSupport Position = "Hard Support"
	PositionUnknown     Position = "Unknown"
)

func (p Position) ToEmoji() string {
	switch p {
	case PositionCarry:
		return "🗡"
	case PositionMid:
		return "🏹"
	case PositionOfflane:
		return "🛡"
	case PositionSoftSupport:
		return "🔮"
	case PositionHardSupport:
		return "✨"
	default:
		return "❌"
	}
}
func (p Position) String() string {
	switch p {
	case PositionCarry:
		return "Carry"
	case PositionMid:
		return "Mid"
	case PositionOfflane:
		return "Offlane"
	case PositionSoftSupport:
		return "Soft Support"
	case PositionHardSupport:
		return "Hard Support"
	default:
		return "Unknown"
	}
}

// Player is an entity representing a player in a Dota 2 match.
type Player struct {
	SteamAccountID SteamID  `db:"player_steam_id"`
	Hero           Hero     `db:"-"`
	Team           Team     `db:"team"`
	Position       Position `db:"position"`
	Items          []Item   `db:"-"`
}
