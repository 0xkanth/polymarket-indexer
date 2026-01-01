# Go Dependencies Guide

Complete explanation of all dependencies in the Polymarket Indexer project.

## Quick Reference Table

| Dependency | Version | Type | What It Does |
|-----------|---------|------|--------------|
| **github.com/ethereum/go-ethereum** | v1.13.14 | Direct | Official Ethereum/Geth library - connects to Polygon RPC, fetches blocks/transactions, decodes events, provides ABI encoding/decoding and cryptographic utilities for blockchain interaction. |
| **github.com/jackc/pgx/v5** | v5.5.5 | Direct | High-performance PostgreSQL driver - 5-10x faster than database/sql, native TimescaleDB support, efficient batch operations and connection pooling for storing event data. |
| **github.com/knadh/koanf/v2** | v2.1.0 | Direct | Flexible configuration management - loads config from TOML files and environment variables, provides type-safe access to settings like RPC endpoints and contract addresses. |
| **github.com/nats-io/nats.go** | v1.34.0 | Direct | NATS messaging client with JetStream - publishes events to persistent streams with deduplication, enables decoupled architecture and horizontal scaling of consumers. |
| **github.com/prometheus/client_golang** | v1.19.0 | Direct | Prometheus metrics instrumentation - tracks performance metrics (blocks/sec, events/sec, errors) and exposes `/metrics` endpoint for monitoring and alerting. |
| **github.com/rs/zerolog** | v1.32.0 | Direct | Zero-allocation JSON logger - 10x faster than stdlib, structured logging for production, outputs machine-readable JSON logs with zero heap allocations. |
| **go.etcd.io/bbolt** | v1.3.9 | Direct | Embedded key-value database - stores block checkpoints for resume capability, single-file BoltDB fork with ACID transactions and no external dependencies. |
| **github.com/btcsuite/btcd/btcec/v2** | v2.3.2 | Indirect | Elliptic curve cryptography (secp256k1) - used by go-ethereum for Ethereum address generation and signature verification. |
| **github.com/decred/dcrd/dcrec/secp256k1/v4** | v4.2.0 | Indirect | Optimized secp256k1 implementation - alternative crypto library for performance-critical operations in Ethereum transactions. |
| **github.com/holiman/uint256** | v1.2.4 | Indirect | 256-bit integer arithmetic - faster than big.Int for Ethereum values like balances, gas prices, and token amounts. |
| **golang.org/x/crypto** | v0.21.0 | Indirect | Go cryptographic library - provides Keccak256 hashing, ECDSA signatures, and other primitives required by Ethereum. |
| **github.com/jackc/pgpassfile** | v1.0.0 | Indirect | PostgreSQL password file parser - reads `.pgpass` file for secure credential management. |
| **github.com/jackc/pgservicefile** | v0.0.0 | Indirect | PostgreSQL service file parser - supports `pg_service.conf` for connection aliases and configuration. |
| **github.com/jackc/puddle/v2** | v2.2.1 | Indirect | Generic resource pool - manages pgx connection pool with health checks and lifecycle management. |
| **github.com/knadh/koanf/maps** | v0.1.1 | Indirect | Map utilities for koanf - handles merging and manipulation of configuration maps. |
| **github.com/knadh/koanf/parsers/toml** | v0.1.0 | Indirect | TOML parser integration - enables koanf to parse TOML configuration files. |
| **github.com/knadh/koanf/providers/env** | v0.1.0 | Indirect | Environment variable provider - allows config override from environment variables for Docker deployment. |
| **github.com/knadh/koanf/providers/file** | v0.1.0 | Indirect | File system provider - loads configuration from files like config.toml. |
| **github.com/pelletier/go-toml** | v1.9.5 | Indirect | TOML parsing library - underlying parser used by koanf for TOML format support. |
| **github.com/fsnotify/fsnotify** | v1.7.0 | Indirect | File system notifications - enables hot-reload of configuration files when they change. |
| **github.com/go-viper/mapstructure/v2** | v2.0.0-alpha.1 | Indirect | Map to struct decoder - converts configuration maps into strongly-typed Go structs. |
| **github.com/mitchellh/copystructure** | v1.2.0 | Indirect | Deep copy utility - creates deep copies of Go data structures for configuration manipulation. |
| **github.com/mitchellh/reflectwalk** | v1.0.2 | Indirect | Reflection walker - traverses Go structures using reflection for copystructure. |
| **github.com/nats-io/nkeys** | v0.4.7 | Indirect | NKeys authentication - Ed25519-based secure authentication for NATS connections. |
| **github.com/nats-io/nuid** | v1.0.1 | Indirect | Unique ID generation - high-performance collision-resistant IDs for NATS messages. |
| **github.com/klauspost/compress** | v1.17.7 | Indirect | Fast compression algorithms - compresses NATS messages for efficient network transmission. |
| **github.com/beorn7/perks** | v1.0.1 | Indirect | Quantile estimation - statistical calculations for Prometheus summary metrics. |
| **github.com/cespare/xxhash/v2** | v2.2.0 | Indirect | Fast hash function - non-cryptographic hashing for Prometheus metric label optimization. |
| **github.com/prometheus/client_model** | v0.6.0 | Indirect | Prometheus data model - defines metric types (Counter, Gauge, Histogram, Summary). |
| **github.com/prometheus/common** | v0.50.0 | Indirect | Common Prometheus utilities - shared code for Prometheus client libraries. |
| **github.com/prometheus/procfs** | v0.13.0 | Indirect | /proc filesystem parser - collects system metrics like CPU, memory, disk from Linux. |
| **google.golang.org/protobuf** | v1.33.0 | Indirect | Protocol Buffers - binary serialization format used by Prometheus for efficient metric transport. |
| **github.com/mattn/go-colorable** | v0.1.13 | Indirect | Colored terminal output - enables ANSI colors in Windows terminals for zerolog. |
| **github.com/mattn/go-isatty** | v0.0.20 | Indirect | TTY detection - checks if output is a terminal to enable/disable colored logs. |
| **golang.org/x/sync** | v0.6.0 | Indirect | Advanced synchronization - provides errgroup for goroutine coordination and error handling. |
| **golang.org/x/sys** | v0.18.0 | Indirect | System calls - low-level OS operations for cross-platform system access. |
| **golang.org/x/text** | v0.14.0 | Indirect | Text processing - Unicode handling, encoding, and internationalization support. |

---

## Direct Dependencies (Explicitly Used)

These are the libraries you'll directly import and use in your code.

### 1. github.com/ethereum/go-ethereum v1.13.14

**Official Ethereum implementation in Go (Geth)**

**Core Functionality:**
- **ethclient** - Connect to Ethereum/Polygon RPC nodes via HTTP/WebSocket
- **types** - Block, Transaction, Receipt, Log data structures
- **common** - Address, Hash types and utilities
- **accounts/abi** - ABI encoding/decoding for smart contracts
- **crypto** - Keccak256 hashing, ECDSA signatures
- **event** - Event subscription and filtering

**Why You Need It:**
- Fetch blocks from Polygon chain
- Decode event logs from CTF Exchange and Conditional Tokens contracts
- Parse transaction receipts for event extraction
- Filter logs by contract address and event signature
- Generate contract bindings with `abigen` tool

**Example Usage:**
```go
client, err := ethclient.Dial("https://polygon-rpc.com")
block, err := client.BlockByNumber(ctx, big.NewInt(20558323))
receipt, err := client.TransactionReceipt(ctx, txHash)
logs, err := client.FilterLogs(ctx, ethereum.FilterQuery{
    Addresses: []common.Address{ctfExchangeAddr},
})
```

**Binary Size:** ~50MB (includes full Ethereum consensus logic)

---

### 2. github.com/jackc/pgx/v5 v5.5.5

**The best PostgreSQL driver for Go**

**Core Functionality:**
- Native PostgreSQL wire protocol implementation
- Connection pooling with `pgxpool`
- Batch operations for high throughput
- Binary protocol support (faster than text)
- COPY protocol for bulk inserts
- Full PostgreSQL type support (JSONB, arrays, timestamps)

**Why You Need It:**
- Store events in TimescaleDB (PostgreSQL extension)
- Insert into `order_fills`, `token_transfers`, `conditions` tables
- Execute efficient queries with proper parameterization
- Batch insert events for maximum performance
- Work with TimescaleDB hypertables and continuous aggregates

**Performance:**
- **5-10x faster** than `database/sql` for bulk operations
- Native support for PostgreSQL arrays and JSONB
- Better error messages and type safety
- Connection pool with automatic health checks

**Example Usage:**
```go
pool, err := pgxpool.New(ctx, "postgres://user:pass@localhost/polymarket")
defer pool.Close()

// Single insert
_, err = pool.Exec(ctx, 
    "INSERT INTO order_fills (time, maker, taker, maker_amount) VALUES ($1, $2, $3, $4)",
    time.Now(), maker, taker, amount)

// Batch insert (fast)
batch := &pgx.Batch{}
for _, fill := range fills {
    batch.Queue("INSERT INTO order_fills (...) VALUES (...)", fill.Time, fill.Maker, ...)
}
results := pool.SendBatch(ctx, batch)
```

**Binary Size:** ~2MB

---

### 3. github.com/knadh/koanf/v2 v2.1.0

**Simple, powerful configuration management**

**Core Functionality:**
- Load config from multiple sources (files, env vars, command-line flags)
- Support for TOML, JSON, YAML, HCL formats
- Environment variable overrides (12-factor app pattern)
- Type-safe config access with `MustString()`, `MustInt()`, etc.
- Hot-reload capability with file watching
- Nested configuration with dot notation

**Why You Need It:**
- Parse `config.toml` file with all settings
- Override with environment variables for Docker deployments
- Access config values safely: `ko.MustString("chain.rpc_endpoint")`
- Production-grade config management (battle-tested, used by eth-tracker)

**Example Usage:**
```go
ko := koanf.New(".")

// Load from TOML file
ko.Load(file.Provider("config.toml"), toml.Parser())

// Override from environment variables (CHAIN_RPC_ENDPOINT -> chain.rpc_endpoint)
ko.Load(env.Provider("", ".", func(s string) string {
    return strings.Replace(strings.ToLower(s), "_", ".", -1)
}), nil)

// Access values
rpcURL := ko.MustString("chain.rpc_endpoint")
chainID := ko.MustInt64("chain.chainid")
startBlock := ko.Int64("chain.start_block") // 0 if not set
```

**Binary Size:** ~1MB

---

### 4. github.com/nats-io/nats.go v1.34.0

**NATS messaging with JetStream persistence**

**Core Functionality:**
- Connect to NATS server
- Publish/Subscribe messaging patterns
- **JetStream** - Persistent, replicated, durable streams
- Message deduplication (prevents duplicate processing)
- Acknowledgements and automatic retries
- Consumer groups for load balancing
- Key-Value and Object stores

**Why You Need It:**
- Publish blockchain events to NATS JetStream
- Deduplication using Message ID (`txHash-logIndex`)
- Consumer subscribes to event streams for database writes
- Decouples indexer from database (failure isolation)
- Enables horizontal scaling of consumers

**JetStream Features:**
- **Persistence**: Messages survive NATS server restarts
- **Replay**: Re-process historical events from any point
- **Deduplication**: 20-minute window prevents duplicate processing
- **At-least-once delivery**: Guarantees no message loss
- **Stream limits**: Automatic message retention policies

**Example Usage:**
```go
// Connect and create stream
nc, _ := nats.Connect("nats://localhost:4222")
js, _ := jetstream.New(nc)

// Create persistent stream with deduplication
js.CreateStream(ctx, jetstream.StreamConfig{
    Name:       "POLYMARKET",
    Subjects:   []string{"POLYMARKET.*"},
    Storage:    jetstream.FileStorage,
    Duplicates: 20 * time.Minute,
})

// Publish with deduplication
msgID := fmt.Sprintf("%s-%d", txHash, logIndex)
js.Publish(ctx, "POLYMARKET.OrderFilled.0x4bFb...", eventData, 
    nats.MsgId(msgID))

// Subscribe
js.Subscribe("POLYMARKET.OrderFilled.*", func(msg *nats.Msg) {
    // Process event
    msg.Ack()
})
```

**Binary Size:** ~5MB

---

### 5. github.com/prometheus/client_golang v1.19.0

**Prometheus metrics instrumentation**

**Core Functionality:**
- **Counter** - Monotonically increasing values (total blocks processed)
- **Gauge** - Values that go up/down (current block height)
- **Histogram** - Distribution of values (latency buckets)
- **Summary** - Quantiles over sliding time window
- HTTP handler for `/metrics` endpoint

**Why You Need It:**
- Monitor indexer performance in real-time
- Track key metrics: blocks/sec, events/sec, errors, processing latency
- Expose metrics at `http://localhost:8080/metrics`
- Integration with Prometheus server + Grafana dashboards
- Production monitoring, alerting, and capacity planning

**Key Metrics You'll Track:**
```go
var (
    blocksProcessed = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "polymarket_blocks_processed_total",
        Help: "Total number of blocks processed",
    })
    
    eventsPublished = prometheus.NewCounterVec(prometheus.CounterOpts{
        Name: "polymarket_events_published_total",
        Help: "Total events published by type",
    }, []string{"event_type"})
    
    processingLatency = prometheus.NewHistogram(prometheus.HistogramOpts{
        Name:    "polymarket_block_processing_seconds",
        Help:    "Block processing latency",
        Buckets: prometheus.DefBuckets,
    })
    
    currentBlock = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "polymarket_current_block",
        Help: "Current block being processed",
    })
)
```

**Example Usage:**
```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

// Register metrics
prometheus.MustRegister(blocksProcessed, eventsPublished, processingLatency)

// Expose endpoint
http.Handle("/metrics", promhttp.Handler())
http.ListenAndServe(":8080", nil)

// Track metrics
start := time.Now()
blocksProcessed.Inc()
eventsPublished.WithLabelValues("OrderFilled").Add(5)
processingLatency.Observe(time.Since(start).Seconds())
currentBlock.Set(20558323)
```

**Binary Size:** ~3MB

---

### 6. github.com/rs/zerolog v1.32.0

**Zero-allocation JSON logger**

**Core Functionality:**
- Structured logging with key-value pairs
- JSON output (machine-readable, easily parsed)
- **Zero heap allocations** (no garbage collection pressure)
- Log levels: debug, info, warn, error, fatal, panic
- Context-aware logging with chained methods
- Pretty console output for development
- Sampling for high-volume logs

**Why You Need It:**
- Production-grade logging infrastructure
- JSON logs for centralized logging (Elasticsearch, Datadog, etc.)
- Performance-critical: no allocations = no GC pauses
- Used by eth-tracker (proven in production at scale)

**Performance Comparison:**
- **10x faster** than standard library `log` package
- **Zero allocations** - crucial for high-throughput indexing
- Can handle 1M+ log messages per second

**Example Usage:**
```go
// Initialize
log := zerolog.New(os.Stdout).With().
    Timestamp().
    Str("service", "polymarket-indexer").
    Logger()

// Structured logging
log.Info().
    Str("contract", "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E").
    Uint64("block", 20558323).
    Int("events", 5).
    Dur("duration", processingTime).
    Msg("processed block")

// Output (JSON):
// {"level":"info","service":"polymarket-indexer","time":"2025-12-31T12:00:00Z",
//  "contract":"0x4bFb...","block":20558323,"events":5,"duration":125,
//  "message":"processed block"}

// Pretty console (dev)
log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
```

**Binary Size:** ~500KB

---

### 7. go.etcd.io/bbolt v1.3.9

**Embedded key-value database (BoltDB fork)**

**Core Functionality:**
- Pure Go embedded database (no C dependencies)
- ACID transactions (Atomicity, Consistency, Isolation, Durability)
- B+tree storage engine
- Single file database (`checkpoints.db`)
- Memory-mapped I/O for fast reads
- Bucket-based organization (like folders)

**Why You Need It:**
- Store checkpoint (last processed block + block hash)
- Persists state across restarts (resume from last block)
- Fast, reliable, battle-tested
- No external database server required
- Used by eth-tracker for checkpoint storage

**Key Advantages:**
- **Embedded**: No separate server process, just a file
- **Fast**: Memory-mapped reads, optimized writes
- **Reliable**: ACID guarantees, crash-safe
- **Simple**: Single 100KB library, no configuration

**Example Usage:**
```go
// Open database
db, err := bbolt.Open("data/checkpoints.db", 0600, nil)
defer db.Close()

// Write checkpoint
db.Update(func(tx *bbolt.Tx) error {
    b, _ := tx.CreateBucketIfNotExists([]byte("checkpoints"))
    checkpoint := Checkpoint{
        LastBlock:     20558323,
        LastBlockHash: "0xabc...",
        Timestamp:     time.Now(),
    }
    data, _ := json.Marshal(checkpoint)
    return b.Put([]byte("indexer"), data)
})

// Read checkpoint
var checkpoint Checkpoint
db.View(func(tx *bbolt.Tx) error {
    b := tx.Bucket([]byte("checkpoints"))
    data := b.Get([]byte("indexer"))
    return json.Unmarshal(data, &checkpoint)
})
```

**Binary Size:** ~100KB

---

## Indirect Dependencies (Transitive)

These are pulled in automatically by your direct dependencies. You typically don't import these directly.

### Ethereum Dependencies

#### github.com/btcsuite/btcd/btcec/v2 v2.3.2
Elliptic curve cryptography (secp256k1) for Ethereum address generation and signature verification.

#### github.com/decred/dcrd/dcrec/secp256k1/v4 v4.2.0
Alternative performance-optimized secp256k1 implementation for crypto operations.

#### github.com/holiman/uint256 v1.2.4
256-bit integer arithmetic, faster than `big.Int` for Ethereum values like balances and gas prices.

#### golang.org/x/crypto v0.21.0
Go's cryptographic library providing Keccak256 hashing (Ethereum's hash function), ECDSA signatures, and other primitives.

---

### PostgreSQL Dependencies

#### github.com/jackc/pgpassfile v1.0.0
Parses PostgreSQL `.pgpass` file for secure credential storage and automatic authentication.

#### github.com/jackc/pgservicefile v0.0.0-20231201235250-de7065d80cb9
Parses `pg_service.conf` for PostgreSQL connection aliases and service definitions.

#### github.com/jackc/puddle/v2 v2.2.1
Generic resource pool implementation used by pgxpool for connection pooling, health checks, and lifecycle management.

---

### Configuration Dependencies

#### github.com/knadh/koanf/maps v0.1.1
Map manipulation utilities for merging configuration from multiple sources.

#### github.com/knadh/koanf/parsers/toml v0.1.0
TOML parser integration for koanf.

#### github.com/knadh/koanf/providers/env v0.1.0
Environment variable provider for 12-factor app configuration.

#### github.com/knadh/koanf/providers/file v0.1.0
File system provider for loading config.toml files.

#### github.com/pelletier/go-toml v1.9.5
Underlying TOML parsing library used by koanf.

#### github.com/fsnotify/fsnotify v1.7.0
File system notifications for hot-reload of configuration files.

#### github.com/go-viper/mapstructure/v2 v2.0.0-alpha.1
Decodes configuration maps into strongly-typed Go structs.

#### github.com/mitchellh/copystructure v1.2.0
Deep copy of Go data structures for configuration manipulation.

#### github.com/mitchellh/reflectwalk v1.0.2
Walks Go data structures using reflection (used by copystructure).

---

### NATS Dependencies

#### github.com/nats-io/nkeys v0.4.7
NKeys authentication using Ed25519 signatures for secure NATS connections.

#### github.com/nats-io/nuid v1.0.1
High-performance unique ID generation for NATS messages (collision-resistant).

#### github.com/klauspost/compress v1.17.7
Fast compression algorithms for NATS message compression (reduces bandwidth).

---

### Prometheus Dependencies

#### github.com/beorn7/perks v1.0.1
Quantile estimation algorithms for Prometheus summary metrics.

#### github.com/cespare/xxhash/v2 v2.2.0
Fast non-cryptographic hash function for metric label hashing.

#### github.com/prometheus/client_model v0.6.0
Prometheus data model defining Counter, Gauge, Histogram, and Summary types.

#### github.com/prometheus/common v0.50.0
Common utilities shared across Prometheus client libraries.

#### github.com/prometheus/procfs v0.13.0
Parses Linux `/proc` filesystem to collect system metrics (CPU, memory, disk).

#### google.golang.org/protobuf v1.33.0
Protocol Buffers for efficient binary serialization (Prometheus exposition format).

---

### Logging Dependencies

#### github.com/mattn/go-colorable v0.1.13
Enables ANSI colored output in Windows terminals for zerolog.

#### github.com/mattn/go-isatty v0.0.20
Detects if output is a TTY to enable/disable colored console logs.

---

### System Dependencies

#### golang.org/x/sync v0.6.0
Advanced synchronization primitives including `errgroup` for coordinating goroutines and error handling.

#### golang.org/x/sys v0.18.0
Low-level OS system calls for cross-platform system access.

#### golang.org/x/text v0.14.0
Text processing, Unicode handling, encoding, and internationalization support.

---

## Dependency Size Summary

| Category | Libraries | Total Size |
|----------|-----------|------------|
| **Blockchain** | go-ethereum + crypto libs | ~50MB |
| **Database** | pgx + pooling | ~2MB |
| **Messaging** | nats.go + nkeys | ~5MB |
| **Config** | koanf + parsers | ~1MB |
| **Logging** | zerolog | ~500KB |
| **Metrics** | prometheus | ~3MB |
| **Storage** | bbolt | ~100KB |
| **System** | golang.org/x/* | ~2MB |
| **Total Binary** | | **~63MB** |

*Note: Actual binary size will be smaller due to dead code elimination during compilation.*

---

## Why These Specific Choices?

### Battle-Tested in Production

All dependencies are used in the reference projects:

1. **go-ethereum** - Only official Ethereum library, 8+ years in production
2. **pgx** - 5-10x faster than database/sql, native TimescaleDB support
3. **koanf** - Used by eth-tracker, simple and flexible
4. **nats.go** - Industry leader for event streaming, JetStream is production-grade
5. **zerolog** - Zero-allocation means zero GC pressure at high throughput
6. **prometheus** - Industry standard for observability
7. **bbolt** - Used by etcd, Kubernetes, proven reliability

### Performance Optimized

- **pgx**: Binary protocol, connection pooling, batch operations
- **zerolog**: Zero heap allocations, 10x faster than stdlib
- **holiman/uint256**: Optimized for Ethereum's 256-bit math
- **bbolt**: Memory-mapped I/O for fast reads

### Production Features

- **Deduplication**: NATS JetStream prevents duplicate event processing
- **Connection pooling**: pgx automatically manages database connections
- **Health checks**: pgx pool, NATS reconnection logic
- **Observability**: Prometheus metrics, structured logging
- **Crash safety**: BoltDB ACID transactions, TimescaleDB hypertables

---

## Installation

All dependencies are already specified in `go.mod`. To download:

```bash
go mod download
go mod tidy
```

Generate vendor directory (optional):
```bash
go mod vendor
```

---

## Upgrading Dependencies

Check for updates:
```bash
go list -u -m all
```

Upgrade all:
```bash
go get -u ./...
go mod tidy
```

Upgrade specific:
```bash
go get github.com/ethereum/go-ethereum@latest
go mod tidy
```

---

## Security Considerations

### Regular Updates

- **go-ethereum**: Security patches for consensus bugs
- **pgx**: SQL injection prevention, connection security
- **golang.org/x/crypto**: Cryptographic vulnerabilities

### Dependency Scanning

Run vulnerability checks:
```bash
# Install govulncheck
go install golang.org/x/vuln/cmd/govulncheck@latest

# Scan dependencies
govulncheck ./...
```

### Minimal Dependencies

This project has **only 7 direct dependencies** (excluding stdlib), minimizing attack surface.

---

## Common Issues & Solutions

### 1. go-ethereum compilation issues

**Problem:** CGO errors or secp256k1 build failures

**Solution:**
```bash
# Ensure C compiler is installed
# macOS
xcode-select --install

# Ubuntu/Debian
sudo apt-get install build-essential

# Set CGO_ENABLED if needed
export CGO_ENABLED=1
go build
```

### 2. pgx connection pool exhaustion

**Problem:** "too many connections" error

**Solution:** Increase pool size in config.toml:
```toml
[timescaledb]
max_connections = 50
max_idle_connections = 10
```

### 3. NATS connection failures

**Problem:** "connection refused" or timeout

**Solution:** Check NATS is running:
```bash
docker-compose up -d nats
docker-compose logs nats
```

---

## Further Reading

- [go-ethereum Documentation](https://geth.ethereum.org/docs/developers/dapp-developer/native)
- [pgx Documentation](https://github.com/jackc/pgx)
- [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/instrumentation/)
- [zerolog Documentation](https://github.com/rs/zerolog)

---

**Last Updated:** December 31, 2025  
**Go Version:** 1.21+  
**Total Dependencies:** 38 (7 direct, 31 indirect)
