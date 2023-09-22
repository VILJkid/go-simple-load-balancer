// server/server.go
package server

import (
	"context"
	"net/http"
	"time"

	"log/slog"

	"github.com/VILJkid/go-simple-load-balancer/example-server/handlers"
)

type Server struct {
	Port string
	http.Server
}

func NewServer(port string) *Server {
	router := http.NewServeMux()
	router.HandleFunc("/", handlers.HelloHandler)

	return &Server{
		Port: port,
		Server: http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.ListenAndServe()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.WarnContext(ctx, "Server shutting down...")

	if err := s.Server.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "Server shutdown error:", "error", err)
		return
	}
	slog.InfoContext(ctx, "Server gracefully stopped")
}
