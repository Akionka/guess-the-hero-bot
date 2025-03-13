package stratz

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/akionka/akionkabot/internal/data"
)

type Position string

const (
	PositionCarry       Position = "Carry"
	PositionMid         Position = "Mid"
	PositionOfflane     Position = "Offlane"
	PositionSoftSupport Position = "Soft Support"
	PositionHardSupport Position = "Hard Support"
	PositionUnknown     Position = "Unknown"
)

func (p *Position) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToUpper(s) {
	case "POSITION_1":
		*p = PositionCarry
	case "POSITION_2":
		*p = PositionMid
	case "POSITION_3":
		*p = PositionOfflane
	case "POSITION_4":
		*p = PositionSoftSupport
	case "POSITION_5":
		*p = PositionHardSupport
	default:
		*p = PositionUnknown
	}

	return nil
}

func (p Position) MarshalJSON() ([]byte, error) {
	var s string
	switch p {
	case PositionCarry:
		s = "POSITION_1"
	case PositionMid:
		s = "POSITION_2"
	case PositionOfflane:
		s = "POSITION_3"
	case PositionSoftSupport:
		s = "POSITION_4"
	case PositionHardSupport:
		s = "POSITION_5"
	default:
		s = "Unknown"
	}

	return json.Marshal(s)
}

func (p Position) toDomain() data.Position {
	switch p {
	case PositionCarry:
		return data.PositionCarry
	case PositionMid:
		return data.PositionMid
	case PositionOfflane:
		return data.PositionOfflane
	case PositionSoftSupport:
		return data.PositionSoftSupport
	case PositionHardSupport:
		return data.PositionHardSupport
	default:
		return data.PositionUnknown
	}
}

type Match struct {
	ID            int64 `json:"id"`
	DidRadiantWin bool  `json:"didRadiantWin"`
	ActualRank    int   `json:"actualRank"`
	StartDateTime int64 `json:"startDateTime"`
	Players       []MatchPlayer
}

func (m Match) toDomain() data.Match {
	players := make([]data.MatchPlayer, len(m.Players))
	for i, p := range m.Players {
		players[i] = p.toDomain()
	}

	return data.Match{
		ID:         m.ID,
		RadiantWon: m.DidRadiantWin,
		StartedAt:  time.Unix(m.StartDateTime, 0),
		AvgMMR:     nil,
		ActualRank: m.ActualRank,
		Players:    players,
	}
}

type MatchPlayer struct {
	Hero         Hero     `json:"hero"`
	Item0Id      int      `json:"item0Id"`
	Item1Id      int      `json:"item1Id"`
	Item2Id      int      `json:"item2Id"`
	Item3Id      int      `json:"item3Id"`
	Item4Id      int      `json:"item4Id"`
	Item5Id      int      `json:"item5Id"`
	IsRadiant    bool     `json:"isRadiant"`
	Position     Position `json:"position"`
	SteamAccount Player   `json:"steamAccount"`
}

func (mp MatchPlayer) toDomain() data.MatchPlayer {
	items := []data.Item{{ID: mp.Item0Id}, {ID: mp.Item1Id}, {ID: mp.Item2Id}, {ID: mp.Item3Id}, {ID: mp.Item4Id}, {ID: mp.Item5Id}}
	return data.MatchPlayer{
		Player:    mp.SteamAccount.toDomain(),
		Hero:      mp.Hero.toDomain(),
		IsRadiant: mp.IsRadiant,
		Position:  mp.Position.toDomain(),
		Items:     items,
	}
}

type Hero struct {
	ID          int    `json:"id"`
	DisplayName string `json:"displayName"`
	ShortName   string `json:"shortName"`
}

func (h Hero) toDomain() data.Hero {
	return data.Hero{
		ID:          h.ID,
		DisplayName: h.DisplayName,
		ShortName:   h.ShortName,
	}
}

type Player struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	ProSteamAccount *ProPlayer `json:"proSteamAccount"`
}

func (p Player) toDomain() data.Player {
	dp := data.Player{
		SteamID: p.ID,
		Name:    p.Name,
		IsPro:   p.ProSteamAccount != nil,
	}
	if dp.IsPro {
		dp.ProName = p.ProSteamAccount.Name
	}
	return dp
}

type ProPlayer struct {
	Name string `json:"name"`
}
