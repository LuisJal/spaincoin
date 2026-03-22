package consensus

import (
	"errors"
	"sort"
	"sync"

	"github.com/spaincoin/spaincoin/core/crypto"
)

// Validator represents a staking participant.
type Validator struct {
	Address crypto.Address
	Stake   uint64 // amount staked in pesetas
	PubKey  *crypto.PublicKey
}

// ValidatorSet is the active set of validators.
type ValidatorSet struct {
	validators map[crypto.Address]*Validator
	mu         sync.RWMutex
}

// NewValidatorSet creates an empty ValidatorSet.
func NewValidatorSet() *ValidatorSet {
	return &ValidatorSet{
		validators: make(map[crypto.Address]*Validator),
	}
}

// Add adds a validator to the set. Returns an error if the address already
// exists or if the stake is zero.
func (vs *ValidatorSet) Add(v *Validator) error {
	if v.Stake == 0 {
		return errors.New("stake must be greater than zero")
	}
	vs.mu.Lock()
	defer vs.mu.Unlock()
	if _, exists := vs.validators[v.Address]; exists {
		return errors.New("validator already exists")
	}
	vs.validators[v.Address] = v
	return nil
}

// Remove removes a validator by address. Returns an error if not found.
func (vs *ValidatorSet) Remove(addr crypto.Address) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	if _, exists := vs.validators[addr]; !exists {
		return errors.New("validator not found")
	}
	delete(vs.validators, addr)
	return nil
}

// Get returns the validator for the given address and a boolean indicating
// whether it was found.
func (vs *ValidatorSet) Get(addr crypto.Address) (*Validator, bool) {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	v, ok := vs.validators[addr]
	return v, ok
}

// UpdateStake sets a new stake for the validator identified by addr. Returns
// an error if the validator is not found.
func (vs *ValidatorSet) UpdateStake(addr crypto.Address, newStake uint64) error {
	vs.mu.Lock()
	defer vs.mu.Unlock()
	v, exists := vs.validators[addr]
	if !exists {
		return errors.New("validator not found")
	}
	v.Stake = newStake
	return nil
}

// TotalStake returns the sum of all validators' stakes.
func (vs *ValidatorSet) TotalStake() uint64 {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	var total uint64
	for _, v := range vs.validators {
		total += v.Stake
	}
	return total
}

// All returns all validators sorted by address for deterministic iteration.
func (vs *ValidatorSet) All() []*Validator {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	list := make([]*Validator, 0, len(vs.validators))
	for _, v := range vs.validators {
		list = append(list, v)
	}
	sort.Slice(list, func(i, j int) bool {
		for k := 0; k < 20; k++ {
			if list[i].Address[k] != list[j].Address[k] {
				return list[i].Address[k] < list[j].Address[k]
			}
		}
		return false
	})
	return list
}

// Size returns the number of validators in the set.
func (vs *ValidatorSet) Size() int {
	vs.mu.RLock()
	defer vs.mu.RUnlock()
	return len(vs.validators)
}
