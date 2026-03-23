package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/spaincoin/spaincoin/core/crypto"
	"github.com/spaincoin/spaincoin/exchange/auth"
	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/database"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// isValidEmail performs a minimal email validation: must contain '@' and the
// domain part must contain at least one '.'.
func isValidEmail(email string) bool {
	at := strings.Index(email, "@")
	if at < 1 {
		return false
	}
	domain := email[at+1:]
	return strings.Contains(domain, ".")
}

// clientIPFromRequest extracts the best-effort client IP for audit logging.
func clientIPFromRequest(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.SplitN(xff, ",", 2)[0]
	}
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}

// ---------------------------------------------------------------------------
// Register
// ---------------------------------------------------------------------------

type registerRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	ImportKey string `json:"import_key,omitempty"` // optional: import existing private key hex
}

type authResponse struct {
	Token   string `json:"token"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

// HandleRegister handles POST /api/auth/register.
func HandleRegister(userDB *database.UserDB, tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIPFromRequest(r)

		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// Audit log — no password logged
		log.Printf("[AUDIT] register email=%s ip=%s", req.Email, ip)

		// Validate email
		if !isValidEmail(req.Email) {
			writeError(w, http.StatusBadRequest, "invalid email address")
			return
		}

		// Validate password length
		if len(req.Password) < 8 {
			writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
			return
		}

		// Hash password
		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			log.Printf("[ERROR] register bcrypt email=%s: %v", req.Email, err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Generate or import SpainCoin wallet
		var privKey *crypto.PrivateKey
		var pubKey *crypto.PublicKey
		if req.ImportKey != "" {
			// Import existing private key (founder wallet)
			var err2 error
			privKey, pubKey, err2 = crypto.PrivateKeyFromHex(req.ImportKey)
			if err2 != nil {
				writeError(w, http.StatusBadRequest, "clave privada inválida")
				return
			}
			log.Printf("[AUDIT] register import_wallet email=%s ip=%s", req.Email, ip)
		} else {
			var err2 error
			privKey, pubKey, err2 = crypto.GenerateKeyPair()
			if err2 != nil {
				log.Printf("[ERROR] register keygen email=%s: %v", req.Email, err2)
				writeError(w, http.StatusInternalServerError, "internal error")
				return
			}
		}
		address := pubKey.ToAddress().String()

		// Encrypt private key with user's password — never store plaintext
		encryptedKey, err := auth.EncryptPrivateKey(privKey.ToHex(), req.Password)
		if err != nil {
			log.Printf("[ERROR] register encrypt key email=%s: %v", req.Email, err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		record := &database.UserRecord{
			ID:            uuid.New().String(),
			Email:         req.Email,
			PasswordHash:  hash,
			WalletAddress: address,
			EncryptedKey:  encryptedKey,
		}

		if err := userDB.CreateUser(record); err != nil {
			if errors.Is(err, database.ErrEmailExists) {
				writeError(w, http.StatusConflict, "email already registered")
				return
			}
			log.Printf("[ERROR] register db email=%s: %v", req.Email, err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Generate JWT
		token, err := auth.GenerateToken(record.ID, record.Email)
		if err != nil {
			log.Printf("[ERROR] register jwt email=%s: %v", req.Email, err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Initialize testnet EUR balance (1000€ de regalo)
		if tradeDB != nil {
			if err := tradeDB.InitUserBalance(record.ID, 1000.0); err != nil {
				log.Printf("[WARN] register init EUR balance email=%s: %v", req.Email, err)
			}
		}

		log.Printf("[AUDIT] register success email=%s ip=%s address=%s", req.Email, ip, address)
		writeJSON(w, http.StatusCreated, authResponse{
			Token:   token,
			Address: address,
			Email:   record.Email,
		})
	}
}

// ---------------------------------------------------------------------------
// Login
// ---------------------------------------------------------------------------

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// HandleLogin handles POST /api/auth/login.
func HandleLogin(userDB *database.UserDB, tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := clientIPFromRequest(r)

		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// Audit log — no password logged
		log.Printf("[AUDIT] login email=%s ip=%s", req.Email, ip)

		record, err := userDB.GetByEmail(req.Email)
		if err != nil {
			// Return 401 for both not-found and wrong-password to avoid enumeration
			log.Printf("[AUDIT] login email=%s fail ip=%s (not found)", req.Email, ip)
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		if !auth.CheckPassword(req.Password, record.PasswordHash) {
			log.Printf("[AUDIT] login email=%s fail ip=%s (bad password)", req.Email, ip)
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		token, err := auth.GenerateToken(record.ID, record.Email)
		if err != nil {
			log.Printf("[ERROR] login jwt email=%s: %v", req.Email, err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}

		// Ensure EUR balance exists (for accounts created before trading)
		if tradeDB != nil {
			tradeDB.InitUserBalance(record.ID, 1000.0)
		}

		log.Printf("[AUDIT] login email=%s success ip=%s", req.Email, ip)
		writeJSON(w, http.StatusOK, authResponse{
			Token:   token,
			Address: record.WalletAddress,
			Email:   record.Email,
		})
	}
}

// ---------------------------------------------------------------------------
// Me (protected)
// ---------------------------------------------------------------------------

type meResponse struct {
	Email      string  `json:"email"`
	Address    string  `json:"address"`
	BalanceSPC float64 `json:"balance_spc"`
	Nonce      uint64  `json:"nonce"`
}

// HandleMe handles GET /api/auth/me.  The route must be wrapped with
// auth.AuthMiddleware before registration.
func HandleMe(userDB *database.UserDB, nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaims(r)
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		record, err := userDB.GetByID(claims.UserID)
		if err != nil {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		var balanceSPC float64
		var nonce uint64
		if balanceInfo, err := nodeClient.GetBalance(record.WalletAddress); err == nil {
			balanceSPC = balanceInfo.BalanceSPC
			nonce = balanceInfo.Nonce
		}

		writeJSON(w, http.StatusOK, meResponse{
			Email:      record.Email,
			Address:    record.WalletAddress,
			BalanceSPC: balanceSPC,
			Nonce:      nonce,
		})
	}
}
