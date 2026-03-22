package chain

import (
	"errors"
	"sync"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
	"github.com/spaincoin/spaincoin/core/mempool"
	"github.com/spaincoin/spaincoin/core/state"
)

// Blockchain mantiene la cadena de bloques completa, el estado global y el mempool.
type Blockchain struct {
	blocks  []*block.Block // todos los bloques en orden
	state   *state.State   // estado actual del mundo
	mempool *mempool.Mempool
	mu      sync.RWMutex
}

// NewBlockchain crea una nueva blockchain inicializando el bloque génesis.
// genesisValidator es la dirección que recibirá el supply inicial.
// initialSupply es la cantidad (en pesetas) acreditada en el bloque génesis.
func NewBlockchain(genesisValidator crypto.Address, initialSupply uint64) (*Blockchain, error) {
	genesis := block.GenesisBlock(genesisValidator, initialSupply)

	s := state.NewState()
	if err := s.ApplyBlock(genesis); err != nil {
		return nil, err
	}

	// Actualizar el StateRoot del bloque génesis tras aplicarlo
	genesis.Header.StateRoot = s.Hash()
	// Recalcular el hash del bloque después de modificar el header
	genesis.Hash = genesis.ComputeHash()

	bc := &Blockchain{
		blocks:  []*block.Block{genesis},
		state:   s,
		mempool: mempool.NewMempool(10000),
	}
	return bc, nil
}

// AddBlock valida y añade un nuevo bloque a la cadena.
// Comprueba: integridad del bloque, altura, hash previo y timestamp.
// Aplica el bloque al estado y limpia el mempool de las txs confirmadas.
func (bc *Blockchain) AddBlock(b *block.Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	// Validar la integridad interna del bloque
	if err := b.Validate(); err != nil {
		return err
	}

	last := bc.blocks[len(bc.blocks)-1]

	// Comprobar altura secuencial
	expectedHeight := last.Header.Height + 1
	if b.Header.Height != expectedHeight {
		return errors.New("altura de bloque incorrecta")
	}

	// Comprobar que el hash previo coincide con el último bloque
	if b.Header.PrevHash != last.Hash {
		return errors.New("hash previo incorrecto")
	}

	// Comprobar que el timestamp es estrictamente mayor que el del bloque anterior
	if b.Header.Timestamp <= last.Header.Timestamp {
		return errors.New("timestamp del bloque no es mayor que el del bloque anterior")
	}

	// Aplicar el bloque al estado
	if err := bc.state.ApplyBlock(b); err != nil {
		return err
	}

	// Actualizar el StateRoot en la cabecera con el hash del nuevo estado
	b.Header.StateRoot = bc.state.Hash()
	// Recalcular el hash del bloque después de modificar el header
	b.Hash = b.ComputeHash()

	// Añadir el bloque a la cadena
	bc.blocks = append(bc.blocks, b)

	// Limpiar del mempool las txs confirmadas en este bloque
	hashes := make([]crypto.Hash, len(b.Transactions))
	for i, tx := range b.Transactions {
		hashes[i] = tx.ID
	}
	bc.mempool.Flush(hashes)

	return nil
}

// GetBlock devuelve el bloque en la altura indicada, o false si no existe.
func (bc *Blockchain) GetBlock(height uint64) (*block.Block, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if height >= uint64(len(bc.blocks)) {
		return nil, false
	}
	return bc.blocks[height], true
}

// GetBlockByHash busca un bloque por su hash. Devuelve false si no se encuentra.
func (bc *Blockchain) GetBlockByHash(hash crypto.Hash) (*block.Block, bool) {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for _, b := range bc.blocks {
		if b.Hash == hash {
			return b, true
		}
	}
	return nil, false
}

// LatestBlock devuelve el bloque más reciente de la cadena.
func (bc *Blockchain) LatestBlock() *block.Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.blocks[len(bc.blocks)-1]
}

// Height devuelve la altura del último bloque (len(blocks) - 1).
func (bc *Blockchain) Height() uint64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return uint64(len(bc.blocks) - 1)
}

// State devuelve el estado global actual.
func (bc *Blockchain) State() *state.State {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.state
}

// Mempool devuelve el mempool de la blockchain.
func (bc *Blockchain) Mempool() *mempool.Mempool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return bc.mempool
}

// IsValid verifica la integridad completa de la cadena:
// - Cada bloque pasa su propia validación interna
// - Las alturas son secuenciales (0, 1, 2, ...)
// - Cada PrevHash coincide con el hash real del bloque anterior
func (bc *Blockchain) IsValid() bool {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	for i, b := range bc.blocks {
		// Validar integridad interna del bloque
		if err := b.Validate(); err != nil {
			return false
		}

		// El bloque génesis tiene altura 0 y PrevHash cero
		if i == 0 {
			if b.Header.Height != 0 {
				return false
			}
			continue
		}

		// Comprobar altura secuencial
		if b.Header.Height != uint64(i) {
			return false
		}

		// Comprobar que PrevHash apunta al bloque anterior
		prev := bc.blocks[i-1]
		if b.Header.PrevHash != prev.Hash {
			return false
		}
	}
	return true
}
