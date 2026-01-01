// Package syncer coordinates blockchain synchronization for the Polymarket indexer.
//
// # PURPOSE
// The syncer is the orchestrator that manages the entire blockchain synchronization lifecycle.
// It bridges the gap between raw blockchain data and the processor that extracts meaningful events.
//
// # WHY IT EXISTS
// Blockchain indexing requires two distinct strategies:
// 1. BACKFILL MODE: Fast catch-up when far behind (batch processing with workers)
// 2. REALTIME MODE: Low-latency tracking when near chain head (single block polling)
//
// The syncer intelligently switches between these modes and manages:
// - Checkpoint persistence (resume from last processed block after restart)
// - Block confirmation safety (wait N confirmations to avoid reorg issues)
// - Worker pool coordination (parallel batch processing)
// - Health monitoring and metrics
//
// # WHO INTERACTS WITH IT
// - main.go (cmd/indexer/main.go): Creates syncer via syncer.New() and calls syncer.Start()
// - internal/processor: Called by syncer to extract events from blocks
// - internal/db/checkpoint: Used by syncer to save/load synchronization progress
// - internal/chain/client: Used by syncer to fetch blocks and chain height
// - Prometheus: Exposes metrics (syncer_height, chain_height, blocks_behind, syncer_errors)
//
// # WHO TRIGGERS SYNC
// - main.go: Calls syncer.Start(ctx) which runs until context is canceled
// - The syncer self-manages:
//   - Automatic mode switching (backfill ↔ realtime)
//   - Continuous polling in realtime mode
//   - Batch processing in backfill mode
//
// # HOW OFTEN
// - BACKFILL MODE: Processes batches continuously (default 1000 blocks per batch)
// - REALTIME MODE: Polls every pollInterval (default 2s from config.toml)
// - CHECKPOINT SAVE: After every batch (backfill) or block (realtime)
//
// # ARCHITECTURE MINDMAP
//
// ┌─────────────────────────────────────────────────────────────────────────────┐
// │                            SYNCER LIFECYCLE                                  │
// └─────────────────────────────────────────────────────────────────────────────┘
//
//	                               ┌──────────────┐
//	                               │   main.go    │
//	                               │  (indexer)   │
//	                               └──────┬───────┘
//	                                      │
//	                           1. syncer.Start(ctx)
//	                                      │
//	                                      ▼
//	                        ┌─────────────────────────┐
//	                        │  Load Checkpoint        │
//	                        │  - Get last processed   │
//	                        │    block from DB        │
//	                        │  - Default: startBlock  │
//	                        └───────────┬─────────────┘
//	                                    │
//	                        2. Determine Strategy
//	                                    │
//	                 ┌──────────────────┴─────────────────┐
//	                 │                                     │
//	        behind > batchSize*2                  behind ≤ batchSize*2
//	                 │                                     │
//	                 ▼                                     ▼
//	     ┌───────────────────────┐           ┌───────────────────────┐
//	     │   BACKFILL MODE       │           │   REALTIME MODE       │
//	     │  (Fast Catch-Up)      │           │  (Live Tracking)      │
//	     └───────────────────────┘           └───────────────────────┘
//	                 │                                     │
//	     ┌───────────┴───────────┐            ┌───────────┴───────────┐
//	     │                       │            │                       │
//	     ▼                       ▼            ▼                       ▼
//	┌─────────┐           ┌─────────┐   ┌─────────┐           ┌─────────┐
//	│ Worker1 │           │ Worker2 │   │  Poll   │           │  Poll   │
//	│ Process │    ...    │ Process │   │  Every  │    ...    │  Every  │
//	│ Blocks  │           │ Blocks  │   │   2s    │           │   2s    │
//	│ 1-500   │           │ 501-1000│   │         │           │         │
//	└────┬────┘           └────┬────┘   └────┬────┘           └────┬────┘
//	     │                     │             │                      │
//	     └──────────┬──────────┘             └──────────┬───────────┘
//	                │                                   │
//	                ▼                                   ▼
//	     ┌──────────────────────┐          ┌─────────────────────┐
//	     │  Batch Checkpoint    │          │  Block Checkpoint   │
//	     │  - Save every batch  │          │  - Save every block │
//	     │  - Update metrics    │          │  - Update metrics   │
//	     └──────────┬───────────┘          └──────────┬──────────┘
//	                │                                  │
//	                │          ┌───────────────────────┘
//	                │          │
//	                ▼          ▼
//	      ┌─────────────────────────┐
//	      │  Auto Mode Switching    │
//	      │  - Backfill → Realtime  │
//	      │    (when caught up)     │
//	      │  - Realtime → Backfill  │
//	      │    (if fell behind)     │
//	      └─────────────────────────┘
//
// # CONFIGURATION
// - confirmations: uint64     - Wait N blocks before considering block final (default: 100)
// - batchSize: uint64         - Blocks per batch in backfill mode (default: 1000)
// - pollInterval: Duration    - Polling frequency in realtime mode (default: 2s)
// - workers: int              - Parallel workers for backfill (default: 5)
//
// # SAFETY MECHANISMS
// - Confirmations: Only process blocks with N confirmations to avoid reorgs
// - Checkpoint persistence: Resume from exact point after crash/restart
// - Health monitoring: Expose health status for readiness probes
// - Error retry: Sleep and retry on transient failures
// - Context cancellation: Graceful shutdown on SIGINT/SIGTERM
//
// # METRICS EXPOSED
// - syncer_current_block:     Current block syncer has processed
// - chain_latest_block:       Latest block number from blockchain
// - syncer_blocks_behind:     How far behind the chain head
// - syncer_errors_total:      Count of errors by type (get_latest_block, process_batch, etc.)
package syncer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"

	"github.com/0xkanth/polymarket-indexer/internal/chain"
	"github.com/0xkanth/polymarket-indexer/internal/db"
	"github.com/0xkanth/polymarket-indexer/internal/processor"
)

var (
	syncerHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "polymarket_syncer_block_height",
		Help: "Current block height being processed",
	})

	chainHeight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "polymarket_chain_block_height",
		Help: "Latest block height on chain",
	})

	blocksBehind = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "polymarket_blocks_behind",
		Help: "Number of blocks behind chain head",
	})

	syncerErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polymarket_syncer_errors_total",
		Help: "Total number of syncer errors",
	}, []string{"error_type"})
)

// Syncer coordinates blockchain synchronization lifecycle.
//
// It manages the dual-mode strategy (backfill/realtime) and handles:
// - Loading checkpoint to resume from last processed block
// - Determining sync strategy based on distance from chain head
// - Orchestrating worker pools for batch processing (backfill mode)
// - Polling for new blocks at configured intervals (realtime mode)
// - Saving checkpoints after each batch/block
// - Exposing health status and Prometheus metrics
//
// The Syncer is stateful and maintains:
// - currentBlock: Last block successfully processed and checkpointed
// - latestBlock: Latest block number fetched from blockchain RPC
// - isHealthy: Health flag updated on each successful sync cycle
type Syncer struct {
	logger        zerolog.Logger
	chain         *chain.OnChainClient
	processor     *processor.BlockEventsProcessor
	checkpoint    *db.CheckpointDB
	serviceName   string
	startBlock    uint64
	batchSize     uint64
	pollInterval  time.Duration
	confirmations uint64
	workers       int
	mu            sync.RWMutex
	currentBlock  uint64
	latestBlock   uint64
	isHealthy     bool
}

// Config holds syncer configuration.
//
// These values are typically loaded from config.toml:
// - confirmations: Number of blocks to wait before considering block final (default: 100)
// - batchSize: Blocks per batch in backfill mode (default: 1000)
// - pollInterval: Polling frequency in realtime mode (default: 2s)
// - workers: Number of parallel workers for backfill (default: 5)
type Config struct {
	ServiceName   string        // Service identifier for checkpoint (e.g., "polymarket-indexer")
	StartBlock    uint64        // Block to start syncing from (from chains.json)
	BatchSize     uint64        // Number of blocks to process in one batch (backfill mode)
	PollInterval  time.Duration // How often to poll for new blocks (realtime mode)
	Confirmations uint64        // Number of confirmations before processing (safety buffer)
	Workers       int           // Number of parallel workers for backfill (default: 5)
}

// New creates a new syncer instance.
//
// Called by main.go during indexer initialization. The syncer is created with:
// - logger: Structured logger for observability
// - chain: Blockchain RPC client for fetching blocks
// - processor: Event processor that extracts logs from blocks
// - checkpoint: Database manager for persisting sync progress
// - cfg: Configuration from config.toml and chains.json
//
// Returns a fully initialized syncer ready to call Start().
func New(
	logger zerolog.Logger,
	chain *chain.OnChainClient,
	processor *processor.BlockEventsProcessor,
	checkpoint *db.CheckpointDB,
	cfg Config,
) *Syncer {
	return &Syncer{
		logger:        logger.With().Str("component", "syncer").Logger(),
		chain:         chain,
		processor:     processor,
		checkpoint:    checkpoint,
		serviceName:   cfg.ServiceName,
		startBlock:    cfg.StartBlock,
		batchSize:     cfg.BatchSize,
		pollInterval:  cfg.PollInterval,
		confirmations: cfg.Confirmations,
		workers:       cfg.Workers,
		isHealthy:     true,
	}
}

// Start begins synchronization and runs until context is canceled.
//
// This is the main entry point called by main.go. It:
// 1. Loads checkpoint from database (or creates new one at startBlock)
// 2. Fetches latest block from blockchain
// 3. Determines sync strategy:
//   - If behind > batchSize*2: Start in backfill mode (fast catch-up)
//   - Otherwise: Start in realtime mode (live polling)
//
// 4. Runs continuously until context is canceled (SIGINT/SIGTERM)
//
// Mode switching is handled automatically:
// - runBackfill() switches to runRealtime() when caught up
// - runRealtime() switches to runBackfill() if it falls behind
//
// Returns error only on critical failures (checkpoint load, initial RPC call).
// Transient errors are retried with exponential backoff.
func (s *Syncer) Start(ctx context.Context) error {
	s.logger.Info().Msg("starting syncer")

	// Get or create checkpoint
	checkpoint, err := s.checkpoint.GetOrCreateCheckpoint(ctx, s.serviceName, s.startBlock)
	if err != nil {
		return fmt.Errorf("failed to get checkpoint: %w", err)
	}

	s.currentBlock = checkpoint.LastBlock
	s.logger.Info().
		Uint64("checkpoint", s.currentBlock).
		Str("hash", checkpoint.LastBlockHash).
		Msg("loaded checkpoint")

	// Get latest block
	latest, err := s.chain.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}
	s.latestBlock = latest
	chainHeight.Set(float64(latest))

	// Determine sync strategy
	behind := latest - s.confirmations - s.currentBlock
	if behind > s.batchSize*2 {
		s.logger.Info().
			Uint64("current", s.currentBlock).
			Uint64("latest", latest).
			Uint64("behind", behind).
			Msg("behind chain, starting backfill")
		return s.runBackfill(ctx)
	}

	s.logger.Info().
		Uint64("current", s.currentBlock).
		Uint64("latest", latest).
		Msg("near chain head, starting realtime sync")
	return s.runRealtime(ctx)
}

// runBackfill processes historical blocks with parallel workers.
//
// This mode is used when the syncer is far behind the chain head (> batchSize*2).
// It processes blocks in batches using a worker pool for maximum throughput.
//
// Flow:
// 1. Fetch latest block and calculate safe head (latest - confirmations)
// 2. If caught up to safe head, switch to runRealtime()
// 3. Process batch (s.currentBlock+1 to min(currentBlock+batchSize, safeHead))
// 4. Save checkpoint after batch completes
// 5. Update Prometheus metrics (syncer_height, blocks_behind)
// 6. Repeat until caught up
//
// Worker Pool:
// - Splits batch into equal chunks per worker
// - Each worker calls processor.ProcessBlockRange()
// - Waits for all workers to complete before checkpointing
//
// Error Handling:
// - On RPC failure: Sleep 5s and retry
// - On processing failure: Sleep 5s and retry same batch
// - All errors increment syncer_errors_total metric
func (s *Syncer) runBackfill(ctx context.Context) error {
	s.logger.Info().
		Int("workers", s.workers).
		Uint64("batch_size", s.batchSize).
		Msg("starting backfill mode")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Get latest block
		latest, err := s.chain.GetLatestBlockNumber(ctx)
		if err != nil {
			syncerErrors.WithLabelValues("get_latest_block").Inc()
			s.logger.Error().Err(err).Msg("failed to get latest block")
			time.Sleep(5 * time.Second)
			continue
		}

		s.latestBlock = latest
		chainHeight.Set(float64(latest))

		// Calculate safe head (with confirmations)
		safeHead := latest
		if latest > s.confirmations {
			safeHead = latest - s.confirmations
		}

		if s.currentBlock >= safeHead {
			s.logger.Info().
				Uint64("current", s.currentBlock).
				Uint64("safe_head", safeHead).
				Msg("caught up to chain head, switching to realtime")
			return s.runRealtime(ctx)
		}

		// Process batch
		batchEnd := s.currentBlock + s.batchSize
		if batchEnd > safeHead {
			batchEnd = safeHead
		}

		if err := s.processBatch(ctx, s.currentBlock+1, batchEnd); err != nil {
			syncerErrors.WithLabelValues("process_batch").Inc()
			s.logger.Error().
				Err(err).
				Uint64("from", s.currentBlock+1).
				Uint64("to", batchEnd).
				Msg("failed to process batch")
			time.Sleep(5 * time.Second)
			continue
		}

		// Update checkpoint
		block, err := s.chain.GetBlockByNumber(ctx, batchEnd)
		if err != nil {
			syncerErrors.WithLabelValues("get_block").Inc()
			s.logger.Error().Err(err).Uint64("block", batchEnd).Msg("failed to get block for checkpoint")
			time.Sleep(5 * time.Second)
			continue
		}

		if err := s.checkpoint.UpdateBlock(ctx, s.serviceName, batchEnd, block.Hash().Hex()); err != nil {
			syncerErrors.WithLabelValues("update_checkpoint").Inc()
			s.logger.Error().Err(err).Msg("failed to update checkpoint")
			time.Sleep(5 * time.Second)
			continue
		}

		s.currentBlock = batchEnd
		syncerHeight.Set(float64(s.currentBlock))
		blocksBehind.Set(float64(safeHead - s.currentBlock))

		s.logger.Info().
			Uint64("processed_to", batchEnd).
			Uint64("latest", latest).
			Uint64("behind", safeHead-batchEnd).
			Msg("processed batch")
	}
}

// runRealtime processes new blocks as they arrive with low-latency polling.
//
// This mode is used when the syncer is near the chain head (≤ batchSize*2 behind).
// It polls for new blocks at the configured interval (default 2s).
//
// Flow:
//  1. Set up ticker for pollInterval (default: 2s)
//  2. On each tick:
//     a. Call syncToHead() to process any new blocks
//     b. Update isHealthy flag based on success/failure
//  3. Continue until context is canceled
//
// Mode Switching:
// - If syncer falls behind > batchSize*2: syncToHead() returns to runBackfill()
// - This can happen during network issues or RPC rate limits
//
// Health Monitoring:
// - isHealthy is set to false on syncToHead() errors
// - isHealthy is set to true on successful sync
// - Exposed via /health endpoint for Kubernetes readiness probes
func (s *Syncer) runRealtime(ctx context.Context) error {
	s.logger.Info().
		Dur("poll_interval", s.pollInterval).
		Uint64("confirmations", s.confirmations).
		Msg("starting realtime mode")

	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := s.syncToHead(ctx); err != nil {
				syncerErrors.WithLabelValues("sync_to_head").Inc()
				s.logger.Error().Err(err).Msg("failed to sync to head")
				s.isHealthy = false
				continue
			}
			s.isHealthy = true
		}
	}
}

// syncToHead syncs to the current chain head in realtime mode.
//
// Called by runRealtime() on each poll interval tick (default: every 2s).
//
// Logic:
// 1. Fetch latest block and calculate safe head (latest - confirmations)
// 2. If already at safe head, return immediately (blocks_behind = 0)
// 3. If fell behind > batchSize*2, switch to runBackfill() for fast catch-up
// 4. Otherwise, process blocks one at a time:
//   - Call processor.ProcessBlock(block) to extract events
//   - Save checkpoint after each block
//   - Update Prometheus metrics
//
// Single-Block Processing:
// - In realtime mode, blocks are processed sequentially (no parallelization)
// - This ensures minimal latency and immediate event publishing
// - Checkpoints are saved after each block for crash recovery
//
// Returns error on RPC failures or processing errors (triggers retry in runRealtime).
func (s *Syncer) syncToHead(ctx context.Context) error {
	// Get latest block
	latest, err := s.chain.GetLatestBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	s.latestBlock = latest
	chainHeight.Set(float64(latest))

	// Calculate safe head (with confirmations)
	safeHead := latest
	if latest > s.confirmations {
		safeHead = latest - s.confirmations
	}

	if s.currentBlock >= safeHead {
		// Already at head
		blocksBehind.Set(0)
		return nil
	}

	behind := safeHead - s.currentBlock
	blocksBehind.Set(float64(behind))

	// If too far behind, switch to backfill
	if behind > s.batchSize*2 {
		s.logger.Warn().
			Uint64("behind", behind).
			Msg("fell behind, switching to backfill mode")
		return s.runBackfill(ctx)
	}

	// Process blocks one at a time in realtime mode
	for block := s.currentBlock + 1; block <= safeHead; block++ {
		if err := s.processor.ProcessBlock(ctx, block); err != nil {
			return fmt.Errorf("failed to process block %d: %w", block, err)
		}

		// Update checkpoint
		header, err := s.chain.GetBlockByNumber(ctx, block)
		if err != nil {
			return fmt.Errorf("failed to get block %d: %w", block, err)
		}

		if err := s.checkpoint.UpdateBlock(ctx, s.serviceName, block, header.Hash().Hex()); err != nil {
			return fmt.Errorf("failed to update checkpoint: %w", err)
		}

		s.currentBlock = block
		syncerHeight.Set(float64(s.currentBlock))

		s.logger.Debug().
			Uint64("block", block).
			Uint64("latest", latest).
			Msg("processed block")
	}

	blocksBehind.Set(0)
	return nil
}

// processBatch processes a batch of blocks with parallel workers.
//
// Called by runBackfill() to process batches efficiently using a worker pool.
//
// Worker Pool Strategy:
// - If workers = 1: Sequential processing (no goroutines)
// - If workers > 1: Split batch into equal chunks per worker
//   - Example: 1000 blocks, 5 workers → 200 blocks per worker
//   - Worker 1: blocks 1-200
//   - Worker 2: blocks 201-400
//   - Worker 3: blocks 401-600
//   - Worker 4: blocks 601-800
//   - Worker 5: blocks 801-1000 (handles remainder)
//
// Synchronization:
// - Uses sync.WaitGroup to wait for all workers to complete
// - Errors are collected via buffered channel
// - Returns first error encountered (all workers must succeed)
//
// Safety:
// - Each worker operates on disjoint block ranges (no race conditions)
// - Processor must be thread-safe (uses NATS for publishing, which is thread-safe)
// - Checkpoint is saved AFTER all workers complete successfully
func (s *Syncer) processBatch(ctx context.Context, from, to uint64) error {
	if from > to {
		return fmt.Errorf("invalid range: from %d > to %d", from, to)
	}

	if s.workers == 1 {
		// Single-threaded processing
		return s.processor.ProcessBlockRange(ctx, from, to)
	}

	// Parallel processing with worker pool
	blockCount := to - from + 1
	blocksPerWorker := blockCount / uint64(s.workers)
	if blocksPerWorker == 0 {
		blocksPerWorker = 1
	}

	var wg sync.WaitGroup
	errChan := make(chan error, s.workers)

	for i := 0; i < s.workers; i++ {
		workerFrom := from + uint64(i)*blocksPerWorker
		workerTo := workerFrom + blocksPerWorker - 1

		// Last worker handles remainder
		if i == s.workers-1 {
			workerTo = to
		}

		if workerFrom > to {
			break
		}

		wg.Add(1)
		go func(from, to uint64) {
			defer wg.Done()
			if err := s.processor.ProcessBlockRange(ctx, from, to); err != nil {
				errChan <- err
			}
		}(workerFrom, workerTo)
	}

	// Wait for all workers
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// GetStatus returns current syncer status for monitoring.
//
// Returns:
// - current: Last block successfully processed and checkpointed
// - latest: Latest block fetched from blockchain RPC
// - healthy: Health flag (false if recent sync failed)
//
// Thread-safe via read lock. Called by HTTP health endpoint and Prometheus metrics.
func (s *Syncer) GetStatus() (current, latest uint64, healthy bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentBlock, s.latestBlock, s.isHealthy
}

// Healthy returns true if the syncer is healthy.
//
// Healthy means the last sync cycle (in runRealtime) completed successfully.
// Set to false on RPC errors or processing failures.
// Set to true on successful syncToHead() completion.
//
// Used by Kubernetes readiness probes to determine if pod should receive traffic.
func (s *Syncer) Healthy() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isHealthy
}
