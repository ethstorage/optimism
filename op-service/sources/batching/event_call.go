package batching

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type EventCall struct {
	query ethereum.FilterQuery
}

func NewEventCall(q ethereum.FilterQuery) *EventCall {
	return &EventCall{query: q}
}

func (b *EventCall) ToBatchElemCreator() (BatchElementCreator, error) {
	arg, err := toFilterArg(b.query)
	if err != nil {
		return nil, err
	}
	return func(block rpcblock.Block) (any, rpc.BatchElem) {
		out := new([]types.Log)
		return out, rpc.BatchElem{
			Method: "eth_getLogs",
			Args:   []interface{}{arg},
			Result: &out,
		}
	}, nil
}

func (c *EventCall) HandleResult(result interface{}) (*CallResult, error) {
	res := result.(*[]types.Log)
	return &CallResult{out: []interface{}{*res}}, nil
}

func toFilterArg(q ethereum.FilterQuery) (interface{}, error) {
	arg := map[string]interface{}{
		"address": q.Addresses,
		"topics":  q.Topics,
	}
	if q.BlockHash != nil {
		arg["blockHash"] = *q.BlockHash
		if q.FromBlock != nil || q.ToBlock != nil {
			return nil, errors.New("cannot specify both BlockHash and FromBlock/ToBlock")
		}
	} else {
		if q.FromBlock == nil {
			arg["fromBlock"] = "0x0"
		} else {
			arg["fromBlock"] = toBlockNumArg(q.FromBlock)
		}
		arg["toBlock"] = toBlockNumArg(q.ToBlock)
	}
	return arg, nil
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	if number.IsInt64() {
		return rpc.BlockNumber(number.Int64()).String()
	}
	// It's negative and large, which is invalid.
	return fmt.Sprintf("<invalid %d>", number)
}
