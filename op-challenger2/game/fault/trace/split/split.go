package split

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace/utils"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
)

var (
	errRefClaimNotDeepEnough = errors.New("reference claim is not deep enough")
)

type ProviderCreator func(ctx context.Context, depth types.Depth, pre types.Claim, post types.Claim) (types.TraceProvider, error)

func NewSplitProviderSelector(topProvider types.TraceProvider, topDepth types.Depth, bottomProviderCreator ProviderCreator) trace.ProviderSelector {
	// pos is next move position of any claim whose ancestor is ref
	return func(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error) {
		if pos.Depth() <= topDepth {
			return topProvider, nil
		}
		if ref.Position.Depth() < topDepth {
			return nil, fmt.Errorf("%w, claim depth: %v, depth required: %v", errRefClaimNotDeepEnough, ref.Position.Depth(), topDepth)
		}

		// Find the ancestor claim at the (splitDepth + nbits) level  for the top game.
		topLeaf, err := trace.FindAncestorAtDepth(game, ref, topDepth)
		if err != nil {
			return nil, err
		}

		ancestorPos := pos
		for ancestorPos.Depth() > topLeaf.Depth() {
			ancestorPos = ancestorPos.ParentN(game.NBits())
		}

		attackBranch := ancestorPos.IndexAtDepth().Int64() % (int64(game.MaxAttackBranch()) + 1)
		postTraceIdx := new(big.Int).Add(topLeaf.TraceIndex(topDepth), big.NewInt(attackBranch))
		preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))

		var outputRootDA = types.DAData{}
		pre, preDA, err := findAncestorWithTraceIndex2(game, topLeaf, topDepth, preTraceIdx)
		if err != nil {
			return nil, fmt.Errorf("failed to find pre claim: %w", err)
		}
		post, postDA, err := findAncestorWithTraceIndex2(game, topLeaf, topDepth, postTraceIdx)
		if err != nil {
			return nil, fmt.Errorf("failed to find post claim: %w", err)
		}

		outputRootDA.PreDA = preDA
		outputRootDA.PostDA = postDA

		// The top game runs from depth 0 to split depth *inclusive*.
		// The - nbits here accounts for the fact that the split depth is included in the top game.
		bottomDepth := game.MaxDepth() - topDepth - types.Depth(game.NBits())
		provider, err := bottomProviderCreator(ctx, bottomDepth, pre, post)
		if err != nil {
			return nil, err
		}

		// Translate such that the root of the bottom game is the level below the top game leaf
		return trace.Translate(provider, topDepth+types.Depth(game.NBits()), outputRootDA), nil
	}
}

// Get the traceIdx's accestor claim and pre/post outputroot proof for MSFDG. traceIdx must be ref's traceIdx ±1.
func findAncestorWithTraceIndex2(game types.Game, ref types.Claim, depth types.Depth, traceIdx *big.Int) (types.Claim, types.DAItem, error) {
	// If traceIdx equals -1, the absolute prestate is returned.
	if traceIdx.Cmp(big.NewInt(-1)) == 0 {
		return types.Claim{}, types.DAItem{}, nil
	}

	// If traceIdx is the right most branch, the root Claim is returned.
	if new(big.Int).Add(traceIdx, big.NewInt(1)).Cmp(new(big.Int).Lsh(big.NewInt(1), uint(depth))) == 0 {
		return game.RootClaim(), types.DAItem{}, nil
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
			return types.Claim{}, types.DAItem{}, fmt.Errorf("failed to get ancestor of claim %v: %w", ancestor.ContractIndex, err)
		}

		ancestor = parent
		minTraceIndex = ancestor.TraceIndex(depth)
		offset = new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(depth-ancestor.Depth()))
		maxTraceIndex = new(big.Int).Add(minTraceIndex, offset)
	}

	branch := new(big.Int).Div(new(big.Int).Sub(traceIdx, minTraceIndex), new(big.Int).Lsh(big.NewInt(1), uint(depth-ancestor.Depth()))).Int64()
	if ancestor.SubValues == nil {
		return types.Claim{}, types.DAItem{}, fmt.Errorf("ancestor's subvalues are nil %v", ancestor.ContractIndex)
	}
	subValues := *ancestor.SubValues
	claim := ancestor
	claim.Value = subValues[branch]
	merkleProof := utils.GenerateProofForSubValues(subValues, uint32(branch))

	outputRootDAItem := types.DAItem{
		DaType:   types.CallDataType,
		DataHash: claim.Value,
		Proof:    merkleProof,
	}

	return claim, outputRootDAItem, nil
}
