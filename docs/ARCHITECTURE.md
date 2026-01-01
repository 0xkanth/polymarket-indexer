# Polymarket Indexer Architecture

This document explains the architecture, design decisions, and data flow of the Polymarket indexer.

## System Overview

The Polymarket indexer is a production-grade blockchain event indexer that tracks Polymarket prediction markets on Polygon. It follows an event-driven architecture with decoupled components communicating via NATS JetStream.

```
┌──────────────┐      ┌──────────────┐      ┌──────────────┐
│   Polygon    │ ───> │   Indexer    │ ───> │     NATS     │
│   Mainnet    │      │   Service    │      │  JetStream   │
└──────────────┘      └──────────────┘      └──────────────┘
                             │                      │
                             │                      ▼
                             ▼              ┌──────────────┐
                      ┌──────────────┐     │   Consumer   │
                      │    BoltDB    │     │   Service    │
                      │ (Checkpoint) │     └──────────────┘
                      └──────────────┘             │
                                                   ▼
                                            ┌──────────────┐
                                            │ TimescaleDB  │
                                            │ (Postgres)   │
                                            └──────────────┘
```

## Architecture Principles

### 1. **Separation of Concerns**
- **Indexer**: Fetches blockchain data, decodes events, publishes to NATS
- **Consumer**: Subscribes to NATS, writes to database
- **Benefits**: Independent scaling, fault isolation, easier testing

### 2. **Event-Driven Design**
- NATS JetStream provides durable event streaming
- At-least-once delivery guarantees
- Message deduplication prevents duplicate processing
- Decouples producer (indexer) from consumer (database writer)

### 3. **Idempotency**
- All operations are idempotent (safe to retry)
- Database uses `ON CONFLICT DO NOTHING` for upserts
- NATS message IDs prevent duplicate event publishing
- Checkpoint stores block hash for reorg detection

### 4. **Observability**
- Prometheus metrics for all components
- Structured JSON logging via zerolog
- Health check endpoints
- Distributed tracing ready (spans for block processing)

### 5. **Resilience**
- Automatic retry with exponential backoff
- Graceful degradation (continues processing on partial failures)
- Checkpoint-based recovery (resume from last known block)
- Reorg detection and handling

## Component Architecture

### Indexer Service

```
┌─────────────────────────────────────────────────┐
│              Indexer Service                    │
├─────────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌─────────────────┐ │
│  │  Main   │  │ Syncer  │  │   Processor     │ │
│  │ (Wire)  │──│         │──│                 │ │
│  └─────────┘  │ • Backfill│ │ • Block Fetch  │ │
│               │ • Realtime│ │ • Log Filter   │ │
│               │ • Worker  │ │ • Route Events │ │
│               │   Pool    │ │ • Publish NATS │ │
│               └─────────┘  └─────────────────┘ │
│                     │              │            │
│                     ▼              ▼            │
│  ┌─────────────────────────────────────────┐   │
│  │         Supporting Components           │   │
│  ├─────────────┬──────────────┬───────────┤   │
│  │ Chain Client│ Checkpoint   │  Router   │   │
│  │ • HTTP RPC  │ • BoltDB     │ • Sig Map │   │
│  │ • WS RPC    │ • Block Hash │ • Handler │   │
│  │ • Receipts  │ • ACID Txn   │   Registry│   │
│  └─────────────┴──────────────┴───────────┘   │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │         Event Handlers                  │   │
│  │ OrderFilled | TransferSingle |          │   │
│  │ ConditionPreparation | PositionSplit    │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

#### Key Components

**Syncer**
- Orchestrates block processing
- Two modes: Backfill (historical) and Realtime (new blocks)
- Auto-switches based on blocks behind chain head
- Parallel worker pool for backfill
- Sequential processing for realtime (maintains order)

**Processor**
- Fetches blocks via chain client
- Filters logs for monitored contracts
- Routes logs to event handlers
- Wraps parsed events in envelope
- Publishes to NATS with deduplication

**Chain Client**
- Dual RPC connections (HTTP + WebSocket)
- HTTP for historical data fetching
- WebSocket for realtime subscriptions
- Automatic reconnection handling
- Chain ID verification

**Router**
- Maps event signatures to handler functions
- Registry pattern for extensibility
- Returns parsed event payload

**Event Handlers**
- Decode ABI-encoded log data
- Parse indexed and non-indexed parameters
- Return strongly-typed event structs
- Handle complex types (arrays, nested structs)

**Checkpoint Store**
- Embedded BoltDB for persistence
- Stores block number + block hash
- ACID transactions
- Enables crash recovery
- Detects chain reorgs

### Consumer Service

```
┌─────────────────────────────────────────────────┐
│              Consumer Service                   │
├─────────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────────┐  ┌────────────┐  │
│  │  Main   │  │    NATS     │  │  Database  │  │
│  │ (Wire)  │──│  Consumer   │──│   Writer   │  │
│  └─────────┘  │             │  │            │  │
│               │ • Subscribe │  │ • Upsert   │  │
│               │ • Ack/Nak   │  │ • Batch    │  │
│               │ • Durable   │  │ • ON CONFLICT│ │
│               └─────────────┘  └────────────┘  │
│                                                 │
│  ┌─────────────────────────────────────────┐   │
│  │      Event Type Handlers                │   │
│  │ storeOrderFilled | storeTokenTransfer   │   │
│  │ storeConditionPreparation | ...         │   │
│  └─────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

#### Key Features

**NATS Consumer**
- Durable consumer (survives restarts)
- Explicit acknowledgment (manual control)
- Filter subject: `POLYMARKET.>`
- Max deliver: 3 (automatic retries)
- Ack wait: 30 seconds

**Database Writer**
- Connection pooling via pgx/v5
- Idempotent inserts (`ON CONFLICT DO NOTHING`)
- Event-specific table mapping
- Raw event + parsed event storage
- Batch operations for TransferBatch events

## Data Flow

### Indexing Flow

```
1. Syncer determines next block(s) to process
   ├─ Backfill: Batch of 100 blocks
   └─ Realtime: Single block at a time

2. Processor fetches block header and logs
   GET /eth_getBlockByNumber
   GET /eth_getLogs (filtered by contracts)

3. For each log:
   ├─ Router maps event signature to handler
   ├─ Handler decodes ABI-encoded data
   ├─ Returns parsed event struct
   └─ Processor wraps in Event envelope

4. Publish to NATS JetStream
   Subject: POLYMARKET.{EventName}.{ContractAddr}
   MessageID: {txHash}-{logIndex}
   Payload: JSON-encoded Event

5. Update checkpoint
   BoltDB.Put(serviceName, {blockNum, blockHash})
```

### Consumption Flow

```
1. Consumer subscribes to NATS stream
   Consumer: polymarket-consumer (durable)
   Filter: POLYMARKET.>

2. Receive message from NATS
   Extract event type from subject
   Unmarshal JSON to Event struct

3. Store raw event
   INSERT INTO events (...)
   ON CONFLICT (tx_hash, log_index) DO NOTHING

4. Store parsed event (type-specific)
   Case OrderFilled:
     INSERT INTO order_fills (...)
   Case TransferSingle:
     INSERT INTO token_transfers (...)
   etc.

5. Acknowledge message
   msg.Ack() → NATS marks as processed
   (or msg.Nak() on error for retry)
```

## Event Types

### CTF Exchange Events

**OrderFilled**
```solidity
event OrderFilled(
    bytes32 indexed orderHash,
    address indexed maker,
    address indexed taker,
    uint256 makerAssetId,
    uint256 takerAssetId,
    uint256 makerAmountFilled,
    uint256 takerAmountFilled,
    uint256 fee
)
```
- Tracks order execution
- Links maker and taker
- Records filled amounts and fees

**TokenRegistered**
```solidity
event TokenRegistered(
    uint256 indexed token0,
    uint256 indexed token1,
    bytes32 indexed conditionId
)
```
- New outcome tokens for a market
- Links tokens to condition

**OrderCancelled**
```solidity
event OrderCancelled(bytes32 indexed orderHash)
```
- Order cancellation tracking

### Conditional Tokens Events

**TransferSingle / TransferBatch**
```solidity
event TransferSingle(
    address indexed operator,
    address indexed from,
    address indexed to,
    uint256 id,
    uint256 value
)
```
- ERC-1155 token transfers
- Tracks outcome token trading

**ConditionPreparation**
```solidity
event ConditionPreparation(
    bytes32 indexed conditionId,
    address indexed oracle,
    bytes32 indexed questionId,
    uint256 outcomeSlotCount
)
```
- New market/condition creation
- Links to oracle and question

**ConditionResolution**
```solidity
event ConditionResolution(
    bytes32 indexed conditionId,
    address indexed oracle,
    bytes32 indexed questionId,
    uint256 outcomeSlotCount,
    uint256[] payoutNumerators
)
```
- Market resolution
- Payout distribution for outcomes

**PositionSplit / PositionsMerge**
```solidity
event PositionSplit(
    address indexed stakeholder,
    address collateralToken,
    bytes32 indexed parentCollectionId,
    bytes32 indexed conditionId,
    uint256[] partition,
    uint256 amount
)
```
- Minting outcome tokens (split)
- Redeeming outcome tokens (merge)

## Database Schema

### Core Tables

**events** (Hypertable - partitioned by time)
- All events in raw form
- JSONB payload for flexibility
- Primary source of truth
- Enables event replay

**order_fills** (Hypertable)
- Parsed OrderFilled events
- Optimized for trading analytics
- Indexed on maker, taker, timestamps

**token_transfers** (Hypertable)
- ERC-1155 transfers
- Tracks token movement
- Indexed on from, to, token_id

**conditions**
- Market definitions
- Oracle and question mapping
- Resolution status and payouts

**token_registrations**
- Outcome token registrations
- Links tokens to conditions

**position_splits / position_merges**
- Token minting and redemption
- Collateral tracking

### Continuous Aggregates

**order_volume_hourly**
```sql
SELECT 
  time_bucket('1 hour', block_timestamp) AS bucket,
  COUNT(*) AS trades,
  SUM(maker_amount_filled::numeric) AS volume
FROM order_fills
GROUP BY bucket
```
- Pre-computed hourly statistics
- Automatically refreshed
- Fast queries without aggregating millions of rows

**daily_active_traders**
```sql
SELECT 
  time_bucket('1 day', block_timestamp) AS bucket,
  maker AS trader,
  COUNT(*) AS trades
FROM order_fills
GROUP BY bucket, maker
UNION ALL
SELECT 
  time_bucket('1 day', block_timestamp) AS bucket,
  taker AS trader,
  COUNT(*) AS trades
FROM order_fills
GROUP BY bucket, taker
```
- Unique traders per day
- Fast analytics queries

## Configuration

### Indexer Config

```toml
[indexer]
start_block = 20558323        # Where to start indexing
batch_size = 100              # Blocks per batch (backfill)
poll_interval = "6s"          # Realtime polling frequency
confirmations = 100           # Blocks to wait before processing
workers = 4                   # Parallel workers for backfill

[chain]
http_url = "https://..."      # Primary RPC
ws_url = "wss://..."          # WebSocket RPC (optional)
chain_id = 137                # Polygon mainnet

[nats]
url = "nats://localhost:4222"
stream_name = "POLYMARKET"
subjects = ["POLYMARKET.>"]
max_age = "24h"               # Message retention
```

### Key Design Decisions

**100 Block Confirmations**
- Polygon has occasional deep reorgs
- 100 blocks ≈ 3-4 minutes delay
- Provides safety against finality issues
- Configurable based on risk tolerance

**20-Minute NATS Deduplication Window**
- Prevents duplicate events on restarts
- Covers typical recovery scenarios
- Balances memory usage vs safety

**Batch Size: 100 Blocks**
- Optimal for Polygon RPC rate limits
- Balances throughput and memory
- Can be adjusted based on event density

**4 Workers for Backfill**
- Utilizes CPU cores efficiently
- Avoids overwhelming RPC endpoints
- Can scale up with dedicated RPC

## Monitoring Metrics

### Indexer Metrics

- `polymarket_blocks_processed_total` - Cumulative blocks indexed
- `polymarket_events_processed_total{event_type}` - Events by type
- `polymarket_syncer_block_height` - Current block being processed
- `polymarket_chain_block_height` - Latest block on chain
- `polymarket_blocks_behind` - How far behind chain head
- `polymarket_block_processing_duration_seconds` - Processing time per block
- `polymarket_processing_errors_total{error_type}` - Error counts

### Consumer Metrics

- `polymarket_events_consumed_total{event_type}` - NATS messages consumed
- `polymarket_events_stored_total{event_type}` - DB inserts completed
- `polymarket_consume_errors_total{error_type}` - Consumer errors
- `polymarket_consumer_lag_seconds` - Time from event to DB write

### Alert Thresholds

```yaml
# Indexer stopped progressing
polymarket_blocks_behind > 1000 for 5 minutes

# Consumer lag too high
polymarket_consumer_lag_seconds > 300 for 5 minutes

# Error rate too high
rate(polymarket_processing_errors_total[5m]) > 10
```

## Scaling Strategies

### Vertical Scaling
- Increase workers for faster backfill
- More CPU cores = more parallel processing
- More memory = larger batches

### Horizontal Scaling
- Run multiple consumers (DB writers)
- Different consumers for different event types
- Load balance with NATS queue groups

### Data Partitioning
- TimescaleDB chunks by time (1 day default)
- Old data compressed automatically
- Drop old chunks for retention policy

### Read Scaling
- TimescaleDB read replicas
- Continuous aggregates for analytics
- Redis cache for hot queries

## Reorg Handling

### Detection
```go
// Compare stored block hash with fetched block hash
checkpoint := checkpoint.GetCheckpoint()
block := chain.GetBlockByNumber(checkpoint.BlockNumber)

if block.Hash() != checkpoint.BlockHash {
    // Reorg detected!
    rollback()
}
```

### Recovery
1. Stop processing new blocks
2. Roll back to last known good block (N-100)
3. Delete affected events from database
4. Clear NATS deduplication window
5. Resume indexing from safe block

## Testing Strategy

### Unit Tests
- Handler ABI decoding
- Router event mapping
- Checkpoint operations
- Event serialization

### Integration Tests
- Chain client with mock RPC
- NATS pub/sub flow
- Database insert operations
- End-to-end event processing

### E2E Tests
- Deploy contracts on local chain
- Emit test events
- Verify indexing and storage
- Test reorg scenarios

## Security Considerations

### RPC Security
- Use authenticated endpoints
- Rate limit protection
- Multiple RPC fallbacks
- Monitor for suspicious responses

### Database Security
- Parameterized queries (prevent SQL injection)
- Least privilege principle
- SSL connections in production
- Audit logging

### NATS Security
- TLS encryption
- Authentication with JWT
- Subject-based authorization
- Network segmentation

## Performance Optimization

### Database
- Batch inserts where possible
- Appropriate indexes (not too many)
- Regular VACUUM and ANALYZE
- Compression policies for old data

### NATS
- Async publishing (don't block on ack)
- Message batching for high throughput
- Stream limits to prevent unbounded growth

### Chain Client
- Connection pooling
- Request caching for repeated queries
- WebSocket for subscriptions (lower latency)
- Parallel log fetching for multiple contracts

## Future Enhancements

1. **GraphQL API** - Query layer on TimescaleDB
2. **Webhooks** - Real-time notifications for events
3. **Event Replay** - Reprocess historical events
4. **Multi-chain** - Support other Polymarket deployments
5. **Advanced Analytics** - ML predictions, anomaly detection
6. **API Rate Limiting** - Protect query endpoints
7. **Data Warehouse** - Export to S3/BigQuery for analytics

## References

- [eth-tracker](../reference/eth-tracker) - NATS patterns inspiration
- [evm-scanner](../reference/evm-scanner) - RPC client patterns
- [TimescaleDB Docs](https://docs.timescale.com)
- [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
- [Polymarket Docs](https://docs.polymarket.com)
