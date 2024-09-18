package batching

import (
	"fmt"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type TxCall struct {
	Abi    *abi.ABI
	TxHash common.Hash
	Method string
}

var _ Call = (*TxCall)(nil)

func NewTxCall(abi *abi.ABI, txhash common.Hash, method string) *TxCall {
	return &TxCall{
		Abi:    abi,
		TxHash: txhash,
		Method: method,
	}
}

func (b *TxCall) ToBatchElemCreator() (BatchElementCreator, error) {
	return func(block rpcblock.Block) (any, rpc.BatchElem) {
		out := new(hexutil.Bytes)
		return out, rpc.BatchElem{
			Method: "eth_getTransactionByHash",
			Args:   []interface{}{b.TxHash, block.ArgValue()},
			Result: &out,
		}
	}, nil
}

func (c *TxCall) HandleResult(result interface{}) (*CallResult, error) {
	out, err := c.Unpack(*result.(*hexutil.Bytes))
	return out, err
}

func (c *TxCall) DecodeTxParams(data []byte) (map[string]interface{}, error) {
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

func (c *TxCall) Unpack(hex hexutil.Bytes) (*CallResult, error) {
	inputs := c.Abi.Methods[c.Method].Inputs

	out, err := inputs.UnpackValues(hex[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to unpack inputs: %w", err)
	}
	return &CallResult{out: out}, nil
}
