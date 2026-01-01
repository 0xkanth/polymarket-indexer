#!/bin/bash

# Generate Go bindings from ABI files
# Usage: ./scripts/generate-bindings.sh

set -e

echo "🔧 Generating Go contract bindings from ABI files..."

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if abigen is installed
if ! command -v abigen &> /dev/null; then
    echo "❌ abigen not found. Installing..."
    go install github.com/ethereum/go-ethereum/cmd/abigen@latest
    echo "✅ abigen installed"
fi

# Project root
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$PROJECT_ROOT"

# ABI directory
ABI_DIR="pkg/contracts/abi"
OUT_DIR="pkg/contracts"

# Ensure output directory exists
mkdir -p "$OUT_DIR"

echo ""
echo "${BLUE}Generating bindings...${NC}"
echo ""

# Generate CTFExchange
echo "📝 Generating CTFExchange..."
abigen \
  --abi "$ABI_DIR/CTFExchange.json" \
  --pkg contracts \
  --type CTFExchange \
  --out "$OUT_DIR/CTFExchange.go"
echo "${GREEN}✅ CTFExchange.go${NC}"

# Generate ConditionalTokens
echo "📝 Generating ConditionalTokens..."
abigen \
  --abi "$ABI_DIR/ConditionalTokens.json" \
  --pkg contracts \
  --type ConditionalTokens \
  --out "$OUT_DIR/ConditionalTokens.go"
echo "${GREEN}✅ ConditionalTokens.go${NC}"

# Generate ERC20
echo "📝 Generating ERC20..."
abigen \
  --abi "$ABI_DIR/ERC20.json" \
  --pkg contracts \
  --type ERC20 \
  --out "$OUT_DIR/ERC20.go"
echo "${GREEN}✅ ERC20.go${NC}"

echo ""
echo "${GREEN}🎉 All contract bindings generated successfully!${NC}"
echo ""
echo "Generated files:"
echo "  - $OUT_DIR/CTFExchange.go"
echo "  - $OUT_DIR/ConditionalTokens.go"
echo "  - $OUT_DIR/ERC20.go"
