package solver

import (
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/test"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/stretchr/testify/require"
)

func TestHonestClaimTracker_RootClaim(t *testing.T) {
	tracker := newHonestClaimTracker()
	nBits := uint64(1)
	splitDepth := types.Depth(3)
	builder := test.NewAlphabetClaimBuilder2(t, big.NewInt(3), 4, nBits, splitDepth)

	claim := builder.Seq().Get()
	require.False(t, tracker.IsHonest(claim))

	tracker.AddHonestClaim(types.Claim{}, claim)
	require.True(t, tracker.IsHonest(claim))
}

func TestHonestClaimTracker_ChildClaim(t *testing.T) {
	tracker := newHonestClaimTracker()
	nBits := uint64(1)
	splitDepth := types.Depth(3)
	builder := test.NewAlphabetClaimBuilder2(t, big.NewInt(3), 4, nBits, splitDepth)

	seq := builder.Seq().Attack2(nil, 0).Attack2(nil, 1)
	parent := seq.Get()
	child := seq.Attack2(nil, 0).Get()
	require.Zero(t, child.ContractIndex, "should work for claims that are not in the game state yet")

	tracker.AddHonestClaim(parent, child)
	require.False(t, tracker.IsHonest(parent))
	require.True(t, tracker.IsHonest(child))
	counter, ok := tracker.HonestCounter(parent)
	require.True(t, ok)
	require.Equal(t, child, counter)
}
