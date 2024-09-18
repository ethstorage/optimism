package batching

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching/test"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestTxCall_ToCallArgs(t *testing.T) {
	addr := common.Address{0xbd}
	testAbi, err := test.ERC20MetaData.GetAbi()
	require.NoError(t, err)
	call := NewTxCall(testAbi, common.Hash{0xcc}, "approve")
	expectedAmount := big.NewInt(1234444)
	expectedSpender := common.Address{0xcc}
	contractCall := NewContractCall(testAbi, addr, "approve", expectedSpender, expectedAmount)
	packed, err := contractCall.Pack()
	require.NoError(t, err)

	unpackedMap, err := call.DecodeTxParams(packed)
	require.NoError(t, err)
	require.Equal(t, expectedAmount, unpackedMap["amount"])
	require.Equal(t, expectedSpender, unpackedMap["spender"])

	unpacked, err := call.Unpack(packed)
	require.NoError(t, err)
	require.Equal(t, expectedSpender, unpacked.GetAddress(0))
	require.Equal(t, expectedAmount, unpacked.GetBigInt(1))
}

func TestGetTxCalldata(t *testing.T) {
	expectedSpender := common.Address{0xcc}
	expectedAmount := big.NewInt(1234444)
	txHash := common.Hash{0x11}
	addr := common.Address{0xbd}

	testAbi, err := test.ERC20MetaData.GetAbi()
	require.NoError(t, err)
	contractCall := NewContractCall(testAbi, addr, "approve", expectedSpender, expectedAmount)
	packed, err := contractCall.Pack()

	stub := test.NewRpcStub(t)
	stub.AddExpectedCall(test.NewGetTxCall(txHash, rpcblock.Latest, &packed))

	caller := NewMultiCaller(stub, DefaultBatchSize)
	txCall := NewTxCall(testAbi, txHash, "approve")
	result, err := caller.SingleCall(context.Background(), rpcblock.Latest, txCall)
	require.NoError(t, err)
	fmt.Println()
	require.Equal(t, expectedSpender, result.GetAddress(0))
	require.Equal(t, expectedAmount, result.GetBigInt(1))
}
