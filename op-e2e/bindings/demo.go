// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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

// DemoMetaData contains all meta data concerning the Demo contract.
var DemoMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DemoABI is the input ABI used to generate the binding from.
// Deprecated: Use DemoMetaData.ABI instead.
var DemoABI = DemoMetaData.ABI

// Demo is an auto generated Go binding around an Ethereum contract.
type Demo struct {
	DemoCaller     // Read-only binding to the contract
	DemoTransactor // Write-only binding to the contract
	DemoFilterer   // Log filterer for contract events
}

// DemoCaller is an auto generated read-only Go binding around an Ethereum contract.
type DemoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DemoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DemoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DemoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DemoSession struct {
	Contract     *Demo             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DemoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DemoCallerSession struct {
	Contract *DemoCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DemoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DemoTransactorSession struct {
	Contract     *DemoTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DemoRaw is an auto generated low-level Go binding around an Ethereum contract.
type DemoRaw struct {
	Contract *Demo // Generic contract binding to access the raw methods on
}

// DemoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DemoCallerRaw struct {
	Contract *DemoCaller // Generic read-only contract binding to access the raw methods on
}

// DemoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DemoTransactorRaw struct {
	Contract *DemoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDemo creates a new instance of Demo, bound to a specific deployed contract.
func NewDemo(address common.Address, backend bind.ContractBackend) (*Demo, error) {
	contract, err := bindDemo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Demo{DemoCaller: DemoCaller{contract: contract}, DemoTransactor: DemoTransactor{contract: contract}, DemoFilterer: DemoFilterer{contract: contract}}, nil
}

// NewDemoCaller creates a new read-only instance of Demo, bound to a specific deployed contract.
func NewDemoCaller(address common.Address, caller bind.ContractCaller) (*DemoCaller, error) {
	contract, err := bindDemo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DemoCaller{contract: contract}, nil
}

// NewDemoTransactor creates a new write-only instance of Demo, bound to a specific deployed contract.
func NewDemoTransactor(address common.Address, transactor bind.ContractTransactor) (*DemoTransactor, error) {
	contract, err := bindDemo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DemoTransactor{contract: contract}, nil
}

// NewDemoFilterer creates a new log filterer instance of Demo, bound to a specific deployed contract.
func NewDemoFilterer(address common.Address, filterer bind.ContractFilterer) (*DemoFilterer, error) {
	contract, err := bindDemo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DemoFilterer{contract: contract}, nil
}

// bindDemo binds a generic wrapper to an already deployed contract.
func bindDemo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DemoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Demo *DemoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Demo.Contract.DemoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Demo *DemoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Demo.Contract.DemoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Demo *DemoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Demo.Contract.DemoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Demo *DemoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Demo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Demo *DemoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Demo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Demo *DemoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Demo.Contract.contract.Transact(opts, method, params...)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Demo *DemoCaller) Retrieve(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Demo.contract.Call(opts, &out, "retrieve")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Demo *DemoSession) Retrieve() (*big.Int, error) {
	return _Demo.Contract.Retrieve(&_Demo.CallOpts)
}

// Retrieve is a free data retrieval call binding the contract method 0x2e64cec1.
//
// Solidity: function retrieve() view returns(uint256)
func (_Demo *DemoCallerSession) Retrieve() (*big.Int, error) {
	return _Demo.Contract.Retrieve(&_Demo.CallOpts)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Demo *DemoTransactor) Store(opts *bind.TransactOpts, num *big.Int) (*types.Transaction, error) {
	return _Demo.contract.Transact(opts, "store", num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Demo *DemoSession) Store(num *big.Int) (*types.Transaction, error) {
	return _Demo.Contract.Store(&_Demo.TransactOpts, num)
}

// Store is a paid mutator transaction binding the contract method 0x6057361d.
//
// Solidity: function store(uint256 num) returns()
func (_Demo *DemoTransactorSession) Store(num *big.Int) (*types.Transaction, error) {
	return _Demo.Contract.Store(&_Demo.TransactOpts, num)
}
