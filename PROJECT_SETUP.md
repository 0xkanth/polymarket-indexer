# Polymarket Indexer - Project Overview

## ✅ What's Been Created

A complete production-ready Polymarket indexer project structure with all necessary configurations, ABIs, and infrastructure setup.

### 📁 Project Structure

```
polymarket-indexer/
├── README.md                          # Comprehensive project documentation
├── config.toml                        # Production configuration (Polygon mainnet)
├── go.mod                             # Go dependencies
├── Makefile                           # Build, test, and deployment commands
├── Dockerfile                         # Multi-stage Docker build
├── docker-compose.yml                 # Full stack (NATS, TimescaleDB, Indexer, Consumer)
├── .gitignore                         # Git ignore rules
│
├── pkg/contracts/abi/                 # Contract ABIs
│   ├── CTFExchange.json              # Polymarket CTF Exchange ABI
│   ├── ConditionalTokens.json        # Conditional Tokens Framework ABI
│   └── ERC20.json                    # ERC-20 token ABI (USDC)
│
└── migrations/
    └── 001_initial_schema.up.sql     # TimescaleDB schema with hypertables
```

## 🎯 Key Information Configured

### Polygon Mainnet Contracts

| Contract | Address | Block | Purpose |
|----------|---------|-------|---------|
| **CTF Exchange** | `0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E` | 20,558,323 | Order matching & trades |
| **Conditional Tokens** | `0x4D97DCd97eC945f40cF65F87097ACe5EA0476045` | 7,534,294 | ERC-1155 position tokens |
| **USDC** | `0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174` | - | Collateral token |

### Events Indexed

**CTF Exchange:**
- `OrderFilled` - Trade executions
- `OrderCancelled` - Order cancellations
- `OrdersMatched` - Order matching
- `TokenRegistered` - New market creation
- `FeeCharged` - Fee events

**Conditional Tokens:**
- `TransferSingle` / `TransferBatch` - Token transfers
- `ConditionPreparation` - New markets/conditions
- `ConditionResolution` - Market resolution
- `PositionSplit` - Position minting
- `PositionsMerge` - Position redemption

### Configuration Highlights

```toml
[chain]
rpc_endpoint = "https://polygon-rpc.com"  # ⚠️ UPDATE THIS
ws_endpoint = "wss://polygon-rpc.com"     # ⚠️ UPDATE THIS
chainid = 137
start_block = 20558323                     # CTF Exchange deployment

[contracts]
ctf_exchange = "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E"
conditional_tokens = "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
usdc = "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174"

[jetstream]
endpoint = "nats://localhost:4222"
persist_duration_hrs = 168  # 7 days
```

## 📊 Database Schema (TimescaleDB)

### Hypertables (time-series optimized)

1. **events** - All blockchain events
2. **order_fills** - Trade data
3. **token_transfers** - Token movements
4. **token_registrations** - Market registrations

### Materialized Views (auto-updating)

1. **trading_volume_hourly** - Hourly volume per market
2. **active_traders_daily** - Daily unique traders
3. **market_activity_daily** - Market-level activity

### Helper Functions

- `get_recent_fills(asset_id, hours)` - Recent trades for a market
- `get_user_stats(address, days)` - User trading statistics

## 🚀 Next Steps

### 1. Update Configuration

```bash
cd polymarket-indexer
nano config.toml
```

**Required changes:**
- Add your Polygon RPC endpoint (get free tier from Alchemy/QuickNode/Infura)
- Update `rpc_endpoint` and `ws_endpoint`

### 2. Set Up Environment

```bash
# Install Go dependencies
make deps

# Check your environment
make check-env

# Start infrastructure (NATS + TimescaleDB)
make infra-up
```

### 3. Generate Contract Bindings

```bash
# Install abigen if not already installed
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# Generate Go bindings from ABIs
make generate-bindings
```

This will create:
- `pkg/contracts/bindings/ctf_exchange.go`
- `pkg/contracts/bindings/conditional_tokens.go`
- `pkg/contracts/bindings/erc20.go`

### 4. Implement Core Code (TODO)

You now need to implement the Go code following patterns from eth-tracker and evm-scanner:

```
internal/
├── chain/              # RPC client with failover
├── syncer/             # Block syncing (WebSocket + HTTP)
├── processor/          # Event processor
├── cache/              # Address cache (optional)
├── router/             # Event router
├── handler/            # Event handlers
│   ├── order_filled.go
│   ├── token_transfer.go
│   └── condition_preparation.go
├── nats/               # NATS publisher/consumer
└── db/                 # Database layer

cmd/
├── indexer/
│   └── main.go        # Main indexer service
└── consumer/
    └── main.go        # Consumer service
```

### 5. Build and Run

```bash
# Build binaries
make build

# Run indexer
./bin/indexer -config config.toml

# In another terminal, run consumer
./bin/consumer -config config.toml
```

Or with Docker:

```bash
# Build and start everything
make docker-up

# View logs
make docker-logs
```

## 🔧 Development Commands

```bash
# Build
make build                 # Build all binaries
make generate-bindings     # Generate contract bindings

# Run
make run-indexer          # Run indexer locally
make run-consumer         # Run consumer locally

# Test
make test                 # Run tests
make lint                 # Run linter

# Docker
make docker-build         # Build Docker images
make docker-up            # Start all services
make docker-logs          # View logs

# Infrastructure
make infra-up             # Start NATS + TimescaleDB
make infra-down           # Stop infrastructure
make infra-reset          # Reset everything

# Database
make migrate-up           # Run migrations
make stats                # Show indexer stats
make health               # Health check
```

## 📖 Implementation References

Follow patterns from:

1. **eth-tracker** (`reference/eth-tracker/`)
   - NATS JetStream integration
   - Event router pattern
   - Checkpoint management
   - Worker pool pattern

2. **evm-scanner** (`reference/evm-scanner/`)
   - Multi-RPC client with failover
   - ABI decoder
   - Sink pattern for outputs
   - Configuration patterns

## 🎓 Learning Resources

- **LEARNING_PATH.md** - Complete learning roadmap
- **eth-tracker README** - NATS patterns
- **evm-scanner README** - RPC resilience patterns
- [Polymarket Docs](https://docs.polymarket.com/)
- [go-ethereum Docs](https://geth.ethereum.org/docs/developers)
- [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
- [TimescaleDB](https://docs.timescale.com/)

## ⚠️ Important Notes

### Polygon RPC Considerations

1. **Free Tier Limits**: Public RPCs have rate limits. Consider:
   - [Alchemy](https://www.alchemy.com/) - 300M CU/month free
   - [QuickNode](https://www.quicknode.com/) - Free tier available
   - [Infura](https://infura.io/) - 100k requests/day free

2. **Multi-RPC Setup**: Implement failover with 2-3 endpoints for production

3. **Archive Node**: Required for backfilling from block 20,558,323

### Reorg Protection

- Polygon experiences reorgs (~50-100 blocks)
- Config sets `confirmation_depth = 100`
- Always store block hash with checkpoint

### Rate Limiting

- Start slow: `rate_limit = 100` blocks/sec
- Monitor RPC provider limits
- Scale up as needed

## 📝 What You Got

✅ **Complete Project Structure**
✅ **Production Configuration** (needs RPC endpoint)
✅ **Contract ABIs** (ready for binding generation)
✅ **TimescaleDB Schema** (hypertables, indexes, aggregations)
✅ **Docker Setup** (NATS, TimescaleDB, Indexer, Consumer)
✅ **Makefile** (50+ commands for development)
✅ **Documentation** (README, LEARNING_PATH)

## 🎯 What's Next

❌ **Implement Go Code** (follow reference patterns)
- Chain RPC client
- Block syncer
- Event processor
- Event handlers
- NATS publisher/consumer
- Database layer

This should take 2-3 weeks following the learning path!

## 🤝 Support

- Study `reference/eth-tracker` for NATS patterns
- Study `reference/evm-scanner` for RPC patterns
- Refer to `LEARNING_PATH.md` for structured guidance

---

**Built with production-grade patterns from battle-tested projects!** 🚀

Ready to index Polymarket? Update your RPC endpoint and start coding! 💪
