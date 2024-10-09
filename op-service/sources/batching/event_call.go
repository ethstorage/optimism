package batching

import (
	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type EventCall struct {
	topics [][]common.Hash
	to     []common.Address
}

func NewEventCall(q ethereum.FilterQuery) *EventCall {
	return &EventCall{
		topics: q.Topics,
		to:     q.Addresses,
	}
}

func (b *EventCall) ToBatchElemCreator() (BatchElementCreator, error) {
	return func(block rpcblock.Block) (any, rpc.BatchElem) {
		out := new([]types.Log)
		return out, rpc.BatchElem{
			Method: "eth_getFilterLogs",
			Args:   []interface{}{b.topics, b.to[0], block.ArgValue()},
			Result: &out,
		}
	}, nil
}

func (c *EventCall) HandleResult(result interface{}) (*CallResult, error) {
	res := result.(*[]types.Log)
	return &CallResult{out: []interface{}{*res}}, nil
}
