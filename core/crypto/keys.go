package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
)

// PrivateKey representa la clave privada ECDSA (secp256k1)
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// PublicKey representa la clave pública derivada
type PublicKey struct {
	X, Y *big.Int
}

// Signature representa una firma ECDSA
type Signature struct {
	R, S *big.Int
}

// GenerateKeyPair genera un nuevo par de claves ECDSA
func GenerateKeyPair() (*PrivateKey, *PublicKey, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return &PrivateKey{key: priv}, &PublicKey{X: priv.PublicKey.X, Y: priv.PublicKey.Y}, nil
}

// Sign firma un hash con la clave privada
func (pk *PrivateKey) Sign(hash []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, pk.key, hash)
	if err != nil {
		return nil, err
	}
	return &Signature{R: r, S: s}, nil
}

// Verify verifica una firma contra una clave pública y un hash
func (pub *PublicKey) Verify(hash []byte, sig *Signature) bool {
	ecPub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: pub.X, Y: pub.Y}
	return ecdsa.Verify(ecPub, hash, sig.R, sig.S)
}

// ToAddress deriva la dirección $SPC de una clave pública (últimos 20 bytes del hash)
func (pub *PublicKey) ToAddress() Address {
	data := append(pub.X.Bytes(), pub.Y.Bytes()...)
	hash := sha256.Sum256(data)
	var addr Address
	copy(addr[:], hash[12:]) // últimos 20 bytes
	return addr
}

// ToHex serializa la clave privada a hex
func (pk *PrivateKey) ToHex() string {
	return hex.EncodeToString(pk.key.D.Bytes())
}

// PrivateKeyFromHex carga una clave privada desde hex
func PrivateKeyFromHex(hexStr string) (*PrivateKey, *PublicKey, error) {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, nil, err
	}
	curve := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.D = new(big.Int).SetBytes(b)
	priv.PublicKey.Curve = curve
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(b)
	if priv.PublicKey.X == nil {
		return nil, nil, errors.New("clave privada inválida")
	}
	return &PrivateKey{key: priv}, &PublicKey{X: priv.PublicKey.X, Y: priv.PublicKey.Y}, nil
}
