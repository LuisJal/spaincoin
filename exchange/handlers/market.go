package handlers

import (
	"net/http"
	"strconv"

	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/market"
)

// priceResponse is the response body for GET /api/market/price.
type priceResponse struct {
	Symbol    string  `json:"symbol"`
	PriceUSD  float64 `json:"price_usd"`
	PriceEUR  float64 `json:"price_eur"`
	Change24h float64 `json:"change_24h"`
	Volume24h float64 `json:"volume_24h"`
	MarketCap float64 `json:"market_cap"`
	High24h   float64 `json:"high_24h"`
	Low24h    float64 `json:"low_24h"`
	Height    uint64  `json:"height"`
	Note      string  `json:"note"`
}

// statsResponse is the response body for GET /api/market/stats.
type statsResponse struct {
	Symbol            string  `json:"symbol"`
	CirculatingSupply float64 `json:"circulating_supply"`
	MaxSupply         float64 `json:"max_supply"`
	PriceEUR          float64 `json:"price_eur"`
	MarketCap         float64 `json:"market_cap"`
	MempoolSize       int     `json:"mempool_size"`
	PeerCount         int     `json:"peer_count"`
	BlockHeight       uint64  `json:"block_height"`
	Note              string  `json:"note"`
}

// HandlePrice handles GET /api/market/price.
// Uses the price simulator to return dynamic prices based on block height.
func HandlePrice(nodeClient *client.NodeClient, sim *market.Simulator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}

		pp := sim.CurrentPrice(nodeStatus.Height)
		change := sim.Change24h(nodeStatus.Height)
		circulatingSupply := float64(nodeStatus.TotalSupply) / 1_000_000_000_000_000.0

		resp := priceResponse{
			Symbol:    "SPC",
			PriceUSD:  pp.Price * 1.08, // EUR to USD approx
			PriceEUR:  pp.Price,
			Change24h: change,
			Volume24h: pp.Volume,
			MarketCap: circulatingSupply * pp.Price,
			High24h:   pp.High,
			Low24h:    pp.Low,
			Height:    nodeStatus.Height,
			Note:      "testnet — precio simulado",
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// HandleStats handles GET /api/market/stats.
func HandleStats(nodeClient *client.NodeClient, sim *market.Simulator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "failed to reach node: "+err.Error())
			return
		}

		const maxSupply = 21_000_000.0
		circulatingSupply := float64(nodeStatus.TotalSupply) / 1_000_000_000_000_000.0
		price := sim.PriceAtHeight(nodeStatus.Height)

		resp := statsResponse{
			Symbol:            "SPC",
			CirculatingSupply: circulatingSupply,
			MaxSupply:         maxSupply,
			PriceEUR:          price,
			MarketCap:         circulatingSupply * price,
			MempoolSize:       nodeStatus.MempoolSize,
			PeerCount:         nodeStatus.PeerCount,
			BlockHeight:       nodeStatus.Height,
			Note:              "testnet — datos simulados",
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

// HandlePriceHistory handles GET /api/market/history?points=100&range=24h.
// Returns an array of price points for chart rendering.
func HandlePriceHistory(nodeClient *client.NodeClient, sim *market.Simulator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}

		// Parse points parameter (default 100)
		points := 100
		if p := r.URL.Query().Get("points"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 && n <= 500 {
				points = n
			}
		}

		// Parse range parameter to determine step
		// 1h = 720 blocks (5s/block), 24h = 17280, 7d = 120960, 30d = 518400
		step := 1
		switch r.URL.Query().Get("range") {
		case "1h":
			step = 720 / points
		case "24h", "":
			step = 17280 / points
		case "7d":
			step = 120960 / points
		case "30d":
			step = 518400 / points
		}
		if step < 1 {
			step = 1
		}

		history := sim.PriceHistory(nodeStatus.Height, points, step)
		writeJSON(w, http.StatusOK, history)
	}
}

// HandleTicker handles GET /api/market/ticker.
// Returns a compact summary for the market overview.
func HandleTicker(nodeClient *client.NodeClient, sim *market.Simulator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}

		pp := sim.CurrentPrice(nodeStatus.Height)
		change := sim.Change24h(nodeStatus.Height)
		circulatingSupply := float64(nodeStatus.TotalSupply) / 1_000_000_000_000_000.0

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"symbol":     "SPC",
			"pair":       "SPC/EUR",
			"price":      pp.Price,
			"change_24h": change,
			"high_24h":   pp.High,
			"low_24h":    pp.Low,
			"volume_24h": pp.Volume,
			"market_cap": circulatingSupply * pp.Price,
			"supply":     circulatingSupply,
			"height":     nodeStatus.Height,
		})
	}
}

// marketTableEntry represents one row in the market overview table.
type marketTableEntry struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	Pair      string  `json:"pair"`
	Price     float64 `json:"price"`
	Change24h float64 `json:"change_24h"`
	Volume    float64 `json:"volume"`
	MarketCap float64 `json:"market_cap"`
	Supply    float64 `json:"supply"`
}

// HandleMarketTable handles GET /api/market/table.
// Returns a list of tokens for the market overview page.
// Currently only SPC, but structured for future expansion.
func HandleMarketTable(nodeClient *client.NodeClient, sim *market.Simulator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}

		pp := sim.CurrentPrice(nodeStatus.Height)
		change := sim.Change24h(nodeStatus.Height)
		circulatingSupply := float64(nodeStatus.TotalSupply) / 1_000_000_000_000_000.0

		table := []marketTableEntry{
			{
				Symbol:    "SPC",
				Name:      "SpainCoin",
				Pair:      "SPC/EUR",
				Price:     pp.Price,
				Change24h: change,
				Volume:    pp.Volume,
				MarketCap: circulatingSupply * pp.Price,
				Supply:    circulatingSupply,
			},
		}

		writeJSON(w, http.StatusOK, table)
	}
}
