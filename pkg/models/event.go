// Package models defines common data structures for events and domain models.
package models

import (
	"math/big"
	"time"
)

// Event represents a generic blockchain event with common fields.
type Event struct {
	Block        uint64    `json:"block"`
	BlockHash    string    `json:"block_hash"`
	TxHash       string    `json:"tx_hash"`
	TxIndex      uint      `json:"tx_index"`
	LogIndex     uint      `json:"log_index"`
	ContractAddr string    `json:"contract_address"`
	EventName    string    `json:"event_name"`
	EventSig     string    `json:"event_signature"`
	Timestamp    uint64    `json:"timestamp"`
	Success      bool      `json:"success"`
	Payload      any       `json:"payload"`
	ProcessedAt  time.Time `json:"processed_at"`
}

// OrderFilled represents a CTF Exchange OrderFilled event.
type OrderFilled struct {
	OrderHash         string   `json:"order_hash"`
	Maker             string   `json:"maker"`
	Taker             string   `json:"taker"`
	MakerAssetID      *big.Int `json:"maker_asset_id"`
	TakerAssetID      *big.Int `json:"taker_asset_id"`
	MakerAmountFilled *big.Int `json:"maker_amount_filled"`
	TakerAmountFilled *big.Int `json:"taker_amount_filled"`
	Fee               *big.Int `json:"fee"`
}

// OrderCancelled represents a CTF Exchange OrderCancelled event.
type OrderCancelled struct {
	OrderHash string `json:"order_hash"`
}

// TokenRegistered represents a CTF Exchange TokenRegistered event.
type TokenRegistered struct {
	Token0      *big.Int `json:"token0"`
	Token1      *big.Int `json:"token1"`
	ConditionID string   `json:"condition_id"`
}

// OrdersMatched represents a CTF Exchange OrdersMatched event.
type OrdersMatched struct {
	TakerOrderHash   string     `json:"taker_order_hash"`
	MakerAddresses   []string   `json:"maker_addresses"`
	MakerOrderHashes []*big.Int `json:"maker_order_hashes"`
	TakerFillAmount  *big.Int   `json:"taker_fill_amount"`
}

// TransferSingle represents a Conditional Tokens TransferSingle event.
type TransferSingle struct {
	Operator string   `json:"operator"`
	From     string   `json:"from"`
	To       string   `json:"to"`
	TokenID  *big.Int `json:"token_id"`
	Amount   *big.Int `json:"amount"`
}

// TransferBatch represents a Conditional Tokens TransferBatch event.
type TransferBatch struct {
	Operator string     `json:"operator"`
	From     string     `json:"from"`
	To       string     `json:"to"`
	TokenIDs []*big.Int `json:"token_ids"`
	Amounts  []*big.Int `json:"amounts"`
}

// ConditionPreparation represents a new condition/market being created.
type ConditionPreparation struct {
	ConditionID      string `json:"condition_id"`
	Oracle           string `json:"oracle"`
	QuestionID       string `json:"question_id"`
	OutcomeSlotCount uint8  `json:"outcome_slot_count"`
}

// ConditionResolution represents a market being resolved.
type ConditionResolution struct {
	ConditionID      string     `json:"condition_id"`
	Oracle           string     `json:"oracle"`
	QuestionID       string     `json:"question_id"`
	OutcomeSlotCount uint8      `json:"outcome_slot_count"`
	PayoutNumerators []*big.Int `json:"payout_numerators"`
}

// PositionSplit represents minting of conditional tokens.
type PositionSplit struct {
	Stakeholder        string     `json:"stakeholder"`
	CollateralToken    string     `json:"collateral_token"`
	ParentCollectionID string     `json:"parent_collection_id"`
	ConditionID        string     `json:"condition_id"`
	Partition          []*big.Int `json:"partition"`
	Amount             *big.Int   `json:"amount"`
}

// PositionsMerge represents redemption of conditional tokens.
type PositionsMerge struct {
	Stakeholder        string     `json:"stakeholder"`
	CollateralToken    string     `json:"collateral_token"`
	ParentCollectionID string     `json:"parent_collection_id"`
	ConditionID        string     `json:"condition_id"`
	Partition          []*big.Int `json:"partition"`
	Amount             *big.Int   `json:"amount"`
}

// Checkpoint represents the indexer's processing state.
type Checkpoint struct {
	ServiceName   string    `json:"service_name"`
	LastBlock     uint64    `json:"last_block"`
	LastBlockHash string    `json:"last_block_hash"`
	UpdatedAt     time.Time `json:"updated_at"`
}
