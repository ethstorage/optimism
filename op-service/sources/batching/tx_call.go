package batching

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type TxGetByHashCall struct {
	Abi    *abi.ABI
	TxHash common.Hash
	Method string
}

func NewTxGetByHash(abi *abi.ABI, txhash common.Hash, method string) *TxGetByHashCall {
	return &TxGetByHashCall{
		Abi:    abi,
		TxHash: txhash,
		Method: method,
	}
}

func (b *TxGetByHashCall) ToBatchElemCreator() (BatchElementCreator, error) {
	return func(block rpcblock.Block) (any, rpc.BatchElem) {
		out := new(types.Transaction)
		return out, rpc.BatchElem{
			Method: "eth_getTransactionByHash",
			Args:   []interface{}{b.TxHash},
			Result: &out,
		}
	}, nil
}

func (c *TxGetByHashCall) HandleResult(result interface{}) (*CallResult, error) {
	txn, ok := result.(*types.Transaction)
	if !ok {
		return nil, fmt.Errorf("result is not types.Transaction")
	}
	return &CallResult{out: []interface{}{txn}}, nil
}

func (c *TxGetByHashCall) UnpackCallData(txn *types.Transaction) (map[string]interface{}, error) {
	data := txn.Data()
	m, err := c.Abi.MethodById(data[:4])
	v := map[string]interface{}{}
	if err != nil {
		return map[string]interface{}{}, err
	}
	if err := m.Inputs.UnpackIntoMap(v, data[4:]); err != nil {
		return map[string]interface{}{}, err
	}
	return v, nil
}
