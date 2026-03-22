package network

import (
	"context"
	"encoding/json"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"github.com/spaincoin/spaincoin/core/block"
)

const (
	// TopicBlocks is the gossipsub topic for block announcements.
	TopicBlocks = "spaincoin/blocks/1.0"
	// TopicTxs is the gossipsub topic for transaction broadcasts.
	TopicTxs = "spaincoin/txs/1.0"
)

// PubSub manages gossipsub publishing and subscribing for blocks and transactions.
type PubSub struct {
	ps         *pubsub.PubSub
	blockTopic *pubsub.Topic
	txTopic    *pubsub.Topic
	blockSub   *pubsub.Subscription
	txSub      *pubsub.Subscription
}

// BlockMessage is broadcast when a new block is produced or discovered.
type BlockMessage struct {
	Height    uint64 `json:"height"`
	Hash      string `json:"hash"`      // hex
	PrevHash  string `json:"prevHash"`  // hex
	Validator string `json:"validator"` // address string
	TxCount   int    `json:"txCount"`
}

// TxMessage is broadcast when a new transaction enters the network.
type TxMessage struct {
	ID     string `json:"id"`   // hex hash
	From   string `json:"from"` // address
	To     string `json:"to"`   // address
	Amount uint64 `json:"amount"`
	Fee    uint64 `json:"fee"`
	Nonce  uint64 `json:"nonce"`
}

// NewPubSub creates a gossipsub instance and joins both block and tx topics.
func NewPubSub(ctx context.Context, h *Host) (*PubSub, error) {
	ps, err := pubsub.NewGossipSub(ctx, h.libp2pHost())
	if err != nil {
		return nil, fmt.Errorf("failed to create gossipsub: %w", err)
	}

	blockTopic, err := ps.Join(TopicBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to join blocks topic: %w", err)
	}

	txTopic, err := ps.Join(TopicTxs)
	if err != nil {
		return nil, fmt.Errorf("failed to join txs topic: %w", err)
	}

	blockSub, err := blockTopic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to blocks topic: %w", err)
	}

	txSub, err := txTopic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to txs topic: %w", err)
	}

	return &PubSub{
		ps:         ps,
		blockTopic: blockTopic,
		txTopic:    txTopic,
		blockSub:   blockSub,
		txSub:      txSub,
	}, nil
}

// PublishBlock serializes a BlockMessage as JSON and publishes it to TopicBlocks.
func (ps *PubSub) PublishBlock(b *block.Block) error {
	msg := &BlockMessage{
		Height:    b.Header.Height,
		Hash:      b.Hash.String(),
		PrevHash:  b.Header.PrevHash.String(),
		Validator: b.Header.Validator.String(),
		TxCount:   len(b.Transactions),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal block message: %w", err)
	}
	return ps.blockTopic.Publish(context.Background(), data)
}

// PublishTx serializes a TxMessage as JSON and publishes it to TopicTxs.
func (ps *PubSub) PublishTx(tx *block.Transaction) error {
	msg := &TxMessage{
		ID:     tx.ID.String(),
		From:   tx.From.String(),
		To:     tx.To.String(),
		Amount: tx.Amount,
		Fee:    tx.Fee,
		Nonce:  tx.Nonce,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal tx message: %w", err)
	}
	return ps.txTopic.Publish(context.Background(), data)
}

// NextBlock reads the next block message from the subscription.
func (ps *PubSub) NextBlock(ctx context.Context) (*BlockMessage, error) {
	msg, err := ps.blockSub.Next(ctx)
	if err != nil {
		return nil, err
	}
	var bm BlockMessage
	if err := json.Unmarshal(msg.Data, &bm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal block message: %w", err)
	}
	return &bm, nil
}

// NextTx reads the next tx message from the subscription.
func (ps *PubSub) NextTx(ctx context.Context) (*TxMessage, error) {
	msg, err := ps.txSub.Next(ctx)
	if err != nil {
		return nil, err
	}
	var tm TxMessage
	if err := json.Unmarshal(msg.Data, &tm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tx message: %w", err)
	}
	return &tm, nil
}

// Close cancels all subscriptions and closes topics.
func (ps *PubSub) Close() {
	ps.blockSub.Cancel()
	ps.txSub.Cancel()
	_ = ps.blockTopic.Close()
	_ = ps.txTopic.Close()
}
