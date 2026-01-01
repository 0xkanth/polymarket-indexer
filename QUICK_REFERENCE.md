# 📇 Quick Reference Card - Print & Keep on Desk

```
┌─────────────────────────────────────────────────────────────────────────┐
│                   POLYMARKET INDEXER QUICK REFERENCE                     │
└─────────────────────────────────────────────────────────────────────────┘

🏗️  ARCHITECTURE (Indexer → Syncer → Processor → Router → Handler → NATS → Consumer → DB)
   
   ┌─────────┐    ┌─────────┐    ┌──────────┐    ┌────────┐    ┌──────┐
   │Polygon  │ -> │ Syncer  │ -> │Processor │ -> │ Router │ -> │ NATS │
   │   RPC   │    │(Orchestr)│    │(Extract) │    │(Route) │    │(Bus) │
   └─────────┘    └─────────┘    └──────────┘    └────────┘    └──┬───┘
                       ↕                                            ↓
                  ┌──────────┐                              ┌──────────┐
                  │Checkpoint│                              │ Consumer │
                  │   DB     │                              └────┬─────┘
                  └──────────┘                                   ↓
                                                          ┌──────────────┐
                                                          │ TimescaleDB  │
                                                          └──────────────┘

🔄 DATA FLOW (10 Steps)
   1. Syncer determines block to process (backfill/realtime)
   2. Syncer calls processor.ProcessBlock(blockNum)
   3. Processor fetches logs via chainClient.FilterLogs()
   4. For each log, processor calls processLog()
   5. processLog() calls router.RouteLog(log, timestamp, blockHash)
   6. Router extracts event signature (log.Topics[0])
   7. Router calls handler to decode ABI
   8. Handler returns payload, router calls callback
   9. Callback publishes to NATS JetStream
   10. Consumer receives, writes to TimescaleDB

💾 DATABASE (TimescaleDB = PostgreSQL + time-series)
   Tables:
   - events (hypertable)     All events (source of truth)
   - order_fills             OrderFilled events (denormalized)
   - token_transfers         ERC1155 transfers
   - conditions              Market conditions
   - token_registrations     Token → condition mapping
   - checkpoints             ⚠️ CRITICAL: Sync progress

   Checkpoints = {service_name, last_block, last_block_hash, updated_at}
   Purpose: Resume from crash, detect reorgs

🔧 KEY COMMANDS
   Infrastructure:
   $ make infra-up           Start NATS + TimescaleDB (Docker)
   $ make infra-down         Stop infrastructure
   $ make migrate-up         Apply database migrations
   
   Build & Run:
   $ make build              Compile indexer + consumer
   $ make run-indexer        Start indexer (producer)
   $ make run-consumer       Start consumer (database writer)
   
   Testing:
   $ make test               Run unit tests
   $ make fork-test          Run fork tests (requires Anvil)

📊 MONITORING
   Indexer:  http://localhost:8080/health
   Metrics:  http://localhost:8080/metrics
   NATS:     http://localhost:8222
   Database: docker exec -it polymarket-timescaledb psql -U polymarket

   Key Metrics:
   - syncer_current_block      Last block processed
   - chain_latest_block        Latest block on chain
   - syncer_blocks_behind      How far behind (should be < 100)
   - syncer_errors_total       Error count by type

🐛 DEBUGGING (When things break)
   Check logs:
   $ docker logs polymarket-indexer --tail 100
   $ docker logs polymarket-consumer --tail 100
   $ docker logs polymarket-nats
   
   Check database:
   $ docker exec -it polymarket-timescaledb psql -U polymarket
   polymarket=# SELECT * FROM checkpoints;
   polymarket=# SELECT COUNT(*) FROM events;
   polymarket=# SELECT block, timestamp FROM events ORDER BY block DESC LIMIT 10;
   
   Check NATS:
   $ curl http://localhost:8222/streaming/channelsz?subs=1
   
   Check metrics:
   $ curl http://localhost:8080/metrics | grep blocks_behind

🚨 ON-CALL QUICK FIXES
   Issue: Indexer stuck
   → Check: docker ps | grep indexer (running?)
   → Check: docker logs polymarket-indexer (errors?)
   → Check: RPC rate limits (switch provider in chains.json)
   → Fix: docker restart polymarket-indexer
   
   Issue: Consumer not writing
   → Check: NATS connection (http://localhost:8222)
   → Check: Database connection (psql)
   → Fix: docker restart polymarket-consumer
   
   Issue: Sync lag growing
   → Check: curl http://localhost:8080/metrics | grep blocks_behind
   → Reduce workers: config.toml → indexer.workers = 3 (from 5)
   → Reduce batch: config.toml → indexer.batch_size = 500 (from 1000)
   
   Issue: Database full
   → Enable compression:
     ALTER TABLE events SET (timescaledb.compress);
     SELECT add_compression_policy('events', INTERVAL '7 days');
   
   Issue: Bad checkpoint
   → Manual rollback:
     UPDATE checkpoints SET last_block = last_block - 1000
     WHERE service_name = 'polymarket-indexer';

📁 KEY FILES (Where to look)
   Entry Points:
   - cmd/indexer/main.go       Indexer entry (creates all components)
   - cmd/consumer/main.go      Consumer entry (NATS → DB)
   
   Core Logic:
   - internal/syncer/syncer.go         Block sync orchestrator
   - internal/processor/processor.go   Event extraction
   - internal/router/router.go         Event routing
   - internal/handler/events.go        ABI decoding
   
   Infrastructure:
   - internal/chain/client.go          RPC client
   - internal/nats/publisher.go        NATS publisher
   - internal/db/checkpoint.go         Checkpoint management
   
   Configuration:
   - config.toml                Runtime config (workers, batch size)
   - config/chains.json         Chain config (RPC, contracts, start block)
   
   Database:
   - migrations/001_initial_schema.up.sql   Database schema

🔐 CONFIGURATION
   config.toml (runtime behavior):
   [chain]
   name = "polygon"                    # Which chain from chains.json
   
   [indexer]
   batch_size = 1000                   # Blocks per batch (backfill)
   poll_interval = "2s"                # Polling frequency (realtime)
   workers = 5                         # Parallel workers (backfill)
   
   [syncer]
   confirmations = 100                 # Wait N blocks before processing
   
   chains.json (chain data):
   {
     "polygon": {
       "rpcUrls": ["https://polygon-rpc.com"],
       "wsUrls": ["wss://polygon-ws.com"],
       "contracts": {
         "ctfExchange": "0x4bFb41d5B3570DeFd03C39a9A4D8dE6Bd8B8982E",
         "conditionalTokens": "0x4D97DCd97eC945f40cF65F87097ACe5EA0476045"
       },
       "startBlock": 20558323,         # Block to start indexing from
       "confirmations": 100
     }
   }

🧪 TESTING
   Unit tests:
   $ go test ./...                     # All packages
   $ go test ./internal/processor/... # Specific package
   
   Fork tests (requires Anvil):
   $ make fork-start                   # Start Anvil fork
   $ make fork-test                    # Run fork tests
   $ make fork-stop                    # Stop Anvil

📖 DOCUMENTATION (Read in this order!)
   1. README.md                        Quick start & setup
   2. docs/MASTER_LEARNING_PATH.md     ⭐ Complete learning guide
   3. docs/ARCHITECTURE.md             System design
   4. docs/COMPONENT_CONNECTIONS.md    Component interactions
   5. docs/DATABASE.md                 Database & checkpoints
   6. docs/SYNCER_ARCHITECTURE.md      Sync strategies
   7. docs/NATS_EXPLAINED.md           Message bus
   8. docs/RUNNING.md                  Production operations
   9. docs/TESTING.md                  Testing guide

   Quick Reference:
   - docs/DOCUMENTATION_INDEX.md       Master index of all docs

🎓 LEARNING PATH (Complete in 1 week)
   Day 1: Setup & Run (2 hours)
   Day 2: Understanding (3 hours) - Architecture, Components, Database
   Day 3: Deep Dive (4 hours) - Syncer code walkthrough
   Day 4: Event Processing (3 hours) - Processor, Router, Handlers
   Day 5: Testing & Debugging (4 hours) - Write tests, simulate failures

🆘 NEED HELP?
   Slack: #polymarket-indexer
   Docs: docs/MASTER_LEARNING_PATH.md (On-Call Runbook section)
   Logs: docker logs polymarket-indexer --tail 100
   Metrics: http://localhost:8080/metrics

💡 REMEMBER
   - Checkpoints are CRITICAL for recovery
   - Syncer switches modes automatically (backfill ↔ realtime)
   - NATS provides durability (7 day retention)
   - Consumer can be scaled (multiple replicas)
   - Database uses hypertables (automatic partitioning by time)
   - Confirmations protect against reorgs (wait 100 blocks)

└─────────────────────────────────────────────────────────────────────────┘
```

**Print this card and keep it on your desk for quick reference!**

**Last Updated**: January 1, 2026  
**Version**: 1.0.0
