package node

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/chain"
	"github.com/spaincoin/spaincoin/core/consensus"
	"github.com/spaincoin/spaincoin/core/crypto"
	"github.com/spaincoin/spaincoin/core/storage"
	"github.com/spaincoin/spaincoin/node/rpc"
)

// Config holds the configuration for a SpainCoin node.
type Config struct {
	// ValidatorAddress is the hex-encoded address of this node's validator.
	// If empty and ValidatorKeyHex is set, the address is derived from the key.
	ValidatorAddress string

	// ValidatorKeyHex is the hex-encoded private key. If empty, the node runs
	// in non-validator (observer) mode.
	ValidatorKeyHex string

	// InitialSupply is the genesis supply in pesetas (1 SPC = 10^18 pesetas).
	// Default: 1_000_000_000_000_000 (1 quadrillion pesetas = 1000 SPC).
	// Fits safely in uint64 (max ~1.84e19).
	InitialSupply uint64

	// BlockTime is the target interval between blocks in seconds (default: 5).
	BlockTime int

	// DataDir is the directory where node data is persisted (currently unused).
	DataDir string

	// RPCPort is the HTTP JSON-RPC port (default: 8545).
	RPCPort int

	// P2PPort is the port for P2P networking (default: 30303).
	P2PPort int

	// LogLevel controls verbosity: "debug", "info", "warn", "error" (default: "info").
	LogLevel string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		InitialSupply: 1_000_000_000_000_000, // 1 quadrillion pesetas = 1000 SPC, fits in uint64
		BlockTime:     5,
		RPCPort:       8545,
		P2PPort:       30303,
		LogLevel:      "info",
		DataDir:       "data",
	}
}

// Node is a running SpainCoin node that participates in the network.
type Node struct {
	config     *Config
	chain      *chain.Blockchain
	consensus  *consensus.PoS
	validators *consensus.ValidatorSet
	privKey    *crypto.PrivateKey
	pubKey     *crypto.PublicKey
	address    crypto.Address
	db         *storage.DB
	rpcServer  *rpc.Server
	stopCh     chan struct{}
	wg         sync.WaitGroup
	logger     *log.Logger
}

// NewNode initialises a new Node from the provided Config.
//
// Steps performed:
//  1. If ValidatorKeyHex is set, load the private key and derive the address.
//  2. Create a ValidatorSet and register self as a validator (stake = InitialSupply/10).
//  3. Create a PoS engine (blockReward = 1_000_000_000_000, minStake = 1_000_000_000_000).
//  4. Create the Blockchain with the genesis block.
//  5. Initialise the logger.
func NewNode(cfg *Config) (*Node, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	logger := log.New(os.Stdout, "[spaincoin] ", log.LstdFlags)

	n := &Node{
		config: cfg,
		stopCh: make(chan struct{}),
		logger: logger,
	}

	// --- Load validator key ---
	if cfg.ValidatorKeyHex != "" {
		priv, pub, err := crypto.PrivateKeyFromHex(cfg.ValidatorKeyHex)
		if err != nil {
			return nil, fmt.Errorf("invalid validator key: %w", err)
		}
		n.privKey = priv
		n.pubKey = pub
		n.address = pub.ToAddress()
	} else if cfg.ValidatorAddress != "" {
		// Observer mode: know the address but have no signing key.
		addr, err := crypto.AddressFromHex(cfg.ValidatorAddress)
		if err != nil {
			return nil, fmt.Errorf("invalid validator address: %w", err)
		}
		n.address = addr
	}

	// --- Build ValidatorSet ---
	vs := consensus.NewValidatorSet()

	if !n.address.IsZero() {
		stake := cfg.InitialSupply / 10
		if stake == 0 {
			stake = 1_000_000_000_000 // fallback minimum
		}
		v := &consensus.Validator{
			Address: n.address,
			Stake:   stake,
			PubKey:  n.pubKey,
		}
		if err := vs.Add(v); err != nil {
			return nil, fmt.Errorf("failed to register validator: %w", err)
		}
	}

	n.validators = vs

	// --- PoS engine ---
	const (
		blockReward = uint64(1_000_000_000_000) // 0.000001 SPC per block
		minStake    = uint64(1_000_000_000_000) // same as blockReward
	)
	n.consensus = consensus.NewPoS(vs, blockReward, minStake)

	// --- Storage (optional) ---
	if cfg.DataDir != "" {
		if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create data dir: %w", err)
		}
		db, err := storage.Open(filepath.Join(cfg.DataDir, "chain.db"))
		if err != nil {
			return nil, fmt.Errorf("failed to open storage: %w", err)
		}
		n.db = db

		storedBlocks, err := db.LoadAllBlocks()
		if err != nil {
			return nil, fmt.Errorf("failed to load blocks from disk: %w", err)
		}

		// Create genesis blockchain first (always needed as the base).
		bc, err := chain.NewBlockchain(n.address, cfg.InitialSupply)
		if err != nil {
			return nil, fmt.Errorf("failed to create blockchain: %w", err)
		}

		if len(storedBlocks) > 0 {
			// Replay stored blocks with height > 0 on top of genesis.
			for _, b := range storedBlocks {
				if b.Header.Height == 0 {
					continue // genesis already created above
				}
				if err := bc.AddBlock(b); err != nil {
					return nil, fmt.Errorf("failed to replay block height=%d: %w", b.Header.Height, err)
				}
			}
			logger.Printf("Loaded %d blocks from disk", len(storedBlocks))
		} else {
			// No stored blocks: persist the genesis block now.
			genesis := bc.LatestBlock()
			if err := db.SaveBlock(genesis); err != nil {
				return nil, fmt.Errorf("failed to save genesis block: %w", err)
			}
		}

		n.chain = bc
	} else {
		// In-memory mode: no persistence.
		bc, err := chain.NewBlockchain(n.address, cfg.InitialSupply)
		if err != nil {
			return nil, fmt.Errorf("failed to create blockchain: %w", err)
		}
		n.chain = bc
	}

	return n, nil
}

// Start launches the block-production loop and the RPC server.
// Returns immediately; use Stop to shut down.
func (n *Node) Start() error {
	// Start RPC server
	addr := fmt.Sprintf(":%d", n.config.RPCPort)
	n.rpcServer = rpc.NewServer(n.chain, addr)
	if err := n.rpcServer.Start(); err != nil {
		return fmt.Errorf("failed to start RPC server: %w", err)
	}

	n.logger.Printf("node started (validator=%s, blockTime=%ds, rpc=%s)",
		n.address.String(), n.config.BlockTime, addr)

	n.wg.Add(1)
	go func() {
		defer n.wg.Done()
		n.runBlockProduction()
	}()

	return nil
}

// Stop signals the node to stop and waits for the block-production goroutine
// to exit cleanly.
func (n *Node) Stop() {
	close(n.stopCh)
	n.wg.Wait()
	if n.rpcServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = n.rpcServer.Stop(ctx)
	}
	n.logger.Printf("node stopped at height %d", n.chain.Height())
	if n.db != nil {
		if err := n.db.Close(); err != nil {
			n.logger.Printf("warning: failed to close storage: %v", err)
		}
	}
}

// runBlockProduction is the main loop that produces blocks when this node is
// the selected validator.
func (n *Node) runBlockProduction() {
	ticker := time.NewTicker(time.Duration(n.config.BlockTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-n.stopCh:
			return
		case <-ticker.C:
			if err := n.tryProduceBlock(); err != nil {
				n.logger.Printf("block production error: %v", err)
			}
		}
	}
}

// tryProduceBlock attempts to produce a block at the next height if this node
// is the selected validator.
func (n *Node) tryProduceBlock() error {
	// We need a private key to sign / propose blocks.
	if n.privKey == nil {
		return nil // observer mode
	}

	latest := n.chain.LatestBlock()
	nextHeight := latest.Header.Height + 1

	selected, err := n.consensus.SelectValidator(nextHeight, latest.Hash)
	if err != nil {
		return fmt.Errorf("SelectValidator: %w", err)
	}

	if selected.Address != n.address {
		// Another validator's turn — nothing to do.
		return nil
	}

	// Collect transactions from the mempool (up to 100 per block).
	mp := n.chain.Mempool()
	txs := mp.SelectTxs(100)

	// Add a coinbase reward transaction for the block producer.
	coinbase := block.NewCoinbaseTx(n.address, n.consensus.BlockReward())
	allTxs := append([]*block.Transaction{coinbase}, txs...)

	newBlock := block.NewBlock(nextHeight, latest.Hash, n.address, allTxs)

	if err := n.chain.AddBlock(newBlock); err != nil {
		return fmt.Errorf("AddBlock height=%d: %w", nextHeight, err)
	}

	if n.db != nil {
		if err := n.db.SaveBlock(newBlock); err != nil {
			n.logger.Printf("warning: failed to persist block height=%d: %v", newBlock.Header.Height, err)
		}
	}

	n.logger.Printf("produced block height=%d hash=%s txs=%d",
		newBlock.Header.Height, newBlock.Hash.String(), len(newBlock.Transactions))

	return nil
}

// Status returns a snapshot of the node's current state.
func (n *Node) Status() map[string]interface{} {
	latest := n.chain.LatestBlock()
	return map[string]interface{}{
		"height":         latest.Header.Height,
		"latestHash":     latest.Hash.String(),
		"totalSupply":    n.chain.State().TotalSupply(),
		"validatorCount": n.validators.Size(),
		"mempoolSize":    n.chain.Mempool().Size(),
	}
}

// SubmitTransaction adds a transaction to the mempool for inclusion in a future
// block.
func (n *Node) SubmitTransaction(tx *block.Transaction) error {
	return n.chain.Mempool().Add(tx)
}
