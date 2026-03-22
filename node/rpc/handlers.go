package rpc

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/chain"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// ---------------------------------------------------------------------------
// Helper types for JSON responses
// ---------------------------------------------------------------------------

type txJSON struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    uint64 `json:"amount"`
	Nonce     uint64 `json:"nonce"`
	Fee       uint64 `json:"fee"`
	Timestamp int64  `json:"timestamp"`
	SigR      string `json:"sig_r,omitempty"`
	SigS      string `json:"sig_s,omitempty"`
}

type blockJSON struct {
	Height       uint64   `json:"height"`
	Hash         string   `json:"hash"`
	PrevHash     string   `json:"prev_hash"`
	MerkleRoot   string   `json:"merkle_root"`
	StateRoot    string   `json:"state_root"`
	Timestamp    int64    `json:"timestamp"`
	Validator    string   `json:"validator"`
	TxCount      int      `json:"tx_count"`
	Transactions []txJSON `json:"transactions"`
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// writeJSON serialises v as JSON, sets Content-Type and the given HTTP status.
// CORS and security headers are applied by the server-level middleware.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a standard {"error": "msg"} response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// blockToJSON converts a *block.Block to the JSON-serialisable struct.
func blockToJSON(b *block.Block) blockJSON {
	txs := make([]txJSON, len(b.Transactions))
	for i, tx := range b.Transactions {
		t := txJSON{
			ID:        tx.ID.String(),
			From:      tx.From.String(),
			To:        tx.To.String(),
			Amount:    tx.Amount,
			Nonce:     tx.Nonce,
			Fee:       tx.Fee,
			Timestamp: tx.Timestamp,
		}
		if tx.Signature != nil {
			t.SigR = hex.EncodeToString(tx.Signature.R.Bytes())
			t.SigS = hex.EncodeToString(tx.Signature.S.Bytes())
		}
		txs[i] = t
	}
	return blockJSON{
		Height:       b.Header.Height,
		Hash:         b.Hash.String(),
		PrevHash:     b.Header.PrevHash.String(),
		MerkleRoot:   b.Header.MerkleRoot.String(),
		StateRoot:    b.Header.StateRoot.String(),
		Timestamp:    b.Header.Timestamp,
		Validator:    b.Header.Validator.String(),
		TxCount:      len(b.Transactions),
		Transactions: txs,
	}
}

// ---------------------------------------------------------------------------
// Handlers
// OPTIONS pre-flight is handled upstream by the security middleware.
// ---------------------------------------------------------------------------

// handleStatus serves GET /status.
func handleStatus(bc *chain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		latest := bc.LatestBlock()
		resp := map[string]interface{}{
			"status":       "ok",
			"height":       latest.Header.Height,
			"latest_hash":  latest.Hash.String(),
			"total_supply": bc.State().TotalSupply(),
			"mempool_size": bc.Mempool().Size(),
			"peer_count":   0,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// handleBlock serves GET /block/{height} and GET /block/latest.
// The mux is registered at "/block/" so we parse the suffix ourselves.
func handleBlock(bc *chain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Strip prefix "/block/" and trim trailing slash.
		suffix := strings.TrimPrefix(r.URL.Path, "/block/")
		suffix = strings.TrimSuffix(suffix, "/")

		if suffix == "latest" {
			handleLatestBlock(bc)(w, r)
			return
		}

		height, err := strconv.ParseUint(suffix, 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid block height")
			return
		}

		b, ok := bc.GetBlock(height)
		if !ok {
			writeError(w, http.StatusNotFound, "block not found")
			return
		}

		writeJSON(w, http.StatusOK, blockToJSON(b))
	}
}

// handleLatestBlock serves GET /block/latest.
func handleLatestBlock(bc *chain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		b := bc.LatestBlock()
		writeJSON(w, http.StatusOK, blockToJSON(b))
	}
}

// handleBalance serves GET /address/{address}/balance.
// The mux is registered at "/address/" so we parse the path ourselves.
func handleBalance(bc *chain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// Expected path: /address/{address}/balance
		// Strip "/address/" prefix.
		suffix := strings.TrimPrefix(r.URL.Path, "/address/")
		// suffix should be "{address}/balance"
		parts := strings.SplitN(suffix, "/", 2)
		if len(parts) != 2 || parts[1] != "balance" {
			writeError(w, http.StatusNotFound, "not found")
			return
		}

		addrStr := parts[0]
		addr, err := crypto.AddressFromHex(addrStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid address")
			return
		}

		st := bc.State()
		balance := st.GetBalance(addr)
		var nonce uint64
		if acc, ok := st.GetAccount(addr); ok {
			nonce = acc.Nonce
		}

		resp := map[string]interface{}{
			"address": addr.String(),
			"balance": balance,
			"nonce":   nonce,
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// sendTxRequest is the JSON body for POST /tx/send.
type sendTxRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
	Nonce  uint64 `json:"nonce"`
	Fee    uint64 `json:"fee"`
	SigR   string `json:"sig_r"`
	SigS   string `json:"sig_s"`
}

// handleSendTx serves POST /tx/send.
func handleSendTx(bc *chain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req sendTxRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
			return
		}

		fromAddr, err := crypto.AddressFromHex(req.From)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid from address: "+err.Error())
			return
		}

		toAddr, err := crypto.AddressFromHex(req.To)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid to address: "+err.Error())
			return
		}

		// Build the transaction.
		tx := &block.Transaction{
			From:      fromAddr,
			To:        toAddr,
			Amount:    req.Amount,
			Nonce:     req.Nonce,
			Fee:       req.Fee,
			Timestamp: time.Now().UnixNano(),
		}
		tx.ID = tx.Hash()

		// Attach signature if provided.
		if req.SigR != "" || req.SigS != "" {
			rBytes, err := hex.DecodeString(req.SigR)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid sig_r: "+err.Error())
				return
			}
			sBytes, err := hex.DecodeString(req.SigS)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid sig_s: "+err.Error())
				return
			}
			tx.Signature = &crypto.Signature{
				R: new(big.Int).SetBytes(rBytes),
				S: new(big.Int).SetBytes(sBytes),
			}
		}

		if err := bc.Mempool().Add(tx); err != nil {
			writeError(w, http.StatusUnprocessableEntity, "mempool rejected tx: "+err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"tx_id":  tx.ID.String(),
			"status": "accepted",
		})
	}
}

// handleGetTx serves GET /tx/{hash}.
// The mux is registered at "/tx/" so we parse the suffix ourselves.
// Note: "/tx/send" is registered separately and takes precedence for that exact path.
func handleGetTx(bc *chain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		hashStr := strings.TrimPrefix(r.URL.Path, "/tx/")
		hashStr = strings.TrimSuffix(hashStr, "/")

		if hashStr == "" {
			writeError(w, http.StatusBadRequest, "missing tx hash")
			return
		}

		h, err := crypto.HashFromHex(hashStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid tx hash")
			return
		}

		tx, ok := bc.Mempool().Get(h)
		if !ok {
			writeError(w, http.StatusNotFound, "transaction not found")
			return
		}

		t := txJSON{
			ID:        tx.ID.String(),
			From:      tx.From.String(),
			To:        tx.To.String(),
			Amount:    tx.Amount,
			Nonce:     tx.Nonce,
			Fee:       tx.Fee,
			Timestamp: tx.Timestamp,
		}
		if tx.Signature != nil {
			t.SigR = hex.EncodeToString(tx.Signature.R.Bytes())
			t.SigS = hex.EncodeToString(tx.Signature.S.Bytes())
		}

		writeJSON(w, http.StatusOK, t)
	}
}

// handleValidators serves GET /validators.
func handleValidators(bc *chain.Blockchain) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		// The blockchain does not expose a ValidatorSet directly.
		// Return an empty list — validators are managed at the node level.
		type validatorJSON struct {
			Address string `json:"address"`
			Stake   uint64 `json:"stake"`
		}
		resp := map[string]interface{}{
			"count":      0,
			"validators": []validatorJSON{},
		}
		writeJSON(w, http.StatusOK, resp)
	}
}
