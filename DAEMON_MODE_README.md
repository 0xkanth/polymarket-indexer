# Daemon Mode - Quick Start

## ✅ Done! You now have daemon mode support

## What Changed?

1. **[scripts/start-fork.sh](scripts/start-fork.sh)** - Now supports `--daemon` flag
2. **[scripts/stop-fork.sh](scripts/stop-fork.sh)** - Stop the daemon gracefully  
3. **[scripts/check-fork.sh](scripts/check-fork.sh)** - Check health and status
4. **[.github/workflows/fork-tests.yml](.github/workflows/fork-tests.yml)** - CI/CD workflow example
5. **[docs/DAEMON_MODE.md](docs/DAEMON_MODE.md)** - Complete daemon mode guide

## Usage

### Start in Daemon Mode

```bash
# Recommended: Runs in background, survives terminal closure
./scripts/start-fork.sh 55000000 --daemon

# Alternative with env var
DAEMON=1 ./scripts/start-fork.sh 55000000

# Auto-detects in CI/CD (when CI=true)
```

### Manage the Daemon

```bash
# Check if running
./scripts/check-fork.sh

# View live logs
tail -f anvil.log

# Stop daemon
./scripts/stop-fork.sh
```

### Run Tests

```bash
# Start daemon once
./scripts/start-fork.sh 55000000 --daemon

# Run tests (daemon keeps running)
go test ./test -v

# Run specific test
go test ./test -v -run TestForkRead

# Stop when done
./scripts/stop-fork.sh
```

## Benefits

✅ **Terminal Resilient** - Won't abort if terminal closes  
✅ **CI/CD Ready** - Auto-detects CI environment  
✅ **Process Management** - PID file tracking  
✅ **Health Monitoring** - Built-in status checks  
✅ **Log Persistence** - All output saved to `anvil.log`  

## CI/CD Integration

### GitHub Actions

```yaml
- name: Start Anvil Fork
  run: DAEMON=1 ./scripts/start-fork.sh 55000000

- name: Run Tests  
  run: go test ./test -v

- name: Stop Anvil
  if: always()
  run: ./scripts/stop-fork.sh
```

See [.github/workflows/fork-tests.yml](.github/workflows/fork-tests.yml) for complete example.

## Files Overview

```
scripts/
  start-fork.sh    # Start fork (foreground or daemon)
  stop-fork.sh     # Stop daemon
  check-fork.sh    # Health check

.anvil.pid         # Process ID (auto-created)
anvil.log          # Output logs (auto-created)

docs/
  DAEMON_MODE.md         # Complete guide
  FORK_TESTING_GUIDE.md  # Original fork testing guide
  
.github/workflows/
  fork-tests.yml   # CI/CD workflow
```

## Troubleshooting

### Port Already in Use

If port 8545 is occupied (e.g., by Docker):

```bash
# Find what's using it
lsof -i :8545

# Kill it
kill $(lsof -t -i:8545)

# Or use different port
PORT=8546 DAEMON=1 ./scripts/start-fork.sh 55000000
```

### Stale PID File

```bash
# Clean up
rm -f .anvil.pid

# Restart
./scripts/start-fork.sh 55000000 --daemon
```

### Check Logs

```bash
# Real-time
tail -f anvil.log

# Full log
cat anvil.log
```

## Next Steps

1. **Read** [docs/DAEMON_MODE.md](docs/DAEMON_MODE.md) for complete guide
2. **Test** locally: `./scripts/start-fork.sh 55000000 --daemon`
3. **Run tests**: `go test ./test -v`
4. **Setup CI/CD** using [.github/workflows/fork-tests.yml](.github/workflows/fork-tests.yml)

Now your fork testing is production-ready! 🚀
