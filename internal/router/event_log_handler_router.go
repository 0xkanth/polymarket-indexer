// Package router provides event routing functionality.
package router

import (
	"context"
	"fmt"

	"github.com/0xkanth/polymarket-indexer/pkg/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// EventCallback is called after an event is processed by a handler.
type EventCallback func(context.Context, models.Event) error

// LogHandlerFunc processes a log event and returns the parsed payload.
type LogHandlerFunc func(context.Context, types.Log, uint64) (any, error)

// EventLogHandlerRouter routes blockchain events to their respective handlers.
type EventLogHandlerRouter struct {
	callback    EventCallback
	logHandlers map[common.Hash]LogHandlerFunc
	eventNames  map[common.Hash]string
}

// New creates a new event router with the specified callback.
func New(callback EventCallback) *EventLogHandlerRouter {
	return &EventLogHandlerRouter{
		callback:    callback,
		logHandlers: make(map[common.Hash]LogHandlerFunc),
		eventNames:  make(map[common.Hash]string),
	}
}

// RegisterLogHandler registers a handler for a specific event signature.
func (r *EventLogHandlerRouter) RegisterLogHandler(eventSignature common.Hash, eventName string, handler LogHandlerFunc) {
	r.logHandlers[eventSignature] = handler
	r.eventNames[eventSignature] = eventName
}

// RouteLog routes a log event to its registered handler.
func (r *EventLogHandlerRouter) RouteLog(ctx context.Context, log types.Log, blockTimestamp uint64, blockHash string) error {
	// Check if we have a handler for this event signature
	if len(log.Topics) == 0 {
		return nil // Skip logs without topics
	}

	eventSig := log.Topics[0]
	handler, exists := r.logHandlers[eventSig]
	if !exists {
		return nil // No handler registered, skip
	}

	// Execute handler to parse the event
	payload, err := handler(ctx, log, blockTimestamp)
	if err != nil {
		return fmt.Errorf("handler failed for event %s: %w", eventSig.Hex(), err)
	}

	// Create the event model
	event := models.Event{
		Block:        log.BlockNumber,
		BlockHash:    blockHash,
		TxHash:       log.TxHash.Hex(),
		TxIndex:      log.TxIndex,
		LogIndex:     log.Index,
		ContractAddr: log.Address.Hex(),
		EventName:    r.eventNames[eventSig],
		EventSig:     eventSig.Hex(),
		Timestamp:    blockTimestamp,
		Success:      !log.Removed, // Removed logs are from reorged blocks
		Payload:      payload,
	}

	// Call the callback (typically NATS publish)
	return r.callback(ctx, event)
}

// RouteLogs routes multiple logs from a receipt.
func (r *EventLogHandlerRouter) RouteLogs(ctx context.Context, logs []types.Log, blockTimestamp uint64, blockHash string) error {
	for _, log := range logs {
		if err := r.RouteLog(ctx, log, blockTimestamp, blockHash); err != nil {
			return err
		}
	}
	return nil
}

// HasHandler checks if a handler is registered for the given event signature.
func (r *EventLogHandlerRouter) HasHandler(eventSignature common.Hash) bool {
	_, exists := r.logHandlers[eventSignature]
	return exists
}

// HandlerCount returns the number of registered handlers.
func (r *EventLogHandlerRouter) HandlerCount() int {
	return len(r.logHandlers)
}
