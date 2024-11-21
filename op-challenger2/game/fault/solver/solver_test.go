package solver

import (
	"context"
	"math/big"
	"testing"

	faulttest "github.com/ethereum-optimism/optimism/op-challenger2/game/fault/test"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace/alphabet"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace/split"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestAttemptStepNary2(t *testing.T) {
	maxDepth := types.Depth(7)
	startingL2BlockNumber := big.NewInt(0)
	nbits := uint64(1)
	splitDepth := types.Depth(3)
	claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, startingL2BlockNumber, maxDepth, nbits, splitDepth)
	traceDepth := maxDepth - splitDepth - types.Depth(nbits)

	// Last accessible leaf is the second last trace index
	// The root node is used for the last trace index and can only be attacked.
	lastLeafTraceIndex := big.NewInt(1<<traceDepth - 2)
	lastLeafTraceIndexPlusOne := big.NewInt(1<<traceDepth - 1)
	ctx := context.Background()
	absolutePrestate, err := claimBuilder.CorrectTraceProvider().AbsolutePreStateCommitment(ctx)
	require.NoError(t, err)

	tests := []struct {
		name                string
		agreeWithOutputRoot bool
		expectedErr         error
		expectNoStep        bool
		expectAttackBranch  uint64
		expectPreState      []byte
		expectProofData     []byte
		expectedOracleData  *types.PreimageOracleData
		expectedVMStateData *types.DAData
		setupGame           func(builder *faulttest.GameBuilder)
	}{
		{
			name:               "AttackFirstTraceIndex",
			expectAttackBranch: 0,
			expectPreState:     claimBuilder.CorrectPreState(common.Big0),
			expectProofData:    claimBuilder.CorrectProofData(common.Big0),
			expectedOracleData: claimBuilder.CorrectOracleData(common.Big0),
			expectedVMStateData: &types.DAData{
				PreDA:  types.DAItem{DataHash: absolutePrestate},
				PostDA: types.DAItem{DataHash: common.Hash{0xbb}},
			},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2([]common.Hash{{0xaa}}, 0).
					Attack2(nil, 0).
					Attack2([]common.Hash{{0xbb}}, 0)
			},
		},
		{
			name:               "DefendFirstTraceIndex",
			expectAttackBranch: 1,
			expectPreState:     claimBuilder.CorrectPreState(big.NewInt(1)),
			expectProofData:    claimBuilder.CorrectProofData(big.NewInt(1)),
			expectedOracleData: claimBuilder.CorrectOracleData(big.NewInt(1)),
			expectedVMStateData: &types.DAData{
				PreDA:  types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(7), common.Big0))},
				PostDA: types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(7), common.Big1))},
			},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2([]common.Hash{{0xaa}}, 0).
					Attack2(nil, 0).
					Attack2(nil, 0)
			},
		},
		{
			name:               "AttackMiddleTraceIndex",
			expectAttackBranch: 0,
			expectPreState:     claimBuilder.CorrectPreState(big.NewInt(4)),
			expectProofData:    claimBuilder.CorrectProofData(big.NewInt(4)),
			expectedOracleData: claimBuilder.CorrectOracleData(big.NewInt(4)),
			expectedVMStateData: &types.DAData{
				PreDA:  types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(5), big.NewInt(0)))},
				PostDA: types.DAItem{DataHash: common.Hash{0xaa}},
			},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 1).
					Attack2([]common.Hash{{0xaa}}, 0)
			},
		},
		{
			name:               "DefendMiddleTraceIndex",
			expectAttackBranch: 1,
			expectPreState:     claimBuilder.CorrectPreState(big.NewInt(5)),
			expectProofData:    claimBuilder.CorrectProofData(big.NewInt(5)),
			expectedOracleData: claimBuilder.CorrectOracleData(big.NewInt(5)),
			expectedVMStateData: &types.DAData{
				PreDA:  types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(7), big.NewInt(4)))},
				PostDA: types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(6), big.NewInt(2)))},
			},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 1).
					Attack2(nil, 0)
			},
		},
		{
			name:               "AttackLastTraceIndex",
			expectAttackBranch: 0,
			expectPreState:     claimBuilder.CorrectPreState(lastLeafTraceIndex),
			expectProofData:    claimBuilder.CorrectProofData(lastLeafTraceIndex),
			expectedOracleData: claimBuilder.CorrectOracleData(lastLeafTraceIndex),
			expectedVMStateData: &types.DAData{
				PreDA:  types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(6), common.Big2))},
				PostDA: types.DAItem{DataHash: common.Hash{0xaa}},
			},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 1).
					Attack2([]common.Hash{{0xaa}}, 1)
			},
		},
		{
			name:               "DefendLastTraceIndex",
			expectAttackBranch: 1,
			expectPreState:     claimBuilder.CorrectPreState(lastLeafTraceIndexPlusOne),
			expectProofData:    claimBuilder.CorrectProofData(lastLeafTraceIndexPlusOne),
			expectedOracleData: claimBuilder.CorrectOracleData(lastLeafTraceIndexPlusOne),
			expectedVMStateData: &types.DAData{
				PreDA:  types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(7), big.NewInt(6)))},
				PostDA: types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(4), big.NewInt(0)))},
			},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 1).
					Attack2(nil, 1)
			},
		},
		{
			name: "CannotStepNonLeaf",
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().Attack2(nil, 0).Attack2(nil, 0)
			},
			expectedErr:         ErrStepNonLeafNode,
			agreeWithOutputRoot: true,
		},
		{
			name: "CannotStepAgreedNode",
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2([]common.Hash{{0xaa}}, 0).
					Attack2(nil, 0)
			},
			expectNoStep:        true,
			agreeWithOutputRoot: true,
		},
		{
			name: "CannotStepInvalidPath",
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2([]common.Hash{{0xaa}}, 0).
					Attack2([]common.Hash{{0xbb}}, 0).
					Attack2([]common.Hash{{0xcc}}, 0)
			},
			expectNoStep:        true,
			agreeWithOutputRoot: true,
		},
		{
			name:               "CannotStepNearlyValidPath",
			expectAttackBranch: 1,
			expectPreState:     claimBuilder.CorrectPreState(big.NewInt(4)),
			expectProofData:    claimBuilder.CorrectProofData(big.NewInt(4)),
			expectedOracleData: claimBuilder.CorrectOracleData(big.NewInt(4)),
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 1).
					Attack2(nil, 1)
			},
			expectNoStep:        true,
			agreeWithOutputRoot: true,
		},
	}

	for _, tableTest := range tests {
		tableTest := tableTest
		t.Run(tableTest.name, func(t *testing.T) {
			builder := claimBuilder.GameBuilder(faulttest.WithInvalidValue(tableTest.agreeWithOutputRoot))
			tableTest.setupGame(builder)
			alphabetSolver := newClaimSolver(maxDepth, trace.NewSimpleTraceAccessor(claimBuilder.CorrectTraceProvider()), types.CallDataType)
			game := builder.Game
			claims := game.Claims()
			lastClaim := claims[len(claims)-1]
			agreedClaims := newHonestClaimTracker()
			if tableTest.agreeWithOutputRoot {
				agreedClaims.AddHonestClaim(types.Claim{}, claims[0])
			}
			if (lastClaim.Depth()%2 == 0) == tableTest.agreeWithOutputRoot {
				parentClaim := claims[lastClaim.ParentContractIndex]
				grandParentClaim := claims[parentClaim.ParentContractIndex]
				agreedClaims.AddHonestClaim(grandParentClaim, parentClaim)
			}
			step, err := alphabetSolver.AttemptStep(ctx, game, lastClaim, agreedClaims, 0)
			require.ErrorIs(t, err, tableTest.expectedErr)
			if !tableTest.expectNoStep && tableTest.expectedErr == nil {
				require.NotNil(t, step)
				require.Equal(t, lastClaim, step.LeafClaim)
				require.Equal(t, tableTest.expectAttackBranch, step.AttackBranch)
				require.Equal(t, tableTest.expectPreState, step.PreState)
				require.Equal(t, tableTest.expectProofData, step.ProofData)
				require.Equal(t, tableTest.expectedOracleData.IsLocal, step.OracleData.IsLocal)
				require.Equal(t, tableTest.expectedOracleData.OracleKey, step.OracleData.OracleKey)
				require.Equal(t, tableTest.expectedOracleData.GetPreimageWithSize(), step.OracleData.GetPreimageWithSize())
				require.Equal(t, tableTest.expectedOracleData.OracleOffset, step.OracleData.OracleOffset)
				require.Equal(t, tableTest.expectedVMStateData.PreDA, step.OracleData.VMStateDA.PreDA)
				require.Equal(t, tableTest.expectedVMStateData.PostDA, step.OracleData.VMStateDA.PostDA)
			} else {
				require.Nil(t, step)
			}
		})
	}
}

type bottomTraceProvider struct {
	pre  types.Claim
	post types.Claim
	*faulttest.AlphabetWithProofProvider
}

func TestAttemptStepNary2WithLocalData(t *testing.T) {
	maxDepth := types.Depth(7)
	startingL2BlockNumber := big.NewInt(0)
	nbits := uint64(1)
	splitDepth := types.Depth(3)
	claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, startingL2BlockNumber, maxDepth, nbits, splitDepth)
	traceDepth := maxDepth - splitDepth - types.Depth(nbits)
	callDataType := types.CallDataType
	oracleKey := faulttest.OraclPostKey
	ctx := context.Background()

	top := alphabet.NewTraceProvider(big.NewInt(0), splitDepth)
	bottomCreator := func(ctx context.Context, depth types.Depth, pre types.Claim, post types.Claim) (types.TraceProvider, error) {
		return &bottomTraceProvider{
			pre:  pre,
			post: post,
			// The oracel key is mocked, and it would be different from correct trace. So, in the following tests only oracle related
			AlphabetWithProofProvider: faulttest.NewAlphabetWithProofProvider(t, big.NewInt(0), traceDepth, nil, splitDepth, oracleKey),
		}, nil
	}
	selector := split.NewSplitProviderSelector(top, splitDepth, bottomCreator)
	accessor := trace.NewAccessor(selector)

	tests := []struct {
		name                string
		agreeWithOutputRoot bool
		expectedErr         error
		expectedOracleKey   []byte
		expectedVMStateData *types.DAData
		expectedLocalData   *types.DAItem
		setupGame           func(builder *faulttest.GameBuilder)
	}{
		{
			name:              "DefendFirstTraceIndex",
			expectedOracleKey: faulttest.LocalDataOracleKeyBytes(oracleKey),
			expectedVMStateData: &types.DAData{
				PreDA:  types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(7), common.Big0))},
				PostDA: types.DAItem{DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(7), common.Big1))},
			},
			expectedLocalData: &types.DAItem{
				DaType:   callDataType,
				DataHash: claimBuilder.CorrectClaimAtPosition(types.NewPosition(types.Depth(3), common.Big0)),
			},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2([]common.Hash{{0xaa}}, 0).
					Attack2(nil, 0).
					Attack2(nil, 0)
			},
		},
	}

	for _, tableTest := range tests {
		tableTest := tableTest
		claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, big.NewInt(0), maxDepth, nbits, splitDepth)
		builder := claimBuilder.GameBuilder(faulttest.WithInvalidValue(tableTest.agreeWithOutputRoot))
		tableTest.setupGame(builder)
		alphabetSolver := newClaimSolver(maxDepth, accessor, callDataType)
		game := builder.Game
		claims := game.Claims()
		lastClaim := claims[len(claims)-1]
		agreedClaims := newHonestClaimTracker()
		if tableTest.agreeWithOutputRoot {
			agreedClaims.AddHonestClaim(types.Claim{}, claims[0])
		}
		if (lastClaim.Depth()%2 == 0) == tableTest.agreeWithOutputRoot {
			parentClaim := claims[lastClaim.ParentContractIndex]
			grandParentClaim := claims[parentClaim.ParentContractIndex]
			agreedClaims.AddHonestClaim(grandParentClaim, parentClaim)
		}
		step, err := alphabetSolver.AttemptStep(ctx, game, lastClaim, agreedClaims, 0)
		require.ErrorIs(t, err, tableTest.expectedErr)

		if tableTest.expectedErr == nil {
			require.NotNil(t, step)
			require.Equal(t, lastClaim, step.LeafClaim)
			require.Equal(t, tableTest.expectedVMStateData.PreDA, step.OracleData.VMStateDA.PreDA)
			require.Equal(t, tableTest.expectedVMStateData.PostDA, step.OracleData.VMStateDA.PostDA)
			require.Equal(t, tableTest.expectedOracleKey, step.OracleData.OracleKey)
			require.Equal(t, tableTest.expectedLocalData.DaType, step.OracleData.OutputRootDAItem.DaType)
			require.Equal(t, tableTest.expectedLocalData.DataHash, step.OracleData.OutputRootDAItem.DataHash)
			require.Equal(t, tableTest.expectedLocalData.Proof, step.OracleData.OutputRootDAItem.Proof)
		} else {
			require.Nil(t, step)
		}
	}
}

func TestAttemptStepNary4(t *testing.T) {
	maxDepth := types.Depth(8)
	startingL2BlockNumber := big.NewInt(0)
	nbits := uint64(2)
	splitDepth := types.Depth(2)
	claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, startingL2BlockNumber, maxDepth, nbits, splitDepth)

	ctx := context.Background()
	absolutePrestate, err := claimBuilder.CorrectTraceProvider().AbsolutePreStateCommitment(ctx)
	require.NoError(t, err)

	claimAt := claimBuilder.CorrectClaimAtPosition
	tests := []struct {
		name                string
		agreeWithOutputRoot bool
		attackBranch        uint64
		expectedErr         error
		expectNoStep        bool
		expectAttackBranch  uint64
		expectPreState      []byte
		expectProofData     []byte
		expectedOracleData  *types.PreimageOracleData
		expectedVMStateData *types.DAData
		// In alphabet game, oracleKey is always blocknumber. So, expectedLocalData is always nil.
		expectedLocalData *types.DAItem
		setupGame         func(builder *faulttest.GameBuilder)
	}{
		{
			name:                "AttackLeftMostBranch",
			expectAttackBranch:  0,
			attackBranch:        0,
			agreeWithOutputRoot: true,
			expectPreState:      claimBuilder.CorrectPreState(common.Big0),
			expectProofData:     claimBuilder.CorrectProofData(common.Big0),
			expectedOracleData:  claimBuilder.CorrectOracleData(common.Big0),
			expectedVMStateData: &types.DAData{
				PreDA: types.DAItem{DataHash: absolutePrestate},
				PostDA: types.DAItem{
					DataHash: common.Hash{0x81},
					Proof:    append(common.Hash{0x82}.Bytes(), common.Hash{0x83}.Bytes()...),
				},
			},
			expectedLocalData: &types.DAItem{},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2([]common.Hash{{0x81}, {0x82}, {0x83}}, 0)
			},
		},
		{
			name:                "AttackFirstBranch",
			expectAttackBranch:  0,
			attackBranch:        0,
			agreeWithOutputRoot: true,
			expectPreState:      claimBuilder.CorrectPreState(big.NewInt(4)),
			expectProofData:     claimBuilder.CorrectProofData(big.NewInt(4)),
			expectedOracleData:  claimBuilder.CorrectOracleData(big.NewInt(4)),
			expectedVMStateData: &types.DAData{
				PreDA: types.DAItem{
					DataHash: claimAt(types.NewPosition(types.Depth(6), common.Big0)),
					Proof: append(
						claimAt(types.NewPosition(types.Depth(6), common.Big1)).Bytes(),
						claimAt(types.NewPosition(types.Depth(6), common.Big2)).Bytes()...,
					),
				},
				PostDA: types.DAItem{
					DataHash: common.Hash{0x81},
					Proof:    append(common.Hash{0x82}.Bytes(), common.Hash{0x83}.Bytes()...),
				},
			},
			expectedLocalData: &types.DAItem{},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2([]common.Hash{{0x81}, {0x82}, {0x83}}, 1)
			},
		},
		{
			name:                "AttackMidBranch",
			expectAttackBranch:  2,
			attackBranch:        2,
			agreeWithOutputRoot: true,
			expectPreState:      claimBuilder.CorrectPreState(big.NewInt(6)),
			expectProofData:     claimBuilder.CorrectProofData(big.NewInt(6)),
			expectedOracleData:  claimBuilder.CorrectOracleData(big.NewInt(6)),
			expectedVMStateData: &types.DAData{
				PreDA: types.DAItem{
					DataHash: claimAt(types.NewPosition(types.Depth(8), big.NewInt(5))),
					Proof:    append(claimAt(types.NewPosition(types.Depth(8), big.NewInt(4))).Bytes(), common.Hash{0x83}.Bytes()...),
				},
				PostDA: types.DAItem{
					DataHash: common.Hash{0x83},
					Proof: crypto.Keccak256(
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(4))).Bytes(),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(5))).Bytes(),
					),
				},
			},
			expectedLocalData: &types.DAItem{},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2([]common.Hash{
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(4))),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(5))),
						{0x83},
					}, 1)
			},
		},
		{ // if subValues at [0, maxAttackBranck -1] are correct, step would attack the maxAttackBranck
			name:                "AttackMaxBranch",
			expectAttackBranch:  3,
			attackBranch:        2,
			agreeWithOutputRoot: true,
			expectPreState:      claimBuilder.CorrectPreState(big.NewInt(7)),
			expectProofData:     claimBuilder.CorrectProofData(big.NewInt(7)),
			expectedOracleData:  claimBuilder.CorrectOracleData(big.NewInt(7)),
			expectedVMStateData: &types.DAData{
				PreDA: types.DAItem{
					DataHash: claimAt(types.NewPosition(types.Depth(8), big.NewInt(6))),
					Proof: crypto.Keccak256(
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(4))).Bytes(),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(5))).Bytes(),
					),
				},
				PostDA: types.DAItem{
					DataHash: claimAt(types.NewPosition(types.Depth(6), common.Big1)),
					Proof: append(
						claimAt(types.NewPosition(types.Depth(6), common.Big0)).Bytes(),
						claimAt(types.NewPosition(types.Depth(6), common.Big2)).Bytes()...,
					),
				},
			},
			expectedLocalData: &types.DAItem{},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2([]common.Hash{
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(4))),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(5))),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(6))),
					}, 1)
			},
		},
		{
			name:                "AttackRightMost",
			expectAttackBranch:  3,
			attackBranch:        2,
			agreeWithOutputRoot: true,
			expectPreState:      claimBuilder.CorrectPreState(big.NewInt(15)),
			expectProofData:     claimBuilder.CorrectProofData(big.NewInt(15)),
			expectedOracleData:  claimBuilder.CorrectOracleData(big.NewInt(15)),
			expectedVMStateData: &types.DAData{
				PreDA: types.DAItem{
					DataHash: claimAt(types.NewPosition(types.Depth(8), big.NewInt(14))),
					Proof: crypto.Keccak256(
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(12))).Bytes(),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(13))).Bytes(),
					),
				},
				PostDA: types.DAItem{
					DataHash: claimAt(types.NewPosition(types.Depth(4), common.Big0)),
				},
			},
			expectedLocalData: &types.DAItem{},
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().
					Attack2(nil, 0). // splitDepth
					Attack2(nil, 0).
					Attack2(nil, 0).
					Attack2([]common.Hash{
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(12))),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(13))),
						claimAt(types.NewPosition(types.Depth(8), big.NewInt(14))),
					}, 3)
			},
		},
	}

	for _, tableTest := range tests {
		tableTest := tableTest
		t.Run(tableTest.name, func(t *testing.T) {
			builder := claimBuilder.GameBuilder(faulttest.WithInvalidValue(tableTest.agreeWithOutputRoot))
			tableTest.setupGame(builder)
			accessor := trace.NewSimpleTraceAccessor(claimBuilder.CorrectTraceProvider())
			alphabetSolver := newClaimSolver(maxDepth, accessor, types.CallDataType)
			game := builder.Game
			claims := game.Claims()
			lastClaim := claims[len(claims)-1]
			agreedClaims := newHonestClaimTracker()
			if tableTest.agreeWithOutputRoot {
				agreedClaims.AddHonestClaim(types.Claim{}, claims[0])
			}
			if (lastClaim.Depth()%(types.Depth(2*nbits)) == 0) == tableTest.agreeWithOutputRoot {
				parentClaim := claims[lastClaim.ParentContractIndex]
				grandParentClaim := claims[parentClaim.ParentContractIndex]
				agreedClaims.AddHonestClaim(grandParentClaim, parentClaim)
			}
			step, err := alphabetSolver.AttemptStep(ctx, game, lastClaim, agreedClaims, tableTest.attackBranch)
			require.ErrorIs(t, err, tableTest.expectedErr)
			_, _, preimage, err := accessor.GetStepData2(ctx, game, lastClaim, lastClaim.MoveRightN(tableTest.expectAttackBranch))
			require.NoError(t, err)
			if !tableTest.expectNoStep && tableTest.expectedErr == nil {
				require.NotNil(t, step)
				require.Equal(t, lastClaim, step.LeafClaim)
				require.Equal(t, tableTest.expectAttackBranch, step.AttackBranch)
				require.Equal(t, tableTest.expectPreState, step.PreState)
				require.Equal(t, tableTest.expectProofData, step.ProofData)
				require.Equal(t, tableTest.expectedOracleData.IsLocal, step.OracleData.IsLocal)
				require.Equal(t, tableTest.expectedOracleData.OracleKey, step.OracleData.OracleKey)
				require.Equal(t, tableTest.expectedOracleData.GetPreimageWithSize(), step.OracleData.GetPreimageWithSize())
				require.Equal(t, tableTest.expectedOracleData.OracleOffset, step.OracleData.OracleOffset)
				require.Equal(t, tableTest.expectedVMStateData.PreDA, step.OracleData.VMStateDA.PreDA)
				require.Equal(t, tableTest.expectedVMStateData.PostDA, step.OracleData.VMStateDA.PostDA)
				require.Equal(t, tableTest.expectedVMStateData.PreDA, preimage.VMStateDA.PreDA)
				require.Equal(t, tableTest.expectedVMStateData.PostDA, preimage.VMStateDA.PostDA)
				require.Equal(t, tableTest.expectedLocalData.DaType, step.OracleData.OutputRootDAItem.DaType)
				require.Equal(t, tableTest.expectedLocalData.DataHash, step.OracleData.OutputRootDAItem.DataHash)
				require.Equal(t, tableTest.expectedLocalData.Proof, step.OracleData.OutputRootDAItem.Proof)
			} else {
				require.Nil(t, step)
			}
		})
	}
}
