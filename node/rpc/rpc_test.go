package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/chain"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// newTestChain creates a Blockchain with a genesis block credited to a fresh
// address. initialSupply is 1_000_000_000_000_000 pesetas (1000 SPC).
func newTestChain(t *testing.T) (*chain.Blockchain, crypto.Address) {
	t.Helper()
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	addr := pub.ToAddress()
	bc, err := chain.NewBlockchain(addr, 1_000_000_000_000_000)
	if err != nil {
		t.Fatalf("NewBlockchain: %v", err)
	}
	return bc, addr
}

// newTestServer spins up an httptest.Server backed by a fresh blockchain and
// returns both the test server and the blockchain (for direct inspection).
func newTestServer(t *testing.T) (*httptest.Server, *chain.Blockchain, crypto.Address) {
	t.Helper()
	bc, addr := newTestChain(t)
	rpcSrv := NewServer(bc, "") // addr unused in tests
	ts := httptest.NewServer(rpcSrv.server.Handler)
	t.Cleanup(ts.Close)
	return ts, bc, addr
}

// get is a helper that issues a GET to url and returns the response.
func get(t *testing.T, url string) *http.Response {
	t.Helper()
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	return resp
}

// decodeJSON decodes the response body into v.
func decodeJSON(t *testing.T, resp *http.Response, v interface{}) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decode JSON: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestStatus(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp := get(t, ts.URL+"/status")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	decodeJSON(t, resp, &body)

	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %v", body["status"])
	}

	height, ok := body["height"].(float64)
	if !ok {
		t.Fatalf("height not a number: %v", body["height"])
	}
	if height < 0 {
		t.Errorf("height should be >= 0, got %v", height)
	}

	if _, ok := body["latest_hash"]; !ok {
		t.Error("missing latest_hash field")
	}
	if _, ok := body["total_supply"]; !ok {
		t.Error("missing total_supply field")
	}
	if _, ok := body["mempool_size"]; !ok {
		t.Error("missing mempool_size field")
	}
	if _, ok := body["peer_count"]; !ok {
		t.Error("missing peer_count field")
	}
}

func TestGetBlock_Genesis(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp := get(t, ts.URL+"/block/0")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	decodeJSON(t, resp, &body)

	if body["height"].(float64) != 0 {
		t.Errorf("expected height=0, got %v", body["height"])
	}
	if body["hash"] == "" {
		t.Error("hash should not be empty")
	}
	if body["tx_count"] == nil {
		t.Error("tx_count missing")
	}
}

func TestGetBlock_NotFound(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp := get(t, ts.URL+"/block/9999")
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestGetBlock_BadParam(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp := get(t, ts.URL+"/block/abc")
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestLatestBlock(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp := get(t, ts.URL+"/block/latest")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	decodeJSON(t, resp, &body)

	if _, ok := body["height"]; !ok {
		t.Error("height field missing")
	}
	if _, ok := body["hash"]; !ok {
		t.Error("hash field missing")
	}
}

func TestBalance_Known(t *testing.T) {
	ts, _, genesisAddr := newTestServer(t)

	url := fmt.Sprintf("%s/address/%s/balance", ts.URL, genesisAddr.String())
	resp := get(t, url)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	decodeJSON(t, resp, &body)

	balance, ok := body["balance"].(float64)
	if !ok {
		t.Fatalf("balance not a number: %v", body["balance"])
	}
	if balance <= 0 {
		t.Errorf("expected positive balance for genesis address, got %v", balance)
	}
}

func TestBalance_Unknown(t *testing.T) {
	ts, _, _ := newTestServer(t)

	// Use a hex address that is unlikely to exist.
	unknownAddr := "SPCffffffffffffffffffffffffffffffffffffffff"
	url := fmt.Sprintf("%s/address/%s/balance", ts.URL, unknownAddr)
	resp := get(t, url)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 (not 404) for unknown address, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	decodeJSON(t, resp, &body)

	if body["balance"].(float64) != 0 {
		t.Errorf("expected balance=0 for unknown address, got %v", body["balance"])
	}
	if body["nonce"].(float64) != 0 {
		t.Errorf("expected nonce=0 for unknown address, got %v", body["nonce"])
	}
}

func TestSendTx_Valid(t *testing.T) {
	ts, _, genesisAddr := newTestServer(t)

	// Build a minimal transfer. The mempool checks for duplicate/expiry/capacity
	// but not signature validity at this layer.
	_, toPub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	toAddr := toPub.ToAddress()

	payload := map[string]interface{}{
		"from":   genesisAddr.String(),
		"to":     toAddr.String(),
		"amount": 1_000_000_000_000,
		"nonce":  0,
		"fee":    1_000_000,
		"sig_r":  "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20",
		"sig_s":  "2122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f40",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(ts.URL+"/tx/send", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("POST /tx/send: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		var errBody map[string]interface{}
		decodeJSON(t, resp, &errBody)
		t.Fatalf("expected 200, got %d: %v", resp.StatusCode, errBody)
	}

	var result map[string]string
	decodeJSON(t, resp, &result)

	if result["tx_id"] == "" {
		t.Error("expected non-empty tx_id")
	}
	if result["status"] != "accepted" {
		t.Errorf("expected status=accepted, got %q", result["status"])
	}
}

func TestSendTx_BadJSON(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp, err := http.Post(ts.URL+"/tx/send", "application/json", bytes.NewReader([]byte("not json")))
	if err != nil {
		t.Fatalf("POST /tx/send: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestGetTx_NotFound(t *testing.T) {
	ts, _, _ := newTestServer(t)

	// 64 hex chars = 32 bytes = valid hash format but not present in mempool.
	unknownHash := "0000000000000000000000000000000000000000000000000000000000000000"
	resp := get(t, ts.URL+"/tx/"+unknownHash)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

// TestGetTx_Found adds a tx to the mempool and then retrieves it by hash.
func TestGetTx_Found(t *testing.T) {
	ts, bc, genesisAddr := newTestServer(t)

	_, toPub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	tx := block.NewTransaction(genesisAddr, toPub.ToAddress(), 1000, 0, 100)
	tx.Timestamp = time.Now().UnixNano() // ensure fresh timestamp
	tx.ID = tx.Hash()

	if err := bc.Mempool().Add(tx); err != nil {
		t.Fatalf("Mempool.Add: %v", err)
	}

	resp := get(t, ts.URL+"/tx/"+tx.ID.String())
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	decodeJSON(t, resp, &body)

	if body["id"] != tx.ID.String() {
		t.Errorf("expected id=%s, got %v", tx.ID.String(), body["id"])
	}
}

// TestValidators checks that GET /validators returns a valid JSON response.
func TestValidators(t *testing.T) {
	ts, _, _ := newTestServer(t)

	resp := get(t, ts.URL+"/validators")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	decodeJSON(t, resp, &body)

	if _, ok := body["count"]; !ok {
		t.Error("count field missing")
	}
	if _, ok := body["validators"]; !ok {
		t.Error("validators field missing")
	}
}
