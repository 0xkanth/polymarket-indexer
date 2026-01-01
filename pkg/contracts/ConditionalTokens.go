// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ConditionalTokensMetaData contains all meta data concerning the ConditionalTokens contract.
var ConditionalTokensMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"conditionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outcomeSlotCount\",\"type\":\"uint256\"}],\"name\":\"ConditionPreparation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"conditionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"questionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"outcomeSlotCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"payoutNumerators\",\"type\":\"uint256[]\"}],\"name\":\"ConditionResolution\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"stakeholder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"collateralToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"parentCollectionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"conditionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"partition\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PositionSplit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"stakeholder\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"contractIERC20\",\"name\":\"collateralToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"parentCollectionId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"conditionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"partition\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PositionsMerge\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"contractIERC20\",\"name\":\"collateralToken\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"parentCollectionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"conditionId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"indexSets\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"payout\",\"type\":\"uint256\"}],\"name\":\"PayoutRedemption\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]",
}

// ConditionalTokensABI is the input ABI used to generate the binding from.
// Deprecated: Use ConditionalTokensMetaData.ABI instead.
var ConditionalTokensABI = ConditionalTokensMetaData.ABI

// ConditionalTokens is an auto generated Go binding around an Ethereum contract.
type ConditionalTokens struct {
	ConditionalTokensCaller     // Read-only binding to the contract
	ConditionalTokensTransactor // Write-only binding to the contract
	ConditionalTokensFilterer   // Log filterer for contract events
}

// ConditionalTokensCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConditionalTokensCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConditionalTokensTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConditionalTokensTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConditionalTokensFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConditionalTokensFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConditionalTokensSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConditionalTokensSession struct {
	Contract     *ConditionalTokens // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ConditionalTokensCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConditionalTokensCallerSession struct {
	Contract *ConditionalTokensCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ConditionalTokensTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConditionalTokensTransactorSession struct {
	Contract     *ConditionalTokensTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ConditionalTokensRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConditionalTokensRaw struct {
	Contract *ConditionalTokens // Generic contract binding to access the raw methods on
}

// ConditionalTokensCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConditionalTokensCallerRaw struct {
	Contract *ConditionalTokensCaller // Generic read-only contract binding to access the raw methods on
}

// ConditionalTokensTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConditionalTokensTransactorRaw struct {
	Contract *ConditionalTokensTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConditionalTokens creates a new instance of ConditionalTokens, bound to a specific deployed contract.
func NewConditionalTokens(address common.Address, backend bind.ContractBackend) (*ConditionalTokens, error) {
	contract, err := bindConditionalTokens(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokens{ConditionalTokensCaller: ConditionalTokensCaller{contract: contract}, ConditionalTokensTransactor: ConditionalTokensTransactor{contract: contract}, ConditionalTokensFilterer: ConditionalTokensFilterer{contract: contract}}, nil
}

// NewConditionalTokensCaller creates a new read-only instance of ConditionalTokens, bound to a specific deployed contract.
func NewConditionalTokensCaller(address common.Address, caller bind.ContractCaller) (*ConditionalTokensCaller, error) {
	contract, err := bindConditionalTokens(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensCaller{contract: contract}, nil
}

// NewConditionalTokensTransactor creates a new write-only instance of ConditionalTokens, bound to a specific deployed contract.
func NewConditionalTokensTransactor(address common.Address, transactor bind.ContractTransactor) (*ConditionalTokensTransactor, error) {
	contract, err := bindConditionalTokens(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensTransactor{contract: contract}, nil
}

// NewConditionalTokensFilterer creates a new log filterer instance of ConditionalTokens, bound to a specific deployed contract.
func NewConditionalTokensFilterer(address common.Address, filterer bind.ContractFilterer) (*ConditionalTokensFilterer, error) {
	contract, err := bindConditionalTokens(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensFilterer{contract: contract}, nil
}

// bindConditionalTokens binds a generic wrapper to an already deployed contract.
func bindConditionalTokens(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConditionalTokensMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConditionalTokens *ConditionalTokensRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConditionalTokens.Contract.ConditionalTokensCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConditionalTokens *ConditionalTokensRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConditionalTokens.Contract.ConditionalTokensTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConditionalTokens *ConditionalTokensRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConditionalTokens.Contract.ConditionalTokensTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConditionalTokens *ConditionalTokensCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConditionalTokens.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConditionalTokens *ConditionalTokensTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConditionalTokens.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConditionalTokens *ConditionalTokensTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConditionalTokens.Contract.contract.Transact(opts, method, params...)
}

// ConditionalTokensApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ConditionalTokens contract.
type ConditionalTokensApprovalForAllIterator struct {
	Event *ConditionalTokensApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensApprovalForAll represents a ApprovalForAll event raised by the ConditionalTokens contract.
type ConditionalTokensApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ConditionalTokensApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensApprovalForAllIterator{contract: _ConditionalTokens.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ConditionalTokensApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensApprovalForAll)
				if err := _ConditionalTokens.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ConditionalTokens *ConditionalTokensFilterer) ParseApprovalForAll(log types.Log) (*ConditionalTokensApprovalForAll, error) {
	event := new(ConditionalTokensApprovalForAll)
	if err := _ConditionalTokens.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensConditionPreparationIterator is returned from FilterConditionPreparation and is used to iterate over the raw logs and unpacked data for ConditionPreparation events raised by the ConditionalTokens contract.
type ConditionalTokensConditionPreparationIterator struct {
	Event *ConditionalTokensConditionPreparation // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensConditionPreparationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensConditionPreparation)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensConditionPreparation)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensConditionPreparationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensConditionPreparationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensConditionPreparation represents a ConditionPreparation event raised by the ConditionalTokens contract.
type ConditionalTokensConditionPreparation struct {
	ConditionId      [32]byte
	Oracle           common.Address
	QuestionId       [32]byte
	OutcomeSlotCount *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterConditionPreparation is a free log retrieval operation binding the contract event 0xab3760c3bd2bb38b5bcf54dc79802ed67338b4cf29f3054ded67ed24661e4177.
//
// Solidity: event ConditionPreparation(bytes32 indexed conditionId, address indexed oracle, bytes32 indexed questionId, uint256 outcomeSlotCount)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterConditionPreparation(opts *bind.FilterOpts, conditionId [][32]byte, oracle []common.Address, questionId [][32]byte) (*ConditionalTokensConditionPreparationIterator, error) {

	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var questionIdRule []interface{}
	for _, questionIdItem := range questionId {
		questionIdRule = append(questionIdRule, questionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "ConditionPreparation", conditionIdRule, oracleRule, questionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensConditionPreparationIterator{contract: _ConditionalTokens.contract, event: "ConditionPreparation", logs: logs, sub: sub}, nil
}

// WatchConditionPreparation is a free log subscription operation binding the contract event 0xab3760c3bd2bb38b5bcf54dc79802ed67338b4cf29f3054ded67ed24661e4177.
//
// Solidity: event ConditionPreparation(bytes32 indexed conditionId, address indexed oracle, bytes32 indexed questionId, uint256 outcomeSlotCount)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchConditionPreparation(opts *bind.WatchOpts, sink chan<- *ConditionalTokensConditionPreparation, conditionId [][32]byte, oracle []common.Address, questionId [][32]byte) (event.Subscription, error) {

	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var questionIdRule []interface{}
	for _, questionIdItem := range questionId {
		questionIdRule = append(questionIdRule, questionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "ConditionPreparation", conditionIdRule, oracleRule, questionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensConditionPreparation)
				if err := _ConditionalTokens.contract.UnpackLog(event, "ConditionPreparation", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConditionPreparation is a log parse operation binding the contract event 0xab3760c3bd2bb38b5bcf54dc79802ed67338b4cf29f3054ded67ed24661e4177.
//
// Solidity: event ConditionPreparation(bytes32 indexed conditionId, address indexed oracle, bytes32 indexed questionId, uint256 outcomeSlotCount)
func (_ConditionalTokens *ConditionalTokensFilterer) ParseConditionPreparation(log types.Log) (*ConditionalTokensConditionPreparation, error) {
	event := new(ConditionalTokensConditionPreparation)
	if err := _ConditionalTokens.contract.UnpackLog(event, "ConditionPreparation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensConditionResolutionIterator is returned from FilterConditionResolution and is used to iterate over the raw logs and unpacked data for ConditionResolution events raised by the ConditionalTokens contract.
type ConditionalTokensConditionResolutionIterator struct {
	Event *ConditionalTokensConditionResolution // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensConditionResolutionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensConditionResolution)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensConditionResolution)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensConditionResolutionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensConditionResolutionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensConditionResolution represents a ConditionResolution event raised by the ConditionalTokens contract.
type ConditionalTokensConditionResolution struct {
	ConditionId      [32]byte
	Oracle           common.Address
	QuestionId       [32]byte
	OutcomeSlotCount *big.Int
	PayoutNumerators []*big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterConditionResolution is a free log retrieval operation binding the contract event 0xb44d84d3289691f71497564b85d4233648d9dbae8cbdbb4329f301c3a0185894.
//
// Solidity: event ConditionResolution(bytes32 indexed conditionId, address indexed oracle, bytes32 indexed questionId, uint256 outcomeSlotCount, uint256[] payoutNumerators)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterConditionResolution(opts *bind.FilterOpts, conditionId [][32]byte, oracle []common.Address, questionId [][32]byte) (*ConditionalTokensConditionResolutionIterator, error) {

	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var questionIdRule []interface{}
	for _, questionIdItem := range questionId {
		questionIdRule = append(questionIdRule, questionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "ConditionResolution", conditionIdRule, oracleRule, questionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensConditionResolutionIterator{contract: _ConditionalTokens.contract, event: "ConditionResolution", logs: logs, sub: sub}, nil
}

// WatchConditionResolution is a free log subscription operation binding the contract event 0xb44d84d3289691f71497564b85d4233648d9dbae8cbdbb4329f301c3a0185894.
//
// Solidity: event ConditionResolution(bytes32 indexed conditionId, address indexed oracle, bytes32 indexed questionId, uint256 outcomeSlotCount, uint256[] payoutNumerators)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchConditionResolution(opts *bind.WatchOpts, sink chan<- *ConditionalTokensConditionResolution, conditionId [][32]byte, oracle []common.Address, questionId [][32]byte) (event.Subscription, error) {

	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}
	var questionIdRule []interface{}
	for _, questionIdItem := range questionId {
		questionIdRule = append(questionIdRule, questionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "ConditionResolution", conditionIdRule, oracleRule, questionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensConditionResolution)
				if err := _ConditionalTokens.contract.UnpackLog(event, "ConditionResolution", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConditionResolution is a log parse operation binding the contract event 0xb44d84d3289691f71497564b85d4233648d9dbae8cbdbb4329f301c3a0185894.
//
// Solidity: event ConditionResolution(bytes32 indexed conditionId, address indexed oracle, bytes32 indexed questionId, uint256 outcomeSlotCount, uint256[] payoutNumerators)
func (_ConditionalTokens *ConditionalTokensFilterer) ParseConditionResolution(log types.Log) (*ConditionalTokensConditionResolution, error) {
	event := new(ConditionalTokensConditionResolution)
	if err := _ConditionalTokens.contract.UnpackLog(event, "ConditionResolution", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensPayoutRedemptionIterator is returned from FilterPayoutRedemption and is used to iterate over the raw logs and unpacked data for PayoutRedemption events raised by the ConditionalTokens contract.
type ConditionalTokensPayoutRedemptionIterator struct {
	Event *ConditionalTokensPayoutRedemption // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensPayoutRedemptionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensPayoutRedemption)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensPayoutRedemption)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensPayoutRedemptionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensPayoutRedemptionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensPayoutRedemption represents a PayoutRedemption event raised by the ConditionalTokens contract.
type ConditionalTokensPayoutRedemption struct {
	Redeemer           common.Address
	CollateralToken    common.Address
	ParentCollectionId [32]byte
	ConditionId        [32]byte
	IndexSets          []*big.Int
	Payout             *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterPayoutRedemption is a free log retrieval operation binding the contract event 0x2682012a4a4f1973119f1c9b90745d1bd91fa2bab387344f044cb3586864d18d.
//
// Solidity: event PayoutRedemption(address indexed redeemer, address indexed collateralToken, bytes32 indexed parentCollectionId, bytes32 conditionId, uint256[] indexSets, uint256 payout)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterPayoutRedemption(opts *bind.FilterOpts, redeemer []common.Address, collateralToken []common.Address, parentCollectionId [][32]byte) (*ConditionalTokensPayoutRedemptionIterator, error) {

	var redeemerRule []interface{}
	for _, redeemerItem := range redeemer {
		redeemerRule = append(redeemerRule, redeemerItem)
	}
	var collateralTokenRule []interface{}
	for _, collateralTokenItem := range collateralToken {
		collateralTokenRule = append(collateralTokenRule, collateralTokenItem)
	}
	var parentCollectionIdRule []interface{}
	for _, parentCollectionIdItem := range parentCollectionId {
		parentCollectionIdRule = append(parentCollectionIdRule, parentCollectionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "PayoutRedemption", redeemerRule, collateralTokenRule, parentCollectionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensPayoutRedemptionIterator{contract: _ConditionalTokens.contract, event: "PayoutRedemption", logs: logs, sub: sub}, nil
}

// WatchPayoutRedemption is a free log subscription operation binding the contract event 0x2682012a4a4f1973119f1c9b90745d1bd91fa2bab387344f044cb3586864d18d.
//
// Solidity: event PayoutRedemption(address indexed redeemer, address indexed collateralToken, bytes32 indexed parentCollectionId, bytes32 conditionId, uint256[] indexSets, uint256 payout)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchPayoutRedemption(opts *bind.WatchOpts, sink chan<- *ConditionalTokensPayoutRedemption, redeemer []common.Address, collateralToken []common.Address, parentCollectionId [][32]byte) (event.Subscription, error) {

	var redeemerRule []interface{}
	for _, redeemerItem := range redeemer {
		redeemerRule = append(redeemerRule, redeemerItem)
	}
	var collateralTokenRule []interface{}
	for _, collateralTokenItem := range collateralToken {
		collateralTokenRule = append(collateralTokenRule, collateralTokenItem)
	}
	var parentCollectionIdRule []interface{}
	for _, parentCollectionIdItem := range parentCollectionId {
		parentCollectionIdRule = append(parentCollectionIdRule, parentCollectionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "PayoutRedemption", redeemerRule, collateralTokenRule, parentCollectionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensPayoutRedemption)
				if err := _ConditionalTokens.contract.UnpackLog(event, "PayoutRedemption", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePayoutRedemption is a log parse operation binding the contract event 0x2682012a4a4f1973119f1c9b90745d1bd91fa2bab387344f044cb3586864d18d.
//
// Solidity: event PayoutRedemption(address indexed redeemer, address indexed collateralToken, bytes32 indexed parentCollectionId, bytes32 conditionId, uint256[] indexSets, uint256 payout)
func (_ConditionalTokens *ConditionalTokensFilterer) ParsePayoutRedemption(log types.Log) (*ConditionalTokensPayoutRedemption, error) {
	event := new(ConditionalTokensPayoutRedemption)
	if err := _ConditionalTokens.contract.UnpackLog(event, "PayoutRedemption", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensPositionSplitIterator is returned from FilterPositionSplit and is used to iterate over the raw logs and unpacked data for PositionSplit events raised by the ConditionalTokens contract.
type ConditionalTokensPositionSplitIterator struct {
	Event *ConditionalTokensPositionSplit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensPositionSplitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensPositionSplit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensPositionSplit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensPositionSplitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensPositionSplitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensPositionSplit represents a PositionSplit event raised by the ConditionalTokens contract.
type ConditionalTokensPositionSplit struct {
	Stakeholder        common.Address
	CollateralToken    common.Address
	ParentCollectionId [32]byte
	ConditionId        [32]byte
	Partition          []*big.Int
	Amount             *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterPositionSplit is a free log retrieval operation binding the contract event 0x2e6bb91f8cbcda0c93623c54d0403a43514fabc40084ec96b6d5379a74786298.
//
// Solidity: event PositionSplit(address indexed stakeholder, address collateralToken, bytes32 indexed parentCollectionId, bytes32 indexed conditionId, uint256[] partition, uint256 amount)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterPositionSplit(opts *bind.FilterOpts, stakeholder []common.Address, parentCollectionId [][32]byte, conditionId [][32]byte) (*ConditionalTokensPositionSplitIterator, error) {

	var stakeholderRule []interface{}
	for _, stakeholderItem := range stakeholder {
		stakeholderRule = append(stakeholderRule, stakeholderItem)
	}

	var parentCollectionIdRule []interface{}
	for _, parentCollectionIdItem := range parentCollectionId {
		parentCollectionIdRule = append(parentCollectionIdRule, parentCollectionIdItem)
	}
	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "PositionSplit", stakeholderRule, parentCollectionIdRule, conditionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensPositionSplitIterator{contract: _ConditionalTokens.contract, event: "PositionSplit", logs: logs, sub: sub}, nil
}

// WatchPositionSplit is a free log subscription operation binding the contract event 0x2e6bb91f8cbcda0c93623c54d0403a43514fabc40084ec96b6d5379a74786298.
//
// Solidity: event PositionSplit(address indexed stakeholder, address collateralToken, bytes32 indexed parentCollectionId, bytes32 indexed conditionId, uint256[] partition, uint256 amount)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchPositionSplit(opts *bind.WatchOpts, sink chan<- *ConditionalTokensPositionSplit, stakeholder []common.Address, parentCollectionId [][32]byte, conditionId [][32]byte) (event.Subscription, error) {

	var stakeholderRule []interface{}
	for _, stakeholderItem := range stakeholder {
		stakeholderRule = append(stakeholderRule, stakeholderItem)
	}

	var parentCollectionIdRule []interface{}
	for _, parentCollectionIdItem := range parentCollectionId {
		parentCollectionIdRule = append(parentCollectionIdRule, parentCollectionIdItem)
	}
	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "PositionSplit", stakeholderRule, parentCollectionIdRule, conditionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensPositionSplit)
				if err := _ConditionalTokens.contract.UnpackLog(event, "PositionSplit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePositionSplit is a log parse operation binding the contract event 0x2e6bb91f8cbcda0c93623c54d0403a43514fabc40084ec96b6d5379a74786298.
//
// Solidity: event PositionSplit(address indexed stakeholder, address collateralToken, bytes32 indexed parentCollectionId, bytes32 indexed conditionId, uint256[] partition, uint256 amount)
func (_ConditionalTokens *ConditionalTokensFilterer) ParsePositionSplit(log types.Log) (*ConditionalTokensPositionSplit, error) {
	event := new(ConditionalTokensPositionSplit)
	if err := _ConditionalTokens.contract.UnpackLog(event, "PositionSplit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensPositionsMergeIterator is returned from FilterPositionsMerge and is used to iterate over the raw logs and unpacked data for PositionsMerge events raised by the ConditionalTokens contract.
type ConditionalTokensPositionsMergeIterator struct {
	Event *ConditionalTokensPositionsMerge // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensPositionsMergeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensPositionsMerge)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensPositionsMerge)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensPositionsMergeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensPositionsMergeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensPositionsMerge represents a PositionsMerge event raised by the ConditionalTokens contract.
type ConditionalTokensPositionsMerge struct {
	Stakeholder        common.Address
	CollateralToken    common.Address
	ParentCollectionId [32]byte
	ConditionId        [32]byte
	Partition          []*big.Int
	Amount             *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterPositionsMerge is a free log retrieval operation binding the contract event 0x6f13ca62553fcc2bcd2372180a43949c1e4cebba603901ede2f4e14f36b282ca.
//
// Solidity: event PositionsMerge(address indexed stakeholder, address collateralToken, bytes32 indexed parentCollectionId, bytes32 indexed conditionId, uint256[] partition, uint256 amount)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterPositionsMerge(opts *bind.FilterOpts, stakeholder []common.Address, parentCollectionId [][32]byte, conditionId [][32]byte) (*ConditionalTokensPositionsMergeIterator, error) {

	var stakeholderRule []interface{}
	for _, stakeholderItem := range stakeholder {
		stakeholderRule = append(stakeholderRule, stakeholderItem)
	}

	var parentCollectionIdRule []interface{}
	for _, parentCollectionIdItem := range parentCollectionId {
		parentCollectionIdRule = append(parentCollectionIdRule, parentCollectionIdItem)
	}
	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "PositionsMerge", stakeholderRule, parentCollectionIdRule, conditionIdRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensPositionsMergeIterator{contract: _ConditionalTokens.contract, event: "PositionsMerge", logs: logs, sub: sub}, nil
}

// WatchPositionsMerge is a free log subscription operation binding the contract event 0x6f13ca62553fcc2bcd2372180a43949c1e4cebba603901ede2f4e14f36b282ca.
//
// Solidity: event PositionsMerge(address indexed stakeholder, address collateralToken, bytes32 indexed parentCollectionId, bytes32 indexed conditionId, uint256[] partition, uint256 amount)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchPositionsMerge(opts *bind.WatchOpts, sink chan<- *ConditionalTokensPositionsMerge, stakeholder []common.Address, parentCollectionId [][32]byte, conditionId [][32]byte) (event.Subscription, error) {

	var stakeholderRule []interface{}
	for _, stakeholderItem := range stakeholder {
		stakeholderRule = append(stakeholderRule, stakeholderItem)
	}

	var parentCollectionIdRule []interface{}
	for _, parentCollectionIdItem := range parentCollectionId {
		parentCollectionIdRule = append(parentCollectionIdRule, parentCollectionIdItem)
	}
	var conditionIdRule []interface{}
	for _, conditionIdItem := range conditionId {
		conditionIdRule = append(conditionIdRule, conditionIdItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "PositionsMerge", stakeholderRule, parentCollectionIdRule, conditionIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensPositionsMerge)
				if err := _ConditionalTokens.contract.UnpackLog(event, "PositionsMerge", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePositionsMerge is a log parse operation binding the contract event 0x6f13ca62553fcc2bcd2372180a43949c1e4cebba603901ede2f4e14f36b282ca.
//
// Solidity: event PositionsMerge(address indexed stakeholder, address collateralToken, bytes32 indexed parentCollectionId, bytes32 indexed conditionId, uint256[] partition, uint256 amount)
func (_ConditionalTokens *ConditionalTokensFilterer) ParsePositionsMerge(log types.Log) (*ConditionalTokensPositionsMerge, error) {
	event := new(ConditionalTokensPositionsMerge)
	if err := _ConditionalTokens.contract.UnpackLog(event, "PositionsMerge", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensTransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the ConditionalTokens contract.
type ConditionalTokensTransferBatchIterator struct {
	Event *ConditionalTokensTransferBatch // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensTransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensTransferBatch)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensTransferBatch)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensTransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensTransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensTransferBatch represents a TransferBatch event raised by the ConditionalTokens contract.
type ConditionalTokensTransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*ConditionalTokensTransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensTransferBatchIterator{contract: _ConditionalTokens.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *ConditionalTokensTransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensTransferBatch)
				if err := _ConditionalTokens.contract.UnpackLog(event, "TransferBatch", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_ConditionalTokens *ConditionalTokensFilterer) ParseTransferBatch(log types.Log) (*ConditionalTokensTransferBatch, error) {
	event := new(ConditionalTokensTransferBatch)
	if err := _ConditionalTokens.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensTransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the ConditionalTokens contract.
type ConditionalTokensTransferSingleIterator struct {
	Event *ConditionalTokensTransferSingle // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensTransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensTransferSingle)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensTransferSingle)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensTransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensTransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensTransferSingle represents a TransferSingle event raised by the ConditionalTokens contract.
type ConditionalTokensTransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*ConditionalTokensTransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensTransferSingleIterator{contract: _ConditionalTokens.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *ConditionalTokensTransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensTransferSingle)
				if err := _ConditionalTokens.contract.UnpackLog(event, "TransferSingle", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_ConditionalTokens *ConditionalTokensFilterer) ParseTransferSingle(log types.Log) (*ConditionalTokensTransferSingle, error) {
	event := new(ConditionalTokensTransferSingle)
	if err := _ConditionalTokens.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConditionalTokensURIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the ConditionalTokens contract.
type ConditionalTokensURIIterator struct {
	Event *ConditionalTokensURI // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ConditionalTokensURIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConditionalTokensURI)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ConditionalTokensURI)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ConditionalTokensURIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConditionalTokensURIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConditionalTokensURI represents a URI event raised by the ConditionalTokens contract.
type ConditionalTokensURI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ConditionalTokens *ConditionalTokensFilterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*ConditionalTokensURIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ConditionalTokens.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &ConditionalTokensURIIterator{contract: _ConditionalTokens.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ConditionalTokens *ConditionalTokensFilterer) WatchURI(opts *bind.WatchOpts, sink chan<- *ConditionalTokensURI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _ConditionalTokens.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConditionalTokensURI)
				if err := _ConditionalTokens.contract.UnpackLog(event, "URI", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_ConditionalTokens *ConditionalTokensFilterer) ParseURI(log types.Log) (*ConditionalTokensURI, error) {
	event := new(ConditionalTokensURI)
	if err := _ConditionalTokens.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
