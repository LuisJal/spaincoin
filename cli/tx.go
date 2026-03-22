package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

func runTx(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: spc tx <subcommand>")
		fmt.Fprintln(os.Stderr, "  new --from <privkey> --to <address> --amount <uint64> --nonce <uint64> --fee <uint64>")
		fmt.Fprintln(os.Stderr, "  hash <privkey> <to-address> <amount> <nonce> <fee>")
		os.Exit(1)
	}

	switch args[0] {
	case "new":
		txNew(args[1:])
	case "hash":
		txHash(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown tx subcommand: %s\n", args[0])
		os.Exit(1)
	}
}

func txNew(args []string) {
	fs := flag.NewFlagSet("tx new", flag.ContinueOnError)
	fromHex := fs.String("from", "", "Sender private key (hex)")
	toAddr := fs.String("to", "", "Recipient address (SPC...)")
	amount := fs.Uint64("amount", 0, "Amount in pesetas")
	nonce := fs.Uint64("nonce", 0, "Transaction nonce")
	fee := fs.Uint64("fee", 0, "Transaction fee in pesetas")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *fromHex == "" || *toAddr == "" {
		fmt.Fprintln(os.Stderr, "Error: --from and --to are required")
		fs.Usage()
		os.Exit(1)
	}

	priv, pub, err := crypto.PrivateKeyFromHex(*fromHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading private key: %v\n", err)
		os.Exit(1)
	}

	fromAddress := pub.ToAddress()

	toAddress, err := crypto.AddressFromHex(*toAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing to-address: %v\n", err)
		os.Exit(1)
	}

	tx := block.NewTransaction(fromAddress, toAddress, *amount, *nonce, *fee)
	if err := tx.Sign(priv); err != nil {
		fmt.Fprintf(os.Stderr, "Error signing transaction: %v\n", err)
		os.Exit(1)
	}

	printTx(tx)
}

func txHash(args []string) {
	if len(args) < 5 {
		fmt.Fprintln(os.Stderr, "Usage: spc tx hash <privkey-hex> <to-address> <amount> <nonce> <fee>")
		os.Exit(1)
	}

	privHex := args[0]
	toAddrStr := args[1]

	amount, err := strconv.ParseUint(args[2], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid amount: %v\n", err)
		os.Exit(1)
	}

	nonce, err := strconv.ParseUint(args[3], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid nonce: %v\n", err)
		os.Exit(1)
	}

	fee, err := strconv.ParseUint(args[4], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid fee: %v\n", err)
		os.Exit(1)
	}

	priv, pub, err := crypto.PrivateKeyFromHex(privHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading private key: %v\n", err)
		os.Exit(1)
	}

	fromAddress := pub.ToAddress()

	toAddress, err := crypto.AddressFromHex(toAddrStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing to-address: %v\n", err)
		os.Exit(1)
	}

	tx := block.NewTransaction(fromAddress, toAddress, amount, nonce, fee)
	if err := tx.Sign(priv); err != nil {
		fmt.Fprintf(os.Stderr, "Error signing transaction: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(tx.ID.String())
}

func printTx(tx *block.Transaction) {
	sigR := ""
	sigS := ""
	if tx.Signature != nil {
		sigR = tx.Signature.R.Text(16)
		sigS = tx.Signature.S.Text(16)
	}

	fmt.Println("Transaction")
	fmt.Println("===========")
	fmt.Printf("  ID:        %s\n", tx.ID.String())
	fmt.Printf("  From:      %s\n", tx.From.String())
	fmt.Printf("  To:        %s\n", tx.To.String())
	fmt.Printf("  Amount:    %d\n", tx.Amount)
	fmt.Printf("  Nonce:     %d\n", tx.Nonce)
	fmt.Printf("  Fee:       %d\n", tx.Fee)
	fmt.Printf("  Timestamp: %d\n", tx.Timestamp)
	fmt.Printf("  Sig.R:     %s\n", sigR)
	fmt.Printf("  Sig.S:     %s\n", sigS)
}
