package tests

import (
	"testing"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/chain"
	"github.com/spaincoin/spaincoin/core/consensus"
	"github.com/spaincoin/spaincoin/core/crypto"
	"github.com/spaincoin/spaincoin/core/mempool"
)

// buildValidBlockOnChain constructs a block with height = chain.Height()+1 whose
// timestamp is strictly greater than the latest block's timestamp.
func buildValidBlockOnChain(t *testing.T, bc *chain.Blockchain, validator crypto.Address, txs []*block.Transaction) *block.Block {
	t.Helper()
	last := bc.LatestBlock()
	prevHash := last.Hash
	height := bc.Height() + 1

	b := block.NewBlock(height, prevHash, validator, txs)
	for b.Header.Timestamp <= last.Header.Timestamp {
		time.Sleep(time.Nanosecond)
		b = block.NewBlock(height, prevHash, validator, txs)
	}
	return b
}

// TestFullBlockchainFlow exercises the complete lifecycle:
// genesis → mempool → block production → state verification.
func TestFullBlockchainFlow(t *testing.T) {
	// 1. Generate 3 key pairs.
	alicePriv, alicePub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair alice: %v", err)
	}
	_, bobPub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair bob: %v", err)
	}
	_, _, err = crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair carol: %v", err)
	}

	aliceAddr := alicePub.ToAddress()
	bobAddr := bobPub.ToAddress()

	// 2. Create blockchain with Alice as genesis validator.
	const initialSupply = uint64(10_000_000_000_000_000)
	bc, err := chain.NewBlockchain(aliceAddr, initialSupply)
	if err != nil {
		t.Fatalf("NewBlockchain: %v", err)
	}

	// 3. Verify genesis block exists at height 0.
	genesis, ok := bc.GetBlock(0)
	if !ok {
		t.Fatal("genesis block not found at height 0")
	}
	if genesis.Header.Height != 0 {
		t.Errorf("genesis block height: got %d, want 0", genesis.Header.Height)
	}

	// 4. Verify Alice's balance equals initial supply.
	aliceBalance := bc.State().GetBalance(aliceAddr)
	if aliceBalance != initialSupply {
		t.Errorf("Alice genesis balance: got %d, want %d", aliceBalance, initialSupply)
	}

	// 5. Add tx to mempool: Alice → Bob, amount = 1_000_000_000_000, fee = 1_000_000, nonce = 0.
	const transferAmount = uint64(1_000_000_000_000)
	const txFee = uint64(1_000_000)
	tx := block.NewTransaction(aliceAddr, bobAddr, transferAmount, 0, txFee)
	if err := tx.Sign(alicePriv); err != nil {
		t.Fatalf("Sign tx: %v", err)
	}

	mp := bc.Mempool()
	if err := mp.Add(tx); err != nil {
		t.Fatalf("Mempool.Add: %v", err)
	}

	// 6. Select validator for height 1 (should be Alice — only validator).
	vs := consensus.NewValidatorSet()
	if err := vs.Add(&consensus.Validator{Address: aliceAddr, Stake: 1000, PubKey: alicePub}); err != nil {
		t.Fatalf("ValidatorSet.Add: %v", err)
	}
	pos := consensus.NewPoS(vs, 100, 0)

	genesisHash := bc.LatestBlock().Hash
	selected, err := pos.SelectValidator(1, genesisHash)
	if err != nil {
		t.Fatalf("SelectValidator: %v", err)
	}
	if selected.Address != aliceAddr {
		t.Errorf("selected validator: got %v, want Alice %v", selected.Address, aliceAddr)
	}

	// 7. Create block 1 with the tx + coinbase reward for Alice.
	const blockReward = uint64(100)
	coinbase := block.NewCoinbaseTx(aliceAddr, blockReward)
	selectedTxs := mp.SelectTxs(100)
	allTxs := append([]*block.Transaction{coinbase}, selectedTxs...)

	blk1 := buildValidBlockOnChain(t, bc, aliceAddr, allTxs)

	// 8. Add block 1 to chain.
	if err := bc.AddBlock(blk1); err != nil {
		t.Fatalf("AddBlock(1): %v", err)
	}

	// 9. Verify height = 1.
	if bc.Height() != 1 {
		t.Errorf("chain height after block 1: got %d, want 1", bc.Height())
	}

	// 9. Verify Bob has balance = transferAmount.
	bobBalance := bc.State().GetBalance(bobAddr)
	if bobBalance != transferAmount {
		t.Errorf("Bob balance after tx: got %d, want %d", bobBalance, transferAmount)
	}

	// 9. Verify Alice's balance decreased (spent amount + fee, gained block reward).
	// Alice started with initialSupply, sent transferAmount + fee, gained blockReward (coinbase).
	expectedAliceBalance := initialSupply - transferAmount - txFee + blockReward
	newAliceBalance := bc.State().GetBalance(aliceAddr)
	if newAliceBalance != expectedAliceBalance {
		t.Errorf("Alice balance after block 1: got %d, want %d", newAliceBalance, expectedAliceBalance)
	}
}

// TestConsensusValidatorSelection verifies weighted PoS selection over 10000 rounds.
func TestConsensusValidatorSelection(t *testing.T) {
	// 1. Create ValidatorSet with 3 validators: stakes 1000, 2000, 7000.
	vs := consensus.NewValidatorSet()

	var addr1, addr2, addr3 crypto.Address
	addr1[0] = 1
	addr2[0] = 2
	addr3[0] = 3

	v1 := &consensus.Validator{Address: addr1, Stake: 1000}
	v2 := &consensus.Validator{Address: addr2, Stake: 2000}
	v3 := &consensus.Validator{Address: addr3, Stake: 7000}

	if err := vs.Add(v1); err != nil {
		t.Fatalf("Add v1: %v", err)
	}
	if err := vs.Add(v2); err != nil {
		t.Fatalf("Add v2: %v", err)
	}
	if err := vs.Add(v3); err != nil {
		t.Fatalf("Add v3: %v", err)
	}

	pos := consensus.NewPoS(vs, 100, 0)

	// 2. Run 10000 selections with different heights/hashes.
	const iterations = 10000
	counts := make(map[crypto.Address]int)
	for i := uint64(0); i < iterations; i++ {
		prevHash := crypto.HashBytes([]byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)})
		v, err := pos.SelectValidator(i, prevHash)
		if err != nil {
			t.Fatalf("SelectValidator at height %d: %v", i, err)
		}
		counts[v.Address]++
	}

	totalStake := uint64(1000 + 2000 + 7000) // 10000

	// 3. Verify distribution is approximately proportional (within 15% of expected).
	type check struct {
		addr  crypto.Address
		stake uint64
		label string
	}
	checks := []check{
		{addr1, 1000, "v1(stake=1000, expected ~10%)"},
		{addr2, 2000, "v2(stake=2000, expected ~20%)"},
		{addr3, 7000, "v3(stake=7000, expected ~70%)"},
	}
	for _, c := range checks {
		got := counts[c.addr]
		expected := float64(iterations) * float64(c.stake) / float64(totalStake)
		tolerance := expected * 0.15
		diff := float64(got) - expected
		if diff < 0 {
			diff = -diff
		}
		if diff > tolerance {
			t.Errorf("%s: got %d selections (expected %.0f ± %.0f)", c.label, got, expected, tolerance)
		}
	}

	// 4. Verify determinism: same inputs always produce same output.
	fixedHash := crypto.HashBytes([]byte("fixed-hash"))
	fixedHeight := uint64(42)
	first, err := pos.SelectValidator(fixedHeight, fixedHash)
	if err != nil {
		t.Fatalf("SelectValidator (determinism check): %v", err)
	}
	for i := 0; i < 20; i++ {
		got, err := pos.SelectValidator(fixedHeight, fixedHash)
		if err != nil {
			t.Fatalf("SelectValidator (determinism iteration %d): %v", i, err)
		}
		if got.Address != first.Address {
			t.Errorf("selection not deterministic on iteration %d: got %v, want %v", i, got.Address, first.Address)
		}
	}
}

// TestMempoolOrdering verifies that SelectTxs returns transactions ordered by
// descending fee.
func TestMempoolOrdering(t *testing.T) {
	mp := mempool.NewMempool(100)

	fees := []uint64{100, 500, 50, 1000, 200}
	var addr crypto.Address
	addr[0] = 1
	var toAddr crypto.Address
	toAddr[0] = 2

	for i, fee := range fees {
		tx := block.NewTransaction(addr, toAddr, 1, uint64(i), fee)
		if err := mp.Add(tx); err != nil {
			t.Fatalf("Mempool.Add tx with fee %d: %v", fee, err)
		}
	}

	// 2. SelectTxs(3) should return the 3 highest fee txs: 1000, 500, 200.
	selected := mp.SelectTxs(3)
	if len(selected) != 3 {
		t.Fatalf("SelectTxs(3) returned %d txs, want 3", len(selected))
	}

	// 3. Verify order is descending by fee.
	expectedFees := []uint64{1000, 500, 200}
	for i, tx := range selected {
		if tx.Fee != expectedFees[i] {
			t.Errorf("selected[%d].Fee = %d, want %d", i, tx.Fee, expectedFees[i])
		}
	}
}

// TestChainIntegrity verifies that tampering with a block's hash causes
// IsValid to return false.
func TestChainIntegrity(t *testing.T) {
	// 1. Build chain with 5 blocks.
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	validatorAddr := pub.ToAddress()

	bc, err := chain.NewBlockchain(validatorAddr, 1_000_000_000)
	if err != nil {
		t.Fatalf("NewBlockchain: %v", err)
	}

	for i := 0; i < 5; i++ {
		b := buildValidBlockOnChain(t, bc, validatorAddr, nil)
		if err := bc.AddBlock(b); err != nil {
			t.Fatalf("AddBlock (height %d): %v", i+1, err)
		}
	}

	if !bc.IsValid() {
		t.Fatal("chain should be valid before tampering")
	}

	// 2. Manually corrupt block 2's hash.
	// GetBlock returns the actual pointer stored in the chain slice.
	blk2, ok := bc.GetBlock(2)
	if !ok {
		t.Fatal("GetBlock(2) returned false")
	}
	blk2.Hash = crypto.HashBytes([]byte("corrupted block 2 hash"))

	// 3. IsValid() should return false.
	if bc.IsValid() {
		t.Error("IsValid() should return false after corrupting block 2 hash")
	}

	// 4. Verify blocks 0 and 1 are still accessible and individually valid.
	blk0, ok := bc.GetBlock(0)
	if !ok {
		t.Fatal("GetBlock(0) returned false")
	}
	if err := blk0.Validate(); err != nil {
		t.Errorf("block 0 individual validation failed: %v", err)
	}

	blk1, ok := bc.GetBlock(1)
	if !ok {
		t.Fatal("GetBlock(1) returned false")
	}
	if err := blk1.Validate(); err != nil {
		t.Errorf("block 1 individual validation failed: %v", err)
	}
}

// TestStateRollback verifies that a block containing an invalid transaction
// is rejected atomically and leaves the state unchanged.
func TestStateRollback(t *testing.T) {
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	validatorAddr := pub.ToAddress()

	// 1. Create chain with some initial balance (1000 pesetas).
	const initialBalance = uint64(1000)
	bc, err := chain.NewBlockchain(validatorAddr, initialBalance)
	if err != nil {
		t.Fatalf("NewBlockchain: %v", err)
	}

	// Record state before the bad block attempt.
	balanceBefore := bc.State().GetBalance(validatorAddr)

	// 2. Build a block with a mix: one coinbase tx + one invalid tx (insufficient balance).
	coinbase := block.NewCoinbaseTx(validatorAddr, 100)

	// Create a recipient address.
	var recipientAddr crypto.Address
	recipientAddr[0] = 42

	// This tx tries to send more than the validator's total balance
	// (initialBalance + coinbase reward) — it should fail.
	// Use 10× the initial balance to guarantee it exceeds any coinbase addition.
	oversizedAmount := initialBalance * 10
	invalidTx := block.NewTransaction(validatorAddr, recipientAddr, oversizedAmount, 0, 0)

	txs := []*block.Transaction{coinbase, invalidTx}
	badBlock := buildValidBlockOnChain(t, bc, validatorAddr, txs)

	// 3. ApplyBlock (via AddBlock) should fail.
	if err := bc.AddBlock(badBlock); err == nil {
		t.Fatal("AddBlock with invalid tx should have returned an error")
	}

	// 4. Verify state is unchanged (rollback worked).
	balanceAfter := bc.State().GetBalance(validatorAddr)
	if balanceAfter != balanceBefore {
		t.Errorf("state was not rolled back: balance before=%d, after=%d", balanceBefore, balanceAfter)
	}

	// Chain height must still be 0.
	if bc.Height() != 0 {
		t.Errorf("chain height after rejected block: got %d, want 0", bc.Height())
	}
}

// TestKeyAndSignatureFlow exercises the full ECDSA sign/verify cycle and
// verifies that tampering with a transaction breaks signature verification.
func TestKeyAndSignatureFlow(t *testing.T) {
	// 1. Generate keypair.
	priv, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}

	from := pub.ToAddress()
	var to crypto.Address
	to[0] = 99

	// 2. Create transaction.
	tx := block.NewTransaction(from, to, 500, 0, 10)

	// 3. Sign with private key.
	if err := tx.Sign(priv); err != nil {
		t.Fatalf("Sign: %v", err)
	}

	// 4. Verify signature passes.
	if !tx.VerifyWithPublicKey(pub) {
		t.Error("signature verification should pass for correctly signed tx")
	}

	// 5. Tamper with amount.
	tx.Amount = 9999999

	// 6. Verify signature now fails (detects tampering).
	// VerifyWithPublicKey recomputes the hash from current fields, so the hash
	// will differ from what was signed.
	if tx.VerifyWithPublicKey(pub) {
		t.Error("signature verification should fail after tampering with amount")
	}
}

// TestMerkleTreeTampering verifies that manually changing a transaction's
// amount in an already-created block makes the block's Validate() fail.
func TestMerkleTreeTampering(t *testing.T) {
	_, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	validatorAddr := pub.ToAddress()

	// 1. Create block with 4 transactions.
	var to crypto.Address
	to[0] = 1
	txs := make([]*block.Transaction, 4)
	for i := range txs {
		txs[i] = block.NewCoinbaseTx(validatorAddr, uint64(100*(i+1)))
		to[0] = byte(i + 1)
	}

	prevHash := crypto.Hash{}
	b := block.NewBlock(1, prevHash, validatorAddr, txs)

	// Sanity check: block should be valid before tampering.
	if err := b.Validate(); err != nil {
		t.Fatalf("block should be valid before tampering: %v", err)
	}

	// 2. Manually change one transaction's amount in the block.
	originalAmount := b.Transactions[1].Amount
	b.Transactions[1].Amount = originalAmount + 999999

	// 3. Recompute merkle root — it should differ from stored root.
	// (We verify this implicitly via block.Validate.)
	// The MerkleRoot in the header is based on original tx hashes; after
	// changing the tx amount the recomputed root will differ.

	// 4. block.Validate() should return an error.
	if err := b.Validate(); err == nil {
		t.Error("block.Validate() should return error after tampering with transaction amount")
	}
}
