package wallet

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

// TestNewWallet verifies that New() produces a wallet with a non-zero address.
func TestNewWallet(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if w.Address().IsZero() {
		t.Fatal("expected non-zero address, got zero address")
	}
	if w.PublicKey == nil {
		t.Fatal("expected non-nil PublicKey")
	}
}

// TestWalletSaveLoad saves a wallet to a temp file, loads it back, and checks
// that the address matches and signing works.
func TestWalletSaveLoad(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "test_wallet.json")
	if err := w.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.Address() != w.Address() {
		t.Errorf("address mismatch: got %s, want %s", loaded.Address(), w.Address())
	}

	// Verify the loaded wallet can sign data.
	sig, err := loaded.Sign([]byte("hello spaincoin"))
	if err != nil {
		t.Fatalf("Sign() after Load error: %v", err)
	}
	if sig == nil {
		t.Fatal("expected non-nil signature after Load")
	}
}

// TestWalletSign signs data and verifies the signature against the public key.
// Sign() hashes the payload with SHA-256 before signing, so we reproduce the
// same digest here for verification.
func TestWalletSign(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	data := []byte("test transaction data")
	sig, err := w.Sign(data)
	if err != nil {
		t.Fatalf("Sign() error: %v", err)
	}
	if sig == nil {
		t.Fatal("expected non-nil signature")
	}

	digest := sha256.Sum256(data)
	if !w.PublicKey.Verify(digest[:], sig) {
		t.Fatal("signature verification failed")
	}
}

// TestWalletFile_Permissions checks that the saved wallet file has 0600 mode.
func TestWalletFile_Permissions(t *testing.T) {
	w, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	path := filepath.Join(t.TempDir(), "perm_wallet.json")
	if err := w.Save(path); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}

	perm := info.Mode().Perm()
	if perm != 0600 {
		t.Errorf("expected file permissions 0600, got %04o", perm)
	}
}
