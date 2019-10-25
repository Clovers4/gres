package gres

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn  net.Conn
	db    *DB // Pointer to currently selected DB
	srv   *Server
	args  []string
	cmd   Command
	reply string
}

func NewClient(conn net.Conn, srv *Server) *Client {
	cli := &Client{
		conn: conn,
		db:   srv.db[0],
		srv:  srv,
	}
	srv.clients = append(srv.clients, cli)
	return cli
}

func (cli *Client) Interact() {
	for buf := make([]byte, readSize); ; {
		// read
		n, err := cli.conn.Read(buf)
		var b []byte
		if n > 0 {
			b = bytes.NewBuffer(buf[:n:n]).Bytes()
			buf = buf[n:]
		}
		if err != nil {
			return
		}

		reply := cli.process(string(b))

		n, err = cli.conn.Write([]byte(reply))
		if err != nil {
			return
		}

		if len(buf) == 0 {
			buf = make([]byte, readSize)
		}
	}
}

// process returns reply.
func (cli *Client) process(input string) string {
	cli.args = strings.Fields(string(input))
	cli.args[0] = strings.ToUpper(cli.args[0])
	cli.processCommand()

	return cli.reply
}

func (cli *Client) processCommand() {
	// lookup cmd
	cli.cmd = cmdDict[cli.args[0]]
	if cli.cmd == nil {
		cli.setReplyError(ErrUnknownCmd)
		return
	}

	// execute cmd
	if err := cli.cmd.Do(cli); err != nil {
		cli.setReplyError(err)
	}
}

func (cli *Client) setReply(s string) {
	cli.reply = s
}

func (cli *Client) setReplyOK() {
	cli.setReply("OK")
}

func (cli *Client) setReplyError(err error) {
	cli.setReply(fmt.Sprintf("(error) %v", err.Error()))
}

func (cli *Client) setReplyNull(kind int) {
	switch kind {
	case OBJ_STRING:
		cli.setReply("(nil)")
	case OBJ_LIST:
		cli.setReply("(empty list)")
	case OBJ_SET:
		cli.setReply("(empty set)")
	case OBJ_ZSET:
		cli.setReply("(empty zset)")
	case OBJ_HASH:
		cli.setReply("(empty hash)")
	default:
		panic(fmt.Sprintf("Unknown obj type:%v", kind))
	}
}

func (cli *Client) setReplyInt(n int) {
	cli.setReply(fmt.Sprintf("(integer) %d", n))
}

func (cli *Client) setReplyList(ls []string) {
	var sb strings.Builder
	for i := 0; i < len(ls); i++ {
		fmt.Fprintf(&sb, "%d)\"%s\"\n", i, ls[i])
	}
	cli.setReply(sb.String())
}

func (cli *Client) Close() {
	cli.conn.Close()
}
