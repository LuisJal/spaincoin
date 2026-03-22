package state

import (
	"encoding/binary"
	"errors"
	"sort"
	"sync"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// State es el libro mayor global — contiene todas las cuentas y sus saldos.
type State struct {
	accounts map[crypto.Address]*Account
	mu       sync.RWMutex
}

// NewState crea un estado vacío.
func NewState() *State {
	return &State{
		accounts: make(map[crypto.Address]*Account),
	}
}

// GetAccount devuelve la cuenta para la dirección dada, o false si no existe.
func (s *State) GetAccount(addr crypto.Address) (*Account, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	acc, ok := s.accounts[addr]
	return acc, ok
}

// GetOrCreate devuelve la cuenta existente o crea una nueva si no existe.
func (s *State) GetOrCreate(addr crypto.Address) *Account {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getOrCreateLocked(addr)
}

// getOrCreateLocked es la versión interna sin lock (debe llamarse con mu ya tomado).
func (s *State) getOrCreateLocked(addr crypto.Address) *Account {
	if acc, ok := s.accounts[addr]; ok {
		return acc
	}
	acc := NewAccount(addr)
	s.accounts[addr] = acc
	return acc
}

// GetBalance devuelve el saldo de la dirección (0 si no existe la cuenta).
func (s *State) GetBalance(addr crypto.Address) uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if acc, ok := s.accounts[addr]; ok {
		return acc.Balance
	}
	return 0
}

// ApplyTransaction valida y aplica una transacción al estado.
// Para coinbase (From == dirección cero) se omiten las comprobaciones de saldo y nonce del emisor.
// Errores posibles: saldo insuficiente, nonce incorrecto, auto-envío.
func (s *State) ApplyTransaction(tx *block.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.applyTransactionLocked(tx)
}

// applyTransactionLocked aplica una tx sin adquirir el lock (debe llamarse con mu tomado).
func (s *State) applyTransactionLocked(tx *block.Transaction) error {
	// Prohibir auto-envíos (excepto coinbase que tiene From == zero)
	if !tx.IsCoinbase() && tx.From == tx.To {
		return errors.New("auto-envío no permitido")
	}

	if tx.IsCoinbase() {
		// Coinbase: simplemente acreditar el importe al destinatario
		receiver := s.getOrCreateLocked(tx.To)
		receiver.Balance += tx.Amount
		return nil
	}

	// Transacción normal
	sender := s.getOrCreateLocked(tx.From)

	// Comprobar nonce
	if tx.Nonce != sender.Nonce {
		return errors.New("nonce incorrecto")
	}

	// Comprobar saldo suficiente para amount + fee
	total := tx.Amount + tx.Fee
	if total < tx.Amount {
		// desbordamiento aritmético
		return errors.New("saldo insuficiente (desbordamiento)")
	}
	if sender.Balance < total {
		return errors.New("saldo insuficiente")
	}

	// Aplicar cambios
	sender.Balance -= total
	sender.Nonce++

	receiver := s.getOrCreateLocked(tx.To)
	receiver.Balance += tx.Amount
	// La fee se quema (se descontará del total supply hasta que consenso la distribuya)

	return nil
}

// ApplyBlock aplica todas las transacciones de un bloque en orden.
// Si alguna transacción falla, se revierten todos los cambios del bloque.
func (s *State) ApplyBlock(b *block.Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Snapshot antes de aplicar para poder revertir
	snapshot := s.cloneLocked()

	for _, tx := range b.Transactions {
		if err := s.applyTransactionLocked(tx); err != nil {
			// Revertir: reemplazar el mapa de cuentas con la copia previa
			s.accounts = snapshot.accounts
			return err
		}
	}
	return nil
}

// Hash calcula un hash determinista del estado completo.
// Ordena las direcciones y hashea todos los datos de cada cuenta.
func (s *State) Hash() crypto.Hash {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Ordenar las direcciones para garantizar determinismo
	addrs := make([]crypto.Address, 0, len(s.accounts))
	for addr := range s.accounts {
		addrs = append(addrs, addr)
	}
	sort.Slice(addrs, func(i, j int) bool {
		for k := 0; k < 20; k++ {
			if addrs[i][k] != addrs[j][k] {
				return addrs[i][k] < addrs[j][k]
			}
		}
		return false
	})

	buf := make([]byte, 0, len(addrs)*(20+8+8))
	b8 := make([]byte, 8)
	for _, addr := range addrs {
		acc := s.accounts[addr]
		buf = append(buf, addr[:]...)
		binary.BigEndian.PutUint64(b8, acc.Balance)
		buf = append(buf, b8...)
		binary.BigEndian.PutUint64(b8, acc.Nonce)
		buf = append(buf, b8...)
	}
	return crypto.HashBytes(buf)
}

// Clone devuelve una copia profunda del estado (útil para rollbacks externos).
func (s *State) Clone() *State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cloneLocked()
}

// cloneLocked realiza la copia sin adquirir el lock.
func (s *State) cloneLocked() *State {
	clone := &State{
		accounts: make(map[crypto.Address]*Account, len(s.accounts)),
	}
	for addr, acc := range s.accounts {
		copied := *acc // copia de valor
		clone.accounts[addr] = &copied
	}
	return clone
}

// TotalSupply devuelve la suma de todos los saldos del estado.
func (s *State) TotalSupply() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var total uint64
	for _, acc := range s.accounts {
		total += acc.Balance
	}
	return total
}
