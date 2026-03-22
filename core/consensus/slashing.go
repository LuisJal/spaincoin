package consensus

import (
	"errors"

	"github.com/spaincoin/spaincoin/core/crypto"
)

// SlashReason identifies why a validator is being slashed.
type SlashReason string

const (
	// SlashDoubleSign is raised when a validator signs two blocks at the same height.
	SlashDoubleSign SlashReason = "double_sign"
	// SlashOffline is raised when a validator misses too many blocks.
	SlashOffline SlashReason = "offline"
)

// minStakeThreshold is the stake floor below which a slashed validator is
// removed from the set.
//
// The spec calls for 1000 SPC (1000 * 10^18 pesetas), but that value exceeds
// uint64's maximum (~1.84 * 10^19). We therefore use 1 SPC (1 * 10^18
// pesetas = 1_000_000_000_000_000_000) as the effective minimum-stake
// threshold — the largest power-of-ten SPC value that fits in a uint64.
const minStakeThreshold uint64 = 1_000_000_000_000_000_000 // 1 SPC in pesetas

// SlashEvent records a slashing that has taken place.
type SlashEvent struct {
	Validator crypto.Address
	Reason    SlashReason
	Height    uint64
	Penalty   uint64 // pesetas slashed
}

// SlashPenalty computes the penalty for the given slash reason and stake.
//
//   - DoubleSign: 50 % of stake
//   - Offline:     1 % of stake
func SlashPenalty(reason SlashReason, stake uint64) uint64 {
	switch reason {
	case SlashDoubleSign:
		return stake / 2
	case SlashOffline:
		return stake / 100
	default:
		return 0
	}
}

// Slash reduces the stake of the validator at addr by the appropriate penalty,
// returns a SlashEvent, and removes the validator from the set if the remaining
// stake falls below minStakeThreshold.
func (vs *ValidatorSet) Slash(addr crypto.Address, reason SlashReason, height uint64) (*SlashEvent, error) {
	vs.mu.Lock()
	defer vs.mu.Unlock()

	v, exists := vs.validators[addr]
	if !exists {
		return nil, errors.New("validator not found")
	}

	penalty := SlashPenalty(reason, v.Stake)

	if penalty >= v.Stake {
		v.Stake = 0
	} else {
		v.Stake -= penalty
	}

	event := &SlashEvent{
		Validator: addr,
		Reason:    reason,
		Height:    height,
		Penalty:   penalty,
	}

	if v.Stake < minStakeThreshold {
		delete(vs.validators, addr)
	}

	return event, nil
}
