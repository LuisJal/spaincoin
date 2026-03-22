package chain

import (
	"testing"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// testValidator genera una dirección de validador para los tests.
func testValidator(t *testing.T) crypto.Address {
	t.Helper()
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	return pub.ToAddress()
}

// newTestBlockchain crea una blockchain de prueba con 1 SPC de supply inicial.
const initialSupply = uint64(1_000_000_000_000_000_000) // 1 SPC en pesetas

func newTestBlockchain(t *testing.T) (*Blockchain, crypto.Address) {
	t.Helper()
	validator := testValidator(t)
	bc, err := NewBlockchain(validator, initialSupply)
	if err != nil {
		t.Fatalf("NewBlockchain: %v", err)
	}
	return bc, validator
}

// buildValidBlock construye un bloque válido sobre el último bloque de la cadena.
func buildValidBlock(t *testing.T, bc *Blockchain, validator crypto.Address, txs []*block.Transaction) *block.Block {
	t.Helper()
	last := bc.LatestBlock()
	prevHash := last.Hash
	height := bc.Height() + 1

	b := block.NewBlock(height, prevHash, validator, txs)

	// Asegurarse de que el timestamp es estrictamente mayor que el del bloque anterior
	for b.Header.Timestamp <= last.Header.Timestamp {
		time.Sleep(time.Nanosecond)
		b = block.NewBlock(height, prevHash, validator, txs)
	}
	return b
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestNewBlockchain(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	// Debe existir el bloque génesis
	if bc.Height() != 0 {
		t.Errorf("Height after genesis: got %d, want 0", bc.Height())
	}

	genesis, ok := bc.GetBlock(0)
	if !ok {
		t.Fatal("GetBlock(0) should return genesis block")
	}
	if genesis.Header.Height != 0 {
		t.Errorf("Genesis height: got %d, want 0", genesis.Header.Height)
	}

	// El supply inicial debe estar en el estado
	balance := bc.State().GetBalance(validator)
	if balance != initialSupply {
		t.Errorf("Genesis balance: got %d, want %d", balance, initialSupply)
	}

	// El supply total debe coincidir
	if bc.State().TotalSupply() != initialSupply {
		t.Errorf("TotalSupply: got %d, want %d", bc.State().TotalSupply(), initialSupply)
	}
}

func TestAddBlock_Valid(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	b := buildValidBlock(t, bc, validator, nil)

	if err := bc.AddBlock(b); err != nil {
		t.Fatalf("AddBlock: %v", err)
	}

	if bc.Height() != 1 {
		t.Errorf("Height after adding block: got %d, want 1", bc.Height())
	}

	// El estado debe seguir siendo coherente
	if bc.State().TotalSupply() != initialSupply {
		t.Errorf("TotalSupply should remain unchanged, got %d", bc.State().TotalSupply())
	}
}

func TestAddBlock_WrongHeight(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	last := bc.LatestBlock()
	// Construir un bloque con altura incorrecta (2 en lugar de 1)
	b := block.NewBlock(2, last.Hash, validator, nil)
	// Asegurar timestamp mayor
	for b.Header.Timestamp <= last.Header.Timestamp {
		time.Sleep(time.Nanosecond)
		b = block.NewBlock(2, last.Hash, validator, nil)
	}

	if err := bc.AddBlock(b); err == nil {
		t.Error("AddBlock with wrong height should return error")
	}

	if bc.Height() != 0 {
		t.Errorf("Height should remain 0 after rejected block, got %d", bc.Height())
	}
}

func TestAddBlock_WrongPrevHash(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	last := bc.LatestBlock()
	wrongHash := crypto.HashBytes([]byte("not the real prev hash"))
	b := block.NewBlock(1, wrongHash, validator, nil)
	for b.Header.Timestamp <= last.Header.Timestamp {
		time.Sleep(time.Nanosecond)
		b = block.NewBlock(1, wrongHash, validator, nil)
	}

	if err := bc.AddBlock(b); err == nil {
		t.Error("AddBlock with wrong prevHash should return error")
	}
}

func TestAddBlock_OldTimestamp(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	last := bc.LatestBlock()
	// Construir un bloque con timestamp viejo manualmente
	b := block.NewBlock(1, last.Hash, validator, nil)
	// Forzar un timestamp igual o anterior al del bloque anterior
	b.Header.Timestamp = last.Header.Timestamp
	b.Hash = b.ComputeHash() // recalcular hash después de modificar

	if err := bc.AddBlock(b); err == nil {
		t.Error("AddBlock with old timestamp should return error")
	}
}

func TestGetBlock(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	b := buildValidBlock(t, bc, validator, nil)
	if err := bc.AddBlock(b); err != nil {
		t.Fatalf("AddBlock: %v", err)
	}

	// Recuperar por altura
	retrieved, ok := bc.GetBlock(1)
	if !ok {
		t.Fatal("GetBlock(1) should return a block")
	}
	if retrieved.Header.Height != 1 {
		t.Errorf("GetBlock(1) height: got %d, want 1", retrieved.Header.Height)
	}

	// Altura inexistente
	_, ok = bc.GetBlock(99)
	if ok {
		t.Error("GetBlock(99) should return false for non-existent height")
	}
}

func TestGetBlockByHash(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	b := buildValidBlock(t, bc, validator, nil)
	if err := bc.AddBlock(b); err != nil {
		t.Fatalf("AddBlock: %v", err)
	}

	// Recuperar por hash
	latestHash := bc.LatestBlock().Hash
	retrieved, ok := bc.GetBlockByHash(latestHash)
	if !ok {
		t.Fatal("GetBlockByHash should find the latest block by hash")
	}
	if retrieved.Hash != latestHash {
		t.Errorf("GetBlockByHash returned wrong block")
	}

	// Hash inexistente
	fakeHash := crypto.HashBytes([]byte("nonexistent"))
	_, ok = bc.GetBlockByHash(fakeHash)
	if ok {
		t.Error("GetBlockByHash should return false for nonexistent hash")
	}
}

func TestIsValid_Clean(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	if !bc.IsValid() {
		t.Error("Fresh blockchain with only genesis should be valid")
	}

	// Añadir un bloque más y comprobar que sigue siendo válida
	b := buildValidBlock(t, bc, validator, nil)
	if err := bc.AddBlock(b); err != nil {
		t.Fatalf("AddBlock: %v", err)
	}

	if !bc.IsValid() {
		t.Error("Blockchain should be valid after adding a valid block")
	}
}

func TestIsValid_Tampered(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	b := buildValidBlock(t, bc, validator, nil)
	if err := bc.AddBlock(b); err != nil {
		t.Fatalf("AddBlock: %v", err)
	}

	// Manipular el hash del bloque génesis directamente
	// (esto rompe la cadena de PrevHash del bloque siguiente)
	bc.mu.Lock()
	bc.blocks[0].Hash = crypto.HashBytes([]byte("tampered genesis hash"))
	bc.mu.Unlock()

	if bc.IsValid() {
		t.Error("IsValid should return false after tampering genesis hash")
	}
}

func TestChain_MempoolFlush(t *testing.T) {
	bc, validator := newTestBlockchain(t)

	// Crear una dirección destino
	to := testValidator(t)

	// El validador tiene fondos del génesis; crear una tx normal
	// Nonce del validator es 0, así que usamos nonce=0
	tx := block.NewTransaction(validator, to, 1000, 0, 10)

	// Añadir la tx al mempool
	mp := bc.Mempool()
	if err := mp.Add(tx); err != nil {
		t.Fatalf("Mempool.Add: %v", err)
	}

	if mp.Size() != 1 {
		t.Fatalf("Mempool size before block: got %d, want 1", mp.Size())
	}

	// Seleccionar txs del mempool para incluir en el bloque
	selectedTxs := mp.SelectTxs(100)

	// Construir y añadir el bloque con esas txs
	b := buildValidBlock(t, bc, validator, selectedTxs)
	if err := bc.AddBlock(b); err != nil {
		t.Fatalf("AddBlock: %v", err)
	}

	// Las txs confirmadas deben haberse eliminado del mempool
	if mp.Size() != 0 {
		t.Errorf("Mempool size after block confirmation: got %d, want 0", mp.Size())
	}
}
