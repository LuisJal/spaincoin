// Package market — Binance public API client for real crypto prices.
// No API key needed. Free tier, ~1200 req/min limit.
package market

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// BinancePrice holds a cached price from Binance.
type BinancePrice struct {
	Symbol    string
	PriceUSD  float64
	UpdatedAt time.Time
}

// PriceCache fetches and caches real prices from Binance.
type PriceCache struct {
	mu     sync.RWMutex
	prices map[string]*BinancePrice
	client *http.Client
}

// Symbols we fetch from Binance (traded against USDT).
var binanceSymbols = map[string]string{
	"BTC":   "BTCUSDT",
	"ETH":   "ETHUSDT",
	"BNB":   "BNBUSDT",
	"SOL":   "SOLUSDT",
	"XRP":   "XRPUSDT",
	"ADA":   "ADAUSDT",
	"DOGE":  "DOGEUSDT",
	"DOT":   "DOTUSDT",
	"AVAX":  "AVAXUSDT",
	"MATIC": "MATICUSDT",
}

// Approximate USD to EUR conversion rate.
const usdToEUR = 0.92

// NewPriceCache creates a new cache and starts background updates.
func NewPriceCache() *PriceCache {
	pc := &PriceCache{
		prices: make(map[string]*BinancePrice),
		client: &http.Client{Timeout: 10 * time.Second},
	}
	// Initial fetch
	pc.fetchAll()
	// Background refresh every 30 seconds
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			pc.fetchAll()
		}
	}()
	return pc
}

// GetPrice returns the EUR price for a symbol. Returns (price, true) if available.
func (pc *PriceCache) GetPrice(symbol string) (float64, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	p, ok := pc.prices[symbol]
	if !ok || time.Since(p.UpdatedAt) > 5*time.Minute {
		return 0, false
	}
	return p.PriceUSD * usdToEUR, true
}

// GetAllPrices returns all cached prices in EUR.
func (pc *PriceCache) GetAllPrices() map[string]float64 {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	result := make(map[string]float64)
	for symbol, p := range pc.prices {
		if time.Since(p.UpdatedAt) < 5*time.Minute {
			result[symbol] = p.PriceUSD * usdToEUR
		}
	}
	return result
}

// binanceTicker is the response from Binance ticker endpoint.
type binanceTicker struct {
	Symbol             string `json:"symbol"`
	LastPrice          string `json:"lastPrice"`
	PriceChangePercent string `json:"priceChangePercent"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	QuoteVolume        string `json:"quoteVolume"`
}

// fetchAll fetches prices for all symbols from Binance in a single API call.
func (pc *PriceCache) fetchAll() {
	url := "https://api.binance.com/api/v3/ticker/24hr"
	resp, err := pc.client.Get(url)
	if err != nil {
		log.Printf("[BINANCE] fetch error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("[BINANCE] HTTP %d", resp.StatusCode)
		return
	}

	var tickers []binanceTicker
	if err := json.NewDecoder(resp.Body).Decode(&tickers); err != nil {
		log.Printf("[BINANCE] decode error: %v", err)
		return
	}

	// Build lookup map: BTCUSDT -> ticker
	tickerMap := make(map[string]*binanceTicker, len(tickers))
	for i := range tickers {
		tickerMap[tickers[i].Symbol] = &tickers[i]
	}

	pc.mu.Lock()
	defer pc.mu.Unlock()

	now := time.Now()
	for symbol, binSymbol := range binanceSymbols {
		t, ok := tickerMap[binSymbol]
		if !ok {
			continue
		}
		price, err := strconv.ParseFloat(t.LastPrice, 64)
		if err != nil {
			continue
		}
		pc.prices[symbol] = &BinancePrice{
			Symbol:    symbol,
			PriceUSD:  price,
			UpdatedAt: now,
		}
	}

	log.Printf("[BINANCE] prices updated: %d symbols", len(pc.prices))
}

// GetTicker24h returns the 24h ticker data for a symbol from the cache.
func (pc *PriceCache) GetTicker24h(symbol string) (*Ticker24h, bool) {
	// We need to refetch for detailed data — use cached price + fetch individual
	binSymbol, ok := binanceSymbols[symbol]
	if !ok {
		return nil, false
	}

	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/24hr?symbol=%s", binSymbol)
	resp, err := pc.client.Get(url)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	var t binanceTicker
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return nil, false
	}

	price, _ := strconv.ParseFloat(t.LastPrice, 64)
	high, _ := strconv.ParseFloat(t.HighPrice, 64)
	low, _ := strconv.ParseFloat(t.LowPrice, 64)
	change, _ := strconv.ParseFloat(t.PriceChangePercent, 64)
	volume, _ := strconv.ParseFloat(t.QuoteVolume, 64)

	return &Ticker24h{
		Price:  price * usdToEUR,
		High:   high * usdToEUR,
		Low:    low * usdToEUR,
		Change: change,
		Volume: volume * usdToEUR,
	}, true
}

// Ticker24h holds 24h ticker data.
type Ticker24h struct {
	Price  float64
	High   float64
	Low    float64
	Change float64
	Volume float64
}
