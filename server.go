package gres

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/clovers4/gres/engine"
	"go.uber.org/zap"
)

const (
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

const (
	readSize = 4096
)

type Server struct {
	opts serverOptions
	// db
	db *engine.DB
	//	networking
	clients []*Client
	log     *zap.Logger

	mu sync.Mutex

	close bool
}

type serverOptions struct {
	configFile        string
	port              int
	connectionTimeout time.Duration
}

var defaultServerOptions = serverOptions{
	port:              9876,
	connectionTimeout: 120 * time.Second,
}

// A ServerOption sets options such as keepalive parameters, etc.
type ServerOption func(opts *serverOptions)

func (opt *serverOptions) readFlag() {
	// todo:
}

func (opt *serverOptions) readConfigFile() {
	// todo:
	if opt.configFile != "" {

	}
}

func ConfigFileOption(f string) ServerOption {
	return func(opts *serverOptions) {
		opts.configFile = f
	}
}

// PortOption set port of server.
// MAYBE just for test, benchmark, etc.
func PortOption(p int) ServerOption {
	return func(opts *serverOptions) {
		opts.port = p
	}
}

// ConnectionTimeoutOption set the time duration of connectionTimeout.
// The connection will be auto closed after the timeout
// MAYBE just for test, benchmark, etc.
func ConnectionTimeoutOption(t time.Duration) ServerOption {
	return func(opts *serverOptions) {
		opts.connectionTimeout = t
	}
}

// NewServer creates a gres server, ready to Serve.
func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o(&opts)
	}
	opts.readFlag()
	opts.readConfigFile()

	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("server-options:%+v", opts))

	srv := &Server{
		opts: opts,
		log:  log,
	}
	srv.db = engine.NewDB(
		engine.PersistOption(true),
		engine.LogOption(log))
	return srv
}

func (srv *Server) Start() {
	srv.listenExist()
	srv.listenAndServe()
}

// Serve creates a listener, accepts incoming connections and creates
// a service goroutine for each. The service goroutines read gres requests
// and then call the registered handlers(commands) to reply to them.
//
// Serve returns when listener.Accept fails with fatal errors. listener
// will be closed when this method returns.
//
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (srv *Server) listenAndServe() {
	// todo:	signal.Notify(quitCh, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", srv.opts.port) // net.ResolveTCPAddr(
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	srv.log.Info("gres start serving.", zap.String("addr", addr))

	for {
		conn, err := lis.Accept()
		if err != nil {
			srv.log.Error("listener.Accpet failed", zap.String("err", err.Error()))
			continue
		}
		go srv.handleConn(conn)
	}
}

func (srv *Server) listenExist() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				srv.log.Warn("[Server Stop] get exit signal, start exit", zap.String("signal", fmt.Sprintf("%v", s)))
				srv.Stop()
				close(c)
				os.Exit(1)
			default:
				srv.log.Warn("[Server Stop] get other signal", zap.String("signal", fmt.Sprintf("%v", s)))
			}
		}
	}()
}

func (srv *Server) handleConn(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			srv.log.Error("[Server] handleConn", zap.String("remote addr", conn.RemoteAddr().String()), zap.String("err", fmt.Sprintf("%v", err)))
		}
	}()

	srv.log.Info("[Server handleConn] Accept new connection", zap.String("remote addr", conn.RemoteAddr().String()))
	//conn.SetDeadline(time.Now().Set(srv.opts.connectionTimeout)) // todo: is correct? client side should be closed after deadline

	cli := NewClient(conn, srv)
	defer func() {
		if err := cli.Close(); err != nil {
			srv.log.Info("[Server handleConn] cli.Close()", zap.String("err", err.Error()))
		}
	}()
	cli.Interact()
}

func (srv *Server) freeMemoryIfNeed() {
	// todo: gopsutil?
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("%+v\n", ms)
}

func (srv *Server) Stop() {
	if srv.close {
		return
	}

	srv.close = true
	var err error
	for _, cli := range srv.clients {
		if err = cli.Close(); err != nil {
			srv.log.Warn("[Server Stop] cli.Close()", zap.String("err", err.Error()))
		}
	}
	if err = srv.log.Sync(); err != nil {
		srv.log.Warn("[Server Stop] log.Sync()", zap.String("err", err.Error()))
	}
	if err = srv.db.Close(); err != nil {
		srv.log.Warn("[Server Stop] db.Close()", zap.String("err", err.Error()))
	}
	srv.log.Info("[Server Stop] finished")
}
