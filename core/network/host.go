package network

import (
	"context"
	"fmt"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// ProtocolID is the protocol identifier for SpainCoin P2P communication.
const ProtocolID = "/spaincoin/1.0.0"

// NetworkConfig holds configuration for the P2P host.
type NetworkConfig struct {
	ListenAddr     string   // e.g. "/ip4/0.0.0.0/tcp/30303"
	BootstrapPeers []string // multiaddrs of bootstrap peers
	EnableMDNS     bool     // local peer discovery
}

// DefaultNetworkConfig returns a NetworkConfig with sensible defaults.
func DefaultNetworkConfig() *NetworkConfig {
	return &NetworkConfig{
		ListenAddr:     "/ip4/0.0.0.0/tcp/30303",
		EnableMDNS:     true,
		BootstrapPeers: []string{},
	}
}

// Host wraps a libp2p host and provides SpainCoin-specific helpers.
type Host struct {
	host   host.Host
	config *NetworkConfig
}

// NewHost creates a new libp2p host with a random identity.
func NewHost(cfg *NetworkConfig) (*Host, error) {
	h, err := libp2p.New(
		libp2p.ListenAddrStrings(cfg.ListenAddr),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}
	return &Host{
		host:   h,
		config: cfg,
	}, nil
}

// ID returns the peer ID of this host.
func (h *Host) ID() peer.ID {
	return h.host.ID()
}

// Addrs returns the listen addresses of this host as multiaddr strings.
func (h *Host) Addrs() []string {
	addrs := h.host.Addrs()
	result := make([]string, len(addrs))
	for i, addr := range addrs {
		result[i] = fmt.Sprintf("%s/p2p/%s", addr.String(), h.host.ID())
	}
	return result
}

// Close shuts down the host.
func (h *Host) Close() error {
	return h.host.Close()
}

// Connect connects to a peer given its full multiaddr string (including /p2p/<peerID>).
func (h *Host) Connect(ctx context.Context, addrStr string) error {
	maddr, err := multiaddr.NewMultiaddr(addrStr)
	if err != nil {
		return fmt.Errorf("invalid multiaddr %q: %w", addrStr, err)
	}
	pi, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return fmt.Errorf("failed to parse peer addr info: %w", err)
	}
	return h.host.Connect(ctx, *pi)
}

// PeerCount returns the number of currently connected peers.
func (h *Host) PeerCount() int {
	return len(h.host.Network().Peers())
}

// Peers returns the peer IDs of all connected peers.
func (h *Host) Peers() []peer.ID {
	return h.host.Network().Peers()
}

// libp2pHost exposes the underlying libp2p host (used by PubSub and Discovery).
func (h *Host) libp2pHost() host.Host {
	return h.host
}
