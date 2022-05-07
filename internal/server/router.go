package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates a new chi router and sets up middleware and routes
func (s *Server) NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Health check endpoint
	r.Get("/health", s.healthCheck)

	return r
}

// healthCheck is a simple handler to check the service status
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
