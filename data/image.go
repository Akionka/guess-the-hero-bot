package data

import (
	"bytes"
	"image"

	"image/png"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/jackc/pgx/v5/pgtype"
	_ "golang.org/x/image/webp"
)

type Image struct {
	image.Image
}

var (
	_ pgtype.BytesScanner = (*Image)(nil)
	_ pgtype.BytesValuer  = (*Image)(nil)
)

func (i *Image) ScanBytes(v []byte) error {
	img, _, err := image.Decode(bytes.NewBuffer(v))
	if err != nil {
		return err
	}
	i.Image = img
	return nil
}

func (i *Image) BytesValue() ([]byte, error) {
	var buf bytes.Buffer

	if err := png.Encode(&buf, i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
