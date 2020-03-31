package gres

import (
	"fmt"
	"github.com/clovers4/gres/engine"
	"go.uber.org/zap"
	"net"
	"runtime"
	"time"
)

const (
	defaultWriteBufSize = 32 * 1024
	defaultReadBufSize  = 32 * 1024
)

const (
	readSize = 4096
)

// todo: default global server

type Server struct {
	opts serverOptions
	// db
	db []*engine.DB
	//	networking
	clients []*Client
	log     *zap.Logger
}

type serverOptions struct {
	configFile        string
	port              int
	dbnum             int
	connectionTimeout time.Duration
}

var defaultServerOptions = serverOptions{
	port:              9876,
	dbnum:             16,
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

// DbnumOption set the number of db in server.
// MAYBE just for test, benchmark, etc.
func DbnumOption(n int) ServerOption {
	return func(opts *serverOptions) {
		opts.dbnum = n
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
	srv.db = make([]*engine.DB, srv.opts.dbnum)
	for i := 0; i < opts.dbnum; i++ {
		srv.db[i] = engine.NewDB()
	}
	return srv
}

// Serve creates a listener, accepts incoming connections and creates
// a service goroutine for each. The service goroutines read gres requests
// and then call the registered handlers(commands) to reply to them.
//
// Serve returns when listener.Accept fails with fatal errors. listener
// will be closed when this method returns.
//
// Serve will return a non-nil error unless Stop or GracefulStop is called.
func (srv *Server) ListenAndServe() error {
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
	return nil
}

func (srv *Server) handleConn(conn net.Conn) {
	//defer func() {
	//	if err := recover(); err != nil {
	//		fmt.Println(err)	// todo: recover
	//		srv.log.Error("", zap.String("remote addr", conn.RemoteAddr().String()))
	//	}
	//}()
	srv.log.Info("Accept new connection", zap.String("remote addr", conn.RemoteAddr().String()))
	//conn.SetDeadline(time.Now().Add(srv.opts.connectionTimeout)) // todo: is correct? client side should be closed after deadline
	cli := NewClient(conn, srv)
	cli.Interact()
}

func (srv *Server) clientsCron() {
	// todo
}

func (srv *Server) serverCron() {

}

func (srv *Server) freeMemoryIfNeed() {
	// todo: gopsutil?
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("%+v\n", ms)
}

func (srv *Server) Stop() {
	for _, cli := range srv.clients {
		cli.Close()
	}
	srv.log.Sync()
}
