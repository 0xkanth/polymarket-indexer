#!/bin/bash

# Check if Anvil fork is running
# Usage: ./scripts/check-fork.sh

set -e

PID_FILE=".anvil.pid"
PORT=8545

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "Checking Anvil fork status..."
echo ""

# Check PID file
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Process Status: Running (PID: $PID)${NC}"
        
        # Get process info
        CPU_MEM=$(ps -p "$PID" -o %cpu,%mem,etime | tail -n 1)
        echo -e "${GREEN}  CPU/MEM/Uptime: $CPU_MEM${NC}"
    else
        echo -e "${RED}✗ Process Status: Not running (stale PID file)${NC}"
        echo -e "${YELLOW}  Run: rm $PID_FILE${NC}"
    fi
else
    echo -e "${YELLOW}○ Process Status: No PID file found${NC}"
fi

echo ""

# Check if port is listening
if lsof -i :$PORT > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Port $PORT: Listening${NC}"
    PORT_INFO=$(lsof -i :$PORT | tail -n 1)
    echo -e "${GREEN}  $PORT_INFO${NC}"
else
    echo -e "${RED}✗ Port $PORT: Not listening${NC}"
fi

echo ""

# Check RPC endpoint
if command -v curl &> /dev/null; then
    RESPONSE=$(curl -s -X POST http://127.0.0.1:$PORT \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' 2>/dev/null || echo "")
    
    if [[ $RESPONSE == *"result"* ]]; then
        BLOCK=$(echo $RESPONSE | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
        BLOCK_DEC=$((16#${BLOCK:2}))
        echo -e "${GREEN}✓ RPC Endpoint: Responding${NC}"
        echo -e "${GREEN}  Current Block: $BLOCK_DEC (hex: $BLOCK)${NC}"
        
        # Get chain ID
        CHAIN_RESPONSE=$(curl -s -X POST http://127.0.0.1:$PORT \
            -H "Content-Type: application/json" \
            -d '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' 2>/dev/null || echo "")
        
        if [[ $CHAIN_RESPONSE == *"result"* ]]; then
            CHAIN_ID=$(echo $CHAIN_RESPONSE | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
            CHAIN_ID_DEC=$((16#${CHAIN_ID:2}))
            echo -e "${GREEN}  Chain ID: $CHAIN_ID_DEC${NC}"
        fi
    else
        echo -e "${RED}✗ RPC Endpoint: Not responding${NC}"
    fi
else
    echo -e "${YELLOW}○ RPC Check: curl not available${NC}"
fi

echo ""

# Summary
if [ -f "$PID_FILE" ] && ps -p "$(cat $PID_FILE)" > /dev/null 2>&1 && lsof -i :$PORT > /dev/null 2>&1; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}Status: ✓ Anvil fork is running properly${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
else
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}Status: ✗ Anvil fork is not running${NC}"
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "To start: ./scripts/start-fork.sh --daemon"
fi
