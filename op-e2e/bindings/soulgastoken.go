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

// SoulGasTokenMetaData contains all meta data concerning the SoulGasToken contract.
var SoulGasTokenMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_isBackedByNative\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addBurners\",\"inputs\":[{\"name\":\"_burners\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addMinters\",\"inputs\":[{\"name\":\"_minters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowSgtValue\",\"inputs\":[{\"name\":\"_contracts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"batchBurnFrom\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchDepositFor\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"batchDepositForAll\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"batchMint\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"batchWithdrawFrom\",\"inputs\":[{\"name\":\"_accounts\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"_values\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"chargeFromOrigin\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"amountCharged_\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delBurners\",\"inputs\":[{\"name\":\"_burners\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delMinters\",\"inputs\":[{\"name\":\"_minters\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deposit\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"disallowSgtValue\",\"inputs\":[{\"name\":\"_contracts\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isBackedByNative\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdrawFrom\",\"inputs\":[{\"name\":\"_account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"AllowSgtValue\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DisallowSgtValue\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
}

// SoulGasTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use SoulGasTokenMetaData.ABI instead.
var SoulGasTokenABI = SoulGasTokenMetaData.ABI

// SoulGasToken is an auto generated Go binding around an Ethereum contract.
type SoulGasToken struct {
	SoulGasTokenCaller     // Read-only binding to the contract
	SoulGasTokenTransactor // Write-only binding to the contract
	SoulGasTokenFilterer   // Log filterer for contract events
}

// SoulGasTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type SoulGasTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SoulGasTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SoulGasTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SoulGasTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SoulGasTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SoulGasTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SoulGasTokenSession struct {
	Contract     *SoulGasToken     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SoulGasTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SoulGasTokenCallerSession struct {
	Contract *SoulGasTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// SoulGasTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SoulGasTokenTransactorSession struct {
	Contract     *SoulGasTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// SoulGasTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type SoulGasTokenRaw struct {
	Contract *SoulGasToken // Generic contract binding to access the raw methods on
}

// SoulGasTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SoulGasTokenCallerRaw struct {
	Contract *SoulGasTokenCaller // Generic read-only contract binding to access the raw methods on
}

// SoulGasTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SoulGasTokenTransactorRaw struct {
	Contract *SoulGasTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSoulGasToken creates a new instance of SoulGasToken, bound to a specific deployed contract.
func NewSoulGasToken(address common.Address, backend bind.ContractBackend) (*SoulGasToken, error) {
	contract, err := bindSoulGasToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SoulGasToken{SoulGasTokenCaller: SoulGasTokenCaller{contract: contract}, SoulGasTokenTransactor: SoulGasTokenTransactor{contract: contract}, SoulGasTokenFilterer: SoulGasTokenFilterer{contract: contract}}, nil
}

// NewSoulGasTokenCaller creates a new read-only instance of SoulGasToken, bound to a specific deployed contract.
func NewSoulGasTokenCaller(address common.Address, caller bind.ContractCaller) (*SoulGasTokenCaller, error) {
	contract, err := bindSoulGasToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenCaller{contract: contract}, nil
}

// NewSoulGasTokenTransactor creates a new write-only instance of SoulGasToken, bound to a specific deployed contract.
func NewSoulGasTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*SoulGasTokenTransactor, error) {
	contract, err := bindSoulGasToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenTransactor{contract: contract}, nil
}

// NewSoulGasTokenFilterer creates a new log filterer instance of SoulGasToken, bound to a specific deployed contract.
func NewSoulGasTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*SoulGasTokenFilterer, error) {
	contract, err := bindSoulGasToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenFilterer{contract: contract}, nil
}

// bindSoulGasToken binds a generic wrapper to an already deployed contract.
func bindSoulGasToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SoulGasTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SoulGasToken *SoulGasTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SoulGasToken.Contract.SoulGasTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SoulGasToken *SoulGasTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SoulGasToken.Contract.SoulGasTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SoulGasToken *SoulGasTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SoulGasToken.Contract.SoulGasTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SoulGasToken *SoulGasTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SoulGasToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SoulGasToken *SoulGasTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SoulGasToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SoulGasToken *SoulGasTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SoulGasToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_SoulGasToken *SoulGasTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_SoulGasToken *SoulGasTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _SoulGasToken.Contract.Allowance(&_SoulGasToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_SoulGasToken *SoulGasTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _SoulGasToken.Contract.Allowance(&_SoulGasToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_SoulGasToken *SoulGasTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_SoulGasToken *SoulGasTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _SoulGasToken.Contract.BalanceOf(&_SoulGasToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_SoulGasToken *SoulGasTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _SoulGasToken.Contract.BalanceOf(&_SoulGasToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SoulGasToken *SoulGasTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SoulGasToken *SoulGasTokenSession) Decimals() (uint8, error) {
	return _SoulGasToken.Contract.Decimals(&_SoulGasToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_SoulGasToken *SoulGasTokenCallerSession) Decimals() (uint8, error) {
	return _SoulGasToken.Contract.Decimals(&_SoulGasToken.CallOpts)
}

// IsBackedByNative is a free data retrieval call binding the contract method 0xbbd10120.
//
// Solidity: function isBackedByNative() view returns(bool)
func (_SoulGasToken *SoulGasTokenCaller) IsBackedByNative(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "isBackedByNative")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBackedByNative is a free data retrieval call binding the contract method 0xbbd10120.
//
// Solidity: function isBackedByNative() view returns(bool)
func (_SoulGasToken *SoulGasTokenSession) IsBackedByNative() (bool, error) {
	return _SoulGasToken.Contract.IsBackedByNative(&_SoulGasToken.CallOpts)
}

// IsBackedByNative is a free data retrieval call binding the contract method 0xbbd10120.
//
// Solidity: function isBackedByNative() view returns(bool)
func (_SoulGasToken *SoulGasTokenCallerSession) IsBackedByNative() (bool, error) {
	return _SoulGasToken.Contract.IsBackedByNative(&_SoulGasToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SoulGasToken *SoulGasTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SoulGasToken *SoulGasTokenSession) Name() (string, error) {
	return _SoulGasToken.Contract.Name(&_SoulGasToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_SoulGasToken *SoulGasTokenCallerSession) Name() (string, error) {
	return _SoulGasToken.Contract.Name(&_SoulGasToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SoulGasToken *SoulGasTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SoulGasToken *SoulGasTokenSession) Owner() (common.Address, error) {
	return _SoulGasToken.Contract.Owner(&_SoulGasToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SoulGasToken *SoulGasTokenCallerSession) Owner() (common.Address, error) {
	return _SoulGasToken.Contract.Owner(&_SoulGasToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SoulGasToken *SoulGasTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SoulGasToken *SoulGasTokenSession) Symbol() (string, error) {
	return _SoulGasToken.Contract.Symbol(&_SoulGasToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_SoulGasToken *SoulGasTokenCallerSession) Symbol() (string, error) {
	return _SoulGasToken.Contract.Symbol(&_SoulGasToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SoulGasToken *SoulGasTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SoulGasToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SoulGasToken *SoulGasTokenSession) TotalSupply() (*big.Int, error) {
	return _SoulGasToken.Contract.TotalSupply(&_SoulGasToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_SoulGasToken *SoulGasTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _SoulGasToken.Contract.TotalSupply(&_SoulGasToken.CallOpts)
}

// AddBurners is a paid mutator transaction binding the contract method 0x3ab84dd9.
//
// Solidity: function addBurners(address[] _burners) returns()
func (_SoulGasToken *SoulGasTokenTransactor) AddBurners(opts *bind.TransactOpts, _burners []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "addBurners", _burners)
}

// AddBurners is a paid mutator transaction binding the contract method 0x3ab84dd9.
//
// Solidity: function addBurners(address[] _burners) returns()
func (_SoulGasToken *SoulGasTokenSession) AddBurners(_burners []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.AddBurners(&_SoulGasToken.TransactOpts, _burners)
}

// AddBurners is a paid mutator transaction binding the contract method 0x3ab84dd9.
//
// Solidity: function addBurners(address[] _burners) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) AddBurners(_burners []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.AddBurners(&_SoulGasToken.TransactOpts, _burners)
}

// AddMinters is a paid mutator transaction binding the contract method 0x71e2a657.
//
// Solidity: function addMinters(address[] _minters) returns()
func (_SoulGasToken *SoulGasTokenTransactor) AddMinters(opts *bind.TransactOpts, _minters []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "addMinters", _minters)
}

// AddMinters is a paid mutator transaction binding the contract method 0x71e2a657.
//
// Solidity: function addMinters(address[] _minters) returns()
func (_SoulGasToken *SoulGasTokenSession) AddMinters(_minters []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.AddMinters(&_SoulGasToken.TransactOpts, _minters)
}

// AddMinters is a paid mutator transaction binding the contract method 0x71e2a657.
//
// Solidity: function addMinters(address[] _minters) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) AddMinters(_minters []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.AddMinters(&_SoulGasToken.TransactOpts, _minters)
}

// AllowSgtValue is a paid mutator transaction binding the contract method 0x674e29ea.
//
// Solidity: function allowSgtValue(address[] _contracts) returns()
func (_SoulGasToken *SoulGasTokenTransactor) AllowSgtValue(opts *bind.TransactOpts, _contracts []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "allowSgtValue", _contracts)
}

// AllowSgtValue is a paid mutator transaction binding the contract method 0x674e29ea.
//
// Solidity: function allowSgtValue(address[] _contracts) returns()
func (_SoulGasToken *SoulGasTokenSession) AllowSgtValue(_contracts []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.AllowSgtValue(&_SoulGasToken.TransactOpts, _contracts)
}

// AllowSgtValue is a paid mutator transaction binding the contract method 0x674e29ea.
//
// Solidity: function allowSgtValue(address[] _contracts) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) AllowSgtValue(_contracts []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.AllowSgtValue(&_SoulGasToken.TransactOpts, _contracts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactor) Approve(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "approve", arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.Approve(&_SoulGasToken.TransactOpts, arg0, arg1)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactorSession) Approve(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.Approve(&_SoulGasToken.TransactOpts, arg0, arg1)
}

// BatchBurnFrom is a paid mutator transaction binding the contract method 0x1b9a7529.
//
// Solidity: function batchBurnFrom(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenTransactor) BatchBurnFrom(opts *bind.TransactOpts, _accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "batchBurnFrom", _accounts, _values)
}

// BatchBurnFrom is a paid mutator transaction binding the contract method 0x1b9a7529.
//
// Solidity: function batchBurnFrom(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenSession) BatchBurnFrom(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchBurnFrom(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BatchBurnFrom is a paid mutator transaction binding the contract method 0x1b9a7529.
//
// Solidity: function batchBurnFrom(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) BatchBurnFrom(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchBurnFrom(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BatchDepositFor is a paid mutator transaction binding the contract method 0x299f8170.
//
// Solidity: function batchDepositFor(address[] _accounts, uint256[] _values) payable returns()
func (_SoulGasToken *SoulGasTokenTransactor) BatchDepositFor(opts *bind.TransactOpts, _accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "batchDepositFor", _accounts, _values)
}

// BatchDepositFor is a paid mutator transaction binding the contract method 0x299f8170.
//
// Solidity: function batchDepositFor(address[] _accounts, uint256[] _values) payable returns()
func (_SoulGasToken *SoulGasTokenSession) BatchDepositFor(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchDepositFor(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BatchDepositFor is a paid mutator transaction binding the contract method 0x299f8170.
//
// Solidity: function batchDepositFor(address[] _accounts, uint256[] _values) payable returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) BatchDepositFor(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchDepositFor(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BatchDepositForAll is a paid mutator transaction binding the contract method 0x84e08810.
//
// Solidity: function batchDepositForAll(address[] _accounts, uint256 _value) payable returns()
func (_SoulGasToken *SoulGasTokenTransactor) BatchDepositForAll(opts *bind.TransactOpts, _accounts []common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "batchDepositForAll", _accounts, _value)
}

// BatchDepositForAll is a paid mutator transaction binding the contract method 0x84e08810.
//
// Solidity: function batchDepositForAll(address[] _accounts, uint256 _value) payable returns()
func (_SoulGasToken *SoulGasTokenSession) BatchDepositForAll(_accounts []common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchDepositForAll(&_SoulGasToken.TransactOpts, _accounts, _value)
}

// BatchDepositForAll is a paid mutator transaction binding the contract method 0x84e08810.
//
// Solidity: function batchDepositForAll(address[] _accounts, uint256 _value) payable returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) BatchDepositForAll(_accounts []common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchDepositForAll(&_SoulGasToken.TransactOpts, _accounts, _value)
}

// BatchMint is a paid mutator transaction binding the contract method 0x68573107.
//
// Solidity: function batchMint(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenTransactor) BatchMint(opts *bind.TransactOpts, _accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "batchMint", _accounts, _values)
}

// BatchMint is a paid mutator transaction binding the contract method 0x68573107.
//
// Solidity: function batchMint(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenSession) BatchMint(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchMint(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BatchMint is a paid mutator transaction binding the contract method 0x68573107.
//
// Solidity: function batchMint(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) BatchMint(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchMint(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BatchWithdrawFrom is a paid mutator transaction binding the contract method 0xb3e2a832.
//
// Solidity: function batchWithdrawFrom(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenTransactor) BatchWithdrawFrom(opts *bind.TransactOpts, _accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "batchWithdrawFrom", _accounts, _values)
}

// BatchWithdrawFrom is a paid mutator transaction binding the contract method 0xb3e2a832.
//
// Solidity: function batchWithdrawFrom(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenSession) BatchWithdrawFrom(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchWithdrawFrom(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BatchWithdrawFrom is a paid mutator transaction binding the contract method 0xb3e2a832.
//
// Solidity: function batchWithdrawFrom(address[] _accounts, uint256[] _values) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) BatchWithdrawFrom(_accounts []common.Address, _values []*big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BatchWithdrawFrom(&_SoulGasToken.TransactOpts, _accounts, _values)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address _account, uint256 _value) returns()
func (_SoulGasToken *SoulGasTokenTransactor) BurnFrom(opts *bind.TransactOpts, _account common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "burnFrom", _account, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address _account, uint256 _value) returns()
func (_SoulGasToken *SoulGasTokenSession) BurnFrom(_account common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BurnFrom(&_SoulGasToken.TransactOpts, _account, _value)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address _account, uint256 _value) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) BurnFrom(_account common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.BurnFrom(&_SoulGasToken.TransactOpts, _account, _value)
}

// ChargeFromOrigin is a paid mutator transaction binding the contract method 0xce25c861.
//
// Solidity: function chargeFromOrigin(uint256 _amount) returns(uint256 amountCharged_)
func (_SoulGasToken *SoulGasTokenTransactor) ChargeFromOrigin(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "chargeFromOrigin", _amount)
}

// ChargeFromOrigin is a paid mutator transaction binding the contract method 0xce25c861.
//
// Solidity: function chargeFromOrigin(uint256 _amount) returns(uint256 amountCharged_)
func (_SoulGasToken *SoulGasTokenSession) ChargeFromOrigin(_amount *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.ChargeFromOrigin(&_SoulGasToken.TransactOpts, _amount)
}

// ChargeFromOrigin is a paid mutator transaction binding the contract method 0xce25c861.
//
// Solidity: function chargeFromOrigin(uint256 _amount) returns(uint256 amountCharged_)
func (_SoulGasToken *SoulGasTokenTransactorSession) ChargeFromOrigin(_amount *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.ChargeFromOrigin(&_SoulGasToken.TransactOpts, _amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_SoulGasToken *SoulGasTokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DecreaseAllowance(&_SoulGasToken.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DecreaseAllowance(&_SoulGasToken.TransactOpts, spender, subtractedValue)
}

// DelBurners is a paid mutator transaction binding the contract method 0xb8de86b4.
//
// Solidity: function delBurners(address[] _burners) returns()
func (_SoulGasToken *SoulGasTokenTransactor) DelBurners(opts *bind.TransactOpts, _burners []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "delBurners", _burners)
}

// DelBurners is a paid mutator transaction binding the contract method 0xb8de86b4.
//
// Solidity: function delBurners(address[] _burners) returns()
func (_SoulGasToken *SoulGasTokenSession) DelBurners(_burners []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DelBurners(&_SoulGasToken.TransactOpts, _burners)
}

// DelBurners is a paid mutator transaction binding the contract method 0xb8de86b4.
//
// Solidity: function delBurners(address[] _burners) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) DelBurners(_burners []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DelBurners(&_SoulGasToken.TransactOpts, _burners)
}

// DelMinters is a paid mutator transaction binding the contract method 0xe04b8180.
//
// Solidity: function delMinters(address[] _minters) returns()
func (_SoulGasToken *SoulGasTokenTransactor) DelMinters(opts *bind.TransactOpts, _minters []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "delMinters", _minters)
}

// DelMinters is a paid mutator transaction binding the contract method 0xe04b8180.
//
// Solidity: function delMinters(address[] _minters) returns()
func (_SoulGasToken *SoulGasTokenSession) DelMinters(_minters []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DelMinters(&_SoulGasToken.TransactOpts, _minters)
}

// DelMinters is a paid mutator transaction binding the contract method 0xe04b8180.
//
// Solidity: function delMinters(address[] _minters) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) DelMinters(_minters []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DelMinters(&_SoulGasToken.TransactOpts, _minters)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_SoulGasToken *SoulGasTokenTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_SoulGasToken *SoulGasTokenSession) Deposit() (*types.Transaction, error) {
	return _SoulGasToken.Contract.Deposit(&_SoulGasToken.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) Deposit() (*types.Transaction, error) {
	return _SoulGasToken.Contract.Deposit(&_SoulGasToken.TransactOpts)
}

// DisallowSgtValue is a paid mutator transaction binding the contract method 0xdc270eb8.
//
// Solidity: function disallowSgtValue(address[] _contracts) returns()
func (_SoulGasToken *SoulGasTokenTransactor) DisallowSgtValue(opts *bind.TransactOpts, _contracts []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "disallowSgtValue", _contracts)
}

// DisallowSgtValue is a paid mutator transaction binding the contract method 0xdc270eb8.
//
// Solidity: function disallowSgtValue(address[] _contracts) returns()
func (_SoulGasToken *SoulGasTokenSession) DisallowSgtValue(_contracts []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DisallowSgtValue(&_SoulGasToken.TransactOpts, _contracts)
}

// DisallowSgtValue is a paid mutator transaction binding the contract method 0xdc270eb8.
//
// Solidity: function disallowSgtValue(address[] _contracts) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) DisallowSgtValue(_contracts []common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.DisallowSgtValue(&_SoulGasToken.TransactOpts, _contracts)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_SoulGasToken *SoulGasTokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.IncreaseAllowance(&_SoulGasToken.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.IncreaseAllowance(&_SoulGasToken.TransactOpts, spender, addedValue)
}

// Initialize is a paid mutator transaction binding the contract method 0x077f224a.
//
// Solidity: function initialize(string _name, string _symbol, address _owner) returns()
func (_SoulGasToken *SoulGasTokenTransactor) Initialize(opts *bind.TransactOpts, _name string, _symbol string, _owner common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "initialize", _name, _symbol, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x077f224a.
//
// Solidity: function initialize(string _name, string _symbol, address _owner) returns()
func (_SoulGasToken *SoulGasTokenSession) Initialize(_name string, _symbol string, _owner common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.Initialize(&_SoulGasToken.TransactOpts, _name, _symbol, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0x077f224a.
//
// Solidity: function initialize(string _name, string _symbol, address _owner) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) Initialize(_name string, _symbol string, _owner common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.Initialize(&_SoulGasToken.TransactOpts, _name, _symbol, _owner)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SoulGasToken *SoulGasTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SoulGasToken *SoulGasTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _SoulGasToken.Contract.RenounceOwnership(&_SoulGasToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SoulGasToken.Contract.RenounceOwnership(&_SoulGasToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactor) Transfer(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "transfer", arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.Transfer(&_SoulGasToken.TransactOpts, arg0, arg1)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactorSession) Transfer(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.Transfer(&_SoulGasToken.TransactOpts, arg0, arg1)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactor) TransferFrom(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "transferFrom", arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.TransferFrom(&_SoulGasToken.TransactOpts, arg0, arg1, arg2)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address , address , uint256 ) returns(bool)
func (_SoulGasToken *SoulGasTokenTransactorSession) TransferFrom(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.TransferFrom(&_SoulGasToken.TransactOpts, arg0, arg1, arg2)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SoulGasToken *SoulGasTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SoulGasToken *SoulGasTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.TransferOwnership(&_SoulGasToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SoulGasToken.Contract.TransferOwnership(&_SoulGasToken.TransactOpts, newOwner)
}

// WithdrawFrom is a paid mutator transaction binding the contract method 0x9470b0bd.
//
// Solidity: function withdrawFrom(address _account, uint256 _value) returns()
func (_SoulGasToken *SoulGasTokenTransactor) WithdrawFrom(opts *bind.TransactOpts, _account common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.contract.Transact(opts, "withdrawFrom", _account, _value)
}

// WithdrawFrom is a paid mutator transaction binding the contract method 0x9470b0bd.
//
// Solidity: function withdrawFrom(address _account, uint256 _value) returns()
func (_SoulGasToken *SoulGasTokenSession) WithdrawFrom(_account common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.WithdrawFrom(&_SoulGasToken.TransactOpts, _account, _value)
}

// WithdrawFrom is a paid mutator transaction binding the contract method 0x9470b0bd.
//
// Solidity: function withdrawFrom(address _account, uint256 _value) returns()
func (_SoulGasToken *SoulGasTokenTransactorSession) WithdrawFrom(_account common.Address, _value *big.Int) (*types.Transaction, error) {
	return _SoulGasToken.Contract.WithdrawFrom(&_SoulGasToken.TransactOpts, _account, _value)
}

// SoulGasTokenAllowSgtValueIterator is returned from FilterAllowSgtValue and is used to iterate over the raw logs and unpacked data for AllowSgtValue events raised by the SoulGasToken contract.
type SoulGasTokenAllowSgtValueIterator struct {
	Event *SoulGasTokenAllowSgtValue // Event containing the contract specifics and raw log

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
func (it *SoulGasTokenAllowSgtValueIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SoulGasTokenAllowSgtValue)
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
		it.Event = new(SoulGasTokenAllowSgtValue)
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
func (it *SoulGasTokenAllowSgtValueIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SoulGasTokenAllowSgtValueIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SoulGasTokenAllowSgtValue represents a AllowSgtValue event raised by the SoulGasToken contract.
type SoulGasTokenAllowSgtValue struct {
	From common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAllowSgtValue is a free log retrieval operation binding the contract event 0xf135aca2ee4483470b8f44f38ab676fc36fc67437777f3c520e5fbeb3706009f.
//
// Solidity: event AllowSgtValue(address indexed from)
func (_SoulGasToken *SoulGasTokenFilterer) FilterAllowSgtValue(opts *bind.FilterOpts, from []common.Address) (*SoulGasTokenAllowSgtValueIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SoulGasToken.contract.FilterLogs(opts, "AllowSgtValue", fromRule)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenAllowSgtValueIterator{contract: _SoulGasToken.contract, event: "AllowSgtValue", logs: logs, sub: sub}, nil
}

// WatchAllowSgtValue is a free log subscription operation binding the contract event 0xf135aca2ee4483470b8f44f38ab676fc36fc67437777f3c520e5fbeb3706009f.
//
// Solidity: event AllowSgtValue(address indexed from)
func (_SoulGasToken *SoulGasTokenFilterer) WatchAllowSgtValue(opts *bind.WatchOpts, sink chan<- *SoulGasTokenAllowSgtValue, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SoulGasToken.contract.WatchLogs(opts, "AllowSgtValue", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SoulGasTokenAllowSgtValue)
				if err := _SoulGasToken.contract.UnpackLog(event, "AllowSgtValue", log); err != nil {
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

// ParseAllowSgtValue is a log parse operation binding the contract event 0xf135aca2ee4483470b8f44f38ab676fc36fc67437777f3c520e5fbeb3706009f.
//
// Solidity: event AllowSgtValue(address indexed from)
func (_SoulGasToken *SoulGasTokenFilterer) ParseAllowSgtValue(log types.Log) (*SoulGasTokenAllowSgtValue, error) {
	event := new(SoulGasTokenAllowSgtValue)
	if err := _SoulGasToken.contract.UnpackLog(event, "AllowSgtValue", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SoulGasTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the SoulGasToken contract.
type SoulGasTokenApprovalIterator struct {
	Event *SoulGasTokenApproval // Event containing the contract specifics and raw log

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
func (it *SoulGasTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SoulGasTokenApproval)
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
		it.Event = new(SoulGasTokenApproval)
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
func (it *SoulGasTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SoulGasTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SoulGasTokenApproval represents a Approval event raised by the SoulGasToken contract.
type SoulGasTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SoulGasToken *SoulGasTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*SoulGasTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SoulGasToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenApprovalIterator{contract: _SoulGasToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SoulGasToken *SoulGasTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *SoulGasTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _SoulGasToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SoulGasTokenApproval)
				if err := _SoulGasToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_SoulGasToken *SoulGasTokenFilterer) ParseApproval(log types.Log) (*SoulGasTokenApproval, error) {
	event := new(SoulGasTokenApproval)
	if err := _SoulGasToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SoulGasTokenDisallowSgtValueIterator is returned from FilterDisallowSgtValue and is used to iterate over the raw logs and unpacked data for DisallowSgtValue events raised by the SoulGasToken contract.
type SoulGasTokenDisallowSgtValueIterator struct {
	Event *SoulGasTokenDisallowSgtValue // Event containing the contract specifics and raw log

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
func (it *SoulGasTokenDisallowSgtValueIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SoulGasTokenDisallowSgtValue)
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
		it.Event = new(SoulGasTokenDisallowSgtValue)
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
func (it *SoulGasTokenDisallowSgtValueIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SoulGasTokenDisallowSgtValueIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SoulGasTokenDisallowSgtValue represents a DisallowSgtValue event raised by the SoulGasToken contract.
type SoulGasTokenDisallowSgtValue struct {
	From common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDisallowSgtValue is a free log retrieval operation binding the contract event 0x42d3350598a4a2ec6e60463e0bffa1aab494a9e8d4484b017270dde628b4edb1.
//
// Solidity: event DisallowSgtValue(address indexed from)
func (_SoulGasToken *SoulGasTokenFilterer) FilterDisallowSgtValue(opts *bind.FilterOpts, from []common.Address) (*SoulGasTokenDisallowSgtValueIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SoulGasToken.contract.FilterLogs(opts, "DisallowSgtValue", fromRule)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenDisallowSgtValueIterator{contract: _SoulGasToken.contract, event: "DisallowSgtValue", logs: logs, sub: sub}, nil
}

// WatchDisallowSgtValue is a free log subscription operation binding the contract event 0x42d3350598a4a2ec6e60463e0bffa1aab494a9e8d4484b017270dde628b4edb1.
//
// Solidity: event DisallowSgtValue(address indexed from)
func (_SoulGasToken *SoulGasTokenFilterer) WatchDisallowSgtValue(opts *bind.WatchOpts, sink chan<- *SoulGasTokenDisallowSgtValue, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _SoulGasToken.contract.WatchLogs(opts, "DisallowSgtValue", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SoulGasTokenDisallowSgtValue)
				if err := _SoulGasToken.contract.UnpackLog(event, "DisallowSgtValue", log); err != nil {
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

// ParseDisallowSgtValue is a log parse operation binding the contract event 0x42d3350598a4a2ec6e60463e0bffa1aab494a9e8d4484b017270dde628b4edb1.
//
// Solidity: event DisallowSgtValue(address indexed from)
func (_SoulGasToken *SoulGasTokenFilterer) ParseDisallowSgtValue(log types.Log) (*SoulGasTokenDisallowSgtValue, error) {
	event := new(SoulGasTokenDisallowSgtValue)
	if err := _SoulGasToken.contract.UnpackLog(event, "DisallowSgtValue", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SoulGasTokenInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SoulGasToken contract.
type SoulGasTokenInitializedIterator struct {
	Event *SoulGasTokenInitialized // Event containing the contract specifics and raw log

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
func (it *SoulGasTokenInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SoulGasTokenInitialized)
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
		it.Event = new(SoulGasTokenInitialized)
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
func (it *SoulGasTokenInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SoulGasTokenInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SoulGasTokenInitialized represents a Initialized event raised by the SoulGasToken contract.
type SoulGasTokenInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SoulGasToken *SoulGasTokenFilterer) FilterInitialized(opts *bind.FilterOpts) (*SoulGasTokenInitializedIterator, error) {

	logs, sub, err := _SoulGasToken.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenInitializedIterator{contract: _SoulGasToken.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SoulGasToken *SoulGasTokenFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SoulGasTokenInitialized) (event.Subscription, error) {

	logs, sub, err := _SoulGasToken.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SoulGasTokenInitialized)
				if err := _SoulGasToken.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SoulGasToken *SoulGasTokenFilterer) ParseInitialized(log types.Log) (*SoulGasTokenInitialized, error) {
	event := new(SoulGasTokenInitialized)
	if err := _SoulGasToken.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SoulGasTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SoulGasToken contract.
type SoulGasTokenOwnershipTransferredIterator struct {
	Event *SoulGasTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SoulGasTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SoulGasTokenOwnershipTransferred)
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
		it.Event = new(SoulGasTokenOwnershipTransferred)
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
func (it *SoulGasTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SoulGasTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SoulGasTokenOwnershipTransferred represents a OwnershipTransferred event raised by the SoulGasToken contract.
type SoulGasTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SoulGasToken *SoulGasTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SoulGasTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SoulGasToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenOwnershipTransferredIterator{contract: _SoulGasToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SoulGasToken *SoulGasTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SoulGasTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SoulGasToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SoulGasTokenOwnershipTransferred)
				if err := _SoulGasToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SoulGasToken *SoulGasTokenFilterer) ParseOwnershipTransferred(log types.Log) (*SoulGasTokenOwnershipTransferred, error) {
	event := new(SoulGasTokenOwnershipTransferred)
	if err := _SoulGasToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SoulGasTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the SoulGasToken contract.
type SoulGasTokenTransferIterator struct {
	Event *SoulGasTokenTransfer // Event containing the contract specifics and raw log

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
func (it *SoulGasTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SoulGasTokenTransfer)
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
		it.Event = new(SoulGasTokenTransfer)
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
func (it *SoulGasTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SoulGasTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SoulGasTokenTransfer represents a Transfer event raised by the SoulGasToken contract.
type SoulGasTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SoulGasToken *SoulGasTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*SoulGasTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SoulGasToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &SoulGasTokenTransferIterator{contract: _SoulGasToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SoulGasToken *SoulGasTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *SoulGasTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _SoulGasToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SoulGasTokenTransfer)
				if err := _SoulGasToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_SoulGasToken *SoulGasTokenFilterer) ParseTransfer(log types.Log) (*SoulGasTokenTransfer, error) {
	event := new(SoulGasTokenTransfer)
	if err := _SoulGasToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
