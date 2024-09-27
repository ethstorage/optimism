package test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

type expectedFilterLogsCall struct {
	topics  [][]common.Hash
	to      common.Address
	block   rpcblock.Block
	outputs []types.Log
	err     error
}

func (c *expectedFilterLogsCall) Matches(rpcMethod string, args ...interface{}) error {
	if rpcMethod != "eth_getFilterLogs" {
		return fmt.Errorf("expected rpcMethod eth_getFilterLogs but was %v", rpcMethod)
	}

	topics := args[0].([][]common.Hash)

	if !reflect.DeepEqual(topics, c.topics) {
		return fmt.Errorf("expected topics %v but was %v", c.topics, topics)
	}

	to := args[1].(common.Address)
	if to != c.to {
		return fmt.Errorf("expected contract address %v but was %v", c.topics, topics)
	}
	return c.err
}

func (c *expectedFilterLogsCall) Execute(t *testing.T, out interface{}) error {
	j, err := json.Marshal((c.outputs))
	require.NoError(t, err)
	json.Unmarshal(j, out)
	return c.err
}

func (c *expectedFilterLogsCall) String() string {
	return fmt.Sprintf("{to: %v, block: %v, outputs: %v}", c.to, c.block, c.outputs)
}

func (l *RpcStub) SetFilterLogResponse(topics [][]common.Hash, to common.Address, block rpcblock.Block, output []types.Log) {
	if output == nil {
		output = []types.Log{}
	}

	l.AddExpectedCall(&expectedFilterLogsCall{
		topics:  topics,
		to:      to,
		block:   block,
		outputs: output,
	})
}
