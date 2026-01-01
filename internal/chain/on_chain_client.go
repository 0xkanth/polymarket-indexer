// Package chain provides Ethereum/Polygon RPC client functionality.
package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

// OnChainClient provides methods to interact with the Ethereum/Polygon blockchain.
type OnChainClient struct {
	rpcClient *ethclient.Client
	wsClient  *ethclient.Client
	chainID   *big.Int
	logger    *zerolog.Logger
}

// NewClient creates a new blockchain client with both HTTP and WebSocket connections.
func NewClient(rpcURL, wsURL string, chainID int64, logger *zerolog.Logger) (*OnChainClient, error) {
	// Connect to HTTP RPC endpoint
	rpcClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC endpoint: %w", err)
	}

	// Connect to WebSocket endpoint (optional, for real-time subscriptions)
	var wsClient *ethclient.Client
	if wsURL != "" {
		wsClient, err = ethclient.Dial(wsURL)
		if err != nil {
			logger.Warn().
				Err(err).
				Str("ws_url", wsURL).
				Msg("failed to connect to WebSocket endpoint, will use HTTP only")
		}
	}

	// Verify chain ID
	actualChainID, err := rpcClient.ChainID(context.Background())
	if err != nil {
		rpcClient.Close()
		if wsClient != nil {
			wsClient.Close()
		}
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	expectedChainID := big.NewInt(chainID)
	if actualChainID.Cmp(expectedChainID) != 0 {
		rpcClient.Close()
		if wsClient != nil {
			wsClient.Close()
		}
		return nil, fmt.Errorf("chain ID mismatch: expected %d, got %d", chainID, actualChainID)
	}

	logger.Info().
		Int64("chain_id", chainID).
		Str("rpc_url", rpcURL).
		Bool("has_websocket", wsClient != nil).
		Msg("blockchain client initialized")

	return &OnChainClient{
		rpcClient: rpcClient,
		wsClient:  wsClient,
		chainID:   expectedChainID,
		logger:    logger,
	}, nil
}

// GetLatestBlockNumber returns the latest block number from the chain.
func (c *OnChainClient) GetLatestBlockNumber(ctx context.Context) (uint64, error) {
	blockNumber, err := c.rpcClient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block number: %w", err)
	}
	return blockNumber, nil
}

// GetBlockByNumber fetches a block by its number.
func (c *OnChainClient) GetBlockByNumber(ctx context.Context, blockNumber uint64) (*types.Block, error) {
	block, err := c.rpcClient.BlockByNumber(ctx, big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block %d: %w", blockNumber, err)
	}
	return block, nil
}

// GetBlockByHash fetches a block by its hash.
func (c *OnChainClient) GetBlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	block, err := c.rpcClient.BlockByHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch block by hash %s: %w", hash.Hex(), err)
	}
	return block, nil
}

// GetTransactionReceipt fetches a transaction receipt.
func (c *OnChainClient) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := c.rpcClient.TransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch receipt for tx %s: %w", txHash.Hex(), err)
	}
	return receipt, nil
}

// GetBlockReceipts fetches all receipts for a given block.
// This is more efficient than fetching receipts individually.
func (c *OnChainClient) GetBlockReceipts(ctx context.Context, blockNumber uint64) ([]*types.Receipt, error) {
	block, err := c.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	receipts := make([]*types.Receipt, 0, len(block.Transactions()))
	for _, tx := range block.Transactions() {
		receipt, err := c.GetTransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("failed to fetch receipt for tx %s in block %d: %w",
				tx.Hash().Hex(), blockNumber, err)
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

// FilterLogs queries for logs matching the given filter.
func (c *OnChainClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	logs, err := c.rpcClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}
	return logs, nil
}

// SubscribeNewHead subscribes to new block headers via WebSocket.
// Returns nil if WebSocket client is not available.
func (c *OnChainClient) SubscribeNewHead(ctx context.Context) (chan *types.Header, ethereum.Subscription, error) {
	if c.wsClient == nil {
		return nil, nil, fmt.Errorf("websocket client not available")
	}

	headers := make(chan *types.Header)
	sub, err := c.wsClient.SubscribeNewHead(ctx, headers)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe to new heads: %w", err)
	}

	return headers, sub, nil
}

// ChainID returns the chain ID.
func (c *OnChainClient) ChainID() *big.Int {
	return c.chainID
}

// Close closes the client connections.
func (c *OnChainClient) Close() {
	c.rpcClient.Close()
	if c.wsClient != nil {
		c.wsClient.Close()
	}
	c.logger.Info().Msg("blockchain client closed")
}
