package proto

import (
	"bufio"
	"fmt"
	"github.com/clovers4/gres/util"
	"io"
	"strconv"
)

const (
	StatusReply     = '+'
	ErrReply        = '-'
	IntReply        = ':'
	BulkStringReply = '$'
	ArraysReply     = '*'
)

type Writer struct {
	wr *bufio.Writer

	lenBuf []byte
	numBuf []byte
}

func NewWriter(wr io.Writer) *Writer {
	return &Writer{
		wr: bufio.NewWriter(wr),

		lenBuf: make([]byte, 128),
		numBuf: make([]byte, 128),
	}
}

func (w *Writer) Reply(reply *Reply) error {
	switch reply.Kind {
	case ReplyKindStatus:
		status := reply.Val.(string)
		return w.ReplyStatus(status)
	case ReplyKindErr:
		err := reply.Err.(error)
		return w.ReplyErr(err)
	case ReplyKindInt:
		i := reply.Val.(int)
		return w.ReplyInt(i)
	case ReplyKindBlukString:
		return w.ReplyBulkStringV(reply.Val)
	case ReplyKindArrays:
		arrays := reply.Val.([]interface{})
		return w.ReplyArrays(arrays)
	default:
		return fmt.Errorf("unknown type of reply")
	}
}

// 状态回复（status reply）的第一个字节是 "+"
func (w *Writer) ReplyStatus(status string) error {
	err := w.wr.WriteByte(StatusReply)
	if err != nil {
		return err
	}

	_, err = w.wr.Write(util.StringToBytes(status))
	if err != nil {
		return err
	}
	return w.crlf()
}

// 错误回复（error reply）的第一个字节是 "-"
func (w *Writer) ReplyErr(reply error) error {
	err := w.wr.WriteByte(ErrReply)
	if err != nil {
		return err
	}

	_, err = w.wr.Write(util.StringToBytes("ERR " + reply.Error()))
	if err != nil {
		return err
	}
	return w.crlf()
}

// 整数回复（integer reply）的第一个字节是 ":"
func (w *Writer) ReplyInt(num int) error {
	err := w.wr.WriteByte(IntReply)
	if err != nil {
		return err
	}

	w.numBuf = strconv.AppendInt(w.numBuf[:0], int64(num), 10)
	_, err = w.wr.Write(w.numBuf)
	if err != nil {
		return err
	}
	return w.crlf()
}

// 批量回复（bulk reply）的第一个字节是 "$"
func (w *Writer) ReplyBulkString(b []byte) error {
	err := w.wr.WriteByte(BulkStringReply)
	if err != nil {
		return err
	}

	err = w.writeLen(len(b))
	if err != nil {
		return err
	}

	_, err = w.wr.Write(b)
	if err != nil {
		return err
	}

	return w.crlf()
}

// 多条批量回复（multi bulk reply）的第一个字节是 "*"
func (w *Writer) ReplyArrays(args []interface{}) error {
	err := w.wr.WriteByte(ArraysReply)
	if err != nil {
		return err
	}

	err = w.writeLen(len(args))
	if err != nil {
		return err
	}

	for _, arg := range args {
		err := w.ReplyBulkStringV(arg)
		if err != nil {
			return err
		}
	}
	return nil
}

// length + n + crlf
func (w *Writer) writeLen(n int) error {
	w.lenBuf = strconv.AppendUint(w.lenBuf[:0], uint64(n), 10)
	w.lenBuf = append(w.lenBuf, '\r', '\n')
	_, err := w.wr.Write(w.lenBuf)
	return err
}

// length + v + crlf
func (w *Writer) ReplyBulkStringV(v interface{}) error {
	switch v := v.(type) {
	case nil:
		return w.string("")
	case string:
		return w.string(v)
	case *string:
		return w.string(*v)
	case []byte:
		return w.ReplyBulkString(v)
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
			return w.int(1)
		}
		return w.int(0)
	case *int:
		return w.int(int64(*v))
	case *int8:
		return w.int(int64(*v))
	case *int16:
		return w.int(int64(*v))
	case *int32:
		return w.int(int64(*v))
	case *int64:
		return w.int(*v)
	case *uint:
		return w.uint(uint64(*v))
	case *uint8:
		return w.uint(uint64(*v))
	case *uint16:
		return w.uint(uint64(*v))
	case *uint32:
		return w.uint(uint64(*v))
	case *uint64:
		return w.uint(*v)
	case *float32:
		return w.float(float64(*v))
	case *float64:
		return w.float(*v)
	case *bool:
		if *v {
			return w.int(1)
		}
		return w.int(0)
	default:
		return fmt.Errorf(
			"redis: can't marshal %T (implement encoding.BinaryMarshaler)", v)
	}
}

func (w *Writer) string(s string) error {
	return w.ReplyBulkString(util.StringToBytes(s))
}

func (w *Writer) uint(n uint64) error {
	w.numBuf = strconv.AppendUint(w.numBuf[:0], n, 10)
	return w.ReplyBulkString(w.numBuf)
}

func (w *Writer) int(n int64) error {
	w.numBuf = strconv.AppendInt(w.numBuf[:0], n, 10)
	return w.ReplyBulkString(w.numBuf)
}

func (w *Writer) float(f float64) error {
	w.numBuf = strconv.AppendFloat(w.numBuf[:0], f, 'f', -1, 64)
	return w.ReplyBulkString(w.numBuf)
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
