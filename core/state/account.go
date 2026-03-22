package state

import (
	"fmt"

	"github.com/spaincoin/spaincoin/core/crypto"
)

// Account representa una cuenta en el estado global de SpainCoin.
// Cada dirección tiene un saldo (en pesetas) y un nonce para evitar replay attacks.
type Account struct {
	Address crypto.Address
	Balance uint64            // en pesetas (1 SPC = 10^18 pesetas)
	Nonce   uint64            // incrementa con cada tx enviada
	PubKey  *crypto.PublicKey // almacenada en la primera tx enviada
}

// NewAccount crea una cuenta vacía para la dirección dada.
func NewAccount(addr crypto.Address) *Account {
	return &Account{
		Address: addr,
		Balance: 0,
		Nonce:   0,
		PubKey:  nil,
	}
}

// String devuelve un resumen legible de la cuenta.
func (a *Account) String() string {
	return fmt.Sprintf("Account{addr=%s, balance=%d, nonce=%d}", a.Address.String(), a.Balance, a.Nonce)
}
