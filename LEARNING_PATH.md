# Production Go EVM Indexer - Complete Learning Path

## 🎯 Your Goal
Build production-grade blockchain indexer expertise for:
- EVM event indexing on Polygon mainnet
- Complex ABI handling for Polymarket contracts
- NATS JetStream integration with deduplication
- TimescaleDB for time-series analytics
- Production-grade resilience (reorg handling, multi-RPC failover, retry logic)
- Polymarket prediction market indexing (CTF Exchange + Conditional Tokens)

## 📋 What This Indexer Does

This is a **real-time blockchain event indexer** for [Polymarket](https://polymarket.com), a prediction markets platform on Polygon. It:

1. **Extracts** - Fetches blocks and event logs from Polygon RPC
2. **Decodes** - Parses ABI-encoded events (trades, token transfers, market creations)
3. **Routes** - Publishes events to NATS JetStream with deduplication
4. **Stores** - Consumer service writes to TimescaleDB for analytics
5. **Scales** - Handles 10k+ blocks/min with <50MB RAM

**Architecture Pattern:** Producer (Indexer) → Message Broker (NATS) → Consumer (DB Writer)

## 🏗️ Current Project Architecture

This indexer already has a **solid foundation** implemented. Here's what exists:

### **Implemented Components** ✅

```
polymarket-indexer/
├── cmd/
│   ├── indexer/main.go       ✅ Entry point with graceful shutdown
│   └── consumer/main.go      ✅ NATS consumer with TimescaleDB writer
├── internal/
│   ├── chain/client.go       ✅ Dual RPC (HTTP + WebSocket)
│   ├── syncer/syncer.go      ✅ Backfill + Realtime modes
│   ├── processor/processor.go ✅ Block fetching, log filtering, routing
│   ├── handler/events.go     ✅ 9 event handlers (OrderFilled, etc.)
│   ├── router/router.go      ✅ Event signature → handler mapping
│   ├── nats/publisher.go     ✅ JetStream with deduplication
│   └── db/checkpoint.go      ✅ BoltDB checkpoint persistence
├── pkg/
│   ├── contracts/            ✅ Generated Go bindings (CTFExchange, ConditionalTokens)
│   ├── models/event.go       ✅ Event data models
│   ├── service/ctf_service.go ✅ High-level CTF service with transaction helpers
│   ├── config/config.go      ✅ Chain configuration loader
│   └── txhelper/transaction.go ✅ Production transaction helpers (simulation, retry)
├── test/fork_test.go         ✅ Anvil fork testing examples
├── migrations/               ✅ TimescaleDB schema with hypertables
├── scripts/                  ✅ Fork testing scripts (daemon mode)
└── docs/                     ✅ Comprehensive documentation
```

### **Implementation Status**

**Core Indexer:** ~80% complete
- ✅ Chain client with dual RPC
- ✅ Syncer orchestration (backfill/realtime switching)
- ✅ Processor with parallel workers
- ✅ Event handlers for all Polymarket events
- ✅ Router with handler registry
- ✅ NATS publisher with deduplication
- ✅ Checkpoint store (BoltDB)
- ⚠️ Needs: Multi-RPC failover, advanced retry logic, reorg detection

**Consumer Service:** ~70% complete
- ✅ NATS consumer with durable subscription
- ✅ TimescaleDB connection pooling
- ✅ Event storage with idempotency
- ✅ Prometheus metrics
- ⚠️ Needs: Batch optimization, error handling polish

**Infrastructure:** 100% complete
- ✅ Docker Compose (NATS, TimescaleDB, Prometheus, Grafana)
- ✅ Database migrations with hypertables
- ✅ Fork testing with Anvil
- ✅ Configuration management
- ✅ Prometheus metrics

## ✅ Reference Projects (For Learning Patterns)

### 1. **eth-tracker** by Grassroots Economics ⭐ PRIMARY REFERENCE
**Location:** `reference/eth-tracker/` (if available)

**Why it's perfect:**
- ✅ **NATS JetStream integration** (exactly what you need!)
- ✅ Production-ready (10k blocks/min, 50MB RAM)
- ✅ Distributed deployment with deduplication
- ✅ Reorg handling patterns
- ✅ Event-driven architecture

**What to study:**
- NATS deduplication patterns
- Worker pool orchestration
- Event handler design
- Checkpoint recovery
- Distributed deployment

### 2. **evm-scanner** by 84hero ⭐ SECONDARY REFERENCE
**Location:** `reference/evm-scanner/` (if available)

**Why it's valuable:**
- ✅ **Multi-RPC load balancing & failover**
- ✅ Production-ready resilience patterns
- ✅ Extensible sink architecture
- ✅ Excellent examples directory

**What to study:**
- RPC failover strategies
- Connection pooling
- Retry with exponential backoff
- Configuration patterns
- Error classification (retryable vs permanent)

## 📚 Core Concepts (Understanding This Indexer)

### 1. **Block Processing Flow**

```
┌─────────────────────────────────────────┐
│  Polygon RPC (Multi-endpoint)           │
│  ├─ BlockNumber() → latest             │
│  ├─ BlockByNumber(N) → block data      │
│  └─ FilterLogs(...) → event logs       │
└─────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────┐
│  Syncer (internal/syncer/syncer.go)     │
│  ├─ GetOrCreateCheckpoint()            │
│  ├─ Mode: Backfill or Realtime?        │
│  │  • Backfill: 100 blocks/batch, 5    │
│  │    workers parallel                  │
│  │  • Realtime: Sequential processing   │
│  └─ UpdateCheckpoint(block, hash)      │
└─────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────┐
│  Processor (internal/processor/)         │
│  ├─ FetchBlockHeader(blockNum)         │
│  ├─ FilterLogs(contracts, fromBlock)   │
│  ├─ RouteLog(log) → Handler            │
│  └─ PublishToNATS(event)               │
└─────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────┐
│  Router (internal/router/router.go)     │
│  ├─ Map: EventSig → HandlerFunc        │
│  ├─ Execute: handler(log) → payload    │
│  └─ Callback: publishToNATS(event)     │
└─────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────┐
│  Handlers (internal/handler/events.go)  │
│  ├─ HandleOrderFilled()                │
│  ├─ HandleTransferSingle()             │
│  ├─ HandleConditionPreparation()       │
│  └─ Returns: Typed event struct        │
└─────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────┐
│  NATS JetStream (Message Broker)        │
│  Subject: POLYMARKET.OrderFilled.0x4bFb│
│  MsgID: {txHash}-{logIndex}            │
│  Deduplication: 20-minute window        │
│  Retention: 7 days                      │
└─────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────┐
│  Consumer (cmd/consumer/main.go)        │
│  ├─ Subscribe: POLYMARKET.>            │
│  ├─ Process: Unmarshal event           │
│  ├─ Store: INSERT INTO events          │
│  │  • events (hypertable)              │
│  │  • order_fills (hypertable)         │
│  │  • token_transfers (hypertable)     │
│  └─ Ack: msg.Ack() or msg.Nak()        │
└─────────────────────────────────────────┘
           ↓
┌─────────────────────────────────────────┐
│  TimescaleDB (Postgres + time-series)   │
│  ├─ Hypertables (auto-partitioned)     │
│  ├─ Continuous Aggregates (hourly vol) │
│  ├─ Compression (old data)             │
│  └─ Indexes (maker, taker, timestamp)  │
└─────────────────────────────────────────┘
```

### 2. **Polymarket Smart Contracts**

#### **CTF Exchange** (`0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E`)
**Purpose:** Order matching and trade execution for prediction markets

**Events We Index:**
```solidity
// Trade execution
event OrderFilled(
    bytes32 indexed orderHash,
    address indexed maker,
    address indexed taker,
    uint256 makerAssetId,      // Position token sold
    uint256 takerAssetId,      // Position token bought
    uint256 makerAmountFilled,
    uint256 takerAmountFilled,
    uint256 fee
);

// Order cancellation
event OrderCancelled(bytes32 indexed orderHash);

// New market registration
event TokenRegistered(
    uint256 indexed token0,
    uint256 indexed token1,
    bytes32 indexed conditionId
);
```

**Handler:** `internal/handler/events.go::HandleOrderFilled()`

#### **Conditional Tokens** (`0x4D97DCd97eC945f40cF65F87097ACe5EA0476045`)
**Purpose:** ERC-1155 tokens representing market outcome positions

**Events We Index:**
```solidity
// Token transfers (ERC-1155)
event TransferSingle(
    address indexed operator,
    address indexed from,
    address indexed to,
    uint256 id,     // Position token ID
    uint256 value   // Amount
);

event TransferBatch(
    address indexed operator,
    address indexed from,
    address indexed to,
    uint256[] ids,
    uint256[] values
);

// New market creation
event ConditionPreparation(
    bytes32 indexed conditionId,
    address indexed oracle,
    bytes32 indexed questionId,
    uint256 outcomeSlotCount
);

// Market resolution
event ConditionResolution(
    bytes32 indexed conditionId,
    address indexed oracle,
    bytes32 indexed questionId,
    uint256 outcomeSlotCount,
    uint256[] payoutNumerators
);

// Minting outcome tokens from collateral
event PositionSplit(
    address indexed stakeholder,
    address collateralToken,
    bytes32 indexed parentCollectionId,
    bytes32 indexed conditionId,
    uint256[] partition,
    uint256 amount
);

// Redeeming collateral from outcome tokens
event PositionsMerge(
    address indexed stakeholder,
    address collateralToken,
    bytes32 indexed parentCollectionId,
    bytes32 indexed conditionId,
    uint256[] partition,
    uint256 amount
);
```

**Handlers:** `internal/handler/events.go` (9 handlers total)

### 3. **Reorg Handling**

**Problem:** Polygon has occasional deep reorgs (~50-100 blocks)

**Solution:**
```go
// 1. Wait for confirmations (config.toml)
confirmation_depth = 100  // ~3-4 minutes on Polygon

// 2. Store block hash with checkpoint
type Checkpoint struct {
    LastBlock     uint64
    LastBlockHash string  // ← For reorg detection
}

// 3. Detect reorg (pseudo-code)
checkpoint := db.GetCheckpoint()
block := chain.GetBlockByNumber(checkpoint.LastBlock)

if block.Hash() != checkpoint.LastBlockHash {
    // Reorg detected!
    // 1. Stop processing
    // 2. Roll back to safe block (N-100)
    // 3. Clear affected NATS messages
    // 4. Resume from safe block
}
```

**Status:** Basic checkpoint implemented, full reorg detection TODO

### 4. **NATS JetStream Integration**

**Why NATS over Kafka?**
- 15MB binary vs ~100MB+ JVM ecosystem
- 5-minute setup vs hours of configuration
- Built-in deduplication (critical for blockchain)
- Native Go client
- Perfect for startup velocity

**Stream Design:**
```go
// Stream configuration (internal/nats/publisher.go)
Stream: POLYMARKET
Subjects: POLYMARKET.*
MaxAge: 168 hours (7 days)
Duplicates: 20 minutes
Storage: FileStorage
```

**Subject Hierarchy:**
```
POLYMARKET.OrderFilled.0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E
POLYMARKET.TransferSingle.0x4D97DCd97eC945f40cF65F87097ACe5EA0476045
POLYMARKET.ConditionPreparation.0x4D97DCd97eC945f40cF65F87097ACe5EA0476045
```

**Deduplication:**
```go
// Message ID prevents duplicates on restart/reprocess
msgID := fmt.Sprintf("%s-%d", txHash, logIndex)
js.Publish(subject, data, jetstream.WithMsgID(msgID))

// NATS checks: "Already have this msgID? Skip storage."
```

**Consumer Pattern:**
```go
// Durable consumer survives restarts
js.CreateOrUpdateConsumer("POLYMARKET", jetstream.ConsumerConfig{
    Name:    "polymarket-consumer",
    Durable: "polymarket-consumer",  // Persists position
    FilterSubjects: []string{"POLYMARKET.>"},
    AckPolicy: jetstream.AckExplicitPolicy,
})

// Message handling
msg := <-subscription
processEvent(msg.Data())
msg.Ack() // ✓ or msg.Nak() ✗
```

**See:** `docs/NATS_EXPLAINED.md` for deep dive

### 5. **Database Schema (TimescaleDB)**

**Why TimescaleDB?**
- Time-series optimized Postgres
- Automatic partitioning by time (hypertables)
- Continuous aggregates (pre-computed analytics)
- Native SQL compatibility

**Hypertables:**
```sql
-- Raw event storage (migrations/001_initial_schema.up.sql)
CREATE TABLE events (
    time TIMESTAMPTZ NOT NULL,
    block_number BIGINT NOT NULL,
    tx_hash TEXT NOT NULL,
    log_index INTEGER NOT NULL,
    contract_address TEXT NOT NULL,
    event_name TEXT NOT NULL,
    event_data JSONB,
    UNIQUE(tx_hash, log_index)
);
SELECT create_hypertable('events', 'time');

-- Trade data
CREATE TABLE order_fills (
    time TIMESTAMPTZ NOT NULL,
    block_number BIGINT NOT NULL,
    tx_hash TEXT NOT NULL,
    maker TEXT NOT NULL,
    taker TEXT NOT NULL,
    maker_asset_id NUMERIC NOT NULL,
    taker_asset_id NUMERIC NOT NULL,
    maker_amount_filled NUMERIC NOT NULL,
    taker_amount_filled NUMERIC NOT NULL,
    fee NUMERIC NOT NULL
);
SELECT create_hypertable('order_fills', 'time');

-- Continuous aggregate for analytics
CREATE MATERIALIZED VIEW trading_volume_hourly
WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('1 hour', time) AS hour,
    COUNT(*) AS trades,
    SUM(maker_amount_filled::numeric) AS volume
FROM order_fills
GROUP BY hour
WITH NO DATA;

-- Auto-refresh policy
SELECT add_continuous_aggregate_policy('trading_volume_hourly',
    start_offset => INTERVAL '2 hours',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');
```

**Idempotency:**
```sql
-- Consumer uses this pattern
INSERT INTO events (...)
VALUES (...)
ON CONFLICT (tx_hash, log_index) DO NOTHING;
```

### 6. **Production-Grade Patterns**

#### **Exponential Backoff Retry**
```go
// pkg/txhelper/transaction.go
config := &TransactionConfig{
    MaxRetries:     3,
    InitialBackoff: 1 * time.Second,
    MaxBackoff:     30 * time.Second,
}

// Retries on: connection errors, timeouts, rate limits
// No retry on: execution reverted, insufficient funds, nonce too low
```

#### **Transaction Simulation**
```go
// Before sending, simulate via eth_call
err := client.CallContract(ctx, msg, nil)
if err != nil {
    return fmt.Errorf("would revert: %w", err)
}
```

#### **Gas Estimation with Buffer**
```go
gasLimit, _ := client.EstimateGas(ctx, msg)
gasLimit = gasLimit * 120 / 100  // +20% buffer
```

**See:** `docs/TRANSACTION_HELPERS.md` for complete patterns

## 🗺️ Your Learning Roadmap (Rust/TS → Golang Web3)

Since you're coming from **Rust and TypeScript**, here's how concepts map:

### **Rust → Go Translations**

| Rust Concept | Go Equivalent | Notes |
|--------------|---------------|-------|
| `Result<T, E>` | `(T, error)` return | Go uses explicit error returns |
| `Option<T>` | Pointers + nil check | `*Type` or sentinel values |
| `impl Trait` | Interfaces | `type Service interface { ... }` |
| `Arc<Mutex<T>>` | Channels / `sync.Mutex` | Go prefers channels for concurrency |
| `async/await` | Goroutines | `go func() { ... }()` |
| `tokio::spawn` | `go func()` | Goroutines are way cheaper |
| Cargo.toml | go.mod | Dependency management |
| `match` | `switch` / `if-else` | Less exhaustive checking |

### **TypeScript/Ethers.js → Go/Go-Ethereum**

| TypeScript/Ethers | Go/Go-Ethereum | Example |
|-------------------|----------------|---------|
| `ethers.Provider` | `*ethclient.Client` | RPC connection |
| `contract.on("event")` | `contract.FilterEvent()` | Event filtering |
| `await tx.wait()` | `receipt, err := WaitMined(ctx, tx)` | Transaction confirmation |
| `BigNumber.from()` | `new(big.Int).SetString()` | Large numbers |
| `utils.parseUnits()` | `new(big.Int).Mul(..., 1e18)` | Token decimals |
| `Contract.filters` | `bind.FilterOpts` | Log filtering |

### **Phase 1: Golang Fundamentals (If Needed) - Week 1**

**Skip if you're comfortable with Go basics**

Study:
1. **Goroutines and Channels**
   ```go
   // Worker pool pattern (used in syncer)
   jobs := make(chan Job, 100)
   results := make(chan Result, 100)
   
   for w := 0; w < workers; w++ {
       go func() {
           for job := range jobs {
               results <- process(job)
           }
       }()
   }
   ```

2. **Interfaces**
   ```go
   // Used throughout for dependency injection
   type ChainClient interface {
       GetBlockNumber(ctx) (uint64, error)
       FilterLogs(ctx, FilterQuery) ([]Log, error)
   }
   ```

3. **Context for Cancellation**
   ```go
   // Everything takes context for graceful shutdown
   func (s *Syncer) Start(ctx context.Context) error {
       for {
           select {
           case <-ctx.Done():
               return ctx.Err()
           default:
               processBlock()
           }
       }
   }
   ```

**Resources:**
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)

### **Phase 2: Study The Implemented Code - Week 1-2**

**Start here!** Understand what's already built:

#### **Day 1-2: Core Data Flow**

Read in this order:

1. **Entry Point:** [cmd/indexer/main.go](cmd/indexer/main.go)
   - See how components wire together
   - Graceful shutdown pattern
   - Signal handling

2. **Configuration:** [config.toml](config.toml) + [pkg/config/config.go](pkg/config/config.go)
   - RPC endpoints, chain ID
   - Contract addresses
   - NATS/TimescaleDB DSNs

3. **Models:** [pkg/models/event.go](pkg/models/event.go)
   - Event envelope structure
   - Typed event payloads

4. **Chain Client:** [internal/chain/client.go](internal/chain/client.go)
   - Dual RPC (HTTP + WebSocket)
   - Chain ID verification
   - Block/log fetching

#### **Day 3-4: Event Processing**

5. **Event Handlers:** [internal/handler/events.go](internal/handler/events.go)
   ```go
   // Study how ABI decoding works
   func HandleOrderFilled(ctx, log, timestamp) (any, error) {
       // 1. Extract indexed params from topics
       orderHash := log.Topics[1].Hex()
       maker := common.BytesToAddress(log.Topics[2].Bytes())
       
       // 2. Extract non-indexed params from data
       makerAssetID := new(big.Int).SetBytes(log.Data[0:32])
       
       // 3. Return typed struct
       return models.OrderFilled{ ... }, nil
   }
   ```

6. **Router:** [internal/router/router.go](internal/router/router.go)
   - Event signature → handler mapping
   - Callback pattern for NATS publish

7. **Processor:** [internal/processor/processor.go](internal/processor/processor.go)
   - Block fetching orchestration
   - Log filtering by contracts
   - Parallel processing

#### **Day 5-7: Infrastructure**

8. **Syncer:** [internal/syncer/syncer.go](internal/syncer/syncer.go)
   - Backfill vs realtime modes
   - Worker pool management
   - Checkpoint updates

9. **NATS Publisher:** [internal/nats/publisher.go](internal/nats/publisher.go)
   - JetStream stream creation
   - Deduplication via message ID
   - Subject hierarchy

10. **Checkpoint:** [internal/db/checkpoint.go](internal/db/checkpoint.go)
    - BoltDB embedded storage
    - Checkpoint persistence
    - Crash recovery

11. **Consumer:** [cmd/consumer/main.go](cmd/consumer/main.go)
    - NATS subscription
    - TimescaleDB writes
    - Idempotent inserts

#### **Practice: Run the Code**

```bash
# 1. Update config
nano config.toml
# Add your Polygon RPC endpoint

# 2. Start dependencies
docker-compose up -d nats timescaledb

# 3. Run migrations
make migrate-up

# 4. Build and run
make build
./bin/indexer

# 5. In another terminal
./bin/consumer

# 6. Watch metrics
curl http://localhost:9090/metrics | grep polymarket
```

### **Phase 3: Fork Testing (Anvil) - Week 2**

**Learn to test against real mainnet state locally**

#### **Setup Foundry/Anvil**

```bash
# Install Foundry (includes Anvil)
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

#### **Study Fork Tests**

Read: [test/fork_test.go](test/fork_test.go)

```go
func TestForkReadBalance(t *testing.T) {
    // 1. Load fork config (localhost:8545)
    cfg, _ := config.LoadConfig("../config/chains.json")
    chainCfg, _ := cfg.GetChain("polygon-fork")
    
    // 2. Create service
    svc, _ := service.NewCTFService(ctx, chainCfg)
    defer svc.Close()
    
    // 3. Read real mainnet data
    balance, _ := svc.BalanceOf(ctx, address, positionId)
}
```

#### **Practice:**

```bash
# Terminal 1: Start fork at block 55M
./scripts/start-fork.sh 55000000 --daemon

# Terminal 2: Run tests
go test ./test -v

# Check status
./scripts/check-fork.sh

# Stop
./scripts/stop-fork.sh
```

**Read:** [docs/FORK_TESTING_GUIDE.md](docs/FORK_TESTING_GUIDE.md)

### **Phase 4: Deep Dive - Go-Ethereum - Week 3**

**Master the go-ethereum library**

Study: [docs/GO_ETHEREUM_GUIDE.md](docs/GO_ETHEREUM_GUIDE.md) (3,300+ lines!)

**Key Topics:**

1. **Generating Go Bindings**
   ```bash
   abigen --abi pkg/contracts/abi/CTFExchange.json \
          --pkg contracts \
          --type CTFExchange \
          --out pkg/contracts/CTFExchange.go
   ```

2. **Reading Contract State**
   ```go
   // View/pure functions (no gas)
   orderStatus, err := ctfExchange.GetOrderStatus(&bind.CallOpts{}, orderHash)
   ```

3. **Sending Transactions**
   ```go
   auth := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
   auth.GasLimit = 300000
   tx, err := ctfExchange.FillOrder(auth, order, amount, sig)
   ```

4. **Waiting for Receipt**
   ```go
   receipt, err := bind.WaitMined(ctx, client, tx)
   if receipt.Status == 0 {
       // Transaction reverted
   }
   ```

5. **Parsing Event Logs**
   ```go
   // Manual parsing (what our handlers do)
   for _, log := range receipt.Logs {
       if log.Topics[0] == OrderFilledSig {
           event := parseOrderFilled(log)
       }
   }
   ```

**Practice:**
- Read [examples/usage_example.go](examples/usage_example.go)
- Study [pkg/service/ctf_service.go](pkg/service/ctf_service.go)

### **Phase 5: Production Patterns - Week 3-4**

**Learn production-grade resilience**

#### **Transaction Helpers**

Study: [docs/TRANSACTION_HELPERS.md](docs/TRANSACTION_HELPERS.md)

```go
// High-level helper (recommended)
tx, err := svc.ExecuteTransaction(ctx, msg, auth, func(opts) {
    return contract.DoSomething(opts, params...)
})

// Internally:
// 1. Simulates via eth_call
// 2. Estimates gas with 20% buffer
// 3. Retries up to 3 times on network errors
// 4. Returns transaction or permanent error
```

**Implementation:** [pkg/txhelper/transaction.go](pkg/txhelper/transaction.go)

#### **Retry with Exponential Backoff**

```go
func WithExponentialBackoff(ctx, maxRetries, initialBackoff) error {
    backoff := initialBackoff
    for attempt := 0; attempt < maxRetries; attempt++ {
        err := operation()
        if err == nil {
            return nil
        }
        
        if !isRetryable(err) {
            return err  // Permanent error
        }
        
        select {
        case <-time.After(backoff):
            backoff *= 2
            if backoff > maxBackoff {
                backoff = maxBackoff
            }
        case <-ctx.Done():
            return ctx.Err()
        }
    }
}
```

#### **Multi-RPC Failover**

Study patterns from `evm-scanner` (if available):

```go
type RPCPool struct {
    endpoints []*ethclient.Client
    current   int
}

func (p *RPCPool) Call(ctx, method, args) (result, error) {
    for attempt := 0; attempt < len(p.endpoints); attempt++ {
        client := p.endpoints[p.current]
        result, err := client.Call(ctx, method, args)
        
        if err == nil {
            return result, nil
        }
        
        // Switch to next RPC
        p.current = (p.current + 1) % len(p.endpoints)
    }
    return nil, errors.New("all RPCs failed")
}
```

### **Phase 6: NATS Deep Dive - Week 4**

**Master NATS JetStream patterns**

Study: [docs/NATS_EXPLAINED.md](docs/NATS_EXPLAINED.md) (1,500+ lines!)

**Key Concepts:**

1. **Why NATS over Kafka?**
   - 15MB binary vs 100MB+ JVM
   - 5-minute setup vs hours
   - Built-in deduplication
   - Perfect for startups

2. **Stream Configuration**
   ```go
   js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
       Name:       "POLYMARKET",
       Subjects:   []string{"POLYMARKET.*"},
       MaxAge:     7 * 24 * time.Hour,  // 7 days retention
       Duplicates: 20 * time.Minute,    // Dedup window
       Storage:    jetstream.FileStorage,
   })
   ```

3. **Publishing with Deduplication**
   ```go
   msgID := fmt.Sprintf("%s-%d", txHash, logIndex)
   js.Publish(ctx, subject, data, jetstream.WithMsgID(msgID))
   ```

4. **Durable Consumer**
   ```go
   cons, _ := js.CreateOrUpdateConsumer(ctx, "POLYMARKET", jetstream.ConsumerConfig{
       Name:    "polymarket-consumer",
       Durable: "polymarket-consumer",  // Survives restarts
   })
   ```

**Practice:**
- Run `docker-compose up nats`
- Access monitoring: http://localhost:8222
- Test with [internal/nats/publisher.go](internal/nats/publisher.go)

### **Phase 7: TimescaleDB & SQL - Week 5**

**Learn time-series database patterns**

Study: [migrations/001_initial_schema.up.sql](migrations/001_initial_schema.up.sql)

**Key Concepts:**

1. **Hypertables**
   ```sql
   CREATE TABLE events (...);
   SELECT create_hypertable('events', 'time');
   -- Automatically partitions by time (1-day chunks)
   ```

2. **Continuous Aggregates**
   ```sql
   CREATE MATERIALIZED VIEW trading_volume_hourly
   WITH (timescaledb.continuous) AS
   SELECT 
       time_bucket('1 hour', time) AS hour,
       COUNT(*) AS trades,
       SUM(maker_amount_filled) AS volume
   FROM order_fills
   GROUP BY hour;
   ```

3. **Idempotent Inserts**
   ```sql
   INSERT INTO events (...)
   VALUES (...)
   ON CONFLICT (tx_hash, log_index) DO NOTHING;
   ```

**Practice:**
```bash
# Connect to DB
docker-compose exec timescaledb psql -U polymarket

# Check hypertables
\dx timescaledb;
SELECT * FROM timescaledb_information.hypertables;

# Query events
SELECT * FROM events ORDER BY time DESC LIMIT 10;

# Check aggregates
SELECT * FROM trading_volume_hourly;
```

### **Phase 8: Build Missing Features - Week 6-8**

**Implement what's needed for production**

#### **Priority 1: Multi-RPC Failover**

```go
// TODO: internal/chain/pool.go
type ClientPool struct {
    clients []*ethclient.Client
    mu      sync.RWMutex
    current int
}

func (p *ClientPool) CallWithFailover(ctx, fn) (result, error) {
    // Try each client in round-robin
    // Track failures
    // Circuit breaker pattern
}
```

#### **Priority 2: Advanced Reorg Detection**

```go
// TODO: internal/syncer/reorg.go
func (s *Syncer) detectReorg(ctx, checkpoint) (bool, uint64) {
    block, err := s.chain.GetBlockByNumber(ctx, checkpoint.LastBlock)
    if err != nil {
        return false, 0
    }
    
    if block.Hash().Hex() != checkpoint.LastBlockHash {
        // Reorg detected!
        safeBlock := checkpoint.LastBlock - s.confirmations
        return true, safeBlock
    }
    
    return false, 0
}

func (s *Syncer) handleReorg(ctx, safeBlock) error {
    // 1. Stop processing
    // 2. Delete affected NATS messages (if possible)
    // 3. Clear affected DB records
    // 4. Update checkpoint to safe block
    // 5. Resume
}
```

#### **Priority 3: Monitoring & Alerts**

```go
// TODO: Add more Prometheus metrics
var (
    reorgDetected = prometheus.NewCounter(...)
    rpcFailures = prometheus.NewCounterVec(...)
    eventLatency = prometheus.NewHistogram(...)
)
```

#### **Priority 4: Graceful Shutdown Improvements**

```go
// Already implemented in main.go, enhance:
- Drain NATS publisher queue
- Finish processing current batch
- Close checkpoint cleanly
```

### **Phase 9: Testing Strategy - Week 9**

**Comprehensive testing**

#### **Unit Tests**

```go
// internal/handler/events_test.go
func TestHandleOrderFilled(t *testing.T) {
    log := types.Log{
        Topics: []common.Hash{OrderFilledSig, orderHash, maker, taker},
        Data:   hexToBytes("..."),
    }
    
    event, err := HandleOrderFilled(ctx, log, 12345678)
    require.NoError(t, err)
    assert.Equal(t, expectedEvent, event)
}
```

#### **Integration Tests**

```go
// test/integration_test.go
func TestEndToEndFlow(t *testing.T) {
    // 1. Start NATS, TimescaleDB (testcontainers)
    // 2. Start indexer
    // 3. Fork Polygon at specific block
    // 4. Verify events published to NATS
    // 5. Verify events stored in DB
}
```

#### **Fork Tests** (Already exist!)

```bash
go test ./test -v
```

**Read:** [docs/TESTING.md](docs/TESTING.md)

### **Phase 10: Production Deployment - Week 10**

**Deploy to production**

1. **Docker Build**
   ```bash
   make docker-build
   docker-compose up -d
   ```

2. **Kubernetes (Optional)**
   - Deploy indexer as Deployment (1-3 replicas)
   - Deploy consumer as Deployment (N replicas for scaling)
   - NATS as StatefulSet
   - TimescaleDB as StatefulSet with PVC

3. **Monitoring**
   - Prometheus + Grafana dashboards
   - Alerting rules
   - Log aggregation (ELK/Loki)

4. **Configuration**
   - Use environment variables in prod
   - Secrets management (Vault/K8s secrets)
   - Multiple RPC endpoints

**Read:** [docs/RUNNING.md](docs/RUNNING.md)

## 🏗️ Complete Project Structure (What You Have)

```
polymarket-indexer/
├── cmd/
│   ├── indexer/              ✅ Main indexer entry point
│   │   └── main.go           ✅ Wiring, signal handling, graceful shutdown
│   └── consumer/             ✅ NATS → TimescaleDB consumer
│       └── main.go           ✅ Event consumption, DB writes, metrics
│
├── internal/                 🔒 Private packages (internal use only)
│   ├── chain/
│   │   └── client.go         ✅ Dual RPC client (HTTP + WebSocket)
│   ├── syncer/
│   │   └── syncer.go         ✅ Block sync orchestration (backfill/realtime)
│   ├── processor/
│   │   └── processor.go      ✅ Block/log processing, event routing
│   ├── handler/
│   │   └── events.go         ✅ 9 event handlers (ABI decoding)
│   ├── router/
│   │   └── router.go         ✅ Event signature → handler mapping
│   ├── nats/
│   │   └── publisher.go      ✅ JetStream publisher with deduplication
│   ├── db/
│   │   └── checkpoint.go     ✅ BoltDB checkpoint storage
│   └── util/
│       └── init.go           ✅ Logger, config initialization
│
├── pkg/                      📦 Public packages (can be imported)
│   ├── models/
│   │   └── event.go          ✅ Event data models
│   ├── contracts/            ✅ Generated Go bindings
│   │   ├── CTFExchange.go    ✅ From abigen
│   │   ├── ConditionalTokens.go ✅
│   │   ├── ERC20.go          ✅
│   │   └── abi/              ✅ ABI JSON files
│   ├── service/
│   │   └── ctf_service.go    ✅ High-level CTF service
│   ├── txhelper/
│   │   └── transaction.go    ✅ Transaction helpers (retry, simulation)
│   └── config/
│       └── config.go         ✅ Chain config loader
│
├── test/
│   └── fork_test.go          ✅ Anvil fork testing examples
│
├── migrations/
│   └── 001_initial_schema.up.sql ✅ TimescaleDB schema
│
├── scripts/
│   ├── start-fork.sh         ✅ Fork testing (daemon mode)
│   ├── stop-fork.sh          ✅
│   ├── check-fork.sh         ✅
│   └── generate-bindings.sh  ✅ Regenerate contract bindings
│
├── examples/
│   ├── usage_example.go      ✅ CTF service examples
│   └── transaction_patterns.go ✅
│
├── docs/                     📚 Comprehensive documentation
│   ├── ARCHITECTURE.md       ✅ System design deep dive (586 lines)
│   ├── GO_ETHEREUM_GUIDE.md  ✅ Go-ethereum tutorial (3,352 lines!)
│   ├── TRANSACTION_HELPERS.md ✅ Production transaction patterns
│   ├── NATS_EXPLAINED.md     ✅ NATS vs Kafka, patterns (1,523 lines)
│   ├── DEPENDENCIES.md       ✅ Every dependency explained (695 lines)
│   ├── FORK_TESTING_GUIDE.md ✅ Anvil testing guide
│   ├── DAEMON_MODE.md        ✅ Background process management
│   ├── TESTING.md            ✅ Testing strategies (682 lines)
│   └── RUNNING.md            ✅ Deployment guide (457 lines)
│
├── config.toml               ✅ Main configuration
├── docker-compose.yml        ✅ Full stack (NATS, TimescaleDB, etc.)
├── Dockerfile                ✅ Multi-stage build
├── Makefile                  ✅ Build, test, deploy commands
├── go.mod                    ✅ Dependencies
├── README.md                 ✅ Project overview (236 lines)
├── PROJECT_SETUP.md          ✅ Setup guide (297 lines)
├── DAEMON_MODE_README.md     ✅ Quick daemon mode ref
├── FORK_TESTING_README.md    ✅ Quick fork testing ref
└── LEARNING_PATH.md          ✅ This file
```

## 🔧 Tech Stack & Dependencies

### Core Dependencies (All Explained in docs/DEPENDENCIES.md)

```bash
# Essential for blockchain interaction
github.com/ethereum/go-ethereum v1.13.14
  • Official Ethereum implementation (Geth)
  • ethclient, types, ABI encoding/decoding
  • Used: Everywhere (chain client, contracts, event parsing)
  • Size: ~50MB binary
  
# PostgreSQL driver (5-10x faster than database/sql)
github.com/jackc/pgx/v5 v5.5.5
  • Native PostgreSQL wire protocol
  • Binary protocol, connection pooling
  • Used: Consumer service, TimescaleDB writes
  • Why: Native TimescaleDB support, batch operations
  
# Configuration management
github.com/knadh/koanf/v2 v2.1.0
  • TOML, env vars, hot-reload
  • Used: Loading config.toml, env overrides
  • Why: 12-factor app pattern, production-ready
  
# NATS messaging
github.com/nats-io/nats.go v1.34.0
  • JetStream client with deduplication
  • Used: Event publishing, consumption
  • Why: Simpler than Kafka, perfect for startups
  
# Metrics
github.com/prometheus/client_golang v1.19.0
  • Prometheus instrumentation
  • Used: /metrics endpoint, counters, gauges, histograms
  • Why: Industry standard monitoring
  
# Logging
github.com/rs/zerolog v1.32.0
  • Zero-allocation JSON logger
  • Used: Structured logging throughout
  • Why: 10x faster than stdlib, JSON output
  
# Checkpoint storage
go.etcd.io/bbolt v1.3.9
  • Embedded key-value store (BoltDB fork)
  • Used: Block checkpoint persistence
  • Why: Single file, ACID, no dependencies
```

**Full dependency guide:** [docs/DEPENDENCIES.md](docs/DEPENDENCIES.md) (695 lines)

### Installation

```bash
# All dependencies managed by go.mod
go mod download

# Or let build handle it
make build
```

## 🎓 Key Learning Points from Implementation

### 1. **Event Parsing with ABI (Already Implemented!)**

**See:** [internal/handler/events.go](internal/handler/events.go)

```go
// How ABI decoding works in handlers:
func HandleOrderFilled(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
    // 1. Check topic count (indexed parameters + signature)
    if len(log.Topics) != 4 {
        return nil, fmt.Errorf("invalid event")
    }
    
    // 2. Extract indexed parameters from topics
    // Topic[0] = event signature (OrderFilledSig)
    // Topic[1] = orderHash (indexed bytes32)
    // Topic[2] = maker (indexed address)
    // Topic[3] = taker (indexed address)
    orderHash := log.Topics[1].Hex()
    maker := common.BytesToAddress(log.Topics[2].Bytes()).Hex()
    taker := common.BytesToAddress(log.Topics[3].Bytes()).Hex()
    
    // 3. Extract non-indexed parameters from data
    // Data contains 5 * 32 bytes = 160 bytes
    // Each uint256 takes 32 bytes
    makerAssetID := new(big.Int).SetBytes(log.Data[0:32])
    takerAssetID := new(big.Int).SetBytes(log.Data[32:64])
    makerAmountFilled := new(big.Int).SetBytes(log.Data[64:96])
    takerAmountFilled := new(big.Int).SetBytes(log.Data[96:128])
    fee := new(big.Int).SetBytes(log.Data[128:160])
    
    // 4. Return strongly-typed struct
    return models.OrderFilled{
        OrderHash:         orderHash,
        Maker:             maker,
        Taker:             taker,
        MakerAssetID:      makerAssetID,
        TakerAssetID:      takerAssetID,
        MakerAmountFilled: makerAmountFilled,
        TakerAmountFilled: takerAmountFilled,
        Fee:               fee,
    }, nil
}
```

**For dynamic arrays (TransferBatch):**
```go
// Need abi.Arguments for complex types
uint256ArrayTy, _ := abi.NewType("uint256[]", "", nil)
args := abi.Arguments{
    {Type: uint256ArrayTy},
    {Type: uint256ArrayTy},
}

unpacked, err := args.Unpack(log.Data)
tokenIDs := unpacked[0].([]*big.Int)
amounts := unpacked[1].([]*big.Int)
```

**Generate bindings from ABI:**
```bash
abigen --abi pkg/contracts/abi/CTFExchange.json \
       --pkg contracts \
       --type CTFExchange \
       --out pkg/contracts/CTFExchange.go
```

### 2. **NATS Deduplication Pattern (Implemented)**

**See:** [internal/nats/publisher.go](internal/nats/publisher.go)

```go
// Message ID construction prevents duplicates
func (p *Publisher) Publish(ctx context.Context, event models.Event) error {
    // 1. Construct hierarchical subject
    subject := fmt.Sprintf("%s.%s.%s", 
        p.prefix,           // "POLYMARKET"
        event.EventName,    // "OrderFilled"
        event.ContractAddr) // "0x4bFb..."
    
    // 2. Create unique message ID
    // Format: {txHash}-{logIndex}
    // Example: "0xabc123...-5"
    msgID := fmt.Sprintf("%s-%d", event.TxHash, event.LogIndex)
    
    // 3. Publish with deduplication
    // NATS checks: "Already have msgID=0xabc123-5? Skip storage."
    _, err = p.js.Publish(ctx, subject, data, jetstream.WithMsgID(msgID))
    
    return err
}
```

**Why this works:**
- Transaction hash + log index is globally unique
- If indexer restarts and reprocesses same block, NATS rejects duplicate
- Deduplication window: 20 minutes (configurable)
- Survives crash recovery, reorg reprocessing

**Stream configuration:**
```go
js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
    Name:       "POLYMARKET",
    Subjects:   []string{"POLYMARKET.*"},     // Wildcard
    MaxAge:     7 * 24 * time.Hour,          // 7 days retention
    Duplicates: 20 * time.Minute,            // Dedup window
    Storage:    jetstream.FileStorage,       // Persistent
})
```

### 3. **Checkpoint Management (Implemented)**

**See:** [internal/db/checkpoint.go](internal/db/checkpoint.go)

```go
// Checkpoint structure
type Checkpoint struct {
    ServiceName   string    // "polymarket-indexer"
    LastBlock     uint64    // Last processed block number
    LastBlockHash string    // Block hash (for reorg detection)
    UpdatedAt     time.Time // Last update timestamp
}

// Get or create on startup
checkpoint, err := db.GetOrCreateCheckpoint(ctx, "polymarket-indexer", startBlock)

// Update after processing block
err = db.UpdateBlock(ctx, "polymarket-indexer", blockNum, blockHash)
```

**BoltDB storage (embedded):**
```go
// Single-file database at data/checkpoints.db
// ACID transactions
// No external dependencies
// Perfect for simple key-value needs
```

**Crash recovery:**
```
1. Service crashes at block 100,000
2. Checkpoint stored: {block: 99,950, hash: "0xdef..."}
3. Service restarts
4. Loads checkpoint: 99,950
5. Resumes from 99,951
6. NATS deduplication prevents duplicate events
```

### 4. **Router Pattern (Implemented)**

**See:** [internal/router/router.go](internal/router/router.go)

```go
// Registry pattern for extensibility
type Router struct {
    callback    EventCallback                    // NATS publish
    logHandlers map[common.Hash]LogHandlerFunc  // Sig → Handler
    eventNames  map[common.Hash]string          // Sig → Name
}

// Registration in processor
router := router.New(logger)
router.RegisterLogHandler(OrderFilledSig, "OrderFilled", HandleOrderFilled)
router.RegisterLogHandler(TransferSingleSig, "TransferSingle", HandleTransferSingle)
// ... 9 handlers total

// Routing
func (r *Router) RouteLog(ctx, log, timestamp, blockHash) error {
    eventSig := log.Topics[0]  // First topic = event signature
    
    handler, exists := r.logHandlers[eventSig]
    if !exists {
        return nil  // No handler, skip
    }
    
    // Execute handler to parse event
    payload, err := handler(ctx, log, timestamp)
    
    // Wrap in Event envelope
    event := models.Event{
        Block:        log.BlockNumber,
        EventName:    r.eventNames[eventSig],
        Payload:      payload,
        // ... other fields
    }
    
    // Call callback (NATS publish)
    return r.callback(ctx, event)
}
```

**Benefits:**
- Extensible: Add new events by registering handlers
- Type-safe: Handlers return any (strongly typed structs)
- Testable: Mock callback for unit tests
- Maintainable: Each handler in separate function

### 5. **Syncer Modes (Implemented)**

**See:** [internal/syncer/syncer.go](internal/syncer/syncer.go)

```go
// Two modes: Backfill and Realtime
func (s *Syncer) Start(ctx context.Context) error {
    for {
        latestBlock := s.chain.GetLatestBlockNumber(ctx)
        behind := latestBlock - s.currentBlock - s.confirmations
        
        if behind > s.batchSize {
            // BACKFILL MODE
            // - Process 100 blocks/batch
            // - Use 5 parallel workers
            // - Max throughput
            s.backfillMode(ctx)
        } else {
            // REALTIME MODE
            // - Process blocks sequentially
            // - Wait for new blocks
            // - Maintain order
            s.realtimeMode(ctx)
        }
    }
}
```

**Backfill worker pool:**
```go
jobs := make(chan uint64, batchSize)
results := make(chan error, batchSize)

// Spawn workers
for w := 0; w < s.workers; w++ {
    go func() {
        for blockNum := range jobs {
            err := s.processor.ProcessBlock(ctx, blockNum)
            results <- err
        }
    }()
}

// Queue blocks
for block := start; block < end; block++ {
    jobs <- block
}
```

### 6. **Production Transaction Patterns (Implemented)**

**See:** [pkg/txhelper/transaction.go](pkg/txhelper/transaction.go)

```go
// High-level helper wraps all best practices
func ExecuteTransaction(ctx, msg, auth, sendFunc) (*types.Transaction, error) {
    config := &TransactionConfig{
        MaxRetries:       3,
        InitialBackoff:   1 * time.Second,
        MaxBackoff:       30 * time.Second,
        GasBufferPercent: 20,
        Simulate:         true,
        TimeoutPerTry:    30 * time.Second,
    }
    
    return SendTransactionWithRetry(ctx, msg, auth, config, sendFunc)
}
```

**Flow:**
1. **Simulate** via `eth_call` to catch reverts early
2. **Estimate gas** and add 20% buffer
3. **Send transaction**
4. On error:
   - Classify: Retryable (network) vs Permanent (revert)
   - Retryable: Wait with exponential backoff, retry
   - Permanent: Return error immediately
5. Return transaction or error

**Error classification:**
```go
func isRetryable(err error) bool {
    // Retry on:
    // - Connection refused/reset
    // - Timeouts
    // - Rate limiting (429)
    // - Gateway errors (502, 503, 504)
    
    // Don't retry on:
    // - execution reverted
    // - insufficient funds
    // - nonce too low
    // - gas too low
}
```

### 7. **Fork Testing (Fully Implemented!)**

**See:** [test/fork_test.go](test/fork_test.go), [scripts/start-fork.sh](scripts/start-fork.sh)

```bash
# Start fork at block 55M (daemon mode)
./scripts/start-fork.sh 55000000 --daemon

# Features:
# - Local copy of Polygon mainnet at specific block
# - Impersonate any address (no private key!)
# - Manipulate state (balances, storage)
# - Time travel (evm_increaseTime)
# - Deterministic testing
```

**Test pattern:**
```go
func TestForkReadData(t *testing.T) {
    // 1. Load fork config (localhost:8545)
    cfg, _ := config.LoadConfig("../config/chains.json")
    chainCfg, _ := cfg.GetChain("polygon-fork")
    
    // 2. Create service
    svc, _ := service.NewCTFService(ctx, chainCfg)
    
    // 3. Read real mainnet data
    balance, _ := svc.BalanceOf(ctx, realAddress, positionId)
    
    // 4. Assert expectations
    assert.True(t, balance.Cmp(big.NewInt(0)) > 0)
}
```

**Guides:**
- [docs/FORK_TESTING_GUIDE.md](docs/FORK_TESTING_GUIDE.md) - Comprehensive guide (473 lines)
- [docs/DAEMON_MODE.md](docs/DAEMON_MODE.md) - Background process management

## 📖 Resources to Study

1. **Go Ethereum (Geth) Docs**
   - https://geth.ethereum.org/docs/developers/geth-developer/dev-guide

2. **NATS JetStream**
   - https://docs.nats.io/nats-concepts/jetstream

3. **TimescaleDB**
   - https://docs.timescale.com/

4. **Polymarket Contracts**
   - CTF Exchange: `0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E`
   - Conditional Tokens: `0x4D97DCd97eC945f40cF65F87097ACe5EA0476045`
   - Polygon RPC: https://polygon-rpc.com

## 🚀 Next Steps

1. **Study eth-tracker code** (focus on NATS integration)
2. **Run eth-tracker locally** to see it in action
3. **Build simple USDC indexer** as practice
4. **Add NATS** to your simple indexer
5. **Build Polymarket indexer** using the same patterns

## 💡 Pro Tips

1. **Start simple** - Don't build everything at once
2. **Copy patterns** - Both reference projects have battle-tested code
3. **Test locally first** - Use local NATS and TimescaleDB
4. **Handle errors properly** - Blockchain data is messy
5. **Monitor everything** - Metrics are crucial in production
6. **Read the code** - The reference projects are your best teachers

## 🤝 Need Help?

- eth-tracker issues: https://github.com/grassrootseconomics/eth-tracker/issues
- evm-scanner issues: https://github.com/84hero/evm-scanner/issues
- Go Ethereum: https://ethereum.stackexchange.com/

---

**Remember:** Whether indexing Polymarket, Uniswap, or any DeFi protocol, the core patterns remain the same. Master these patterns from the reference projects, and you can index anything!
