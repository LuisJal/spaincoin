package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/spaincoin/spaincoin/node"
)

func main() {
	cfg := node.DefaultConfig()

	// Override from environment variables when present.
	if v := os.Getenv("SPC_VALIDATOR_KEY"); v != "" {
		cfg.ValidatorKeyHex = v
	}
	if v := os.Getenv("SPC_VALIDATOR_ADDRESS"); v != "" {
		cfg.ValidatorAddress = v
	}
	if v := os.Getenv("SPC_RPC_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid SPC_RPC_PORT %q: %v", v, err)
		}
		cfg.RPCPort = port
	}
	if v := os.Getenv("SPC_P2P_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid SPC_P2P_PORT %q: %v", v, err)
		}
		cfg.P2PPort = port
	}
	if v := os.Getenv("SPC_BLOCK_TIME"); v != "" {
		bt, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid SPC_BLOCK_TIME %q: %v", v, err)
		}
		cfg.BlockTime = bt
	}
	if v := os.Getenv("SPC_DATA_DIR"); v != "" {
		cfg.DataDir = v
	}
	if v := os.Getenv("SPC_LOG_LEVEL"); v != "" {
		cfg.LogLevel = v
	}

	n, err := node.NewNode(cfg)
	if err != nil {
		log.Fatalf("failed to create node: %v", err)
	}

	status := n.Status()
	log.Printf("SpainCoin node starting...")
	log.Printf("chain height: %v", status["height"])
	log.Printf("total supply: %v pesetas", status["totalSupply"])
	log.Printf("validators:   %v", status["validatorCount"])

	if err := n.Start(); err != nil {
		log.Fatalf("failed to start node: %v", err)
	}

	// Block until SIGINT or SIGTERM.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Printf("shutting down...")
	n.Stop()
	log.Printf("node stopped.")
}
