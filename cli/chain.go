package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

func runChain(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spc chain <subcommand>")
		fmt.Fprintln(os.Stderr, "  genesis --validator <address> --supply <uint64>")
		fmt.Fprintln(os.Stderr, "  validate-key <privkey-hex>")
		os.Exit(1)
	}

	switch args[0] {
	case "genesis":
		chainGenesis(args[1:])
	case "validate-key":
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "Usage: spc chain validate-key <privkey-hex>")
			os.Exit(1)
		}
		chainValidateKey(args[1])
	default:
		fmt.Fprintf(os.Stderr, "Unknown chain subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func chainGenesis(args []string) {
	fs := flag.NewFlagSet("chain genesis", flag.ContinueOnError)
	validatorStr := fs.String("validator", "", "Validator address (SPC...)")
	supply := fs.Uint64("supply", 0, "Initial supply in pesetas")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *validatorStr == "" {
		fmt.Fprintln(os.Stderr, "Error: --validator is required")
		fs.Usage()
		os.Exit(1)
	}

	validatorAddr, err := crypto.AddressFromHex(*validatorStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing validator address: %v\n", err)
		os.Exit(1)
	}

	genesis := block.GenesisBlock(validatorAddr, *supply)

	ts := time.Unix(0, genesis.Header.Timestamp).UTC().Format(time.RFC3339)

	fmt.Println("Genesis Block")
	fmt.Println("=============")
	fmt.Printf("  Height:    %d\n", genesis.Header.Height)
	fmt.Printf("  Hash:      %s\n", genesis.Hash.String())
	fmt.Printf("  Timestamp: %s\n", ts)
	fmt.Printf("  Validator: %s\n", genesis.Header.Validator.String())
	fmt.Printf("  PrevHash:  %s\n", genesis.Header.PrevHash.String())
	fmt.Printf("  MerkleRoot:%s\n", genesis.Header.MerkleRoot.String())
	fmt.Println()
	fmt.Println("Coinbase Transaction")
	fmt.Println("--------------------")
	if len(genesis.Transactions) > 0 {
		coinbase := genesis.Transactions[0]
		fmt.Printf("  ID:        %s\n", coinbase.ID.String())
		fmt.Printf("  To:        %s\n", coinbase.To.String())
		fmt.Printf("  Amount:    %d\n", coinbase.Amount)
	}
}

func chainValidateKey(privHex string) {
	priv, pub, err := crypto.PrivateKeyFromHex(privHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "INVALID: cannot load key: %v\n", err)
		os.Exit(1)
	}
	_ = priv // key loaded successfully

	addr := pub.ToAddress()

	fmt.Println("Key is valid.")
	fmt.Printf("  Address:    %s\n", addr.String())
	fmt.Printf("  PrivateKey: %s\n", priv.ToHex())
}
