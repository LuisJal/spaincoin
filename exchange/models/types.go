package models

// NodeStatus represents the response from GET /status on the SpainCoin node.
type NodeStatus struct {
	Status      string `json:"status"`
	Height      uint64 `json:"height"`
	LatestHash  string `json:"latest_hash"`
	TotalSupply uint64 `json:"total_supply"`
	MempoolSize int    `json:"mempool_size"`
	PeerCount   int    `json:"peer_count"`
}

// BlockInfo represents the response from GET /block/{height} or GET /block/latest.
type BlockInfo struct {
	Height       uint64   `json:"height"`
	Hash         string   `json:"hash"`
	PrevHash     string   `json:"prev_hash"`
	Timestamp    int64    `json:"timestamp"`
	Validator    string   `json:"validator"`
	TxCount      int      `json:"tx_count"`
	Transactions []TxInfo `json:"transactions"`
}

// TxInfo represents a single transaction embedded in a block or returned by GET /tx/{hash}.
type TxInfo struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    uint64 `json:"amount"`
	Fee       uint64 `json:"fee"`
	Nonce     uint64 `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
}

// BalanceInfo represents the response from GET /address/{address}/balance.
// BalanceSPC is a human-readable value: balance / 1e15 (3 decimal places).
type BalanceInfo struct {
	Address    string  `json:"address"`
	Balance    uint64  `json:"balance"`
	Nonce      uint64  `json:"nonce"`
	BalanceSPC float64 `json:"balance_spc"`
}

// SendTxRequest is the payload for POST /tx/send on the SpainCoin node.
type SendTxRequest struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint64 `json:"amount"`
	Nonce  uint64 `json:"nonce"`
	Fee    uint64 `json:"fee"`
	SigR   string `json:"sig_r"`
	SigS   string `json:"sig_s"`
}

// SendTxResponse is the response from POST /tx/send on the SpainCoin node.
type SendTxResponse struct {
	TxID   string `json:"tx_id"`
	Status string `json:"status"`
}

// ExchangeStatus is returned by the exchange API GET /api/status.
type ExchangeStatus struct {
	Exchange         string      `json:"exchange"`
	Version          string      `json:"version"`
	Node             *NodeStatus `json:"node"`
	BlockTimeSeconds int         `json:"block_time_seconds"`
}
