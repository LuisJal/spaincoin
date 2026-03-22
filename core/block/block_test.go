package block

import (
	"testing"
	"time"

	"github.com/spaincoin/spaincoin/core/crypto"
)

// ---------------------------------------------------------------------------
// Transaction tests
// ---------------------------------------------------------------------------

func TestNewTransaction(t *testing.T) {
	priv, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	from := pub.ToAddress()

	_, pub2, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	to := pub2.ToAddress()

	_ = priv // used in sign tests below

	tx := NewTransaction(from, to, 1_000_000_000_000_000_000, 1, 1000)

	if tx.From != from {
		t.Errorf("From mismatch: got %v, want %v", tx.From, from)
	}
	if tx.To != to {
		t.Errorf("To mismatch: got %v, want %v", tx.To, to)
	}
	if tx.Amount != 1_000_000_000_000_000_000 {
		t.Errorf("Amount mismatch: got %d", tx.Amount)
	}
	if tx.Nonce != 1 {
		t.Errorf("Nonce mismatch: got %d", tx.Nonce)
	}
	if tx.Fee != 1000 {
		t.Errorf("Fee mismatch: got %d", tx.Fee)
	}
	if tx.Timestamp <= 0 {
		t.Errorf("Timestamp should be positive, got %d", tx.Timestamp)
	}
	if tx.ID.IsZero() {
		t.Error("ID should not be zero after NewTransaction")
	}
	if tx.Signature != nil {
		t.Error("Signature should be nil before signing")
	}
}

func TestTransactionHash(t *testing.T) {
	from := crypto.Address{1, 2, 3}
	to := crypto.Address{4, 5, 6}

	tx := &Transaction{
		From:      from,
		To:        to,
		Amount:    500,
		Nonce:     7,
		Fee:       10,
		Timestamp: 1234567890,
	}

	h1 := tx.Hash()
	h2 := tx.Hash()

	if h1 != h2 {
		t.Error("Hash() is not deterministic")
	}
	if h1.IsZero() {
		t.Error("Hash() returned zero hash")
	}

	// Changing a field must change the hash
	tx2 := *tx
	tx2.Amount = 501
	if tx.Hash() == tx2.Hash() {
		t.Error("Different amounts should produce different hashes")
	}
}

func TestTransactionSignVerify(t *testing.T) {
	priv, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	from := pub.ToAddress()

	_, pub2, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	to := pub2.ToAddress()

	tx := NewTransaction(from, to, 100, 1, 5)

	if err := tx.Sign(priv); err != nil {
		t.Fatalf("Sign: %v", err)
	}
	if tx.Signature == nil {
		t.Fatal("Signature should not be nil after Sign")
	}

	// Verify with correct public key — must pass
	if !tx.VerifyWithPublicKey(pub) {
		t.Error("VerifyWithPublicKey should return true with correct key")
	}

	// Verify with wrong key — must fail
	_, wrongPub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	if tx.VerifyWithPublicKey(wrongPub) {
		t.Error("VerifyWithPublicKey should return false with wrong key")
	}
}

func TestCoinbaseTx(t *testing.T) {
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	to := pub.ToAddress()

	// 5 SPC expressed in pesetas (5 * 10^18), within uint64 range
	const fiveSPC = uint64(5_000_000_000_000_000_000)
	tx := NewCoinbaseTx(to, fiveSPC)

	if !tx.IsCoinbase() {
		t.Error("IsCoinbase() should return true for coinbase tx")
	}
	if !tx.From.IsZero() {
		t.Error("Coinbase From address should be zero")
	}
	if tx.To != to {
		t.Errorf("To mismatch: got %v, want %v", tx.To, to)
	}
	if tx.Amount != fiveSPC {
		t.Errorf("Amount mismatch: got %d", tx.Amount)
	}
	if tx.ID.IsZero() {
		t.Error("ID should not be zero")
	}
}

// ---------------------------------------------------------------------------
// MerkleTree tests
// ---------------------------------------------------------------------------

func TestMerkleTree_Empty(t *testing.T) {
	tree := NewMerkleTree(nil)
	if !tree.RootHash().IsZero() {
		t.Error("Empty tree should have zero root hash")
	}

	tree2 := NewMerkleTree([]crypto.Hash{})
	if !tree2.RootHash().IsZero() {
		t.Error("Empty slice should have zero root hash")
	}
}

func TestMerkleTree_Single(t *testing.T) {
	h := crypto.HashBytes([]byte("single tx"))
	tree := NewMerkleTree([]crypto.Hash{h})

	root := tree.RootHash()
	if root.IsZero() {
		t.Error("Single-node tree root should not be zero")
	}
	// A single-node tree has no pairs to combine; the root is the leaf itself.
	if root != h {
		t.Errorf("Single-node root: got %v, want %v", root, h)
	}
}

func TestMerkleTree_Multiple(t *testing.T) {
	hashes := make([]crypto.Hash, 4)
	for i := range hashes {
		hashes[i] = crypto.HashBytes([]byte{byte(i)})
	}

	tree := NewMerkleTree(hashes)
	root := tree.RootHash()

	if root.IsZero() {
		t.Error("Multi-node tree root should not be zero")
	}

	// Build expected manually:
	// Layer 0: h0, h1, h2, h3
	// Layer 1: hash(h0,h1), hash(h2,h3)
	// Layer 2: hash(hash(h0,h1), hash(h2,h3))
	l1a := hashPair(hashes[0], hashes[1])
	l1b := hashPair(hashes[2], hashes[3])
	expected := hashPair(l1a, l1b)

	if root != expected {
		t.Errorf("4-node root: got %v, want %v", root, expected)
	}
}

func TestMerkleTree_OddCount(t *testing.T) {
	hashes := make([]crypto.Hash, 3)
	for i := range hashes {
		hashes[i] = crypto.HashBytes([]byte{byte(i + 10)})
	}

	// Should not panic
	tree := NewMerkleTree(hashes)
	root := tree.RootHash()

	if root.IsZero() {
		t.Error("Odd-count tree root should not be zero")
	}

	// Build expected manually:
	// Layer 0: h0, h1, h2  → pad to h0, h1, h2, h2
	// Layer 1: hash(h0,h1), hash(h2,h2)
	// Layer 2: hash(hash(h0,h1), hash(h2,h2))
	l1a := hashPair(hashes[0], hashes[1])
	l1b := hashPair(hashes[2], hashes[2])
	expected := hashPair(l1a, l1b)

	if root != expected {
		t.Errorf("3-node root: got %v, want %v", root, expected)
	}
}

// ---------------------------------------------------------------------------
// Block tests
// ---------------------------------------------------------------------------

func TestNewBlock(t *testing.T) {
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	validator := pub.ToAddress()

	_, pub2, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	to := pub2.ToAddress()

	txs := []*Transaction{
		NewCoinbaseTx(to, 1000),
	}

	prevHash := crypto.HashBytes([]byte("previous block"))
	b := NewBlock(1, prevHash, validator, txs)

	if b.Header == nil {
		t.Fatal("Block header should not be nil")
	}
	if b.Header.Height != 1 {
		t.Errorf("Height: got %d, want 1", b.Header.Height)
	}
	if b.Header.PrevHash != prevHash {
		t.Errorf("PrevHash mismatch")
	}
	if b.Header.Validator != validator {
		t.Errorf("Validator mismatch")
	}
	if b.Header.Timestamp <= 0 {
		t.Errorf("Timestamp should be positive")
	}
	if b.Hash.IsZero() {
		t.Error("Block hash should not be zero")
	}
	if len(b.Transactions) != 1 {
		t.Errorf("Transactions count: got %d, want 1", len(b.Transactions))
	}
}

func TestBlockHash(t *testing.T) {
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	validator := pub.ToAddress()

	header := &Header{
		Height:     5,
		PrevHash:   crypto.Hash{1, 2, 3},
		MerkleRoot: crypto.Hash{4, 5, 6},
		Timestamp:  time.Now().UnixNano(),
		Validator:  validator,
		StateRoot:  crypto.Hash{},
		Nonce:      0,
	}
	b := &Block{Header: header, Transactions: nil}

	h1 := b.ComputeHash()
	h2 := b.ComputeHash()

	if h1 != h2 {
		t.Error("ComputeHash() is not deterministic")
	}
	if h1.IsZero() {
		t.Error("ComputeHash() should not return zero hash")
	}
}

func TestBlockValidate(t *testing.T) {
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	validator := pub.ToAddress()

	_, pub2, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	to := pub2.ToAddress()

	txs := []*Transaction{NewCoinbaseTx(to, 1000)}
	b := NewBlock(2, crypto.Hash{}, validator, txs)

	// Valid block should pass
	if err := b.Validate(); err != nil {
		t.Errorf("Validate() on valid block failed: %v", err)
	}

	// Tamper with the hash — should fail
	tampered := *b
	tampered.Hash = crypto.HashBytes([]byte("wrong"))
	if err := tampered.Validate(); err == nil {
		t.Error("Validate() should fail when hash is tampered")
	}

	// Tamper with MerkleRoot — should fail
	tampered2 := *b
	header2 := *b.Header
	header2.MerkleRoot = crypto.HashBytes([]byte("wrong merkle"))
	tampered2.Header = &header2
	tampered2.Hash = tampered2.ComputeHash() // recalculate hash so hash check passes
	if err := tampered2.Validate(); err == nil {
		t.Error("Validate() should fail when MerkleRoot is tampered")
	}

	// Invalid timestamp — should fail
	tampered3 := *b
	header3 := *b.Header
	header3.Timestamp = 0
	tampered3.Header = &header3
	tampered3.Hash = tampered3.ComputeHash()
	if err := tampered3.Validate(); err == nil {
		t.Error("Validate() should fail with zero timestamp")
	}
}

func TestGenesisBlock(t *testing.T) {
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	validator := pub.ToAddress()

	// 1 SPC expressed in pesetas (10^18); 1M SPC would overflow uint64
	const initialSupply = uint64(1_000_000_000_000_000_000)

	genesis := GenesisBlock(validator, initialSupply)

	if genesis.Header.Height != 0 {
		t.Errorf("Genesis height: got %d, want 0", genesis.Header.Height)
	}
	if !genesis.Header.PrevHash.IsZero() {
		t.Error("Genesis PrevHash should be zero")
	}
	if len(genesis.Transactions) != 1 {
		t.Fatalf("Genesis should have exactly 1 transaction, got %d", len(genesis.Transactions))
	}

	coinbase := genesis.Transactions[0]
	if !coinbase.IsCoinbase() {
		t.Error("Genesis transaction should be coinbase")
	}
	if coinbase.To != validator {
		t.Errorf("Coinbase To: got %v, want %v", coinbase.To, validator)
	}
	if coinbase.Amount != initialSupply {
		t.Errorf("Coinbase Amount: got %d, want %d", coinbase.Amount, initialSupply)
	}

	// Genesis block must validate correctly
	if err := genesis.Validate(); err != nil {
		t.Errorf("GenesisBlock.Validate() failed: %v", err)
	}
}
