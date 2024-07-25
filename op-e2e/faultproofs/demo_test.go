package faultproofs

import (
	"context"
	"math/big"
	"testing"

	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/disputegame"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/stretchr/testify/require"
)

func TestSendDemoTx(t *testing.T) {
	op_e2e.InitParallel(t)
	ctx := context.Background()
	sys, l1Client := StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)

	demoContract := disputegame.NewDemoHelper(t, ctx, sys)
	num := big.NewInt(3)
	demoContract.Store(ctx, num)
	require.NoError(t, wait.ForNextBlock(ctx, l1Client))

	_num, _ := demoContract.Retrieve(ctx)
	require.Equal(t, num, _num)
}
