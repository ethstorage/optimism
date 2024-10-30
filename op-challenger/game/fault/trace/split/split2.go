package split

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

// Get the traceIdx's accestor claims for multi-sec fault proof. traceIdx can be ref's traceIdx ±1.
func findAncestorWithTraceIndex2(game types.Game, ref types.Claim, depth types.Depth, traceIdx *big.Int) (types.Claim, error) {
	// If traceIdx equals -1, the absolute prestate is returned.
	if traceIdx.Cmp(big.NewInt(-1)) == 0 {
		return types.Claim{}, nil
	}

	// If traceIdx is the right most branch, the root Claim is returned.
	if new(big.Int).Add(traceIdx, big.NewInt(1)).Cmp(new(big.Int).Lsh(big.NewInt(1), uint(depth))) == 0 {
		return game.RootClaim(), nil
	}

	ancestor := ref
	minTraceIndex := ancestor.TraceIndex(depth)
	offset := new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(depth-ancestor.Depth()))
	maxTraceIndex := new(big.Int).Add(minTraceIndex, offset)
	// If Ancestor's left most branch traceIdx <= traceId <= Ancestor's rightmost branch traceIdx, we do nothing;
	// otherwise, we should trace the correct ancestor
	for minTraceIndex.Cmp(traceIdx) == 1 || maxTraceIndex.Cmp(traceIdx) == -1 {
		parent, err := game.GetParent(ancestor)
		if err != nil {
			return types.Claim{}, fmt.Errorf("failed to get ancestor of claim %v: %w", ancestor.ContractIndex, err)
		}

		ancestor = parent
		minTraceIndex = ancestor.TraceIndex(depth)
		offset = new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(depth-ancestor.Depth()))
		maxTraceIndex = new(big.Int).Add(minTraceIndex, offset)
	}

	branch := new(big.Int).Div(new(big.Int).Sub(traceIdx, minTraceIndex), new(big.Int).Lsh(big.NewInt(1), uint(depth-ancestor.Depth()))).Int64()
	subValues := ancestor.SubValues()

	claim := types.Claim{
		ClaimData: types.ClaimData{
			Value:    subValues[branch],
			Position: types.NewPositionFromGIndex(traceIdx),
		},
	}

	return claim, nil
}

func NewSplitProviderSelector2(topProvider types.TraceProvider, topDepth types.Depth, bottomProviderCreator ProviderCreator) trace.ProviderSelector {
	return func(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error) {
		if pos.Depth() <= topDepth {
			return topProvider, nil
		}
		if ref.Position.Depth() < topDepth {
			return nil, fmt.Errorf("%w, claim depth: %v, depth required: %v", errRefClaimNotDeepEnough, ref.Position.Depth(), topDepth)
		}

		// Find the ancestor claim at the (splitDepth + 1) level.
		splitLeaf, err := trace.FindAncestorAtDepth(game, ref, types.Depth(topDepth+game.NBits()))
		if err != nil {
			return nil, err
		}

		// Find the ancestor claim at the leaf level splitDepth for the top game.
		topLeaf, err := game.GetParent(splitLeaf)
		if err != nil {
			return nil, err
		}

		attackBranch := int64(splitLeaf.AttackBranch)
		preTraceIdx := new(big.Int).Add(topLeaf.TraceIndex(topDepth), big.NewInt(attackBranch))
		postTraceIdx := new(big.Int).Add(preTraceIdx, big.NewInt(1))
		pre, err := findAncestorWithTraceIndex2(game, topLeaf, topDepth, preTraceIdx)
		if err != nil {
			return nil, fmt.Errorf("failed to find pre claim: %w", err)
		}
		post, err := findAncestorWithTraceIndex2(game, topLeaf, topDepth, postTraceIdx)
		if err != nil {
			return nil, fmt.Errorf("failed to find post claim: %w", err)
		}

		// The top game runs from depth 0 to split depth *inclusive*.
		// The - 1 here accounts for the fact that the split depth is included in the top game.
		bottomDepth := game.MaxDepth() - topDepth - game.NBits()
		provider, err := bottomProviderCreator(ctx, bottomDepth, pre, post)
		if err != nil {
			return nil, err
		}
		// Translate such that the root of the bottom game is the level below the top game leaf
		return trace.Translate(provider, topDepth+game.NBits()), nil
	}
}
