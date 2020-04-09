package gres

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/clovers4/gres/commands"
	"github.com/clovers4/gres/engine"
	"github.com/clovers4/gres/proto"
	"go.uber.org/zap"
)

const (
	crlf = "\r\n"
)

type Client struct {
	ctx context.Context

	conn *proto.Conn
	db   *engine.DB // Pointer to currently selected DB
	srv  *Server

	log *zap.Logger
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
		log:  srv.log,
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

			var args []string
			args, ok := input.([]string)
			if !ok || len(args) == 0 {
				return fmt.Errorf("incorrect input: %v", input)
			}

			args[0] = strings.ToLower(args[0])
			if args[0] == "quit" || args[0] == "exit" {
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

		if quit == true {
			break
		}

		err = conn.WithWriter(ctx, 0, func(wr *proto.Writer) error { // todo:time
			if err != nil {
				return wr.ReplyErr(err)
			}
			return wr.Reply(reply)
		})

		if err != nil {
			cli.log.Warn("[Client Interact] conn.WithWriter()", zap.String("err", err.Error()))
			break
		}
	}
}

func (cli *Client) Close() error {
	err := cli.conn.Close()

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
	cli.log.Info("End connection", zap.String("remote addr", cli.conn.RemoteAddr().String()))
	return err
}
