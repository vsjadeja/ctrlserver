package ctrlserver

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const defaultControlPort int = 8080

type Server struct {
	httpLis net.Listener
	httpMux *http.ServeMux
	httpSrv *http.Server

	logLevelCb    LogLevelFunc
	logLevelSetCb SetLogLevelFunc
}

func New(addr string, opts ...Option) (*Server, error) {
	if addr == `` {
		addr = `:` + strconv.Itoa(defaultControlPort)
	}
	lis, err := net.Listen(`tcp`, addr)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		httpLis: lis,
		httpMux: http.NewServeMux(),
		httpSrv: &http.Server{
			Addr: addr,
		},
	}
	for _, opt := range opts {
		opt.apply(srv)
	}

	srv.httpMux.Handle(MetricsPath, promhttp.Handler())
	srv.httpMux.Handle(LogLevelPath, LogLevelHandler(srv.logLevelCb, srv.logLevelSetCb))
	// pprof handlers
	srv.AddHandler(PprofPath, pprof.Index)
	srv.AddHandler(PprofCmdlinePath, pprof.Cmdline)
	srv.AddHandler(PprofProfilePath, pprof.Profile)
	srv.AddHandler(PprofSymbolPath, pprof.Symbol)
	srv.AddHandler(PprofTracePath, pprof.Trace)

	return srv, nil
}

func (srv *Server) AddHandler(path string, h http.HandlerFunc) {
	srv.httpMux.Handle(path, h)
}

func (srv *Server) Serve() error {
	srv.httpSrv.Handler = srv.httpMux

	err := srv.httpSrv.Serve(srv.httpLis)
	if err == http.ErrServerClosed {
		return ErrServerClosed
	}
	return err
}

func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.httpSrv.Shutdown(ctx)
}

var (
	ErrServerClosed = errors.New(`control server closed`)
)
