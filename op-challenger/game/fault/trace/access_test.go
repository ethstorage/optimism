package trace

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/test"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/alphabet"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestAccessor_UsesSelector(t *testing.T) {
	ctx := context.Background()
	depth := types.Depth(4)
	provider1 := test.NewAlphabetWithProofProvider(t, big.NewInt(0), depth, nil)
	provider2 := alphabet.NewTraceProvider(big.NewInt(0), depth)
	claim := types.Claim{}
	game := types.NewGameState([]types.Claim{claim}, depth)
	pos1 := types.NewPositionFromGIndex(big.NewInt(4))
	pos2 := types.NewPositionFromGIndex(big.NewInt(6))

	accessor := &Accessor{
		selector: func(ctx context.Context, actualGame types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error) {
			require.Equal(t, game, actualGame)
			require.Equal(t, claim, ref)

			if pos == pos1 {
				return provider1, nil
			} else if pos == pos2 {
				return provider2, nil
			}
			return nil, fmt.Errorf("incorrect position requested: %v", pos)
		},
	}

	t.Run("Get", func(t *testing.T) {
		actual, err := accessor.Get(ctx, game, claim, pos1)
		require.NoError(t, err)
		expected, err := provider1.Get(ctx, pos1)
		require.NoError(t, err)
		require.Equal(t, expected, actual)

		actual, err = accessor.Get(ctx, game, claim, pos2)
		require.NoError(t, err)
		expected, err = provider2.Get(ctx, pos2)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("GetStepData", func(t *testing.T) {
		actualPrestate, actualProofData, actualPreimageData, err := accessor.GetStepData(ctx, game, claim, pos1)
		require.NoError(t, err)
		expectedPrestate, expectedProofData, expectedPreimageData, err := provider1.GetStepData(ctx, pos1)
		require.NoError(t, err)
		require.Equal(t, expectedPrestate, actualPrestate)
		require.Equal(t, expectedProofData, actualProofData)
		require.Equal(t, expectedPreimageData, actualPreimageData)

		actualPrestate, actualProofData, actualPreimageData, err = accessor.GetStepData(ctx, game, claim, pos2)
		require.NoError(t, err)
		expectedPrestate, expectedProofData, expectedPreimageData, err = provider2.GetStepData(ctx, pos2)
		require.NoError(t, err)
		require.Equal(t, expectedPrestate, actualPrestate)
		require.Equal(t, expectedProofData, actualProofData)
		require.Equal(t, expectedPreimageData, actualPreimageData)
	})

	t.Run("GetL2BlockNumberChallenge", func(t *testing.T) {
		provider := &ChallengeTraceProvider{
			TraceProvider: provider1,
		}
		accessor := &Accessor{
			selector: func(ctx context.Context, actualGame types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error) {
				require.Equal(t, game, actualGame)
				require.Equal(t, game.Claims()[0], ref)
				require.Equal(t, types.RootPosition, pos)
				return provider, nil
			},
		}
		challenge, err := accessor.GetL2BlockNumberChallenge(ctx, game)
		require.NoError(t, err)
		require.NotNil(t, challenge)
		require.Equal(t, eth.Bytes32{0xaa, 0xbb}, challenge.Output.OutputRoot)
	})
}

func TestGetStepData2(t *testing.T) {
	ctx := context.Background()
	traceDepth := types.Depth(4)
	splitDepth := types.Depth(2)
	nbits := uint64(2)
	nary := int64(1) << nbits
	maxDepth := traceDepth + splitDepth + types.Depth(nbits)
	provider := test.NewAlphabetWithProofProvider2(t, big.NewInt(0), traceDepth, splitDepth, nil)
	translatedTracerovider := Translate(provider, splitDepth+types.Depth(nary))

	claimBuilder := test.NewAlphabetClaimBuilder(t, big.NewInt(0), maxDepth)
	gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
	seq := gameBuilder.Seq()
	claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], nbits, 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], nbits, 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], nbits, 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], nbits, 1)

	game := gameBuilder.Game
	claim := game.Claims()[len(game.Claims())-1]

	pos := claim.Position.MoveRight()

	accessor := &Accessor{
		selector: func(ctx context.Context, actualGame types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error) {
			require.Equal(t, game, actualGame)
			require.Equal(t, claim, ref)
			return translatedTracerovider, nil
		},
	}

	expectedStateDA := types.DAData{
		Prestate: types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[0][:],
			Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
		},
		PostState: types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[1][:],
			Proof:    append(subClaimsDep8[0][:], subClaimsDep8[2][:]...),
		},
	}

	expectedPrestate, expectedProofData, expectedPreimageData, err := translatedTracerovider.GetStepData(ctx, pos)
	expectedPreimageData.DAData = expectedStateDA
	require.NoError(t, err)

	actualPrestate, actualProofData, actualPreimageData, err := accessor.GetStepData2(ctx, game, claim, pos)
	require.NoError(t, err)

	require.Equal(t, expectedPrestate, actualPrestate)
	require.Equal(t, expectedProofData, actualProofData)
	require.Equal(t, expectedPreimageData, actualPreimageData)
}

func TestFindAncestorProofAtDepth2(t *testing.T) {
	nbits := uint64(2)
	nary := int64(1) << nbits
	maxGameDepth := types.Depth(8)
	splitDepth := types.Depth(2)
	claimBuilder := test.NewAlphabetClaimBuilder(t, big.NewInt(0), maxGameDepth)
	traceDepth := maxGameDepth - splitDepth - types.Depth(nbits)
	ctx := context.Background()
	provider := test.NewAlphabetWithProofProvider(t, big.NewInt(0), traceDepth, nil)

	{ // Attack the left most claim
		depth8Branch := uint64(0)
		stepBranch := int64(0)

		// output roots
		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)
		// traces
		subClaimsDep4 := []common.Hash{{0x04}}
		seq = seq.Attack2(subClaimsDep4[:], nbits, 0)
		subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
		seq = seq.Attack2(subClaimsDep6[:], nbits, 0)
		subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
		seq.Attack2(subClaimsDep8[:], nbits, depth8Branch)

		game := gameBuilder.Game
		ref := game.Claims()[4]
		pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
		traceRootDepth := game.RootDepth()
		relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
		require.NoError(t, err)
		postTraceIdx := relativePos.TraceIndex(traceDepth)
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))

		absolutePre, err := provider.AbsolutePreStateCommitment(ctx)
		require.NoError(t, err)
		expectedPre := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: absolutePre.Bytes(),
			Proof:    []byte{},
		}
		preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, preStateDaItem)

		expectedPost := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[0][:],
			Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
		}
		postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, postStateDaItem)
	}

	{ // Attack the left claim
		depth8Branch := uint64(1)
		stepBranch := int64(0)

		// output roots
		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)
		// traces
		subClaimsDep4 := []common.Hash{{0x04}}
		seq = seq.Attack2(subClaimsDep4[:], nbits, 0)
		subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
		seq = seq.Attack2(subClaimsDep6[:], nbits, 0)
		subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
		seq.Attack2(subClaimsDep8[:], nbits, depth8Branch)

		game := gameBuilder.Game
		ref := game.Claims()[4]
		pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
		traceRootDepth := game.RootDepth()
		relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
		require.NoError(t, err)
		postTraceIdx := relativePos.TraceIndex(traceDepth)
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))

		expectedPre := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep6[0][:],
			Proof:    append(subClaimsDep6[1][:], subClaimsDep6[2][:]...),
		}
		preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, preStateDaItem)

		expectedPost := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[0][:],
			Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
		}
		postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, postStateDaItem)
	}

	{ // Attack the mid claim
		depth8Branch := uint64(1)
		stepBranch := int64(1)

		// output roots
		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)
		// traces
		subClaimsDep4 := []common.Hash{{0x04}}
		seq = seq.Attack2(subClaimsDep4[:], nbits, 0)
		subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
		seq = seq.Attack2(subClaimsDep6[:], nbits, 0)
		subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
		seq.Attack2(subClaimsDep8[:], nbits, depth8Branch)

		game := gameBuilder.Game
		ref := game.Claims()[4]
		pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
		traceRootDepth := game.RootDepth()
		relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
		require.NoError(t, err)
		postTraceIdx := relativePos.TraceIndex(traceDepth)
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))

		expectedPre := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[0][:],
			Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
		}
		preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, preStateDaItem)

		expectedPost := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[1][:],
			Proof:    append(subClaimsDep8[0][:], subClaimsDep8[2][:]...),
		}
		postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, postStateDaItem)
	}

	{ // Attack the right claim
		depth8Branch := uint64(1)
		stepBranch := int64(3)

		// output roots
		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)
		// traces
		subClaimsDep4 := []common.Hash{{0x04}}
		seq = seq.Attack2(subClaimsDep4[:], nbits, 0)
		subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
		seq = seq.Attack2(subClaimsDep6[:], nbits, 0)
		subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
		seq.Attack2(subClaimsDep8[:], nbits, depth8Branch)

		game := gameBuilder.Game
		ref := game.Claims()[4]
		pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
		traceRootDepth := game.RootDepth()
		relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
		require.NoError(t, err)
		postTraceIdx := relativePos.TraceIndex(traceDepth)
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))

		expectedPre := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[2][:],
			Proof:    append(subClaimsDep8[0][:], subClaimsDep8[1][:]...),
		}
		preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, preStateDaItem)

		expectedPost := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep6[1][:],
			Proof:    append(subClaimsDep6[0][:], subClaimsDep6[2][:]...),
		}
		postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, postStateDaItem)
	}

	{ // Attack the right most claim
		depth8Branch := uint64(3)
		stepBranch := int64(3)

		// output roots
		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)
		// traces
		subClaimsDep4 := []common.Hash{{0x04}}
		seq = seq.Attack2(subClaimsDep4[:], nbits, 0)
		subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
		seq = seq.Attack2(subClaimsDep6[:], nbits, 0)
		subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
		seq.Attack2(subClaimsDep8[:], nbits, depth8Branch)

		game := gameBuilder.Game
		ref := game.Claims()[4]
		pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
		traceRootDepth := game.RootDepth()
		relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
		require.NoError(t, err)
		postTraceIdx := relativePos.TraceIndex(traceDepth)
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))

		expectedPre := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: subClaimsDep8[2][:],
			Proof:    append(subClaimsDep8[0][:], subClaimsDep8[1][:]...),
		}
		preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, preStateDaItem)

		splitLeaf := game.Claims()[2]
		expectedPost := types.DAItem{
			DaType:   types.CallDataType,
			DataHash: splitLeaf.Value[:],
			Proof:    []byte{},
		}
		postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, postStateDaItem)
	}

}

type ChallengeTraceProvider struct {
	types.TraceProvider
}

func (c *ChallengeTraceProvider) GetL2BlockNumberChallenge(_ context.Context) (*types.InvalidL2BlockNumberChallenge, error) {
	return &types.InvalidL2BlockNumberChallenge{
		Output: &eth.OutputResponse{OutputRoot: eth.Bytes32{0xaa, 0xbb}},
	}, nil
}
