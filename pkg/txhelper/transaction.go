package txhelper

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// TransactionHelper provides reusable transaction utilities for any Ethereum client
type TransactionHelper struct {
	client        *ethclient.Client
	blockTime     int // seconds
	confirmations int
}

// NewTransactionHelper creates a new transaction helper
func NewTransactionHelper(client *ethclient.Client, blockTime, confirmations int) *TransactionHelper {
	return &TransactionHelper{
		client:        client,
		blockTime:     blockTime,
		confirmations: confirmations,
	}
}

// TransactionConfig holds configuration for sending transactions
type TransactionConfig struct {
	MaxRetries       int           // Maximum retry attempts (default: 3)
	InitialBackoff   time.Duration // Initial backoff duration (default: 1s)
	MaxBackoff       time.Duration // Maximum backoff duration (default: 30s)
	GasBufferPercent int           // Gas limit buffer % (default: 20)
	Simulate         bool          // Simulate before sending (default: true)
	TimeoutPerTry    time.Duration // Timeout per attempt (default: 30s)
}

// DefaultTransactionConfig returns safe defaults for transaction execution
func DefaultTransactionConfig() *TransactionConfig {
	return &TransactionConfig{
		MaxRetries:       3,
		InitialBackoff:   1 * time.Second,
		MaxBackoff:       30 * time.Second,
		GasBufferPercent: 20,
		Simulate:         true,
		TimeoutPerTry:    30 * time.Second,
	}
}

// SimulateTransaction simulates a transaction using eth_call before sending
// Returns nil if simulation succeeds, error if it would revert
func (h *TransactionHelper) SimulateTransaction(ctx context.Context, msg ethereum.CallMsg) error {
	// Override gas limit for simulation (set high value)
	msg.Gas = 30000000

	result, err := h.client.CallContract(ctx, msg, nil)
	if err != nil {
		// Check if it's a revert with data
		if strings.Contains(err.Error(), "execution reverted") {
			return fmt.Errorf("simulation failed: %w", err)
		}
		return fmt.Errorf("simulation error: %w", err)
	}

	log.Printf("Simulation successful, result length: %d bytes", len(result))
	return nil
}

// EstimateGasWithBuffer estimates gas and adds a buffer percentage
func (h *TransactionHelper) EstimateGasWithBuffer(ctx context.Context, msg ethereum.CallMsg, bufferPercent int) (uint64, error) {
	// Estimate base gas
	gasEstimate, err := h.client.EstimateGas(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("gas estimation failed: %w", err)
	}

	// Add buffer
	buffer := gasEstimate * uint64(bufferPercent) / 100
	gasWithBuffer := gasEstimate + buffer

	// Cap at chain-specific limits
	maxGasLimit := uint64(30000000) // 30M default
	if gasWithBuffer > maxGasLimit {
		gasWithBuffer = maxGasLimit
	}

	log.Printf("Gas estimated: %d, with %d%% buffer: %d", gasEstimate, bufferPercent, gasWithBuffer)
	return gasWithBuffer, nil
}

// IsRetryableError checks if an error is retryable (RPC/network issues)
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// RPC errors (retryable)
	retryableErrors := []string{
		"connection refused",
		"connection reset",
		"EOF",
		"timeout",
		"TLS handshake timeout",
		"no such host",
		"network is unreachable",
		"429", // Rate limit
		"502", // Bad gateway
		"503", // Service unavailable
		"504", // Gateway timeout
	}

	for _, retryable := range retryableErrors {
		if strings.Contains(errStr, retryable) {
			return true
		}
	}

	// Non-retryable errors (permanent failures)
	permanentErrors := []string{
		"execution reverted",
		"insufficient funds",
		"gas too low",
		"nonce too low",
		"replacement transaction underpriced",
		"already known",
	}

	for _, permanent := range permanentErrors {
		if strings.Contains(errStr, permanent) {
			return false
		}
	}

	// Check for RPC error codes
	var rpcErr rpc.Error
	if errors.As(err, &rpcErr) {
		code := rpcErr.ErrorCode()
		// Retryable RPC codes
		if code == -32000 || code == -32603 { // Internal error, may be transient
			return true
		}
	}

	// Default: retry on unknown errors
	return true
}

// SendTransactionWithRetry sends a transaction with exponential backoff retry
func (h *TransactionHelper) SendTransactionWithRetry(
	ctx context.Context,
	msg ethereum.CallMsg,
	auth *bind.TransactOpts,
	config *TransactionConfig,
	sendFunc func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, error) {
	if config == nil {
		config = DefaultTransactionConfig()
	}

	// Step 1: Simulate transaction if enabled
	if config.Simulate {
		log.Println("Simulating transaction...")
		if err := h.SimulateTransaction(ctx, msg); err != nil {
			return nil, fmt.Errorf("simulation failed, aborting: %w", err)
		}
	}

	// Step 2: Estimate gas with buffer
	log.Println("Estimating gas...")
	gasLimit, err := h.EstimateGasWithBuffer(ctx, msg, config.GasBufferPercent)
	if err != nil {
		return nil, fmt.Errorf("gas estimation failed: %w", err)
	}
	auth.GasLimit = gasLimit

	// Step 3: Send transaction with retry logic
	var tx *types.Transaction
	backoff := config.InitialBackoff

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("Retry attempt %d/%d after %v", attempt, config.MaxRetries, backoff)
			time.Sleep(backoff)

			// Exponential backoff with jitter
			backoff = backoff * 2
			if backoff > config.MaxBackoff {
				backoff = config.MaxBackoff
			}
		}

		// Create timeout context for this attempt
		attemptCtx, cancel := context.WithTimeout(ctx, config.TimeoutPerTry)
		auth.Context = attemptCtx

		// Send transaction
		tx, err = sendFunc(auth)
		cancel()

		if err == nil {
			log.Printf("Transaction sent successfully: %s", tx.Hash().Hex())
			return tx, nil
		}

		log.Printf("Attempt %d failed: %v", attempt+1, err)

		// Check if error is retryable
		if !IsRetryableError(err) {
			return nil, fmt.Errorf("non-retryable error: %w", err)
		}

		// Last attempt, don't retry
		if attempt == config.MaxRetries {
			return nil, fmt.Errorf("max retries (%d) reached: %w", config.MaxRetries, err)
		}
	}

	return nil, fmt.Errorf("transaction failed after %d attempts", config.MaxRetries)
}

// ExecuteTransaction is a high-level helper that combines simulation, gas estimation, and retry
// This is the recommended way to send transactions in production
func (h *TransactionHelper) ExecuteTransaction(
	ctx context.Context,
	msg ethereum.CallMsg,
	auth *bind.TransactOpts,
	sendFunc func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, error) {
	config := DefaultTransactionConfig()
	return h.SendTransactionWithRetry(ctx, msg, auth, config, sendFunc)
}

// WaitForTransaction waits for a transaction to be mined and returns the receipt
func (h *TransactionHelper) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	timeout := time.Duration(h.blockTime*h.confirmations) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout*2) // 2x safety margin
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for transaction %s", tx.Hash().Hex())
		default:
		}

		receipt, err := h.client.TransactionReceipt(ctx, tx.Hash())
		if err == nil {
			if receipt.Status == 0 {
				return receipt, fmt.Errorf("transaction reverted: %s", tx.Hash().Hex())
			}

			log.Printf("Transaction mined in block %d with status %d", receipt.BlockNumber.Uint64(), receipt.Status)
			return receipt, nil
		}

		// Wait before next poll
		time.Sleep(time.Duration(h.blockTime) * time.Second)
	}
}

// EstimateTotalGasCost estimates the total cost (gas * gasPrice) for a transaction
func (h *TransactionHelper) EstimateTotalGasCost(ctx context.Context, msg ethereum.CallMsg, bufferPercent int) (*big.Int, error) {
	// Get gas limit
	gasLimit, err := h.EstimateGasWithBuffer(ctx, msg, bufferPercent)
	if err != nil {
		return nil, err
	}

	// Get current gas price
	gasPrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// Calculate total cost
	totalCost := new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPrice)
	log.Printf("Estimated total cost: %s wei (gas: %d, gasPrice: %s)", totalCost.String(), gasLimit, gasPrice.String())

	return totalCost, nil
}

// SuggestGasPriceWithTip suggests gas price with optional priority fee for EIP-1559
func (h *TransactionHelper) SuggestGasPriceWithTip(ctx context.Context, tipPercent int) (*big.Int, error) {
	basePrice, err := h.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get base gas price: %w", err)
	}

	if tipPercent > 0 {
		tip := new(big.Int).Mul(basePrice, big.NewInt(int64(tipPercent)))
		tip.Div(tip, big.NewInt(100))
		basePrice.Add(basePrice, tip)
	}

	log.Printf("Suggested gas price: %s (with %d%% tip)", basePrice.String(), tipPercent)
	return basePrice, nil
}
