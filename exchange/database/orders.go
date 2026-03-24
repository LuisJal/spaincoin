package database

import (
	"encoding/json"
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

var (
	bucketOrders  = []byte("orders")
	bucketAdminKV = []byte("admin_kv")
	bucketWallets = []byte("wallets")
)

// OrderStatus represents the state of a buy/sell order.
type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderConfirmed OrderStatus = "confirmed"
	OrderCancelled OrderStatus = "cancelled"
)

// Order represents a P2P buy or sell request.
type Order struct {
	ID          int64       `json:"id"`
	Type        string      `json:"type"`         // "buy" or "sell"
	UserName    string      `json:"user_name"`    // Telegram username or name
	UserChatID  int64       `json:"user_chat_id"` // Telegram chat ID
	WalletAddr  string      `json:"wallet_addr"`  // SPC address
	AmountEUR   float64     `json:"amount_eur"`   // EUR amount
	AmountSPC   float64     `json:"amount_spc"`   // SPC amount (calculated at current price)
	PriceEUR    float64     `json:"price_eur"`    // Price at time of order
	Status      OrderStatus `json:"status"`
	CreatedAt   int64       `json:"created_at"`
	ConfirmedAt int64       `json:"confirmed_at,omitempty"`
}

// OrderDB handles orders and admin key-value storage.
type OrderDB struct {
	db *bbolt.DB
}

// NewOrderDB initializes the orders and admin buckets.
func NewOrderDB(db *bbolt.DB) (*OrderDB, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(bucketOrders); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketAdminKV); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(bucketWallets); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &OrderDB{db: db}, nil
}

// SetPrice stores the admin-set SPC price in EUR.
func (o *OrderDB) SetPrice(priceEUR float64) error {
	return o.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAdminKV)
		data, _ := json.Marshal(priceEUR)
		return b.Put([]byte("spc_price_eur"), data)
	})
}

// GetPrice returns the admin-set SPC price in EUR.
func (o *OrderDB) GetPrice() (float64, error) {
	var price float64
	err := o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAdminKV)
		data := b.Get([]byte("spc_price_eur"))
		if data == nil {
			price = 0.10 // default
			return nil
		}
		return json.Unmarshal(data, &price)
	})
	return price, err
}

// SetAdminValue stores a key-value pair in admin storage.
func (o *OrderDB) SetAdminValue(key string, value string) error {
	return o.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAdminKV)
		return b.Put([]byte(key), []byte(value))
	})
}

// GetAdminValue retrieves a value from admin storage.
func (o *OrderDB) GetAdminValue(key string) (string, error) {
	var val string
	err := o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketAdminKV)
		data := b.Get([]byte(key))
		if data != nil {
			val = string(data)
		}
		return nil
	})
	return val, err
}

// CreateOrder stores a new order and returns its ID.
func (o *OrderDB) CreateOrder(order *Order) error {
	return o.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketOrders)
		id, _ := b.NextSequence()
		order.ID = int64(id)
		order.CreatedAt = time.Now().Unix()
		data, err := json.Marshal(order)
		if err != nil {
			return err
		}
		return b.Put([]byte(fmt.Sprintf("%010d", order.ID)), data)
	})
}

// GetOrder retrieves an order by ID.
func (o *OrderDB) GetOrder(id int64) (*Order, error) {
	var order Order
	err := o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketOrders)
		data := b.Get([]byte(fmt.Sprintf("%010d", id)))
		if data == nil {
			return fmt.Errorf("order %d not found", id)
		}
		return json.Unmarshal(data, &order)
	})
	return &order, err
}

// ConfirmOrder marks an order as confirmed.
func (o *OrderDB) ConfirmOrder(id int64) error {
	return o.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketOrders)
		key := []byte(fmt.Sprintf("%010d", id))
		data := b.Get(key)
		if data == nil {
			return fmt.Errorf("order %d not found", id)
		}
		var order Order
		if err := json.Unmarshal(data, &order); err != nil {
			return err
		}
		order.Status = OrderConfirmed
		order.ConfirmedAt = time.Now().Unix()
		newData, err := json.Marshal(order)
		if err != nil {
			return err
		}
		return b.Put(key, newData)
	})
}

// CancelOrder marks an order as cancelled.
func (o *OrderDB) CancelOrder(id int64) error {
	return o.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketOrders)
		key := []byte(fmt.Sprintf("%010d", id))
		data := b.Get(key)
		if data == nil {
			return fmt.Errorf("order %d not found", id)
		}
		var order Order
		if err := json.Unmarshal(data, &order); err != nil {
			return err
		}
		order.Status = OrderCancelled
		newData, err := json.Marshal(order)
		if err != nil {
			return err
		}
		return b.Put(key, newData)
	})
}

// GetPendingOrders returns all pending orders.
func (o *OrderDB) GetPendingOrders() ([]*Order, error) {
	var orders []*Order
	err := o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketOrders)
		return b.ForEach(func(k, v []byte) error {
			var order Order
			if err := json.Unmarshal(v, &order); err != nil {
				return nil
			}
			if order.Status == OrderPending {
				orders = append(orders, &order)
			}
			return nil
		})
	})
	return orders, err
}

// GetStats returns total SPC sold and total EUR received.
func (o *OrderDB) GetStats() (totalSPC float64, totalEUR float64, count int, err error) {
	err = o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketOrders)
		return b.ForEach(func(k, v []byte) error {
			var order Order
			if err := json.Unmarshal(v, &order); err != nil {
				return nil
			}
			if order.Status == OrderConfirmed && order.Type == "buy" {
				totalSPC += order.AmountSPC
				totalEUR += order.AmountEUR
				count++
			}
			return nil
		})
	})
	return
}

// RegisterWallet adds a wallet address to the registry (deduplicates).
func (o *OrderDB) RegisterWallet(address string) error {
	return o.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketWallets)
		existing := b.Get([]byte(address))
		if existing != nil {
			return nil // already registered
		}
		data, _ := json.Marshal(time.Now().Unix())
		return b.Put([]byte(address), data)
	})
}

// WalletCount returns the total number of registered wallets.
func (o *OrderDB) WalletCount() (int, error) {
	var count int
	err := o.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucketWallets)
		count = b.Stats().KeyN
		return nil
	})
	return count, err
}
