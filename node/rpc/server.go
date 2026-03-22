package rpc

import (
	"context"
	"net/http"

	"github.com/spaincoin/spaincoin/core/chain"
)

// Server is an HTTP JSON-RPC server that exposes the blockchain state and
// transaction submission to external clients (exchange app, CLI).
type Server struct {
	chain  *chain.Blockchain
	addr   string
	server *http.Server
}

// NewServer creates a new RPC Server bound to addr (e.g. ":8545").
func NewServer(bc *chain.Blockchain, addr string) *Server {
	s := &Server{
		chain: bc,
		addr:  addr,
	}

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/status", handleStatus(bc))
	// /block/latest must be registered before /block/{height} because the
	// pattern "/block/" is a prefix match; we disambiguate inside handleBlock.
	mux.HandleFunc("/block/", handleBlock(bc))
	mux.HandleFunc("/address/", handleBalance(bc))
	mux.HandleFunc("/tx/send", handleSendTx(bc))
	mux.HandleFunc("/tx/", handleGetTx(bc))
	mux.HandleFunc("/validators", handleValidators(bc))

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s
}

// Start begins listening in a background goroutine. It returns immediately.
// Errors during serving (other than http.ErrServerClosed) are silently dropped;
// the caller should use Stop to shut down cleanly.
func (s *Server) Start() error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Non-fatal: log is not imported here to keep the package lean.
			_ = err
		}
	}()
	return nil
}

// Stop shuts the HTTP server down gracefully, waiting at most until ctx is
// cancelled.
func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
