package batching

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching/test"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestDecodeTxGetByHash(t *testing.T) {
	addr := common.Address{0xbd}
	testAbi, err := test.ERC20MetaData.GetAbi()
	require.NoError(t, err)
	call := NewTxGetByHash(testAbi, common.Hash{0xcc}, "approve")
	expectedAmount := big.NewInt(1234444)
	expectedSpender := common.Address{0xcc}
	contractCall := NewContractCall(testAbi, addr, "approve", expectedSpender, expectedAmount)
	packed, err := contractCall.Pack()
	require.NoError(t, err)

	unpackedMap, err := call.DecodeTxParams(packed)
	require.NoError(t, err)
	require.Equal(t, expectedAmount, unpackedMap["amount"])
	require.Equal(t, expectedSpender, unpackedMap["spender"])
}

func TestUnpackTxCalldata(t *testing.T) {
	expectedSpender := common.Address{0xcc}
	expectedAmount := big.NewInt(1234444)
	txHash := common.Hash{0x11}
	addr := common.Address{0xbd}

	testAbi, err := test.ERC20MetaData.GetAbi()
	require.NoError(t, err)
	contractCall := NewContractCall(testAbi, addr, "approve", expectedSpender, expectedAmount)
	inputData, err := contractCall.Pack()
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		GasPrice: big.NewInt(11111),
		Gas:      1111,
		To:       &addr,
		Value:    big.NewInt(111),
		Data:     inputData,
	})
	require.NoError(t, err)
	packed, err := tx.MarshalBinary()
	require.NoError(t, err)

	stub := test.NewRpcStub(t)
	stub.AddExpectedCall(test.NewGetTxCall(txHash, rpcblock.Latest, &packed))

	caller := NewMultiCaller(stub, DefaultBatchSize)
	txCall := NewTxGetByHash(testAbi, txHash, "approve")
	result, err := caller.SingleCall(context.Background(), rpcblock.Latest, txCall)
	require.NoError(t, err)

	decodedTx, err := txCall.DecodeToTx(result)
	require.NoError(t, err)
	unpackedMap, err := txCall.UnpackCallData(decodedTx)
	require.NoError(t, err)
	require.Equal(t, expectedSpender, unpackedMap["spender"])
	require.Equal(t, expectedAmount, unpackedMap["amount"])
}
