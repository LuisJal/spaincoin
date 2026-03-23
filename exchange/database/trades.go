package database

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

var (
	bucketTrades      = []byte("trades")
	bucketUserTrades  = []byte("user_trades")
	bucketEURBalances = []byte("eur_balances")
)

// TradeRecord represents a completed trade (buy or sell).
type TradeRecord struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Type      string  `json:"type"` // "buy" or "sell"
	Pair      string  `json:"pair"` // "SPC/EUR"
	AmountSPC float64 `json:"amount_spc"`
	PriceEUR  float64 `json:"price_eur"`
	TotalEUR  float64 `json:"total_eur"`
	Status    string  `json:"status"` // "completed"
	CreatedAt int64   `json:"created_at"`
}

// TradeDB handles trade and EUR balance storage.
type TradeDB struct {
	db *bbolt.DB
}

// NewTradeDB initializes the trades and balances buckets.
func NewTradeDB(db *bbolt.DB) (*TradeDB, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketTrades); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketUserTrades); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketEURBalances); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &TradeDB{db: db}, nil
}

// SaveTrade stores a trade record.
func (t *TradeDB) SaveTrade(trade *TradeRecord) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		data, err := json.Marshal(trade)
		if err != nil {
			return err
		}

		// Save in trades bucket
		b := tx.Bucket(bucketTrades)
		if err := b.Put([]byte(trade.ID), data); err != nil {
			return err
		}

		// Save in user index: key = userID + timestamp for ordering
		idx := tx.Bucket(bucketUserTrades)
		key := fmt.Sprintf("%s:%020d:%s", trade.UserID, trade.CreatedAt, trade.ID)
		return idx.Put([]byte(key), []byte(trade.ID))
	})
}

// GetUserTrades returns the most recent trades for a user, ordered by most recent first.
func (t *TradeDB) GetUserTrades(userID string, limit int) ([]*TradeRecord, error) {
	var trades []*TradeRecord

	err := t.db.View(func(tx *bbolt.Tx) error {
		idx := tx.Bucket(bucketUserTrades)
		tradesBucket := tx.Bucket(bucketTrades)
		c := idx.Cursor()

		prefix := []byte(userID + ":")
		// Seek to end of prefix range and iterate backwards
		var keys [][]byte
		for k, _ := c.Seek(prefix); k != nil && len(k) >= len(prefix) && string(k[:len(prefix)]) == string(prefix); k, _ = c.Next() {
			keyCopy := make([]byte, len(k))
			copy(keyCopy, k)
			keys = append(keys, keyCopy)
		}

		// Reverse to get most recent first
		start := len(keys) - limit
		if start < 0 {
			start = 0
		}
		for i := len(keys) - 1; i >= start; i-- {
			tradeID := idx.Get(keys[i])
			if tradeID == nil {
				continue
			}
			data := tradesBucket.Get(tradeID)
			if data == nil {
				continue
			}
			var trade TradeRecord
			if err := json.Unmarshal(data, &trade); err != nil {
				continue
			}
			trades = append(trades, &trade)
		}

		return nil
	})

	return trades, err
}

// GetEURBalance returns the EUR balance for a user.
func (t *TradeDB) GetEURBalance(userID string) (float64, error) {
	var balance float64
	err := t.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketEURBalances)
		data := b.Get([]byte(userID))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &balance)
	})
	return balance, err
}

// SetEURBalance sets the EUR balance for a user.
func (t *TradeDB) SetEURBalance(userID string, balance float64) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketEURBalances)
		data, err := json.Marshal(balance)
		if err != nil {
			return err
		}
		return b.Put([]byte(userID), data)
	})
}

// CreditEUR adds EUR to a user's balance.
func (t *TradeDB) CreditEUR(userID string, amount float64) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketEURBalances)
		var current float64
		data := b.Get([]byte(userID))
		if data != nil {
			json.Unmarshal(data, &current)
		}
		current += amount
		newData, err := json.Marshal(current)
		if err != nil {
			return err
		}
		return b.Put([]byte(userID), newData)
	})
}

// DebitEUR subtracts EUR from a user's balance. Returns error if insufficient funds.
func (t *TradeDB) DebitEUR(userID string, amount float64) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketEURBalances)
		var current float64
		data := b.Get([]byte(userID))
		if data != nil {
			json.Unmarshal(data, &current)
		}
		if current < amount {
			return fmt.Errorf("insufficient EUR balance: have %.2f, need %.2f", current, amount)
		}
		current -= amount
		newData, err := json.Marshal(current)
		if err != nil {
			return err
		}
		return b.Put([]byte(userID), newData)
	})
}

// InitUserBalance sets initial EUR balance for a new user (testnet money).
func (t *TradeDB) InitUserBalance(userID string, initialEUR float64) error {
	return t.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketEURBalances)
		existing := b.Get([]byte(userID))
		if existing != nil {
			return nil // already initialized
		}
		data, err := json.Marshal(initialEUR)
		if err != nil {
			return err
		}
		return b.Put([]byte(userID), data)
	})
}

// GenerateTradeID creates a unique trade ID from timestamp and user.
func GenerateTradeID(userID string) string {
	return fmt.Sprintf("trade_%d_%s", time.Now().UnixNano(), userID[:8])
}
