// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

// The contract is located at https://github.com/ethstorage/es-op-batchinbox/blob/main/src/BatchInbox.sol

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

// BatchInboxMetaData contains all meta data concerning the BatchInbox contract.
var BatchInboxMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_esStorageContract\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BalanceNotEnough\",\"type\":\"error\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"esStorageContract\",\"outputs\":[{\"internalType\":\"contractStorageContract\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60a060405234801561000f575f5ffd5b50604051610984380380610984833981810160405281019061003191906100c9565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1681525050506100f4565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100988261006f565b9050919050565b6100a88161008e565b81146100b2575f5ffd5b50565b5f815190506100c38161009f565b92915050565b5f602082840312156100de576100dd61006b565b5b5f6100eb848285016100b5565b91505092915050565b60805161086a61011a5f395f81816101990152818161022901526103b3015261086a5ff3fe608060405260043610610042575f3560e01c806322eb767d1461006f57806327e235e314610099578063f340fa01146100d5578063f3fef3a3146100f15761005b565b3661005b576100513334610119565b610059610178565b005b6100653334610119565b61006d610178565b005b34801561007a575f5ffd5b506100836103b1565b604051610090919061057f565b60405180910390f35b3480156100a4575f5ffd5b506100bf60048036038101906100ba91906105d7565b6103d5565b6040516100cc919061061a565b60405180910390f35b6100ef60048036038101906100ea91906105d7565b6103e9565b005b3480156100fc575f5ffd5b506101176004803603810190610112919061065d565b6103f6565b005b5f81031561017457805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461016c91906106c8565b925050819055505b5050565b5f5f90505f5f90505f5b824990505f5f1b8103156102c5575f8203610227577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631ccbc6da6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610200573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610224919061070f565b91505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634581a920838386620200006040518563ffffffff1660e01b81526004016102889392919061078b565b5f604051808303818588803b15801561029f575f5ffd5b505af11580156102b1573d5f5f3e3d5ffd5b505050505082806001019350506001610182575b5f83036102d4575050506103af565b5f83836102e191906107c0565b90505f5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205490508181101561035d576040517f9882883500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81816103699190610801565b5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555050505050505b565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f602052805f5260405f205f915090505481565b6103f38134610119565b50565b5f5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610470576040517f9882883500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818161047c9190610801565b5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508273ffffffffffffffffffffffffffffffffffffffff166108fc8390811502906040515f60405180830381858888f193505050501580156104ff573d5f5f3e3d5ffd5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f61054761054261053d84610505565b610524565b610505565b9050919050565b5f6105588261052d565b9050919050565b5f6105698261054e565b9050919050565b6105798161055f565b82525050565b5f6020820190506105925f830184610570565b92915050565b5f5ffd5b5f6105a682610505565b9050919050565b6105b68161059c565b81146105c0575f5ffd5b50565b5f813590506105d1816105ad565b92915050565b5f602082840312156105ec576105eb610598565b5b5f6105f9848285016105c3565b91505092915050565b5f819050919050565b61061481610602565b82525050565b5f60208201905061062d5f83018461060b565b92915050565b61063c81610602565b8114610646575f5ffd5b50565b5f8135905061065781610633565b92915050565b5f5f6040838503121561067357610672610598565b5b5f610680858286016105c3565b925050602061069185828601610649565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6106d282610602565b91506106dd83610602565b92508282019050808211156106f5576106f461069b565b5b92915050565b5f8151905061070981610633565b92915050565b5f6020828403121561072457610723610598565b5b5f610731848285016106fb565b91505092915050565b5f819050919050565b61074c8161073a565b82525050565b5f819050919050565b5f61077561077061076b84610752565b610524565b610602565b9050919050565b6107858161075b565b82525050565b5f60608201905061079e5f830186610743565b6107ab602083018561060b565b6107b8604083018461077c565b949350505050565b5f6107ca82610602565b91506107d583610602565b92508282026107e381610602565b915082820484148315176107fa576107f961069b565b5b5092915050565b5f61080b82610602565b915061081683610602565b925082820390508181111561082e5761082d61069b565b5b9291505056fea2646970667358221220ad769cb754fdad1ddf58c9c13e0a15d938b6139b499ebc6a85cd815b50b05f8b64736f6c634300081c0033",
}

// BatchInboxABI is the input ABI used to generate the binding from.
// Deprecated: Use BatchInboxMetaData.ABI instead.
var BatchInboxABI = BatchInboxMetaData.ABI

// BatchInboxBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BatchInboxMetaData.Bin instead.
var BatchInboxBin = BatchInboxMetaData.Bin

// DeployBatchInbox deploys a new Ethereum contract, binding an instance of BatchInbox to it.
func DeployBatchInbox(auth *bind.TransactOpts, backend bind.ContractBackend, _esStorageContract common.Address) (common.Address, *types.Transaction, *BatchInbox, error) {
	parsed, err := BatchInboxMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchInboxBin), backend, _esStorageContract)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchInbox{BatchInboxCaller: BatchInboxCaller{contract: contract}, BatchInboxTransactor: BatchInboxTransactor{contract: contract}, BatchInboxFilterer: BatchInboxFilterer{contract: contract}}, nil
}

// BatchInbox is an auto generated Go binding around an Ethereum contract.
type BatchInbox struct {
	BatchInboxCaller     // Read-only binding to the contract
	BatchInboxTransactor // Write-only binding to the contract
	BatchInboxFilterer   // Log filterer for contract events
}

// BatchInboxCaller is an auto generated read-only Go binding around an Ethereum contract.
type BatchInboxCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchInboxTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BatchInboxTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchInboxFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BatchInboxFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchInboxSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BatchInboxSession struct {
	Contract     *BatchInbox       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BatchInboxCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BatchInboxCallerSession struct {
	Contract *BatchInboxCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// BatchInboxTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BatchInboxTransactorSession struct {
	Contract     *BatchInboxTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// BatchInboxRaw is an auto generated low-level Go binding around an Ethereum contract.
type BatchInboxRaw struct {
	Contract *BatchInbox // Generic contract binding to access the raw methods on
}

// BatchInboxCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BatchInboxCallerRaw struct {
	Contract *BatchInboxCaller // Generic read-only contract binding to access the raw methods on
}

// BatchInboxTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BatchInboxTransactorRaw struct {
	Contract *BatchInboxTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBatchInbox creates a new instance of BatchInbox, bound to a specific deployed contract.
func NewBatchInbox(address common.Address, backend bind.ContractBackend) (*BatchInbox, error) {
	contract, err := bindBatchInbox(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchInbox{BatchInboxCaller: BatchInboxCaller{contract: contract}, BatchInboxTransactor: BatchInboxTransactor{contract: contract}, BatchInboxFilterer: BatchInboxFilterer{contract: contract}}, nil
}

// NewBatchInboxCaller creates a new read-only instance of BatchInbox, bound to a specific deployed contract.
func NewBatchInboxCaller(address common.Address, caller bind.ContractCaller) (*BatchInboxCaller, error) {
	contract, err := bindBatchInbox(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchInboxCaller{contract: contract}, nil
}

// NewBatchInboxTransactor creates a new write-only instance of BatchInbox, bound to a specific deployed contract.
func NewBatchInboxTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchInboxTransactor, error) {
	contract, err := bindBatchInbox(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchInboxTransactor{contract: contract}, nil
}

// NewBatchInboxFilterer creates a new log filterer instance of BatchInbox, bound to a specific deployed contract.
func NewBatchInboxFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchInboxFilterer, error) {
	contract, err := bindBatchInbox(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchInboxFilterer{contract: contract}, nil
}

// bindBatchInbox binds a generic wrapper to an already deployed contract.
func bindBatchInbox(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchInboxMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchInbox *BatchInboxRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchInbox.Contract.BatchInboxCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchInbox *BatchInboxRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchInbox.Contract.BatchInboxTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchInbox *BatchInboxRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchInbox.Contract.BatchInboxTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchInbox *BatchInboxCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchInbox.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchInbox *BatchInboxTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchInbox.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchInbox *BatchInboxTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchInbox.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_BatchInbox *BatchInboxCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BatchInbox.contract.Call(opts, &out, "balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_BatchInbox *BatchInboxSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _BatchInbox.Contract.Balances(&_BatchInbox.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) view returns(uint256)
func (_BatchInbox *BatchInboxCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _BatchInbox.Contract.Balances(&_BatchInbox.CallOpts, arg0)
}

// EsStorageContract is a free data retrieval call binding the contract method 0x22eb767d.
//
// Solidity: function esStorageContract() view returns(address)
func (_BatchInbox *BatchInboxCaller) EsStorageContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BatchInbox.contract.Call(opts, &out, "esStorageContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EsStorageContract is a free data retrieval call binding the contract method 0x22eb767d.
//
// Solidity: function esStorageContract() view returns(address)
func (_BatchInbox *BatchInboxSession) EsStorageContract() (common.Address, error) {
	return _BatchInbox.Contract.EsStorageContract(&_BatchInbox.CallOpts)
}

// EsStorageContract is a free data retrieval call binding the contract method 0x22eb767d.
//
// Solidity: function esStorageContract() view returns(address)
func (_BatchInbox *BatchInboxCallerSession) EsStorageContract() (common.Address, error) {
	return _BatchInbox.Contract.EsStorageContract(&_BatchInbox.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _to) payable returns()
func (_BatchInbox *BatchInboxTransactor) Deposit(opts *bind.TransactOpts, _to common.Address) (*types.Transaction, error) {
	return _BatchInbox.contract.Transact(opts, "deposit", _to)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _to) payable returns()
func (_BatchInbox *BatchInboxSession) Deposit(_to common.Address) (*types.Transaction, error) {
	return _BatchInbox.Contract.Deposit(&_BatchInbox.TransactOpts, _to)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address _to) payable returns()
func (_BatchInbox *BatchInboxTransactorSession) Deposit(_to common.Address) (*types.Transaction, error) {
	return _BatchInbox.Contract.Deposit(&_BatchInbox.TransactOpts, _to)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _to, uint256 _amount) returns()
func (_BatchInbox *BatchInboxTransactor) Withdraw(opts *bind.TransactOpts, _to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BatchInbox.contract.Transact(opts, "withdraw", _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _to, uint256 _amount) returns()
func (_BatchInbox *BatchInboxSession) Withdraw(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BatchInbox.Contract.Withdraw(&_BatchInbox.TransactOpts, _to, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _to, uint256 _amount) returns()
func (_BatchInbox *BatchInboxTransactorSession) Withdraw(_to common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _BatchInbox.Contract.Withdraw(&_BatchInbox.TransactOpts, _to, _amount)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BatchInbox *BatchInboxTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _BatchInbox.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BatchInbox *BatchInboxSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BatchInbox.Contract.Fallback(&_BatchInbox.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BatchInbox *BatchInboxTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BatchInbox.Contract.Fallback(&_BatchInbox.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BatchInbox *BatchInboxTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchInbox.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BatchInbox *BatchInboxSession) Receive() (*types.Transaction, error) {
	return _BatchInbox.Contract.Receive(&_BatchInbox.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BatchInbox *BatchInboxTransactorSession) Receive() (*types.Transaction, error) {
	return _BatchInbox.Contract.Receive(&_BatchInbox.TransactOpts)
}
