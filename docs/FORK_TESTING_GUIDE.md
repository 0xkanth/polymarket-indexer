# Fork Testing Guide for Go Ethereum Services

## Overview

Fork testing allows you to test your Go services against a **local copy of Polygon mainnet** at a specific block number. This is similar to Hardhat's fork mode but uses **Anvil** (from Foundry).

## Why Fork Testing?

- ✅ Test against **real mainnet state** (actual balances, contract storage)
- ✅ Test with **real user addresses** without needing private keys
- ✅ **No gas costs** - it's a local chain
- ✅ **Deterministic** - fork at specific block for reproducible tests
- ✅ **Fast** - no network latency
- ✅ **Safe** - can't accidentally send real transactions

## Installation

### Install Foundry (includes Anvil)

```bash
# Install Foundry
curl -L https://foundry.paradigm.xyz | bash

# Reload shell or run:
source ~/.bashrc  # or ~/.zshrc

# Install/update Foundry tools
foundryup

# Verify installation
anvil --version
```

## Starting a Forked Network

### Basic Fork (Latest Block)

```bash
anvil --fork-url https://polygon-rpc.com
```

### Fork at Specific Block Number

```bash
# Fork Polygon at block 55,000,000
anvil --fork-url https://polygon-rpc.com --fork-block-number 55000000 --chain-id 137
```

### Using the Script

```bash
# Make script executable
chmod +x scripts/start-fork.sh

# Fork at specific block
./scripts/start-fork.sh 55000000

# Or use default block
./scripts/start-fork.sh
```

### Anvil Configuration Options

```bash
anvil \
  --fork-url https://polygon-rpc.com \        # RPC to fork from
  --fork-block-number 55000000 \              # Specific block
  --chain-id 137 \                            # Polygon mainnet chain ID
  --port 8545 \                               # Local port
  --host 127.0.0.1 \                         # Local host
  --block-time 2 \                            # 2 second blocks (like Polygon)
  --accounts 10 \                             # 10 test accounts
  --balance 10000 \                           # 10000 ETH per account
  --gas-limit 30000000                        # Block gas limit
```

## Default Test Accounts

Anvil provides 10 pre-funded accounts:

```
Account #0: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
Private Key: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

Account #1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8
Private Key: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

... (8 more accounts)
```

Each has 10,000 ETH by default.

## Running Go Tests Against Fork

### 1. Start Fork in Terminal 1

```bash
./scripts/start-fork.sh 55000000
```

Keep this running!

### 2. Run Tests in Terminal 2

```bash
# Run all fork tests
go test ./test -v

# Run specific test
go test ./test -v -run TestForkRead

# Skip fork tests (for CI/CD)
go test ./test -short
```

## Common Fork Testing Patterns

### Pattern 1: Reading Real Mainnet State

```go
func TestForkReadRealData(t *testing.T) {
    ctx := context.Background()
    
    // Connect to fork
    cfg, _ := config.LoadConfig("../config/chains.json")
    chainCfg, _ := cfg.GetChain("polygon-fork")
    svc, _ := service.NewCTFService(ctx, chainCfg)
    defer svc.Close()
    
    // Read actual balance from a known Polygon address
    realUserAddr := common.HexToAddress("0xABCD...") // Real Polygon address
    positionId := big.NewInt(12345)
    
    balance, err := svc.BalanceOf(ctx, realUserAddr, positionId)
    require.NoError(t, err)
    
    t.Logf("User %s has balance: %s", realUserAddr.Hex(), balance.String())
}
```

### Pattern 2: Impersonating Any Address

Anvil lets you send transactions **as any address** without the private key!

```go
func TestForkImpersonate(t *testing.T) {
    ctx := context.Background()
    client, _ := ethclient.Dial("http://127.0.0.1:8545")
    
    // Impersonate a rich address from mainnet
    richAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb") // Example whale
    
    // Use anvil_impersonateAccount RPC call
    err := client.Client().Call(nil, "anvil_impersonateAccount", richAddr.Hex())
    require.NoError(t, err)
    
    // Now send transaction as that address!
    auth := &bind.TransactOpts{
        From:     richAddr,
        Signer:   nil, // No signature needed when impersonating
        GasLimit: 500000,
    }
    
    // Execute transaction...
}
```

### Pattern 3: Manipulating Fork State

```go
func TestForkStateManipulation(t *testing.T) {
    client, _ := ethclient.Dial("http://127.0.0.1:8545")
    
    // Set ETH balance for any address
    testAddr := common.HexToAddress("0x1234...")
    balance := "0xDE0B6B3A7640000" // 1 ETH in hex
    
    err := client.Client().Call(nil, "anvil_setBalance", testAddr.Hex(), balance)
    require.NoError(t, err)
    
    // Set storage slot (advanced)
    contractAddr := common.HexToAddress("0xCTFExchange...")
    storageSlot := common.HexToHash("0x0")
    value := common.HexToHash("0x1")
    
    err = client.Client().Call(nil, "anvil_setStorageAt", 
        contractAddr.Hex(), 
        storageSlot.Hex(), 
        value.Hex())
    require.NoError(t, err)
}
```

### Pattern 4: Time Travel

```go
func TestForkTimeTravel(t *testing.T) {
    client, _ := ethclient.Dial("http://127.0.0.1:8545")
    
    // Increase time by 1 hour
    seconds := 3600
    err := client.Client().Call(nil, "evm_increaseTime", seconds)
    require.NoError(t, err)
    
    // Mine a new block to apply the time change
    err = client.Client().Call(nil, "evm_mine")
    require.NoError(t, err)
}
```

### Pattern 5: Testing Transaction Reverts

```go
func TestForkRevertScenario(t *testing.T) {
    ctx := context.Background()
    svc, _ := service.NewCTFService(ctx, chainCfg)
    
    // Try to fill an order that should revert
    auth := getTestAuth()
    
    tx, err := svc.FillOrder(ctx, auth, invalidOrder, fillAmount, sig)
    
    // Transaction submission might succeed
    if err != nil {
        t.Logf("Transaction rejected: %v", err)
        return
    }
    
    // Wait for receipt - should show revert
    receipt, err := svc.WaitForTransaction(ctx, tx)
    require.Error(t, err, "Expected revert error")
    require.Equal(t, uint64(0), receipt.Status, "Status should be 0 for revert")
}
```

## Complete Test Example

```go
package test

import (
    "context"
    "math/big"
    "testing"
    
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/stretchr/testify/require"
    
    "github.com/yourusername/polymarket-indexer/pkg/config"
    "github.com/yourusername/polymarket-indexer/pkg/service"
)

func TestCompleteOrderFlow(t *testing.T) {
    // Setup
    ctx := context.Background()
    cfg, err := config.LoadConfig("../config/chains.json")
    require.NoError(t, err)
    
    chainCfg, err := cfg.GetChain("polygon-fork")
    require.NoError(t, err)
    
    svc, err := service.NewCTFService(ctx, chainCfg)
    require.NoError(t, err)
    defer svc.Close()
    
    // Use Anvil default account
    privateKey, _ := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
    auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(137))
    
    // Get initial balance
    userAddr := auth.From
    positionId := big.NewInt(1)
    balanceBefore, err := svc.BalanceOf(ctx, userAddr, positionId)
    require.NoError(t, err)
    
    t.Logf("Balance before: %s", balanceBefore.String())
    
    // Execute order (with valid test data)
    // ... fill order logic ...
    
    // Check balance after
    balanceAfter, err := svc.BalanceOf(ctx, userAddr, positionId)
    require.NoError(t, err)
    
    t.Logf("Balance after: %s", balanceAfter.String())
    require.NotEqual(t, balanceBefore, balanceAfter, "Balance should change")
}
```

## Useful Anvil RPC Methods

### Account Impersonation
- `anvil_impersonateAccount(address)` - Impersonate address
- `anvil_stopImpersonatingAccount(address)` - Stop impersonation

### State Manipulation
- `anvil_setBalance(address, balance)` - Set ETH balance
- `anvil_setCode(address, code)` - Set contract code
- `anvil_setStorageAt(address, slot, value)` - Set storage slot
- `anvil_setNonce(address, nonce)` - Set nonce

### Time Control
- `evm_increaseTime(seconds)` - Fast forward time
- `evm_setNextBlockTimestamp(timestamp)` - Set next block time
- `evm_mine()` - Mine a block immediately

### Snapshots
- `evm_snapshot()` - Create state snapshot (returns ID)
- `evm_revert(id)` - Revert to snapshot

## Best Practices

### 1. Use Specific Block Numbers

```bash
# Good - reproducible
./scripts/start-fork.sh 55000000

# Bad - changes every time
./scripts/start-fork.sh  # Uses latest
```

### 2. Create Helper Functions

```go
// test/helpers.go
func GetForkService(t *testing.T) *service.CTFService {
    ctx := context.Background()
    cfg, err := config.LoadConfig("../config/chains.json")
    require.NoError(t, err)
    
    chainCfg, err := cfg.GetChain("polygon-fork")
    require.NoError(t, err)
    
    svc, err := service.NewCTFService(ctx, chainCfg)
    require.NoError(t, err)
    
    return svc
}
```

### 3. Use Table-Driven Tests

```go
func TestMultipleAddresses(t *testing.T) {
    tests := []struct {
        name    string
        address string
        want    *big.Int
    }{
        {"Alice", "0x123...", big.NewInt(1000)},
        {"Bob", "0x456...", big.NewInt(2000)},
    }
    
    svc := GetForkService(t)
    defer svc.Close()
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            balance, err := svc.BalanceOf(ctx, common.HexToAddress(tt.address), posId)
            require.NoError(t, err)
            require.Equal(t, tt.want, balance)
        })
    }
}
```

### 4. Snapshot and Revert

```go
func TestWithSnapshot(t *testing.T) {
    client, _ := ethclient.Dial("http://127.0.0.1:8545")
    
    // Take snapshot
    var snapshotID string
    err := client.Client().Call(&snapshotID, "evm_snapshot")
    require.NoError(t, err)
    
    // Make changes
    // ... test logic ...
    
    // Revert to snapshot
    var reverted bool
    err = client.Client().Call(&reverted, "evm_revert", snapshotID)
    require.NoError(t, err)
    require.True(t, reverted)
}
```

## Troubleshooting

### Fork Not Working

```bash
# Check if Anvil is running
curl -X POST http://127.0.0.1:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'

# Should return: {"jsonrpc":"2.0","id":1,"result":"0x..."}
```

### Wrong Chain ID

Make sure both Anvil and your config use same chain ID:

```bash
# Anvil
anvil --chain-id 137

# Config (chains.json)
"polygon-fork": {
  "chainId": 137,
  ...
}
```

### RPC Rate Limiting

If you see "429 Too Many Requests" when forking:

```bash
# Use Alchemy or Infura with API key
anvil --fork-url https://polygon-mainnet.g.alchemy.com/v2/YOUR_API_KEY
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Fork Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Foundry
        uses: foundry-rs/foundry-toolchain@v1
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Start Anvil Fork
        run: |
          anvil --fork-url ${{ secrets.POLYGON_RPC }} \
            --fork-block-number 55000000 \
            --port 8545 &
          sleep 5
      
      - name: Run Fork Tests
        run: go test ./test -v
```

## Summary

Fork testing with Anvil gives you:
1. **Real mainnet state** at specific block
2. **Unlimited test ETH** for gas
3. **Address impersonation** without private keys
4. **State manipulation** for edge case testing
5. **Fast, deterministic, safe** testing environment

Start with the provided `start-fork.sh` script and `fork_test.go` examples, then expand based on your needs!
