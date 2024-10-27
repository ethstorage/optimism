package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

type expectedGetTxByHashCall struct {
	txHash  common.Hash
	outputs []byte
	err     error
}

func (c *expectedGetTxByHashCall) Matches(rpcMethod string, args ...interface{}) error {
	if rpcMethod != "eth_getTransactionByHash" {
		return fmt.Errorf("expected rpcMethod eth_getTransactionByHash but was %v", rpcMethod)
	}

	txhash, ok := args[0].(common.Hash)
	if !ok {
		return fmt.Errorf("arg 0 is not common.Hash")
	}

	if txhash != c.txHash {
		return fmt.Errorf("expected txHash %v but was %v", c.txHash, txhash)
	}

	return c.err
}

func (c *expectedGetTxByHashCall) Execute(t *testing.T, out interface{}) error {
	j, err := json.Marshal(hexutil.Bytes(c.outputs))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(j, out))
	return c.err
}

func (c *expectedGetTxByHashCall) String() string {
	return fmt.Sprintf("{txHash: %v, outputs: %v}", c.txHash, c.outputs)
}
