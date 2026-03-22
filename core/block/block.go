package block

import (
	"errors"
	"time"

	"github.com/spaincoin/spaincoin/core/crypto"
)

// Header contiene los metadatos del bloque. Su hash identifica unívocamente al bloque.
type Header struct {
	Height     uint64
	PrevHash   crypto.Hash
	MerkleRoot crypto.Hash
	Timestamp  int64
	Validator  crypto.Address
	StateRoot  crypto.Hash
	Nonce      uint64
}

// Block es la unidad fundamental de la blockchain SpainCoin.
// Contiene una cabecera y la lista de transacciones incluidas.
type Block struct {
	Header       *Header
	Transactions []*Transaction
	Hash         crypto.Hash
}

// NewBlock crea un nuevo bloque a partir de los datos proporcionados.
// Calcula automáticamente el MerkleRoot a partir de las transacciones
// y el hash del bloque a partir de la cabecera.
func NewBlock(height uint64, prevHash crypto.Hash, validator crypto.Address, txs []*Transaction) *Block {
	merkleRoot := computeMerkleRoot(txs)

	header := &Header{
		Height:     height,
		PrevHash:   prevHash,
		MerkleRoot: merkleRoot,
		Timestamp:  time.Now().UnixNano(),
		Validator:  validator,
		StateRoot:  crypto.Hash{}, // se establece cuando se aplica el bloque al estado
		Nonce:      0,
	}

	b := &Block{
		Header:       header,
		Transactions: txs,
	}
	b.Hash = b.ComputeHash()
	return b
}

// ComputeHash calcula el hash SHA-256 de la cabecera del bloque.
// Usa codificación binaria determinista.
func (b *Block) ComputeHash() crypto.Hash {
	h := b.Header
	buf := make([]byte, 0, 8+32+32+8+20+32+8)
	buf = appendUint64(buf, h.Height)
	buf = append(buf, h.PrevHash[:]...)
	buf = append(buf, h.MerkleRoot[:]...)
	buf = appendInt64(buf, h.Timestamp)
	buf = append(buf, h.Validator[:]...)
	buf = append(buf, h.StateRoot[:]...)
	buf = appendUint64(buf, h.Nonce)
	return crypto.HashBytes(buf)
}

// Validate verifica la integridad del bloque:
// - El hash almacenado coincide con el hash calculado de la cabecera
// - El MerkleRoot coincide con las transacciones actuales
// - El Timestamp es positivo (mayor que cero)
func (b *Block) Validate() error {
	if b.Header == nil {
		return errors.New("bloque sin cabecera")
	}

	if b.Header.Timestamp <= 0 {
		return errors.New("timestamp del bloque inválido")
	}

	// Verificar que el hash almacenado es correcto
	computed := b.ComputeHash()
	if computed != b.Hash {
		return errors.New("hash del bloque incorrecto")
	}

	// Verificar que el MerkleRoot coincide con las transacciones
	expectedMerkle := computeMerkleRoot(b.Transactions)
	if expectedMerkle != b.Header.MerkleRoot {
		return errors.New("merkle root incorrecto")
	}

	return nil
}

// GenesisBlock crea el bloque génesis (altura 0) de la blockchain SpainCoin.
// Contiene una transacción coinbase que acredita el suministro inicial a validatorAddr.
// El hash previo es el hash cero, ya que no hay bloque anterior.
func GenesisBlock(validatorAddr crypto.Address, initialSupply uint64) *Block {
	coinbase := NewCoinbaseTx(validatorAddr, initialSupply)

	merkleRoot := computeMerkleRoot([]*Transaction{coinbase})

	header := &Header{
		Height:     0,
		PrevHash:   crypto.Hash{}, // hash cero — no hay bloque anterior
		MerkleRoot: merkleRoot,
		Timestamp:  time.Now().UnixNano(),
		Validator:  validatorAddr,
		StateRoot:  crypto.Hash{},
		Nonce:      0,
	}

	b := &Block{
		Header:       header,
		Transactions: []*Transaction{coinbase},
	}
	b.Hash = b.ComputeHash()
	return b
}

// computeMerkleRoot calcula el MerkleRoot de una lista de transacciones.
func computeMerkleRoot(txs []*Transaction) crypto.Hash {
	hashes := make([]crypto.Hash, len(txs))
	for i, tx := range txs {
		hashes[i] = tx.Hash()
	}
	tree := NewMerkleTree(hashes)
	return tree.RootHash()
}
