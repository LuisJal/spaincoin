package network

import (
	"context"
	"testing"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

// makeTestConfig returns a NetworkConfig that listens on a random OS-assigned port
// (port 0) and disables mDNS to keep unit tests fast and isolated.
func makeTestConfig() *NetworkConfig {
	return &NetworkConfig{
		ListenAddr:     "/ip4/127.0.0.1/tcp/0",
		EnableMDNS:     false,
		BootstrapPeers: []string{},
	}
}

// TestNewHost verifies that a new host is created with a non-empty ID and at
// least one listen address.
func TestNewHost(t *testing.T) {
	h, err := NewHost(makeTestConfig())
	if err != nil {
		t.Fatalf("NewHost failed: %v", err)
	}
	defer h.Close()

	if h.ID() == "" {
		t.Error("host ID should not be empty")
	}

	addrs := h.Addrs()
	if len(addrs) == 0 {
		t.Error("host should have at least one listen address")
	}
	t.Logf("Host ID: %s", h.ID())
	t.Logf("Host addrs: %v", addrs)
}

// TestHostConnect creates two hosts and connects them, then verifies that each
// sees the other as a connected peer.
func TestHostConnect(t *testing.T) {
	h1, err := NewHost(makeTestConfig())
	if err != nil {
		t.Fatalf("NewHost h1 failed: %v", err)
	}
	defer h1.Close()

	h2, err := NewHost(makeTestConfig())
	if err != nil {
		t.Fatalf("NewHost h2 failed: %v", err)
	}
	defer h2.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Connect h2 -> h1 using h1's first full multiaddr.
	addrs := h1.Addrs()
	if len(addrs) == 0 {
		t.Fatal("h1 has no addresses to connect to")
	}
	if err := h2.Connect(ctx, addrs[0]); err != nil {
		t.Fatalf("Connect failed: %v", err)
	}

	if h2.PeerCount() != 1 {
		t.Errorf("h2 expected 1 peer, got %d", h2.PeerCount())
	}
	if h1.PeerCount() != 1 {
		t.Errorf("h1 expected 1 peer, got %d", h1.PeerCount())
	}
}

// TestPubSub_BlockBroadcast creates two networked nodes, connects them,
// and verifies that a block published by node1 is received by node2.
func TestPubSub_BlockBroadcast(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h1, err := NewHost(makeTestConfig())
	if err != nil {
		t.Fatalf("NewHost h1: %v", err)
	}
	defer h1.Close()

	h2, err := NewHost(makeTestConfig())
	if err != nil {
		t.Fatalf("NewHost h2: %v", err)
	}
	defer h2.Close()

	ps1, err := NewPubSub(ctx, h1)
	if err != nil {
		t.Fatalf("NewPubSub ps1: %v", err)
	}
	defer ps1.Close()

	ps2, err := NewPubSub(ctx, h2)
	if err != nil {
		t.Fatalf("NewPubSub ps2: %v", err)
	}
	defer ps2.Close()

	// Connect the two hosts.
	addrs := h1.Addrs()
	if err := h2.Connect(ctx, addrs[0]); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Give gossipsub time to propagate mesh membership.
	time.Sleep(300 * time.Millisecond)

	// Build a minimal block to broadcast.
	var validator crypto.Address
	copy(validator[:], []byte("testvalidator------"))
	b := block.NewBlock(42, crypto.Hash{}, validator, nil)

	// Publish from node1.
	if err := ps1.PublishBlock(b); err != nil {
		t.Fatalf("PublishBlock: %v", err)
	}

	// Receive on node2 (skip self-messages until we get one from a different peer).
	recvCtx, recvCancel := context.WithTimeout(ctx, 5*time.Second)
	defer recvCancel()

	for {
		msg, err := ps2.NextBlock(recvCtx)
		if err != nil {
			t.Fatalf("NextBlock timed out or errored: %v", err)
		}
		if msg.Height == 42 {
			t.Logf("Received block message: height=%d hash=%s", msg.Height, msg.Hash)
			break
		}
	}
}

// TestPubSub_TxBroadcast creates two networked nodes, connects them,
// and verifies that a transaction published by node1 is received by node2.
func TestPubSub_TxBroadcast(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h1, err := NewHost(makeTestConfig())
	if err != nil {
		t.Fatalf("NewHost h1: %v", err)
	}
	defer h1.Close()

	h2, err := NewHost(makeTestConfig())
	if err != nil {
		t.Fatalf("NewHost h2: %v", err)
	}
	defer h2.Close()

	ps1, err := NewPubSub(ctx, h1)
	if err != nil {
		t.Fatalf("NewPubSub ps1: %v", err)
	}
	defer ps1.Close()

	ps2, err := NewPubSub(ctx, h2)
	if err != nil {
		t.Fatalf("NewPubSub ps2: %v", err)
	}
	defer ps2.Close()

	// Connect the two hosts.
	addrs := h1.Addrs()
	if err := h2.Connect(ctx, addrs[0]); err != nil {
		t.Fatalf("Connect: %v", err)
	}

	// Give gossipsub time to propagate mesh membership.
	time.Sleep(300 * time.Millisecond)

	// Build a minimal transaction to broadcast.
	var from, to crypto.Address
	copy(from[:], []byte("sender--------------"))
	copy(to[:], []byte("receiver------------"))
	tx := block.NewTransaction(from, to, 1000, 1, 10)

	// Publish from node1.
	if err := ps1.PublishTx(tx); err != nil {
		t.Fatalf("PublishTx: %v", err)
	}

	// Receive on node2.
	recvCtx, recvCancel := context.WithTimeout(ctx, 5*time.Second)
	defer recvCancel()

	for {
		msg, err := ps2.NextTx(recvCtx)
		if err != nil {
			t.Fatalf("NextTx timed out or errored: %v", err)
		}
		if msg.Amount == 1000 && msg.Nonce == 1 {
			t.Logf("Received tx message: id=%s from=%s amount=%d", msg.ID, msg.From, msg.Amount)
			break
		}
	}
}

// TestNetwork_StartStop verifies that NewNetwork, Start, and Stop complete
// without panicking.
func TestNetwork_StartStop(t *testing.T) {
	ctx := context.Background()
	cfg := makeTestConfig()

	n, err := NewNetwork(ctx, cfg)
	if err != nil {
		t.Fatalf("NewNetwork failed: %v", err)
	}

	if err := n.Start(ctx); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	info := n.NodeInfo()
	if info["id"] == "" {
		t.Error("NodeInfo id should not be empty")
	}
	t.Logf("NodeInfo: %v", info)

	n.Stop()
}
