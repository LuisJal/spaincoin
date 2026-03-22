package handlers

import (
	"net/http"

	"github.com/spaincoin/spaincoin/exchange/client"
)

// priceResponse is the response body for GET /api/market/price.
type priceResponse struct {
	Symbol    string  `json:"symbol"`
	PriceUSD  float64 `json:"price_usd"`
	PriceEUR  float64 `json:"price_eur"`
	Change24h float64 `json:"change_24h"`
	Volume24h float64 `json:"volume_24h"`
	MarketCap float64 `json:"market_cap"`
	Note      string  `json:"note"`
}

// statsResponse is the response body for GET /api/market/stats.
type statsResponse struct {
	Symbol            string  `json:"symbol"`
	CirculatingSupply float64 `json:"circulating_supply"`
	MaxSupply         float64 `json:"max_supply"`
	PriceUSD          float64 `json:"price_usd"`
	MarketCap         float64 `json:"market_cap"`
	MempoolSize       int     `json:"mempool_size"`
	PeerCount         int     `json:"peer_count"`
	BlockHeight       uint64  `json:"block_height"`
	Note              string  `json:"note"`
}

// HandlePrice handles GET /api/market/price.
// Returns hardcoded testnet price data.
func HandlePrice(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := priceResponse{
			Symbol:    "SPC",
			PriceUSD:  0.10,
			PriceEUR:  0.09,
			Change24h: 0.0,
			Volume24h: 0,
			MarketCap: 0,
			Note:      "testnet — precio de referencia",
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// HandleStats handles GET /api/market/stats.
// Calls the node /status and computes supply stats.
func HandleStats(nodeClient *client.NodeClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to reach node: "+err.Error())
			return
		}

		const maxSupply = 21_000_000.0
		circulatingSupply := float64(nodeStatus.TotalSupply) / 1_000_000_000_000_000.0
		const priceUSD = 0.10

		resp := statsResponse{
			Symbol:            "SPC",
			CirculatingSupply: circulatingSupply,
			MaxSupply:         maxSupply,
			PriceUSD:          priceUSD,
			MarketCap:         circulatingSupply * priceUSD,
			MempoolSize:       nodeStatus.MempoolSize,
			PeerCount:         nodeStatus.PeerCount,
			BlockHeight:       nodeStatus.Height,
			Note:              "testnet — datos de referencia",
		}
		writeJSON(w, http.StatusOK, resp)
	}
}
