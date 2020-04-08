package gres

import (
	"context"
	"fmt"
	"github.com/clovers4/gres/commands"
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/proto"
	"net"
	"strings"
)

const (
	crlf = "\r\n"
)

type Client struct {
	ctx context.Context

	conn *proto.Conn
	db   *engine.DB // Pointer to currently selected DB
	srv  *Server
}

func NewClient(netConn net.Conn, srv *Server) *Client {
	conn := proto.NewConn(netConn)
	ctx := context.Background()
	ctx = engine.CtxWithDB(ctx, srv.db)

	cli := &Client{
		ctx:  ctx,
		conn: conn,
		db:   srv.db,
		srv:  srv,
	}

	srv.mu.Lock()
	srv.clients = append(srv.clients, cli)
	srv.mu.Unlock()
	return cli
}

func (cli *Client) Interact() {
	var err error
	conn := cli.conn
	ctx := cli.ctx
	for {
		var reply *proto.Reply
		var quit bool
		err = conn.WithReader(ctx, 0, func(rd *proto.Reader) error { // todo:time
			input, err := rd.ReadReply()
			if err != nil {
				return err
			}

			// todo
			fmt.Printf("read %v args: %v\n", conn.RemoteAddr(), input)

			// todo: cmd like "quit"
			var args []string
			args, ok := input.([]string)
			if !ok || len(args) == 0 {
				return fmt.Errorf("incorrect input: %v", input)
			}

			args[0] = strings.ToLower(args[0])
			if args[0] == "quit" {
				quit = true
				return nil
			}

			cmd := commands.GetCmd(args[0])
			if cmd == nil {
				return fmt.Errorf("cannot find cmd: %v", args[0])
			}
			reply = cmd.Do(ctx, args)

			return nil
		})

		fmt.Println("read finish") // todo
		if quit == true {
			break
		}

		err = conn.WithWriter(ctx, 0, func(wr *proto.Writer) error { // todo:time
			if err != nil {
				return wr.ReplyErr(err)
			}
			return wr.Reply(reply)
		})
		fmt.Printf("write %v reply: %v\n", conn.RemoteAddr(), reply) // todo

		// todo: log err
	}
}

func (cli *Client) Close() {
	cli.conn.Close()

	srv := cli.srv
	srv.mu.Lock()
	for i, c := range srv.clients {
		if c == cli {
			// we dont care the order
			srv.clients[i] = nil
			srv.clients[i] = srv.clients[len(srv.clients)-1]
			srv.clients = srv.clients[:len(srv.clients)-1]
			break
		}
	}
	srv.mu.Unlock()
}
