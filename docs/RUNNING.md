# Running the Polymarket Indexer

This guide will walk you through running the Polymarket indexer from scratch.

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Make (optional, for convenience)

## Quick Start with Docker Compose

The easiest way to run the full stack:

```bash
# Start all services (NATS, TimescaleDB, Indexer, Consumer)
docker-compose up -d

# View logs
docker-compose logs -f indexer
docker-compose logs -f consumer

# Stop all services
docker-compose down
```

This will start:
- **NATS JetStream** on port 4222 (HTTP monitoring: 8222)
- **TimescaleDB** on port 5432
- **Indexer service** (metrics on :9090, health on :8080)
- **Consumer service** (metrics on :9091)
- **Prometheus** on port 9090
- **Grafana** on port 3000

## Step-by-Step Setup

### 1. Configure the Indexer

Edit `config.toml` to set your RPC endpoints:

```toml
[chain]
http_url = "https://polygon-rpc.com"  # Your Polygon RPC
ws_url = "wss://polygon-rpc.com"      # WebSocket RPC (optional)
chain_id = 137
```

**Important**: Use a reliable RPC provider (Alchemy, Infura, QuickNode) for production.

### 2. Start Infrastructure Services

```bash
# Start NATS and TimescaleDB only
docker-compose up -d nats timescaledb

# Wait for services to be ready
sleep 5

# Verify NATS is running
curl http://localhost:8222/healthz

# Verify PostgreSQL is running
docker-compose exec timescaledb psql -U polymarket -d polymarket -c "SELECT version();"
```

### 3. Run Database Migrations

```bash
# Install migrate CLI (if not already installed)
# macOS
brew install golang-migrate

# Or download from: https://github.com/golang-migrate/migrate

# Run migrations
migrate -path ./migrations -database "postgresql://polymarket:polymarket@localhost:5432/polymarket?sslmode=disable" up

# Verify tables were created
docker-compose exec timescaledb psql -U polymarket -d polymarket -c "\dt"
```

Expected tables:
- `events` (hypertable)
- `order_fills` (hypertable)
- `token_transfers` (hypertable)
- `token_registrations`
- `conditions`
- `position_splits`
- `position_merges`

### 4. Build the Services

```bash
# Build indexer
go build -o bin/indexer ./cmd/indexer

# Build consumer
go build -o bin/consumer ./cmd/consumer

# Or use Make
make build
```

### 5. Run the Indexer

```bash
# Run indexer (foreground)
./bin/indexer

# Or with Make
make run-indexer

# Check health
curl http://localhost:8080/health

# Check metrics
curl http://localhost:9090/metrics
```

Expected log output:
```
{"level":"info","time":"...","message":"starting polymarket indexer"}
{"level":"info","message":"initialized chain client","http":"https://polygon-rpc.com","chain_id":137}
{"level":"info","message":"initialized checkpoint store","path":"data/checkpoints.db"}
{"level":"info","message":"loaded checkpoint","checkpoint":20558323}
{"level":"info","message":"behind chain, starting backfill","behind":100000}
{"level":"info","message":"processed batch","processed_to":20559323,"behind":99000}
```

### 6. Run the Consumer

```bash
# Run consumer (foreground)
./bin/consumer

# Or with Make
make run-consumer

# Check metrics
curl http://localhost:9091/metrics
```

Expected log output:
```
{"level":"info","time":"...","message":"starting polymarket consumer"}
{"level":"info","message":"connected to database","host":"localhost","database":"polymarket"}
{"level":"info","message":"connected to nats","url":"nats://localhost:4222"}
{"level":"info","message":"created consumer","stream":"POLYMARKET","consumer":"polymarket-consumer"}
{"level":"info","message":"consumer started, waiting for messages"}
{"level":"debug","message":"processing event","event":"OrderFilled","block":20560000}
```

## Monitoring

### Health Checks

```bash
# Indexer health
curl http://localhost:8080/health

# Expected output:
# healthy
# current: 20560000
# latest: 20560100
# behind: 100
```

### Metrics (Prometheus)

Access Prometheus at http://localhost:9090

Key metrics:
- `polymarket_blocks_processed_total` - Total blocks processed
- `polymarket_events_processed_total` - Events by type
- `polymarket_syncer_block_height` - Current indexer position
- `polymarket_chain_block_height` - Latest chain block
- `polymarket_blocks_behind` - How far behind the indexer is
- `polymarket_events_consumed_total` - Consumer metrics
- `polymarket_consumer_lag_seconds` - Time lag from event to DB

Example queries:
```promql
# Events processed per second
rate(polymarket_events_processed_total[1m])

# Blocks behind chain
polymarket_blocks_behind

# Consumer lag
polymarket_consumer_lag_seconds
```

### Grafana Dashboards

Access Grafana at http://localhost:3000 (admin/admin)

1. Add Prometheus data source: http://prometheus:9090
2. Import dashboards from `grafana/dashboards/` (if available)
3. Create custom dashboards for:
   - Block processing rate
   - Event distribution by type
   - Consumer throughput
   - System lag

## Querying Data

### Connect to TimescaleDB

```bash
# Using docker
docker-compose exec timescaledb psql -U polymarket -d polymarket

# Or directly
psql -h localhost -p 5432 -U polymarket -d polymarket
```

### Example Queries

```sql
-- Total events indexed
SELECT COUNT(*) FROM events;

-- Events by type
SELECT event_signature, COUNT(*) 
FROM events 
GROUP BY event_signature 
ORDER BY COUNT(*) DESC;

-- Recent order fills
SELECT 
  block_timestamp,
  maker,
  taker,
  maker_amount_filled,
  taker_amount_filled
FROM order_fills
ORDER BY block_timestamp DESC
LIMIT 10;

-- Trading volume by hour (uses continuous aggregate)
SELECT * FROM order_volume_hourly
ORDER BY bucket DESC
LIMIT 24;

-- Active traders today
SELECT COUNT(DISTINCT trader) 
FROM daily_active_traders
WHERE bucket = date_trunc('day', NOW());

-- Token transfers for a specific address
SELECT 
  block_timestamp,
  from_address,
  to_address,
  token_id,
  amount
FROM token_transfers
WHERE from_address = '0x...' OR to_address = '0x...'
ORDER BY block_timestamp DESC;

-- Market conditions
SELECT 
  condition_id,
  oracle,
  outcome_slot_count,
  resolved,
  block_timestamp
FROM conditions
ORDER BY block_timestamp DESC;
```

## Troubleshooting

### Indexer not syncing

```bash
# Check RPC connection
curl -X POST https://your-rpc-url \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Check logs
docker-compose logs indexer | tail -100

# Verify checkpoint
sqlite3 data/checkpoints.db "SELECT * FROM checkpoints;"
```

### Consumer not processing events

```bash
# Check NATS stream
docker-compose exec nats nats stream info POLYMARKET

# Check consumer lag
docker-compose exec nats nats consumer info POLYMARKET polymarket-consumer

# Check database connection
docker-compose exec consumer psql -h timescaledb -U polymarket -d polymarket -c "SELECT 1;"
```

### High memory usage

```bash
# Reduce batch size in config.toml
[indexer]
batch_size = 50  # Default: 100

# Reduce workers
workers = 2  # Default: 4
```

### Database performance issues

```sql
-- Check table sizes
SELECT 
  schemaname,
  tablename,
  pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Check index usage
SELECT 
  schemaname,
  tablename,
  indexname,
  idx_scan,
  pg_size_pretty(pg_relation_size(indexrelid)) AS size
FROM pg_stat_user_indexes
ORDER BY idx_scan;

-- Reindex if needed
REINDEX TABLE events;
```

### Reorg detected

The indexer automatically handles reorgs by checking block hashes. If a reorg is detected:

```bash
# Check logs for reorg messages
docker-compose logs indexer | grep -i reorg

# The indexer will automatically:
# 1. Stop processing
# 2. Roll back to last known good block
# 3. Resume from safe checkpoint
```

## Production Considerations

### 1. RPC Provider
- Use a dedicated RPC provider (Alchemy, Infura, QuickNode)
- Enable WebSocket for realtime mode
- Set up rate limit handling

### 2. Database
- Use managed TimescaleDB (Timescale Cloud)
- Enable compression policies for old data
- Set up automated backups
- Configure retention policies

### 3. Monitoring
- Set up alerting in Prometheus
- Monitor disk usage
- Track consumer lag
- Alert on indexer health failures

### 4. Scaling
- Increase workers for faster backfill
- Run multiple consumers for parallel DB writes
- Use read replicas for queries
- Partition by time if data grows large

### 5. Security
- Change default passwords
- Use secrets management (Vault, AWS Secrets Manager)
- Enable SSL for PostgreSQL
- Restrict network access

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/handler -v
```

### Local Development

```bash
# Start only infrastructure
docker-compose up -d nats timescaledb

# Run indexer locally with hot reload (using air)
air -c .air.toml

# Or manually
go run ./cmd/indexer
```

### Adding New Event Handlers

1. Add event signature to `internal/handler/events.go`
2. Implement handler function
3. Register handler in `internal/processor/processor.go`
4. Add database table/columns in new migration
5. Add consumer logic in `cmd/consumer/main.go`

Example:
```go
// 1. Add signature
var NewEventSig = common.HexToHash("0x...")

// 2. Implement handler
func HandleNewEvent(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
    // Parse event
    return models.NewEvent{...}, nil
}

// 3. Register in processor.New()
r.RegisterLogHandler(handler.NewEventSig, handler.HandleNewEvent)

// 4. Create migration (migrations/002_add_new_event.up.sql)
CREATE TABLE new_events (...);

// 5. Add consumer logic
case "NewEvent":
    return storeNewEvent(ctx, pool, event)
```

## Support

For issues, questions, or contributions:
- Check logs: `docker-compose logs`
- Review configuration: `config.toml`
- Check metrics: http://localhost:9090/metrics
- Query database: `psql -h localhost -U polymarket polymarket`

## Next Steps

- Set up alerting in Grafana
- Create custom dashboards
- Build API layer on top of TimescaleDB
- Add GraphQL interface
- Implement webhooks for real-time notifications
