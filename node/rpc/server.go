package rpc

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spaincoin/spaincoin/core/chain"
)

// ---------------------------------------------------------------------------
// Rate limiter (token bucket per IP, stdlib only)
// ---------------------------------------------------------------------------

const (
	rpcRateLimit   = 60 // max requests per window
	rpcRateWindow  = 1 * time.Minute
	maxRequestBody = 1 << 20 // 1 MB
)

type bucket struct {
	tokens    int
	lastReset time.Time
	mu        sync.Mutex
}

type rateLimiter struct {
	buckets sync.Map
}

func newRateLimiter() *rateLimiter {
	rl := &rateLimiter{}
	// Cleanup goroutine — removes stale buckets every minute.
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.buckets.Range(func(key, value interface{}) bool {
				b := value.(*bucket)
				b.mu.Lock()
				stale := time.Since(b.lastReset) > 2*rpcRateWindow
				b.mu.Unlock()
				if stale {
					rl.buckets.Delete(key)
				}
				return true
			})
		}
	}()
	return rl
}

// allow returns true if the given IP is within the rate limit.
func (rl *rateLimiter) allow(ip string) bool {
	val, _ := rl.buckets.LoadOrStore(ip, &bucket{tokens: rpcRateLimit, lastReset: time.Now()})
	b := val.(*bucket)
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	if now.Sub(b.lastReset) >= rpcRateWindow {
		b.tokens = rpcRateLimit
		b.lastReset = now
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// ---------------------------------------------------------------------------
// Middleware helpers
// ---------------------------------------------------------------------------

// clientIP extracts the real client IP from the request.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	// Strip port from RemoteAddr.
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// setSecurityHeaders adds hardened security headers to every response.
func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Content-Security-Policy", "default-src 'none'")
}

// allowedOrigin returns the origin that is permitted to call this RPC,
// read from SPC_ALLOWED_ORIGIN (default: http://localhost:3001).
func allowedOrigin() string {
	if v := os.Getenv("SPC_ALLOWED_ORIGIN"); v != "" {
		return v
	}
	return "http://localhost:3001"
}

// ---------------------------------------------------------------------------
// Server
// ---------------------------------------------------------------------------

// Server is an HTTP JSON-RPC server that exposes the blockchain state and
// transaction submission to external clients (exchange app, CLI).
type Server struct {
	chain  *chain.Blockchain
	addr   string
	server *http.Server
	rl     *rateLimiter
}

// NewServer creates a new RPC Server bound to addr (e.g. ":8545").
func NewServer(bc *chain.Blockchain, addr string) *Server {
	s := &Server{
		chain: bc,
		addr:  addr,
		rl:    newRateLimiter(),
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
		Handler: s.securityMiddleware(mux),
	}

	return s
}

// securityMiddleware wraps the mux with rate limiting, CORS enforcement,
// security headers, and request body size limiting.
func (s *Server) securityMiddleware(next http.Handler) http.Handler {
	origin := allowedOrigin()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Security headers on every response.
		setSecurityHeaders(w)

		// 2. Enforce body size limit (1 MB).
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBody)

		// 3. Rate limiting per IP.
		ip := clientIP(r)
		if !s.rl.allow(ip) {
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}

		// 4. CORS — only allow the configured exchange server origin.
		//    Preflight OPTIONS requests must also be validated.
		reqOrigin := r.Header.Get("Origin")
		if reqOrigin != "" {
			if reqOrigin != origin {
				log.Printf("CORS: rejected origin %q from %s", reqOrigin, ip)
				http.Error(w, `{"error":"origin not allowed"}`, http.StatusForbidden)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Vary", "Origin")
		}

		// Handle pre-flight OPTIONS.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
