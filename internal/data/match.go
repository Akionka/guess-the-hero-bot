package data

import (
	"math"
	"time"
)

// Rank is a value object representing the rank of a player in Dota 2.
// The rank is represented as an integer, where the first digit represents the
// rank (1-8) and the second digit represents the sub-rank (1-5).
type Rank int

const (
	RankHeraldI Rank = 11 + iota
	RankHeraldII
	RankHeraldIII
	RankHeraldIV
	RankHeraldV

	RankGuardianI Rank = 21 + iota - 5
	RankGuardianII
	RankGuardianIII
	RankGuardianIV
	RankGuardianV

	RankCrusaderI Rank = 31 + iota - 10
	RankCrusaderII
	RankCrusaderIII
	RankCrusaderIV
	RankCrusaderV

	RankArchonI Rank = 41 + iota - 15
	RankArchonII
	RankArchonIII
	RankArchonIV
	RankArchonV

	RankLegendI Rank = 51 + iota - 20
	RankLegendII
	RankLegendIII
	RankLegendIV
	RankLegendV

	RankAncientI Rank = 61 + iota - 25
	RankAncientII
	RankAncientIII
	RankAncientIV
	RankAncientV

	RankDivineI Rank = 71 + iota - 30
	RankDivineII
	RankDivineIII
	RankDivineIV
	RankDivineV

	RankImmortal Rank = 80 + iota - 35
)

func (r Rank) RatingRange() [2]int {
	rank := int(r / 10)
	subRank := int(r % 10)

	subRankMMR := 154
	if rank == 7 {
		subRankMMR = 200
	}

	low := (rank-1)*154*5 + (subRank-1)*subRankMMR
	high := low + subRankMMR

	if r == RankImmortal {
		high = math.MaxInt
	}

	return [2]int{low, high}
}

type MatchID int64

// Match is an aggregate root representing a Dota 2 match.
type Match struct {
	ID          MatchID   `db:"match_id"`
	WinningTeam Team      `db:"winning_team"`
	StartedAt   time.Time `db:"started_at"`
	AvgMMR      *int      `db:"avg_mmr"`
	ActualRank  Rank      `db:"actual_rank"`

	Players []Player `db:"-"`
}

func (m Match) PlayerWon(steamID SteamID) bool {
	for _, player := range m.Players {
		if player.SteamAccountID == steamID {
			return player.Team == m.WinningTeam
		}
	}
	return false
}

func (m Match) HeroWon(heroID HeroID) bool {
	for _, player := range m.Players {
		if player.Hero.ID == heroID {
			return player.Team == m.WinningTeam
		}
	}
	return false
}

func (m Match) Teams() (radiant, dire []Player) {
	radiant = make([]Player, 0, 5)
	dire = make([]Player, 0, 5)

	for _, player := range m.Players {
		if player.Team == TeamRadiant {
			radiant = append(radiant, player)
		} else {
			dire = append(dire, player)
		}
	}
	return radiant, dire
}
