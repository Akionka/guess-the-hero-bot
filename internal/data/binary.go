package data

import (
	"encoding/binary"
	"io"
	"unsafe"
)

func writeBinaryString(w io.Writer, order binary.ByteOrder, src string) {
	binary.Write(w, order, int64(len(src)))
	binary.Write(w, order, unsafe.Slice(unsafe.StringData(src), len(src)))
}

func readBinaryString(r io.Reader, order binary.ByteOrder, dst *string) error {
	var length int64

	if err := binary.Read(r, order, &length); err != nil {
		return err
	}

	if length == 0 {
		return nil
	}

	buf := make([]byte, length)

	err := binary.Read(r, order, buf)

	*dst = unsafe.String(&buf[0], length)
	return err
}
