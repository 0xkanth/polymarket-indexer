// Package handler provides event handlers for Polymarket contracts.
package handler

import (
	"context"
	"fmt"
	"math/big"

	"github.com/0xkanth/polymarket-indexer/pkg/models"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Event signatures for CTF Exchange
var (
	// OrderFilled(bytes32 indexed orderHash, address indexed maker, address indexed taker,
	//             uint256 makerAssetId, uint256 takerAssetId, uint256 makerAmountFilled,
	//             uint256 takerAmountFilled, uint256 fee)
	OrderFilledSig = common.HexToHash("0xd0a08e8c493f9c94f29311604c9de0fa40fe441d0d4d6e8b87b3e1a4cbadba5c")

	// OrderCancelled(bytes32 indexed orderHash)
	OrderCancelledSig = common.HexToHash("0x5152abf959f6564662358c2e52b702259b78bac5ee7842a0f01937e670efcc7d")

	// TokenRegistered(uint256 indexed token0, uint256 indexed token1, bytes32 indexed conditionId)
	TokenRegisteredSig = common.HexToHash("0xd0cba75e58a31a78e930fa8243a934dd8ed3c9d25f8c82e5c2bc7d0fdd1975f8")
)

// Event signatures for Conditional Tokens
var (
	// TransferSingle(address indexed operator, address indexed from, address indexed to,
	//                uint256 id, uint256 value)
	TransferSingleSig = common.HexToHash("0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62")

	// TransferBatch(address indexed operator, address indexed from, address indexed to,
	//               uint256[] ids, uint256[] values)
	TransferBatchSig = common.HexToHash("0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb")

	// ConditionPreparation(bytes32 indexed conditionId, address indexed oracle,
	//                       bytes32 indexed questionId, uint256 outcomeSlotCount)
	ConditionPreparationSig = common.HexToHash("0xcc914d01b5c6aa4ed0f1ce5d86badddf5cce7dc7740b28f5dbbc3dda0dff45b6")

	// ConditionResolution(bytes32 indexed conditionId, address indexed oracle,
	//                      bytes32 indexed questionId, uint256 outcomeSlotCount, uint256[] payoutNumerators)
	ConditionResolutionSig = common.HexToHash("0xb3574d9e77eea35b4c597c1ea75c16cb1c2cd18308085b42fc29dcf8bc8c0e3b")

	// PositionSplit(address indexed stakeholder, address collateralToken,
	//               bytes32 indexed parentCollectionId, bytes32 indexed conditionId,
	//               uint256[] partition, uint256 amount)
	PositionSplitSig = common.HexToHash("0x708228a5bb6c5c05fb64e66e1ef1fbbf4cf3ba9ec0c8fb333e8df26f7098c81d")

	// PositionsMerge(address indexed stakeholder, address collateralToken,
	//                bytes32 indexed parentCollectionId, bytes32 indexed conditionId,
	//                uint256[] partition, uint256 amount)
	PositionsMergeSig = common.HexToHash("0x5c2a65c3f6c72c9fb63c29b54c7f21e2cb10f60de87b9e42b90e7bdd76b6f26c")
)

// HandleOrderFilled processes OrderFilled events from CTF Exchange.
func HandleOrderFilled(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid OrderFilled event: expected 4 topics, got %d", len(log.Topics))
	}

	// Parse indexed parameters from topics
	orderHash := log.Topics[1].Hex()
	maker := common.BytesToAddress(log.Topics[2].Bytes()).Hex()
	taker := common.BytesToAddress(log.Topics[3].Bytes()).Hex()

	// Parse non-indexed parameters from data
	// Data contains: makerAssetId, takerAssetId, makerAmountFilled, takerAmountFilled, fee
	if len(log.Data) < 160 { // 5 * 32 bytes
		return nil, fmt.Errorf("invalid OrderFilled data length: %d", len(log.Data))
	}

	makerAssetID := new(big.Int).SetBytes(log.Data[0:32])
	takerAssetID := new(big.Int).SetBytes(log.Data[32:64])
	makerAmountFilled := new(big.Int).SetBytes(log.Data[64:96])
	takerAmountFilled := new(big.Int).SetBytes(log.Data[96:128])
	fee := new(big.Int).SetBytes(log.Data[128:160])

	return models.OrderFilled{
		OrderHash:         orderHash,
		Maker:             maker,
		Taker:             taker,
		MakerAssetID:      makerAssetID,
		TakerAssetID:      takerAssetID,
		MakerAmountFilled: makerAmountFilled,
		TakerAmountFilled: takerAmountFilled,
		Fee:               fee,
	}, nil
}

// HandleOrderCancelled processes OrderCancelled events from CTF Exchange.
func HandleOrderCancelled(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 2 {
		return nil, fmt.Errorf("invalid OrderCancelled event: expected 2 topics, got %d", len(log.Topics))
	}

	orderHash := log.Topics[1].Hex()

	return models.OrderCancelled{
		OrderHash: orderHash,
	}, nil
}

// HandleTokenRegistered processes TokenRegistered events from CTF Exchange.
func HandleTokenRegistered(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid TokenRegistered event: expected 4 topics, got %d", len(log.Topics))
	}

	token0 := new(big.Int).SetBytes(log.Topics[1].Bytes())
	token1 := new(big.Int).SetBytes(log.Topics[2].Bytes())
	conditionID := log.Topics[3].Hex()

	return models.TokenRegistered{
		Token0:      token0,
		Token1:      token1,
		ConditionID: conditionID,
	}, nil
}

// HandleTransferSingle processes TransferSingle events from Conditional Tokens.
func HandleTransferSingle(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid TransferSingle event: expected 4 topics, got %d", len(log.Topics))
	}

	operator := common.BytesToAddress(log.Topics[1].Bytes()).Hex()
	from := common.BytesToAddress(log.Topics[2].Bytes()).Hex()
	to := common.BytesToAddress(log.Topics[3].Bytes()).Hex()

	// Parse data: id and value
	if len(log.Data) < 64 {
		return nil, fmt.Errorf("invalid TransferSingle data length: %d", len(log.Data))
	}

	tokenID := new(big.Int).SetBytes(log.Data[0:32])
	amount := new(big.Int).SetBytes(log.Data[32:64])

	return models.TransferSingle{
		Operator: operator,
		From:     from,
		To:       to,
		TokenID:  tokenID,
		Amount:   amount,
	}, nil
}

// HandleTransferBatch processes TransferBatch events from Conditional Tokens.
func HandleTransferBatch(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid TransferBatch event: expected 4 topics, got %d", len(log.Topics))
	}

	operator := common.BytesToAddress(log.Topics[1].Bytes()).Hex()
	from := common.BytesToAddress(log.Topics[2].Bytes()).Hex()
	to := common.BytesToAddress(log.Topics[3].Bytes()).Hex()

	// Parse data: dynamic arrays of ids and values
	// We need to decode the ABI-encoded dynamic arrays
	uint256ArrayTy, _ := abi.NewType("uint256[]", "", nil)
	args := abi.Arguments{
		{Type: uint256ArrayTy},
		{Type: uint256ArrayTy},
	}

	unpacked, err := args.Unpack(log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack TransferBatch data: %w", err)
	}

	tokenIDs := unpacked[0].([]*big.Int)
	amounts := unpacked[1].([]*big.Int)

	return models.TransferBatch{
		Operator: operator,
		From:     from,
		To:       to,
		TokenIDs: tokenIDs,
		Amounts:  amounts,
	}, nil
}

// HandleConditionPreparation processes ConditionPreparation events.
func HandleConditionPreparation(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid ConditionPreparation event: expected 4 topics, got %d", len(log.Topics))
	}

	conditionID := log.Topics[1].Hex()
	oracle := common.BytesToAddress(log.Topics[2].Bytes()).Hex()
	questionID := log.Topics[3].Hex()

	// Parse outcomeSlotCount from data
	if len(log.Data) < 32 {
		return nil, fmt.Errorf("invalid ConditionPreparation data length: %d", len(log.Data))
	}

	outcomeSlotCount := uint8(new(big.Int).SetBytes(log.Data[0:32]).Uint64())

	return models.ConditionPreparation{
		ConditionID:      conditionID,
		Oracle:           oracle,
		QuestionID:       questionID,
		OutcomeSlotCount: outcomeSlotCount,
	}, nil
}

// HandleConditionResolution processes ConditionResolution events.
func HandleConditionResolution(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid ConditionResolution event: expected 4 topics, got %d", len(log.Topics))
	}

	conditionID := log.Topics[1].Hex()
	oracle := common.BytesToAddress(log.Topics[2].Bytes()).Hex()
	questionID := log.Topics[3].Hex()

	// Parse data: outcomeSlotCount and payoutNumerators array
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	uint256ArrayTy, _ := abi.NewType("uint256[]", "", nil)
	args := abi.Arguments{
		{Type: uint256Ty},      // outcomeSlotCount
		{Type: uint256ArrayTy}, // payoutNumerators
	}

	unpacked, err := args.Unpack(log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack ConditionResolution data: %w", err)
	}

	outcomeSlotCount := uint8(unpacked[0].(*big.Int).Uint64())
	payoutNumerators := unpacked[1].([]*big.Int)

	return models.ConditionResolution{
		ConditionID:      conditionID,
		Oracle:           oracle,
		QuestionID:       questionID,
		OutcomeSlotCount: outcomeSlotCount,
		PayoutNumerators: payoutNumerators,
	}, nil
}

// HandlePositionSplit processes PositionSplit events.
func HandlePositionSplit(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid PositionSplit event: expected 4 topics, got %d", len(log.Topics))
	}

	stakeholder := common.BytesToAddress(log.Topics[1].Bytes()).Hex()
	parentCollectionID := log.Topics[2].Hex()
	conditionID := log.Topics[3].Hex()

	// Parse data: collateralToken, partition array, amount
	addressTy, _ := abi.NewType("address", "", nil)
	uint256ArrayTy, _ := abi.NewType("uint256[]", "", nil)
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	args := abi.Arguments{
		{Type: addressTy},      // collateralToken
		{Type: uint256ArrayTy}, // partition
		{Type: uint256Ty},      // amount
	}

	unpacked, err := args.Unpack(log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack PositionSplit data: %w", err)
	}

	collateralToken := unpacked[0].(common.Address).Hex()
	partition := unpacked[1].([]*big.Int)
	amount := unpacked[2].(*big.Int)

	return models.PositionSplit{
		Stakeholder:        stakeholder,
		CollateralToken:    collateralToken,
		ParentCollectionID: parentCollectionID,
		ConditionID:        conditionID,
		Partition:          partition,
		Amount:             amount,
	}, nil
}

// HandlePositionsMerge processes PositionsMerge events.
func HandlePositionsMerge(ctx context.Context, log types.Log, timestamp uint64) (any, error) {
	if len(log.Topics) != 4 {
		return nil, fmt.Errorf("invalid PositionsMerge event: expected 4 topics, got %d", len(log.Topics))
	}

	stakeholder := common.BytesToAddress(log.Topics[1].Bytes()).Hex()
	parentCollectionID := log.Topics[2].Hex()
	conditionID := log.Topics[3].Hex()

	// Parse data: collateralToken, partition array, amount
	addressTy, _ := abi.NewType("address", "", nil)
	uint256ArrayTy, _ := abi.NewType("uint256[]", "", nil)
	uint256Ty, _ := abi.NewType("uint256", "", nil)
	args := abi.Arguments{
		{Type: addressTy},      // collateralToken
		{Type: uint256ArrayTy}, // partition
		{Type: uint256Ty},      // amount
	}

	unpacked, err := args.Unpack(log.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack PositionsMerge data: %w", err)
	}

	collateralToken := unpacked[0].(common.Address).Hex()
	partition := unpacked[1].([]*big.Int)
	amount := unpacked[2].(*big.Int)

	return models.PositionsMerge{
		Stakeholder:        stakeholder,
		CollateralToken:    collateralToken,
		ParentCollectionID: parentCollectionID,
		ConditionID:        conditionID,
		Partition:          partition,
		Amount:             amount,
	}, nil
}
