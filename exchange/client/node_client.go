package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spaincoin/spaincoin/exchange/models"
)

// NodeClient is an HTTP client for the SpainCoin node RPC API.
type NodeClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewNodeClient creates a NodeClient targeting nodeURL (e.g. "http://204.168.176.40:8545")
// with a 10-second timeout.
func NewNodeClient(nodeURL string) *NodeClient {
	return &NodeClient{
		baseURL: nodeURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// get performs a GET request and JSON-decodes the response body into dest.
func (c *NodeClient) get(path string, dest interface{}) error {
	url := c.baseURL + path
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: unexpected status %d", path, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("GET %s: decode: %w", path, err)
	}
	return nil
}

// post performs a POST request with a JSON body and JSON-decodes the response into dest.
func (c *NodeClient) post(path string, body interface{}, dest interface{}) error {
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("POST %s: marshal: %w", path, err)
	}

	url := c.baseURL + path
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("POST %s: %w", path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("POST %s: unexpected status %d", path, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(dest); err != nil {
		return fmt.Errorf("POST %s: decode: %w", path, err)
	}
	return nil
}

// Status calls GET /status on the node.
func (c *NodeClient) Status() (*models.NodeStatus, error) {
	var s models.NodeStatus
	if err := c.get("/status", &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetBlock calls GET /block/{height} on the node.
func (c *NodeClient) GetBlock(height uint64) (*models.BlockInfo, error) {
	var b models.BlockInfo
	if err := c.get(fmt.Sprintf("/block/%d", height), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// GetLatestBlock calls GET /block/latest on the node.
func (c *NodeClient) GetLatestBlock() (*models.BlockInfo, error) {
	var b models.BlockInfo
	if err := c.get("/block/latest", &b); err != nil {
		return nil, err
	}
	return &b, nil
}

// GetBalance calls GET /address/{address}/balance on the node and populates BalanceSPC.
func (c *NodeClient) GetBalance(address string) (*models.BalanceInfo, error) {
	var info models.BalanceInfo
	if err := c.get(fmt.Sprintf("/address/%s/balance", address), &info); err != nil {
		return nil, err
	}
	info.BalanceSPC = float64(info.Balance) / 1_000_000_000_000.0
	return &info, nil
}

// SendTx calls POST /tx/send on the node.
func (c *NodeClient) SendTx(req *models.SendTxRequest) (*models.SendTxResponse, error) {
	var resp models.SendTxResponse
	if err := c.post("/tx/send", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetRecentBlocks fetches the last count blocks from the chain.
// It first retrieves the latest block height, then fetches descending blocks.
// Handles the case where the chain height is less than count.
func (c *NodeClient) GetRecentBlocks(count int) ([]*models.BlockInfo, error) {
	latest, err := c.GetLatestBlock()
	if err != nil {
		return nil, fmt.Errorf("GetRecentBlocks: get latest: %w", err)
	}

	height := int(latest.Height)
	if count > height+1 {
		// height is 0-indexed: blocks 0..height exist
		count = height + 1
	}

	blocks := make([]*models.BlockInfo, 0, count)
	// Include the latest block we already fetched
	blocks = append(blocks, latest)

	for i := 1; i < count; i++ {
		h := uint64(height - i)
		b, err := c.GetBlock(h)
		if err != nil {
			return nil, fmt.Errorf("GetRecentBlocks: get block %d: %w", h, err)
		}
		blocks = append(blocks, b)
	}

	return blocks, nil
}
