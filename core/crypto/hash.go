package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash representa un hash SHA-256 (32 bytes)
type Hash [32]byte

// String devuelve el hash como string hex
func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

// IsZero comprueba si el hash es el hash vacío
func (h Hash) IsZero() bool {
	return h == Hash{}
}

// HashBytes calcula SHA-256 de un slice de bytes
func HashBytes(data []byte) Hash {
	return sha256.Sum256(data)
}

// DoubleHash calcula SHA-256(SHA-256(data)) — usado en bloques (como Bitcoin)
func DoubleHash(data []byte) Hash {
	first := sha256.Sum256(data)
	return sha256.Sum256(first[:])
}

// HashFromHex parsea un hash desde string hex
func HashFromHex(s string) (Hash, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, err
	}
	var h Hash
	copy(h[:], b)
	return h, nil
}
