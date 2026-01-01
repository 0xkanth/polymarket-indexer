# Polymarket Indexer 🔮

A production-grade blockchain indexer for Polymarket prediction markets on Polygon. Built with Go, following battle-tested patterns from eth-tracker and evm-scanner.

## 📚 Documentation Hub - START HERE!

> **New to the project?** Read the [**Master Learning Path**](docs/MASTER_LEARNING_PATH.md) - a complete guide for naive developers to understand, debug, and maintain this codebase.

### 🎯 Quick Links by Role

**Backend Developer (new to blockchain)**
1. [README](#) (you are here) → Setup & run locally
2. [ARCHITECTURE.md](docs/ARCHITECTURE.md) → System design
3. [DATABASE.md](docs/DATABASE.md) → Database & checkpoints (NEW!)
4. [COMPONENT_CONNECTIONS.md](docs/COMPONENT_CONNECTIONS.md) → Component interactions

**Blockchain Developer (new to indexers)**  
1. [README](#) → Quick start
2. [ARCHITECTURE.md](docs/ARCHITECTURE.md) → High-level design
3. [SYNCER_ARCHITECTURE.md](docs/SYNCER_ARCHITECTURE.md) → Sync strategies (NEW state diagrams!)
4. [COMPONENT_CONNECTIONS.md](docs/COMPONENT_CONNECTIONS.md) → Data flow (NEW mermaid diagrams!)

**DevOps/SRE (operations focus)**  
1. [README](#) → Setup
2. [RUNNING.md](docs/RUNNING.md) → Deployment & monitoring
3. [DATABASE.md](docs/DATABASE.md) → Backup, recovery, tuning
4. [MASTER_LEARNING_PATH.md](docs/MASTER_LEARNING_PATH.md) → On-call runbook

**On-Call Engineer (emergency prep)**  
1. [README](#) → 15 min setup
2. [MASTER_LEARNING_PATH.md](docs/MASTER_LEARNING_PATH.md) → Skip to "On-Call Runbook" section
3. [RUNNING.md](docs/RUNNING.md) → Monitoring dashboard
4. [DATABASE.md](docs/DATABASE.md) → Checkpoint recovery

### 📖 All Documentation

| Document | Read Time | Purpose | Diagrams |
|----------|-----------|---------|----------|
| **[MASTER_LEARNING_PATH.md](docs/MASTER_LEARNING_PATH.md)** | 30 min | **Complete learning guide** for naive developers. Logical doc order, mental models, on-call runbook | ✅ Mermaid |
| [ARCHITECTURE.md](docs/ARCHITECTURE.md) | 15 min | System design, why each component exists | ASCII |
| [COMPONENT_CONNECTIONS.md](docs/COMPONENT_CONNECTIONS.md) | 20 min | Component interactions, data flow | ✅ **NEW Mermaid** (graph, sequence) |
| [DATABASE.md](docs/DATABASE.md) | 25 min | **NEW**: Schema, why TimescaleDB, checkpoint recovery | ✅ **NEW Mermaid** (ER, sequence, flow) |
| [SYNCER_ARCHITECTURE.md](docs/SYNCER_ARCHITECTURE.md) | 20 min | Block sync (backfill/realtime), worker pools | ✅ **NEW Mermaid** (state, sequence) + ASCII |
| [NATS_EXPLAINED.md](docs/NATS_EXPLAINED.md) | 15 min | Message bus, producer-consumer decoupling | ASCII |
| [RUNNING.md](docs/RUNNING.md) | 15 min | Production deployment, monitoring, scaling | - |
| [TESTING.md](docs/TESTING.md) | 20 min | Unit tests, fork tests, integration tests | - |
| [DEPENDENCIES.md](docs/DEPENDENCIES.md) | 10 min | External libraries (go-ethereum, NATS, etc.) | - |

---

## Features

- ⚡ **High Performance** - Process 10k+ blocks/min with <50MB RAM
- 🔄 **Real-time + Historical** - WebSocket subscription + backfill
- 📨 **NATS JetStream** - Deduplication & persistent messaging
- 📊 **TimescaleDB** - Time-series data with hypertables
- 🛡️ **Production Ready** - Reorg handling, multi-RPC failover, graceful shutdown
- 🎯 **Event-Driven** - CTF Exchange orders, token transfers, market events

## Architecture

```
┌─────────────────────────────────────────┐
│  Polygon RPC (Multi-endpoint failover)  │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Syncer (Real-time + Backfill)          │
│  - WebSocket subscription               │
│  - Checkpoint management                │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Processor (Worker Pool)                │
│  - Decode events via ABI                │
│  - Filter by contract addresses         │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Event Router & Handlers                │
│  - OrderFilled, OrderCancelled          │
│  - TransferSingle, TransferBatch        │
│  - MarketCreated                        │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  NATS JetStream (Deduplication)         │
│  Subject: POLYMARKET.{Event}.{Contract} │
└─────────────────────────────────────────┘
              ↓
┌─────────────────────────────────────────┐
│  Consumer → TimescaleDB                 │
│  - Event storage (hypertable)           │
│  - Aggregations (volume, trades)        │
└─────────────────────────────────────────┘
```

## Contracts (Polygon Mainnet)

| Contract | Address | Deployment Block |
|----------|---------|------------------|
| CTF Exchange | `0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E` | 20,558,323 |
| Conditional Tokens | `0x4D97DCd97eC945f40cF65F87097ACe5EA0476045` | 7,534,294 |
| USDC (Collateral) | `0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174` | - |

## Quick Start

### Prerequisites

- **Go 1.24+** (latest stable - for go-ethereum v1.16.7 compatibility)
- **Docker & Docker Compose** (for infrastructure)
- **Polygon RPC endpoint** (Alchemy, QuickNode, or Infura - see setup below)

### 1. Clone and Install Dependencies

```bash
git clone <your-repo>
cd polymarket-indexer

# Install Go 1.24 (if needed)
# Using gvm:
gvm install go1.24.11
gvm use go1.24.11 --default

# Verify installation
go version  # Should show go1.24.11
```

### 2. Get Production RPC Access

⚠️ **CRITICAL:** Public RPCs (`https://polygon-rpc.com`) are rate-limited and unreliable.

**Get a FREE production RPC from:**

| Provider | Free Tier | Sign Up |
|----------|-----------|---------|
| **Alchemy** | 300M compute units/month | https://www.alchemy.com/ |
| **QuickNode** | 50M credits/month | https://www.quicknode.com/ |
| **Infura** | 100k requests/day | https://infura.io/ |

**After signing up:**
1. Create a Polygon Mainnet endpoint
2. Copy both HTTP and WebSocket URLs
3. Update `config/chains.json` (see step 3)

### 3. Configure Chain Settings

Edit `config/chains.json` to add your RPC endpoint:

```bash
# Open config file
vim config/chains.json

# Update the "polygon" section:
{
  "chains": {
    "polygon": {
      "rpcUrls": [
        "https://polygon-mainnet.g.alchemy.com/v2/YOUR-API-KEY"  # ← Add your key
      ],
      "wsUrls": [
        "wss://polygon-mainnet.g.alchemy.com/v2/YOUR-API-KEY"   # ← Add your key
      ],
      "chainId": 137,
      "startBlock": 20558323,
      "confirmations": 100,
      "contracts": {
        "ctfExchange": "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E",
        "conditionalTokens": "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
      }
    }
  }
}
```

**Chain Configuration Explained:**
- `rpcUrls` - HTTP endpoints for batch requests (eth_getLogs)
- `wsUrls` - WebSocket endpoints for real-time block subscriptions
- `startBlock: 20558323` - CTF Exchange deployment block (Sept 2021)
- `confirmations: 100` - Reorg protection (Polygon has 50-100 block reorgs)
- `contracts` - Polymarket contracts to monitor

**Switch chains easily:**
```bash
# config.toml:
[chain]
name = "polygon"         # Production
# name = "polygon-fork"  # Local Anvil fork for testing
# name = "mumbai"        # Testnet
```

### 4. Start Infrastructure

```bash
# Start NATS JetStream + TimescaleDB
make infra-up

# Verify containers are running
docker ps
# Should see: polymarket-nats, polymarket-timescaledb
```

**Troubleshooting Infrastructure:**

<details>
<summary>Port 5432 already in use</summary>

```bash
# Find what's using the port
lsof -i :5432

# Stop conflicting PostgreSQL/TimescaleDB
docker stop <container-id>
# or
brew services stop postgresql

# Retry
make infra-up
```
</details>

<details>
<summary>NATS container restarting</summary>

```bash
# Check NATS logs
docker logs polymarket-nats

# Restart infrastructure
make infra-down
make infra-up
```
</details>

**Health Check:**
```bash
# NATS monitoring endpoint
curl http://localhost:8222/healthz

# TimescaleDB connection
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c "SELECT 1;"
```

### 5. Initialize Database Schema

```bash
# Apply database migrations (creates tables, indexes, hypertables)
make migrate-up
```

**What this does:**
- ✅ Enables TimescaleDB extension
- ✅ Creates `events` hypertable (partitioned by time)
- ✅ Creates specialized tables: `order_fills`, `token_transfers`, `conditions`
- ✅ Creates indexes for fast queries
- ✅ Sets up continuous aggregates for analytics

**Verify migration:**
```bash
# List created tables
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c "\dt"

# Expected output:
# checkpoints, events, order_fills, token_transfers, conditions, token_registrations
```

⚠️ **Note:** `make migrate-up` is not idempotent. Run only once after `make infra-up`.

### 6. Build Binaries

```bash
# Build indexer and consumer
make build

# Verify binaries exist
ls -lh bin/
# Should see: indexer, consumer
```

### 7. Run Indexer (Terminal 1)

```bash
# Start indexing Polygon blockchain
make run-indexer

# Or run binary directly:
./bin/indexer
```

**Expected Output:**
```json
{"level":"info","message":"starting polymarket indexer"}
{"level":"info","chain":"Polygon Mainnet","chain_id":137,"rpc":"https://polygon-mainnet.g.alchemy.com/...","message":"loaded chain configuration"}
{"level":"info","contracts":["0x4bFb...","0x4D97..."],"start_block":20558323,"message":"initialized processor"}
{"level":"info","message":"syncer started"}
{"level":"info","block":20558323,"message":"processing block"}
{"level":"info","block":20558324,"events":3,"message":"published events to NATS"}
```

**Monitor Progress:**
```bash
# Check NATS stream
curl http://localhost:8222/jsz | jq '.streams[0]'

# Check Prometheus metrics
curl http://localhost:9090/metrics | grep polymarket_blocks_processed_total
```

### 8. Run Consumer (Terminal 2)

```bash
# Start consuming from NATS → writing to TimescaleDB
make run-consumer

# Or run binary directly:
./bin/consumer
```

**Expected Output:**
```json
{"level":"info","message":"starting polymarket consumer"}
{"level":"info","host":"localhost","database":"polymarket","message":"connected to database"}
{"level":"info","url":"nats://localhost:4222","message":"connected to nats"}
{"level":"info","stream":"POLYMARKET_EVENTS","message":"consumer started"}
{"level":"info","event":"OrderFilled","block":20558500,"message":"processed event"}
```

**Verify Data:**
```bash
# Check events in database
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c \
  "SELECT COUNT(*) FROM order_fills;"

# View recent trades
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c \
  "SELECT * FROM order_fills ORDER BY time DESC LIMIT 10;"
```

---

## Configuration Files

### `config.toml` - Runtime Behavior

Controls indexer/consumer runtime settings:

```toml
[chain]
name = "polygon"  # Selects chain from config/chains.json

[db]
checkpoint_path = "data/checkpoints.db"  # BoltDB for resume tracking

[nats]
url = "nats://localhost:4222"
stream_name = "POLYMARKET_EVENTS"
max_age = "168h"  # Keep messages for 7 days

[indexer]
batch_size = 100        # Blocks per batch
poll_interval = "2s"    # Block polling frequency
workers = 5             # Concurrent processing workers

[postgres]
host = "localhost"
port = 5432
user = "polymarket"
password = "polymarket"
database = "polymarket"

[logging]
level = "info"  # debug, info, warn, error
```

### `config/chains.json` - Chain-Specific Data

Multi-chain configuration (network, contracts, RPC endpoints):

```json
{
  "chains": {
    "polygon": {
      "chainId": 137,
      "name": "Polygon Mainnet",
      "rpcUrls": ["https://your-rpc-url"],
      "wsUrls": ["wss://your-ws-url"],
      "startBlock": 20558323,
      "confirmations": 100,
      "contracts": {
        "ctfExchange": "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E",
        "conditionalTokens": "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
      }
    }
  }
}
```

**Benefits of chains.json:**
- ✅ Switch networks by changing `chain.name` in config.toml
- ✅ Support multiple chains (polygon, polygon-fork, mumbai)
- ✅ Centralized contract addresses (no duplication)
- ✅ RPC failover support (multiple URLs per chain)

## Project Structure

```
polymarket-indexer/
├── cmd/
│   ├── indexer/          # Main indexer service
│   └── consumer/         # NATS → TimescaleDB consumer
├── internal/
│   ├── chain/            # Multi-RPC client
│   ├── syncer/           # Block syncing (real-time + backfill)
│   ├── processor/        # Event processing
│   ├── cache/            # Address filtering
│   ├── router/           # Event routing
│   ├── handler/          # Event handlers
│   ├── nats/             # NATS publisher/consumer
│   └── db/               # Database (TimescaleDB + BoltDB)
├── pkg/
│   ├── models/           # Domain models
│   └── contracts/        # ABIs + Go bindings
├── migrations/           # SQL migrations
├── docker-compose.yml
└── Makefile
```

## Events Indexed

### CTF Exchange (`0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E`)

- `OrderFilled` - Trades executed
- `OrderCancelled` - Orders cancelled
- `OrdersMatched` - Order matching
- `TokenRegistered` - New market tokens

### Conditional Tokens (`0x4D97DCd97eC945f40cF65F87097ACe5EA0476045`)

- `TransferSingle` - Single token transfer
- `TransferBatch` - Batch token transfer
- `ConditionPreparation` - New condition/market
- `ConditionResolution` - Market resolution
- `PositionSplit` - Position minting
- `PositionsMerge` - Position redemption

## Development

### Generate Contract Bindings

```bash
make generate-bindings
```

### Run Tests

```bash
make test
```

### Lint

```bash
make lint
```

## Monitoring

### Infrastructure Health

**NATS JetStream:**
```bash
# Health check
curl http://localhost:8222/healthz

# Stream statistics
curl http://localhost:8222/jsz | jq

# View messages in stream
curl http://localhost:8222/jsz | jq '.streams[] | {name, messages, bytes}'
```

**TimescaleDB:**
```bash
# Connection test
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c "SELECT 1;"

# List hypertables
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c \
  "SELECT * FROM timescaledb_information.hypertables;"

# Table sizes
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c \
  "SELECT 
    schemaname, 
    tablename, 
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
  FROM pg_tables 
  WHERE schemaname = 'public' 
  ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```

### Application Health

**Indexer Health Check:**
```bash
curl http://localhost:9090/health
```

**Indexer Metrics (Prometheus):**
```bash
curl http://localhost:9090/metrics

# Key metrics:
# polymarket_blocks_processed_total - Total blocks indexed
# polymarket_events_processed_total{event_type="OrderFilled"} - Events by type
# polymarket_block_processing_duration_seconds - Processing latency
# polymarket_processing_errors_total - Error counter
```

**Consumer Metrics:**
```bash
curl http://localhost:9091/metrics

# Key metrics:
# polymarket_events_consumed_total{event_type="OrderFilled"} - Events consumed
# polymarket_consume_errors_total - Consumer errors
# polymarket_processing_lag_seconds - Time lag between event and processing
```

**Database Statistics:**
```bash
# Events indexed
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c \
  "SELECT 
    event_name, 
    COUNT(*) as count,
    MIN(time) as first_event,
    MAX(time) as last_event
  FROM events 
  GROUP BY event_name 
  ORDER BY count DESC;"

# Recent trades
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c \
  "SELECT * FROM order_fills ORDER BY time DESC LIMIT 10;"

# Indexer progress
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c \
  "SELECT * FROM checkpoints;"
```

### Troubleshooting

<details>
<summary><b>Indexer not processing blocks</b></summary>

```bash
# Check RPC connection
curl -X POST https://your-rpc-url \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Check logs
docker logs polymarket-indexer --tail 100

# Verify chain config loaded
grep "loaded chain configuration" logs

# Check NATS connection
docker logs polymarket-nats | grep -i error
```
</summary>

<details>
<summary><b>Consumer not writing to database</b></summary>

```bash
# Check database connection
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c "SELECT 1;"

# Check NATS stream has messages
curl http://localhost:8222/jsz | jq '.streams[0].state.messages'

# Check consumer logs
docker logs polymarket-consumer --tail 100

# Verify tables exist
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c "\dt"
```
</details>

<details>
<summary><b>RPC rate limit errors</b></summary>

```bash
# Update config/chains.json with production RPC
# Alchemy: 300M compute units/month
# QuickNode: 50M credits/month
# Infura: 100k requests/day

# Reduce batch size in config.toml
[indexer]
batch_size = 50  # Reduce from 100
poll_interval = "5s"  # Slow down polling

# Add multiple RPC URLs for failover
"rpcUrls": [
  "https://primary-rpc-url",
  "https://backup-rpc-url"
]
```
</details>

<details>
<summary><b>Database migration errors</b></summary>

```bash
# Reset database (WARNING: deletes all data)
make infra-down
docker volume rm polymarket-indexer_timescale_data
make infra-up
make migrate-up

# Check for existing tables
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket -c "\dt"

# Re-run specific parts of migration manually
docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket < migrations/001_initial_schema.up.sql
```
</details>

## Production Deployment

See `docs/DEPLOYMENT.md` for production deployment guide including:
- Multi-RPC setup
- Kubernetes manifests
- Monitoring & alerting
- Backup strategies

## Contributing

Contributions welcome! Please see `CONTRIBUTING.md`.

## License

MIT

## Acknowledgments

Built on patterns from:
- [eth-tracker](https://github.com/grassrootseconomics/eth-tracker) - NATS integration
- [evm-scanner](https://github.com/84hero/evm-scanner) - RPC resilience

---

**Built with ❤️ for prediction markets**
