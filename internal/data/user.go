package data

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `db:"user_id"`
	TelegramID int64     `db:"telegram_id"`
	Username   string    `db:"username"`
	FirstName  string    `db:"first_name"`
	LastName   string    `db:"last_name"`
	CreatedAt  time.Time `db:"created_at"`
	SteamAcc    *Player   `db:"-"`
}

func (u *User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("uuid", u.ID.String()),
		slog.Int64("telegram_id", u.TelegramID),
		slog.String("telegram_username", u.Username),
		slog.Any("steam", u.SteamAcc),
	)
}

var _ encoding.BinaryMarshaler = (*User)(nil)
var _ encoding.BinaryUnmarshaler = (*User)(nil)

func (u *User) MarshalBinary() (data []byte, err error) {
	order := binary.LittleEndian
	buf := bytes.NewBuffer(make([]byte, 0, 16+8+8*3+len(u.Username)+len(u.FirstName)+len(u.LastName)+1+16))

	binary.Write(buf, order, u.ID)
	binary.Write(buf, order, u.TelegramID)

	writeBinaryString(buf, order, u.Username)
	writeBinaryString(buf, order, u.FirstName)
	writeBinaryString(buf, order, u.LastName)

	b, _ := u.CreatedAt.MarshalBinary()
	binary.Write(buf, order, uint8(len(b)))
	binary.Write(buf, order, b)

	if u.SteamAcc == nil {
		binary.Write(buf, order, uint32(0))
	} else {
		b, _ = u.SteamAcc.MarshalBinary()
		binary.Write(buf, order, uint32(len(b)))
		binary.Write(buf, order, b)
	}

	return buf.Bytes(), nil
}

func (u *User) UnmarshalBinary(data []byte) error {
	order := binary.LittleEndian

	r := bytes.NewReader(data)

	binary.Read(r, order, &u.ID)
	binary.Read(r, order, &u.TelegramID)
	readBinaryString(r, order, &u.Username)
	readBinaryString(r, order, &u.FirstName)
	readBinaryString(r, order, &u.LastName)

	var timeLen uint8
	binary.Read(r, order, &timeLen)
	timeBytes := make([]byte, timeLen)
	r.Read(timeBytes)
	u.CreatedAt.UnmarshalBinary(timeBytes)

	var steamAccLen uint32
	binary.Read(r, order, &steamAccLen)
	if steamAccLen != 0 {
		steamAccBytes := make([]byte, steamAccLen)
		r.Read(steamAccBytes)
		u.SteamAcc.UnmarshalBinary(steamAccBytes)
	}

	return nil
}
