package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/models"
)

// isValidSPCAddress returns true if the address starts with "SPC" and is exactly 43 characters.
func isValidSPCAddress(addr string) bool {
	return strings.HasPrefix(addr, "SPC") && len(addr) == 43
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
			writeError(w, http.StatusBadRequest, "invalid SPC address: must start with 'SPC' and be 43 characters long")
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
// Validates the request and proxies it to the node POST /tx/send.
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

		if !isValidSPCAddress(req.From) {
			writeError(w, http.StatusBadRequest, "invalid 'from' address: must start with 'SPC' and be 43 characters long")
			return
		}
		if !isValidSPCAddress(req.To) {
			writeError(w, http.StatusBadRequest, "invalid 'to' address: must start with 'SPC' and be 43 characters long")
			return
		}

		resp, err := nodeClient.SendTx(&req)
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to send transaction: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, resp)
	}
}
