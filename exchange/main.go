package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spaincoin/spaincoin/exchange/server"
)

func main() {
	nodeURL := os.Getenv("SPC_NODE_URL")
	if nodeURL == "" {
		nodeURL = "http://localhost:8545"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	srv := server.NewServer(":"+port, nodeURL)

	log.Printf("SpainCoin Exchange API starting on :%s", port)
	log.Printf("Connected to node: %s", nodeURL)

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Stop(ctx)
	log.Println("Exchange API stopped.")
}
