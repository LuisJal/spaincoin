package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/spaincoin/spaincoin/core/block"
	"github.com/spaincoin/spaincoin/core/crypto"
)

func runSend(args []string) {
	fs := flag.NewFlagSet("send", flag.ContinueOnError)
	keyHex := fs.String("key", "", "Sender private key (hex)")
	toAddr := fs.String("to", "", "Recipient address (SPC...)")
	amountSPC := fs.Float64("amount", 0, "Amount in SPC (e.g. 200)")
	fee := fs.Uint64("fee", 1000, "Transaction fee in pesetas")
	nodeURL := fs.String("node", "http://localhost:8545", "Node RPC URL")

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if *keyHex == "" {
		// Try environment variable
		*keyHex = os.Getenv("SPC_VALIDATOR_KEY")
	}

	if *keyHex == "" || *toAddr == "" || *amountSPC <= 0 {
		fmt.Fprintln(os.Stderr, "Usage: spc send --key <hex> --to <SPC...> --amount <SPC>")
		fmt.Fprintln(os.Stderr, "  Or set SPC_VALIDATOR_KEY env var instead of --key")
		os.Exit(1)
	}

	// Convert SPC to pesetas (1 SPC = 10^12 pesetas)
	amountPesetas := uint64(*amountSPC * 1_000_000_000_000)

	// Load key
	priv, pub, err := crypto.PrivateKeyFromHex(*keyHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid private key: %v\n", err)
		os.Exit(1)
	}
	fromAddress := pub.ToAddress()

	// Get nonce from node
	nonce, err := getNonce(*nodeURL, fromAddress.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting nonce: %v\n", err)
		os.Exit(1)
	}

	toAddress, err := crypto.AddressFromHex(*toAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid to address: %v\n", err)
		os.Exit(1)
	}

	// Create and sign transaction
	tx := block.NewTransaction(fromAddress, toAddress, amountPesetas, nonce, *fee)
	if err := tx.Sign(priv); err != nil {
		fmt.Fprintf(os.Stderr, "Error signing: %v\n", err)
		os.Exit(1)
	}

	// Send to node
	sigR := tx.Signature.R.Text(16)
	sigS := tx.Signature.S.Text(16)

	body := map[string]interface{}{
		"from":   fromAddress.String(),
		"to":     toAddress.String(),
		"amount": amountPesetas,
		"nonce":  nonce,
		"fee":    *fee,
		"sig_r":  sigR,
		"sig_s":  sigS,
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(*nodeURL+"/tx/send", "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending tx: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		fmt.Fprintf(os.Stderr, "Error: %s\n", string(respBody))
		os.Exit(1)
	}

	fmt.Println("✅ Transacción enviada")
	fmt.Printf("  De:       %s\n", fromAddress.String())
	fmt.Printf("  A:        %s\n", toAddress.String())
	fmt.Printf("  Cantidad: %.4f SPC\n", *amountSPC)
	fmt.Printf("  Nonce:    %d\n", nonce)

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	if txID, ok := result["tx_id"]; ok {
		fmt.Printf("  TX ID:    %s\n", txID)
	}
}

func getNonce(nodeURL, address string) (uint64, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/address/%s/balance", nodeURL, address))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Nonce uint64 `json:"nonce"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Nonce, nil
}
