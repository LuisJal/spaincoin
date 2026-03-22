package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/spaincoin/spaincoin/core/crypto"
)

func runWallet(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spc wallet <subcommand>")
		fmt.Fprintln(os.Stderr, "  new")
		fmt.Fprintln(os.Stderr, "  address <privkey-hex>")
		fmt.Fprintln(os.Stderr, "  verify <address>")
		os.Exit(1)
	}

	switch args[0] {
	case "new":
		walletNew()
	case "address":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: spc wallet address <privkey-hex>")
			os.Exit(1)
		}
		walletAddress(args[1])
	case "verify":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: spc wallet verify <address>")
			os.Exit(1)
		}
		walletVerify(args[1])
	default:
		fmt.Fprintf(os.Stderr, "Unknown wallet subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func walletNew() {
	priv, pub, err := crypto.GenerateKeyPair()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating key pair: %v\n", err)
		os.Exit(1)
	}

	addr := pub.ToAddress()
	pubHex := hex.EncodeToString(append(pub.X.Bytes(), pub.Y.Bytes()...))

	fmt.Println("New SpainCoin Wallet")
	fmt.Println("====================")
	fmt.Printf("Address:     %s\n", addr.String())
	fmt.Printf("Public Key:  %s\n", pubHex)
	fmt.Printf("Private Key: %s\n", priv.ToHex())
	fmt.Println()
	fmt.Println("WARNING: KEEP YOUR PRIVATE KEY SECRET. IT CANNOT BE RECOVERED IF LOST.")
}

func walletAddress(privHex string) {
	_, pub, err := crypto.PrivateKeyFromHex(privHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading private key: %v\n", err)
		os.Exit(1)
	}

	addr := pub.ToAddress()
	fmt.Printf("Address: %s\n", addr.String())
}

func walletVerify(address string) {
	// Address must start with "SPC" followed by 40 hex chars (20 bytes)
	const prefix = "SPC"
	const hexLen = 40 // 20 bytes * 2

	if len(address) != len(prefix)+hexLen {
		fmt.Printf("INVALID: address has wrong length (got %d, want %d)\n", len(address), len(prefix)+hexLen)
		os.Exit(1)
	}

	if address[:len(prefix)] != prefix {
		fmt.Printf("INVALID: address must start with %q\n", prefix)
		os.Exit(1)
	}

	hexPart := address[len(prefix):]
	_, err := hex.DecodeString(hexPart)
	if err != nil {
		fmt.Printf("INVALID: address contains non-hex characters: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("VALID: %s\n", address)
}
