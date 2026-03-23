package handlers

import (
	"math"
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

// externalCoin defines a reference crypto for the market overview.
type externalCoin struct {
	Symbol    string
	Name      string
	BasePrice float64
	Supply    float64
	WaveDiv   float64 // unique wave divisor per coin
}

var referenceCryptos = []externalCoin{
	{"BTC", "Bitcoin", 82000.0, 19_800_000, 97},
	{"ETH", "Ethereum", 1850.0, 120_500_000, 73},
	{"BNB", "BNB", 610.0, 145_900_000, 61},
	{"SOL", "Solana", 130.0, 440_000_000, 53},
	{"XRP", "XRP", 2.15, 55_000_000_000, 43},
	{"ADA", "Cardano", 0.70, 35_000_000_000, 67},
	{"DOGE", "Dogecoin", 0.16, 144_000_000_000, 31},
	{"DOT", "Polkadot", 4.20, 1_400_000_000, 59},
	{"AVAX", "Avalanche", 20.0, 400_000_000, 47},
	{"MATIC", "Polygon", 0.22, 10_000_000_000, 37},
}

// HandleMarketTable handles GET /api/market/table.
// Returns SPC plus reference cryptos for a full market overview.
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

		// SPC first
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

		// Reference cryptos with simulated variation
		h := float64(nodeStatus.Height)
		for _, c := range referenceCryptos {
			// Each coin has unique wave pattern seeded by its WaveDiv
			wave1 := math.Sin(h/c.WaveDiv) * 0.02
			wave2 := math.Sin(h/(c.WaveDiv*2.7)) * 0.015
			price := c.BasePrice * (1 + wave1 + wave2)
			price = math.Round(price*100) / 100

			// 24h change from block height offset
			prevH := h - 17280
			if prevH < 0 {
				prevH = 0
			}
			prevWave1 := math.Sin(prevH/c.WaveDiv) * 0.02
			prevWave2 := math.Sin(prevH/(c.WaveDiv*2.7)) * 0.015
			prevPrice := c.BasePrice * (1 + prevWave1 + prevWave2)
			ch := ((price - prevPrice) / prevPrice) * 100
			ch = math.Round(ch*100) / 100

			vol := c.BasePrice * c.Supply * 0.001 * (1 + math.Sin(h/50)*0.3)
			vol = math.Round(vol)

			table = append(table, marketTableEntry{
				Symbol:    c.Symbol,
				Name:      c.Name,
				Pair:      c.Symbol + "/EUR",
				Price:     price,
				Change24h: ch,
				Volume:    vol,
				MarketCap: price * c.Supply,
				Supply:    c.Supply,
			})
		}

		writeJSON(w, http.StatusOK, table)
	}
}
