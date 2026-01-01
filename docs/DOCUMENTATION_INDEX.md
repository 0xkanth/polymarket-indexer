# Documentation Index

**Complete guide to all documentation files in the Polymarket Indexer project.**

## 📋 Quick Reference

| File | Purpose | Read Time | Has Diagrams | For Whom |
|------|---------|-----------|--------------|----------|
| [README.md](../README.md) | **START HERE**: Setup, quick start, features | 10 min | ✅ Mermaid | Everyone |
| **[MASTER_LEARNING_PATH.md](MASTER_LEARNING_PATH.md)** | **Complete learning guide** for naive developers | 30 min | ✅ Mermaid | **New developers** |
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design, component overview | 15 min | ASCII | Backend devs |
| [COMPONENT_CONNECTIONS.md](COMPONENT_CONNECTIONS.md) | Component interactions, data flow | 20 min | ✅ Mermaid | Backend devs |
| [DATABASE.md](DATABASE.md) | Schema, checkpoints, why TimescaleDB | 25 min | ✅ Mermaid | Backend/Data devs |
| [SYNCER_ARCHITECTURE.md](SYNCER_ARCHITECTURE.md) | Block sync strategies (backfill/realtime) | 20 min | ✅ Mermaid + ASCII | Backend devs |
| [NATS_EXPLAINED.md](NATS_EXPLAINED.md) | Message bus, producer-consumer pattern | 15 min | ASCII | Backend devs |
| [RUNNING.md](RUNNING.md) | Production deployment, monitoring | 15 min | - | DevOps/SRE |
| [TESTING.md](TESTING.md) | Unit/integration/fork tests | 20 min | - | Developers |
| [DEPENDENCIES.md](DEPENDENCIES.md) | External library docs | 10 min | - | Developers |
| [GO_ETHEREUM_GUIDE.md](GO_ETHEREUM_GUIDE.md) | go-ethereum library usage | 15 min | - | Blockchain devs |
| [TRANSACTION_HELPERS.md](TRANSACTION_HELPERS.md) | Transaction parsing utilities | 10 min | - | Blockchain devs |
| [FORK_TESTING_GUIDE.md](FORK_TESTING_GUIDE.md) | Advanced fork testing with Anvil | 20 min | - | Developers |
| [DAEMON_MODE.md](DAEMON_MODE.md) | Running as system service | 10 min | - | DevOps/SRE |

**Total Reading Time**: ~4 hours to become proficient

---

## 🎯 Reading Order by Goal

### Goal: Get System Running (30 minutes)
1. [README.md](../README.md) - Setup steps
2. Follow Quick Start guide
3. Run `make infra-up && make migrate-up && make build`

---

### Goal: Understand Architecture (1 hour)
1. [README.md](../README.md) - Overview
2. [ARCHITECTURE.md](ARCHITECTURE.md) - System design
3. [COMPONENT_CONNECTIONS.md](COMPONENT_CONNECTIONS.md) - How components talk
4. [DATABASE.md](DATABASE.md) - Data storage

---

### Goal: Become On-Call Ready (2 hours)
1. [README.md](../README.md) - Setup
2. [ARCHITECTURE.md](ARCHITECTURE.md) - High-level understanding
3. [MASTER_LEARNING_PATH.md](MASTER_LEARNING_PATH.md) - Jump to "On-Call Runbook" section
4. [RUNNING.md](RUNNING.md) - Monitoring dashboards
5. [DATABASE.md](DATABASE.md) - Checkpoint recovery

---

### Goal: Master the Codebase (1 week)
Follow the [MASTER_LEARNING_PATH.md](MASTER_LEARNING_PATH.md) Day 1-5 exercises:
- Day 1: Setup & Run
- Day 2: Understanding (Architecture, Components, Database)
- Day 3: Deep Dive (Syncer)
- Day 4: Event Processing (Processor, Router, Handlers)
- Day 5: Testing & Debugging

---

## 📚 Document Details

### Level 1: Essential (Read First)

#### [README.md](../README.md)
**Purpose**: Project overview, quick start, setup instructions  
**Contains**:
- What is Polymarket indexer?
- Architecture diagram (Mermaid)
- Quick start (8 steps)
- Configuration guide
- Monitoring basics

**Key Diagrams**:
- ✅ High-level architecture (Mermaid flowchart)

**When to Read**: First thing before doing anything

---

#### [MASTER_LEARNING_PATH.md](MASTER_LEARNING_PATH.md) ⭐ NEW
**Purpose**: Complete learning guide for naive developers  
**Contains**:
- Logical document reading order
- Learning paths by role (Backend, Blockchain, DevOps, On-Call)
- Mental models (Pipeline, State Machine, Dependency Tree)
- Hands-on exercises (Day 1-5)
- On-call runbook
- Debugging mindmap
- Success criteria (Beginner → Expert)

**Key Diagrams**:
- ✅ Learning roadmap (Mermaid graph)
- ✅ State machine (Mermaid)
- ✅ Debugging mindmap (Mermaid)

**When to Read**: After setup, before diving into code

---

#### [ARCHITECTURE.md](ARCHITECTURE.md)
**Purpose**: High-level system design  
**Contains**:
- Producer-consumer architecture
- Why NATS JetStream?
- Why TimescaleDB?
- Component responsibilities
- Design decisions

**Key Diagrams**:
- ASCII architecture diagram

**When to Read**: After setup, to understand the big picture

---

### Level 2: Core Understanding

#### [COMPONENT_CONNECTIONS.md](COMPONENT_CONNECTIONS.md) ⭐ UPDATED
**Purpose**: Detailed component interactions  
**Contains**:
- How Indexer → Syncer → Processor → Router → Handler → NATS → Consumer → DB
- Dependency injection patterns
- Data flow (10 steps)
- Component relationship table
- Why consumer is separate

**Key Diagrams**:
- ✅ **NEW** High-level architecture (Mermaid graph)
- ✅ **NEW** Sequence diagram: Complete event flow (Mermaid)
- ✅ **NEW** Component dependency graph (Mermaid)
- ASCII architecture diagram

**When to Read**: When you want to understand how components connect

---

#### [DATABASE.md](DATABASE.md) ⭐ NEW
**Purpose**: Database architecture, schema, checkpoints  
**Contains**:
- What is TimescaleDB? (PostgreSQL + time-series)
- Why TimescaleDB vs MongoDB/Cassandra/PostgreSQL?
- Database schema (ER diagram)
- How checkpoints work (sequence diagram)
- Query examples
- Backup & recovery
- Performance tuning

**Key Diagrams**:
- ✅ **NEW** TimescaleDB hypertable diagram (Mermaid)
- ✅ **NEW** Database ER diagram (Mermaid)
- ✅ **NEW** Checkpoint flow sequence (Mermaid)
- ✅ **NEW** Data flow diagram (Mermaid)

**When to Read**: When you want to understand data storage and recovery

---

#### [SYNCER_ARCHITECTURE.md](SYNCER_ARCHITECTURE.md) ⭐ UPDATED
**Purpose**: Block synchronization strategies  
**Contains**:
- Why dual-mode strategy (backfill vs realtime)?
- How syncer orchestrates sync
- Worker pool coordination
- Automatic mode switching
- Configuration options
- Metrics exposed

**Key Diagrams**:
- ✅ **NEW** State machine diagram (Mermaid)
- ✅ **NEW** Sequence diagram: Syncer lifecycle (Mermaid)
- ASCII syncer lifecycle flow
- ASCII worker pool distribution

**When to Read**: When you want to understand block sync orchestration

---

#### [NATS_EXPLAINED.md](NATS_EXPLAINED.md)
**Purpose**: Message bus architecture  
**Contains**:
- What is NATS JetStream?
- Why use message bus (decoupling)
- Streams vs subjects
- Message retention & replay
- Consumer groups
- Exactly-once delivery

**Key Diagrams**:
- ASCII producer-consumer flow

**When to Read**: When you want to understand the messaging layer

---

### Level 3: Operations & Testing

#### [RUNNING.md](RUNNING.md)
**Purpose**: Production deployment and monitoring  
**Contains**:
- How to deploy
- Prometheus metrics
- Grafana dashboards
- Health checks
- Scaling strategies
- Troubleshooting

**When to Read**: Before deploying to production

---

#### [TESTING.md](TESTING.md)
**Purpose**: Testing strategies  
**Contains**:
- Unit testing patterns
- Fork testing with Anvil
- Integration tests
- Mocking strategies
- Test coverage
- CI/CD integration

**When to Read**: When writing or fixing tests

---

### Level 4: Deep Dives

#### [DEPENDENCIES.md](DEPENDENCIES.md)
**Purpose**: External library documentation  
**Contains**:
- go-ethereum (blockchain interaction)
- NATS (messaging)
- pgx (PostgreSQL driver)
- zerolog (logging)
- prometheus (metrics)

**When to Read**: When you need to modify core functionality

---

#### [GO_ETHEREUM_GUIDE.md](GO_ETHEREUM_GUIDE.md)
**Purpose**: go-ethereum library usage  
**Contains**:
- RPC client setup
- Filtering logs
- ABI encoding/decoding
- Transaction parsing
- Block fetching

**When to Read**: When working with blockchain interactions

---

#### [TRANSACTION_HELPERS.md](TRANSACTION_HELPERS.md)
**Purpose**: Transaction parsing utilities  
**Contains**:
- Helper functions for transaction parsing
- Event extraction
- ABI utilities

**When to Read**: When adding new event handlers

---

#### [FORK_TESTING_GUIDE.md](FORK_TESTING_GUIDE.md)
**Purpose**: Advanced fork testing  
**Contains**:
- Setting up Anvil
- Forking mainnet
- Testing against real state
- Debugging fork tests

**When to Read**: When testing against mainnet state

---

#### [DAEMON_MODE.md](DAEMON_MODE.md)
**Purpose**: Running as system service  
**Contains**:
- systemd service setup
- Logging configuration
- Auto-restart on failure
- Security hardening

**When to Read**: When deploying to production servers

---

## 🎨 Diagram Types by Document

### Mermaid Diagrams (Interactive, GitHub-rendered)

**README.md**:
- ✅ High-level architecture (graph LR)

**MASTER_LEARNING_PATH.md**:
- ✅ Learning roadmap (graph TD)
- ✅ Path by role (graph LR)
- ✅ State machine (stateDiagram-v2)
- ✅ Debugging mindmap (mindmap)

**COMPONENT_CONNECTIONS.md**:
- ✅ High-level architecture (graph TB)
- ✅ Sequence diagram (sequenceDiagram)
- ✅ Dependency graph (graph TD)

**DATABASE.md**:
- ✅ TimescaleDB hypertables (graph LR)
- ✅ Database ER diagram (erDiagram)
- ✅ Checkpoint flow (sequenceDiagram)
- ✅ Data flow (flowchart TB)

**SYNCER_ARCHITECTURE.md**:
- ✅ State machine (stateDiagram-v2)
- ✅ Lifecycle sequence (sequenceDiagram)

---

### ASCII Diagrams (Text-based, always visible)

**ARCHITECTURE.md**:
- Producer-consumer architecture
- Component overview

**COMPONENT_CONNECTIONS.md**:
- Detailed ASCII architecture (fallback)

**SYNCER_ARCHITECTURE.md**:
- Syncer lifecycle flow (box drawing)
- Worker pool distribution

**NATS_EXPLAINED.md**:
- Message flow
- Stream architecture

---

## 📊 Documentation Metrics

| Metric | Value |
|--------|-------|
| **Total Documents** | 15 |
| **Total Reading Time** | ~4 hours |
| **Mermaid Diagrams** | 15+ (NEW!) |
| **ASCII Diagrams** | 10+ |
| **Code Examples** | 50+ |
| **Links to Code** | 100+ |

---

## 🔄 Document Dependencies

```
README.md (Start)
    ↓
MASTER_LEARNING_PATH.md (Learning guide)
    ↓
    ├─→ ARCHITECTURE.md (System design)
    │       ↓
    │   COMPONENT_CONNECTIONS.md (Interactions)
    │       ↓
    │       ├─→ DATABASE.md (Storage)
    │       ├─→ SYNCER_ARCHITECTURE.md (Sync)
    │       └─→ NATS_EXPLAINED.md (Messaging)
    │
    ├─→ RUNNING.md (Operations)
    │       ↓
    │   DAEMON_MODE.md (Deployment)
    │
    └─→ TESTING.md (Testing)
            ↓
        FORK_TESTING_GUIDE.md (Advanced)
```

---

## 🎓 Certification Checklist

Use this checklist to track your progress through the documentation:

### Bronze Level (Beginner)
- [ ] Read README.md
- [ ] Read MASTER_LEARNING_PATH.md
- [ ] Read ARCHITECTURE.md
- [ ] Setup local environment
- [ ] Query database successfully

### Silver Level (Intermediate)
- [ ] Read COMPONENT_CONNECTIONS.md
- [ ] Read DATABASE.md
- [ ] Read SYNCER_ARCHITECTURE.md
- [ ] Understand data flow
- [ ] Can explain architecture to someone else

### Gold Level (Advanced)
- [ ] Read NATS_EXPLAINED.md
- [ ] Read RUNNING.md
- [ ] Read TESTING.md
- [ ] Write a unit test
- [ ] Fix a bug

### Diamond Level (Expert)
- [ ] Read all documentation
- [ ] Complete Day 1-5 hands-on exercises
- [ ] Handle a production incident
- [ ] Add a new feature
- [ ] Mentor a new team member

---

## 🆘 Need Help?

**Can't find what you're looking for?**

1. Check [MASTER_LEARNING_PATH.md](MASTER_LEARNING_PATH.md) - It has a comprehensive index
2. Search this file (Ctrl+F / Cmd+F)
3. Check [README.md](../README.md) for quick links
4. Ask in Slack #polymarket-indexer

**Want to improve documentation?**

1. Open an issue describing what's unclear
2. Submit a PR with improvements
3. Update this index if you add new docs

---

**Last Updated**: January 1, 2026  
**Maintained By**: Polymarket Indexer Team  
**Total Documents**: 15  
**New Documents This Release**: 2 (DATABASE.md, MASTER_LEARNING_PATH.md)  
**Updated Documents**: 3 (README.md, COMPONENT_CONNECTIONS.md, SYNCER_ARCHITECTURE.md)
