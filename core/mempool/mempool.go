package mempool

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// ErrDuplicate se devuelve cuando se intenta añadir una tx ya presente en el mempool.
var ErrDuplicate = errors.New("transacción duplicada")

// ErrFull se devuelve cuando el mempool ha alcanzado su capacidad máxima.
var ErrFull = errors.New("mempool lleno")

// ErrExpired se devuelve cuando la transacción es demasiado antigua.
var ErrExpired = errors.New("transacción expirada")

// defaultMaxAgeSeconds es la antigüedad máxima por defecto para una tx (1 hora).
const defaultMaxAgeSeconds int64 = 3600

// Mempool almacena transacciones pendientes de ser incluidas en un bloque.
// Las transacciones se priorizan por fee (mayor fee = mayor prioridad).
type Mempool struct {
	txs     map[crypto.Hash]*block.Transaction
	mu      sync.RWMutex
	maxSize int
}

// NewMempool crea un nuevo mempool con el tamaño máximo indicado.
// maxSize es el número máximo de transacciones pendientes admitidas.
func NewMempool(maxSize int) *Mempool {
	return &Mempool{
		txs:     make(map[crypto.Hash]*block.Transaction),
		maxSize: maxSize,
	}
}

// Add añade una transacción al mempool tras validarla.
// Rechaza duplicados, transacciones expiradas y si el mempool está lleno.
func (m *Mempool) Add(tx *block.Transaction) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Comprobar duplicado
	if _, exists := m.txs[tx.ID]; exists {
		return ErrDuplicate
	}

	// Comprobar expiración (timestamp en nanosegundos)
	now := time.Now().Unix()
	txTime := tx.Timestamp / 1_000_000_000 // nanosegundos → segundos
	if now-txTime > defaultMaxAgeSeconds {
		return ErrExpired
	}

	// Comprobar capacidad
	if len(m.txs) >= m.maxSize {
		return ErrFull
	}

	m.txs[tx.ID] = tx
	return nil
}

// Remove elimina una transacción del mempool por su hash.
func (m *Mempool) Remove(hash crypto.Hash) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.txs, hash)
}

// Has indica si el mempool contiene una transacción con el hash dado.
func (m *Mempool) Has(hash crypto.Hash) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.txs[hash]
	return ok
}

// Get devuelve la transacción con el hash dado, o false si no existe.
func (m *Mempool) Get(hash crypto.Hash) (*block.Transaction, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tx, ok := m.txs[hash]
	return tx, ok
}

// Pending devuelve todas las transacciones pendientes ordenadas por fee descendente.
func (m *Mempool) Pending() []*block.Transaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	txs := make([]*block.Transaction, 0, len(m.txs))
	for _, tx := range m.txs {
		txs = append(txs, tx)
	}

	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Fee > txs[j].Fee
	})
	return txs
}

// SelectTxs devuelve las top maxCount transacciones con mayor fee.
// Si hay menos de maxCount txs pendientes, devuelve todas.
func (m *Mempool) SelectTxs(maxCount int) []*block.Transaction {
	pending := m.Pending()
	if maxCount >= len(pending) {
		return pending
	}
	return pending[:maxCount]
}

// Size devuelve el número de transacciones actualmente en el mempool.
func (m *Mempool) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.txs)
}

// Flush elimina del mempool las transacciones cuyos hashes se proporcionan.
// Se llama típicamente tras confirmar un bloque para limpiar las txs incluidas.
func (m *Mempool) Flush(hashes []crypto.Hash) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, h := range hashes {
		delete(m.txs, h)
	}
}

// PruneExpired elimina las transacciones más antiguas que maxAgeSeconds segundos.
func (m *Mempool) PruneExpired(maxAgeSeconds int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().Unix()
	for hash, tx := range m.txs {
		txTime := tx.Timestamp / 1_000_000_000 // nanosegundos → segundos
		if now-txTime > maxAgeSeconds {
			delete(m.txs, hash)
		}
	}
}
