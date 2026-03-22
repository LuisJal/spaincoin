package crypto

import (
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	priv, pub, err := GenerateKeyPair()
	if err != nil {
		t.Fatalf("error generando keys: %v", err)
	}
	if priv == nil || pub == nil {
		t.Fatal("keys nil")
	}
}

func TestSignAndVerify(t *testing.T) {
	priv, pub, _ := GenerateKeyPair()
	hash := HashBytes([]byte("hola spaincoin"))

	sig, err := priv.Sign(hash[:])
	if err != nil {
		t.Fatalf("error firmando: %v", err)
	}

	if !pub.Verify(hash[:], sig) {
		t.Fatal("verificación fallida con clave correcta")
	}
}

func TestVerifyFails_WrongKey(t *testing.T) {
	priv, _, _ := GenerateKeyPair()
	_, pub2, _ := GenerateKeyPair()

	hash := HashBytes([]byte("hola spaincoin"))
	sig, _ := priv.Sign(hash[:])

	if pub2.Verify(hash[:], sig) {
		t.Fatal("verificación debería fallar con clave incorrecta")
	}
}

func TestAddress(t *testing.T) {
	_, pub, _ := GenerateKeyPair()
	addr := pub.ToAddress()

	if addr.IsZero() {
		t.Fatal("dirección no debería ser cero")
	}
	if len(addr.String()) != 43 { // "SPC" + 40 hex chars
		t.Fatalf("formato de dirección incorrecto: %s", addr.String())
	}
}

func TestPrivateKeyHexRoundtrip(t *testing.T) {
	priv, pub, _ := GenerateKeyPair()
	hexStr := priv.ToHex()

	priv2, pub2, err := PrivateKeyFromHex(hexStr)
	if err != nil {
		t.Fatalf("error cargando key desde hex: %v", err)
	}
	_ = priv2

	if pub.X.Cmp(pub2.X) != 0 || pub.Y.Cmp(pub2.Y) != 0 {
		t.Fatal("claves públicas no coinciden tras roundtrip")
	}
}
