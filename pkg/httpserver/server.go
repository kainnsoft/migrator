package httpserver

import (
	"time"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/kainnsoft/migrator/config"
)

const (
	_defaultReadTimeout     = 50 * time.Second
	_defaultWriteTimeout    = 50 * time.Second
	_defaultAddr            = ":8080"
	_defaultShutdownTimeout = 30 * time.Second
)

// Server -.
type Server struct {
	server          *fasthttp.Server
	addr            string
	notify          chan error
	shutdownTimeout time.Duration
}

// New -.
func New(mux *router.Router, cfg config.HTTP) *Server {
	var (
		addr       = _defaultAddr
		httpServer *fasthttp.Server
		s          *Server
	)
	httpServer = &fasthttp.Server{
		Handler:      mux.Handler,
		ReadTimeout:  _defaultReadTimeout,
		WriteTimeout: _defaultWriteTimeout,
	}

	if cfg.HTTPAddr != _defaultAddr {
		addr = cfg.HTTPAddr
	}

	s = &Server{
		server:          httpServer,
		addr:            addr,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe(s.addr)
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}

func (s *Server) GetAddr() string {
	return s.addr
}
