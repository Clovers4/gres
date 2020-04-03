package proto

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/clovers4/gres/engine/object"
	"io"
	"strconv"

	"github.com/clovers4/gres/util"
)

type Writer struct {
	wr *bufio.Writer

	lenBuf []byte
	numBuf []byte
}

func NewWriter(wr io.Writer) *Writer {
	return &Writer{
		wr: bufio.NewWriter(wr),

		lenBuf: make([]byte, 64),
		numBuf: make([]byte, 64),
	}
}

func (w *Writer) WriteObj(obj *object.Object) error {
	err := w.wr.WriteByte(byte(obj.Kind))
	if err != nil {
		return err
	}

	err = w.writeData(obj.Data)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) writeData(v interface{}) error {
	switch v := v.(type) {
	case string:
		return w.string(v)
	case []byte:
		return w.bytes(v)
	case int:
		return w.int(int64(v))
	case int8:
		return w.int(int64(v))
	case int16:
		return w.int(int64(v))
	case int32:
		return w.int(int64(v))
	case int64:
		return w.int(v)
	case uint:
		return w.uint(uint64(v))
	case uint8:
		return w.uint(uint64(v))
	case uint16:
		return w.uint(uint64(v))
	case uint32:
		return w.uint(uint64(v))
	case uint64:
		return w.uint(v)
	case float32:
		return w.float(float64(v))
	case float64:
		return w.float(v)
	case bool:
		if v {
			return w.wr.WriteByte(1)
		}
		return w.wr.WriteByte(0)
	case list.List:
		ls := list.List(v)
		for e := ls.Front(); e != nil; e = e.Next() {
			w.writeData(e.Value)
		}
	}
	return fmt.Errorf("redis: can't marshal %T (implement encoding.BinaryMarshaler)", v)
}

func (w *Writer) writeLen(n int) error {
	w.lenBuf = strconv.AppendUint(w.lenBuf[:0], uint64(n), 10)
	w.lenBuf = append(w.lenBuf, '\r', '\n')
	_, err := w.wr.Write(w.lenBuf)
	return err
}

func (w *Writer) bytes(b []byte) error {
	err := w.writeLen(len(b))
	if err != nil {
		return err
	}

	_, err = w.wr.Write(b)
	if err != nil {
		return err
	}

	return w.crlf()
}

func (w *Writer) string(s string) error {
	return w.bytes(util.StringToBytes(s))
}

func (w *Writer) uint(n uint64) error {
	w.numBuf = strconv.AppendUint(w.numBuf[:0], n, 10)
	return w.bytes(w.numBuf)
}

func (w *Writer) int(n int64) error {
	w.numBuf = strconv.AppendInt(w.numBuf[:0], n, 10)
	return w.bytes(w.numBuf)
}

func (w *Writer) float(f float64) error {
	w.numBuf = strconv.AppendFloat(w.numBuf[:0], f, 'f', -1, 64)
	return w.bytes(w.numBuf)
}

func (w *Writer) crlf() error {
	err := w.wr.WriteByte('\r')
	if err != nil {
		return err
	}
	return w.wr.WriteByte('\n')
}

func (w *Writer) Buffered() int {
	return w.wr.Buffered()
}

func (w *Writer) Reset(wr io.Writer) {
	w.wr.Reset(wr)
}

func (w *Writer) Flush() error {
	return w.wr.Flush()
}
