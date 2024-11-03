package batching

import (
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	batchingTest "github.com/ethereum-optimism/optimism/op-service/sources/batching/test"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

func TestEventLogFilter(t *testing.T) {
	addr := common.Address{0xbd}
	stub := batchingTest.NewRpcStub(t)
	owner := []common.Address{{0xaa}}
	spender := []common.Address{{0xbb}}

	testAbi, err := batchingTest.ERC20MetaData.GetAbi()
	require.NoError(t, err)
	name := "Approval"
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}
	query := [][]interface{}{ownerRule, spenderRule}
	query = append([][]interface{}{{testAbi.Events[name].ID}}, query...)

	topics, err := abi.MakeTopics(query...)
	require.NoError(t, err)

	txHash := common.Hash{0x11}
	block := rpcblock.Latest
	require.NoError(t, err)
	event := testAbi.Events[name]
	inputs := event.Inputs

	amount := big.NewInt(3)
	packedData, err := inputs.NonIndexed().Pack(amount)
	require.NoError(t, err)
	_out := []types.Log{
		{
			Address: addr,
			Topics:  []common.Hash{topics[0][0], topics[1][0], topics[2][0]},
			Data:    packedData,
			TxHash:  txHash,
		},
	}
	out := make([]interface{}, len(_out))
	for i, r := range _out {
		out[i] = r
	}

	stub.SetFilterLogResponse(topics, addr, block, _out)
	caller := NewMultiCaller(stub, DefaultBatchSize)

	filter, err := batchingTest.NewERC20Filterer(addr, caller)
	require.NoError(t, err)

	iterator, err := filter.FilterApproval(nil, owner, spender)
	require.NoError(t, err)

	res := iterator.Next()
	require.True(t, res, iterator.Error())
	require.Equal(t, _out[0].Address, iterator.Event.Raw.Address)
	require.Equal(t, _out[0].Topics, iterator.Event.Raw.Topics)
	require.Equal(t, _out[0].Data, iterator.Event.Raw.Data)
	require.Equal(t, txHash, iterator.Event.Raw.TxHash)
}
