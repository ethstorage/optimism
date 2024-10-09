package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

type expectedTxCall struct {
	txHash  common.Hash
	outputs []byte
	err     error
}

func (c *expectedTxCall) Matches(rpcMethod string, args ...interface{}) error {
	if rpcMethod != "eth_getTransactionByHash" {
		return fmt.Errorf("expected rpcMethod eth_getTransactionByHash but was %v", rpcMethod)
	}

	txhash := args[0].(common.Hash)

	if txhash != c.txHash {
		return fmt.Errorf("expected txHash %v but was %v", c.txHash, txhash)
	}

	return c.err
}

func (c *expectedTxCall) Execute(t *testing.T, out interface{}) error {
	j, err := json.Marshal(hexutil.Bytes(c.outputs))
	require.NoError(t, err)
	json.Unmarshal(j, out)
	return c.err
}

func (c *expectedTxCall) String() string {
	return fmt.Sprintf("{txHash: %v, outputs: %v}", c.txHash, c.outputs)
}
