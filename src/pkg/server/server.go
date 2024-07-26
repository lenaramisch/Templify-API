package server

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

const READHEADERTIMEOUT = 5 * time.Second

func NewServer(cfg *Config, mux http.Handler) *Server {
	server := &http.Server{
		ReadHeaderTimeout: READHEADERTIMEOUT,
		Addr:              cfg.Address,
		Handler:           mux,
	}

	return &Server{
		server: server,
	}
}

func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctxShutDown context.Context) error {
	return s.server.Shutdown(ctxShutDown)
}
