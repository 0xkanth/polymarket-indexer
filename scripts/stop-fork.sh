#!/bin/bash

# Stop the Anvil fork daemon
# Usage: ./scripts/stop-fork.sh

set -e

PID_FILE=".anvil.pid"
LOG_FILE="anvil.log"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

if [ ! -f "$PID_FILE" ]; then
    echo -e "${YELLOW}No Anvil process found (PID file doesn't exist)${NC}"
    exit 0
fi

PID=$(cat "$PID_FILE")

if ! ps -p "$PID" > /dev/null 2>&1; then
    echo -e "${YELLOW}Anvil process (PID: $PID) is not running${NC}"
    rm -f "$PID_FILE"
    exit 0
fi

echo -e "${YELLOW}Stopping Anvil (PID: $PID)...${NC}"
kill "$PID" 2>/dev/null || true

# Wait for process to stop (max 5 seconds)
for i in {1..10}; do
    if ! ps -p "$PID" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Anvil stopped successfully${NC}"
        rm -f "$PID_FILE"
        exit 0
    fi
    sleep 0.5
done

# Force kill if still running
if ps -p "$PID" > /dev/null 2>&1; then
    echo -e "${YELLOW}Force killing Anvil...${NC}"
    kill -9 "$PID" 2>/dev/null || true
    rm -f "$PID_FILE"
    echo -e "${GREEN}✓ Anvil force stopped${NC}"
else
    rm -f "$PID_FILE"
    echo -e "${GREEN}✓ Anvil stopped${NC}"
fi
