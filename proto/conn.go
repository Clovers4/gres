package proto

import (
	"context"
	"net"
	"time"
)

var noDeadline = time.Time{}

type Conn struct {
	netConn net.Conn

	rd *Reader
	wr *Writer
}

func NewConn(netConn net.Conn) *Conn {
	cn := &Conn{
		netConn: netConn,
	}
	cn.rd = NewReader(netConn)
	cn.wr = NewWriter(netConn)
	return cn
}

func (cn *Conn) SetNetConn(netConn net.Conn) {
	cn.netConn = netConn
	cn.rd.Reset(netConn)
	cn.wr.Reset(netConn)
}

func (cn *Conn) Write(b []byte) (int, error) {
	return cn.netConn.Write(b)
}

func (cn *Conn) RemoteAddr() net.Addr {
	return cn.netConn.RemoteAddr()
}

// todo:test
func (cn *Conn) WithReader(ctx context.Context, timeout time.Duration, fn func(rd *Reader) error) error {
	err := cn.netConn.SetReadDeadline(cn.deadline(ctx, timeout))
	if err != nil {
		return err
	}
	return fn(cn.rd)
}

// todo:test
func (cn *Conn) WithWriter(ctx context.Context, timeout time.Duration, fn func(wr *Writer) error) error {
	err := cn.netConn.SetWriteDeadline(cn.deadline(ctx, timeout))
	if err != nil {
		return err
	}

	if cn.wr.Buffered() > 0 {
		cn.wr.Reset(cn.netConn)
	}

	err = fn(cn.wr)
	if err != nil {
		return err
	}

	return cn.wr.Flush()
}

func (cn *Conn) Close() error {
	return cn.netConn.Close()
}

func (cn *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	tm := time.Now()
	if timeout > 0 {
		tm = tm.Add(timeout)
	}

	if ctx != nil {
		deadline, ok := ctx.Deadline()
		if ok {
			if timeout == 0 {
				return deadline
			}
			if deadline.Before(tm) {
				return deadline
			}
			return tm
		}
	}

	if timeout > 0 {
		return tm
	}
	return noDeadline
}
