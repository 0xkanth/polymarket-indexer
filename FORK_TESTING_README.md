# Fork Testing Quick Reference

## Installation
```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

## Start Fork

### Daemon Mode (Recommended - Survives Terminal Closure)
```bash
# Start in background
./scripts/start-fork.sh 55000000 --daemon

# Check status
./scripts/check-fork.sh

# View logs
tail -f anvil.log

# Stop
./scripts/stop-fork.sh
```

### Foreground Mode (Old Way)
```bash
# Terminal 1: Start fork at specific block
./scripts/start-fork.sh 55000000

# Or manually
anvil --fork-url https://polygon-rpc.com --fork-block-number 55000000 --chain-id 137
```

## Run Tests
```bash
# Terminal 2: Run tests
go test ./test -v

# Specific test
go test ./test -v -run TestForkRead

# Skip fork tests (CI)
go test ./test -short
```

## Default Anvil Accounts
```
Account 0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Private:   0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

Account 1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
Private:   0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
```

## Useful RPC Methods (from Go)
```go
client, _ := ethclient.Dial("http://127.0.0.1:8545")

// Impersonate any address
client.Client().Call(nil, "anvil_impersonateAccount", "0xAddress")

// Set balance
client.Client().Call(nil, "anvil_setBalance", "0xAddress", "0xDE0B6B3A7640000") // 1 ETH

// Time travel
client.Client().Call(nil, "evm_increaseTime", 3600) // +1 hour
client.Client().Call(nil, "evm_mine")              // Mine block

// Snapshot
var snapshotID string
client.Client().Call(&snapshotID, "evm_snapshot")
client.Client().Call(nil, "evm_revert", snapshotID)
```

## Test Pattern
```go
func TestFork(t *testing.T) {
    // 1. Load config
    cfg, _ := config.LoadConfig("../config/chains.json")
    chainCfg, _ := cfg.GetChain("polygon-fork")
    
    // 2. Create service
    svc, _ := service.NewCTFService(context.Background(), chainCfg)
    defer svc.Close()
    
    // 3. Test read
    balance, _ := svc.BalanceOf(ctx, addr, positionId)
    
    // 4. Test write
    auth := getAuth() // Use Anvil account
    tx, _ := svc.FillOrder(ctx, auth, order, amount, sig)
    receipt, _ := svc.WaitForTransaction(ctx, tx)
}
```

## Troubleshooting
```bash
# Check if Anvil is running
curl -X POST http://127.0.0.1:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Should return block number in hex
```

## Management Scripts

- `./scripts/start-fork.sh [BLOCK] [--daemon]` - Start fork
- `./scripts/stop-fork.sh` - Stop daemon
- `./scripts/check-fork.sh` - Check status

## Files Created
- `pkg/config/config.go` - Chain configuration loader
- `pkg/service/ctf_service.go` - CTF service with read/write/listen
- `test/fork_test.go` - Fork test examples
- `scripts/start-fork.sh` - Fork startup script (daemon mode support)
- `scripts/stop-fork.sh` - Stop daemon script
- `scripts/check-fork.sh` - Status check script
- `docs/FORK_TESTING_GUIDE.md` - Complete guide
- `docs/DAEMON_MODE.md` - Daemon mode documentation
- `.github/workflows/fork-tests.yml` - CI/CD workflow
- `config/chains.json` - Chain configs (already created)
