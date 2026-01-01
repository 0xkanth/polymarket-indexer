# Testing the Polymarket Indexer

This guide covers testing strategies, test execution, and verification procedures.

## Test Structure

```
polymarket-indexer/
├── internal/
│   ├── handler/
│   │   ├── events.go
│   │   └── events_test.go          # Handler unit tests
│   ├── router/
│   │   ├── router.go
│   │   └── router_test.go          # Router unit tests
│   ├── chain/
│   │   ├── client.go
│   │   └── client_test.go          # Chain client tests
│   ├── db/
│   │   ├── checkpoint.go
│   │   └── checkpoint_test.go      # BoltDB tests
│   └── processor/
│       ├── processor.go
│       └── processor_test.go       # Processor integration tests
└── test/
    ├── integration/                # Integration tests
    ├── e2e/                        # End-to-end tests
    └── fixtures/                   # Test data
```

## Quick Start

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/handler -v

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

## Unit Tests

### Testing Event Handlers

Create `internal/handler/events_test.go`:

```go
package handler

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/0xkanth/polymarket-indexer/pkg/models"
)

func TestHandleOrderFilled(t *testing.T) {
	tests := []struct {
		name      string
		log       types.Log
		timestamp uint64
		want      models.OrderFilled
		wantErr   bool
	}{
		{
			name: "valid order filled event",
			log: types.Log{
				Address: common.HexToAddress("0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E"),
				Topics: []common.Hash{
					OrderFilledSig,
					common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"), // orderHash
					common.HexToAddress("0x1111111111111111111111111111111111111111").Hash(),                   // maker
					common.HexToAddress("0x2222222222222222222222222222222222222222").Hash(),                   // taker
				},
				Data: hexToBytes(
					"0000000000000000000000000000000000000000000000000000000000000001" + // makerAssetId
						"0000000000000000000000000000000000000000000000000000000000000002" + // takerAssetId
						"00000000000000000000000000000000000000000000000000000000000003e8" + // makerAmountFilled (1000)
						"00000000000000000000000000000000000000000000000000000000000007d0" + // takerAmountFilled (2000)
						"0000000000000000000000000000000000000000000000000000000000000064", // fee (100)
				),
			},
			timestamp: 1234567890,
			want: models.OrderFilled{
				OrderHash:         "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Maker:             "0x1111111111111111111111111111111111111111",
				Taker:             "0x2222222222222222222222222222222222222222",
				MakerAssetID:      big.NewInt(1),
				TakerAssetID:      big.NewInt(2),
				MakerAmountFilled: big.NewInt(1000),
				TakerAmountFilled: big.NewInt(2000),
				Fee:               big.NewInt(100),
			},
			wantErr: false,
		},
		{
			name: "invalid topic count",
			log: types.Log{
				Topics: []common.Hash{OrderFilledSig},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HandleOrderFilled(context.Background(), tt.log, tt.timestamp)
			
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func hexToBytes(s string) []byte {
	b, _ := common.FromHex(s)
	return b
}
```

### Testing Router

Create `internal/router/router_test.go`:

```go
package router

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	logger := zerolog.Nop()
	router := New(logger)

	// Mock handler
	mockSig := common.HexToHash("0x1234")
	mockPayload := struct{ Value string }{Value: "test"}
	mockHandler := func(ctx context.Context, log types.Log, ts uint64) (any, error) {
		return mockPayload, nil
	}

	// Register handler
	router.RegisterLogHandler(mockSig, mockHandler)

	// Test routing
	log := types.Log{
		Topics: []common.Hash{mockSig},
	}

	payload, err := router.RouteLog(context.Background(), log, 1234567890)
	require.NoError(t, err)
	assert.Equal(t, mockPayload, payload)

	// Test unknown event
	unknownLog := types.Log{
		Topics: []common.Hash{common.HexToHash("0xunknown")},
	}
	_, err = router.RouteLog(context.Background(), unknownLog, 1234567890)
	assert.ErrorIs(t, err, ErrNoHandler)
}
```

### Testing Checkpoint Store

Create `internal/db/checkpoint_test.go`:

```go
package db

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckpointStore(t *testing.T) {
	// Create temp file
	tmpFile := "/tmp/test_checkpoint.db"
	defer os.Remove(tmpFile)

	store, err := NewCheckpointStore(tmpFile)
	require.NoError(t, err)
	defer store.Close()

	ctx := context.Background()

	t.Run("GetOrCreateCheckpoint creates new", func(t *testing.T) {
		cp, err := store.GetOrCreateCheckpoint(ctx, "test-service", 1000)
		require.NoError(t, err)
		assert.Equal(t, uint64(1000), cp.LastBlock)
		assert.Equal(t, "test-service", cp.ServiceName)
	})

	t.Run("UpdateBlock updates checkpoint", func(t *testing.T) {
		err := store.UpdateBlock(ctx, "test-service", 2000, "0xabcd")
		require.NoError(t, err)

		cp, err := store.GetCheckpoint(ctx, "test-service")
		require.NoError(t, err)
		assert.Equal(t, uint64(2000), cp.LastBlock)
		assert.Equal(t, "0xabcd", cp.LastBlockHash)
	})

	t.Run("SaveCheckpoint overwrites", func(t *testing.T) {
		newCp := Checkpoint{
			ServiceName:   "test-service",
			LastBlock:     3000,
			LastBlockHash: "0xbeef",
			UpdatedAt:     time.Now(),
		}
		err := store.SaveCheckpoint(ctx, newCp)
		require.NoError(t, err)

		cp, err := store.GetCheckpoint(ctx, "test-service")
		require.NoError(t, err)
		assert.Equal(t, uint64(3000), cp.LastBlock)
	})
}
```

## Integration Tests

### Testing with Mock RPC

Create `test/integration/processor_test.go`:

```go
package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	
	"github.com/0xkanth/polymarket-indexer/internal/chain"
	"github.com/0xkanth/polymarket-indexer/internal/processor"
	// ... other imports
)

func TestProcessorIntegration(t *testing.T) {
	// Mock RPC server
	mockRPC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string `json:"method"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		switch req.Method {
		case "eth_blockNumber":
			json.NewEncoder(w).Encode(map[string]any{
				"result": "0x13a1c43", // 20,558,323
			})
		case "eth_getBlockByNumber":
			json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{
					"number":    "0x13a1c43",
					"hash":      "0xabcd1234",
					"timestamp": "0x65abcd12",
				},
			})
		case "eth_getLogs":
			json.NewEncoder(w).Encode(map[string]any{
				"result": []any{}, // No logs for test
			})
		}
	}))
	defer mockRPC.Close()

	// Create chain client with mock RPC
	chainClient, err := chain.New(chain.Config{
		HTTPURL: mockRPC.URL,
		ChainID: 137,
	})
	require.NoError(t, err)

	// Test processing
	// ... rest of test
}
```

### Testing NATS Integration

```go
func TestNATSIntegration(t *testing.T) {
	// Start embedded NATS server
	srv, err := server.NewServer(&server.Options{
		JetStream: true,
		Port:      -1, // random port
	})
	require.NoError(t, err)
	defer srv.Shutdown()

	go srv.Start()
	if !srv.ReadyForConnections(5 * time.Second) {
		t.Fatal("nats server not ready")
	}

	// Connect and test publishing
	publisher, err := nats.NewPublisher(nats.PublisherConfig{
		URL:        srv.ClientURL(),
		StreamName: "TEST",
		Subjects:   []string{"TEST.>"},
	})
	require.NoError(t, err)
	defer publisher.Close()

	// Publish test event
	event := models.Event{
		BlockNumber:     1000,
		TransactionHash: "0xtest",
		Payload:         models.OrderFilled{},
	}

	err = publisher.Publish(context.Background(), "OrderFilled", "0xcontract", event)
	require.NoError(t, err)
}
```

### Testing Database Operations

```go
func TestDatabaseIntegration(t *testing.T) {
	// Use testcontainers for real PostgreSQL
	ctx := context.Background()
	
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "timescale/timescaledb:latest-pg15",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "test",
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	})
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	// Get connection details
	host, _ := postgresContainer.Host(ctx)
	port, _ := postgresContainer.MappedPort(ctx, "5432")

	// Run migrations
	connStr := fmt.Sprintf("postgresql://test:test@%s:%s/test?sslmode=disable", host, port.Port())
	// ... run migrations

	// Test database operations
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	// Insert test data
	_, err = pool.Exec(ctx, `
		INSERT INTO events (block_number, block_hash, block_timestamp, transaction_hash, log_index, contract_address, event_signature, payload)
		VALUES (1000, '0xhash', to_timestamp(1234567890), '0xtx', 0, '0xcontract', '0xsig', '{}')
	`)
	require.NoError(t, err)

	// Query test data
	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM events").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
```

## End-to-End Tests

### Local Chain E2E Test

```go
func TestE2EWithLocalChain(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	// Start local Ethereum node (anvil, ganache, etc.)
	// Deploy contracts
	// Emit test events
	// Run indexer
	// Verify data in database
}
```

### Manual E2E Testing

```bash
# 1. Start infrastructure
docker-compose up -d nats timescaledb

# 2. Run migrations
make migrate-up

# 3. Start indexer (in one terminal)
go run ./cmd/indexer

# 4. Start consumer (in another terminal)
go run ./cmd/consumer

# 5. Verify events are being indexed
# Watch indexer logs
docker-compose logs -f indexer

# Check database
psql -h localhost -U polymarket -d polymarket -c "SELECT COUNT(*) FROM events;"

# 6. Check metrics
curl http://localhost:9090/metrics | grep polymarket

# 7. Verify health
curl http://localhost:8080/health
```

## Test Data Generation

### Creating Test Fixtures

```go
// test/fixtures/events.go
package fixtures

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func OrderFilledLog() types.Log {
	return types.Log{
		Address: common.HexToAddress("0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E"),
		Topics: []common.Hash{
			common.HexToHash("0xd0a08e8c493f9c94f29311604c9de0fa40fe441d0d4d6e8b87b3e1a4cbadba5c"),
			common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
			common.HexToAddress("0x1111111111111111111111111111111111111111").Hash(),
			common.HexToAddress("0x2222222222222222222222222222222222222222").Hash(),
		},
		Data:        makeOrderFilledData(1, 2, 1000, 2000, 100),
		BlockNumber: 1000,
		TxHash:      common.HexToHash("0xtxhash"),
		Index:       0,
	}
}

func makeOrderFilledData(makerAsset, takerAsset, makerAmount, takerAmount, fee int64) []byte {
	data := make([]byte, 160)
	// Encode fields...
	return data
}
```

## Performance Testing

### Benchmark Tests

```go
func BenchmarkHandleOrderFilled(b *testing.B) {
	log := fixtures.OrderFilledLog()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HandleOrderFilled(ctx, log, 1234567890)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessBlock(b *testing.B) {
	// Setup processor with mock dependencies
	// ...

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := processor.ProcessBlock(context.Background(), uint64(1000+i))
		if err != nil {
			b.Fatal(err)
		}
	}
}
```

### Load Testing

```bash
# Generate high-volume test data
go run test/loadgen/main.go --events 10000

# Run indexer with profiling
go run -pprof=:6060 ./cmd/indexer

# Analyze performance
go tool pprof http://localhost:6060/debug/pprof/profile
```

## Continuous Integration

### GitHub Actions Workflow

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: timescale/timescaledb:latest-pg15
        env:
          POSTGRES_PASSWORD: test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
      
      nats:
        image: nats:latest
        options: --name nats-test

    steps:
    - uses: actions/checkout@v3
    
    - uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Run tests
      run: |
        go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
    
    - name: Upload coverage
      uses: codecov/codecov-action@v3
      with:
        files: ./coverage.txt
```

## Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# Open in browser
open coverage.html
```

### Coverage Goals

- **Overall**: >80%
- **Critical paths**: >90%
  - Event handlers
  - Checkpoint operations
  - NATS publishing
  - Database inserts

## Testing Checklist

Before deploying:

- [ ] All unit tests passing
- [ ] Integration tests passing
- [ ] E2E test with local chain
- [ ] Load test with 10k+ events
- [ ] Reorg scenario tested
- [ ] Error handling tested (RPC failures, DB failures)
- [ ] Metrics verified (Prometheus)
- [ ] Health checks working
- [ ] Graceful shutdown tested
- [ ] Memory leak test (run for 24h)
- [ ] Database migration tested (up and down)
- [ ] Configuration validation tested

## Common Test Patterns

### Table-Driven Tests

```go
tests := []struct {
    name    string
    input   Input
    want    Output
    wantErr bool
}{
    {"case 1", input1, output1, false},
    {"case 2", input2, output2, true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := Function(tt.input)
        if tt.wantErr {
            require.Error(t, err)
            return
        }
        require.NoError(t, err)
        assert.Equal(t, tt.want, got)
    })
}
```

### Test Helpers

```go
func setupTest(t *testing.T) (*TestEnv, func()) {
    t.Helper()
    
    env := &TestEnv{
        // Setup test environment
    }
    
    cleanup := func() {
        // Cleanup
    }
    
    return env, cleanup
}

func TestSomething(t *testing.T) {
    env, cleanup := setupTest(t)
    defer cleanup()
    
    // Test code
}
```

## Resources

- [Go Testing Documentation](https://pkg.go.dev/testing)
- [Testify](https://github.com/stretchr/testify) - Assertion library
- [Testcontainers](https://golang.testcontainers.org/) - Container-based integration tests
- [go-ethereum Testing](https://geth.ethereum.org/docs/developers/dapp-developer/native-bindings)
