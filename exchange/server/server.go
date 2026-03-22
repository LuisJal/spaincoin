package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/handlers"
)

// Server is the Exchange API HTTP server.
type Server struct {
	router     *http.ServeMux
	addr       string
	nodeURL    string
	httpServer *http.Server
}

// NewServer creates a new Server listening on addr and connecting to nodeURL.
func NewServer(addr, nodeURL string) *Server {
	s := &Server{
		router:  http.NewServeMux(),
		addr:    addr,
		nodeURL: nodeURL,
	}
	s.registerRoutes()
	return s
}

// registerRoutes sets up all API routes.
func (s *Server) registerRoutes() {
	nodeClient := client.NewNodeClient(s.nodeURL)

	mux := s.router

	// Wrap the mux with CORS + logging middleware
	// Routes:
	mux.HandleFunc("/api/status", handlers.HandleStatus(nodeClient))
	mux.HandleFunc("/api/blocks/latest", handlers.HandleLatestBlocks(nodeClient))
	mux.HandleFunc("/api/blocks/", handlers.HandleBlock(nodeClient))
	mux.HandleFunc("/api/explorer", handlers.HandleExplorer(nodeClient))
	mux.HandleFunc("/api/wallet/send", handlers.HandleSend(nodeClient))
	mux.HandleFunc("/api/wallet/", handlers.HandleWallet(nodeClient))
	mux.HandleFunc("/api/market/price", handlers.HandlePrice(nodeClient))
	mux.HandleFunc("/api/market/stats", handlers.HandleStats(nodeClient))
	mux.HandleFunc("/health", handleHealth)

	handler := corsMiddleware(loggingMiddleware(mux))

	s.httpServer = &http.Server{
		Addr:         s.addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// handleHealth is a simple liveness check endpoint.
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"status":"ok"}`)
}

// Start begins listening for requests. It blocks until the server stops.
func (s *Server) Start() error {
	log.Printf("Exchange API listening on %s", s.addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully shuts down the server using the provided context.
func (s *Server) Stop(ctx context.Context) error {
	log.Println("Shutting down Exchange API...")
	return s.httpServer.Shutdown(ctx)
}

// responseWriter wraps http.ResponseWriter to capture the status code for logging.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// loggingMiddleware logs the HTTP method, path, response status code, and duration.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
	})
}

// corsMiddleware adds CORS headers to allow all origins (for React dev).
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
