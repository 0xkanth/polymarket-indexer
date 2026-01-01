package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/0xkanth/polymarket-indexer/pkg/config"
	"github.com/0xkanth/polymarket-indexer/pkg/service"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Example of using the CTF service with config
func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig("config/chains.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Get chain configuration (polygon, polygon-fork, or mumbai)
	chainName := "polygon" // or "polygon-fork" for testing
	chainCfg, err := cfg.GetChain(chainName)
	if err != nil {
		log.Fatalf("Failed to get chain config: %v", err)
	}

	// Print configuration
	fmt.Printf("Chain: %s (ID: %d)\n", chainCfg.Name, chainCfg.ChainID)
	fmt.Printf("RPC URLs: %v\n", chainCfg.RPCUrls)
	fmt.Printf("CTFExchange: %s\n", chainCfg.GetCTFExchangeAddress().Hex())
	fmt.Printf("ConditionalTokens: %s\n", chainCfg.GetConditionalTokensAddress().Hex())
	fmt.Printf("Start Block: %d\n", chainCfg.StartBlock)
	fmt.Printf("Block Time: %d seconds\n", chainCfg.BlockTime)
	fmt.Printf("Confirmations: %d\n", chainCfg.Confirmations)
	fmt.Println()

	// Create CTF service
	svc, err := service.NewCTFService(ctx, chainCfg)
	if err != nil {
		log.Fatalf("Failed to create CTF service: %v", err)
	}
	defer svc.Close()

	fmt.Println("CTF Service initialized successfully!")

	// Example 1: Read contract data
	fmt.Println("\n=== Example 1: Reading Contract Data ===")
	exampleOrderHash := [32]byte{} // Replace with real order hash
	orderStatus, err := svc.GetOrderStatus(ctx, exampleOrderHash)
	if err != nil {
		log.Printf("Failed to get order status: %v", err)
	} else {
		fmt.Printf("Order Status: %+v\n", orderStatus)
	}

	// Example 2: Get position complement
	fmt.Println("\n=== Example 2: Get Position Complement ===")
	positionId := big.NewInt(12345)
	complement, err := svc.GetComplement(ctx, positionId)
	if err != nil {
		log.Printf("Failed to get complement: %v", err)
	} else {
		fmt.Printf("Position %s complement: %s\n", positionId.String(), complement.String())
	}

	// Example 3: Listen to events
	fmt.Println("\n=== Example 3: Listening to OrderFilled Events ===")
	fromBlock := chainCfg.StartBlock
	toBlock := fromBlock + 1000 // Scan first 1000 blocks

	fmt.Printf("Scanning events from block %d to %d\n", fromBlock, toBlock)

	// Filter OrderFilled events
	ctfExchangeAddr := chainCfg.GetCTFExchangeAddress()
	filterOpts := &bind.FilterOpts{
		Start:   fromBlock,
		End:     &toBlock,
		Context: ctx,
	}

	// Note: This is a simplified example. In production, use the processor/handler pattern
	fmt.Println("Scanning for OrderFilled events...")
	fmt.Printf("Note: Use the full indexer for production event scanning\n")
	fmt.Printf("This example shows how to use the service for reading contract data\n")

	// Example 4: Get condition ID for a token
	fmt.Println("\n=== Example 4: Get Condition ID ===")
	tokenId := big.NewInt(67890)
	conditionId, err := svc.GetConditionId(ctx, tokenId)
	if err != nil {
		log.Printf("Failed to get condition ID: %v", err)
	} else {
		fmt.Printf("Token %s has condition ID: %x\n", tokenId.String(), conditionId)
	}

	fmt.Println("\n=== Service Examples Complete ===")
	fmt.Println("For full event indexing, see cmd/indexer/main.go")
	fmt.Println("For transaction examples, see examples/transaction_patterns.go")
}
