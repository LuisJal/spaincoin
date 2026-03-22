package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// tempDB crea una base de datos temporal y devuelve el handle y una función de limpieza.
func tempDB(t *testing.T) (*DB, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "spaincoin-storage-test-*")
	if err != nil {
		t.Fatalf("no se pudo crear directorio temporal: %v", err)
	}
	dbPath := filepath.Join(dir, "test.db")
	db, err := Open(dbPath)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatalf("no se pudo abrir la base de datos: %v", err)
	}
	cleanup := func() {
		db.Close()
		os.RemoveAll(dir)
	}
	return db, cleanup
}

// makeAddr crea una dirección de prueba con el byte indicado relleno.
func makeAddr(b byte) crypto.Address {
	var addr crypto.Address
	for i := range addr {
		addr[i] = b
	}
	return addr
}

func TestDB_OpenClose(t *testing.T) {
	dir, err := os.MkdirTemp("", "spaincoin-openclose-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(dir)

	db, err := Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestDB_SaveAndGetByHeight(t *testing.T) {
	db, cleanup := tempDB(t)
	defer cleanup()

	addr := makeAddr(0xAB)
	genesis := block.GenesisBlock(addr, 1000)

	if err := db.SaveBlock(genesis); err != nil {
		t.Fatalf("SaveBlock: %v", err)
	}

	got, err := db.GetBlockByHeight(0)
	if err != nil {
		t.Fatalf("GetBlockByHeight: %v", err)
	}

	if got.Header.Height != genesis.Header.Height {
		t.Errorf("Height: quería %d, obtuve %d", genesis.Header.Height, got.Header.Height)
	}
	if got.Hash != genesis.Hash {
		t.Errorf("Hash: quería %s, obtuve %s", genesis.Hash, got.Hash)
	}
	if got.Header.Timestamp != genesis.Header.Timestamp {
		t.Errorf("Timestamp: quería %d, obtuve %d", genesis.Header.Timestamp, got.Header.Timestamp)
	}
	if got.Header.Validator != genesis.Header.Validator {
		t.Errorf("Validator: quería %s, obtuve %s", genesis.Header.Validator, got.Header.Validator)
	}
	if len(got.Transactions) != len(genesis.Transactions) {
		t.Errorf("len(Transactions): quería %d, obtuve %d", len(genesis.Transactions), len(got.Transactions))
	}
}

func TestDB_GetByHash(t *testing.T) {
	db, cleanup := tempDB(t)
	defer cleanup()

	addr := makeAddr(0x01)
	genesis := block.GenesisBlock(addr, 500)

	if err := db.SaveBlock(genesis); err != nil {
		t.Fatalf("SaveBlock: %v", err)
	}

	hashHex := genesis.Hash.String()
	got, err := db.GetBlockByHash(hashHex)
	if err != nil {
		t.Fatalf("GetBlockByHash: %v", err)
	}

	if got.Hash != genesis.Hash {
		t.Errorf("Hash: quería %s, obtuve %s", genesis.Hash, got.Hash)
	}
	if got.Header.Height != genesis.Header.Height {
		t.Errorf("Height: quería %d, obtuve %d", genesis.Header.Height, got.Header.Height)
	}
}

func TestDB_GetHeight(t *testing.T) {
	db, cleanup := tempDB(t)
	defer cleanup()

	// Base de datos vacía → -1
	h, err := db.GetHeight()
	if err != nil {
		t.Fatalf("GetHeight (vacío): %v", err)
	}
	if h != -1 {
		t.Errorf("GetHeight vacío: quería -1, obtuve %d", h)
	}

	addr := makeAddr(0x02)
	genesis := block.GenesisBlock(addr, 1000)
	if err := db.SaveBlock(genesis); err != nil {
		t.Fatalf("SaveBlock genesis: %v", err)
	}

	h, err = db.GetHeight()
	if err != nil {
		t.Fatalf("GetHeight tras génesis: %v", err)
	}
	if h != 0 {
		t.Errorf("GetHeight tras génesis: quería 0, obtuve %d", h)
	}

	// Añadir bloque 1
	block1 := block.NewBlock(1, genesis.Hash, addr, nil)
	if err := db.SaveBlock(block1); err != nil {
		t.Fatalf("SaveBlock bloque 1: %v", err)
	}

	h, err = db.GetHeight()
	if err != nil {
		t.Fatalf("GetHeight tras bloque 1: %v", err)
	}
	if h != 1 {
		t.Errorf("GetHeight tras bloque 1: quería 1, obtuve %d", h)
	}
}

func TestDB_LoadAllBlocks(t *testing.T) {
	db, cleanup := tempDB(t)
	defer cleanup()

	addr := makeAddr(0x03)
	genesis := block.GenesisBlock(addr, 1000)
	b1 := block.NewBlock(1, genesis.Hash, addr, nil)
	b2 := block.NewBlock(2, b1.Hash, addr, nil)

	for _, b := range []*block.Block{genesis, b1, b2} {
		if err := db.SaveBlock(b); err != nil {
			t.Fatalf("SaveBlock altura %d: %v", b.Header.Height, err)
		}
	}

	all, err := db.LoadAllBlocks()
	if err != nil {
		t.Fatalf("LoadAllBlocks: %v", err)
	}

	if len(all) != 3 {
		t.Fatalf("LoadAllBlocks: quería 3 bloques, obtuve %d", len(all))
	}

	// Verificar orden ascendente
	for i, b := range all {
		if b.Header.Height != uint64(i) {
			t.Errorf("bloque[%d]: altura esperada %d, obtuve %d", i, i, b.Header.Height)
		}
	}

	// Verificar que los hashes coinciden
	expected := []*block.Block{genesis, b1, b2}
	for i, b := range all {
		if b.Hash != expected[i].Hash {
			t.Errorf("bloque[%d]: hash esperado %s, obtuve %s", i, expected[i].Hash, b.Hash)
		}
	}
}

func TestDB_Persistence(t *testing.T) {
	dir, err := os.MkdirTemp("", "spaincoin-persist-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	defer os.RemoveAll(dir)

	dbPath := filepath.Join(dir, "persist.db")
	addr := makeAddr(0xCC)
	genesis := block.GenesisBlock(addr, 42_000_000)

	// Primera sesión: guardar y cerrar
	{
		db, err := Open(dbPath)
		if err != nil {
			t.Fatalf("Open (primera sesión): %v", err)
		}
		if err := db.SaveBlock(genesis); err != nil {
			t.Fatalf("SaveBlock: %v", err)
		}
		if err := db.Close(); err != nil {
			t.Fatalf("Close: %v", err)
		}
	}

	// Segunda sesión: reabrir y verificar
	{
		db, err := Open(dbPath)
		if err != nil {
			t.Fatalf("Open (segunda sesión): %v", err)
		}
		defer db.Close()

		got, err := db.GetBlockByHeight(0)
		if err != nil {
			t.Fatalf("GetBlockByHeight tras reabrir: %v", err)
		}

		if got.Hash != genesis.Hash {
			t.Errorf("Persistencia: hash no coincide. Quería %s, obtuve %s", genesis.Hash, got.Hash)
		}
		if got.Header.Height != genesis.Header.Height {
			t.Errorf("Persistencia: altura no coincide. Quería %d, obtuve %d", genesis.Header.Height, got.Header.Height)
		}
		if len(got.Transactions) != len(genesis.Transactions) {
			t.Errorf("Persistencia: número de transacciones no coincide")
		}

		h, err := db.GetHeight()
		if err != nil {
			t.Fatalf("GetHeight tras reabrir: %v", err)
		}
		if h != 0 {
			t.Errorf("Persistencia: GetHeight quería 0, obtuve %d", h)
		}
	}
}
