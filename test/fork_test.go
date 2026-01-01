package test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/0xkanth/polymarket-indexer/pkg/config"
	"github.com/0xkanth/polymarket-indexer/pkg/service"
)

// TestForkRead tests reading from forked Polygon mainnet
func TestForkRead(t *testing.T) {
	// Load config
	cfg, err := config.LoadConfig("../config/chains.json")
	require.NoError(t, err)

	chainCfg, err := cfg.GetChain("polygon-fork")
	require.NoError(t, err)

	// Create service
	ctx := context.Background()
	svc, err := service.NewCTFService(ctx, chainCfg)
	require.NoError(t, err)
	defer svc.Close()

	// Test 1: Read a known order hash status
	// Replace with actual order hash from Polygon mainnet
	var orderHash [32]byte
	copy(orderHash[:], common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000001"))

	status, err := svc.GetOrderStatus(ctx, orderHash)
	require.NoError(t, err)
	t.Logf("Order status: %d", status)

	// Test 2: Check balance of a known address
	// Replace with actual address that holds tokens on mainnet
	testAddr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	positionId := big.NewInt(1)

	balance, err := svc.BalanceOf(ctx, testAddr, positionId)
	require.NoError(t, err)
	t.Logf("Balance of %s for position %s: %s", testAddr.Hex(), positionId.String(), balance.String())
}

// TestForkWrite tests sending transactions to forked chain
func TestForkWrite(t *testing.T) {
	// Skip if not running fork tests
	if testing.Short() {
		t.Skip("Skipping fork test in short mode")
	}

	// Load config
	cfg, err := config.LoadConfig("../config/chains.json")
	require.NoError(t, err)

	chainCfg, err := cfg.GetChain("polygon-fork")
	require.NoError(t, err)

	// Create service
	ctx := context.Background()
	svc, err := service.NewCTFService(ctx, chainCfg)
	require.NoError(t, err)
	defer svc.Close()

	// Create a test wallet (Anvil default account)
	// Private key for Anvil's first default account
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	require.NoError(t, err)

	// Create auth
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainCfg.ChainID))
	require.NoError(t, err)

	// Set gas parameters
	auth.GasLimit = 500000
	auth.GasPrice = big.NewInt(30000000000) // 30 gwei

	// Example: Fill an order
	// You need to create a valid order struct based on mainnet data
	order := struct {
		// Fill this based on CTFExchange Order struct
		// Check pkg/contracts/CTFExchange.go for the exact struct
	}{}

	fillAmount := big.NewInt(1000000) // 1 USDC (6 decimals)
	signature := []byte{}             // Valid signature

	// Note: This will fail without valid order data
	// Use it as template for your actual tests
	/*
		tx, err := svc.FillOrder(ctx, auth, order, fillAmount, signature)
		if err != nil {
			t.Logf("Expected error with dummy data: %v", err)
			return
		}

		// Wait for transaction
		receipt, err := svc.WaitForTransaction(ctx, tx)
		require.NoError(t, err)
		require.Equal(t, uint64(1), receipt.Status)
		t.Logf("Transaction mined: %s", tx.Hash().Hex())
	*/
}

// TestForkEvents tests listening to events from forked chain
func TestForkEvents(t *testing.T) {
	// Load config
	cfg, err := config.LoadConfig("../config/chains.json")
	require.NoError(t, err)

	chainCfg, err := cfg.GetChain("polygon-fork")
	require.NoError(t, err)

	// Create service
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	svc, err := service.NewCTFService(ctx, chainCfg)
	require.NoError(t, err)
	defer svc.Close()

	// Filter historical events from a specific block range
	// Use startBlock from config
	fromBlock := chainCfg.StartBlock
	toBlock := fromBlock + 1000 // Scan 1000 blocks

	t.Logf("Scanning events from block %d to %d", fromBlock, toBlock)

	iter, err := svc.FilterOrderFilled(ctx, fromBlock, toBlock, nil, nil, nil)
	require.NoError(t, err)
	defer iter.Close()

	eventCount := 0
	for iter.Next() {
		event := iter.Event
		t.Logf("OrderFilled: hash=%x, maker=%s, taker=%s, makerAmount=%s",
			event.OrderHash,
			event.Maker.Hex(),
			event.Taker.Hex(),
			event.MakerAssetFilledAmount.String(),
		)
		eventCount++
	}

	require.NoError(t, iter.Error())
	t.Logf("Found %d OrderFilled events", eventCount)
}
