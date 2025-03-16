package data

import (
	"bytes"
	"encoding/binary"
	"log/slog"
)

type Player struct {
	SteamID int64  `db:"player_steam_id"`
	Name    string `db:"name"`
	IsPro   bool   `db:"is_pro"`
	ProName string `db:"pro_name"`
}

func (p *Player) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("steam_id", p.SteamID),
		slog.String("steam_name", p.Name),
	)
}

func (p *Player) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 8+4+len(p.Name)+1+4+len(p.ProName)))

	binary.Write(buf, binary.LittleEndian, p.SteamID)
	writeBinaryString(buf, binary.LittleEndian, p.Name)
	binary.Write(buf, binary.LittleEndian, p.IsPro)
	writeBinaryString(buf, binary.LittleEndian, p.ProName)

	return buf.Bytes(), nil
}

func (p *Player) UnmarshalBinary(b []byte) error {
	r := bytes.NewReader(b)

	binary.Read(r, binary.LittleEndian, p.SteamID)
	readBinaryString(r, binary.LittleEndian, &p.Name)
	binary.Read(r, binary.LittleEndian, p.IsPro)
	readBinaryString(r, binary.LittleEndian, &p.ProName)

	return nil
}

type MatchPlayer struct {
	Player    Player   `db:"-"`
	Hero      Hero     `db:"-"`
	IsRadiant bool     `db:"is_radiant"`
	Position  Position `db:"position"`
	Items     []Item   `db:"-"`
}
