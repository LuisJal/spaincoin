package crypto

import "encoding/hex"

// Address representa una dirección SpainCoin (20 bytes, como Ethereum)
type Address [20]byte

// String devuelve la dirección como string hex con prefijo "SPC"
func (a Address) String() string {
	return "SPC" + hex.EncodeToString(a[:])
}

// IsZero comprueba si la dirección es la dirección vacía
func (a Address) IsZero() bool {
	return a == Address{}
}

// AddressFromHex parsea una dirección desde string hex (con o sin prefijo "SPC")
func AddressFromHex(s string) (Address, error) {
	if len(s) > 3 && s[:3] == "SPC" {
		s = s[3:]
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return Address{}, err
	}
	var addr Address
	copy(addr[:], b)
	return addr, nil
}
