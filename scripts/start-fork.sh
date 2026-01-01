#!/bin/bash

# Fork Polygon mainnet at specific block number
# Usage: 
#   ./scripts/start-fork.sh [BLOCK_NUMBER]           - Run in foreground
#   ./scripts/start-fork.sh [BLOCK_NUMBER] --daemon  - Run in background
#   DAEMON=1 ./scripts/start-fork.sh [BLOCK_NUMBER]  - Run in background (env var)

set -e

# Configuration
FORK_BLOCK_NUMBER=${1:-55000000}
POLYGON_RPC="https://polygon-rpc.com"
CHAIN_ID=137
PORT=8545
PID_FILE=".anvil.pid"
LOG_FILE="anvil.log"

# Check for daemon mode
DAEMON_MODE=false
if [[ "$2" == "--daemon" ]] || [[ "$DAEMON" == "1" ]] || [[ "$CI" == "true" ]]; then
    DAEMON_MODE=true
fi

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Check if Anvil is installed
if ! command -v anvil &> /dev/null; then
    echo -e "${RED}Error: Anvil not found${NC}"
    echo "Install Foundry first:"
    echo "  curl -L https://foundry.paradigm.xyz | bash"
    echo "  foundryup"
    exit 1
fi

# Check if already running
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        echo -e "${YELLOW}Anvil is already running (PID: $PID)${NC}"
        echo "Use ./scripts/stop-fork.sh to stop it first"
        exit 1
    else
        echo -e "${YELLOW}Removing stale PID file${NC}"
        rm -f "$PID_FILE"
    fi
fi

echo -e "${GREEN}Starting Anvil fork of Polygon mainnet...${NC}"
echo -e "${YELLOW}Fork Block: ${FORK_BLOCK_NUMBER}${NC}"
echo -e "${YELLOW}Chain ID: ${CHAIN_ID}${NC}"
echo -e "${YELLOW}Port: ${PORT}${NC}"
echo -e "${YELLOW}Mode: $([ "$DAEMON_MODE" = true ] && echo 'Daemon' || echo 'Foreground')${NC}"
echo ""

if [ "$DAEMON_MODE" = true ]; then
    # Start in background
    nohup anvil \
      --fork-url "${POLYGON_RPC}" \
      --fork-block-number "${FORK_BLOCK_NUMBER}" \
      --chain-id "${CHAIN_ID}" \
      --port "${PORT}" \
      --host "127.0.0.1" \
      --block-time 2 \
      --accounts 10 \
      --balance 10000 \
      > "$LOG_FILE" 2>&1 &
    
    ANVIL_PID=$!
    echo $ANVIL_PID > "$PID_FILE"
    
    # Wait a bit and check if it started successfully
    sleep 2
    if ps -p "$ANVIL_PID" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Anvil started successfully in background${NC}"
        echo -e "${GREEN}  PID: $ANVIL_PID${NC}"
        echo -e "${GREEN}  Log: $LOG_FILE${NC}"
        echo ""
        echo "To stop: ./scripts/stop-fork.sh"
        echo "To check status: ./scripts/check-fork.sh"
        echo "To view logs: tail -f $LOG_FILE"
    else
        echo -e "${RED}✗ Failed to start Anvil${NC}"
        rm -f "$PID_FILE"
        echo "Check $LOG_FILE for details"
        exit 1
    fi
else
    # Start in foreground
    anvil \
      --fork-url "${POLYGON_RPC}" \
      --fork-block-number "${FORK_BLOCK_NUMBER}" \
      --chain-id "${CHAIN_ID}" \
      --port "${PORT}" \
      --host "127.0.0.1" \
      --block-time 2 \
      --accounts 10 \
      --balance 10000
fi
