package health

import (
	"context"
	"log"
	"net/http"
	"time"
)

// Server wraps an HTTP server that exposes the health endpoint.
type Server struct {
	addr   string
	mux    *http.ServeMux
	server *http.Server
}

// NewServer creates a Server listening on addr, serving the health handler at /health.
func NewServer(addr string, h *Handler) *Server {
	mux := http.NewServeMux()
	mux.Handle("/health", h)

	return &Server{
		addr: addr,
		mux:  mux,
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

// Start begins listening in the background. It returns immediately.
func (s *Server) Start() {
	go func() {
		log.Printf("health: listening on %s", s.addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("health: server error: %v", err)
		}
	}()
}

// Shutdown gracefully stops the server with a 5-second deadline.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("health: shutdown error: %v", err)
	}
}
