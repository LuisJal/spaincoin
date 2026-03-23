package handlers

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/spaincoin/spaincoin/exchange/auth"
	"github.com/spaincoin/spaincoin/exchange/client"
	"github.com/spaincoin/spaincoin/exchange/database"
	"github.com/spaincoin/spaincoin/exchange/market"
)

// buyRequest is the request body for POST /api/trade/buy.
type buyRequest struct {
	Symbol    string  `json:"symbol"`
	Amount    float64 `json:"amount"`
	AmountSPC float64 `json:"amount_spc"` // backward compat
	AmountEUR float64 `json:"amount_eur"`
}

// sellRequest is the request body for POST /api/trade/sell.
type sellRequest struct {
	Symbol    string  `json:"symbol"`
	Amount    float64 `json:"amount"`
	AmountSPC float64 `json:"amount_spc"` // backward compat
}

// tradeResponse is the response for buy/sell operations.
type tradeResponse struct {
	TradeID   string  `json:"trade_id"`
	Type      string  `json:"type"`
	Symbol    string  `json:"symbol"`
	Amount    float64 `json:"amount"`
	AmountSPC float64 `json:"amount_spc"` // backward compat
	PriceEUR  float64 `json:"price_eur"`
	TotalEUR  float64 `json:"total_eur"`
	Status    string  `json:"status"`
}

// depositRequest is the request body for POST /api/trade/deposit-eur.
type depositRequest struct {
	Amount float64 `json:"amount"`
}

// holdingEntry is one row in the portfolio.
type holdingEntry struct {
	Symbol   string  `json:"symbol"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	ValueEUR float64 `json:"value_eur"`
}

// portfolioResponse is the response for GET /api/trade/portfolio.
type portfolioResponse struct {
	EUR           float64        `json:"eur"`
	Holdings      []holdingEntry `json:"holdings"`
	TotalValueEUR float64        `json:"total_value_eur"`
}

// balanceResponse is the response for GET /api/trade/balance.
type balanceResponse struct {
	EUR float64 `json:"eur"`
	SPC float64 `json:"spc"`
}

// HandleBuy handles POST /api/trade/buy (auth required). Supports any symbol.
func HandleBuy(nodeClient *client.NodeClient, sim *market.Simulator, userDB *database.UserDB, tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		claims := auth.GetClaims(r)
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req buyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		// Backward compat: default to SPC
		symbol := req.Symbol
		if symbol == "" {
			symbol = "SPC"
		}
		amount := req.Amount
		if amount == 0 {
			amount = req.AmountSPC
		}

		// Get current price
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}
		price, ok := GetSimulatedPrice(symbol, sim, nodeStatus.Height)
		if !ok {
			writeError(w, http.StatusBadRequest, "símbolo no soportado: "+symbol)
			return
		}

		// Calculate amounts
		var totalEUR float64
		if amount > 0 {
			totalEUR = amount * price
		} else if req.AmountEUR > 0 {
			totalEUR = req.AmountEUR
			amount = totalEUR / price
		} else {
			writeError(w, http.StatusBadRequest, "specify amount or amount_eur")
			return
		}

		if amount <= 0 || amount > 1_000_000 {
			writeError(w, http.StatusBadRequest, "invalid amount")
			return
		}

		amount = math.Round(amount*1_000_000) / 1_000_000
		totalEUR = math.Round(amount*price*100) / 100

		// Check EUR balance
		eurBalance, err := tradeDB.GetEURBalance(claims.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "balance check failed")
			return
		}
		if eurBalance < totalEUR {
			writeError(w, http.StatusBadRequest, "saldo EUR insuficiente")
			return
		}

		// Debit EUR
		if err := tradeDB.DebitEUR(claims.UserID, totalEUR); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Credit crypto
		if err := tradeDB.CreditCrypto(claims.UserID, symbol, amount); err != nil {
			tradeDB.CreditEUR(claims.UserID, totalEUR)
			writeError(w, http.StatusInternalServerError, "credit failed")
			return
		}

		// Record trade
		trade := &database.TradeRecord{
			ID:        database.GenerateTradeID(claims.UserID),
			UserID:    claims.UserID,
			Type:      "buy",
			Pair:      symbol + "/EUR",
			Symbol:    symbol,
			Amount:    amount,
			AmountSPC: amount, // backward compat
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
			CreatedAt: time.Now().UnixNano(),
		}
		if err := tradeDB.SaveTrade(trade); err != nil {
			tradeDB.CreditEUR(claims.UserID, totalEUR)
			tradeDB.DebitCrypto(claims.UserID, symbol, amount)
			writeError(w, http.StatusInternalServerError, "trade failed")
			return
		}

		log.Printf("[TRADE] BUY user=%s symbol=%s amount=%.6f price=%.2f total=%.2f EUR",
			claims.UserID, symbol, amount, price, totalEUR)

		writeJSON(w, http.StatusOK, tradeResponse{
			TradeID:   trade.ID,
			Type:      "buy",
			Symbol:    symbol,
			Amount:    amount,
			AmountSPC: amount,
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
		})
	}
}

// HandleSell handles POST /api/trade/sell (auth required). Supports any symbol.
func HandleSell(nodeClient *client.NodeClient, sim *market.Simulator, userDB *database.UserDB, tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		claims := auth.GetClaims(r)
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req sellRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		symbol := req.Symbol
		if symbol == "" {
			symbol = "SPC"
		}
		amount := req.Amount
		if amount == 0 {
			amount = req.AmountSPC
		}

		if amount <= 0 || amount > 1_000_000 {
			writeError(w, http.StatusBadRequest, "invalid amount")
			return
		}

		// Get current price
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}
		price, ok := GetSimulatedPrice(symbol, sim, nodeStatus.Height)
		if !ok {
			writeError(w, http.StatusBadRequest, "símbolo no soportado: "+symbol)
			return
		}

		amount = math.Round(amount*1_000_000) / 1_000_000
		totalEUR := math.Round(amount*price*100) / 100

		// Check crypto balance
		cryptoBalance, err := tradeDB.GetCryptoBalance(claims.UserID, symbol)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "balance check failed")
			return
		}
		if cryptoBalance < amount {
			writeError(w, http.StatusBadRequest, "saldo "+symbol+" insuficiente")
			return
		}

		// Debit crypto
		if err := tradeDB.DebitCrypto(claims.UserID, symbol, amount); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Credit EUR
		if err := tradeDB.CreditEUR(claims.UserID, totalEUR); err != nil {
			tradeDB.CreditCrypto(claims.UserID, symbol, amount)
			writeError(w, http.StatusInternalServerError, "credit failed")
			return
		}

		// Record trade
		trade := &database.TradeRecord{
			ID:        database.GenerateTradeID(claims.UserID),
			UserID:    claims.UserID,
			Type:      "sell",
			Pair:      symbol + "/EUR",
			Symbol:    symbol,
			Amount:    amount,
			AmountSPC: amount,
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
			CreatedAt: time.Now().UnixNano(),
		}
		if err := tradeDB.SaveTrade(trade); err != nil {
			tradeDB.CreditCrypto(claims.UserID, symbol, amount)
			tradeDB.DebitEUR(claims.UserID, totalEUR)
			writeError(w, http.StatusInternalServerError, "trade failed")
			return
		}

		log.Printf("[TRADE] SELL user=%s symbol=%s amount=%.6f price=%.2f total=%.2f EUR",
			claims.UserID, symbol, amount, price, totalEUR)

		writeJSON(w, http.StatusOK, tradeResponse{
			TradeID:   trade.ID,
			Type:      "sell",
			Symbol:    symbol,
			Amount:    amount,
			AmountSPC: amount,
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
		})
	}
}

// HandleDepositEUR handles POST /api/trade/deposit-eur (auth required).
func HandleDepositEUR(tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		claims := auth.GetClaims(r)
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req depositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Amount <= 0 || req.Amount > 10000 {
			writeError(w, http.StatusBadRequest, "amount must be between 0 and 10,000 EUR")
			return
		}

		if err := tradeDB.CreditEUR(claims.UserID, req.Amount); err != nil {
			writeError(w, http.StatusInternalServerError, "deposit failed")
			return
		}

		newBalance, _ := tradeDB.GetEURBalance(claims.UserID)

		log.Printf("[TRADE] DEPOSIT user=%s amount=%.2f EUR new_balance=%.2f", claims.UserID, req.Amount, newBalance)

		writeJSON(w, http.StatusOK, map[string]interface{}{
			"deposited":   req.Amount,
			"new_balance": math.Round(newBalance*100) / 100,
		})
	}
}

// HandlePortfolio handles GET /api/trade/portfolio (auth required).
func HandlePortfolio(nodeClient *client.NodeClient, sim *market.Simulator, userDB *database.UserDB, tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		claims := auth.GetClaims(r)
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}

		eurBalance, _ := tradeDB.GetEURBalance(claims.UserID)
		cryptoBalances, _ := tradeDB.GetAllCryptoBalances(claims.UserID)

		var holdings []holdingEntry
		var totalValue float64

		// Build holdings from all crypto balances
		for symbol, amount := range cryptoBalances {
			if amount <= 0 {
				continue
			}
			price, ok := GetSimulatedPrice(symbol, sim, nodeStatus.Height)
			if !ok {
				continue
			}
			value := amount * price
			name := symbolName(symbol)
			holdings = append(holdings, holdingEntry{
				Symbol:   symbol,
				Name:     name,
				Amount:   amount,
				Price:    price,
				ValueEUR: math.Round(value*100) / 100,
			})
			totalValue += value
		}

		totalValue += eurBalance

		writeJSON(w, http.StatusOK, portfolioResponse{
			EUR:           math.Round(eurBalance*100) / 100,
			Holdings:      holdings,
			TotalValueEUR: math.Round(totalValue*100) / 100,
		})
	}
}

// HandleTradeHistory handles GET /api/trade/history (auth required).
func HandleTradeHistory(tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		claims := auth.GetClaims(r)
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		trades, err := tradeDB.GetUserTrades(claims.UserID, 50)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to load trades")
			return
		}

		if trades == nil {
			trades = []*database.TradeRecord{}
		}

		writeJSON(w, http.StatusOK, trades)
	}
}

// HandleTradeBalance handles GET /api/trade/balance (auth required).
func HandleTradeBalance(nodeClient *client.NodeClient, userDB *database.UserDB, tradeDB *database.TradeDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		claims := auth.GetClaims(r)
		if claims == nil {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		eurBalance, _ := tradeDB.GetEURBalance(claims.UserID)

		// SPC from exchange holdings
		spcBalance, _ := tradeDB.GetCryptoBalance(claims.UserID, "SPC")

		writeJSON(w, http.StatusOK, balanceResponse{
			EUR: math.Round(eurBalance*100) / 100,
			SPC: math.Round(spcBalance*1_000_000) / 1_000_000,
		})
	}
}

// symbolName returns the human-readable name for a symbol.
func symbolName(symbol string) string {
	names := map[string]string{
		"SPC": "SpainCoin", "BTC": "Bitcoin", "ETH": "Ethereum",
		"BNB": "BNB", "SOL": "Solana", "XRP": "XRP",
		"ADA": "Cardano", "DOGE": "Dogecoin", "DOT": "Polkadot",
		"AVAX": "Avalanche", "MATIC": "Polygon",
	}
	if n, ok := names[symbol]; ok {
		return n
	}
	return symbol
}
