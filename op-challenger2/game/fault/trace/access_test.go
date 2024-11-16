package trace

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/test"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace/alphabet"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestAccessor_UsesSelector(t *testing.T) {
	ctx := context.Background()
	depth := types.Depth(4)
	provider1 := test.NewAlphabetWithProofProvider(t, big.NewInt(0), depth, nil, 0, test.OracleDefaultKey)
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

type ChallengeTraceProvider struct {
	types.TraceProvider
}

func (c *ChallengeTraceProvider) GetL2BlockNumberChallenge(_ context.Context) (*types.InvalidL2BlockNumberChallenge, error) {
	return &types.InvalidL2BlockNumberChallenge{
		Output: &eth.OutputResponse{OutputRoot: eth.Bytes32{0xaa, 0xbb}},
	}, nil
}

func TestGetStepData2(t *testing.T) {
	ctx := context.Background()
	traceDepth := types.Depth(4)
	splitDepth := types.Depth(2)
	nbits := uint64(2)
	nary := int64(1) << nbits
	maxDepth := traceDepth + splitDepth + types.Depth(nbits)
	provider := test.NewAlphabetWithProofProvider(t, big.NewInt(0), traceDepth, nil, splitDepth, test.OracleDefaultKey)
	translatedTracerovider := Translate(provider, splitDepth+types.Depth(nary), types.DAData{})

	claimBuilder := test.NewAlphabetClaimBuilder2(t, big.NewInt(0), maxDepth, nbits, splitDepth)
	gameBuilder := claimBuilder.GameBuilder()
	seq := gameBuilder.Seq()
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], 1)

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
		PreDA: types.DAItem{
			DataHash: subClaimsDep8[0],
			Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
		},
		PostDA: types.DAItem{
			DataHash: subClaimsDep8[1],
			Proof:    append(subClaimsDep8[0][:], subClaimsDep8[2][:]...),
		},
	}

	expectedPrestate, expectedProofData, expectedPreimageData, err := translatedTracerovider.GetStepData(ctx, pos)
	expectedPreimageData.VMStateDA = expectedStateDA
	require.NoError(t, err)

	actualPrestate, actualProofData, actualPreimageData, err := accessor.GetStepData2(ctx, game, claim, pos)
	require.NoError(t, err)

	require.Equal(t, expectedPrestate, actualPrestate)
	require.Equal(t, expectedProofData, actualProofData)
	require.Equal(t, expectedPreimageData, actualPreimageData)
}

func TestGetStepDataWithOutputRoot(t *testing.T) {
	ctx := context.Background()
	traceDepth := types.Depth(4)
	splitDepth := types.Depth(2)
	nbits := uint64(2)
	nary := int64(1) << nbits
	maxDepth := traceDepth + splitDepth + types.Depth(nbits)

	claimBuilder := test.NewAlphabetClaimBuilder2(t, big.NewInt(0), maxDepth, nbits, splitDepth)
	gameBuilder := claimBuilder.GameBuilder()
	seq := gameBuilder.Seq()
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], 1)

	game := gameBuilder.Game
	claim := game.Claims()[len(game.Claims())-1]

	pos := claim.Position.MoveRight()

	outputRootDA := types.DAData{
		PreDA: types.DAItem{},
		PostDA: types.DAItem{
			DataHash: subClaimsDep2[0],
			Proof:    append(subClaimsDep2[1][:], subClaimsDep2[2][:]...),
		},
	}

	tests := []struct {
		oracalKeyType  test.OracleKeyType
		expectedDAItem types.DAItem
	}{
		{test.OraclePreKey, outputRootDA.PreDA},
		{test.OraclPostKey, outputRootDA.PostDA},
	}

	for _, tCase := range tests {
		provider := test.NewAlphabetWithProofProvider(t, big.NewInt(0), traceDepth, nil, splitDepth, tCase.oracalKeyType)
		translatedTracerovider := Translate(provider, splitDepth+types.Depth(nary), outputRootDA)
		accessor := &Accessor{
			selector: func(ctx context.Context, actualGame types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error) {
				require.Equal(t, game, actualGame)
				require.Equal(t, claim, ref)
				return translatedTracerovider, nil
			},
		}

		_, _, actualPreimageData, err := accessor.GetStepData2(ctx, game, claim, pos)
		require.NoError(t, err)
		require.Equal(t, tCase.expectedDAItem, actualPreimageData.OutputRootDAItem)
	}
}

func TestFindAncestorProofAtDepth2(t *testing.T) {
	nbits := uint64(2)
	maxGameDepth := types.Depth(8)
	splitDepth := types.Depth(2)
	traceDepth := maxGameDepth - splitDepth - types.Depth(nbits)
	ctx := context.Background()
	provider := test.NewAlphabetWithProofProvider(t, big.NewInt(0), traceDepth, nil, splitDepth+types.Depth(nbits), test.OracleDefaultKey)

	tests := []struct {
		name  string
		setup func(t *testing.T, ctx context.Context, gameBuilder *test.GameBuilder, provider *test.AlphabetWithProofProvider) (game types.Game, ref types.Claim, preTraceIdx *big.Int, postTraceIdx *big.Int, expectPreDA types.DAItem, expectPostDA types.DAItem)
	}{
		// attack every claims's first branch (index at 0)
		{"attackLeftMost", attackLeftMost},
		// attack one claim's branch between [1,maxBranch-1] and leaf claims's first branch
		{"attackFirstBranch", attackFirstBranch},
		// attack one claim's branch between [1,maxBranch-1] and leaf claim's branch between [1,maxBranch-1]
		{"attackMidBranch", attackMidBranch},
		// attack one claim's branch between [1,maxBranch-1] and leaf claim's last branch (index at maxBranch)
		{"attackMaxBranch", attackMaxBranch},
		// attack every claims's last branch
		{"attackRightMost", attackRightMost},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			claimBuilder := test.NewAlphabetClaimBuilder2(t, big.NewInt(0), maxGameDepth, nbits, splitDepth)
			game, ref, preTraceIdx, postTraceIdx, expectPreDA, expectPostDA := tCase.setup(t, ctx, claimBuilder.GameBuilder(), provider)
			preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
			require.NoError(t, err)
			require.Equal(t, expectPreDA, preStateDaItem)
			postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
			require.NoError(t, err)
			require.Equal(t, expectPostDA, postStateDaItem)
		})
	}
}

func attackLeftMost(t *testing.T, ctx context.Context, gameBuilder *test.GameBuilder, provider *test.AlphabetWithProofProvider) (game types.Game, ref types.Claim, preTraceIdx *big.Int, postTraceIdx *big.Int, expectPreDA types.DAItem, expectPostDA types.DAItem) {
	depth8Branch := uint64(0)
	stepBranch := int64(0)
	// output roots
	seq := gameBuilder.Seq()
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], depth8Branch)

	game = gameBuilder.Game
	ref = game.Claims()[4]
	pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
	traceRootDepth := game.TraceRootDepth()
	relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
	require.NoError(t, err)
	traceDepth := game.MaxDepth() - game.TraceRootDepth()
	postTraceIdx = relativePos.TraceIndex(traceDepth)
	preTraceIdx = new(big.Int).Sub(postTraceIdx, big.NewInt(1))

	absolutePre, err := provider.AbsolutePreStateCommitment(ctx)
	require.NoError(t, err)
	expectPreDA = types.DAItem{
		DataHash: absolutePre,
	}

	expectPostDA = types.DAItem{
		DataHash: subClaimsDep8[0],
		Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
	}
	return
}

func attackFirstBranch(t *testing.T, ctx context.Context, gameBuilder *test.GameBuilder, provider *test.AlphabetWithProofProvider) (game types.Game, ref types.Claim, preTraceIdx *big.Int, postTraceIdx *big.Int, expectPreDA types.DAItem, expectPostDA types.DAItem) {
	depth8Branch := uint64(1)
	stepBranch := int64(0)
	// output roots
	seq := gameBuilder.Seq()
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], depth8Branch)

	game = gameBuilder.Game
	ref = game.Claims()[4]
	pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
	traceRootDepth := game.TraceRootDepth()
	relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
	require.NoError(t, err)

	traceDepth := game.MaxDepth() - game.TraceRootDepth()
	postTraceIdx = relativePos.TraceIndex(traceDepth)
	preTraceIdx = new(big.Int).Sub(postTraceIdx, big.NewInt(1))

	expectPreDA = types.DAItem{
		DataHash: subClaimsDep6[0],
		Proof:    append(subClaimsDep6[1][:], subClaimsDep6[2][:]...),
	}

	expectPostDA = types.DAItem{
		DataHash: subClaimsDep8[0],
		Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
	}
	return
}
func attackMidBranch(t *testing.T, ctx context.Context, gameBuilder *test.GameBuilder, provider *test.AlphabetWithProofProvider) (game types.Game, ref types.Claim, preTraceIdx *big.Int, postTraceIdx *big.Int, expectPreDA types.DAItem, expectPostDA types.DAItem) {
	depth8Branch := uint64(1)
	stepBranch := int64(1)
	// output roots
	seq := gameBuilder.Seq()
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], depth8Branch)

	game = gameBuilder.Game
	ref = game.Claims()[4]
	pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
	traceRootDepth := game.TraceRootDepth()
	relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
	require.NoError(t, err)
	traceDepth := game.MaxDepth() - game.TraceRootDepth()
	postTraceIdx = relativePos.TraceIndex(traceDepth)
	preTraceIdx = new(big.Int).Sub(postTraceIdx, big.NewInt(1))

	expectPreDA = types.DAItem{
		DataHash: subClaimsDep8[0],
		Proof:    append(subClaimsDep8[1][:], subClaimsDep8[2][:]...),
	}

	expectPostDA = types.DAItem{
		DataHash: subClaimsDep8[1],
		Proof:    append(subClaimsDep8[0][:], subClaimsDep8[2][:]...),
	}
	return
}

func attackMaxBranch(t *testing.T, ctx context.Context, gameBuilder *test.GameBuilder, provider *test.AlphabetWithProofProvider) (game types.Game, ref types.Claim, preTraceIdx *big.Int, postTraceIdx *big.Int, expectPreDA types.DAItem, expectPostDA types.DAItem) {
	depth8Branch := uint64(1)
	stepBranch := int64(3)
	// output roots
	seq := gameBuilder.Seq()
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], depth8Branch)

	game = gameBuilder.Game
	ref = game.Claims()[4]
	pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
	traceRootDepth := game.TraceRootDepth()
	relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
	require.NoError(t, err)
	traceDepth := game.MaxDepth() - game.TraceRootDepth()
	postTraceIdx = relativePos.TraceIndex(traceDepth)
	preTraceIdx = new(big.Int).Sub(postTraceIdx, big.NewInt(1))

	expectPreDA = types.DAItem{
		DataHash: subClaimsDep8[2],
		Proof:    crypto.Keccak256(subClaimsDep8[0][:], subClaimsDep8[1][:]),
	}

	expectPostDA = types.DAItem{
		DataHash: subClaimsDep6[1],
		Proof:    append(subClaimsDep6[0][:], subClaimsDep6[2][:]...),
	}
	return
}
func attackRightMost(t *testing.T, ctx context.Context, gameBuilder *test.GameBuilder, provider *test.AlphabetWithProofProvider) (game types.Game, ref types.Claim, preTraceIdx *big.Int, postTraceIdx *big.Int, expectPreDA types.DAItem, expectPostDA types.DAItem) {
	depth8Branch := uint64(3)
	stepBranch := int64(3)
	// output roots
	seq := gameBuilder.Seq()
	subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}} // claim at splitDepth
	seq = seq.Attack2(subClaimsDep2[:], 0)
	// traces
	subClaimsDep4 := []common.Hash{{0x04}}
	seq = seq.Attack2(subClaimsDep4[:], 0)
	subClaimsDep6 := []common.Hash{{0x61}, {0x62}, {0x63}}
	seq = seq.Attack2(subClaimsDep6[:], 0)
	subClaimsDep8 := []common.Hash{{0x81}, {0x82}, {0x83}}
	seq.Attack2(subClaimsDep8[:], depth8Branch)

	game = gameBuilder.Game
	ref = game.Claims()[4]
	pos := types.NewPositionFromGIndex(big.NewInt(ref.ToGIndex().Int64() + stepBranch))
	traceRootDepth := game.TraceRootDepth()
	relativePos, err := pos.RelativeToAncestorAtDepth(traceRootDepth)
	require.NoError(t, err)
	traceDepth := game.MaxDepth() - game.TraceRootDepth()
	postTraceIdx = relativePos.TraceIndex(traceDepth)
	preTraceIdx = new(big.Int).Sub(postTraceIdx, big.NewInt(1))

	expectPreDA = types.DAItem{
		DataHash: subClaimsDep8[2],
		Proof:    crypto.Keccak256(subClaimsDep8[0][:], subClaimsDep8[1][:]),
	}

	splitLeaf := game.Claims()[2]
	expectPostDA = types.DAItem{
		DataHash: splitLeaf.Value,
	}
	return
}
