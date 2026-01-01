# Production-Grade Transaction Helpers

## Overview

The CTF service now includes production-ready, reusable transaction helpers with:
- ✅ **Transaction simulation** via `eth_call` before sending
- ✅ **Gas estimation with buffer** to prevent out-of-gas failures
- ✅ **Exponential backoff retry** on RPC/network failures
- ✅ **Smart error detection** (retryable vs permanent)
- ✅ **Configurable timeouts** and retry limits

## Quick Start

### Simple Read (No Transaction)

```go
// No gas, no retry needed for reads
status, err := svc.GetOrderStatus(ctx, orderHash)
balance, err := svc.BalanceOf(ctx, address, positionId)
```

### High-Level Write (Recommended)

```go
// ExecuteTransaction does everything automatically:
// 1. Simulates via eth_call
// 2. Estimates gas with 20% buffer
// 3. Retries up to 3 times with exponential backoff
// 4. Returns transaction or error

tx, err := svc.ExecuteTransaction(ctx, msg, auth, func(opts *bind.TransactOpts) (*types.Transaction, error) {
    return svc.ctfExchange.FillOrder(opts, order, amount, signature)
})
```

## Detailed API

### 1. SimulateTransaction

Simulates a transaction using `eth_call` before sending. Catches reverts early.

```go
msg := ethereum.CallMsg{
    From:  auth.From,
    To:    &contractAddr,
    Value: big.NewInt(0),
    Data:  encodedData,
}

err := svc.SimulateTransaction(ctx, msg)
if err != nil {
    // Transaction would revert - don't send
    log.Printf("Simulation failed: %v", err)
    return
}

// Safe to proceed
```

**Use when:** You want to verify transaction won't revert before spending gas

### 2. EstimateGasWithBuffer

Estimates gas and adds a safety buffer percentage.

```go
// Estimate with 20% buffer
gasLimit, err := svc.EstimateGasWithBuffer(ctx, msg, 20)
auth.GasLimit = gasLimit
```

**Parameters:**
- `bufferPercent`: Extra % to add (typical: 10-30%)

**Returns:**
- Gas limit with buffer applied
- Capped at 30M maximum

### 3. SendTransactionWithRetry

Sends transaction with full retry logic and configuration.

```go
config := &service.TransactionConfig{
    MaxRetries:       5,                // Try up to 5 times
    InitialBackoff:   2 * time.Second,  // Start with 2s backoff
    MaxBackoff:       60 * time.Second, // Cap at 60s
    GasBufferPercent: 25,               // 25% gas buffer
    Simulate:         true,             // Simulate first
    TimeoutPerTry:    30 * time.Second, // 30s per attempt
}

tx, err := svc.SendTransactionWithRetry(
    ctx,
    msg,
    auth,
    config,
    func(opts *bind.TransactOpts) (*types.Transaction, error) {
        return svc.ctfExchange.FillOrder(opts, order, amount, sig)
    },
)
```

**Flow:**
1. Simulates (if enabled)
2. Estimates gas with buffer
3. Tries to send transaction
4. On RPC/network error: waits (exponential backoff) and retries
5. On permanent error: returns immediately
6. Returns after success or max retries

### 4. ExecuteTransaction (Recommended)

High-level helper with sensible defaults.

```go
tx, err := svc.ExecuteTransaction(ctx, msg, auth, sendFunc)
```

**Defaults:**
- MaxRetries: 3
- InitialBackoff: 1s
- MaxBackoff: 30s
- GasBuffer: 20%
- Simulate: true
- Timeout: 30s per try

## Error Handling

### Retryable Errors (Will Retry)

- Connection refused/reset
- Timeouts
- Rate limiting (429)
- Gateway errors (502, 503, 504)
- Network unreachable

### Permanent Errors (Won't Retry)

- `execution reverted` - Transaction logic failed
- `insufficient funds` - Not enough ETH
- `gas too low` - Gas price too low
- `nonce too low` - Nonce already used
- `already known` - Transaction already in mempool

## Usage Patterns

### Pattern 1: Simple Transaction

```go
// Just send, let ExecuteTransaction handle everything
tx, err := svc.ExecuteTransaction(ctx, msg, auth, func(opts *bind.TransactOpts) (*types.Transaction, error) {
    return contract.Transfer(opts, to, amount)
})
```

### Pattern 2: Custom Retry Configuration

```go
config := &service.TransactionConfig{
    MaxRetries:       10,              // High retry for important txs
    InitialBackoff:   500 * time.Millisecond,
    MaxBackoff:       30 * time.Second,
    GasBufferPercent: 30,              // Higher buffer for complex txs
    Simulate:         true,
    TimeoutPerTry:    60 * time.Second,
}

tx, err := svc.SendTransactionWithRetry(ctx, msg, auth, config, sendFunc)
```

### Pattern 3: Manual Control

```go
// 1. Simulate first
if err := svc.SimulateTransaction(ctx, msg); err != nil {
    return fmt.Errorf("would revert: %w", err)
}

// 2. Estimate gas
gasLimit, err := svc.EstimateGasWithBuffer(ctx, msg, 20)
auth.GasLimit = gasLimit

// 3. Send with manual retry logic
for attempt := 0; attempt < 3; attempt++ {
    tx, err := contract.DoSomething(auth, params)
    if err == nil {
        break
    }
    time.Sleep(time.Second * time.Duration(attempt+1))
}
```

### Pattern 4: Batch Transactions

```go
// Prepare multiple transactions
transactions := []struct{
    msg ethereum.CallMsg
    sendFunc func(*bind.TransactOpts) (*types.Transaction, error)
}{
    {msg1, sendFunc1},
    {msg2, sendFunc2},
    {msg3, sendFunc3},
}

// Execute each with retry
for i, txData := range transactions {
    tx, err := svc.ExecuteTransaction(ctx, txData.msg, auth, txData.sendFunc)
    if err != nil {
        log.Printf("Transaction %d failed: %v", i, err)
        continue
    }
    log.Printf("Transaction %d sent: %s", i, tx.Hash().Hex())
}
```

## Configuration Recommendations

### Development/Testing
```go
config := &service.TransactionConfig{
    MaxRetries:       2,
    InitialBackoff:   500 * time.Millisecond,
    MaxBackoff:       5 * time.Second,
    GasBufferPercent: 50, // Higher buffer for safety
    Simulate:         true,
    TimeoutPerTry:    10 * time.Second,
}
```

### Production (Standard)
```go
config := service.DefaultTransactionConfig() // Use defaults
```

### Production (Mission Critical)
```go
config := &service.TransactionConfig{
    MaxRetries:       10,              // More retries
    InitialBackoff:   1 * time.Second,
    MaxBackoff:       60 * time.Second,
    GasBufferPercent: 25,              // Standard buffer
    Simulate:         true,            // Always simulate
    TimeoutPerTry:    60 * time.Second, // Longer timeout
}
```

### Production (Fast Operations)
```go
config := &service.TransactionConfig{
    MaxRetries:       3,
    InitialBackoff:   500 * time.Millisecond,
    MaxBackoff:       10 * time.Second, // Fail fast
    GasBufferPercent: 15,               // Lower buffer
    Simulate:         true,
    TimeoutPerTry:    15 * time.Second,
}
```

## Real-World Example

```go
func FillOrderWithRetry(
    svc *service.CTFService,
    ctx context.Context,
    auth *bind.TransactOpts,
    order Order,
    amount *big.Int,
    signature []byte,
) (*types.Receipt, error) {
    // 1. Prepare transaction message
    msg := ethereum.CallMsg{
        From:  auth.From,
        To:    &contractAddress,
        Value: big.NewInt(0),
        // Data: encoded function call (would need ABI encoder)
    }

    // 2. Execute with defaults
    tx, err := svc.ExecuteTransaction(ctx, msg, auth, func(opts *bind.TransactOpts) (*types.Transaction, error) {
        return svc.FillOrder(ctx, opts, order, amount, signature)
    })
    if err != nil {
        return nil, fmt.Errorf("failed to send transaction: %w", err)
    }

    log.Printf("Transaction sent: %s", tx.Hash().Hex())

    // 3. Wait for confirmation
    receipt, err := svc.WaitForTransaction(ctx, tx)
    if err != nil {
        return nil, fmt.Errorf("transaction failed: %w", err)
    }

    log.Printf("Transaction confirmed in block %d", receipt.BlockNumber.Uint64())
    return receipt, nil
}
```

## Best Practices

1. **Always use ExecuteTransaction for production writes**
   - Handles simulation, gas, retry automatically
   - Sensible defaults for most use cases

2. **Increase retries for important transactions**
   - Financial operations: 5-10 retries
   - User-facing operations: 3-5 retries
   - Background tasks: 2-3 retries

3. **Adjust gas buffer based on complexity**
   - Simple transfers: 10-15%
   - Standard calls: 20%
   - Complex multi-call: 30-40%

4. **Always simulate in production**
   - Catches logic errors before gas spent
   - Only disable for read-only calls

5. **Set appropriate timeouts**
   - Polygon: 5-15s per try (fast blocks)
   - Ethereum: 30-60s per try (slower blocks)
   - During congestion: increase timeout

6. **Monitor and log**
   - Log all retry attempts
   - Track success rates
   - Alert on consistent failures

## Migration Guide

### Old Code
```go
auth.GasLimit = 500000 // Manual gas
tx, err := contract.DoSomething(auth, params)
// No retry, no simulation
```

### New Code (Recommended)
```go
msg := ethereum.CallMsg{From: auth.From, To: &addr}
tx, err := svc.ExecuteTransaction(ctx, msg, auth, func(opts *bind.TransactOpts) (*types.Transaction, error) {
    return contract.DoSomething(opts, params)
})
```

## Summary

The production-grade helpers provide:

| Feature | Benefit |
|---------|---------|
| Simulation | Catch reverts before spending gas |
| Gas estimation | Prevent out-of-gas failures |
| Buffer | Handle gas usage variations |
| Retry logic | Survive RPC outages |
| Exponential backoff | Reduce load during failures |
| Error detection | Don't retry permanent failures |
| Configurability | Tune for your use case |

Use `ExecuteTransaction` for 99% of cases - it just works! 🚀
