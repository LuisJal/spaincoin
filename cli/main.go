package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "wallet":
		runWallet(os.Args[2:])
	case "send":
		runSend(os.Args[2:])
	case "tx":
		runTx(os.Args[2:])
	case "chain":
		runChain(os.Args[2:])
	case "help", "--help", "-h":
		printHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}
