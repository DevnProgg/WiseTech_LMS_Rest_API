package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"wisetech-lms-api/internal/config"
)

// Server holds the dependencies for the HTTP server
type Server struct {
	DB  *sql.DB
	Cfg *config.Config
}

// New creates a new Server instance
func New(db *sql.DB, cfg *config.Config) *Server {
	return &Server{
		DB:  db,
		Cfg: cfg,
	}
}

// Start runs the HTTP server
func (s *Server) Start() error {
	outer := s.NewRouter()

	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(s.Cfg.ServerPort),
		Handler:      outer,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Server listening on port %d\n", s.Cfg.ServerPort)
	return httpServer.ListenAndServe()
}
