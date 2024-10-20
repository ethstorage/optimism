package split

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

var (
	errRefClaimNotDeepEnough = errors.New("reference claim is not deep enough")
)

type ProviderCreator func(ctx context.Context, depth types.Depth, pre types.Claim, post types.Claim) (types.TraceProvider, error)

func NewSplitProviderSelector(topProvider types.TraceProvider, topDepth types.Depth, bottomProviderCreator ProviderCreator) trace.ProviderSelector {
	return func(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error) {
		if pos.Depth() <= topDepth {
			return topProvider, nil
		}
		if ref.Position.Depth() < topDepth {
			return nil, fmt.Errorf("%w, claim depth: %v, depth required: %v", errRefClaimNotDeepEnough, ref.Position.Depth(), topDepth)
		}

		// Find the ancestor claim at the leaf level for the top game.
		topLeaf, err := findAncestorAtDepth(game, ref, topDepth)
		if err != nil {
			return nil, err
		}

		var pre, post types.Claim
		// If pos is to the right of the leaf from the top game, we must be defending that output root
		// otherwise, we're attacking it.
		if pos.TraceIndex(pos.Depth()).Cmp(topLeaf.TraceIndex(pos.Depth())) > 0 {
			// Defending the top leaf claim, so use it as the pre-claim and find the post
			pre = topLeaf
			postTraceIdx := new(big.Int).Add(pre.TraceIndex(topDepth), big.NewInt(1))
			post, err = findAncestorWithTraceIndex(game, topLeaf, topDepth, postTraceIdx)
			if err != nil {
				return nil, fmt.Errorf("failed to find post claim: %w", err)
			}
		} else {
			// Attacking the top leaf claim, so use it as the post-claim and find the pre
			post = topLeaf
			postTraceIdx := post.TraceIndex(topDepth)
			if postTraceIdx.Cmp(big.NewInt(0)) == 0 {
				pre = types.Claim{}
			} else {
				preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
				pre, err = findAncestorWithTraceIndex(game, topLeaf, topDepth, preTraceIdx)
				if err != nil {
					return nil, fmt.Errorf("failed to find pre claim: %w", err)
				}
			}
		}
		// The top game runs from depth 0 to split depth *inclusive*.
		// The - 1 here accounts for the fact that the split depth is included in the top game.
		bottomDepth := game.MaxDepth() - topDepth - 1
		provider, err := bottomProviderCreator(ctx, bottomDepth, pre, post)
		if err != nil {
			return nil, err
		}
		// Translate such that the root of the bottom game is the level below the top game leaf
		return trace.Translate(provider, topDepth+1), nil
	}
}

func findAncestorAtDepth(game types.Game, claim types.Claim, depth types.Depth) (types.Claim, error) {
	for claim.Depth() > depth {
		parent, err := game.GetParent(claim)
		if err != nil {
			return types.Claim{}, fmt.Errorf("failed to find ancestor at depth %v: %w", depth, err)
		}
		claim = parent
	}
	return claim, nil
}

func findAncestorWithTraceIndex(game types.Game, ref types.Claim, depth types.Depth, traceIdx *big.Int) (types.Claim, error) {
	candidate := ref
	for candidate.TraceIndex(depth).Cmp(traceIdx) != 0 {
		parent, err := game.GetParent(candidate)
		if err != nil {
			return types.Claim{}, fmt.Errorf("failed to get parent of claim %v: %w", candidate.ContractIndex, err)
		}
		candidate = parent
	}
	return candidate, nil
}

// Get the traceIdx's accestor claims hash and its merkel proof in subclaims. traceIdx can be ref's traceIdx ±1.
func findAncestorProofAtDepth2(game types.Game, ref types.Claim, traceIdx *big.Int) (types.DAItem, error) {
	// If traceIdx equals -1, the absolute prestate is returned.
	if traceIdx.Cmp(big.NewInt(-1)) == 0 {
		return types.DAItem{
			DaType:   types.CallDataType,
			DataHash: game.AbsolutePrestate().Value.Bytes(),
			Proof:    []byte{},
		}, fmt.Errorf("absolute prestate is not implemented")
	}

	// If traceIdx is the right most branch, the root Claim is returned.
	if new(big.Int).Add(traceIdx, big.NewInt(1)).Cmp(new(big.Int).Lsh(big.NewInt(1), uint(game.MaxDepth()))) == 0 {
		return types.DAItem{
			DaType:   types.CallDataType,
			DataHash: game.RootClaim().Value.Bytes(),
			Proof:    []byte{},
		}, nil
	}

	ancestor := ref
	minTraceIndex := ancestor.TraceIndex(game.MaxDepth())
	offset := new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(game.MaxDepth()-ancestor.Depth()))
	maxTraceIndex := new(big.Int).Add(minTraceIndex, offset)
	// If Ancestor's left most branch traceIdx <= traceId <= Ancestor's rightmost branch traceIdx, we do nothing;
	// otherwise, we should trace the correct ancestor
	for minTraceIndex.Cmp(traceIdx) == 1 || maxTraceIndex.Cmp(traceIdx) == -1 {
		parent, err := game.GetParent(ancestor)
		if err != nil {
			return types.DAItem{}, fmt.Errorf("failed to get ancestor of claim %v: %w", ancestor.ContractIndex, err)
		}

		ancestor = parent
		minTraceIndex = ancestor.TraceIndex(game.MaxDepth())
		offset = new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(game.MaxDepth()-ancestor.Depth()))
		maxTraceIndex = new(big.Int).Add(minTraceIndex, offset)
	}

	branch := new(big.Int).Div(new(big.Int).Sub(traceIdx, minTraceIndex), new(big.Int).Lsh(big.NewInt(1), uint(game.MaxDepth()-ancestor.Depth()))).Int64()
	subclaims, err := game.GetSubClaims(ancestor)
	if err != nil {
		return types.DAItem{}, fmt.Errorf("failed to get subclaims of claim %v: %w", ref.ContractIndex, err)
	}

	ancestorClaim := subclaims[branch]
	merkleProof := make([]byte, 0)
	for i := int64(0); i < branch; i++ {
		merkleProof = append(merkleProof, subclaims[i].Bytes()...)
	}

	for i := branch + 1; i < int64(game.MaxAttackBranch()); i++ {
		merkleProof = append(merkleProof, subclaims[i].Bytes()...)
	}

	return types.DAItem{
		DaType:   types.CallDataType,
		DataHash: ancestorClaim.Bytes(),
		Proof:    merkleProof,
	}, nil
}
