// Package db provides database abstractions for checkpoint storage.
package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xkanth/polymarket-indexer/pkg/models"
	"go.etcd.io/bbolt"
)

const (
	// checkpointBucket is the BoltDB bucket name for storing checkpoints
	checkpointBucket = "checkpoints"
)

// CheckpointDB provides checkpoint persistence using BoltDB.
type CheckpointDB struct {
	db *bbolt.DB
}

// NewCheckpointDB creates a new checkpoint database instance.
func NewCheckpointDB(dbPath string) (*CheckpointDB, error) {
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open checkpoint db: %w", err)
	}

	// Create bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(checkpointBucket))
		return err
	})
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create checkpoint bucket: %w", err)
	}

	return &CheckpointDB{db: db}, nil
}

// SaveCheckpoint saves or updates a checkpoint for a service.
func (c *CheckpointDB) SaveCheckpoint(ctx context.Context, checkpoint models.Checkpoint) error {
	checkpoint.UpdatedAt = time.Now()

	return c.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(checkpointBucket))
		if b == nil {
			return fmt.Errorf("checkpoint bucket not found")
		}

		data, err := json.Marshal(checkpoint)
		if err != nil {
			return fmt.Errorf("failed to marshal checkpoint: %w", err)
		}

		return b.Put([]byte(checkpoint.ServiceName), data)
	})
}

// GetCheckpoint retrieves a checkpoint for a service.
func (c *CheckpointDB) GetCheckpoint(ctx context.Context, serviceName string) (*models.Checkpoint, error) {
	var checkpoint models.Checkpoint

	err := c.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(checkpointBucket))
		if b == nil {
			return fmt.Errorf("checkpoint bucket not found")
		}

		data := b.Get([]byte(serviceName))
		if data == nil {
			return fmt.Errorf("checkpoint not found for service: %s", serviceName)
		}

		return json.Unmarshal(data, &checkpoint)
	})

	if err != nil {
		return nil, err
	}

	return &checkpoint, nil
}

// GetOrCreateCheckpoint gets an existing checkpoint or creates a new one with the start block.
func (c *CheckpointDB) GetOrCreateCheckpoint(ctx context.Context, serviceName string, startBlock uint64) (*models.Checkpoint, error) {
	checkpoint, err := c.GetCheckpoint(ctx, serviceName)
	if err == nil {
		return checkpoint, nil
	}

	// Create new checkpoint
	checkpoint = &models.Checkpoint{
		ServiceName:   serviceName,
		LastBlock:     startBlock,
		LastBlockHash: "0x0000000000000000000000000000000000000000000000000000000000000000",
		UpdatedAt:     time.Now(),
	}

	if err := c.SaveCheckpoint(ctx, *checkpoint); err != nil {
		return nil, fmt.Errorf("failed to create checkpoint: %w", err)
	}

	return checkpoint, nil
}

// UpdateBlock updates just the block number and hash in the checkpoint.
func (c *CheckpointDB) UpdateBlock(ctx context.Context, serviceName string, blockNumber uint64, blockHash string) error {
	checkpoint, err := c.GetCheckpoint(ctx, serviceName)
	if err != nil {
		return err
	}

	checkpoint.LastBlock = blockNumber
	checkpoint.LastBlockHash = blockHash

	return c.SaveCheckpoint(ctx, *checkpoint)
}

// Close closes the database connection.
func (c *CheckpointDB) Close() error {
	return c.db.Close()
}

// Stats returns database statistics.
func (c *CheckpointDB) Stats() bbolt.Stats {
	return c.db.Stats()
}
