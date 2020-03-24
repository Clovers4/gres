package gres

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

const (
	crlf = "\r\n"
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
	cli.args[0] = strings.ToLower(cli.args[0])
	cli.reply = ""

	cli.processCommand()

	return cli.reply
}

func (cli *Client) processCommand() {
	// lookup cmd
	cli.cmd = commands[cli.args[0]]

	if cli.cmd == nil {
		cli.addReplyError(ErrUnknownCmd)
		return
	}

	// execute cmd
	cli.cmd.Do(cli)
}

func (cli *Client) addReply(s string) {
	cli.reply += s
}

func (cli *Client) addReplyStatus() {
	// todo
	cli.addReply(fmt.Sprintf("+"))
}

func (cli *Client) addReplyError(err error) {
	cli.addReply(fmt.Sprintf("(error) %v", err.Error()))
}

func (cli *Client) setReplyInt(n int) {
	cli.addReply(fmt.Sprintf(":%d", n))
}

func (cli *Client) addReplyBulk() {

	// todo
}

func (cli *Client) addReplyMultiBulk() {
	// todo
}

func (cli *Client) setReplyOK() {
	cli.addReply("OK")
}

func (cli *Client) setReplyNull() {

}

func (cli *Client) setReplyList(ls []string) {
	var sb strings.Builder
	for i := 0; i < len(ls); i++ {
		fmt.Fprintf(&sb, "%d)\"%s\"\n", i, ls[i])
	}
	cli.addReply(sb.String())
}

func (cli *Client) Close() {
	cli.conn.Close()
}
