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
	AmountSPC float64 `json:"amount_spc"`
	AmountEUR float64 `json:"amount_eur"`
}

// sellRequest is the request body for POST /api/trade/sell.
type sellRequest struct {
	AmountSPC float64 `json:"amount_spc"`
}

// tradeResponse is the response for buy/sell operations.
type tradeResponse struct {
	TradeID   string  `json:"trade_id"`
	Type      string  `json:"type"`
	AmountSPC float64 `json:"amount_spc"`
	PriceEUR  float64 `json:"price_eur"`
	TotalEUR  float64 `json:"total_eur"`
	Status    string  `json:"status"`
}

// balanceResponse is the response for GET /api/trade/balance.
type balanceResponse struct {
	EUR float64 `json:"eur"`
	SPC float64 `json:"spc"`
}

// HandleBuy handles POST /api/trade/buy (auth required).
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

		// Get current price
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}
		price := sim.PriceAtHeight(nodeStatus.Height)

		// Calculate amounts
		var amountSPC, totalEUR float64
		if req.AmountSPC > 0 {
			amountSPC = req.AmountSPC
			totalEUR = amountSPC * price
		} else if req.AmountEUR > 0 {
			totalEUR = req.AmountEUR
			amountSPC = totalEUR / price
		} else {
			writeError(w, http.StatusBadRequest, "specify amount_spc or amount_eur")
			return
		}

		// Validate amounts
		if amountSPC <= 0 || amountSPC > 1_000_000 {
			writeError(w, http.StatusBadRequest, "invalid amount")
			return
		}

		// Round to 6 decimals
		amountSPC = math.Round(amountSPC*1_000_000) / 1_000_000
		totalEUR = math.Round(totalEUR*100) / 100

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

		// Record trade
		trade := &database.TradeRecord{
			ID:        database.GenerateTradeID(claims.UserID),
			UserID:    claims.UserID,
			Type:      "buy",
			Pair:      "SPC/EUR",
			AmountSPC: amountSPC,
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
			CreatedAt: time.Now().UnixNano(),
		}
		if err := tradeDB.SaveTrade(trade); err != nil {
			// Refund EUR on failure
			tradeDB.CreditEUR(claims.UserID, totalEUR)
			writeError(w, http.StatusInternalServerError, "trade failed")
			return
		}

		log.Printf("[TRADE] BUY user=%s amount=%.6f SPC price=%.6f EUR total=%.2f EUR",
			claims.UserID, amountSPC, price, totalEUR)

		writeJSON(w, http.StatusOK, tradeResponse{
			TradeID:   trade.ID,
			Type:      "buy",
			AmountSPC: amountSPC,
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
		})
	}
}

// HandleSell handles POST /api/trade/sell (auth required).
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

		if req.AmountSPC <= 0 || req.AmountSPC > 1_000_000 {
			writeError(w, http.StatusBadRequest, "invalid amount")
			return
		}

		// Get current price
		nodeStatus, err := nodeClient.Status()
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot reach node")
			return
		}
		price := sim.PriceAtHeight(nodeStatus.Height)

		amountSPC := math.Round(req.AmountSPC*1_000_000) / 1_000_000
		totalEUR := math.Round(amountSPC*price*100) / 100

		// Check SPC balance on-chain
		user, err := userDB.GetByID(claims.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "user not found")
			return
		}
		balanceInfo, err := nodeClient.GetBalance(user.WalletAddress)
		if err != nil {
			writeError(w, http.StatusBadGateway, "cannot check SPC balance")
			return
		}
		balanceSPC := float64(balanceInfo.Balance) / 1_000_000_000_000_000.0
		if balanceSPC < amountSPC {
			writeError(w, http.StatusBadRequest, "saldo SPC insuficiente")
			return
		}

		// Credit EUR
		if err := tradeDB.CreditEUR(claims.UserID, totalEUR); err != nil {
			writeError(w, http.StatusInternalServerError, "credit failed")
			return
		}

		// Record trade
		trade := &database.TradeRecord{
			ID:        database.GenerateTradeID(claims.UserID),
			UserID:    claims.UserID,
			Type:      "sell",
			Pair:      "SPC/EUR",
			AmountSPC: amountSPC,
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
			CreatedAt: time.Now().UnixNano(),
		}
		if err := tradeDB.SaveTrade(trade); err != nil {
			tradeDB.DebitEUR(claims.UserID, totalEUR)
			writeError(w, http.StatusInternalServerError, "trade failed")
			return
		}

		log.Printf("[TRADE] SELL user=%s amount=%.6f SPC price=%.6f EUR total=%.2f EUR",
			claims.UserID, amountSPC, price, totalEUR)

		writeJSON(w, http.StatusOK, tradeResponse{
			TradeID:   trade.ID,
			Type:      "sell",
			AmountSPC: amountSPC,
			PriceEUR:  price,
			TotalEUR:  totalEUR,
			Status:    "completed",
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

		eurBalance, err := tradeDB.GetEURBalance(claims.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "balance check failed")
			return
		}

		// Get SPC balance from node
		user, err := userDB.GetByID(claims.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "user not found")
			return
		}
		var spcBalance float64
		balanceInfo, err := nodeClient.GetBalance(user.WalletAddress)
		if err == nil {
			spcBalance = float64(balanceInfo.Balance) / 1_000_000_000_000_000.0
		}

		writeJSON(w, http.StatusOK, balanceResponse{
			EUR: math.Round(eurBalance*100) / 100,
			SPC: math.Round(spcBalance*1_000_000) / 1_000_000,
		})
	}
}
