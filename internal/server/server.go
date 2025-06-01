package server

import (
	"canary/internal/types"
	"context"
	"net/http"
	"time"

	"canary/internal/logger"
)

type Server types.Server

func New(addr string, handler http.Handler, log *logger.Logger) *Server {
	srv := Server{
		HttpServer: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		Logger: log,
	}
	return &srv
}

func (s *Server) Start() error {
	logger.Info("starting server", "addr", s.HttpServer.Addr)
	return s.HttpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Info("shutting down server")
	return s.HttpServer.Shutdown(ctx)
}
