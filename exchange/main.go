package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.etcd.io/bbolt"

	"github.com/spaincoin/spaincoin/exchange/database"
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

	dataDir := os.Getenv("SPC_DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	if err := os.MkdirAll(dataDir, 0700); err != nil {
		log.Fatalf("create data dir %s: %v", dataDir, err)
	}

	boltDB, err := bbolt.Open(filepath.Join(dataDir, "users.db"), 0600, nil)
	if err != nil {
		log.Fatalf("open users.db: %v", err)
	}
	defer boltDB.Close()

	userDB, err := database.NewUserDB(boltDB)
	if err != nil {
		log.Fatalf("init user db: %v", err)
	}

	srv := server.NewServer(":"+port, nodeURL, userDB)

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
