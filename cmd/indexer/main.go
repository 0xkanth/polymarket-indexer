// Main indexer service.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/0xkanth/polymarket-indexer/internal/chain"
	"github.com/0xkanth/polymarket-indexer/internal/db"
	"github.com/0xkanth/polymarket-indexer/internal/nats"
	"github.com/0xkanth/polymarket-indexer/internal/processor"
	"github.com/0xkanth/polymarket-indexer/internal/syncer"
	"github.com/0xkanth/polymarket-indexer/internal/util"
	"github.com/0xkanth/polymarket-indexer/pkg/config"
)

const (
	serviceName = "polymarket-indexer"
)

func main() {
	// Initialize logger
	logger := util.InitLogger()
	logger.Info().Msg("starting polymarket indexer")

	// Load configuration
	cfg := util.InitConfig(logger, "config.toml")

	// Update log level from config
	util.UpdateLogLevel(cfg, logger)

	// Load chain configuration from chains.json
	chainConfigs, err := config.LoadConfig("config/chains.json")
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load chains.json")
	}

	// Get selected chain (from config.toml)
	chainName := cfg.String("chain.name")
	selectedChain, err := chainConfigs.GetChain(chainName)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("chain", chainName).
			Msg("chain not found in chains.json")
	}

	logger.Info().
		Str("chain", selectedChain.Name).
		Int64("chain_id", selectedChain.ChainID).
		Strs("rpc_urls", selectedChain.RPCUrls).
		Strs("contracts", selectedChain.GetAllContractAddressStrings()).
		Uint64("start_block", selectedChain.StartBlock).
		Int("confirmations", selectedChain.Confirmations).
		Msg("loaded chain configuration")

	// Initialize chain client
	httpURL := selectedChain.RPCUrls[0]
	wsURL := ""
	if len(selectedChain.WSUrls) > 0 {
		wsURL = selectedChain.WSUrls[0]
	}

	chainClient, err := chain.NewClient(
		httpURL,
		wsURL,
		selectedChain.ChainID,
		logger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create chain client")
	}
	logger.Info().
		Str("http", httpURL).
		Str("ws", wsURL).
		Int64("chain_id", selectedChain.ChainID).
		Msg("initialized chain client")

	// Initialize checkpoint store
	checkpointStore, err := db.NewCheckpointDB(cfg.String("db.checkpoint_path"))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create checkpoint store")
	}
	defer checkpointStore.Close()
	logger.Info().
		Str("path", cfg.String("db.checkpoint_path")).
		Msg("initialized checkpoint store")

	// Initialize NATS publisher
	publisher, err := nats.NewPublisher(
		cfg.String("nats.url"),
		cfg.Duration("nats.max_age"),
		cfg.String("nats.stream_name"),
		logger,
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create nats publisher")
	}
	defer publisher.Close()
	logger.Info().
		Str("url", cfg.String("nats.url")).
		Str("stream", cfg.String("nats.stream_name")).
		Msg("initialized nats publisher")

	// Initialize processor
	proc, err := processor.New(
		*logger,
		chainClient,
		publisher,
		processor.BlockEventProcessingConfig{
			Contracts:  selectedChain.GetAllContractAddressStrings(),
			StartBlock: selectedChain.StartBlock,
		},
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create processor")
	}
	logger.Info().
		Strs("contracts", selectedChain.GetAllContractAddressStrings()).
		Uint64("start_block", selectedChain.StartBlock).
		Msg("initialized processor")

	// Initialize syncer
	sync := syncer.New(
		*logger,
		chainClient,
		proc,
		checkpointStore,
		syncer.Config{
			ServiceName:   serviceName,
			StartBlock:    selectedChain.StartBlock,
			BatchSize:     uint64(cfg.Int64("indexer.batch_size")),
			PollInterval:  cfg.Duration("indexer.poll_interval"),
			Confirmations: uint64(selectedChain.Confirmations),
			Workers:       cfg.Int("indexer.workers"),
		},
	)
	logger.Info().
		Uint64("batch_size", uint64(cfg.Int64("indexer.batch_size"))).
		Dur("poll_interval", cfg.Duration("indexer.poll_interval")).
		Uint64("confirmations", uint64(selectedChain.Confirmations)).
		Int("workers", cfg.Int("indexer.workers")).
		Msg("initialized syncer")

	// Start metrics server
	metricsAddr := cfg.String("metrics.address")
	metricsServer := &http.Server{
		Addr:    metricsAddr,
		Handler: promhttp.Handler(),
	}

	go func() {
		logger.Info().Str("address", metricsAddr).Msg("starting metrics server")
		if err := metricsServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("metrics server error")
		}
	}()

	// Start health check server
	healthAddr := cfg.String("health.address")
	healthServer := &http.Server{
		Addr:    healthAddr,
		Handler: http.HandlerFunc(healthCheckHandler(sync, publisher)),
	}

	go func() {
		logger.Info().Str("address", healthAddr).Msg("starting health check server")
		if err := healthServer.ListenAndServe(); err != http.ErrServerClosed {
			logger.Error().Err(err).Msg("health check server error")
		}
	}()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start syncer in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- sync.Start(ctx)
	}()

	// Wait for shutdown signal or error
	select {
	case sig := <-sigChan:
		logger.Info().Str("signal", sig.String()).Msg("received shutdown signal")
	case err := <-errChan:
		if err != nil {
			logger.Error().Err(err).Msg("syncer error")
		}
	}

	// Graceful shutdown
	logger.Info().Msg("shutting down")
	cancel()

	// Shutdown metrics server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("metrics server shutdown error")
	}

	if err := healthServer.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("health server shutdown error")
	}

	logger.Info().Msg("shutdown complete")
}

// healthCheckHandler returns a health check handler.
func healthCheckHandler(sync *syncer.Syncer, pub *nats.Publisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sync.Healthy() || !pub.Healthy() {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "unhealthy\n")
			return
		}

		current, latest, _ := sync.GetStatus()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "healthy\ncurrent: %d\nlatest: %d\nbehind: %d\n",
			current, latest, latest-current)
	}
}
