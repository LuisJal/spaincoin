package consensus

import (
	"encoding/binary"
	"errors"
	"math/rand"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// PoS implements a simple Proof of Stake validator selector.
type PoS struct {
	validators  *ValidatorSet
	blockReward uint64 // pesetas awarded per block
	minStake    uint64 // minimum stake to be eligible as validator
}

// NewPoS creates a new PoS instance.
func NewPoS(validators *ValidatorSet, blockReward, minStake uint64) *PoS {
	return &PoS{
		validators:  validators,
		blockReward: blockReward,
		minStake:    minStake,
	}
}

// SelectValidator picks a validator for the given block height using weighted
// random selection. The random source is seeded deterministically from the
// SHA-256 of (prevBlockHash bytes || blockHeight as big-endian uint64), so
// the same inputs always produce the same result.
//
// A validator's probability of being selected equals stake/totalStake.
func (p *PoS) SelectValidator(blockHeight uint64, prevBlockHash crypto.Hash) (*Validator, error) {
	all := p.validators.All()
	if len(all) == 0 {
		return nil, errors.New("no validators available")
	}

	totalStake := p.validators.TotalStake()
	if totalStake == 0 {
		return nil, errors.New("total stake is zero")
	}

	// Build deterministic seed: SHA-256(prevBlockHash || blockHeight)
	seedBuf := make([]byte, 32+8)
	copy(seedBuf[:32], prevBlockHash[:])
	binary.BigEndian.PutUint64(seedBuf[32:], blockHeight)
	seedHash := crypto.HashBytes(seedBuf)

	// Use the first 8 bytes of the hash as an int64 seed.
	seed := int64(binary.BigEndian.Uint64(seedHash[:8]))
	//nolint:gosec // deterministic pseudo-random for consensus, not security
	rng := rand.New(rand.NewSource(seed))

	// Pick a random point in [0, totalStake).
	point := rng.Uint64() % totalStake

	// Walk cumulative weights to find the selected validator.
	var cumulative uint64
	for _, v := range all {
		cumulative += v.Stake
		if point < cumulative {
			return v, nil
		}
	}

	// Fallback — should not happen if totalStake is consistent.
	return all[len(all)-1], nil
}

// ValidateBlockProposer returns true if the block's Header.Validator is the
// validator that SelectValidator would have chosen for that block.
func (p *PoS) ValidateBlockProposer(b *block.Block, prevBlockHash crypto.Hash) bool {
	selected, err := p.SelectValidator(b.Header.Height, prevBlockHash)
	if err != nil {
		return false
	}
	return selected.Address == b.Header.Validator
}

// BlockReward returns the number of pesetas awarded per block.
func (p *PoS) BlockReward() uint64 {
	return p.blockReward
}

// MinStake returns the minimum stake required to be eligible.
func (p *PoS) MinStake() uint64 {
	return p.minStake
}

// IsEligible returns true if a validator with the given address exists in the
// set and has a stake greater than or equal to minStake.
func (p *PoS) IsEligible(addr crypto.Address) bool {
	v, ok := p.validators.Get(addr)
	if !ok {
		return false
	}
	return v.Stake >= p.minStake
}
