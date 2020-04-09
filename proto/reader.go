package proto

import (
	"bufio"
	"fmt"
	"io"

	"github.com/clovers4/gres/util"
)

const Nil = RedisError("redis: nil")

type ArraysHook func(*Reader, int64) (interface{}, error)

type Reader struct {
	rd *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReader(rd),
	}
}

// if err==nil, interface{} will be string/int/[]string
func (r *Reader) ReadReply() (interface{}, error) {
	line, err := r.ReadLine()
	if err != nil {
		return nil, err
	}

	switch line[0] {
	case ErrReply:
		return parseErrorReply(line), nil
	case StatusReply:
		return string(line[1:]), nil
	case IntReply:
		return util.ParseInt(line[1:], 10, 64)
	case BulkStringReply:
		return r.readBulkStringReply(line)
	case ArraysReply:
		n, err := parseArrayLen(line)
		if err != nil {
			return nil, err
		}
		vals := make([]string, n, n)
		for i := 0; i < int(n); i++ {
			v, err := r.ReadReply()
			if err != nil {
				return nil, err
			}
			vals[i] = v.(string)
		}
		return vals, nil
	}
	return nil, fmt.Errorf("redis: type is incorrect %.100q", line)
}

func (r *Reader) ReadLine() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	if isNilReply(line) {
		return nil, Nil
	}
	return line, nil
}

// readLine that returns an error if:
//   - there is a pending read error;
//   - or line does not end with \r\n.
func (r *Reader) readLine() ([]byte, error) {
	b, err := r.rd.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("redis: invalid reply: %v", b)
	}
	b = b[:len(b)-2]
	return b, nil
}

func (r *Reader) readBulkStringReply(line []byte) (string, error) {
	if isNilReply(line) {
		return Nil.Error(), nil
	}

	replyLen, err := util.Atoi(line[1:])
	if err != nil {
		return "", err
	}

	b := make([]byte, replyLen+2)
	_, err = io.ReadFull(r.rd, b)
	if err != nil {
		return "", err
	}

	return util.BytesToString(b[:replyLen]), nil
}

func (r *Reader) Buffered() int {
	return r.rd.Buffered()
}

func (r *Reader) Peek(n int) ([]byte, error) {
	return r.rd.Peek(n)
}

func (r *Reader) Reset(rd io.Reader) {
	r.rd.Reset(rd)
}

func isNilReply(b []byte) bool {
	return len(b) == 3 &&
		(b[0] == BulkStringReply || b[0] == ArraysReply) &&
		b[1] == '-' && b[2] == '1'
}

func parseArrayLen(line []byte) (int64, error) {
	if isNilReply(line) {
		return 0, Nil
	}
	return util.ParseInt(line[1:], 10, 64)
}

type RedisError string

func (e RedisError) Error() string { return string(e) }

func parseErrorReply(line []byte) error {
	return RedisError(string(line[1:]))
}
