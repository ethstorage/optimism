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
	to      []common.Address
	block   rpcblock.Block
	outputs []types.Log
	err     error
}

func (c *expectedFilterLogsCall) Matches(rpcMethod string, args ...interface{}) error {
	if rpcMethod != "eth_getLogs" {
		return fmt.Errorf("expected rpcMethod eth_getLogs but was %v", rpcMethod)
	}

	query, ok := args[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("arg 0 is not map[string]interface{}")
	}
	topics := query["topics"].([][]common.Hash)
	if !ok {
		return fmt.Errorf("topics is not [][]common.Hash")
	}

	if !reflect.DeepEqual(topics, c.topics) {
		return fmt.Errorf("expected topics %v but was %v", c.topics, topics)
	}

	to := query["address"].([]common.Address)
	if !reflect.DeepEqual(to, c.to) {
		return fmt.Errorf("expected contract address %v but was %v", c.to, to)
	}
	return c.err
}

func (c *expectedFilterLogsCall) Execute(t *testing.T, out interface{}) error {
	j, err := json.Marshal(c.outputs)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(j, out))
	return c.err
}

func (c *expectedFilterLogsCall) String() string {
	return fmt.Sprintf("{to: %v, block: %v, outputs: %v}", c.to, c.block, c.outputs)
}

func (l *RpcStub) SetFilterLogResponse(topics [][]common.Hash, to []common.Address, block rpcblock.Block, output []types.Log) {
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
