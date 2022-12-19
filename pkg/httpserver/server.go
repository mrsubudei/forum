package httpserver

import (
	"context"
	"net/http"
	"time"

	v1 "forum/internal/controller/http/v1"
)

type Server struct {
	httpServer *http.Server
	h          *v1.Handler
}

const DefaultTime = int(time.Second)

func NewServer(handler *v1.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         handler.Cfg.Server.Port,
			ReadTimeout:  time.Duration(handler.Cfg.Server.ReadTimeout * DefaultTime),
			WriteTimeout: time.Duration(handler.Cfg.Server.WriteTimeout * DefaultTime),
			Handler:      handler.Mux,
		},
		h: handler,
	}
}

func (s *Server) Run() error {
	s.h.RegisterRoutes(s.h.Mux)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.h.Cfg.Server.ShutDownTimeout*DefaultTime))
	defer cancel()
	return s.httpServer.Shutdown(ctx)
}
