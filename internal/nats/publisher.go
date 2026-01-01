// Package nats provides NATS JetStream publishing functionality.
package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xkanth/polymarket-indexer/pkg/models"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog"
)

const (
	// streamName is the NATS JetStream stream name
	streamName = "POLYMARKET"

	// streamSubjectPattern is the subject pattern for all Polymarket events
	streamSubjectPattern = "POLYMARKET.*"

	// streamCreateTimeout is the timeout for stream creation
	streamCreateTimeout = 10 * time.Second
)

// Publisher publishes events to NATS JetStream with deduplication.
type Publisher struct {
	js     jetstream.JetStream
	nc     *nats.Conn
	logger *zerolog.Logger
	prefix string
}

// NewPublisher creates a new NATS JetStream publisher.
func NewPublisher(natsURL string, persistDuration time.Duration, subjectPrefix string, logger *zerolog.Logger) (*Publisher, error) {
	// Connect to NATS
	nc, err := nats.Connect(natsURL,
		nats.Name("polymarket-indexer"),
		nats.MaxReconnects(-1), // Unlimited reconnects
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				logger.Error().Err(err).Msg("nats disconnected")
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.Info().Msg("nats reconnected")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	// Create or update the stream
	ctx, cancel := context.WithTimeout(context.Background(), streamCreateTimeout)
	defer cancel()

	duplicateWindow := 20 * time.Minute
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:       streamName,
		Subjects:   []string{streamSubjectPattern},
		MaxAge:     persistDuration,
		Storage:    jetstream.FileStorage,
		Duplicates: duplicateWindow,
		Retention:  jetstream.LimitsPolicy,
	})
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	logger.Info().
		Str("stream", streamName).
		Str("subjects", streamSubjectPattern).
		Dur("max_age", persistDuration).
		Dur("duplicate_window", duplicateWindow).
		Msg("NATS publisher initialized")

	return &Publisher{
		js:     js,
		nc:     nc,
		logger: logger,
		prefix: subjectPrefix,
	}, nil
}

// Publish publishes an event to NATS JetStream with deduplication.
// The message ID is constructed from txHash and logIndex to prevent duplicates.
func (p *Publisher) Publish(ctx context.Context, event models.Event) error {
	// Construct subject: POLYMARKET.{EventName}.{ContractAddress}
	subject := fmt.Sprintf("%s.%s.%s", p.prefix, event.EventName, event.ContractAddr)

	// Marshal event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Create message ID for deduplication: txHash-logIndex
	msgID := fmt.Sprintf("%s-%d", event.TxHash, event.LogIndex)

	// Publish with deduplication
	_, err = p.js.Publish(ctx, subject, data, jetstream.WithMsgID(msgID))
	if err != nil {
		p.logger.Error().
			Err(err).
			Str("subject", subject).
			Str("msg_id", msgID).
			Uint64("block", event.Block).
			Msg("failed to publish event")
		return fmt.Errorf("failed to publish to NATS: %w", err)
	}

	p.logger.Debug().
		Str("subject", subject).
		Str("event", event.EventName).
		Uint64("block", event.Block).
		Str("tx", event.TxHash).
		Msg("event published")

	return nil
}

// PublishBatch publishes multiple events in a batch for better performance.
func (p *Publisher) PublishBatch(ctx context.Context, events []models.Event) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the NATS connection.
func (p *Publisher) Close() {
	if p.nc != nil {
		p.nc.Close()
		p.logger.Info().Msg("NATS publisher closed")
	}
}

// Healthy checks if the NATS connection is healthy.
func (p *Publisher) Healthy() bool {
	return p.nc != nil && p.nc.IsConnected()
}
