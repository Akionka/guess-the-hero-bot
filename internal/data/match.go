package data

import "time"

type Match struct {
	ID         int64     `db:"match_db"`
	RadiantWon bool      `db:"radiant_won"`
	StartedAt  time.Time `db:"started_at"`
	AvgMMR     *int      `db:"avg_mmr"`

	// Encodes rank as  adecimal number.
	// 11 = Herald I
	// 23 = Guadiang III etc
	// 80 = Immortal
	ActualRank int `db:"actual_rank"`

	Players []MatchPlayer `db:"-"`
}

func (m Match) PlayerWon(steamID int64) bool {
	for _, player := range m.Players {
		if player.IsRadiant == m.RadiantWon {
			return true
		}
	}
	return false
}

func (m Match) HeroWon(heroID int) bool {
	for _, player := range m.Players {
		if player.Hero.ID == heroID && player.IsRadiant == m.RadiantWon {
			return true
		}
	}
	return false
}

func (m Match) Teams() (radiant, dire []MatchPlayer) {
	radiant = make([]MatchPlayer, 0, 5)
	dire = make([]MatchPlayer, 0, 5)

	for _, player := range m.Players {
		if player.IsRadiant {
			radiant = append(radiant, player)
		} else {
			dire = append(dire, player)
		}
	}
	return radiant, dire
}
