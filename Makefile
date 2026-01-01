.PHONY: help build test lint clean docker-build docker-push run-indexer run-consumer generate-bindings migrate

# Variables
BINARY_NAME=polymarket-indexer
DOCKER_IMAGE=polymarket-indexer
VERSION?=latest
BUILD=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.build=$(BUILD)"

help: ## Show this help message
	@echo "Polymarket Indexer - Makefile Commands"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## Build Commands

build: ## Build the indexer and consumer binaries
	@echo "Building binaries..."
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/indexer cmd/indexer/main.go
	go build $(LDFLAGS) -o bin/consumer cmd/consumer/main.go
	@echo "✅ Binaries built: bin/indexer, bin/consumer"

build-all: build ## Build all binaries

install: ## Install binaries to $GOPATH/bin
	@echo "Installing binaries..."
	go install $(LDFLAGS) ./cmd/indexer
	go install $(LDFLAGS) ./cmd/consumer
	@echo "✅ Installed to $(shell go env GOPATH)/bin"

## Development Commands

run-indexer: ## Run the indexer locally
	@echo "Starting indexer..."
	go run cmd/indexer/main.go -config config.toml

run-consumer: ## Run the consumer locally
	@echo "Starting consumer..."
	go run cmd/consumer/main.go -config config.toml

dev: ## Run indexer with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

## Testing Commands

test: ## Run all tests
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests passed"

test-coverage: test ## Run tests and show coverage
	go tool cover -html=coverage.out

test-short: ## Run short tests only
	go test -v -short ./...

bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

## Code Quality Commands

lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run --timeout=5m

fmt: ## Format code
	go fmt ./...
	goimports -w .

vet: ## Run go vet
	go vet ./...

## Contract Commands

generate-bindings: ## Generate Go bindings from ABIs
	@echo "Generating contract bindings..."
	@mkdir -p pkg/contracts/bindings
	abigen --abi pkg/contracts/abi/CTFExchange.json --pkg bindings --type CTFExchange --out pkg/contracts/bindings/ctf_exchange.go
	abigen --abi pkg/contracts/abi/ConditionalTokens.json --pkg bindings --type ConditionalTokens --out pkg/contracts/bindings/conditional_tokens.go
	abigen --abi pkg/contracts/abi/ERC20.json --pkg bindings --type ERC20 --out pkg/contracts/bindings/erc20.go
	@echo "✅ Bindings generated"

download-abis: ## Download ABIs from PolygonScan
	@echo "Downloading ABIs from PolygonScan..."
	@mkdir -p pkg/contracts/abi
	curl -s "https://api.polygonscan.com/api?module=contract&action=getabi&address=0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E&apikey=YourApiKey" | jq -r .result > pkg/contracts/abi/CTFExchange.json
	curl -s "https://api.polygonscan.com/api?module=contract&action=getabi&address=0x4D97DCd97eC945f40cF65F87097ACe5EA0476045&apikey=YourApiKey" | jq -r .result > pkg/contracts/abi/ConditionalTokens.json
	@echo "✅ ABIs downloaded"

## Database Commands

migrate-up: ## Run database migrations
	@echo "Running migrations..."
	@docker exec -i polymarket-timescaledb psql -U polymarket -d polymarket < migrations/001_initial_schema.up.sql
	@echo "✅ Migrations applied"

migrate-down: ## Rollback last migration
	@echo "Rolling back migration..."
	@echo "⚠️  Manual rollback required - no down migration implemented"
	go run cmd/migrate/main.go down
	@echo "✅ Migration rolled back"

migrate-create: ## Create a new migration (usage: make migrate-create NAME=add_markets_table)
	@if [ -z "$(NAME)" ]; then echo "❌ NAME is required. Usage: make migrate-create NAME=add_markets_table"; exit 1; fi
	@echo "Creating migration: $(NAME)"
	@mkdir -p migrations
	@touch migrations/$(shell date +%Y%m%d%H%M%S)_$(NAME).up.sql
	@touch migrations/$(shell date +%Y%m%d%H%M%S)_$(NAME).down.sql
	@echo "✅ Migration created"

## Docker Commands

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):$(BUILD) .
	@echo "✅ Docker image built: $(DOCKER_IMAGE):$(VERSION)"

docker-build-indexer: ## Build indexer Docker image
	docker build -t $(DOCKER_IMAGE)-indexer:$(VERSION) --target indexer .

docker-build-consumer: ## Build consumer Docker image
	docker build -t $(DOCKER_IMAGE)-consumer:$(VERSION) --target consumer .

docker-push: docker-build ## Push Docker images to registry
	docker push $(DOCKER_IMAGE):$(VERSION)
	docker push $(DOCKER_IMAGE):$(BUILD)

docker-up: ## Start all services with docker-compose
	docker-compose up -d

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show logs from all services
	docker-compose logs -f

docker-logs-indexer: ## Show indexer logs
	docker-compose logs -f indexer

docker-logs-consumer: ## Show consumer logs
	docker-compose logs -f consumer

## Infrastructure Commands

infra-up: ## Start infrastructure (NATS + TimescaleDB)
	docker-compose up -d nats timescaledb
	@echo "⏳ Waiting for services to be ready..."
	@sleep 5
	@echo "✅ Infrastructure ready"

infra-down: ## Stop infrastructure
	docker-compose down nats timescaledb

infra-reset: infra-down ## Reset infrastructure (delete volumes)
	docker-compose down -v
	rm -rf data/
	@echo "✅ Infrastructure reset"

## Maintenance Commands

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf dist/
	rm -f coverage.out
	rm -rf data/checkpoints.db
	@echo "✅ Cleaned"

clean-all: clean ## Clean everything including caches
	go clean -cache -testcache -modcache
	docker-compose down -v

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies downloaded"

deps-upgrade: ## Upgrade all dependencies
	@echo "Upgrading dependencies..."
	go get -u ./...
	go mod tidy
	@echo "✅ Dependencies upgraded"

## Monitoring Commands

stats: ## Show indexer stats
	@curl -s http://localhost:8080/stats | jq

health: ## Check service health
	@curl -s http://localhost:8080/health | jq

metrics: ## Show Prometheus metrics
	@curl -s http://localhost:8080/metrics

## Production Commands

release: clean test lint build ## Create a release (test + lint + build)
	@echo "✅ Release ready: bin/"

production-deploy: ## Deploy to production (placeholder)
	@echo "🚀 Deploying to production..."
	@echo "⚠️  Implement your deployment strategy here"

## Utility Commands

check-env: ## Check environment setup
	@echo "Checking environment..."
	@which go > /dev/null && echo "✅ Go installed: $(shell go version)" || echo "❌ Go not found"
	@which docker > /dev/null && echo "✅ Docker installed" || echo "❌ Docker not found"
	@which docker-compose > /dev/null && echo "✅ Docker Compose installed" || echo "❌ Docker Compose not found"
	@which abigen > /dev/null && echo "✅ abigen installed" || echo "⚠️  abigen not found (go install github.com/ethereum/go-ethereum/cmd/abigen@latest)"
	@which golangci-lint > /dev/null && echo "✅ golangci-lint installed" || echo "⚠️  golangci-lint not found"

setup: deps check-env infra-up ## Initial setup (deps + infra)
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Update config.toml with your Polygon RPC endpoint"
	@echo "  2. Run: make run-indexer"
	@echo "  3. Run: make run-consumer (in another terminal)"

.DEFAULT_GOAL := help
