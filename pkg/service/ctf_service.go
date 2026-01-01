package service

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"

	"github.com/0xkanth/polymarket-indexer/pkg/config"
	"github.com/0xkanth/polymarket-indexer/pkg/contracts"
	"github.com/0xkanth/polymarket-indexer/pkg/txhelper"
)

// CTFService provides methods to interact with CTFExchange contract
type CTFService struct {
	client                *ethclient.Client
	chainConfig           *config.ChainConfig
	ctfExchange           *contracts.CTFExchange
	conditionalTokens     *contracts.ConditionalTokens
	ctfExchangeAddr       common.Address
	conditionalTokensAddr common.Address
	txHelper              *txhelper.TransactionHelper
}

// NewCTFService creates a new CTFService instance
func NewCTFService(ctx context.Context, chainConfig *config.ChainConfig) (*CTFService, error) {
	// Try connecting to RPC endpoints with fallback
	var client *ethclient.Client
	var err error

	for i, rpcURL := range chainConfig.RPCUrls {
		client, err = ethclient.DialContext(ctx, rpcURL)
		if err != nil {
			log.Printf("Failed to connect to RPC %d (%s): %v", i, rpcURL, err)
			continue
		}

		// Test connection
		_, err = client.ChainID(ctx)
		if err != nil {
			log.Printf("RPC %d responded but chain ID failed: %v", i, err)
			client.Close()
			continue
		}

		log.Printf("Connected to %s via RPC %d", chainConfig.Name, i)
		break
	}

	if client == nil {
		return nil, fmt.Errorf("failed to connect to any RPC endpoint")
	}

	ctfExchangeAddr := chainConfig.GetCTFExchangeAddress()
	conditionalTokensAddr := chainConfig.GetConditionalTokensAddress()

	// Bind to CTFExchange contract
	ctfExchange, err := contracts.NewCTFExchange(ctfExchangeAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind CTFExchange: %w", err)
	}

	// Bind to ConditionalTokens contract
	conditionalTokens, err := contracts.NewConditionalTokens(conditionalTokensAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind ConditionalTokens: %w", err)
	}

	// Create transaction helper
	txHelper := txhelper.NewTransactionHelper(client, chainConfig.BlockTime, chainConfig.Confirmations)

	return &CTFService{
		client:                client,
		chainConfig:           chainConfig,
		ctfExchange:           ctfExchange,
		conditionalTokens:     conditionalTokens,
		ctfExchangeAddr:       ctfExchangeAddr,
		conditionalTokensAddr: conditionalTokensAddr,
		txHelper:              txHelper,
	}, nil
}

// Close closes the underlying client connection
func (s *CTFService) Close() {
	s.client.Close()
}

// ============================================================================
// READ METHODS (View/Pure functions - No gas cost)
// ============================================================================

// GetOrderStatus returns the status of an order by its hash
func (s *CTFService) GetOrderStatus(ctx context.Context, orderHash [32]byte) (contracts.OrderStatus, error) {
	status, err := s.ctfExchange.GetOrderStatus(&bind.CallOpts{Context: ctx}, orderHash)
	if err != nil {
		return contracts.OrderStatus{}, fmt.Errorf("failed to get order status: %w", err)
	}
	return status, nil
}

// GetComplement returns the complement of a position ID
func (s *CTFService) GetComplement(ctx context.Context, token *big.Int) (*big.Int, error) {
	complement, err := s.ctfExchange.GetComplement(&bind.CallOpts{Context: ctx}, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get complement: %w", err)
	}
	return complement, nil
}

// Note: GetConditionId and BalanceOf methods are available on CTFExchange, not ConditionalTokens
// CTFExchange has getConditionId(uint256 token) view method
// For ERC1155 balances, you need to use the ConditionalTokens contract directly with proper ABI

// GetConditionId returns the condition ID for a token
func (s *CTFService) GetConditionId(ctx context.Context, token *big.Int) ([32]byte, error) {
	conditionId, err := s.ctfExchange.GetConditionId(&bind.CallOpts{Context: ctx}, token)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to get condition ID: %w", err)
	}
	return conditionId, nil
}

// ============================================================================
// TRANSACTION HELPERS (Delegated to txhelper package)
// ============================================================================

// ExecuteTransaction executes a transaction with simulation, gas estimation, and retry
func (s *CTFService) ExecuteTransaction(
	ctx context.Context,
	msg ethereum.CallMsg,
	auth *bind.TransactOpts,
	sendFunc func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, error) {
	return s.txHelper.ExecuteTransaction(ctx, msg, auth, sendFunc)
}

// SendTransactionWithRetry sends a transaction with custom retry configuration
func (s *CTFService) SendTransactionWithRetry(
	ctx context.Context,
	msg ethereum.CallMsg,
	auth *bind.TransactOpts,
	config *txhelper.TransactionConfig,
	sendFunc func(*bind.TransactOpts) (*types.Transaction, error),
) (*types.Transaction, error) {
	return s.txHelper.SendTransactionWithRetry(ctx, msg, auth, config, sendFunc)
}

// SimulateTransaction simulates a transaction before sending
func (s *CTFService) SimulateTransaction(ctx context.Context, msg ethereum.CallMsg) error {
	return s.txHelper.SimulateTransaction(ctx, msg)
}

// EstimateGasWithBuffer estimates gas with a buffer percentage
func (s *CTFService) EstimateGasWithBuffer(ctx context.Context, msg ethereum.CallMsg, bufferPercent int) (uint64, error) {
	return s.txHelper.EstimateGasWithBuffer(ctx, msg, bufferPercent)
}

// WaitForTransaction waits for a transaction to be mined
func (s *CTFService) WaitForTransaction(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	return s.txHelper.WaitForTransaction(ctx, tx)
}

// ============================================================================
// WRITE METHODS (State-changing transactions - Require gas)
// ============================================================================

// FillOrder fills an order on CTFExchange with production-grade retry logic
func (s *CTFService) FillOrder(
	ctx context.Context,
	auth *bind.TransactOpts,
	order contracts.Order,
	fillAmount *big.Int,
	signature []byte,
) (*types.Transaction, error) {
	// Prepare call message for simulation and gas estimation
	msg := ethereum.CallMsg{
		From:  auth.From,
		To:    &s.ctfExchangeAddr,
		Value: auth.Value,
		Data:  nil, // Would need ABI-encoded FillOrder call data
	}

	// Use production-grade transaction execution
	// Note: Signature is embedded in the Order struct, not a separate parameter
	return s.ExecuteTransaction(ctx, msg, auth, func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return s.ctfExchange.FillOrder(opts, order, fillAmount)
	})
}

// FillOrderSimple is a simple version without retry logic (for testing)
func (s *CTFService) FillOrderSimple(
	ctx context.Context,
	auth *bind.TransactOpts,
	order contracts.Order,
	fillAmount *big.Int,
	signature []byte,
) (*types.Transaction, error) {
	auth.GasLimit = 0 // Let it estimate

	// Signature is embedded in the Order struct
	tx, err := s.ctfExchange.FillOrder(auth, order, fillAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to fill order: %w", err)
	}

	log.Printf("FillOrder transaction sent: %s", tx.Hash().Hex())
	return tx, nil
}

// ============================================================================
// EVENT LISTENING
// ============================================================================

// WatchOrderFilled watches for OrderFilled events
func (s *CTFService) WatchOrderFilled(
	ctx context.Context,
	sink chan<- *contracts.CTFExchangeOrderFilled,
	orderHash [][32]byte,
	maker []common.Address,
	taker []common.Address,
) (event.Subscription, error) {
	opts := &bind.WatchOpts{
		Context: ctx,
		Start:   nil, // Start from current block
	}

	sub, err := s.ctfExchange.WatchOrderFilled(opts, sink, orderHash, maker, taker)
	if err != nil {
		return nil, fmt.Errorf("failed to watch OrderFilled: %w", err)
	}

	return sub, nil
}

// FilterOrderFilled filters historical OrderFilled events
func (s *CTFService) FilterOrderFilled(
	ctx context.Context,
	fromBlock, toBlock uint64,
	orderHash [][32]byte,
	maker []common.Address,
	taker []common.Address,
) (*contracts.CTFExchangeOrderFilledIterator, error) {
	opts := &bind.FilterOpts{
		Context: ctx,
		Start:   fromBlock,
		End:     &toBlock,
	}

	iter, err := s.ctfExchange.FilterOrderFilled(opts, orderHash, maker, taker)
	if err != nil {
		return nil, fmt.Errorf("failed to filter OrderFilled: %w", err)
	}

	return iter, nil
}
