package storage

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
	bolt "go.etcd.io/bbolt"
)

var (
	bucketBlocks = []byte("blocks")
	bucketHashes = []byte("hashes")
	bucketMeta   = []byte("meta")
	keyHeight    = []byte("height")
)

// storedBlock es la representación JSON de un bloque para persistencia.
type storedBlock struct {
	Height       uint64     `json:"height"`
	PrevHash     string     `json:"prev_hash"`
	Hash         string     `json:"hash"`
	MerkleRoot   string     `json:"merkle_root"`
	StateRoot    string     `json:"state_root"`
	Timestamp    int64      `json:"timestamp"`
	Validator    string     `json:"validator"`
	Nonce        uint64     `json:"nonce"`
	Transactions []storedTx `json:"transactions"`
}

// storedTx es la representación JSON de una transacción para persistencia.
type storedTx struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    uint64 `json:"amount"`
	Nonce     uint64 `json:"nonce"`
	Fee       uint64 `json:"fee"`
	Timestamp int64  `json:"timestamp"`
	SigR      string `json:"sig_r,omitempty"`
	SigS      string `json:"sig_s,omitempty"`
}

// DB es el handle principal de la base de datos persistente de SpainCoin.
type DB struct {
	db *bolt.DB
}

// Open abre (o crea) la base de datos en la ruta indicada y crea los buckets necesarios.
func Open(path string) (*DB, error) {
	bdb, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("storage: no se pudo abrir la base de datos: %w", err)
	}

	// Crear buckets si no existen
	err = bdb.Update(func(tx *bolt.Tx) error {
		for _, name := range [][]byte{bucketBlocks, bucketHashes, bucketMeta} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("storage: no se pudo crear bucket %s: %w", name, err)
			}
		}
		return nil
	})
	if err != nil {
		bdb.Close()
		return nil, err
	}

	return &DB{db: bdb}, nil
}

// Close cierra la base de datos.
func (d *DB) Close() error {
	return d.db.Close()
}

// SaveBlock guarda un bloque indexado por altura y por hash.
// Actualiza también la meta-clave "height" si es el bloque más alto visto.
func (d *DB) SaveBlock(b *block.Block) error {
	sb := blockToStored(b)
	data, err := json.Marshal(sb)
	if err != nil {
		return fmt.Errorf("storage: no se pudo serializar bloque: %w", err)
	}

	heightKey := encodeHeight(b.Header.Height)
	hashKey := []byte(b.Hash.String())

	return d.db.Update(func(tx *bolt.Tx) error {
		// Guardar en blocks bucket
		if err := tx.Bucket(bucketBlocks).Put(heightKey, data); err != nil {
			return err
		}

		// Guardar índice hash → height
		if err := tx.Bucket(bucketHashes).Put(hashKey, heightKey); err != nil {
			return err
		}

		// Actualizar meta height si es mayor que el actual
		meta := tx.Bucket(bucketMeta)
		current := meta.Get(keyHeight)
		if current == nil || b.Header.Height > binary.BigEndian.Uint64(current) {
			if err := meta.Put(keyHeight, heightKey); err != nil {
				return err
			}
		}

		return nil
	})
}

// GetBlockByHeight recupera un bloque por su altura.
// Devuelve error si no existe.
func (d *DB) GetBlockByHeight(height uint64) (*block.Block, error) {
	var b *block.Block

	err := d.db.View(func(tx *bolt.Tx) error {
		data := tx.Bucket(bucketBlocks).Get(encodeHeight(height))
		if data == nil {
			return fmt.Errorf("storage: bloque en altura %d no encontrado", height)
		}
		var sb storedBlock
		if err := json.Unmarshal(data, &sb); err != nil {
			return fmt.Errorf("storage: no se pudo deserializar bloque: %w", err)
		}
		var err error
		b, err = storedToBlock(&sb)
		return err
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetBlockByHash recupera un bloque por su hash hex.
// Devuelve error si no existe.
func (d *DB) GetBlockByHash(hashHex string) (*block.Block, error) {
	var b *block.Block

	err := d.db.View(func(tx *bolt.Tx) error {
		heightBytes := tx.Bucket(bucketHashes).Get([]byte(hashHex))
		if heightBytes == nil {
			return fmt.Errorf("storage: bloque con hash %s no encontrado", hashHex)
		}
		height := binary.BigEndian.Uint64(heightBytes)
		data := tx.Bucket(bucketBlocks).Get(encodeHeight(height))
		if data == nil {
			return fmt.Errorf("storage: datos de bloque en altura %d no encontrados", height)
		}
		var sb storedBlock
		if err := json.Unmarshal(data, &sb); err != nil {
			return fmt.Errorf("storage: no se pudo deserializar bloque: %w", err)
		}
		var err error
		b, err = storedToBlock(&sb)
		return err
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

// GetHeight devuelve la mayor altura almacenada. Devuelve -1 si la base de datos está vacía.
func (d *DB) GetHeight() (int64, error) {
	var height int64 = -1

	err := d.db.View(func(tx *bolt.Tx) error {
		v := tx.Bucket(bucketMeta).Get(keyHeight)
		if v == nil {
			return nil
		}
		height = int64(binary.BigEndian.Uint64(v))
		return nil
	})
	return height, err
}

// LoadAllBlocks devuelve todos los bloques ordenados por altura ascendente.
// Útil para reconstruir la cadena al arrancar el nodo.
func (d *DB) LoadAllBlocks() ([]*block.Block, error) {
	var blocks []*block.Block

	err := d.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucketBlocks).Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var sb storedBlock
			if err := json.Unmarshal(v, &sb); err != nil {
				return fmt.Errorf("storage: no se pudo deserializar bloque: %w", err)
			}
			b, err := storedToBlock(&sb)
			if err != nil {
				return err
			}
			blocks = append(blocks, b)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

// ── Helpers de conversión ────────────────────────────────────────────────────

// encodeHeight codifica un uint64 como 8 bytes big-endian.
// Esto garantiza que las claves BoltDB se iteren en orden de altura.
func encodeHeight(h uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, h)
	return b
}

// blockToStored convierte un *block.Block a su representación JSON-friendly.
func blockToStored(b *block.Block) *storedBlock {
	sb := &storedBlock{
		Height:     b.Header.Height,
		PrevHash:   b.Header.PrevHash.String(),
		Hash:       b.Hash.String(),
		MerkleRoot: b.Header.MerkleRoot.String(),
		StateRoot:  b.Header.StateRoot.String(),
		Timestamp:  b.Header.Timestamp,
		Validator:  hex.EncodeToString(b.Header.Validator[:]),
		Nonce:      b.Header.Nonce,
	}

	sb.Transactions = make([]storedTx, len(b.Transactions))
	for i, tx := range b.Transactions {
		st := storedTx{
			ID:        tx.ID.String(),
			From:      hex.EncodeToString(tx.From[:]),
			To:        hex.EncodeToString(tx.To[:]),
			Amount:    tx.Amount,
			Nonce:     tx.Nonce,
			Fee:       tx.Fee,
			Timestamp: tx.Timestamp,
		}
		if tx.Signature != nil {
			st.SigR = tx.Signature.R.Text(16)
			st.SigS = tx.Signature.S.Text(16)
		}
		sb.Transactions[i] = st
	}

	return sb
}

// storedToBlock convierte una storedBlock a *block.Block.
func storedToBlock(sb *storedBlock) (*block.Block, error) {
	prevHash, err := crypto.HashFromHex(sb.PrevHash)
	if err != nil {
		return nil, fmt.Errorf("storage: prev_hash inválido: %w", err)
	}
	blockHash, err := crypto.HashFromHex(sb.Hash)
	if err != nil {
		return nil, fmt.Errorf("storage: hash inválido: %w", err)
	}
	merkleRoot, err := crypto.HashFromHex(sb.MerkleRoot)
	if err != nil {
		return nil, fmt.Errorf("storage: merkle_root inválido: %w", err)
	}
	stateRoot, err := crypto.HashFromHex(sb.StateRoot)
	if err != nil {
		return nil, fmt.Errorf("storage: state_root inválido: %w", err)
	}
	validator, err := crypto.AddressFromHex(sb.Validator)
	if err != nil {
		return nil, fmt.Errorf("storage: validator inválido: %w", err)
	}

	txs := make([]*block.Transaction, len(sb.Transactions))
	for i, st := range sb.Transactions {
		txID, err := crypto.HashFromHex(st.ID)
		if err != nil {
			return nil, fmt.Errorf("storage: tx.ID inválido: %w", err)
		}
		from, err := crypto.AddressFromHex(st.From)
		if err != nil {
			return nil, fmt.Errorf("storage: tx.From inválido: %w", err)
		}
		to, err := crypto.AddressFromHex(st.To)
		if err != nil {
			return nil, fmt.Errorf("storage: tx.To inválido: %w", err)
		}

		tx := &block.Transaction{
			ID:        txID,
			From:      from,
			To:        to,
			Amount:    st.Amount,
			Nonce:     st.Nonce,
			Fee:       st.Fee,
			Timestamp: st.Timestamp,
		}

		if st.SigR != "" && st.SigS != "" {
			r, ok1 := new(big.Int).SetString(st.SigR, 16)
			s, ok2 := new(big.Int).SetString(st.SigS, 16)
			if !ok1 || !ok2 {
				return nil, errors.New("storage: firma inválida en transacción")
			}
			tx.Signature = &crypto.Signature{R: r, S: s}
		}

		txs[i] = tx
	}

	b := &block.Block{
		Header: &block.Header{
			Height:     sb.Height,
			PrevHash:   prevHash,
			MerkleRoot: merkleRoot,
			Timestamp:  sb.Timestamp,
			Validator:  validator,
			StateRoot:  stateRoot,
			Nonce:      sb.Nonce,
		},
		Transactions: txs,
		Hash:         blockHash,
	}

	return b, nil
}
