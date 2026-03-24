package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// walletRegistry tracks all registered wallet addresses.
type walletRegistry struct {
	mu      sync.RWMutex
	wallets map[string]int64 // address → timestamp
	path    string
}

var registry *walletRegistry

// InitWalletRegistry loads or creates the wallet registry.
func InitWalletRegistry(dataDir string) {
	path := dataDir + "/wallets.json"
	r := &walletRegistry{
		wallets: make(map[string]int64),
		path:    path,
	}
	// Load existing
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &r.wallets)
	}
	registry = r
}

func (r *walletRegistry) register(address string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.wallets[address]; exists {
		return false
	}
	r.wallets[address] = time.Now().Unix()
	// Save to disk
	data, _ := json.Marshal(r.wallets)
	os.WriteFile(r.path, data, 0600)
	return true
}

func (r *walletRegistry) count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.wallets)
}

// HandleRegisterWallet handles POST /api/wallets/register.
func HandleRegisterWallet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req struct {
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid body")
			return
		}

		addr := strings.TrimSpace(req.Address)
		if len(addr) != 43 || !strings.HasPrefix(addr, "SPC") {
			writeError(w, http.StatusBadRequest, "invalid SPC address")
			return
		}

		if registry == nil {
			writeError(w, http.StatusInternalServerError, "registry not initialized")
			return
		}

		isNew := registry.register(addr)
		count := registry.count()

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"registered": isNew,
			"total":      count,
		})
	}
}

// HandleWalletCount handles GET /api/wallets/count.
func HandleWalletCount() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count := 0
		if registry != nil {
			count = registry.count()
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"total": count,
		})
	}
}
