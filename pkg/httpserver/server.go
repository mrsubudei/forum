package httpserver

import (
	"context"
	v1 "forum/internal/controller/http/v1"
	"net/http"
	"time"
)

const (
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 5 * time.Second
	DefaultAddr            = ":8087"
	DefaultShutdownTimeout = 3 * time.Second
	ShutdownTimeout        = 5 * time.Second
)

type Server struct {
	httpServer *http.Server
	handler    *v1.Handler
}

func NewServer(handler *v1.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         DefaultAddr,
			ReadTimeout:  DefaultReadTimeout,
			WriteTimeout: DefaultWriteTimeout,
		},
		handler: handler,
	}
}

func (s *Server) Run() error {
	http.HandleFunc("/", s.handler.MainHandler)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
