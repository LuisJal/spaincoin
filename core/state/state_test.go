package state

import (
	"testing"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// helpers

func randomAddress(seed byte) crypto.Address {
	var addr crypto.Address
	for i := range addr {
		addr[i] = seed + byte(i)
	}
	return addr
}

func TestNewState_Empty(t *testing.T) {
	s := NewState()
	if s == nil {
		t.Fatal("NewState devolvió nil")
	}
	if s.TotalSupply() != 0 {
		t.Errorf("estado vacío: supply esperado 0, obtenido %d", s.TotalSupply())
	}
	_, ok := s.GetAccount(randomAddress(1))
	if ok {
		t.Error("estado vacío no debería tener cuentas")
	}
}

func TestApplyTransaction_Coinbase(t *testing.T) {
	s := NewState()
	to := randomAddress(1)
	tx := block.NewCoinbaseTx(to, 1_000_000)

	if err := s.ApplyTransaction(tx); err != nil {
		t.Fatalf("coinbase falló: %v", err)
	}

	bal := s.GetBalance(to)
	if bal != 1_000_000 {
		t.Errorf("saldo esperado 1_000_000, obtenido %d", bal)
	}
}

func TestApplyTransaction_Transfer(t *testing.T) {
	s := NewState()
	sender := randomAddress(1)
	receiver := randomAddress(2)

	// Dar saldo inicial al emisor con coinbase
	coinbase := block.NewCoinbaseTx(sender, 500)
	s.ApplyTransaction(coinbase)

	// Transferencia: 100 pesetas + 10 de fee
	tx := block.NewTransaction(sender, receiver, 100, 0, 10)

	if err := s.ApplyTransaction(tx); err != nil {
		t.Fatalf("transferencia falló: %v", err)
	}

	senderBal := s.GetBalance(sender)
	receiverBal := s.GetBalance(receiver)

	if senderBal != 390 { // 500 - 100 - 10
		t.Errorf("saldo emisor esperado 390, obtenido %d", senderBal)
	}
	if receiverBal != 100 {
		t.Errorf("saldo receptor esperado 100, obtenido %d", receiverBal)
	}
}

func TestApplyTransaction_InsufficientBalance(t *testing.T) {
	s := NewState()
	sender := randomAddress(1)
	receiver := randomAddress(2)

	// Dar solo 50 pesetas
	coinbase := block.NewCoinbaseTx(sender, 50)
	s.ApplyTransaction(coinbase)

	// Intentar enviar 100
	tx := block.NewTransaction(sender, receiver, 100, 0, 0)
	err := s.ApplyTransaction(tx)
	if err == nil {
		t.Error("debería haber retornado error por saldo insuficiente")
	}
}

func TestApplyTransaction_WrongNonce(t *testing.T) {
	s := NewState()
	sender := randomAddress(1)
	receiver := randomAddress(2)

	coinbase := block.NewCoinbaseTx(sender, 1000)
	s.ApplyTransaction(coinbase)

	// Nonce incorrecto (esperado 0, enviamos 5)
	tx := block.NewTransaction(sender, receiver, 100, 5, 0)
	err := s.ApplyTransaction(tx)
	if err == nil {
		t.Error("debería haber retornado error por nonce incorrecto")
	}
}

func TestApplyBlock_Rollback(t *testing.T) {
	s := NewState()
	sender := randomAddress(1)
	receiver := randomAddress(2)
	badAddr := randomAddress(3)

	// Financiar al emisor
	coinbase := block.NewCoinbaseTx(sender, 200)
	s.ApplyTransaction(coinbase)

	initialBalance := s.GetBalance(sender)

	// Bloque con dos txs: la primera válida, la segunda fallará (saldo insuficiente)
	tx1 := block.NewTransaction(sender, receiver, 50, 0, 0)
	// tx2 intenta enviar más de lo que queda tras tx1
	tx2 := block.NewTransaction(badAddr, receiver, 9999, 0, 0) // badAddr sin saldo

	b := block.NewBlock(1, crypto.Hash{}, sender, []*block.Transaction{tx1, tx2})

	err := s.ApplyBlock(b)
	if err == nil {
		t.Fatal("ApplyBlock debería haber fallado")
	}

	// El estado debe haberse revertido: emisor conserva su saldo original
	if s.GetBalance(sender) != initialBalance {
		t.Errorf("rollback fallido: saldo emisor esperado %d, obtenido %d",
			initialBalance, s.GetBalance(sender))
	}
	if s.GetBalance(receiver) != 0 {
		t.Errorf("rollback fallido: saldo receptor debería ser 0, obtenido %d",
			s.GetBalance(receiver))
	}
}

func TestStateHash_Deterministic(t *testing.T) {
	s1 := NewState()
	s2 := NewState()

	addr1 := randomAddress(10)
	addr2 := randomAddress(20)

	// Aplicar mismas operaciones en el mismo orden
	for _, s := range []*State{s1, s2} {
		s.ApplyTransaction(block.NewCoinbaseTx(addr1, 1000))
		s.ApplyTransaction(block.NewCoinbaseTx(addr2, 2000))
	}

	h1 := s1.Hash()
	h2 := s2.Hash()

	if h1 != h2 {
		t.Errorf("hashes distintos para el mismo estado: %s vs %s", h1, h2)
	}

	// Modificar s2 y comprobar que el hash cambia
	s2.ApplyTransaction(block.NewCoinbaseTx(addr2, 1))
	h3 := s2.Hash()
	if h1 == h3 {
		t.Error("el hash no cambió tras modificar el estado")
	}
}

func TestTotalSupply(t *testing.T) {
	s := NewState()
	addr1 := randomAddress(1)
	addr2 := randomAddress(2)

	// Dos coinbases: supply total = 1000 + 500 = 1500
	s.ApplyTransaction(block.NewCoinbaseTx(addr1, 1000))
	s.ApplyTransaction(block.NewCoinbaseTx(addr2, 500))

	if s.TotalSupply() != 1500 {
		t.Errorf("supply esperado 1500, obtenido %d", s.TotalSupply())
	}

	// Transferencia de addr1 a addr2: 100 pesetas + 10 fee (fee se quema)
	tx := block.NewTransaction(addr1, addr2, 100, 0, 10)
	s.ApplyTransaction(tx)

	// Supply disminuye por la fee quemada
	expectedSupply := uint64(1500 - 10)
	if s.TotalSupply() != expectedSupply {
		t.Errorf("supply tras fee esperado %d, obtenido %d", expectedSupply, s.TotalSupply())
	}
}
