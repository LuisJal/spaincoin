package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/models"
)

// writeJSON encodes v as JSON and writes it to w with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// HandleStatus handles GET /api/status.
// Returns combined exchange + node status.
func HandleStatus(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to reach node: "+err.Error())
			return
		}

		resp := models.ExchangeStatus{
			Exchange:         "SpainCoin Exchange",
			Version:          "0.1.0",
			Node:             nodeStatus,
			BlockTimeSeconds: 5,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// HandleLatestBlocks handles GET /api/blocks/latest.
// Returns the last 10 blocks.
func HandleLatestBlocks(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blocks, err := nodeClient.GetRecentBlocks(10)
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to fetch blocks: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, blocks)
	}
}

// HandleBlock handles GET /api/blocks/{height}.
// Returns the block at the given height.
func HandleBlock(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract {height} from the URL path.
		// Path is /api/blocks/{height}
		path := r.URL.Path
		parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
		// parts: ["api", "blocks", "{height}"]
		if len(parts) < 3 {
			writeError(w, http.StatusBadRequest, "missing block height")
			return
		}
		heightStr := parts[2]
		height, err := strconv.ParseUint(heightStr, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid block height: "+heightStr)
			return
		}

		block, err := nodeClient.GetBlock(height)
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to fetch block: "+err.Error())
			return
		}
		writeJSON(w, http.StatusOK, block)
	}
}

// explorerResponse is the response body for GET /api/explorer.
type explorerResponse struct {
	Blocks         []*models.BlockInfo `json:"blocks"`
	Height         uint64              `json:"height"`
	TotalSupply    uint64              `json:"total_supply"`
	TotalSupplySPC float64             `json:"total_supply_spc"`
}

// HandleExplorer handles GET /api/explorer.
// Returns the last 10 blocks plus chain-level stats.
func HandleExplorer(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to reach node: "+err.Error())
			return
		}

		blocks, err := nodeClient.GetRecentBlocks(10)
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to fetch blocks: "+err.Error())
			return
		}

		resp := explorerResponse{
			Blocks:         blocks,
			Height:         nodeStatus.Height,
			TotalSupply:    nodeStatus.TotalSupply,
			TotalSupplySPC: float64(nodeStatus.TotalSupply) / 1_000_000_000_000_000.0,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}
