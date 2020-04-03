package util

import (
	"encoding/binary"
	"io"
)

var DefaultByteOrder = binary.BigEndian

func Write(w io.Writer, data interface{}) error {
	if s, ok := data.(string); ok {
		bs := StringToBytes(s)
		if err := binary.Write(w, DefaultByteOrder, int64(len(bs))); err != nil {
			return err
		}
		if err := binary.Write(w, DefaultByteOrder, bs); err != nil {
			return err
		}
		return nil
	}
	return binary.Write(w, DefaultByteOrder, data)
}

func Read(r io.Reader, data interface{}) error {
	if s, ok := data.(*string); ok {
		var len int64
		if err := binary.Read(r, DefaultByteOrder, &len); err != nil {
			return err
		}

		bs := make([]byte, len)
		if _, err := io.ReadFull(r, bs); err != nil {
			return err
		}

		*s = BytesToString(bs)
		return nil
	}
	return binary.Read(r, DefaultByteOrder, data)
}
