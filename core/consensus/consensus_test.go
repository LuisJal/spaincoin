package consensus

import (
	"testing"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// makeAddress creates a deterministic crypto.Address from a single byte value.
func makeAddress(b byte) crypto.Address {
	var addr crypto.Address
	addr[0] = b
	return addr
}

// makeValidator creates a Validator with the given address byte and stake.
func makeValidator(b byte, stake uint64) *Validator {
	return &Validator{
		Address: makeAddress(b),
		Stake:   stake,
	}
}

// -----------------------------------------------------------------------------
// ValidatorSet tests
// -----------------------------------------------------------------------------

func TestValidatorSet_AddAndGet(t *testing.T) {
	vs := NewValidatorSet()
	v := makeValidator(1, 1000)

	if err := vs.Add(v); err != nil {
		t.Fatalf("unexpected error adding validator: %v", err)
	}

	got, ok := vs.Get(v.Address)
	if !ok {
		t.Fatal("validator not found after add")
	}
	if got.Stake != v.Stake {
		t.Errorf("got stake %d, want %d", got.Stake, v.Stake)
	}
}

func TestValidatorSet_NoDuplicates(t *testing.T) {
	vs := NewValidatorSet()
	v := makeValidator(1, 1000)

	if err := vs.Add(v); err != nil {
		t.Fatalf("first add failed: %v", err)
	}
	if err := vs.Add(v); err == nil {
		t.Fatal("expected error adding duplicate validator, got nil")
	}
}

func TestValidatorSet_TotalStake(t *testing.T) {
	vs := NewValidatorSet()
	_ = vs.Add(makeValidator(1, 1000))
	_ = vs.Add(makeValidator(2, 2000))
	_ = vs.Add(makeValidator(3, 3000))

	want := uint64(6000)
	if got := vs.TotalStake(); got != want {
		t.Errorf("TotalStake() = %d, want %d", got, want)
	}
}

// -----------------------------------------------------------------------------
// PoS tests
// -----------------------------------------------------------------------------

func TestPoS_SelectValidator_Deterministic(t *testing.T) {
	vs := NewValidatorSet()
	_ = vs.Add(makeValidator(1, 1000))
	_ = vs.Add(makeValidator(2, 2000))
	_ = vs.Add(makeValidator(3, 3000))

	pos := NewPoS(vs, 100, 0)

	prevHash := crypto.HashBytes([]byte("test-block"))
	height := uint64(42)

	first, err := pos.SelectValidator(height, prevHash)
	if err != nil {
		t.Fatalf("SelectValidator error: %v", err)
	}

	for i := 0; i < 20; i++ {
		got, err := pos.SelectValidator(height, prevHash)
		if err != nil {
			t.Fatalf("SelectValidator error on iteration %d: %v", i, err)
		}
		if got.Address != first.Address {
			t.Errorf("selection not deterministic: got %v, want %v", got.Address, first.Address)
		}
	}
}

func TestPoS_SelectValidator_WeightedDistribution(t *testing.T) {
	vs := NewValidatorSet()
	v1 := makeValidator(1, 1000)
	v2 := makeValidator(2, 2000) // 2× the stake of v1
	_ = vs.Add(v1)
	_ = vs.Add(v2)

	pos := NewPoS(vs, 100, 0)

	counts := make(map[crypto.Address]int)
	const iterations = 10000

	for i := uint64(0); i < iterations; i++ {
		prevHash := crypto.HashBytes([]byte{byte(i), byte(i >> 8), byte(i >> 16)})
		v, err := pos.SelectValidator(i, prevHash)
		if err != nil {
			t.Fatalf("SelectValidator error: %v", err)
		}
		counts[v.Address]++
	}

	c1 := counts[v1.Address]
	c2 := counts[v2.Address]

	if c1 == 0 || c2 == 0 {
		t.Fatalf("one validator never selected: v1=%d v2=%d", c1, c2)
	}

	// v2 should be selected ~2× as often as v1.
	// Allow 20% margin: ratio should be in [1.6, 2.4].
	ratio := float64(c2) / float64(c1)
	if ratio < 1.6 || ratio > 2.4 {
		t.Errorf("weighted distribution off: v1=%d v2=%d ratio=%.2f (want ~2.0)", c1, c2, ratio)
	}
}

func TestPoS_ValidateBlockProposer(t *testing.T) {
	vs := NewValidatorSet()
	_ = vs.Add(makeValidator(1, 1000))
	_ = vs.Add(makeValidator(2, 2000))

	pos := NewPoS(vs, 100, 0)

	prevHash := crypto.HashBytes([]byte("prev"))
	height := uint64(5)

	correct, err := pos.SelectValidator(height, prevHash)
	if err != nil {
		t.Fatalf("SelectValidator error: %v", err)
	}

	// Build a block with the correct validator.
	validBlock := block.NewBlock(height, prevHash, correct.Address, nil)
	if !pos.ValidateBlockProposer(validBlock, prevHash) {
		t.Error("ValidateBlockProposer returned false for correct proposer")
	}

	// Build a block with the wrong validator (pick any address that is not correct).
	wrongAddr := makeAddress(99) // not in the validator set but used as wrong proposer
	invalidBlock := block.NewBlock(height, prevHash, wrongAddr, nil)
	if pos.ValidateBlockProposer(invalidBlock, prevHash) {
		t.Error("ValidateBlockProposer returned true for wrong proposer")
	}
}

// -----------------------------------------------------------------------------
// Slashing tests
// -----------------------------------------------------------------------------

func TestSlashPenalty_DoubleSign(t *testing.T) {
	stake := uint64(1000)
	penalty := SlashPenalty(SlashDoubleSign, stake)
	want := uint64(500) // 50%
	if penalty != want {
		t.Errorf("DoubleSign penalty = %d, want %d", penalty, want)
	}
}

func TestSlashPenalty_Offline(t *testing.T) {
	stake := uint64(1000)
	penalty := SlashPenalty(SlashOffline, stake)
	want := uint64(10) // 1%
	if penalty != want {
		t.Errorf("Offline penalty = %d, want %d", penalty, want)
	}
}

func TestValidatorSet_Slash_ReducesStake(t *testing.T) {
	vs := NewValidatorSet()
	// Use a stake large enough to survive the slash (above minStakeThreshold after slash).
	// minStakeThreshold = 1000 * 1e18; after 1% slash we need stake*0.99 >= threshold.
	// Use 2 * minStakeThreshold so after 1% slash the remaining stake > threshold.
	stake := 2 * minStakeThreshold
	v := &Validator{Address: makeAddress(1), Stake: stake}
	_ = vs.Add(v)

	event, err := vs.Slash(v.Address, SlashOffline, 10)
	if err != nil {
		t.Fatalf("Slash error: %v", err)
	}

	expectedPenalty := stake / 100
	if event.Penalty != expectedPenalty {
		t.Errorf("Penalty = %d, want %d", event.Penalty, expectedPenalty)
	}

	remaining := stake - expectedPenalty
	got, ok := vs.Get(v.Address)
	if !ok {
		t.Fatal("validator was unexpectedly removed")
	}
	if got.Stake != remaining {
		t.Errorf("remaining stake = %d, want %d", got.Stake, remaining)
	}
}

func TestIsEligible(t *testing.T) {
	minStake := uint64(1000)
	vs := NewValidatorSet()
	pos := NewPoS(vs, 100, minStake)

	// Validator with stake below minStake.
	v := makeValidator(1, minStake-1)
	_ = vs.Add(v)

	if pos.IsEligible(v.Address) {
		t.Error("validator below minStake should not be eligible")
	}

	// Raise stake to exactly minStake — now eligible.
	_ = vs.UpdateStake(v.Address, minStake)
	if !pos.IsEligible(v.Address) {
		t.Error("validator at minStake should be eligible")
	}

	// Non-existent address should not be eligible.
	if pos.IsEligible(makeAddress(99)) {
		t.Error("non-existent address should not be eligible")
	}
}
