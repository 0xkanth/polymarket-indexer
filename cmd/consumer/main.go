// Consumer service - reads events from NATS and writes to TimescaleDB.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"

	"github.com/0xkanth/polymarket-indexer/internal/util"
	"github.com/0xkanth/polymarket-indexer/pkg/models"
)

var (
	eventsConsumed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polymarket_events_consumed_total",
		Help: "Total number of events consumed from NATS",
	}, []string{"event_type"})

	eventsStored = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polymarket_events_stored_total",
		Help: "Total number of events stored in database",
	}, []string{"event_type"})

	consumeErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "polymarket_consume_errors_total",
		Help: "Total number of consume errors",
	}, []string{"error_type"})

	processingLag = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "polymarket_consumer_lag_seconds",
		Help: "Time lag between event occurrence and processing",
	})
)

const (
	serviceName = "polymarket-consumer"
)

func main() {
	// Initialize logger
	logger := util.InitLogger()
	logger.Info().Msg("starting polymarket consumer")

	// Load configuration
	cfg := util.InitConfig(logger, "config.toml")

	// Update log level from config
	util.UpdateLogLevel(cfg, logger)

	// Connect to PostgreSQL
	dbConfig := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.String("postgres.host"),
		cfg.Int("postgres.port"),
		cfg.String("postgres.user"),
		cfg.String("postgres.password"),
		cfg.String("postgres.database"),
		cfg.String("postgres.sslmode"),
	)

	pool, err := pgxpool.New(context.Background(), dbConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}
	logger.Info().
		Str("host", cfg.String("postgres.host")).
		Str("database", cfg.String("postgres.database")).
		Msg("connected to database")

	// Connect to NATS
	nc, err := nats.Connect(cfg.String("nats.url"))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to nats")
	}
	defer nc.Close()
	logger.Info().Str("url", cfg.String("nats.url")).Msg("connected to nats")

	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create jetstream context")
	}

	// Create durable consumer
	streamName := cfg.String("nats.stream_name")
	consumerName := cfg.String("nats.consumer_name")

	consumer, err := js.CreateOrUpdateConsumer(context.Background(), streamName, jetstream.ConsumerConfig{
		Name:          consumerName,
		Durable:       consumerName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    3,
		AckWait:       30 * time.Second,
		FilterSubject: "POLYMARKET.>",
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create consumer")
	}
	logger.Info().
		Str("stream", streamName).
		Str("consumer", consumerName).
		Msg("created consumer")

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

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consuming messages
	consCtx, err := consumer.Consume(func(msg jetstream.Msg) {
		if err := processMessage(ctx, pool, msg, *logger); err != nil {
			consumeErrors.WithLabelValues("process_message").Inc()
			logger.Error().Err(err).Str("subject", msg.Subject()).Msg("failed to process message")
			// Negative acknowledgment to retry
			msg.Nak()
			return
		}
		// Acknowledge message
		msg.Ack()
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start consuming")
	}
	defer consCtx.Stop()

	logger.Info().Msg("consumer started, waiting for messages")

	// Wait for shutdown signal
	sig := <-sigChan
	logger.Info().Str("signal", sig.String()).Msg("received shutdown signal")

	// Graceful shutdown
	logger.Info().Msg("shutting down")
	cancel()

	// Shutdown metrics server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("metrics server shutdown error")
	}

	logger.Info().Msg("shutdown complete")
}

// processMessage processes a single NATS message.
func processMessage(ctx context.Context, pool *pgxpool.Pool, msg jetstream.Msg, logger zerolog.Logger) error {
	// Parse event
	var event models.Event
	if err := json.Unmarshal(msg.Data(), &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	// Calculate processing lag
	eventTime := time.Unix(int64(event.Timestamp), 0)
	lag := time.Since(eventTime)
	processingLag.Set(lag.Seconds())

	// Extract event type from subject (POLYMARKET.{EventType}.{ContractAddress})
	eventType := extractEventType(msg.Subject())
	eventsConsumed.WithLabelValues(eventType).Inc()

	logger.Debug().
		Str("event", eventType).
		Uint64("block", event.Block).
		Str("tx", event.TxHash).
		Msg("processing event")

	// Store event in appropriate table based on type
	if err := storeEvent(ctx, pool, eventType, event); err != nil {
		return fmt.Errorf("failed to store event: %w", err)
	}

	eventsStored.WithLabelValues(eventType).Inc()
	return nil
}

// extractEventType extracts event type from NATS subject.
func extractEventType(subject string) string {
	// Subject format: POLYMARKET.{EventType}.{ContractAddress}
	// Extract middle part
	parts := []byte(subject)
	firstDot := -1
	secondDot := -1
	for i, b := range parts {
		if b == '.' {
			if firstDot == -1 {
				firstDot = i
			} else {
				secondDot = i
				break
			}
		}
	}
	if firstDot >= 0 && secondDot > firstDot {
		return subject[firstDot+1 : secondDot]
	}
	return "Unknown"
}

// storeEvent stores an event in the database.
func storeEvent(ctx context.Context, pool *pgxpool.Pool, eventType string, event models.Event) error {
	// Store raw event
	if err := storeRawEvent(ctx, pool, event); err != nil {
		return fmt.Errorf("failed to store raw event: %w", err)
	}

	// Store parsed event based on type
	switch eventType {
	case "OrderFilled":
		return storeOrderFilled(ctx, pool, event)
	case "TokenRegistered":
		return storeTokenRegistered(ctx, pool, event)
	case "TransferSingle":
		return storeTokenTransfer(ctx, pool, event)
	case "TransferBatch":
		return storeTokenTransferBatch(ctx, pool, event)
	case "ConditionPreparation":
		return storeConditionPreparation(ctx, pool, event)
	case "ConditionResolution":
		return storeConditionResolution(ctx, pool, event)
	case "PositionSplit":
		return storePositionSplit(ctx, pool, event)
	case "PositionsMerge":
		return storePositionsMerge(ctx, pool, event)
	default:
		// Unknown event type, already stored as raw event
		return nil
	}
}

// storeRawEvent stores the raw event in the events table.
func storeRawEvent(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, err := json.Marshal(event.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	query := `
		INSERT INTO events (
			block_number, block_hash, block_timestamp, transaction_hash, log_index,
			contract_address, event_signature, payload
		) VALUES ($1, $2, to_timestamp($3), $4, $5, $6, $7, $8)
		ON CONFLICT (transaction_hash, log_index) DO NOTHING
	`

	_, err = pool.Exec(ctx, query,
		event.Block,
		event.BlockHash,
		event.Timestamp,
		event.TxHash,
		event.LogIndex,
		event.ContractAddr,
		event.EventSig,
		payloadJSON,
	)

	return err
}

// storeOrderFilled stores an OrderFilled event.
func storeOrderFilled(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var order models.OrderFilled
	if err := json.Unmarshal(payloadJSON, &order); err != nil {
		return err
	}

	query := `
		INSERT INTO order_fills (
			block_number, block_timestamp, transaction_hash, log_index,
			order_hash, maker, taker, maker_asset_id, taker_asset_id,
			maker_amount_filled, taker_amount_filled, fee
		) VALUES ($1, to_timestamp($2), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (transaction_hash, log_index) DO NOTHING
	`

	_, err := pool.Exec(ctx, query,
		event.Block,
		event.Timestamp,
		event.TxHash,
		event.LogIndex,
		order.OrderHash,
		order.Maker,
		order.Taker,
		order.MakerAssetID.String(),
		order.TakerAssetID.String(),
		order.MakerAmountFilled.String(),
		order.TakerAmountFilled.String(),
		order.Fee.String(),
	)

	return err
}

// storeTokenRegistered stores a TokenRegistered event.
func storeTokenRegistered(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var token models.TokenRegistered
	if err := json.Unmarshal(payloadJSON, &token); err != nil {
		return err
	}

	query := `
		INSERT INTO token_registrations (
			block_number, block_timestamp, transaction_hash, log_index,
			token0, token1, condition_id
		) VALUES ($1, to_timestamp($2), $3, $4, $5, $6, $7)
		ON CONFLICT (transaction_hash, log_index) DO NOTHING
	`

	_, err := pool.Exec(ctx, query,
		event.Block,
		event.Timestamp,
		event.TxHash,
		event.LogIndex,
		token.Token0.String(),
		token.Token1.String(),
		token.ConditionID,
	)

	return err
}

// storeTokenTransfer stores a TransferSingle event.
func storeTokenTransfer(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var transfer models.TransferSingle
	if err := json.Unmarshal(payloadJSON, &transfer); err != nil {
		return err
	}

	query := `
		INSERT INTO token_transfers (
			block_number, block_timestamp, transaction_hash, log_index,
			operator, from_address, to_address, token_id, amount
		) VALUES ($1, to_timestamp($2), $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (transaction_hash, log_index) DO NOTHING
	`

	_, err := pool.Exec(ctx, query,
		event.Block,
		event.Timestamp,
		event.TxHash,
		event.LogIndex,
		transfer.Operator,
		transfer.From,
		transfer.To,
		transfer.TokenID.String(),
		transfer.Amount.String(),
	)

	return err
}

// storeTokenTransferBatch stores TransferBatch events (creates multiple records).
func storeTokenTransferBatch(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var transfer models.TransferBatch
	if err := json.Unmarshal(payloadJSON, &transfer); err != nil {
		return err
	}

	// Insert each token transfer separately
	for i := range transfer.TokenIDs {
		query := `
			INSERT INTO token_transfers (
				block_number, block_timestamp, transaction_hash, log_index,
				operator, from_address, to_address, token_id, amount
			) VALUES ($1, to_timestamp($2), $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (transaction_hash, log_index, token_id) DO NOTHING
		`

		if _, err := pool.Exec(ctx, query,
			event.Block,
			event.Timestamp,
			event.TxHash,
			event.LogIndex,
			transfer.Operator,
			transfer.From,
			transfer.To,
			transfer.TokenIDs[i].String(),
			transfer.Amounts[i].String(),
		); err != nil {
			return err
		}
	}

	return nil
}

// storeConditionPreparation stores a ConditionPreparation event.
func storeConditionPreparation(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var condition models.ConditionPreparation
	if err := json.Unmarshal(payloadJSON, &condition); err != nil {
		return err
	}

	query := `
		INSERT INTO conditions (
			condition_id, oracle, question_id, outcome_slot_count,
			block_number, block_timestamp, transaction_hash
		) VALUES ($1, $2, $3, $4, $5, to_timestamp($6), $7)
		ON CONFLICT (condition_id) DO NOTHING
	`

	_, err := pool.Exec(ctx, query,
		condition.ConditionID,
		condition.Oracle,
		condition.QuestionID,
		condition.OutcomeSlotCount,
		event.Block,
		event.Timestamp,
		event.TxHash,
	)

	return err
}

// storeConditionResolution stores a ConditionResolution event.
func storeConditionResolution(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var resolution models.ConditionResolution
	if err := json.Unmarshal(payloadJSON, &resolution); err != nil {
		return err
	}

	// Convert payout numerators to string array
	payouts := make([]string, len(resolution.PayoutNumerators))
	for i, p := range resolution.PayoutNumerators {
		payouts[i] = p.String()
	}

	query := `
		UPDATE conditions
		SET resolved = true,
		    payout_numerators = $1,
		    resolution_block = $2,
		    resolution_timestamp = to_timestamp($3),
		    resolution_tx = $4
		WHERE condition_id = $5
	`

	_, err := pool.Exec(ctx, query,
		payouts,
		event.Block,
		event.Timestamp,
		event.TxHash,
		resolution.ConditionID,
	)

	return err
}

// storePositionSplit stores a PositionSplit event.
func storePositionSplit(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var split models.PositionSplit
	if err := json.Unmarshal(payloadJSON, &split); err != nil {
		return err
	}

	partition := make([]string, len(split.Partition))
	for i, p := range split.Partition {
		partition[i] = p.String()
	}

	query := `
		INSERT INTO position_splits (
			block_number, block_timestamp, transaction_hash, log_index,
			stakeholder, collateral_token, parent_collection_id, condition_id,
			partition, amount
		) VALUES ($1, to_timestamp($2), $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (transaction_hash, log_index) DO NOTHING
	`

	_, err := pool.Exec(ctx, query,
		event.Block,
		event.Timestamp,
		event.TxHash,
		event.LogIndex,
		split.Stakeholder,
		split.CollateralToken,
		split.ParentCollectionID,
		split.ConditionID,
		partition,
		split.Amount.String(),
	)

	return err
}

// storePositionsMerge stores a PositionsMerge event.
func storePositionsMerge(ctx context.Context, pool *pgxpool.Pool, event models.Event) error {
	payloadJSON, _ := json.Marshal(event.Payload)
	var merge models.PositionsMerge
	if err := json.Unmarshal(payloadJSON, &merge); err != nil {
		return err
	}

	partition := make([]string, len(merge.Partition))
	for i, p := range merge.Partition {
		partition[i] = p.String()
	}

	query := `
		INSERT INTO position_merges (
			block_number, block_timestamp, transaction_hash, log_index,
			stakeholder, collateral_token, parent_collection_id, condition_id,
			partition, amount
		) VALUES ($1, to_timestamp($2), $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (transaction_hash, log_index) DO NOTHING
	`

	_, err := pool.Exec(ctx, query,
		event.Block,
		event.Timestamp,
		event.TxHash,
		event.LogIndex,
		merge.Stakeholder,
		merge.CollateralToken,
		merge.ParentCollectionID,
		merge.ConditionID,
		partition,
		merge.Amount.String(),
	)

	return err
}

// bigIntFromString parses a big.Int from string.
func bigIntFromString(s string) *big.Int {
	n := new(big.Int)
	n.SetString(s, 10)
	return n
}
