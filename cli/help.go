package main

import "fmt"

func printHelp() {
	fmt.Print(`SpainCoin CLI — $SPC

Usage: spc <command> [subcommand] [flags]

Commands:
  wallet new                           Generate a new wallet
  wallet address <privkey>             Get address from private key
  wallet verify <address>              Verify address format

  tx new --from --to --amount --nonce --fee   Create and sign a transaction
  tx hash <from> <to> <amount> <nonce> <fee>  Print transaction hash

  chain genesis --validator --supply    Create genesis block
  chain validate-key <privkey>          Validate a private key

Examples:
  spc wallet new
  spc wallet address abc123...
  spc tx new --from abc123 --to SPC456... --amount 1000000000000 --nonce 0 --fee 1000000
`)
}
