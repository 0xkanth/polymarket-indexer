// Package processor handles block and event processing for Polymarket contracts.
//
// WHAT THIS DOES:
// This is the CORE ENGINE of the indexer. It continuously polls Polygon blockchain,
// extracts events from CTF Exchange and Conditional Tokens contracts, and publishes
// them to NATS JetStream for the consumer to write to TimescaleDB.
//
// ARCHITECTURE FLOW:
// 1. ProcessBlocks() runs in a loop polling for new blocks
// 2. For each block, calls FilterLogs() to get all events from monitored contracts
// 3. Calls processLog() which routes each event to the correct handler (OrderFilled, OrdersMatched, etc.)
// 4. Handler decodes the event and publishes it to NATS as JSON
// 5. Consumer picks up from NATS and writes to TimescaleDB
//
// KEY COMPONENTS:
// - chain.OnChainClient: Ethereum JSON-RPC client wrapper (go-ethereum)
// - router.EventLogHandlerRouter: Maps event signatures to handler functions
// - nats.Publisher: Publishes events to NATS JetStream
// - handler.Events: Decodes ABI events into Go structs
//
// PROMETHEUS METRICS:
// - polymarket_blocks_processed_total: Blocks processed
// - polymarket_events_processed_total: Events by type (OrderFilled, OrdersMatched, etc.)
// - polymarket_block_processing_duration_seconds: Performance tracking
// - polymarket_processing_errors_total: Error monitoring
//
// USAGE:
// p := processor.New(logger, chainClient, natsPublisher, cfg)
// go p.ProcessBlocks(ctx, currentBlock)
package processor

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"

	"github.com/0xkanth/polymarket-indexer/internal/chain"
	"github.com/0xkanth/polymarket-indexer/internal/handler"
	"github.com/0xkanth/polymarket-indexer/internal/nats"
	"github.com/0xkanth/polymarket-indexer/internal/router"
	"github.com/0xkanth/polymarket-indexer/pkg/models"
)

var (
	blocksProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "polymarket_blocks_processed_total",
		Help: "Total number of blocks processed",
	})

	eventsProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polymarket_events_processed_total",
		Help: "Total number of events processed by type",
	}, []string{"event_type"})

	processingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "polymarket_block_processing_duration_seconds",
		Help:    "Time taken to process a block",
		Buckets: prometheus.DefBuckets,
	})

	processingErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polymarket_processing_errors_total",
		Help: "Total number of processing errors",
	}, []string{"error_type"})
)

// BlockEventsProcessor handles block and event processing.
type BlockEventsProcessor struct {
	logger                zerolog.Logger
	chain                 *chain.OnChainClient
	eventLogHandlerRouter *router.EventLogHandlerRouter
	natsEventPublisher    *nats.Publisher
	contracts             []common.Address
	startBlock            uint64
}

// BlockEventProcessingConfig holds processor configuration.
type BlockEventProcessingConfig struct {
	Contracts  []string // Contract addresses to monitor
	StartBlock uint64   // Block to start processing from
}

// New creates a new processor.
func New(
	logger zerolog.Logger,
	chain *chain.OnChainClient,
	natsEventPublisher *nats.Publisher,
	cfg BlockEventProcessingConfig,
) (*BlockEventsProcessor, error) {
	// Parse contract addresses
	contracts := make([]common.Address, len(cfg.Contracts))
	for i, addr := range cfg.Contracts {
		if !common.IsHexAddress(addr) {
			return nil, fmt.Errorf("invalid contract address: %s", addr)
		}
		contracts[i] = common.HexToAddress(addr)
	}

	// Create event callback that publishes to NATS
	eventCallback := func(ctx context.Context, event models.Event) error {
		return natsEventPublisher.Publish(ctx, event)
	}

	// Create eventLogHandlerRouter with callback
	r := router.New(eventCallback)

	// Register CTF Exchange handlers
	r.RegisterLogHandler(handler.OrderFilledSig, "OrderFilled", handler.HandleOrderFilled)
	r.RegisterLogHandler(handler.OrderCancelledSig, "OrderCancelled", handler.HandleOrderCancelled)
	r.RegisterLogHandler(handler.TokenRegisteredSig, "TokenRegistered", handler.HandleTokenRegistered)

	// Register Conditional Tokens handlers
	r.RegisterLogHandler(handler.TransferSingleSig, "TransferSingle", handler.HandleTransferSingle)
	r.RegisterLogHandler(handler.TransferBatchSig, "TransferBatch", handler.HandleTransferBatch)
	r.RegisterLogHandler(handler.ConditionPreparationSig, "ConditionPreparation", handler.HandleConditionPreparation)
	r.RegisterLogHandler(handler.ConditionResolutionSig, "ConditionResolution", handler.HandleConditionResolution)
	r.RegisterLogHandler(handler.PositionSplitSig, "PositionSplit", handler.HandlePositionSplit)
	r.RegisterLogHandler(handler.PositionsMergeSig, "PositionsMerge", handler.HandlePositionsMerge)

	return &BlockEventsProcessor{
		logger:                logger.With().Str("component", "processor").Logger(),
		chain:                 chain,
		eventLogHandlerRouter: r,
		natsEventPublisher:    natsEventPublisher,
		contracts:             contracts,
		startBlock:            cfg.StartBlock,
	}, nil
}

// ProcessBlock processes a single block.
func (p *BlockEventsProcessor) ProcessBlock(ctx context.Context, blockNumber uint64) error {
	start := time.Now()
	defer func() {
		processingDuration.Observe(time.Since(start).Seconds())
	}()

	p.logger.Debug().Uint64("block", blockNumber).Msg("processing block")

	// Fetch block header
	block, err := p.chain.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		processingErrors.WithLabelValues("fetch_block").Inc()
		return fmt.Errorf("failed to get block %d: %w", blockNumber, err)
	}

	// Filter logs for monitored contracts
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(blockNumber)),
		ToBlock:   big.NewInt(int64(blockNumber)),
		Addresses: p.contracts,
	}
	logs, err := p.chain.FilterLogs(ctx, query)
	if err != nil {
		processingErrors.WithLabelValues("filter_logs").Inc()
		return fmt.Errorf("failed to filter logs for block %d: %w", blockNumber, err)
	}

	if len(logs) == 0 {
		p.logger.Debug().
			Uint64("block", blockNumber).
			Uint64("timestamp", block.Time()).
			Msg("no events in block")
		blocksProcessed.Inc()
		return nil
	}

	p.logger.Info().
		Uint64("block", blockNumber).
		Uint64("timestamp", block.Time()).
		Int("events", len(logs)).
		Msg("processing block with events")

	// Process each log
	for _, log := range logs {
		if err := p.processLog(ctx, log, block.Header(), block.Hash().Hex()); err != nil {
			processingErrors.WithLabelValues("process_log").Inc()
			p.logger.Error().
				Err(err).
				Str("tx", log.TxHash.Hex()).
				Uint("log_index", log.Index).
				Msg("failed to process log")
			// Continue processing other logs
			continue
		}
	}

	blocksProcessed.Inc()
	return nil
}

// processLog processes a single log entry.
func (p *BlockEventsProcessor) processLog(ctx context.Context, log types.Log, header *types.Header, blockHash string) error {
	if log.Removed {
		p.logger.Warn().
			Str("tx", log.TxHash.Hex()).
			Uint("log_index", log.Index).
			Msg("skipping removed log")
		return nil
	}

	// Route log to appropriate handler (this publishes via callback)
	err := p.eventLogHandlerRouter.RouteLog(ctx, log, header.Time, blockHash)
	if err != nil {
		// Check if it's just an unknown event (no handler registered)
		if len(log.Topics) > 0 && !p.eventLogHandlerRouter.HasHandler(log.Topics[0]) {
			// Unknown event type, skip silently
			p.logger.Debug().
				Str("tx", log.TxHash.Hex()).
				Uint("log_index", log.Index).
				Str("topic0", log.Topics[0].Hex()).
				Msg("no handler for event")
			return nil
		}
		return fmt.Errorf("failed to route log: %w", err)
	}

	// Count event (event name is handled in eventLogHandlerRouter callback)
	var eventName string
	if len(log.Topics) > 0 {
		eventName = p.getEventName(log.Topics[0])
		eventsProcessed.WithLabelValues(eventName).Inc()
	}

	p.logger.Debug().
		Str("event", eventName).
		Str("tx", log.TxHash.Hex()).
		Uint("log_index", log.Index).
		Msg("processed event")

	return nil
}

// getEventName returns a human-readable name for an event signature.
func (p *BlockEventsProcessor) getEventName(sig common.Hash) string {
	switch sig {
	case handler.OrderFilledSig:
		return "OrderFilled"
	case handler.OrderCancelledSig:
		return "OrderCancelled"
	case handler.TokenRegisteredSig:
		return "TokenRegistered"
	case handler.TransferSingleSig:
		return "TransferSingle"
	case handler.TransferBatchSig:
		return "TransferBatch"
	case handler.ConditionPreparationSig:
		return "ConditionPreparation"
	case handler.ConditionResolutionSig:
		return "ConditionResolution"
	case handler.PositionSplitSig:
		return "PositionSplit"
	case handler.PositionsMergeSig:
		return "PositionsMerge"
	default:
		return "Unknown"
	}
}

// ProcessBlockRange processes a range of blocks.
func (p *BlockEventsProcessor) ProcessBlockRange(ctx context.Context, from, to uint64) error {
	p.logger.Info().
		Uint64("from", from).
		Uint64("to", to).
		Uint64("count", to-from+1).
		Msg("processing block range")

	for block := from; block <= to; block++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := p.ProcessBlock(ctx, block); err != nil {
			return fmt.Errorf("failed to process block %d: %w", block, err)
		}
	}

	return nil
}
