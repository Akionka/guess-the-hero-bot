package data

import (
	"bytes"
	"encoding/binary"
	"log/slog"
)

type SteamAccount struct {
	SteamID int64  `db:"player_steam_id"`
	Name    string `db:"name"`
	IsPro   bool   `db:"is_pro"`
	ProName string `db:"pro_name"`
}

func (a *SteamAccount) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("steam_id", a.SteamID),
		slog.String("steam_name", a.Name),
	)
}

func (a *SteamAccount) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 8+4+len(a.Name)+1+4+len(a.ProName)))

	binary.Write(buf, binary.LittleEndian, a.SteamID)
	writeBinaryString(buf, binary.LittleEndian, a.Name)
	binary.Write(buf, binary.LittleEndian, a.IsPro)
	writeBinaryString(buf, binary.LittleEndian, a.ProName)

	return buf.Bytes(), nil
}

func (a *SteamAccount) UnmarshalBinary(b []byte) error {
	r := bytes.NewReader(b)

	binary.Read(r, binary.LittleEndian, a.SteamID)
	readBinaryString(r, binary.LittleEndian, &a.Name)
	binary.Read(r, binary.LittleEndian, a.IsPro)
	readBinaryString(r, binary.LittleEndian, &a.ProName)

	return nil
}

type MatchPlayer struct {
	SteamAccount SteamAccount `db:"-"`
	Hero         Hero         `db:"-"`
	IsRadiant    bool         `db:"is_radiant"`
	Position     Position     `db:"position"`
	Items        []Item       `db:"-"`
}
