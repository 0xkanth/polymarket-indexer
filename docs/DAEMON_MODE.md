# Daemon Mode Fork Testing

## Quick Start

### Start in Daemon Mode (Background)

```bash
# Method 1: Using flag
./scripts/start-fork.sh 55000000 --daemon

# Method 2: Using environment variable
DAEMON=1 ./scripts/start-fork.sh 55000000

# Automatically uses daemon mode in CI/CD (when CI=true)
```

### Manage the Daemon

```bash
# Check status
./scripts/check-fork.sh

# View logs
tail -f anvil.log

# Stop daemon
./scripts/stop-fork.sh
```

## Why Daemon Mode?

✅ **Terminal resilient** - Won't stop if terminal closes accidentally  
✅ **CI/CD friendly** - Runs in background for automated testing  
✅ **Process management** - PID file tracking and health checks  
✅ **Log persistence** - All output saved to anvil.log  
✅ **Easy control** - Simple start/stop/check scripts  

## Usage Examples

### Local Development

```bash
# Terminal 1: Start daemon
./scripts/start-fork.sh 55000000 --daemon

# Check it's running
./scripts/check-fork.sh

# Terminal 2: Run tests
go test ./test -v

# When done
./scripts/stop-fork.sh
```

### Development Workflow

```bash
# Start once in the morning
./scripts/start-fork.sh 55000000 --daemon

# Run tests throughout the day
go test ./test -v -run TestForkRead
go test ./test -v -run TestForkWrite

# Check if still running
./scripts/check-fork.sh

# View logs if issues
tail -f anvil.log

# Stop at end of day
./scripts/stop-fork.sh
```

### CI/CD Integration

#### GitHub Actions

See [.github/workflows/fork-tests.yml](.github/workflows/fork-tests.yml) for complete example.

```yaml
steps:
  - name: Install Foundry
    uses: foundry-rs/foundry-toolchain@v1
  
  - name: Start Anvil fork
    env:
      DAEMON: 1  # or CI=true (auto-detected)
    run: ./scripts/start-fork.sh 55000000
  
  - name: Run tests
    run: go test ./test -v
  
  - name: Stop Anvil
    if: always()
    run: ./scripts/stop-fork.sh
```

#### GitLab CI

```yaml
fork-tests:
  image: golang:1.21
  before_script:
    - curl -L https://foundry.paradigm.xyz | bash
    - source ~/.bashrc
    - foundryup
  script:
    - DAEMON=1 ./scripts/start-fork.sh 55000000
    - go test ./test -v
  after_script:
    - ./scripts/stop-fork.sh
```

#### Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('Setup') {
            steps {
                sh 'curl -L https://foundry.paradigm.xyz | bash'
                sh 'foundryup'
            }
        }
        stage('Fork Tests') {
            steps {
                sh 'DAEMON=1 ./scripts/start-fork.sh 55000000'
                sh 'go test ./test -v'
            }
        }
    }
    post {
        always {
            sh './scripts/stop-fork.sh'
            archiveArtifacts artifacts: 'anvil.log', allowEmptyArchive: true
        }
    }
}
```

## Script Details

### start-fork.sh

**Features:**
- Daemon mode via `--daemon` flag or `DAEMON=1` env var
- Auto-detects CI environment (`CI=true`)
- PID file management (`.anvil.pid`)
- Log file output (`anvil.log`)
- Prevents duplicate instances
- Health check after startup

**Exit Codes:**
- `0` - Started successfully
- `1` - Anvil not installed or startup failed
- `1` - Already running (duplicate prevention)

### stop-fork.sh

**Features:**
- Graceful shutdown (SIGTERM)
- Force kill after 5 seconds (SIGKILL)
- Cleans up PID file
- Handles stale PID files

**Exit Codes:**
- `0` - Stopped successfully or not running

### check-fork.sh

**Features:**
- Process status (PID check)
- CPU/Memory usage
- Port listening check
- RPC endpoint health check
- Current block number
- Chain ID verification
- Color-coded status summary

**Output Example:**
```
Checking Anvil fork status...

✓ Process Status: Running (PID: 12345)
  CPU/MEM/Uptime:  0.5  0.8 00:05:23

✓ Port 8545: Listening
  anvil    12345 user   4u  IPv4 0x... TCP 127.0.0.1:8545 (LISTEN)

✓ RPC Endpoint: Responding
  Current Block: 55000123 (hex: 0x347007b)
  Chain ID: 137

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Status: ✓ Anvil fork is running properly
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

## Files Created

- `.anvil.pid` - Process ID file (auto-managed)
- `anvil.log` - Output logs (rotates with each start)
- `nohup.out` - Fallback log file (if nohup is used)

## Troubleshooting

### Check if Running

```bash
./scripts/check-fork.sh
```

### View Real-time Logs

```bash
tail -f anvil.log
```

### Find Stale Process

```bash
# If PID file is missing but port is in use
lsof -i :8545

# Kill manually
kill $(lsof -t -i:8545)
```

### Clean Restart

```bash
./scripts/stop-fork.sh
rm -f .anvil.pid anvil.log
./scripts/start-fork.sh 55000000 --daemon
```

### Port Already in Use

```bash
# Find what's using port 8545
lsof -i :8545

# Kill it
kill $(lsof -t -i:8545)

# Or use different port (edit start-fork.sh)
PORT=8546 ./scripts/start-fork.sh 55000000 --daemon
```

## Advanced Usage

### Multiple Forks (Different Blocks)

Edit `start-fork.sh` to support custom port:

```bash
PORT=${PORT:-8545}
```

Then run:

```bash
# Fork at block 55M on port 8545
PORT=8545 DAEMON=1 ./scripts/start-fork.sh 55000000

# Fork at block 56M on port 8546  
PORT=8546 DAEMON=1 ./scripts/start-fork.sh 56000000
```

Update `config/chains.json`:

```json
{
  "polygon-fork-55m": {
    "rpcUrls": ["http://127.0.0.1:8545"]
  },
  "polygon-fork-56m": {
    "rpcUrls": ["http://127.0.0.1:8546"]
  }
}
```

### Auto-restart on Crash

Create systemd service (Linux) or launchd (macOS):

**systemd** (`/etc/systemd/system/anvil-fork.service`):

```ini
[Unit]
Description=Anvil Polygon Fork
After=network.target

[Service]
Type=simple
User=youruser
WorkingDirectory=/path/to/polymarket-indexer
ExecStart=/bin/bash -c 'anvil --fork-url https://polygon-rpc.com --fork-block-number 55000000 --chain-id 137 --port 8545'
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable anvil-fork
sudo systemctl start anvil-fork
sudo systemctl status anvil-fork
```

### Monitoring

Add to monitoring script:

```bash
#!/bin/bash
# monitor-fork.sh

while true; do
    if ! ./scripts/check-fork.sh > /dev/null 2>&1; then
        echo "$(date): Anvil fork is down, restarting..."
        ./scripts/stop-fork.sh
        DAEMON=1 ./scripts/start-fork.sh 55000000
    fi
    sleep 60
done
```

## Best Practices

1. **Always use daemon mode in CI/CD**
   ```bash
   DAEMON=1 ./scripts/start-fork.sh 55000000
   ```

2. **Check status before running tests**
   ```bash
   ./scripts/check-fork.sh && go test ./test -v
   ```

3. **Clean up in CI/CD** (use `if: always()`)
   ```yaml
   - name: Stop Anvil
     if: always()
     run: ./scripts/stop-fork.sh
   ```

4. **Archive logs for debugging**
   ```yaml
   - name: Upload logs
     if: always()
     uses: actions/upload-artifact@v3
     with:
       name: anvil-logs
       path: anvil.log
   ```

5. **Use specific block numbers for reproducibility**
   ```bash
   # Good - deterministic
   ./scripts/start-fork.sh 55000000 --daemon
   
   # Bad - non-deterministic
   ./scripts/start-fork.sh --daemon  # Uses default
   ```

## Summary

Daemon mode makes fork testing production-ready:
- ✅ Survives terminal disconnections
- ✅ Perfect for CI/CD pipelines
- ✅ Easy process management
- ✅ Health monitoring built-in
- ✅ Log persistence for debugging

Use it as the default for all automated testing!
