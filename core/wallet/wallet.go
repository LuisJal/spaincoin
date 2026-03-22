package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spaincoin/spaincoin/core/crypto"
)

// Wallet holds an ECDSA keypair and the derived SpainCoin address.
// The private key is unexported; use Sign to produce signatures.
type Wallet struct {
	addr       crypto.Address
	PublicKey  *crypto.PublicKey
	privateKey *crypto.PrivateKey
	filepath   string
}

// WalletFile is the JSON structure stored on disk.
// TODO: in production, PrivateKey should be encrypted with a user-supplied
// password (e.g. AES-256-GCM with an Argon2id-derived key) rather than stored
// as plaintext hex.
type WalletFile struct {
	Address    string `json:"address"`     // "SPC..."
	PublicKeyX string `json:"pub_key_x"`   // hex
	PublicKeyY string `json:"pub_key_y"`   // hex
	PrivateKey string `json:"private_key"` // hex — plaintext for now, see TODO above
	CreatedAt  int64  `json:"created_at"`  // unix timestamp
	Version    string `json:"version"`     // "1.0"
}

// New generates a new ECDSA keypair and returns a Wallet.
// The wallet is not saved to disk; call Save to persist it.
func New() (*Wallet, error) {
	priv, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("wallet: failed to generate keypair: %w", err)
	}
	return &Wallet{
		addr:       pub.ToAddress(),
		PublicKey:  pub,
		privateKey: priv,
	}, nil
}

// Save writes the wallet as JSON to the given filepath with 0600 permissions
// (owner read/write only).
func (w *Wallet) Save(path string) error {
	wf := WalletFile{
		Address:    w.addr.String(),
		PublicKeyX: hex.EncodeToString(w.PublicKey.X.Bytes()),
		PublicKeyY: hex.EncodeToString(w.PublicKey.Y.Bytes()),
		PrivateKey: w.privateKey.ToHex(),
		CreatedAt:  time.Now().Unix(),
		Version:    "1.0",
	}

	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return fmt.Errorf("wallet: failed to marshal wallet file: %w", err)
	}

	// Write with 0600 so that only the owner can read/write the file.
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("wallet: failed to write wallet file: %w", err)
	}

	w.filepath = path
	return nil
}

// Load reads a WalletFile from disk and reconstructs the Wallet,
// including the private key.
func Load(path string) (*Wallet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("wallet: failed to read wallet file: %w", err)
	}

	var wf WalletFile
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("wallet: failed to parse wallet file: %w", err)
	}

	priv, pub, err := crypto.PrivateKeyFromHex(wf.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("wallet: failed to load private key: %w", err)
	}

	addr, err := crypto.AddressFromHex(wf.Address)
	if err != nil {
		return nil, fmt.Errorf("wallet: failed to parse address: %w", err)
	}

	return &Wallet{
		addr:       addr,
		PublicKey:  pub,
		privateKey: priv,
		filepath:   path,
	}, nil
}

// Sign hashes data with SHA-256 and signs the digest with the wallet's private key.
func (w *Wallet) Sign(data []byte) (*crypto.Signature, error) {
	if w.privateKey == nil {
		return nil, fmt.Errorf("wallet: no private key available")
	}
	digest := sha256.Sum256(data)
	return w.privateKey.Sign(digest[:])
}

// Address returns the SpainCoin address derived from the wallet's public key.
func (w *Wallet) Address() crypto.Address {
	return w.addr
}

// String returns a human-readable representation of the wallet.
func (w *Wallet) String() string {
	return fmt.Sprintf("SpainCoin Wallet\nAddress: %s\n", w.addr.String())
}
