package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pablo-puyat/elinor/internal/monitor"
)

type Server struct {
	port    int
	monitor *monitor.Monitor
	logger  *log.Logger
	server  *http.Server
}

func New(port int, monitor *monitor.Monitor, logger *log.Logger) *Server {
	return &Server{
		port:    port,
		monitor: monitor,
		logger:  logger,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", s.handleStats)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	s.logger.Printf("Starting API server on port %d", s.port)
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Printf("HTTP server error: %v", err)
	}
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Printf("HTTP server shutdown error: %v", err)
	}
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.monitor.GetStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
