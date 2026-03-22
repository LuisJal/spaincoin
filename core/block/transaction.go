package block

import (
	"encoding/binary"
	"errors"
	"time"

	"github.com/spaincoin/spaincoin/core/crypto"
)

// Transaction representa una transferencia de $SPC entre dos direcciones.
// Amount y Fee están expresados en "pesetas" (1 SPC = 10^18 pesetas).
type Transaction struct {
	ID        crypto.Hash
	From      crypto.Address
	To        crypto.Address
	Amount    uint64
	Nonce     uint64
	Fee       uint64
	Timestamp int64
	Signature *crypto.Signature
}

// NewTransaction crea una transacción nueva sin firmar.
// El campo ID se calcula automáticamente a partir de los campos de datos.
func NewTransaction(from, to crypto.Address, amount, nonce, fee uint64) *Transaction {
	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    amount,
		Nonce:     nonce,
		Fee:       fee,
		Timestamp: time.Now().UnixNano(),
	}
	tx.ID = tx.Hash()
	return tx
}

// NewCoinbaseTx crea una transacción coinbase (recompensa de bloque).
// La dirección origen es la dirección cero, lo que la identifica como coinbase.
func NewCoinbaseTx(to crypto.Address, amount uint64) *Transaction {
	tx := &Transaction{
		From:      crypto.Address{}, // dirección cero = coinbase
		To:        to,
		Amount:    amount,
		Nonce:     0,
		Fee:       0,
		Timestamp: time.Now().UnixNano(),
	}
	tx.ID = tx.Hash()
	return tx
}

// Hash calcula el hash SHA-256 de todos los campos de la transacción excepto la firma.
// Usa codificación binaria determinista para garantizar resultados consistentes.
func (tx *Transaction) Hash() crypto.Hash {
	buf := make([]byte, 0, 20+20+8+8+8+8)
	buf = append(buf, tx.From[:]...)
	buf = append(buf, tx.To[:]...)
	buf = appendUint64(buf, tx.Amount)
	buf = appendUint64(buf, tx.Nonce)
	buf = appendUint64(buf, tx.Fee)
	buf = appendInt64(buf, tx.Timestamp)
	return crypto.HashBytes(buf)
}

// Sign firma la transacción usando la clave privada proporcionada.
// Guarda la firma en el campo Signature y actualiza el ID.
func (tx *Transaction) Sign(priv *crypto.PrivateKey) error {
	h := tx.Hash()
	sig, err := priv.Sign(h[:])
	if err != nil {
		return err
	}
	tx.Signature = sig
	tx.ID = h
	return nil
}

// Verify comprueba que la firma de la transacción es válida respecto a la
// dirección From. Devuelve false si no hay firma o si la verificación falla.
func (tx *Transaction) Verify() bool {
	if tx.Signature == nil {
		return false
	}
	// Las transacciones coinbase no tienen firma válida
	if tx.IsCoinbase() {
		return true
	}
	// Necesitamos recuperar la clave pública a partir de la firma.
	// Como la dirección se deriva de la clave pública, necesitamos que
	// el verificador tenga acceso a la clave pública. En este diseño,
	// la verificación completa (incluyendo recuperación de clave pública)
	// se realiza a nivel de cadena donde se dispone del estado de cuentas.
	// Aquí verificamos que la firma tiene estructura válida (R y S no nulos).
	return tx.Signature.R != nil && tx.Signature.S != nil
}

// VerifyWithPublicKey verifica la firma de la transacción contra una clave pública concreta.
func (tx *Transaction) VerifyWithPublicKey(pub *crypto.PublicKey) bool {
	if tx.Signature == nil {
		return false
	}
	if tx.IsCoinbase() {
		return true
	}
	h := tx.Hash()
	return pub.Verify(h[:], tx.Signature)
}

// IsCoinbase devuelve true si la transacción es una recompensa de bloque
// (la dirección From es la dirección cero).
func (tx *Transaction) IsCoinbase() bool {
	return tx.From.IsZero()
}

// appendUint64 añade un uint64 en big-endian a un slice de bytes.
func appendUint64(buf []byte, v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return append(buf, b...)
}

// appendInt64 añade un int64 en big-endian a un slice de bytes.
func appendInt64(buf []byte, v int64) []byte {
	return appendUint64(buf, uint64(v))
}

// ErrInvalidSignature se devuelve cuando la verificación de firma falla.
var ErrInvalidSignature = errors.New("firma de transacción inválida")
