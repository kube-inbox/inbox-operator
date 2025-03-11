package api

import (
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Server represents the HTTP server for the API
type Server struct {
	handler *InboxHandler
	logger  logr.Logger
}

// NewServer creates a new API server
func NewServer(client client.Client, logger logr.Logger) *Server {
	return &Server{
		handler: NewInboxHandler(client, logger),
		logger:  logger,
	}
}

// setupCORS adds CORS headers to the response
func setupCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Start starts the HTTP server
func (s *Server) Start(port int) error {
	mux := http.NewServeMux()

	// Register API routes with CORS middleware
	mux.HandleFunc("/api/v1/inboxes", setupCORS(s.handler.ListInboxes))
	mux.HandleFunc("/api/v1/inbox", setupCORS(s.handler.GetInbox))

	addr := fmt.Sprintf(":%d", port)
	s.logger.Info("Starting API server", "addr", addr)

	return http.ListenAndServe(addr, mux)
}
