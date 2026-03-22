package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/models"
)

const (
	// maxSafeAmount is the maximum allowed transfer amount (uint64 safe max for 18-decimal token).
	maxSafeAmount uint64 = 18_000_000_000_000_000_000
	// maxReasonableNonce prevents obviously bogus replay-protection values.
	maxReasonableNonce uint64 = 1_000_000_000
)

// isValidSPCAddress returns true if addr satisfies all of:
//   - starts with "SPC"
//   - is exactly 43 characters total (SPC + 40 hex chars)
//   - the 40 trailing chars are valid lowercase hex (0-9, a-f)
func isValidSPCAddress(addr string) bool {
	if !strings.HasPrefix(addr, "SPC") {
		return false
	}
	if len(addr) != 43 {
		return false
	}
	hex := addr[3:]
	for _, c := range hex {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}

// isHexString returns true if s is a non-empty string containing only hex characters.
func isHexString(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// clientIP extracts the best-effort client IP from a request.
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

// HandleWallet handles GET /api/wallet/{address}.
// Returns balance info with BalanceSPC computed.
func HandleWallet(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract {address} from /api/wallet/{address}
		path := r.URL.Path
		parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
		// parts: ["api", "wallet", "{address}"]
		if len(parts) < 3 || parts[2] == "" {
			writeError(w, http.StatusBadRequest, "missing address")
			return
		}
		address := parts[2]

		if !isValidSPCAddress(address) {
			writeError(w, http.StatusBadRequest, "invalid SPC address: must start with 'SPC' followed by exactly 40 lowercase hex characters")
			return
		}

		balance, err := nodeClient.GetBalance(address)
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to fetch balance: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, balance)
	}
}

// HandleSend handles POST /api/wallet/send.
// Validates the request thoroughly and proxies it to the node POST /tx/send.
func HandleSend(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req models.SendTxRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body: "+err.Error())
			return
		}

		// --- Address validation ---
		if !isValidSPCAddress(req.From) {
			writeError(w, http.StatusBadRequest, "invalid 'from' address: must start with 'SPC' followed by exactly 40 lowercase hex characters")
			return
		}
		if !isValidSPCAddress(req.To) {
			writeError(w, http.StatusBadRequest, "invalid 'to' address: must start with 'SPC' followed by exactly 40 lowercase hex characters")
			return
		}

		// --- Amount validation ---
		if req.Amount == 0 {
			writeError(w, http.StatusBadRequest, "amount must be greater than 0")
			return
		}
		if req.Amount > maxSafeAmount {
			writeError(w, http.StatusBadRequest, "amount exceeds maximum allowed value")
			return
		}

		// --- Fee validation ---
		// Fee is uint64 so it can't be negative; only document the constraint.
		// (A zero fee is technically valid — the mempool may reject it separately.)

		// --- Nonce validation ---
		if req.Nonce >= maxReasonableNonce {
			writeError(w, http.StatusBadRequest, "nonce value is unreasonably large")
			return
		}

		// --- Signature validation ---
		if !isHexString(req.SigR) {
			writeError(w, http.StatusBadRequest, "sig_r must be a non-empty hex string")
			return
		}
		if !isHexString(req.SigS) {
			writeError(w, http.StatusBadRequest, "sig_s must be a non-empty hex string")
			return
		}

		// Audit log: record every send attempt with IP, sender, and amount.
		ip := clientIP(r)
		log.Printf("[AUDIT] send attempt ip=%s from=%s to=%s amount=%d fee=%d nonce=%d",
			ip, req.From, req.To, req.Amount, req.Fee, req.Nonce)

		resp, err := nodeClient.SendTx(&req)
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to send transaction: "+err.Error())
			return
		}

		log.Printf("[AUDIT] send accepted ip=%s from=%s tx_id=%s", ip, req.From, resp.TxID)
		writeJSON(w, http.StatusOK, resp)
	}
}
