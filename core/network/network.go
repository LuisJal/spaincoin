package network

import (
	"context"
	"fmt"

	"github.com/spaincoin/spaincoin/core/block"
)

// Network is the top-level coordinator for SpainCoin P2P networking.
// It combines host identity, gossipsub messaging, and optional mDNS discovery.
type Network struct {
	host      *Host
	pubsub    *PubSub
	discovery *Discovery
	config    *NetworkConfig
}

// NewNetwork creates a Host, PubSub, and optionally starts mDNS discovery.
func NewNetwork(ctx context.Context, cfg *NetworkConfig) (*Network, error) {
	h, err := NewHost(cfg)
	if err != nil {
		return nil, fmt.Errorf("network: failed to create host: %w", err)
	}

	ps, err := NewPubSub(ctx, h)
	if err != nil {
		_ = h.Close()
		return nil, fmt.Errorf("network: failed to create pubsub: %w", err)
	}

	return &Network{
		host:   h,
		pubsub: ps,
		config: cfg,
	}, nil
}

// Start starts optional services (currently mDNS discovery when EnableMDNS is true).
func (n *Network) Start(ctx context.Context) error {
	if n.config.EnableMDNS {
		d, err := NewMDNSDiscovery(ctx, n.host)
		if err != nil {
			return fmt.Errorf("network: failed to start mDNS discovery: %w", err)
		}
		n.discovery = d
	}
	return nil
}

// BroadcastBlock publishes a block announcement to the blocks gossipsub topic.
func (n *Network) BroadcastBlock(b *block.Block) error {
	return n.pubsub.PublishBlock(b)
}

// BroadcastTx publishes a transaction announcement to the txs gossipsub topic.
func (n *Network) BroadcastTx(tx *block.Transaction) error {
	return n.pubsub.PublishTx(tx)
}

// ReceiveBlocks returns a channel that emits BlockMessage values as they arrive
// from the network. The goroutine exits when ctx is cancelled.
func (n *Network) ReceiveBlocks(ctx context.Context) (<-chan *BlockMessage, error) {
	ch := make(chan *BlockMessage, 64)
	go func() {
		defer close(ch)
		for {
			msg, err := n.pubsub.NextBlock(ctx)
			if err != nil {
				return
			}
			select {
			case ch <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}

// ReceiveTxs returns a channel that emits TxMessage values as they arrive
// from the network. The goroutine exits when ctx is cancelled.
func (n *Network) ReceiveTxs(ctx context.Context) (<-chan *TxMessage, error) {
	ch := make(chan *TxMessage, 64)
	go func() {
		defer close(ch)
		for {
			msg, err := n.pubsub.NextTx(ctx)
			if err != nil {
				return
			}
			select {
			case ch <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch, nil
}

// Stop closes all network services cleanly.
func (n *Network) Stop() {
	n.pubsub.Close()
	if n.discovery != nil {
		n.discovery.Close()
	}
	_ = n.host.Close()
}

// PeerCount returns the number of currently connected peers.
func (n *Network) PeerCount() int {
	return n.host.PeerCount()
}

// NodeInfo returns a map with the node's ID, listen addresses, and connected peer count.
func (n *Network) NodeInfo() map[string]interface{} {
	return map[string]interface{}{
		"id":    n.host.ID().String(),
		"addrs": n.host.Addrs(),
		"peers": n.host.PeerCount(),
	}
}
