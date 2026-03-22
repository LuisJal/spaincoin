package mempool

import (
	"testing"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// helpers

func makeAddr(seed byte) crypto.Address {
	var addr crypto.Address
	for i := range addr {
		addr[i] = seed + byte(i)
	}
	return addr
}

// newTxWithFee crea una transacción con la fee indicada y timestamp actual.
func newTxWithFee(fee uint64, nonce uint64) *block.Transaction {
	from := makeAddr(1)
	to := makeAddr(2)
	return block.NewTransaction(from, to, 100, nonce, fee)
}

func TestMempool_AddAndGet(t *testing.T) {
	m := NewMempool(100)
	tx := newTxWithFee(10, 0)

	if err := m.Add(tx); err != nil {
		t.Fatalf("Add falló: %v", err)
	}

	got, ok := m.Get(tx.ID)
	if !ok {
		t.Fatal("Get debería encontrar la tx añadida")
	}
	if got.ID != tx.ID {
		t.Error("tx recuperada con ID distinto")
	}

	if !m.Has(tx.ID) {
		t.Error("Has debería devolver true")
	}
	if m.Size() != 1 {
		t.Errorf("Size esperado 1, obtenido %d", m.Size())
	}
}

func TestMempool_NoDuplicates(t *testing.T) {
	m := NewMempool(100)
	tx := newTxWithFee(10, 0)

	if err := m.Add(tx); err != nil {
		t.Fatalf("primer Add falló: %v", err)
	}
	err := m.Add(tx)
	if err == nil {
		t.Fatal("segundo Add debería haber fallado con ErrDuplicate")
	}
	if err != ErrDuplicate {
		t.Errorf("error esperado ErrDuplicate, obtenido: %v", err)
	}
	if m.Size() != 1 {
		t.Errorf("Size esperado 1 tras duplicado, obtenido %d", m.Size())
	}
}

func TestMempool_Full(t *testing.T) {
	maxSize := 3
	m := NewMempool(maxSize)

	for i := 0; i < maxSize; i++ {
		tx := newTxWithFee(uint64(i), uint64(i))
		if err := m.Add(tx); err != nil {
			t.Fatalf("Add %d falló: %v", i, err)
		}
	}

	// La siguiente tx debe fallar
	extraTx := newTxWithFee(999, uint64(maxSize))
	err := m.Add(extraTx)
	if err == nil {
		t.Fatal("Add debería haber fallado con ErrFull")
	}
	if err != ErrFull {
		t.Errorf("error esperado ErrFull, obtenido: %v", err)
	}
}

func TestMempool_SelectTxs_ByFee(t *testing.T) {
	m := NewMempool(100)

	// Añadir txs con distintas fees
	fees := []uint64{5, 30, 10, 50, 20}
	for i, fee := range fees {
		tx := newTxWithFee(fee, uint64(i))
		m.Add(tx)
	}

	selected := m.SelectTxs(3)
	if len(selected) != 3 {
		t.Fatalf("SelectTxs esperaba 3 txs, obtuvo %d", len(selected))
	}

	// Las 3 primeras deben ser las de mayor fee: 50, 30, 20
	expectedFees := []uint64{50, 30, 20}
	for i, tx := range selected {
		if tx.Fee != expectedFees[i] {
			t.Errorf("tx[%d]: fee esperada %d, obtenida %d", i, expectedFees[i], tx.Fee)
		}
	}
}

func TestMempool_Flush(t *testing.T) {
	m := NewMempool(100)

	tx1 := newTxWithFee(10, 0)
	tx2 := newTxWithFee(20, 1)
	tx3 := newTxWithFee(30, 2)

	m.Add(tx1)
	m.Add(tx2)
	m.Add(tx3)

	if m.Size() != 3 {
		t.Fatalf("Size esperado 3, obtenido %d", m.Size())
	}

	// Flush tx1 y tx2
	m.Flush([]crypto.Hash{tx1.ID, tx2.ID})

	if m.Size() != 1 {
		t.Errorf("Size esperado 1 tras flush, obtenido %d", m.Size())
	}
	if m.Has(tx1.ID) || m.Has(tx2.ID) {
		t.Error("tx1 o tx2 siguen presentes tras flush")
	}
	if !m.Has(tx3.ID) {
		t.Error("tx3 debería seguir presente")
	}
}

func TestMempool_PruneExpired(t *testing.T) {
	m := NewMempool(100)

	// Tx fresca (ahora mismo)
	freshTx := newTxWithFee(10, 0)
	m.Add(freshTx)

	// Tx artificial con timestamp muy antiguo (hace 2 horas)
	oldTx := &block.Transaction{
		From:      makeAddr(10),
		To:        makeAddr(20),
		Amount:    50,
		Nonce:     0,
		Fee:       5,
		Timestamp: time.Now().Add(-2 * time.Hour).UnixNano(),
	}
	oldTx.ID = oldTx.Hash()

	// Insertar directamente para evitar la validación de expiración de Add
	m.mu.Lock()
	m.txs[oldTx.ID] = oldTx
	m.mu.Unlock()

	if m.Size() != 2 {
		t.Fatalf("Size esperado 2, obtenido %d", m.Size())
	}

	// Podar las txs con más de 3600 segundos (1 hora)
	m.PruneExpired(3600)

	if m.Size() != 1 {
		t.Errorf("Size esperado 1 tras pruneExpired, obtenido %d", m.Size())
	}
	if !m.Has(freshTx.ID) {
		t.Error("la tx fresca no debería haber sido eliminada")
	}
	if m.Has(oldTx.ID) {
		t.Error("la tx antigua debería haber sido eliminada")
	}
}
