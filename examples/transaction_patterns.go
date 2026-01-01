package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/0xkanth/polymarket-indexer/pkg/config"
	"github.com/0xkanth/polymarket-indexer/pkg/service"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig("config/chains.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	chainCfg, err := cfg.GetChain("polygon")
	if err != nil {
		log.Fatalf("Failed to get chain config: %v", err)
	}

	// Create service
	svc, err := service.NewCTFService(ctx, chainCfg)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	defer svc.Close()

	// Example 1: Read operation (no transaction)
	fmt.Println("=== Example 1: Read Order Status ===")
	var orderHash [32]byte
	copy(orderHash[:], common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"))

	status, err := svc.GetOrderStatus(ctx, orderHash)
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Order status: %d\n", status)
	}

	// Example 2: Simulate a transaction before sending
	fmt.Println("\n=== Example 2: Simulate Transaction ===")

	// Create a test wallet (use your own private key)
	privateKey, err := crypto.HexToECDSA("your_private_key_here")
	if err != nil {
		log.Printf("Private key error: %v", err)
		return
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainCfg.ChainID))
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	// Prepare transaction message for simulation
	msg := ethereum.CallMsg{
		From:  auth.From,
		To:    &svc.GetCTFExchangeAddress(), // Helper method you'd add
		Value: big.NewInt(0),
		Data:  nil, // ABI-encoded function call
	}

	// Simulate before sending
	err = svc.SimulateTransaction(ctx, msg)
	if err != nil {
		fmt.Printf("Simulation failed: %v\n", err)
		fmt.Println("Transaction would revert - NOT sending")
		return
	}
	fmt.Println("Simulation successful - safe to proceed")

	// Example 3: Estimate gas with buffer
	fmt.Println("\n=== Example 3: Estimate Gas ===")

	gasLimit, err := svc.EstimateGasWithBuffer(ctx, msg, 20)
	if err != nil {
		log.Printf("Gas estimation error: %v", err)
	} else {
		fmt.Printf("Estimated gas with 20%% buffer: %d\n", gasLimit)
	}

	// Example 4: Send transaction with production-grade retry
	fmt.Println("\n=== Example 4: Send Transaction with Retry ===")

	// Configure custom retry settings
	txConfig := &service.TransactionConfig{
		MaxRetries:       5,                // Try up to 5 times
		InitialBackoff:   2 * time.Second,  // Start with 2s
		MaxBackoff:       60 * time.Second, // Max 60s between retries
		GasBufferPercent: 25,               // 25% gas buffer
		Simulate:         true,             // Simulate first
		TimeoutPerTry:    30 * time.Second, // 30s per attempt
	}

	// Send transaction with retry
	tx, err := svc.SendTransactionWithRetry(
		ctx,
		msg,
		auth,
		txConfig,
		func(opts *bind.TransactOpts) (*types.Transaction, error) {
			// Your actual transaction function
			// Example: return svc.FillOrder(...)
			return nil, fmt.Errorf("example - not sending real tx")
		},
	)

	if err != nil {
		fmt.Printf("Transaction failed after retries: %v\n", err)
	} else {
		fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

		// Wait for confirmation
		receipt, err := svc.WaitForTransaction(ctx, tx)
		if err != nil {
			log.Printf("Error waiting for tx: %v", err)
		} else {
			fmt.Printf("Transaction mined in block %d\n", receipt.BlockNumber.Uint64())
		}
	}

	// Example 5: High-level ExecuteTransaction (recommended)
	fmt.Println("\n=== Example 5: High-level Execute (Recommended) ===")

	// This combines simulation, gas estimation, and retry in one call
	tx, err = svc.ExecuteTransaction(ctx, msg, auth, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		// Your actual transaction function
		// Example: return svc.FillOrder(ctx, opts, order, amount, sig)
		return nil, fmt.Errorf("example - not sending real tx")
	})

	if err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("Production-grade transaction helpers provide:")
	fmt.Println("✓ Transaction simulation before sending")
	fmt.Println("✓ Automatic gas estimation with buffer")
	fmt.Println("✓ Exponential backoff retry on RPC failures")
	fmt.Println("✓ Smart error detection (retryable vs permanent)")
	fmt.Println("✓ Configurable timeout and retry limits")
}
