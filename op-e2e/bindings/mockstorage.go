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

// MockEthStorageMetaData contains all meta data concerning the MockEthStorage contract.
var MockEthStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_cost\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"kvIdx\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"kvSize\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"PutBlob\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"kvEntryCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_key\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_blobIdx\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_length\",\"type\":\"uint256\"}],\"name\":\"putBlob\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upfrontPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a0604052348015600e575f5ffd5b506040516104fc3803806104fc8339818101604052810190602e9190606d565b8060808181525050506093565b5f5ffd5b5f819050919050565b604f81603f565b81146058575f5ffd5b50565b5f815190506067816048565b92915050565b5f60208284031215607f57607e603b565b5b5f608a84828501605b565b91505092915050565b6080516104526100aa5f395f60aa01526104525ff3fe608060405260043610610033575f3560e01c80631ccbc6da146100375780634581a92014610061578063638ba9e91461007d575b5f5ffd5b348015610042575f5ffd5b5061004b6100a7565b60405161005891906101c6565b60405180910390f35b61007b60048036038101906100769190610240565b6100ce565b005b348015610088575f5ffd5b506100916101a9565b60405161009e91906101c6565b60405180910390f35b5f7f0000000000000000000000000000000000000000000000000000000000000000905090565b5f824990505f5f1b8103610117576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161010e90610310565b60405180910390fd5b5f5f1b840361015b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101529061039e565b60405180910390fd5b5f5f54905060015f5461016e91906103e9565b5f819055508183827f8b7a21215282409938287ae262331bfe6411d35d3d46aa7e505ef02000524ac260405160405180910390a45050505050565b5f5481565b5f819050919050565b6101c0816101ae565b82525050565b5f6020820190506101d95f8301846101b7565b92915050565b5f5ffd5b5f819050919050565b6101f5816101e3565b81146101ff575f5ffd5b50565b5f81359050610210816101ec565b92915050565b61021f816101ae565b8114610229575f5ffd5b50565b5f8135905061023a81610216565b92915050565b5f5f5f60608486031215610257576102566101df565b5b5f61026486828701610202565b93505060206102758682870161022c565b92505060406102868682870161022c565b9150509250925092565b5f82825260208201905092915050565b7f45746853746f72616765436f6e74726163743a206661696c656420746f2067655f8201527f7420626c6f622068617368000000000000000000000000000000000000000000602082015250565b5f6102fa602b83610290565b9150610305826102a0565b604082019050919050565b5f6020820190508181035f830152610327816102ee565b9050919050565b7f45746853746f72616765436f6e74726163743a206661696c656420746f2067655f8201527f7420626c6f62206b657900000000000000000000000000000000000000000000602082015250565b5f610388602a83610290565b91506103938261032e565b604082019050919050565b5f6020820190508181035f8301526103b58161037c565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6103f3826101ae565b91506103fe836101ae565b9250828201905080821115610416576104156103bc565b5b9291505056fea264697066735822122063d49cb27e579c573ed1bf982b7c8ed08b4425e513d5b8c01c50402ab28e701964736f6c634300081c0033",
}

// MockEthStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use MockEthStorageMetaData.ABI instead.
var MockEthStorageABI = MockEthStorageMetaData.ABI

// MockEthStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockEthStorageMetaData.Bin instead.
var MockEthStorageBin = MockEthStorageMetaData.Bin

// DeployMockEthStorage deploys a new Ethereum contract, binding an instance of MockEthStorage to it.
func DeployMockEthStorage(auth *bind.TransactOpts, backend bind.ContractBackend, _cost *big.Int) (common.Address, *types.Transaction, *MockEthStorage, error) {
	parsed, err := MockEthStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockEthStorageBin), backend, _cost)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MockEthStorage{MockEthStorageCaller: MockEthStorageCaller{contract: contract}, MockEthStorageTransactor: MockEthStorageTransactor{contract: contract}, MockEthStorageFilterer: MockEthStorageFilterer{contract: contract}}, nil
}

// MockEthStorage is an auto generated Go binding around an Ethereum contract.
type MockEthStorage struct {
	MockEthStorageCaller     // Read-only binding to the contract
	MockEthStorageTransactor // Write-only binding to the contract
	MockEthStorageFilterer   // Log filterer for contract events
}

// MockEthStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockEthStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockEthStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockEthStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockEthStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockEthStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockEthStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockEthStorageSession struct {
	Contract     *MockEthStorage   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MockEthStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockEthStorageCallerSession struct {
	Contract *MockEthStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// MockEthStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockEthStorageTransactorSession struct {
	Contract     *MockEthStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// MockEthStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockEthStorageRaw struct {
	Contract *MockEthStorage // Generic contract binding to access the raw methods on
}

// MockEthStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockEthStorageCallerRaw struct {
	Contract *MockEthStorageCaller // Generic read-only contract binding to access the raw methods on
}

// MockEthStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockEthStorageTransactorRaw struct {
	Contract *MockEthStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockEthStorage creates a new instance of MockEthStorage, bound to a specific deployed contract.
func NewMockEthStorage(address common.Address, backend bind.ContractBackend) (*MockEthStorage, error) {
	contract, err := bindMockEthStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockEthStorage{MockEthStorageCaller: MockEthStorageCaller{contract: contract}, MockEthStorageTransactor: MockEthStorageTransactor{contract: contract}, MockEthStorageFilterer: MockEthStorageFilterer{contract: contract}}, nil
}

// NewMockEthStorageCaller creates a new read-only instance of MockEthStorage, bound to a specific deployed contract.
func NewMockEthStorageCaller(address common.Address, caller bind.ContractCaller) (*MockEthStorageCaller, error) {
	contract, err := bindMockEthStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockEthStorageCaller{contract: contract}, nil
}

// NewMockEthStorageTransactor creates a new write-only instance of MockEthStorage, bound to a specific deployed contract.
func NewMockEthStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*MockEthStorageTransactor, error) {
	contract, err := bindMockEthStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockEthStorageTransactor{contract: contract}, nil
}

// NewMockEthStorageFilterer creates a new log filterer instance of MockEthStorage, bound to a specific deployed contract.
func NewMockEthStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*MockEthStorageFilterer, error) {
	contract, err := bindMockEthStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockEthStorageFilterer{contract: contract}, nil
}

// bindMockEthStorage binds a generic wrapper to an already deployed contract.
func bindMockEthStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockEthStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockEthStorage *MockEthStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockEthStorage.Contract.MockEthStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockEthStorage *MockEthStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockEthStorage.Contract.MockEthStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockEthStorage *MockEthStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockEthStorage.Contract.MockEthStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockEthStorage *MockEthStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockEthStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockEthStorage *MockEthStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockEthStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockEthStorage *MockEthStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockEthStorage.Contract.contract.Transact(opts, method, params...)
}

// KvEntryCount is a free data retrieval call binding the contract method 0x638ba9e9.
//
// Solidity: function kvEntryCount() view returns(uint256)
func (_MockEthStorage *MockEthStorageCaller) KvEntryCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockEthStorage.contract.Call(opts, &out, "kvEntryCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// KvEntryCount is a free data retrieval call binding the contract method 0x638ba9e9.
//
// Solidity: function kvEntryCount() view returns(uint256)
func (_MockEthStorage *MockEthStorageSession) KvEntryCount() (*big.Int, error) {
	return _MockEthStorage.Contract.KvEntryCount(&_MockEthStorage.CallOpts)
}

// KvEntryCount is a free data retrieval call binding the contract method 0x638ba9e9.
//
// Solidity: function kvEntryCount() view returns(uint256)
func (_MockEthStorage *MockEthStorageCallerSession) KvEntryCount() (*big.Int, error) {
	return _MockEthStorage.Contract.KvEntryCount(&_MockEthStorage.CallOpts)
}

// UpfrontPayment is a free data retrieval call binding the contract method 0x1ccbc6da.
//
// Solidity: function upfrontPayment() view returns(uint256)
func (_MockEthStorage *MockEthStorageCaller) UpfrontPayment(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockEthStorage.contract.Call(opts, &out, "upfrontPayment")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UpfrontPayment is a free data retrieval call binding the contract method 0x1ccbc6da.
//
// Solidity: function upfrontPayment() view returns(uint256)
func (_MockEthStorage *MockEthStorageSession) UpfrontPayment() (*big.Int, error) {
	return _MockEthStorage.Contract.UpfrontPayment(&_MockEthStorage.CallOpts)
}

// UpfrontPayment is a free data retrieval call binding the contract method 0x1ccbc6da.
//
// Solidity: function upfrontPayment() view returns(uint256)
func (_MockEthStorage *MockEthStorageCallerSession) UpfrontPayment() (*big.Int, error) {
	return _MockEthStorage.Contract.UpfrontPayment(&_MockEthStorage.CallOpts)
}

// PutBlob is a paid mutator transaction binding the contract method 0x4581a920.
//
// Solidity: function putBlob(bytes32 _key, uint256 _blobIdx, uint256 _length) payable returns()
func (_MockEthStorage *MockEthStorageTransactor) PutBlob(opts *bind.TransactOpts, _key [32]byte, _blobIdx *big.Int, _length *big.Int) (*types.Transaction, error) {
	return _MockEthStorage.contract.Transact(opts, "putBlob", _key, _blobIdx, _length)
}

// PutBlob is a paid mutator transaction binding the contract method 0x4581a920.
//
// Solidity: function putBlob(bytes32 _key, uint256 _blobIdx, uint256 _length) payable returns()
func (_MockEthStorage *MockEthStorageSession) PutBlob(_key [32]byte, _blobIdx *big.Int, _length *big.Int) (*types.Transaction, error) {
	return _MockEthStorage.Contract.PutBlob(&_MockEthStorage.TransactOpts, _key, _blobIdx, _length)
}

// PutBlob is a paid mutator transaction binding the contract method 0x4581a920.
//
// Solidity: function putBlob(bytes32 _key, uint256 _blobIdx, uint256 _length) payable returns()
func (_MockEthStorage *MockEthStorageTransactorSession) PutBlob(_key [32]byte, _blobIdx *big.Int, _length *big.Int) (*types.Transaction, error) {
	return _MockEthStorage.Contract.PutBlob(&_MockEthStorage.TransactOpts, _key, _blobIdx, _length)
}

// MockEthStoragePutBlobIterator is returned from FilterPutBlob and is used to iterate over the raw logs and unpacked data for PutBlob events raised by the MockEthStorage contract.
type MockEthStoragePutBlobIterator struct {
	Event *MockEthStoragePutBlob // Event containing the contract specifics and raw log

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
func (it *MockEthStoragePutBlobIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockEthStoragePutBlob)
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
		it.Event = new(MockEthStoragePutBlob)
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
func (it *MockEthStoragePutBlobIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockEthStoragePutBlobIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockEthStoragePutBlob represents a PutBlob event raised by the MockEthStorage contract.
type MockEthStoragePutBlob struct {
	KvIdx    *big.Int
	KvSize   *big.Int
	DataHash [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPutBlob is a free log retrieval operation binding the contract event 0x8b7a21215282409938287ae262331bfe6411d35d3d46aa7e505ef02000524ac2.
//
// Solidity: event PutBlob(uint256 indexed kvIdx, uint256 indexed kvSize, bytes32 indexed dataHash)
func (_MockEthStorage *MockEthStorageFilterer) FilterPutBlob(opts *bind.FilterOpts, kvIdx []*big.Int, kvSize []*big.Int, dataHash [][32]byte) (*MockEthStoragePutBlobIterator, error) {

	var kvIdxRule []interface{}
	for _, kvIdxItem := range kvIdx {
		kvIdxRule = append(kvIdxRule, kvIdxItem)
	}
	var kvSizeRule []interface{}
	for _, kvSizeItem := range kvSize {
		kvSizeRule = append(kvSizeRule, kvSizeItem)
	}
	var dataHashRule []interface{}
	for _, dataHashItem := range dataHash {
		dataHashRule = append(dataHashRule, dataHashItem)
	}

	logs, sub, err := _MockEthStorage.contract.FilterLogs(opts, "PutBlob", kvIdxRule, kvSizeRule, dataHashRule)
	if err != nil {
		return nil, err
	}
	return &MockEthStoragePutBlobIterator{contract: _MockEthStorage.contract, event: "PutBlob", logs: logs, sub: sub}, nil
}

// WatchPutBlob is a free log subscription operation binding the contract event 0x8b7a21215282409938287ae262331bfe6411d35d3d46aa7e505ef02000524ac2.
//
// Solidity: event PutBlob(uint256 indexed kvIdx, uint256 indexed kvSize, bytes32 indexed dataHash)
func (_MockEthStorage *MockEthStorageFilterer) WatchPutBlob(opts *bind.WatchOpts, sink chan<- *MockEthStoragePutBlob, kvIdx []*big.Int, kvSize []*big.Int, dataHash [][32]byte) (event.Subscription, error) {

	var kvIdxRule []interface{}
	for _, kvIdxItem := range kvIdx {
		kvIdxRule = append(kvIdxRule, kvIdxItem)
	}
	var kvSizeRule []interface{}
	for _, kvSizeItem := range kvSize {
		kvSizeRule = append(kvSizeRule, kvSizeItem)
	}
	var dataHashRule []interface{}
	for _, dataHashItem := range dataHash {
		dataHashRule = append(dataHashRule, dataHashItem)
	}

	logs, sub, err := _MockEthStorage.contract.WatchLogs(opts, "PutBlob", kvIdxRule, kvSizeRule, dataHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockEthStoragePutBlob)
				if err := _MockEthStorage.contract.UnpackLog(event, "PutBlob", log); err != nil {
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

// ParsePutBlob is a log parse operation binding the contract event 0x8b7a21215282409938287ae262331bfe6411d35d3d46aa7e505ef02000524ac2.
//
// Solidity: event PutBlob(uint256 indexed kvIdx, uint256 indexed kvSize, bytes32 indexed dataHash)
func (_MockEthStorage *MockEthStorageFilterer) ParsePutBlob(log types.Log) (*MockEthStoragePutBlob, error) {
	event := new(MockEthStoragePutBlob)
	if err := _MockEthStorage.contract.UnpackLog(event, "PutBlob", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
