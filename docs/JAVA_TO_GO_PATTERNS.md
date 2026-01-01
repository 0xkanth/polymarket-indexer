# 🔄 Java/Spring Boot to Go Pattern Translation Guide

> **For Java developers migrating to Go: How to achieve polymorphism, inheritance, and Spring Boot patterns in Go**

## 📚 Table of Contents

1. [Core Philosophy Differences](#core-philosophy-differences)
2. [Polymorphism in Go](#polymorphism-in-go)
3. [Inheritance vs Composition](#inheritance-vs-composition)
4. [Spring Boot Patterns in Go](#spring-boot-patterns-in-go)
5. [Service Layer Patterns](#service-layer-patterns)
6. [Dependency Injection](#dependency-injection)
7. [Real Examples from Your Codebase](#real-examples-from-your-codebase)

---

## 🧠 Core Philosophy Differences

| Concept | Java/Spring Boot | Go |
|---------|------------------|-----|
| **Inheritance** | Classes extend classes | Composition + interfaces |
| **Polymorphism** | Abstract classes, interfaces | Interfaces only (implicit) |
| **DI Framework** | Spring IoC container | Manual DI or libraries like Wire/Fx |
| **Annotations** | `@Service`, `@Autowired`, etc. | Struct tags (JSON/YAML only) |
| **Generics** | Full generics support | Type parameters (Go 1.18+) |
| **Exception Handling** | try/catch/throw | Explicit error returns |

**Key Mindset Shift**: 
- Java: "Inherit behavior from parent"
- Go: "Compose small pieces, implement interfaces"

---

## 🎭 Polymorphism in Go

### Java Way: Abstract Classes + Inheritance

```java
// Java
public abstract class EventHandler {
    protected Logger logger;
    
    public EventHandler(Logger logger) {
        this.logger = logger;
    }
    
    // Template method pattern
    public void handle(Event event) {
        logger.info("Handling " + event.getType());
        processEvent(event);
        afterProcess(event);
    }
    
    protected abstract void processEvent(Event event);
    
    protected void afterProcess(Event event) {
        // Default implementation
    }
}

@Service
public class OrderFilledHandler extends EventHandler {
    @Override
    protected void processEvent(Event event) {
        OrderFilled order = (OrderFilled) event;
        // Process order...
    }
}

@Service
public class TransferHandler extends EventHandler {
    @Override
    protected void processEvent(Event event) {
        Transfer transfer = (Transfer) event;
        // Process transfer...
    }
}
```

### Go Way: Interfaces + Composition

```go
// Go - Interface-based polymorphism
package handler

import (
    "context"
    "github.com/ethereum/go-ethereum/core/types"
)

// Handler is the interface (like Java interface)
type Handler interface {
    Handle(ctx context.Context, log types.Log, timestamp uint64) (any, error)
}

// BaseHandler provides common functionality (composition, not inheritance)
type BaseHandler struct {
    logger *zap.Logger
}

func NewBaseHandler(logger *zap.Logger) BaseHandler {
    return BaseHandler{logger: logger}
}

// Common methods available to all handlers
func (h *BaseHandler) LogProcessing(eventType string) {
    h.logger.Info("Handling event", zap.String("type", eventType))
}

// OrderFilledHandler - embeds BaseHandler (composition)
type OrderFilledHandler struct {
    BaseHandler  // Embedded struct (composition, not inheritance!)
    exchange     *contracts.CTFExchange
}

func NewOrderFilledHandler(logger *zap.Logger, exchange *contracts.CTFExchange) *OrderFilledHandler {
    return &OrderFilledHandler{
        BaseHandler: NewBaseHandler(logger),
        exchange:    exchange,
    }
}

// Implements Handler interface
func (h *OrderFilledHandler) Handle(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
    h.LogProcessing("OrderFilled")  // Use embedded BaseHandler method
    
    // Parse ABI
    event, err := h.exchange.ParseOrderFilled(log)
    if err != nil {
        return nil, err
    }
    
    // Process...
    return event, nil
}

// TransferHandler - also embeds BaseHandler
type TransferHandler struct {
    BaseHandler
    ctf *contracts.ConditionalTokens
}

func NewTransferHandler(logger *zap.Logger, ctf *contracts.ConditionalTokens) *TransferHandler {
    return &TransferHandler{
        BaseHandler: NewBaseHandler(logger),
        ctf:         ctf,
    }
}

func (h *TransferHandler) Handle(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
    h.LogProcessing("TransferSingle")
    
    event, err := h.ctf.ParseTransferSingle(log)
    if err != nil {
        return nil, err
    }
    
    return event, nil
}
```

### Usage: Polymorphic Collection

```go
// Store different handlers polymorphically
handlers := map[string]Handler{
    "OrderFilled":    NewOrderFilledHandler(logger, exchange),
    "TransferSingle": NewTransferHandler(logger, ctf),
}

// Use polymorphically
for eventType, handler := range handlers {
    result, err := handler.Handle(ctx, log, timestamp)  // Polymorphic call!
    // ...
}
```

**Key Differences**:
1. **Interfaces are implicit** - No `implements` keyword
2. **Composition over inheritance** - Embed structs instead of extending
3. **No virtual methods** - All methods are "virtual" (can be overridden via embedding)

---

## 🧩 Inheritance vs Composition

### Java: Class Hierarchy

```java
// Java - Deep inheritance hierarchy
public abstract class BaseService {
    @Autowired
    protected Logger logger;
    
    @Autowired
    protected MetricsService metrics;
    
    protected void recordMetric(String name, long value) {
        metrics.record(name, value);
    }
}

public abstract class EventService extends BaseService {
    @Autowired
    protected EventRepository eventRepo;
    
    protected Event saveEvent(Event event) {
        recordMetric("events.saved", 1);
        return eventRepo.save(event);
    }
}

@Service
public class OrderService extends EventService {
    @Autowired
    private OrderRepository orderRepo;
    
    public void processOrder(Order order) {
        Event event = new Event(order);
        saveEvent(event);  // From EventService
        orderRepo.save(order);
    }
}
```

### Go: Flat Composition

```go
// Go - Composition pattern (idiomatic)
package service

// Small, focused interfaces
type Logger interface {
    Info(msg string, fields ...zap.Field)
    Error(msg string, fields ...zap.Field)
}

type MetricsRecorder interface {
    RecordMetric(name string, value int64)
}

type EventRepository interface {
    Save(ctx context.Context, event *models.Event) error
}

type OrderRepository interface {
    Save(ctx context.Context, order *models.Order) error
}

// OrderService composes dependencies (no inheritance!)
type OrderService struct {
    logger    Logger
    metrics   MetricsRecorder
    eventRepo EventRepository
    orderRepo OrderRepository
}

// Constructor (like Spring @Autowired, but explicit)
func NewOrderService(
    logger Logger,
    metrics MetricsRecorder,
    eventRepo EventRepository,
    orderRepo OrderRepository,
) *OrderService {
    return &OrderService{
        logger:    logger,
        metrics:   metrics,
        eventRepo: eventRepo,
        orderRepo: orderRepo,
    }
}

// Business logic method
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    // Use composed dependencies
    s.logger.Info("Processing order", zap.String("id", order.ID))
    
    event := models.NewEventFromOrder(order)
    
    // Record metric
    s.metrics.RecordMetric("events.saved", 1)
    
    // Save event
    if err := s.eventRepo.Save(ctx, event); err != nil {
        return fmt.Errorf("save event: %w", err)
    }
    
    // Save order
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return fmt.Errorf("save order: %w", err)
    }
    
    return nil
}
```

**Why Composition is Better in Go**:
1. **Explicit dependencies** - No hidden `@Autowired` magic
2. **Easier testing** - Inject mocks directly
3. **Better readability** - See all dependencies in struct
4. **Avoid deep hierarchies** - Flat is better than nested

---

#### 🎯 Key Question: Where did EventService go?

**Java Developer Asks:**
> "In Java, OrderService extends EventService, so I can call `saveEvent()` from OrderService.  
> In Go, there's no EventService struct. Where did it go? How do I call `saveEvent()`?"

**Answer: EventService became EventRepository interface + direct composition**

Let's break this down:

##### **Java Approach (Inheritance Chain)**

```java
// 3-level inheritance hierarchy
public abstract class BaseService {
    protected Logger logger;
    protected MetricsService metrics;
    
    protected void recordMetric(String name, long value) {
        metrics.record(name, value);  // Helper method
    }
}

public abstract class EventService extends BaseService {
    protected EventRepository eventRepo;  // ← EventService has-a EventRepository
    
    // EventService provides saveEvent() method
    protected Event saveEvent(Event event) {
        recordMetric("events.saved", 1);  // From BaseService
        return eventRepo.save(event);      // Delegates to repository
    }
}

public class OrderService extends EventService {
    private OrderRepository orderRepo;
    
    public void processOrder(Order order) {
        Event event = new Event(order);
        
        // Call inherited method from EventService
        saveEvent(event);  // ← This looks like my own method, but it's from parent
        
        orderRepo.save(order);
    }
}
```

**What's happening:**
1. `OrderService` inherits from `EventService`
2. `EventService` provides `saveEvent()` method (wraps `eventRepo.save()`)
3. `OrderService` calls `saveEvent()` as if it's its own method
4. **Inheritance is used to share the `saveEvent()` helper method**

##### **Go Approach (Direct Composition - No Middle Layer)**

```go
// NO 3-level hierarchy! Flat structure instead.

type OrderService struct {
    logger    Logger               // Composed directly
    metrics   MetricsRecorder      // Composed directly
    eventRepo EventRepository      // ← Composed directly (no EventService wrapper!)
    orderRepo OrderRepository      // Composed directly
}

func NewOrderService(
    logger Logger,
    metrics MetricsRecorder,
    eventRepo EventRepository,  // ← EventRepository is injected directly
    orderRepo OrderRepository,
) *OrderService {
    return &OrderService{
        logger:    logger,
        metrics:   metrics,
        eventRepo: eventRepo,  // Store it directly
        orderRepo: orderRepo,
    }
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    event := models.NewEventFromOrder(order)
    
    // No saveEvent() helper method - call repository directly!
    s.metrics.RecordMetric("events.saved", 1)
    if err := s.eventRepo.Save(ctx, event); err != nil {  // ← Direct call
        return fmt.Errorf("save event: %w", err)
    }
    
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return fmt.Errorf("save order: %w", err)
    }
    
    return nil
}
```

**What's happening:**
1. `OrderService` composes `EventRepository` **directly** (no EventService parent)
2. No `saveEvent()` helper method - just call `s.eventRepo.Save()` directly
3. If you want a helper, create it in `OrderService` itself
4. **Composition replaces the inheritance chain**

##### **Side-by-Side: Where Did EventService Go?**

```
┌────────────────────── JAVA ──────────────────────┐  ┌──────────────────────── GO ─────────────────────┐
│                                                   │  │                                                  │
│  BaseService                                      │  │  (No BaseService)                                │
│    ├─ logger                                      │  │                                                  │
│    ├─ metrics                                     │  │                                                  │
│    └─ recordMetric()                              │  │                                                  │
│         ↑ extends                                 │  │                                                  │
│  EventService                                     │  │  (No EventService - removed the middle layer!)   │
│    ├─ eventRepo                                   │  │                                                  │
│    └─ saveEvent() ← provides helper               │  │                                                  │
│         ↑ extends                                 │  │                                                  │
│  OrderService                                     │  │  OrderService                                    │
│    ├─ orderRepo                                   │  │    ├─ logger      ← composed directly            │
│    └─ processOrder()                              │  │    ├─ metrics     ← composed directly            │
│         calls: saveEvent()  (from parent)         │  │    ├─ eventRepo   ← composed directly (no parent)│
│                orderRepo.save() (own field)       │  │    ├─ orderRepo   ← composed directly            │
│                                                   │  │    └─ ProcessOrder()                             │
│                                                   │  │         calls: eventRepo.Save() (direct)         │
│                                                   │  │                orderRepo.Save() (direct)         │
│                                                   │  │                                                  │
│  3 levels deep                                    │  │  1 level (flat)                                  │
└───────────────────────────────────────────────────┘  └──────────────────────────────────────────────────┘
```

##### **Your Understanding is 100% Correct!**

> **Java**: OrderService **inherits** EventService → gets `saveEvent()` method for free  
> **Go**: OrderService **composes** EventRepository → calls `eventRepo.Save()` directly

**Key Insight:**
- **Java**: Created `EventService` class to provide `saveEvent()` helper method
- **Go**: Skipped the middle layer entirely - compose `EventRepository` directly
- **Why?**: In Go, you don't need inheritance to share behavior. Just compose the dependency you actually need!

##### **What if I Want a Helper Method in Go?**

If you really want a `saveEvent()` helper (like Java's EventService provides), you have **3 options**:

**Option 1: Add helper method directly to OrderService (most common)**
```go
type OrderService struct {
    eventRepo EventRepository
    orderRepo OrderRepository
}

// Helper method in OrderService itself (no parent needed)
func (s *OrderService) saveEvent(ctx context.Context, event *models.Event) error {
    return s.eventRepo.Save(ctx, event)
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    event := models.NewEventFromOrder(order)
    if err := s.saveEvent(ctx, event); err != nil {  // Use helper
        return err
    }
    return s.orderRepo.Save(ctx, order)
}
```

**Option 2: Create a separate EventService type (less common)**
```go
// EventService provides event-related helpers
type EventService struct {
    eventRepo EventRepository
}

func (es *EventService) SaveEvent(ctx context.Context, event *models.Event) error {
    return es.eventRepo.Save(ctx, event)
}

// OrderService composes EventService (composition, not inheritance!)
type OrderService struct {
    eventService *EventService  // Composed, not inherited
    orderRepo    OrderRepository
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    event := models.NewEventFromOrder(order)
    if err := s.eventService.SaveEvent(ctx, event); err != nil {  // Delegate to composed service
        return err
    }
    return s.orderRepo.Save(ctx, order)
}
```

**Option 3: Embed EventService struct (closest to inheritance)**
```go
type EventService struct {
    eventRepo EventRepository
}

func (es *EventService) SaveEvent(ctx context.Context, event *models.Event) error {
    return es.eventRepo.Save(ctx, event)
}

// OrderService embeds EventService (similar to inheritance)
type OrderService struct {
    EventService  // Embedded (promoted methods)
    orderRepo OrderRepository
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    event := models.NewEventFromOrder(order)
    // Can call SaveEvent directly due to embedding (looks like inheritance)
    if err := s.SaveEvent(ctx, event); err != nil {
        return err
    }
    return s.orderRepo.Save(ctx, order)
}
```

**Which is most idiomatic Go?**
- **Option 1** (helper in same struct) - Most common, most Go-like
- **Option 2** (separate composed service) - Good for shared logic across multiple services
- **Option 3** (embedded struct) - Works but feels Java-like; less common in Go code

##### **Summary for Java Developers**

| Aspect | Java | Go | Translation |
|--------|------|-----|-------------|
| **EventService class** | Exists as parent class | **Doesn't exist** - just EventRepository interface | Middle layer removed |
| **saveEvent() method** | In EventService parent | Option 1: In OrderService itself<br>Option 2: In composed EventService<br>Option 3: In embedded EventService | Helper methods live in the class that needs them |
| **Access to eventRepo** | `this.eventRepo` (inherited) | `s.eventRepo` (composed field) | Direct field access, not inherited |
| **Relationship** | `extends` (is-a) | Composition (has-a) | Has-a replaces is-a |
| **Why this way?** | Share saveEvent() via inheritance | Compose what you need, create helpers where needed | Favor composition over inheritance |

**You nailed it!** 🎯 In Go, we **compose EventRepository directly** instead of creating an EventService parent class. This is the Go way!

---

## 🌱 Spring Boot Patterns in Go

### 1. @Service / @Component Pattern

#### Java
```java
@Service
public class BlockchainSyncService {
    @Autowired
    private ChainClient chainClient;
    
    @Autowired
    private EventProcessor processor;
    
    @PostConstruct
    public void init() {
        // Initialization
    }
    
    public void syncBlocks(long from, long to) {
        // ...
    }
}
```

#### Go Equivalent
```go
// service/syncer.go
package service

type BlockchainSyncService struct {
    chainClient ChainClient
    processor   EventProcessor
}

// Constructor acts as factory (no @Service annotation needed)
func NewBlockchainSyncService(
    chainClient ChainClient,
    processor EventProcessor,
) *BlockchainSyncService {
    s := &BlockchainSyncService{
        chainClient: chainClient,
        processor:   processor,
    }
    
    // @PostConstruct equivalent - call init in constructor
    s.init()
    
    return s
}

func (s *BlockchainSyncService) init() {
    // Initialization logic
}

func (s *BlockchainSyncService) SyncBlocks(ctx context.Context, from, to uint64) error {
    // ...
    return nil
}
```

---

### 2. @Configuration / @Bean Pattern

#### Java
```java
@Configuration
public class AppConfig {
    @Bean
    public ChainClient chainClient(@Value("${rpc.url}") String rpcUrl) {
        return new ChainClient(rpcUrl);
    }
    
    @Bean
    public EventProcessor eventProcessor(ChainClient client, Logger logger) {
        return new EventProcessor(client, logger);
    }
    
    @Bean
    public SyncService syncService(EventProcessor processor) {
        return new SyncService(processor);
    }
}
```

#### Go Equivalent (Manual Wire)
```go
// cmd/indexer/main.go
package main

import (
    "github.com/0xkanth/polymarket-indexer/internal/chain"
    "github.com/0xkanth/polymarket-indexer/internal/processor"
    "github.com/0xkanth/polymarket-indexer/internal/syncer"
    "github.com/0xkanth/polymarket-indexer/pkg/config"
)

func main() {
    // Load config (equivalent to @Value)
    cfg, err := config.Load("config.toml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Wire dependencies manually (like @Bean methods)
    logger := setupLogger()
    
    // Bean 1: ChainClient
    chainClient, err := chain.NewClient(cfg.RPC.URL)
    if err != nil {
        log.Fatal(err)
    }
    
    // Bean 2: EventProcessor (depends on ChainClient)
    eventProcessor := processor.NewEventProcessor(chainClient, logger)
    
    // Bean 3: SyncService (depends on EventProcessor)
    syncService := syncer.NewSyncService(eventProcessor)
    
    // Run
    if err := syncService.Start(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

#### Go with Wire (Google's DI Framework)
```go
// wire.go
//go:build wireinject
// +build wireinject

package main

import (
    "github.com/google/wire"
    "github.com/0xkanth/polymarket-indexer/internal/chain"
    "github.com/0xkanth/polymarket-indexer/internal/processor"
    "github.com/0xkanth/polymarket-indexer/internal/syncer"
)

// Provider sets (like @Configuration classes)
var chainClientSet = wire.NewSet(
    chain.NewClient,
    wire.Bind(new(ChainClientInterface), new(*chain.Client)),
)

var processorSet = wire.NewSet(
    processor.NewEventProcessor,
)

var syncerSet = wire.NewSet(
    syncer.NewSyncService,
)

// Wire everything together (like ApplicationContext)
func InitializeApp(cfg *config.Config) (*App, error) {
    wire.Build(
        chainClientSet,
        processorSet,
        syncerSet,
        NewApp,
    )
    return nil, nil
}

// Generated code (wire_gen.go) will look like:
// func InitializeApp(cfg *config.Config) (*App, error) {
//     client := chain.NewClient(cfg.RPC.URL)
//     processor := processor.NewEventProcessor(client)
//     syncer := syncer.NewSyncService(processor)
//     app := NewApp(syncer)
//     return app, nil
// }
```

---

### 3. @Transactional Pattern

#### Java Version
```java
@Service
public class OrderService {
    @Autowired
    private OrderRepository orderRepository;
    
    @Autowired
    private EventRepository eventRepository;
    
    @Transactional
    public void processOrder(Order order) {
        orderRepository.save(order);
        eventRepository.save(new Event(order));
        // Automatic rollback on exception
    }
}
```

#### Go Version
```go
// service/order.go
package service

import (
    "context"
    "database/sql"
    "fmt"
)

type OrderService struct {
    db *sql.DB
}

func NewOrderService(db *sql.DB) *OrderService {
    return &OrderService{db: db}
}

func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    // 1. Begin transaction (explicit, not annotation)
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    
    // 2. Defer rollback (executes on panic or error return)
    defer tx.Rollback()
    
    // 3. Execute operations
    if err := s.saveOrder(ctx, tx, order); err != nil {
        return err  // Rollback happens automatically via defer
    }
    
    event := models.NewEventFromOrder(order)
    if err := s.saveEvent(ctx, tx, event); err != nil {
        return err  // Rollback happens automatically
    }
    
    // 4. Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit tx: %w", err)
    }
    
    return nil
}

func (s *OrderService) saveOrder(ctx context.Context, tx *sql.Tx, order *models.Order) error {
    _, err := tx.ExecContext(ctx, "INSERT INTO orders (...) VALUES (...)", order.ID, order.Amount)
    return err
}

func (s *OrderService) saveEvent(ctx context.Context, tx *sql.Tx, event *models.Event) error {
    _, err := tx.ExecContext(ctx, "INSERT INTO events (...) VALUES (...)", event.ID, event.Type)
    return err
}
```

---

#### 🤔 Java Developer Perspective: Why This Translation?

Let me walk through each part as a Java developer would think about it:

##### **1. Method Signature Differences**

**Java:**
```java
@Transactional
public void processOrder(Order order)
```

**Go:**
```go
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error
```

**Why the differences?**

| Aspect | Java | Go | Reasoning |
|--------|------|-----|-----------|
| **Return type** | `void` | `error` | Go doesn't have exceptions. **All errors must be returned explicitly**. Java throws exceptions that bubble up. |
| **Context parameter** | Hidden (ThreadLocal) | Explicit `ctx` | Go makes cancellation/timeout explicit. In Java, transaction timeout is in annotation: `@Transactional(timeout=30)`. In Go, you control it: `ctx, cancel := context.WithTimeout(ctx, 30*time.Second)` |
| **Order parameter** | `Order` | `*models.Order` | Go uses pointers for structs to avoid copying. Java always passes objects by reference (except primitives). |
| **Receiver** | Implicit `this` | Explicit `(s *OrderService)` | Go methods have explicit receiver. Like `public void processOrder(this OrderService s, ...)` if Java made `this` explicit. |

##### **2. Transaction Management**

**In Java, you think:**
```java
@Transactional  // ← Spring proxy intercepts this method
public void processOrder(Order order) {
    // Spring proxy BEFORE: begins transaction
    
    orderRepository.save(order);
    eventRepository.save(new Event(order));
    
    // Spring proxy AFTER:
    // - If no exception: commit
    // - If exception: rollback
}
```

**In Go, you must write what Spring does for you:**
```go
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    // YOU manually begin transaction (what Spring proxy does before)
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    
    // YOU set up rollback (what Spring proxy does on exception)
    defer tx.Rollback()  // ← Like try/catch/rollback in Java
    
    // Your business logic (same as Java)
    if err := s.saveOrder(ctx, tx, order); err != nil {
        return err  // ← Rollback via defer (like throw exception)
    }
    
    event := models.NewEventFromOrder(order)
    if err := s.saveEvent(ctx, tx, event); err != nil {
        return err  // ← Rollback via defer (like throw exception)
    }
    
    // YOU manually commit (what Spring proxy does on success)
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit tx: %w", err)
    }
    
    return nil
}
```

**Java Developer Mental Model:**
- **Java**: Spring AOP proxy wraps your method → handles tx lifecycle → you write business logic
- **Go**: No magic → you are the proxy → you handle tx lifecycle + business logic

##### **3. The `defer tx.Rollback()` Pattern - Understanding from Java**

**This is the hardest part for Java developers!**

**Java equivalent would be:**
```java
public void processOrder(Order order) {
    EntityTransaction tx = entityManager.getTransaction();
    tx.begin();
    
    try {
        saveOrder(tx, order);
        saveEvent(tx, new Event(order));
        tx.commit();
    } catch (Exception e) {
        tx.rollback();  // ← You must call this on error
        throw e;
    }
}
```

**Go uses `defer` to simplify this:**
```go
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    
    defer tx.Rollback()  // ← Scheduled to run when function exits
    
    // If we return error here → function exits → defer runs → rollback!
    if err := s.saveOrder(ctx, tx, order); err != nil {
        return err  // Triggers rollback via defer
    }
    
    // If we panic here → function exits → defer runs → rollback!
    event := models.NewEventFromOrder(order)
    if err := s.saveEvent(ctx, tx, event); err != nil {
        return err  // Triggers rollback via defer
    }
    
    // Success path: commit
    if err := tx.Commit(); err != nil {
        return err
    }
    
    return nil  // Even though defer runs tx.Rollback(), commit already succeeded, so rollback is a no-op
}
```

**Key Insight**: `defer tx.Rollback()` is Go's version of `finally { rollback(); }` but cleaner.

- **If commit succeeds**: Rollback is called but does nothing (tx already closed)
- **If commit fails**: Rollback properly cleans up
- **If any error before commit**: Rollback happens automatically
- **If panic**: Rollback happens automatically

##### **4. Why Pass `tx *sql.Tx` to Helper Methods?**

**Java:**
```java
@Transactional
public void processOrder(Order order) {
    orderRepository.save(order);  // ← Transaction is ThreadLocal, repo uses it automatically
    eventRepository.save(event);   // ← Same transaction, no need to pass
}
```

**Go:**
```go
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    tx, err := s.db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    // YOU must pass tx explicitly to ensure same transaction
    if err := s.saveOrder(ctx, tx, order); err != nil {
        return err
    }
    
    if err := s.saveEvent(ctx, tx, event); err != nil {
        return err
    }
    
    return tx.Commit()
}

func (s *OrderService) saveOrder(ctx context.Context, tx *sql.Tx, order *models.Order) error {
    // Uses tx, not s.db, to be part of same transaction
    _, err := tx.ExecContext(ctx, "INSERT INTO orders (...) VALUES (...)", order.ID)
    return err
}
```

**Why?**

| Java | Go | Reason |
|------|-----|---------|
| Transaction is **ThreadLocal** (hidden) | Transaction is **explicit parameter** | Go has no ThreadLocal. Everything is explicit. |
| Spring manages transaction in background | You manage transaction explicitly | No framework magic → you control lifecycle |
| All DB calls in same thread use same tx | Must pass `tx` to ensure same transaction | No implicit context → pass explicitly |

**Java Developer Analogy:**
- Java: Spring puts transaction in invisible backpack everyone can access
- Go: You must hand the transaction object to each person who needs it

##### **5. Error Handling Differences**

**Java:**
```java
@Transactional(rollbackFor = Exception.class)
public void processOrder(Order order) throws OrderException {
    orderRepository.save(order);  // If throws → Spring catches → rollback → re-throw
    eventRepository.save(event);   // If throws → Spring catches → rollback → re-throw
}
```

**Go:**
```go
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    tx, err := s.db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    // Check error, return immediately → defer triggers rollback
    if err := s.saveOrder(ctx, tx, order); err != nil {
        return fmt.Errorf("save order: %w", err)  // Wrap error with context
    }
    
    if err := s.saveEvent(ctx, tx, event); err != nil {
        return fmt.Errorf("save event: %w", err)  // Wrap error with context
    }
    
    return tx.Commit()
}
```

**Key Differences:**

| Java | Go | Java Developer Translation |
|------|-----|---------------------------|
| `throw new OrderException()` | `return fmt.Errorf("save order: %w", err)` | Returning error = throwing exception |
| Spring catches exception | Caller checks `if err != nil` | No try/catch, check errors explicitly |
| Stack trace automatic | Use `%w` to wrap errors | `%w` preserves error chain (like `initCause()`) |
| `@throws OrderException` in javadoc | Return type is `error` | Errors are values, not exceptions |

##### **6. Dependency Injection**

**Java:**
```java
@Service
public class OrderService {
    @Autowired
    private OrderRepository orderRepository;  // Spring injects
    
    @Autowired
    private EventRepository eventRepository;  // Spring injects
}
```

**Go:**
```go
type OrderService struct {
    db *sql.DB  // You inject manually via constructor
}

func NewOrderService(db *sql.DB) *OrderService {
    return &OrderService{db: db}  // Explicit injection
}
```

**In main.go (manual wiring):**
```go
db := connectDatabase(cfg)
orderService := NewOrderService(db)  // You control dependencies
```

**Java Developer Translation:**
- **Java**: Spring IoC container manages beans → automatic wiring
- **Go**: You are the container → manual wiring (or use Wire/Fx library)

---

##### **7. Visual Flow Comparison**

**Java Transaction Flow (Hidden):**
```
Your Code:        @Transactional processOrder(order)
                          ↓
Spring AOP Proxy: ┌─────────────────────┐
                  │ 1. Begin Transaction│
                  └─────────────────────┘
                          ↓
Your Code:        save(order)
Your Code:        save(event)
                          ↓
Spring AOP Proxy: ┌─────────────────────┐
                  │ 2. Commit/Rollback  │  ← Automatic based on exception
                  └─────────────────────┘
```

**Go Transaction Flow (Explicit):**
```
Your Code:        ProcessOrder(ctx, order)
                          ↓
Your Code:        ┌─────────────────────┐
                  │ 1. tx.Begin()       │  ← You do this
                  └─────────────────────┘
                          ↓
Your Code:        ┌─────────────────────┐
                  │ 2. defer Rollback() │  ← You set this up
                  └─────────────────────┘
                          ↓
Your Code:        save(tx, order)
Your Code:        save(tx, event)
                          ↓
Your Code:        ┌─────────────────────┐
                  │ 3. tx.Commit()      │  ← You do this
                  └─────────────────────┘
                          ↓
Go Runtime:       ┌─────────────────────┐
                  │ 4. defer runs       │  ← Auto rollback if commit failed
                  └─────────────────────┘
```

**Side-by-Side Execution:**

```
┌─────────────────────── JAVA ─────────────────────┐  ┌─────────────────────── GO ──────────────────────┐
│                                                   │  │                                                  │
│  @Transactional                                   │  │  func ProcessOrder(ctx, order) error {          │
│  public void processOrder(Order order) {          │  │      tx, err := db.BeginTx(ctx, nil)  ← EXPLICIT│
│                                                   │  │      if err != nil { return err }               │
│      // ← SPRING BEGINS TX (HIDDEN)               │  │                                                  │
│                                                   │  │      defer tx.Rollback()  ← YOU SET THIS UP     │
│      orderRepository.save(order);                 │  │                                                  │
│                                                   │  │      if err := saveOrder(tx, order); err != nil {│
│      // If exception here:                        │  │          return err  ← ERROR = ROLLBACK VIA DEFER│
│      //   - Spring catches                        │  │      }                                           │
│      //   - Rolls back                            │  │                                                  │
│      //   - Re-throws                             │  │      event := NewEvent(order)                    │
│                                                   │  │      if err := saveEvent(tx, event); err != nil {│
│      eventRepository.save(new Event(order));      │  │          return err  ← ERROR = ROLLBACK VIA DEFER│
│                                                   │  │      }                                           │
│      // ← SPRING COMMITS TX (HIDDEN)              │  │                                                  │
│  }                                                │  │      return tx.Commit()  ← EXPLICIT COMMIT       │
│                                                   │  │  }  ← DEFER RUNS HERE (ROLLBACK IF COMMIT FAILED)│
└───────────────────────────────────────────────────┘  └──────────────────────────────────────────────────┘
```

**Key Insight for Java Developers:**
> In Java, Spring is like an invisible butler managing transactions behind the scenes.  
> In Go, **you are the butler**. But Go gives you `defer` to make the job easier.

---

##### **8. Common Java Developer Mistakes in Go**

**❌ Mistake 1: Forgetting to check errors**
```go
// BAD (Java mindset: exceptions bubble up)
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    tx, _ := s.db.BeginTx(ctx, nil)  // Ignoring error!
    defer tx.Rollback()
    
    s.saveOrder(ctx, tx, order)  // Ignoring error!
    return tx.Commit()
}

// GOOD (Go mindset: check every error)
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)  // Handle it!
    }
    defer tx.Rollback()
    
    if err := s.saveOrder(ctx, tx, order); err != nil {
        return fmt.Errorf("save order: %w", err)  // Handle it!
    }
    
    return tx.Commit()
}
```

**❌ Mistake 2: Using `s.db` instead of `tx` in helper methods**
```go
// BAD (uses s.db, not in same transaction!)
func (s *OrderService) saveOrder(ctx context.Context, tx *sql.Tx, order *models.Order) error {
    _, err := s.db.ExecContext(ctx, "INSERT INTO orders ...")  // WRONG: uses s.db
    return err
}

// GOOD (uses tx parameter)
func (s *OrderService) saveOrder(ctx context.Context, tx *sql.Tx, order *models.Order) error {
    _, err := tx.ExecContext(ctx, "INSERT INTO orders ...")  // CORRECT: uses tx
    return err
}
```

**❌ Mistake 3: Not understanding `defer` execution order**
```go
// Defers run in LIFO (Last In, First Out) order
func example() {
    defer fmt.Println("1")  // Runs 3rd
    defer fmt.Println("2")  // Runs 2nd
    defer fmt.Println("3")  // Runs 1st
    fmt.Println("body")     // Runs first
}
// Output: body, 3, 2, 1

// For transactions, only one defer needed:
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    tx, err := s.db.BeginTx(ctx, nil)
    defer tx.Rollback()  // ← Only need this one defer
    
    // Do work...
    
    return tx.Commit()  // Commit succeeds → rollback becomes no-op
}
```

**❌ Mistake 4: Trying to create abstract base classes**
```go
// BAD (Java mindset: inheritance)
type BaseService struct {
    db *sql.DB
}

func (b *BaseService) WithTransaction(fn func(*sql.Tx) error) error {
    // transaction logic
}

type OrderService struct {
    BaseService  // Trying to inherit
}

// GOOD (Go mindset: composition)
type OrderService struct {
    db     *sql.DB
    txMgr  TransactionManager  // Compose, don't inherit
}

type TransactionManager interface {
    WithTransaction(ctx context.Context, fn func(*sql.Tx) error) error
}
```

---

**Reusable Transaction Helper**:
```go
// db/transaction.go
package db

import (
    "context"
    "database/sql"
)

// WithTransaction is a helper for @Transactional pattern
func WithTransaction(ctx context.Context, db *sql.DB, fn func(*sql.Tx) error) error {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin tx: %w", err)
    }
    
    defer func() {
        if p := recover(); p != nil {
            tx.Rollback()
            panic(p)  // Re-throw panic
        } else if err != nil {
            tx.Rollback()
        }
    }()
    
    err = fn(tx)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

// Usage
func (s *OrderService) ProcessOrder(ctx context.Context, order *models.Order) error {
    return db.WithTransaction(ctx, s.db, func(tx *sql.Tx) error {
        if err := s.saveOrder(ctx, tx, order); err != nil {
            return err
        }
        
        event := models.NewEventFromOrder(order)
        return s.saveEvent(ctx, tx, event)
    })
}
```

---

### 4. @Scheduled / @Async Pattern

#### Java
```java
@Service
public class SyncScheduler {
    @Scheduled(fixedDelay = 5000)
    public void syncBlocks() {
        // Runs every 5 seconds
    }
    
    @Async
    public CompletableFuture<Result> asyncOperation() {
        // Runs in background thread pool
        return CompletableFuture.completedFuture(result);
    }
}
```

#### Go Equivalent
```go
// service/scheduler.go
package service

import (
    "context"
    "time"
)

type SyncScheduler struct {
    syncService *SyncService
    ticker      *time.Ticker
    done        chan struct{}
}

func NewSyncScheduler(syncService *SyncService) *SyncScheduler {
    return &SyncScheduler{
        syncService: syncService,
        done:        make(chan struct{}),
    }
}

// Start scheduled task (@Scheduled equivalent)
func (s *SyncScheduler) Start(ctx context.Context) {
    s.ticker = time.NewTicker(5 * time.Second)
    
    go func() {
        for {
            select {
            case <-s.ticker.C:
                s.syncBlocks(ctx)  // Run every 5 seconds
            case <-ctx.Done():
                return
            case <-s.done:
                return
            }
        }
    }()
}

func (s *SyncScheduler) Stop() {
    if s.ticker != nil {
        s.ticker.Stop()
    }
    close(s.done)
}

func (s *SyncScheduler) syncBlocks(ctx context.Context) {
    // Scheduled task logic
    if err := s.syncService.Sync(ctx); err != nil {
        log.Error("sync failed", zap.Error(err))
    }
}

// Async operation (@Async equivalent)
func (s *SyncScheduler) AsyncOperation(ctx context.Context) <-chan Result {
    resultChan := make(chan Result, 1)
    
    // Launch goroutine (like @Async thread pool)
    go func() {
        defer close(resultChan)
        
        // Do async work
        result := s.doWork(ctx)
        resultChan <- result
    }()
    
    return resultChan
}

// Usage of async
func main() {
    scheduler := NewSyncScheduler(syncService)
    
    // Call async operation
    resultChan := scheduler.AsyncOperation(ctx)
    
    // Wait for result (like .get() on CompletableFuture)
    result := <-resultChan
    fmt.Println(result)
}
```

**Better: Use errgroup for parallel operations**
```go
import "golang.org/x/sync/errgroup"

func (s *SyncScheduler) ProcessBatchAsync(ctx context.Context, batches []Batch) error {
    g, ctx := errgroup.WithContext(ctx)
    
    // Launch parallel goroutines (like @Async pool)
    for _, batch := range batches {
        batch := batch  // Capture loop var
        g.Go(func() error {
            return s.processBatch(ctx, batch)
        })
    }
    
    // Wait for all to complete (like CompletableFuture.allOf())
    return g.Wait()
}
```

---

## 🏗️ Service Layer Patterns

### Java: Layered Architecture

```java
// Controller layer
@RestController
@RequestMapping("/api/orders")
public class OrderController {
    @Autowired
    private OrderService orderService;
    
    @PostMapping
    public ResponseEntity<Order> createOrder(@RequestBody OrderRequest req) {
        Order order = orderService.createOrder(req);
        return ResponseEntity.ok(order);
    }
}

// Service layer
@Service
public class OrderService {
    @Autowired
    private OrderRepository orderRepository;
    
    @Autowired
    private EventPublisher eventPublisher;
    
    @Transactional
    public Order createOrder(OrderRequest req) {
        Order order = new Order(req);
        order = orderRepository.save(order);
        eventPublisher.publish(new OrderCreatedEvent(order));
        return order;
    }
}

// Repository layer
@Repository
public interface OrderRepository extends JpaRepository<Order, Long> {
    List<Order> findByStatus(OrderStatus status);
}
```

### Go: Layered Architecture

```go
// Handler layer (like Controller)
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type OrderHandler struct {
    orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
    return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
    var req OrderRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    order, err := h.orderService.CreateOrder(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, order)
}

// Service layer
package service

type OrderService struct {
    orderRepo      OrderRepository
    eventPublisher EventPublisher
}

func NewOrderService(orderRepo OrderRepository, eventPublisher EventPublisher) *OrderService {
    return &OrderService{
        orderRepo:      orderRepo,
        eventPublisher: eventPublisher,
    }
}

func (s *OrderService) CreateOrder(ctx context.Context, req *OrderRequest) (*models.Order, error) {
    order := models.NewOrder(req)
    
    // Save order
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return nil, fmt.Errorf("save order: %w", err)
    }
    
    // Publish event
    event := models.NewOrderCreatedEvent(order)
    if err := s.eventPublisher.Publish(ctx, event); err != nil {
        return nil, fmt.Errorf("publish event: %w", err)
    }
    
    return order, nil
}

// Repository layer
package repository

type OrderRepository interface {
    Save(ctx context.Context, order *models.Order) error
    FindByStatus(ctx context.Context, status models.OrderStatus) ([]*models.Order, error)
}

type PostgresOrderRepository struct {
    db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
    return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Save(ctx context.Context, order *models.Order) error {
    query := `INSERT INTO orders (id, amount, status) VALUES ($1, $2, $3)`
    _, err := r.db.ExecContext(ctx, query, order.ID, order.Amount, order.Status)
    return err
}

func (r *PostgresOrderRepository) FindByStatus(ctx context.Context, status models.OrderStatus) ([]*models.Order, error) {
    query := `SELECT id, amount, status FROM orders WHERE status = $1`
    rows, err := r.db.QueryContext(ctx, query, status)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var orders []*models.Order
    for rows.Next() {
        var order models.Order
        if err := rows.Scan(&order.ID, &order.Amount, &order.Status); err != nil {
            return nil, err
        }
        orders = append(orders, &order)
    }
    
    return orders, rows.Err()
}
```

---

## 💉 Dependency Injection

### Option 1: Manual Constructor Injection (Simple)

```go
// main.go
func main() {
    // Load config
    cfg := loadConfig()
    
    // Create infrastructure
    db := connectDB(cfg.Database)
    logger := setupLogger()
    
    // Create repositories
    orderRepo := repository.NewPostgresOrderRepository(db)
    eventRepo := repository.NewEventRepository(db)
    
    // Create services
    eventPublisher := nats.NewPublisher(cfg.NATS)
    orderService := service.NewOrderService(orderRepo, eventPublisher, logger)
    
    // Create handlers
    orderHandler := handler.NewOrderHandler(orderService)
    
    // Setup routes
    router := gin.Default()
    router.POST("/orders", orderHandler.CreateOrder)
    
    router.Run(":8080")
}
```

### Option 2: Google Wire (Automatic DI)

```go
// wire.go
//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

// Provider sets
var repositorySet = wire.NewSet(
    repository.NewPostgresOrderRepository,
    repository.NewEventRepository,
    wire.Bind(new(service.OrderRepository), new(*repository.PostgresOrderRepository)),
)

var serviceSet = wire.NewSet(
    service.NewOrderService,
    nats.NewPublisher,
    wire.Bind(new(service.EventPublisher), new(*nats.Publisher)),
)

var handlerSet = wire.NewSet(
    handler.NewOrderHandler,
)

// Wire function (like Spring ApplicationContext)
func InitializeApp(cfg *config.Config, db *sql.DB, logger *zap.Logger) (*App, error) {
    wire.Build(
        repositorySet,
        serviceSet,
        handlerSet,
        NewApp,
    )
    return nil, nil
}

// Run: wire gen ./...
// Generates: wire_gen.go with all dependencies wired
```

### Option 3: Uber Fx (Runtime DI)

```go
// main.go
package main

import (
    "go.uber.org/fx"
    "github.com/0xkanth/polymarket-indexer/internal/handler"
    "github.com/0xkanth/polymarket-indexer/internal/repository"
    "github.com/0xkanth/polymarket-indexer/internal/service"
)

func main() {
    fx.New(
        // Provide dependencies (like @Bean)
        fx.Provide(
            loadConfig,
            connectDB,
            setupLogger,
            repository.NewPostgresOrderRepository,
            repository.NewEventRepository,
            nats.NewPublisher,
            service.NewOrderService,
            handler.NewOrderHandler,
            setupRouter,
        ),
        // Invoke application startup (like @PostConstruct)
        fx.Invoke(func(router *gin.Engine) {
            router.Run(":8080")
        }),
    ).Run()
}
```

---

## 🎯 Real Examples from Your Codebase

### Example 1: Event Handler Pattern (Polymorphism)

**Your Current Code** ([internal/handler/events.go](../internal/handler/events.go)):
```go
// Already using Go polymorphism correctly!

// Interface
type LogHandlerFunc func(context.Context, types.Log, uint64) (any, error)

// Different implementations
func HandleOrderFilled(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
    // ...
}

func HandleTransferSingle(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
    // ...
}

// Polymorphic storage and dispatch
handlers := map[common.Hash]LogHandlerFunc{
    OrderFilledSig:    HandleOrderFilled,
    TransferSingleSig: HandleTransferSingle,
}
```

**Java Equivalent** (what this would look like):
```java
// Java equivalent
public interface LogHandler {
    Object handle(Context ctx, Log log, long timestamp) throws Exception;
}

@Component
public class OrderFilledHandler implements LogHandler {
    @Override
    public Object handle(Context ctx, Log log, long timestamp) {
        // ...
    }
}

@Component
public class TransferSingleHandler implements LogHandler {
    @Override
    public Object handle(Context ctx, Log log, long timestamp) {
        // ...
    }
}

@Service
public class EventRouter {
    @Autowired
    private Map<Hash, LogHandler> handlers;  // Spring auto-wires
    
    public void route(Log log) {
        handlers.get(log.getTopic()).handle(ctx, log, timestamp);
    }
}
```

---

### Example 2: Service Composition Pattern

**Your Current Code** ([internal/processor/block_events_processor.go](../internal/processor/block_events_processor.go)):
```go
// Already using composition correctly!

type BlockEventsProcessor struct {
    chainClient         ChainClient
    natsEventPublisher  NATSEventPublisher
    router              *router.EventLogHandlerRouter
    contracts           []common.Address
    logger              *zap.Logger
}

// Constructor injection (like @Autowired)
func NewBlockEventsProcessor(
    chainClient ChainClient,
    natsEventPublisher NATSEventPublisher,
    cfg *config.ChainConfig,
    logger *zap.Logger,
) (*BlockEventsProcessor, error) {
    // ...
}
```

**Java Equivalent**:
```java
@Service
public class BlockEventsProcessor {
    @Autowired
    private ChainClient chainClient;
    
    @Autowired
    private NATSEventPublisher natsEventPublisher;
    
    @Autowired
    private EventRouter router;
    
    @Autowired
    private Logger logger;
    
    // Spring auto-injects dependencies
}
```

---

### Example 3: Interface Segregation

**Your Current Code** - Multiple small interfaces:
```go
// Small, focused interfaces (Interface Segregation Principle)

type ChainClient interface {
    FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
    GetLatestBlockNumber(ctx context.Context) (uint64, error)
}

type NATSEventPublisher interface {
    Publish(ctx context.Context, event models.Event) error
}

type CheckpointDB interface {
    GetOrCreateCheckpoint(ctx context.Context, serviceName string, startBlock uint64) (*Checkpoint, error)
    UpdateBlock(ctx context.Context, serviceName string, block uint64, hash string) error
}
```

**Java Equivalent**:
```java
// Small interfaces (good practice in both languages)

public interface ChainClient {
    List<Log> filterLogs(Context ctx, FilterQuery query);
    long getLatestBlockNumber(Context ctx);
}

public interface NATSEventPublisher {
    void publish(Context ctx, Event event);
}

public interface CheckpointDB {
    Checkpoint getOrCreateCheckpoint(Context ctx, String serviceName, long startBlock);
    void updateBlock(Context ctx, String serviceName, long block, String hash);
}
```

---

## 📋 Quick Reference: Pattern Translation

| Java Pattern | Go Pattern | Example |
|--------------|------------|---------|
| `class extends Parent` | Embed struct | `type Child struct { Parent }` |
| `implements Interface` | Implicit satisfaction | Just implement methods |
| `@Service` | Constructor function | `func NewService(...) *Service` |
| `@Autowired` | Constructor parameters | Explicit params |
| `@Configuration + @Bean` | `main.go` wiring | Manual or Wire/Fx |
| `@Transactional` | `tx, err := db.Begin()` | Explicit transaction |
| `@Scheduled` | `time.Ticker` | Goroutine with ticker |
| `@Async` | `go func()` | Launch goroutine |
| `CompletableFuture` | `<-chan Result` | Channel for async result |
| `try/catch` | `if err != nil` | Explicit error handling |
| `throw new Exception` | `return fmt.Errorf(...)` | Return error |
| `Optional<T>` | Pointer or zero value | `*T` or `T` with check |
| `Stream API` | Loops or slices package | No lazy streams |

---

## ✅ Best Practices for Java → Go

### 1. **Embrace Composition**
```go
// DON'T try to recreate Java inheritance
type BaseService struct {
    logger *zap.Logger
}

type OrderService struct {
    BaseService  // Avoid this pattern
}

// DO use composition explicitly
type OrderService struct {
    logger    *zap.Logger  // Compose, don't inherit
    orderRepo OrderRepository
}
```

### 2. **Prefer Small Interfaces**
```go
// DON'T create large interfaces (like Java abstract classes)
type Service interface {
    Init()
    Start()
    Stop()
    GetStatus() Status
    Configure(cfg Config)
    // ... 10 more methods
}

// DO create focused interfaces
type Startable interface {
    Start(ctx context.Context) error
}

type Stoppable interface {
    Stop() error
}

type Configurable interface {
    Configure(cfg Config) error
}
```

### 3. **Return Errors, Don't Panic**
```go
// DON'T use panic like exceptions
func ProcessOrder(order Order) {
    if order.Amount <= 0 {
        panic("invalid amount")  // BAD
    }
}

// DO return errors
func ProcessOrder(order Order) error {
    if order.Amount <= 0 {
        return errors.New("invalid amount")  // GOOD
    }
    return nil
}
```

### 4. **Use Context for Cancellation**
```go
// DON'T ignore context
func LongOperation() {
    // Long running work without cancellation
}

// DO use context
func LongOperation(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()  // Respect cancellation
        default:
            // Do work
        }
    }
}
```

### 5. **Test with Interfaces**
```go
// Service depends on interface
type OrderService struct {
    repo OrderRepository  // Interface, not concrete type
}

// Easy to mock in tests
type MockOrderRepository struct {
    SaveFunc func(ctx context.Context, order *Order) error
}

func (m *MockOrderRepository) Save(ctx context.Context, order *Order) error {
    return m.SaveFunc(ctx, order)
}

// Test
func TestOrderService(t *testing.T) {
    mockRepo := &MockOrderRepository{
        SaveFunc: func(ctx context.Context, order *Order) error {
            return nil
        },
    }
    
    service := NewOrderService(mockRepo)
    // ...
}
```

---

## 🔗 Related Resources

- **Effective Go**: https://go.dev/doc/effective_go (Must read!)
- **Uber Go Style Guide**: https://github.com/uber-go/guide
- **Google Wire**: https://github.com/google/wire (DI framework)
- **Uber Fx**: https://github.com/uber-go/fx (Runtime DI)
- **Go Patterns**: https://github.com/tmrts/go-patterns
- **From Java to Go**: https://yourbasic.org/golang/java-to-go/

---

## 📖 Learning Path for Java Developers

### Week 1: Unlearn Inheritance
- Read: "Composition over Inheritance" in Go
- Exercise: Refactor a Java class hierarchy to Go composition
- Practice: Rewrite Spring Boot service as Go service

### Week 2: Master Interfaces
- Read: Go interface semantics (implicit satisfaction)
- Exercise: Create small, focused interfaces
- Practice: Mock interfaces for testing

### Week 3: Dependency Injection
- Read: Manual DI vs Wire vs Fx
- Exercise: Wire a small app manually
- Practice: Convert Spring Boot @Configuration to Go

### Week 4: Concurrency
- Read: Goroutines vs Java threads
- Exercise: Convert `@Async` to goroutines
- Practice: Use `errgroup` for parallel operations

---

**Pro Tip**: Don't try to write Java in Go. Embrace Go's simplicity and explicitness. Your code will be better for it! 🚀
