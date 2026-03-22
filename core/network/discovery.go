package network

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// Discovery manages local peer discovery via mDNS.
type Discovery struct {
	mdnsService mdns.Service
}

// MDNSNotifee implements mdns.Notifee. When a peer is found via mDNS, it
// connects to that peer automatically.
type MDNSNotifee struct {
	h *Host
}

// HandlePeerFound is called by the mDNS service when a new peer is discovered.
// It attempts to connect using a background context and silently ignores errors
// (the peer may already be connected or unreachable).
func (n *MDNSNotifee) HandlePeerFound(pi peer.AddrInfo) {
	_ = n.h.libp2pHost().Connect(context.Background(), pi)
}

// NewMDNSDiscovery starts mDNS discovery on the local network and returns a
// Discovery that can be closed later.
func NewMDNSDiscovery(ctx context.Context, h *Host) (*Discovery, error) {
	notifee := &MDNSNotifee{h: h}
	svc := mdns.NewMdnsService(h.libp2pHost(), mdns.ServiceName, notifee)
	if err := svc.Start(); err != nil {
		return nil, err
	}
	return &Discovery{mdnsService: svc}, nil
}

// Close stops the mDNS discovery service.
func (d *Discovery) Close() {
	_ = d.mdnsService.Close()
}
