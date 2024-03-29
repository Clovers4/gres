package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/clovers4/gres/proto"
	"go.uber.org/zap"
)

const readSize = 4096

var (
	port = flag.Int("p", 9876, "specify port to use.  defaults to 9876.")
	host = flag.String("h", "127.0.0.1", "specify host to use.  defaults to 127.0.0.1.")
)

func init() {
	flag.Parse() //todo
}

// Client represents another side of Server, and is not the same as
// gres/Client.
type Client struct {
	ctx  context.Context
	opts clientOptions
	conn *proto.Conn
	scan *bufio.Scanner
	rb   []byte

	log *zap.Logger
}

type clientOptions struct {
	remoteAddr        string        // Record the remote address, the form is host:port
	connectionTimeout time.Duration // Max timeout of connection
}

var defaultClientOptions = clientOptions{
	remoteAddr:        "127.0.0.1:9876",
	connectionTimeout: 5 * time.Second,
}

// getRemoteHostAndPort returns host and port.
// It uses strings.Split as a shortcut to get this two strings host and port
func (opt *clientOptions) getRemoteHostAndPort() (host, port string) {
	s := strings.Split(opt.remoteAddr, ":")
	return s[0], s[1]
}

func NewClient() *Client {
	opts := defaultClientOptions
	initFlag(&opts)

	netConn, err := net.DialTimeout("tcp", opts.remoteAddr, opts.connectionTimeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dail err=%v\n", err)
		os.Exit(0)
	}

	conn := proto.NewConn(netConn)

	return &Client{
		ctx:  context.Background(),
		opts: opts,
		conn: conn,
		scan: bufio.NewScanner(os.Stdin),
		rb:   make([]byte, readSize),
	}
}

func initFlag(opts *clientOptions) {
	opts.remoteAddr = fmt.Sprintf("%s:%d", *host, *port)
}

func (cli *Client) Start() {
	cli.listenExist()
	cli.Interact()
}

func (cli *Client) listenExist() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Printf("\n[Client listenExist] get exit signal=%v, start exit\n", fmt.Sprintf("%v", s))
				cli.GracefulExit()
				close(c)
				os.Exit(1)
			default:
				fmt.Printf("[Client listenExist] get other signal=%v, start exit\n", fmt.Sprintf("%v", s))
			}
		}
	}()
}

func (cli *Client) Interact() {
	for {
		// write
		fmt.Printf("%v> ", cli.conn.RemoteAddr())
		if !cli.scan.Scan() {
			break
		}
		command := cli.scan.Text()

		output, err := cli.interact(command)

		if err != nil {
			fmt.Println("client has errors:", err)
			break
		} else {
			fmt.Println(output)
		}

		name := strings.ToLower(command)
		if name == "exit" || name == "quit" {
			break
		}
	}

	if err := cli.scan.Err(); err != nil {
		fmt.Printf("reading standard input failed, err=%v", err)
	}
}

func (cli *Client) interact(input string) (output string, err error) {
	input = strings.TrimSpace(input)
	vals := strings.Split(input, " ")
	var args []interface{}
	for _, val := range vals {
		args = append(args, val)
	}

	err = cli.conn.WithWriter(cli.ctx, 0, func(wr *proto.Writer) error {
		return wr.ReplyArrays(args)
	})

	if err != nil {
		return "", err
	}

	err = cli.conn.WithReader(cli.ctx, 0, func(rd *proto.Reader) error {
		reply, err := rd.ReadReply()
		if err != nil {
			return err
		}
		if reply == "" {
			output = "(nil)"
		} else if _, ok := reply.(string); ok {
			output = fmt.Sprintf("\"%v\"", reply)
		} else if ss, ok := reply.([]string); ok {
			for i, s := range ss {
				output += fmt.Sprintf("%v) \"%v\"", i+1, s)
				if i < len(ss)-1 {
					output += "\n"
				}
			}
			if len(ss) == 0 {
				output += "(empty list or set)"
			}
		} else if i, ok := reply.(int64); ok {
			output += fmt.Sprintf("(integer) %v", i)
		} else {
			output = fmt.Sprintf("%v", reply)
		}
		return nil
	})
	return output, err
}

// GracefulExit does some remaining work and will exit gracefully.
// It will close the connection, ... , etc.
func (cli *Client) GracefulExit() {
	cli.conn.Close()
	fmt.Println("CLOSE connection done...")
	os.Exit(0)
}

func (cli *Client) BenchSet(t int) {
	for i := 0; i < t; i++ {
		v := strconv.FormatInt(int64(i), 10)
		command := fmt.Sprintf("set %v %v", v, v)

		output, err := cli.interact(command)

		if err != nil {
			fmt.Println("client has errors:", err)
			break
		} else {
			fmt.Println(output)
		}

		name := strings.ToLower(command)
		if name == "exit" || name == "quit" {
			break
		}
	}
}
