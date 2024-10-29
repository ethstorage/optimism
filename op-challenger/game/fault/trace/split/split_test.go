package split

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/test"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/alphabet"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	gameDepth  = 7
	splitDepth = 3
)

func TestUseTopProvider(t *testing.T) {
	ctx := context.Background()
	topProvider, selector, gameBuilder := setupAlphabetSplitSelector(t)

	ref := gameBuilder.Game.Claims()[0]

	pos := ref.Position
	for pos.Depth() <= splitDepth {
		provider, err := selector(ctx, gameBuilder.Game, ref, ref.Position)
		require.NoError(t, err)
		require.Same(t, topProvider, provider)
		_, err = topProvider.Get(ctx, pos)
		require.NoError(t, err, "should be able to use provider for position")
		pos = pos.Attack()
	}
}

func TestErrorWhenRefAboveTopGameLeafButPositionInBottom(t *testing.T) {
	ctx := context.Background()
	_, selector, gameBuilder := setupAlphabetSplitSelector(t)

	// Generate claims at depths up to but not including the leaf of the top providers
	createClaimsToDepth(gameBuilder, splitDepth-1)
	for _, ref := range gameBuilder.Game.Claims() {
		pos := types.NewPosition(splitDepth+1, big.NewInt(0))
		provider, err := selector(ctx, gameBuilder.Game, ref, pos)
		require.ErrorIsf(t, err, errRefClaimNotDeepEnough, "should not get provider with ref claim at depth: %v", ref.Depth())
		require.Nil(t, provider)
	}
}

func TestTranslatePositionsForBottomProvider(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim)
	}{
		// There are 4 leaf nodes that can be accessed in the top tree of depth 3: 8, 10, 12, 14
		// Then you can attack and defend any of those to challenge all blocks
		{"attackTopLeafGIndex8", attackTopLeafGIndex8},
		{"defendTopLeafGIndex8", defendTopLeafGIndex8},
		{"attackTopLeafGIndex10", attackTopLeafGIndex10},
		{"defendTopLeafGIndex10", defendTopLeafGIndex10},
		{"attackTopLeafGIndex12", attackTopLeafGIndex12},
		{"defendTopLeafGIndex12", defendTopLeafGIndex12},
		{"attackTopLeafGIndex14", attackTopLeafGIndex14},
		{"attackTopLeafGIndex14", defendTopLeafGIndex14},
	}
	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			_, selector, gameBuilder := setupAlphabetSplitSelector(t)
			ref, pos, _, _ := tCase.setup(t, gameBuilder)
			provider, err := selector(context.Background(), gameBuilder.Game, ref, pos)
			require.NoError(t, err)

			claimPos := pos
			localClaimPos := types.NewPositionFromGIndex(big.NewInt(1))
			requireSameValue(t, provider, claimPos, asBottomTraceProvider(t, provider).AlphabetTraceProvider, localClaimPos)
			requireSameValue(t, provider, claimPos.Attack(), asBottomTraceProvider(t, provider).AlphabetTraceProvider, localClaimPos.Attack())
			requireSameValue(t, provider, claimPos.Attack().Defend(), asBottomTraceProvider(t, provider).AlphabetTraceProvider, localClaimPos.Attack().Defend())
		})
	}
}

func requireSameValue(t *testing.T, a types.TraceProvider, aPos types.Position, b types.TraceProvider, bPos types.Position) {
	// Check Get returns the same results
	aValue, err := a.Get(context.Background(), aPos)
	require.NoError(t, err)
	bValue, err := b.Get(context.Background(), bPos)
	require.NoError(t, err)
	require.Equal(t, aValue, bValue)

	// Check GetStepData returns the same results
	aPrestate, aProofData, aPreimageData, err := a.GetStepData(context.Background(), aPos)
	require.NoError(t, err)
	bPrestate, bProofData, bPreimageData, err := b.GetStepData(context.Background(), bPos)
	require.NoError(t, err)
	require.Equal(t, aPrestate, bPrestate)
	require.Equal(t, aProofData, bProofData)
	require.Equal(t, aPreimageData, bPreimageData)
}

func TestBottomProviderAttackingTopLeaf(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim)
	}{
		// There are 4 leaf nodes that can be accessed in the top tree of depth 3: 8, 10, 12, 14
		// Then you can attack and defend any of those to challenge all blocks
		// We can then use these setups to test any other reference claim descending from what these setup since
		// that whole subtree should have the same pre and post claim from the top provider.
		{"attackTopLeafGIndex8", attackTopLeafGIndex8},
		{"defendTopLeafGIndex8", defendTopLeafGIndex8},
		{"attackTopLeafGIndex10", attackTopLeafGIndex10},
		{"defendTopLeafGIndex10", defendTopLeafGIndex10},
		{"attackTopLeafGIndex12", attackTopLeafGIndex12},
		{"defendTopLeafGIndex12", defendTopLeafGIndex12},
		{"attackTopLeafGIndex14", attackTopLeafGIndex14},
		{"attackTopLeafGIndex14", defendTopLeafGIndex14},
	}
	for _, tCase := range tests {
		tCase := tCase
		t.Run(tCase.name, func(t *testing.T) {
			_, selector, gameBuilder := setupAlphabetSplitSelector(t)

			ref, pos, expectedPre, expectedPost := tCase.setup(t, gameBuilder)

			runTest := func(ref types.Claim, pos types.Position) {
				t.Run(fmt.Sprintf("Ref-d%vi%v_Pos-d%vi%v", ref.Depth(), ref.IndexAtDepth(), pos.Depth(), pos.IndexAtDepth()), func(t *testing.T) {
					provider, err := selector(context.Background(), gameBuilder.Game, ref, pos)
					require.NoError(t, err)
					requireBottomProviderForClaims(t, provider, expectedPre, expectedPost)
				})
			}

			// Check we get the same pre and post for any reference claim lower in the game
			var testDescendantClaims func(ref types.Claim, pos types.Position)
			testDescendantClaims = func(ref types.Claim, pos types.Position) {
				// For each reference claim, check it works with the claim position, or attacking or defending the claim
				runTest(ref, pos)
				runTest(ref, pos.Attack())
				runTest(ref, pos.Defend())
				if pos.Depth() >= gameDepth {
					return
				}

				// If the ref is the leaf of the top claim, ensure we respect whether the test is setup
				// to attack or defend the top leaf claim.
				if ref.Depth() != splitDepth || !pos.RightOf(ref.Position) {
					gameBuilder.SeqFrom(ref).Attack()
					attackRef := latestClaim(gameBuilder)
					testDescendantClaims(attackRef, attackRef.Position)
				}
				if ref.Depth() != splitDepth || pos.RightOf(ref.Position) {
					gameBuilder.SeqFrom(ref).Defend()
					defendRef := latestClaim(gameBuilder)
					testDescendantClaims(defendRef, defendRef.Position)
				}
			}
			testDescendantClaims(ref, pos)
		})
	}
}

func TestFindAncestorProofAtDepth2(t *testing.T) {
	var nbits uint64 = 2
	nary := int64(1) << nbits
	maxGameDepth := types.Depth(8)
	splitDepth := types.Depth(4)
	claimBuilder := test.NewAlphabetClaimBuilder(t, big.NewInt(0), maxGameDepth)

	{ // Attack the left most claim
		expectedTopLeafBranch := uint64(0)

		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}}
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)

		subClaimsDep4 := []common.Hash{{0x04}, {0x05}, {0x06}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep4[:], nbits, expectedTopLeafBranch)
		depth4Claim := gameBuilder.Game.Claims()[2]

		subClaimsDep6 := []common.Hash{{0x07}} // claim at splitDepth + 1
		seq = seq.Attack2(subClaimsDep6[:], nbits, 0)

		subClaimsDep8 := []common.Hash{{0x11, 0x12, 0x13}}
		seq.Attack2(subClaimsDep8[:], nbits, 0)
		depth8Claim := gameBuilder.Game.Claims()[4]

		game := gameBuilder.Game
		splitLeaf, err := trace.FindAncestorAtDepth(game, depth8Claim, types.Depth(splitDepth+types.Depth(nbits)))
		require.NoError(t, err)
		// Find the ancestor claim at the leaf level splitDepth for the top game.
		topLeaf, err := game.GetParent(splitLeaf)
		require.NoError(t, err)
		require.Equal(t, expectedTopLeafBranch, topLeaf.AttackBranch)

		attackBranch := int64(splitLeaf.AttackBranch)
		postTraceIdx := new(big.Int).Add(depth4Claim.TraceIndex(splitDepth), big.NewInt(attackBranch))
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
		expectedPre := types.Claim{}
		pre, err := findAncestorWithTraceIndex2(game, topLeaf, splitDepth, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, pre)

		expectedPost := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep4[0],
				Position: types.NewPositionFromGIndex(postTraceIdx),
			},
		}
		post, err := findAncestorWithTraceIndex2(gameBuilder.Game, topLeaf, splitDepth, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, post)
	}

	{ // Attack the left claim
		expectedTopLeafBranch := uint64(1)

		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}}
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)

		subClaimsDep4 := []common.Hash{{0x04}, {0x05}, {0x06}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep4[:], nbits, expectedTopLeafBranch)
		depth4Claim := gameBuilder.Game.Claims()[2]

		subClaimsDep6 := []common.Hash{{0x07}} // claim at splitDepth + 1
		seq = seq.Attack2(subClaimsDep6[:], nbits, 0)

		subClaimsDep8 := []common.Hash{{0x11, 0x12, 0x13}}
		seq.Attack2(subClaimsDep8[:], nbits, 0)
		depth8Claim := gameBuilder.Game.Claims()[4]

		game := gameBuilder.Game
		splitLeaf, err := trace.FindAncestorAtDepth(game, depth8Claim, types.Depth(splitDepth+types.Depth(nbits)))
		require.NoError(t, err)
		// Find the ancestor claim at the leaf level splitDepth for the top game.
		topLeaf, err := game.GetParent(splitLeaf)
		require.NoError(t, err)
		require.Equal(t, expectedTopLeafBranch, topLeaf.AttackBranch)

		attackBranch := int64(splitLeaf.AttackBranch)
		postTraceIdx := new(big.Int).Add(depth4Claim.TraceIndex(splitDepth), big.NewInt(attackBranch))
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
		expectedPre := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep2[0],
				Position: types.NewPositionFromGIndex(preTraceIdx),
			},
		}
		pre, err := findAncestorWithTraceIndex2(game, topLeaf, splitDepth, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, pre)

		expectedPost := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep4[0],
				Position: types.NewPositionFromGIndex(postTraceIdx),
			},
		}
		post, err := findAncestorWithTraceIndex2(gameBuilder.Game, topLeaf, splitDepth, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, post)
	}

	{ // Attack the middle claim
		expectedTopLeafBranch := uint64(1)
		expectedSplitBranch := uint64(1)

		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}}
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)

		subClaimsDep4 := []common.Hash{{0x04}, {0x05}, {0x06}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep4[:], nbits, expectedTopLeafBranch)
		depth4Claim := gameBuilder.Game.Claims()[2]

		subClaimsDep6 := []common.Hash{{0x07}} // claim at splitDepth + 1
		seq = seq.Attack2(subClaimsDep6[:], nbits, expectedSplitBranch)

		subClaimsDep8 := []common.Hash{{0x11, 0x12, 0x13}}
		seq.Attack2(subClaimsDep8[:], nbits, 0)
		depth8Claim := gameBuilder.Game.Claims()[4]

		game := gameBuilder.Game
		splitLeaf, err := trace.FindAncestorAtDepth(game, depth8Claim, types.Depth(splitDepth+types.Depth(nbits)))
		require.NoError(t, err)
		// Find the ancestor claim at the leaf level splitDepth for the top game.
		topLeaf, err := game.GetParent(splitLeaf)
		require.NoError(t, err)
		require.Equal(t, expectedTopLeafBranch, topLeaf.AttackBranch)
		require.Equal(t, expectedSplitBranch, splitLeaf.AttackBranch)

		attackBranch := int64(splitLeaf.AttackBranch)
		postTraceIdx := new(big.Int).Add(depth4Claim.TraceIndex(splitDepth), big.NewInt(attackBranch))
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
		expectedPre := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep4[0],
				Position: types.NewPositionFromGIndex(preTraceIdx),
			},
		}
		pre, err := findAncestorWithTraceIndex2(game, topLeaf, splitDepth, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, pre)

		expectedPost := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep4[1],
				Position: types.NewPositionFromGIndex(postTraceIdx),
			},
		}
		post, err := findAncestorWithTraceIndex2(gameBuilder.Game, topLeaf, splitDepth, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, post)
	}

	{ // Attack the right claim
		expectedTopLeafBranch := uint64(1)
		expectedSplitBranch := uint64(3)

		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}}
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)

		subClaimsDep4 := []common.Hash{{0x04}, {0x05}, {0x06}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep4[:], nbits, expectedTopLeafBranch)
		depth4Claim := gameBuilder.Game.Claims()[2]

		subClaimsDep6 := []common.Hash{{0x07}} // claim at splitDepth + 1
		seq = seq.Attack2(subClaimsDep6[:], nbits, expectedSplitBranch)

		subClaimsDep8 := []common.Hash{{0x11, 0x12, 0x13}}
		seq.Attack2(subClaimsDep8[:], nbits, 0)
		depth8Claim := gameBuilder.Game.Claims()[4]

		game := gameBuilder.Game
		splitLeaf, err := trace.FindAncestorAtDepth(game, depth8Claim, types.Depth(splitDepth+types.Depth(nbits)))
		require.NoError(t, err)
		// Find the ancestor claim at the leaf level splitDepth for the top game.
		topLeaf, err := game.GetParent(splitLeaf)
		require.NoError(t, err)
		require.Equal(t, expectedTopLeafBranch, topLeaf.AttackBranch)
		require.Equal(t, expectedSplitBranch, splitLeaf.AttackBranch)

		attackBranch := int64(splitLeaf.AttackBranch)
		postTraceIdx := new(big.Int).Add(depth4Claim.TraceIndex(splitDepth), big.NewInt(attackBranch))
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
		expectedPre := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep4[expectedSplitBranch-1],
				Position: types.NewPositionFromGIndex(preTraceIdx),
			},
		}
		pre, err := findAncestorWithTraceIndex2(game, topLeaf, splitDepth, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, pre)

		expectedPost := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep2[expectedTopLeafBranch],
				Position: types.NewPositionFromGIndex(postTraceIdx),
			},
		}
		post, err := findAncestorWithTraceIndex2(gameBuilder.Game, topLeaf, splitDepth, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, post)
	}

	{ // Attack the right most claim
		expectedTopLeafBranch := uint64(3)
		expectedSplitBranch := uint64(3)

		gameBuilder := claimBuilder.GameBuilder2(nary, splitDepth)
		seq := gameBuilder.Seq()
		claimBuilder.CreateRootClaim(test.WithValue(common.Hash{0x00}))
		subClaimsDep2 := []common.Hash{{0x01}, {0x02}, {0x03}}
		seq = seq.Attack2(subClaimsDep2[:], nbits, 0)

		subClaimsDep4 := []common.Hash{{0x04}, {0x05}, {0x06}} // claim at splitDepth
		seq = seq.Attack2(subClaimsDep4[:], nbits, expectedTopLeafBranch)
		depth4Claim := gameBuilder.Game.Claims()[2]

		subClaimsDep6 := []common.Hash{{0x07}} // claim at splitDepth + 1
		seq = seq.Attack2(subClaimsDep6[:], nbits, expectedSplitBranch)

		subClaimsDep8 := []common.Hash{{0x11, 0x12, 0x13}}
		seq.Attack2(subClaimsDep8[:], nbits, 0)
		depth8Claim := gameBuilder.Game.Claims()[4]

		game := gameBuilder.Game
		splitLeaf, err := trace.FindAncestorAtDepth(game, depth8Claim, types.Depth(splitDepth+types.Depth(nbits)))
		require.NoError(t, err)
		// Find the ancestor claim at the leaf level splitDepth for the top game.
		topLeaf, err := game.GetParent(splitLeaf)
		require.NoError(t, err)
		require.Equal(t, expectedTopLeafBranch, topLeaf.AttackBranch)
		require.Equal(t, expectedSplitBranch, splitLeaf.AttackBranch)

		attackBranch := int64(splitLeaf.AttackBranch)
		postTraceIdx := new(big.Int).Add(depth4Claim.TraceIndex(splitDepth), big.NewInt(attackBranch))
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
		expectedPre := types.Claim{
			ClaimData: types.ClaimData{
				Value:    subClaimsDep4[expectedSplitBranch-1],
				Position: types.NewPositionFromGIndex(preTraceIdx),
			},
		}
		pre, err := findAncestorWithTraceIndex2(game, topLeaf, splitDepth, preTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPre, pre)

		expectedPost := game.RootClaim()
		post, err := findAncestorWithTraceIndex2(gameBuilder.Game, topLeaf, splitDepth, postTraceIdx)
		require.NoError(t, err)
		require.Equal(t, expectedPost, post)
	}
}

func attackTopLeafGIndex8(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	// Generate claims down to the top provider's leaf
	seq := gameBuilder.Seq() // gindex 1, trace 7
	seq = seq.Attack()       // gindex 2, trace 3
	seq = seq.Attack()       // gindex 4, trace 1
	seq.Attack()             // gindex 8, trace 0
	expectPost = latestClaim(gameBuilder)

	// No pre-claim as the first output root is being challenged.
	expectPre = types.Claim{}

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Attack()
	return
}

func defendTopLeafGIndex8(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	// Generate claims down to the top provider's leaf
	seq := gameBuilder.Seq() // gindex 1, trace 7
	seq = seq.Attack()       // gindex 2, trace 3
	seq = seq.Attack()       // gindex 4, trace 1
	expectPost = latestClaim(gameBuilder)
	seq.Attack() // gindex 8, trace 0
	expectPre = latestClaim(gameBuilder)

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Defend()
	return
}

func attackTopLeafGIndex10(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	seq := gameBuilder.Seq() // gindex 1, trace 7
	seq = seq.Attack()       // gindex 2, trace 3
	seq = seq.Attack()       // gindex 4, trace 1
	expectPre = latestClaim(gameBuilder)
	seq.Defend() // gindex 10, trace 2
	expectPost = latestClaim(gameBuilder)

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Attack()
	return
}

func defendTopLeafGIndex10(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	seq := gameBuilder.Seq() // gindex 1, trace 7
	seq = seq.Attack()       // gindex 2, trace 3
	expectPost = latestClaim(gameBuilder)
	seq = seq.Attack() // gindex 4, trace 1
	seq.Defend()       // gindex 10, trace 2
	expectPre = latestClaim(gameBuilder)

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Defend()
	return
}

func attackTopLeafGIndex12(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	seq := gameBuilder.Seq() // gindex 1, trace 7
	seq = seq.Attack()       // gindex 2, trace 3
	expectPre = latestClaim(gameBuilder)
	seq = seq.Defend() // gindex 6, trace 5
	seq.Attack()       // gindex 12, trace 4
	expectPost = latestClaim(gameBuilder)

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Attack()
	return
}

func defendTopLeafGIndex12(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	seq := gameBuilder.Seq() // gindex 1, trace 7
	seq = seq.Attack()       // gindex 2, trace 3
	seq = seq.Defend()       // gindex 6, trace 5
	expectPost = latestClaim(gameBuilder)
	seq.Attack() // gindex 12, trace 4
	expectPre = latestClaim(gameBuilder)

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Defend()
	return
}

func attackTopLeafGIndex14(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	seq := gameBuilder.Seq() // gindex 1, trace 7
	seq = seq.Attack()       // gindex 2, trace 3
	seq = seq.Defend()       // gindex 6, trace 5
	expectPre = latestClaim(gameBuilder)
	seq.Defend() // gindex 14, trace 6
	expectPost = latestClaim(gameBuilder)

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Attack()
	return
}

func defendTopLeafGIndex14(_ *testing.T, gameBuilder *test.GameBuilder) (ref types.Claim, pos types.Position, expectPre types.Claim, expectPost types.Claim) {
	seq := gameBuilder.Seq() // gindex 1, trace 7
	expectPost = latestClaim(gameBuilder)
	seq = seq.Attack() // gindex 2, trace 3
	seq = seq.Defend() // gindex 6, trace 5
	seq.Defend()       // gindex 14, trace 6
	expectPre = latestClaim(gameBuilder)

	ref = latestClaim(gameBuilder)
	pos = ref.Position.Defend()
	return
}

func latestClaim(gameBuilder *test.GameBuilder) types.Claim {
	return gameBuilder.Game.Claims()[len(gameBuilder.Game.Claims())-1]
}

func createClaimsToDepth(gameBuilder *test.GameBuilder, depth int) {
	seq := gameBuilder.Seq()
	for i := 0; i < depth; i++ {
		seq = seq.Attack()
	}
}

func requireBottomProviderForClaims(t *testing.T, actual types.TraceProvider, expectedPre types.Claim, expectedPost types.Claim) {
	if expectedPre != (types.Claim{}) {
		require.Equal(t,
			new(big.Int).Add(expectedPre.TraceIndex(splitDepth), big.NewInt(1)),
			expectedPost.TraceIndex(splitDepth),
			"should expect adjacent top level trace indices")
	}

	bottomProvider := asBottomTraceProvider(t, actual)
	require.Equal(t, expectedPre, bottomProvider.pre, "Incorrect pre claim")
	require.Equal(t, expectedPost, bottomProvider.post, "Incorrect post claim")
}

func asBottomTraceProvider(t *testing.T, actual types.TraceProvider) *bottomTraceProvider {
	translatingProvider, ok := actual.(*trace.TranslatingProvider)
	require.True(t, ok)
	bottomProvider, ok := translatingProvider.Original().(*bottomTraceProvider)
	require.True(t, ok)
	return bottomProvider
}

func setupAlphabetSplitSelector(t *testing.T) (*alphabet.AlphabetTraceProvider, trace.ProviderSelector, *test.GameBuilder) {
	top := alphabet.NewTraceProvider(big.NewInt(0), splitDepth)
	bottomCreator := func(ctx context.Context, depth types.Depth, pre types.Claim, post types.Claim) (types.TraceProvider, error) {
		return &bottomTraceProvider{
			pre:                   pre,
			post:                  post,
			AlphabetTraceProvider: alphabet.NewTraceProvider(big.NewInt(0), depth),
		}, nil
	}
	selector := NewSplitProviderSelector(top, splitDepth, bottomCreator)

	claimBuilder := test.NewAlphabetClaimBuilder(t, big.NewInt(0), gameDepth)
	gameBuilder := claimBuilder.GameBuilder()
	return top, selector, gameBuilder
}

type bottomTraceProvider struct {
	pre  types.Claim
	post types.Claim
	*alphabet.AlphabetTraceProvider
}
