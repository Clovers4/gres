package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const readSize = 4096

var (
	port = flag.Int("p", 9876, "specify port to use.  defaults to 9876.")
	host = flag.String("h", "127.0.0.1", "specify host to use.  defaults to 127.0.0.1.")
)

func init(){
	flag.Parse()//todo
}

// Client represents another side of Server, and is not the same as
// gres/Client.
type Client struct {
	opts clientOptions
	conn net.Conn
	scan *bufio.Scanner
	rb   []byte
	done chan struct{} //todo:
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

	conn, err := net.DialTimeout("tcp", opts.remoteAddr, opts.connectionTimeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dail err=%v", err)
		os.Exit(0)
	}

	return &Client{
		opts: opts,
		conn: conn,
		scan: bufio.NewScanner(os.Stdin),
		rb:   make([]byte, readSize),
		done: make(chan struct{}),
	}
}

func initFlag(opts *clientOptions) {
	opts.remoteAddr = fmt.Sprintf("%s:%d", *host, *port)
}

func (cli *Client) Interact() {
	for {
		// write
		fmt.Printf("%v> ", cli.conn.RemoteAddr())
		if !cli.scan.Scan() {
			break
		}
		command := cli.scan.Text()
		if strings.ToLower(command) == "exit" { //todo
			return
		}
		output, err := cli.interact(command) //todo
		if err != nil {
			panic(err) //todo:??
		}
		fmt.Println(output)
	}
	if err := cli.scan.Err(); err != nil {
		fmt.Printf("reading standard input failed, err=%v", err)
	}
}

func (cli *Client) interact(input string) (output string, err error) {
	cli.conn.Write([]byte(input))
	n, err := cli.conn.Read(cli.rb)
	if err != nil {
		return "", err
	}
	if n > 0 {
		output = string(cli.rb[:n])
		cli.rb = cli.rb[n:]
	}

	if len(cli.rb) == 0 {
		cli.rb = make([]byte, readSize)
	}
	return output, nil
}

// GracefulExit does some remaining work and will exit gracefully.
// It will close the connection, ... , etc.
func (cli *Client) GracefulExit() {
	// todo: modify graceful stop
	os.Exit(0)
}
