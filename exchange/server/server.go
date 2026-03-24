package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spaincoin/spaincoin/exchange/auth"
	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/database"
	"github.com/spaincoin/spaincoin/exchange/handlers"
	"github.com/spaincoin/spaincoin/exchange/market"
)

// ---------------------------------------------------------------------------
// Rate limiter (token bucket per IP, stdlib only)
// ---------------------------------------------------------------------------

const (
	exchangeRateLimit  = 100 // max requests per window
	exchangeRateWindow = 1 * time.Minute
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
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			rl.buckets.Range(func(key, value interface{}) bool {
				b := value.(*bucket)
				b.mu.Lock()
				stale := time.Since(b.lastReset) > 2*exchangeRateWindow
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

func (rl *rateLimiter) allow(ip string) bool {
	val, _ := rl.buckets.LoadOrStore(ip, &bucket{tokens: exchangeRateLimit, lastReset: time.Now()})
	b := val.(*bucket)
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	if now.Sub(b.lastReset) >= exchangeRateWindow {
		b.tokens = exchangeRateLimit
		b.lastReset = now
	}

	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

// clientIP extracts the real client IP from the request.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// ---------------------------------------------------------------------------
// Server
// ---------------------------------------------------------------------------

// Server is the Exchange API HTTP server.
type Server struct {
	router     *http.ServeMux
	addr       string
	nodeURL    string
	httpServer *http.Server
	rl         *rateLimiter
	userDB     *database.UserDB
	tradeDB    *database.TradeDB
}

// NewServer creates a new Server listening on addr and connecting to nodeURL.
func NewServer(addr, nodeURL string, userDB *database.UserDB, tradeDB *database.TradeDB) *Server {
	s := &Server{
		router:  http.NewServeMux(),
		addr:    addr,
		nodeURL: nodeURL,
		rl:      newRateLimiter(),
		userDB:  userDB,
		tradeDB: tradeDB,
	}
	s.registerRoutes()
	return s
}

// registerRoutes sets up all API routes.
func (s *Server) registerRoutes() {
	nodeClient := client.NewNodeClient(s.nodeURL)
	sim := market.NewSimulator(0.09) // base price: €0.09

	// Start Binance price cache for real crypto prices
	handlers.InitPriceCache()

	// Wallet registry
	dataDir := os.Getenv("SPC_DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}
	handlers.InitWalletRegistry(dataDir)

	mux := s.router

	mux.HandleFunc("/api/status", handlers.HandleStatus(nodeClient))
	mux.HandleFunc("/api/blocks/latest", handlers.HandleLatestBlocks(nodeClient))
	mux.HandleFunc("/api/blocks/", handlers.HandleBlock(nodeClient))
	mux.HandleFunc("/api/explorer", handlers.HandleExplorer(nodeClient))
	mux.HandleFunc("/api/wallet/send", handlers.HandleSend(nodeClient))
	mux.HandleFunc("/api/wallet/", handlers.HandleWallet(nodeClient))
	mux.HandleFunc("/api/market/price", handlers.HandlePrice(nodeClient, sim))
	mux.HandleFunc("/api/market/stats", handlers.HandleStats(nodeClient, sim))
	mux.HandleFunc("/api/market/history", handlers.HandlePriceHistory(nodeClient, sim))
	mux.HandleFunc("/api/market/ticker", handlers.HandleTicker(nodeClient, sim))
	mux.HandleFunc("/api/market/table", handlers.HandleMarketTable(nodeClient, sim))
	mux.HandleFunc("/api/wallets/register", handlers.HandleRegisterWallet())
	mux.HandleFunc("/api/wallets/count", handlers.HandleWalletCount())
	mux.HandleFunc("/health", handleHealth)

	// Auth routes
	if s.userDB != nil {
		mux.HandleFunc("/api/auth/register", handlers.HandleRegister(s.userDB, s.tradeDB))
		mux.HandleFunc("/api/auth/login", handlers.HandleLogin(s.userDB, s.tradeDB))
		mux.HandleFunc("/api/auth/me", auth.AuthMiddleware(handlers.HandleMe(s.userDB, nodeClient)))
	}

	// Trading routes (auth required)
	if s.userDB != nil && s.tradeDB != nil {
		mux.HandleFunc("/api/trade/buy", auth.AuthMiddleware(handlers.HandleBuy(nodeClient, sim, s.userDB, s.tradeDB)))
		mux.HandleFunc("/api/trade/sell", auth.AuthMiddleware(handlers.HandleSell(nodeClient, sim, s.userDB, s.tradeDB)))
		mux.HandleFunc("/api/trade/history", auth.AuthMiddleware(handlers.HandleTradeHistory(s.tradeDB)))
		mux.HandleFunc("/api/trade/balance", auth.AuthMiddleware(handlers.HandleTradeBalance(nodeClient, s.userDB, s.tradeDB)))
		mux.HandleFunc("/api/trade/deposit-eur", auth.AuthMiddleware(handlers.HandleDepositEUR(s.tradeDB)))
		mux.HandleFunc("/api/trade/portfolio", auth.AuthMiddleware(handlers.HandlePortfolio(nodeClient, sim, s.userDB, s.tradeDB)))
	}

	handler := s.rateLimitMiddleware(
		securityHeadersMiddleware(
			corsMiddleware(
				jsonContentTypeMiddleware(
					loggingMiddleware(mux),
				),
			),
		),
	)

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

// ---------------------------------------------------------------------------
// Middleware
// ---------------------------------------------------------------------------

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

// corsMiddleware adds CORS headers. Keeps * for the React app but adds Vary: Origin.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Vary", "Origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// securityHeadersMiddleware adds hardened HTTP security headers to every response.
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "SAMEORIGIN")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		next.ServeHTTP(w, r)
	})
}

// jsonContentTypeMiddleware rejects POST requests whose Content-Type is not application/json.
func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			ct := r.Header.Get("Content-Type")
			if !strings.HasPrefix(ct, "application/json") {
				http.Error(w, `{"error":"Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware enforces the per-IP token-bucket rate limit.
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		if !s.rl.allow(ip) {
			http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
