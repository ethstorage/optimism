package trace

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

// Get the traceIdx's accestor claims hash and its merkel proof in subclaims. traceIdx can be ref's traceIdx ±1.
// Params:
// - relativeTraceIdx: the trace index relative to splitDepth
func findAncestorProofAtDepth2(ctx context.Context, provider types.TraceProvider, game types.Game, ref types.Claim, relativeTraceIdx *big.Int) (types.DAItem, error) {
	maxTraceDepth := game.MaxDepth() - game.SplitDepth()
	// Find the ancestor claim at the (splitDepth + 1) level.
	splitLeaf, err := FindAncestorAtDepth(game, ref, types.Depth(game.SplitDepth()+game.NBits()))
	if err != nil {
		return types.DAItem{}, err
	}
	// If traceIdx equals -1, the absolute prestate is returned.
	if relativeTraceIdx.Cmp(big.NewInt(-1)) == 0 {
		absolutePresate, err := provider.AbsolutePreStateCommitment(ctx)
		if err != nil {
			return types.DAItem{}, fmt.Errorf("failed to get absolutePrestate: %w", err)
		}
		return types.DAItem{
			DaType:   types.CallDataType,
			DataHash: absolutePresate.Bytes(),
			Proof:    []byte{},
		}, nil
	}

	// If traceIdx is the right most branch, the root Claim is returned.
	if new(big.Int).Add(relativeTraceIdx, big.NewInt(1)).Cmp(new(big.Int).Lsh(big.NewInt(1), uint(maxTraceDepth-game.NBits()))) == 0 {
		return types.DAItem{
			DaType:   types.CallDataType,
			DataHash: splitLeaf.Value.Bytes(),
			Proof:    []byte{},
		}, nil
	}

	ancestor := ref
	relativeAncestorPos, err := ancestor.RelativeToAncestorAtDepth(game.SplitDepth())
	if err != nil {
		return types.DAItem{}, err
	}
	minTraceIndex := relativeAncestorPos.TraceIndex(maxTraceDepth)
	offset := new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(maxTraceDepth-relativeAncestorPos.Depth()))
	maxTraceIndex := new(big.Int).Add(minTraceIndex, offset)
	// If Ancestor's left most branch traceIdx <= traceId <= Ancestor's rightmost branch traceIdx, we do nothing;
	// otherwise, we should trace the correct ancestor
	for minTraceIndex.Cmp(relativeTraceIdx) == 1 || maxTraceIndex.Cmp(relativeTraceIdx) == -1 {
		parent, err := game.GetParent(ancestor)
		if err != nil {
			return types.DAItem{}, fmt.Errorf("failed to get ancestor of claim %v: %w", ancestor.ContractIndex, err)
		}

		ancestor = parent
		relativeAncestorPos, err = ancestor.RelativeToAncestorAtDepth(game.SplitDepth())
		if err != nil {
			return types.DAItem{}, err
		}
		minTraceIndex = relativeAncestorPos.TraceIndex(maxTraceDepth)
		offset = new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(maxTraceDepth-relativeAncestorPos.Depth()))
		maxTraceIndex = new(big.Int).Add(minTraceIndex, offset)
	}

	branch := new(big.Int).Div(new(big.Int).Sub(relativeTraceIdx, minTraceIndex), new(big.Int).Lsh(big.NewInt(1), uint(maxTraceDepth-relativeAncestorPos.Depth()))).Int64()
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

// Params:
// - ref: the attacked claim
// - pos: next attacking position (global) for the ref claim.
func (t *Accessor) GetStepData2(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (prestate []byte, proofData []byte, preimageData *types.PreimageOracleData, err error) {
	// Get oracle data
	// prestate, proofData, preimageData, err = t.GetStepData(ctx, game, ref, pos)
	provider, err := t.selector(ctx, game, ref, pos)
	if err != nil {
		return nil, nil, nil, err
	}
	prestate, proofData, preimageData, err = provider.GetStepData2(ctx, pos, pos)
	if err != nil {
		return nil, nil, nil, err
	}
	// Get stepProof for stepV2
	relativePos, err := pos.RelativeToAncestorAtDepth(game.SplitDepth())
	if err != nil {
		return nil, nil, nil, err
	}
	postTraceIdx := relativePos.TraceIndex(game.MaxDepth())
	preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
	preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get postStateDaItem at trace index %v: %w", preTraceIdx, err)
	}
	postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get postStateDaItem at trace index %v: %w", postTraceIdx, err)
	}
	StepProof := types.StepProofData{
		Prestate:  preStateDaItem,
		PostState: postStateDaItem,
	}

	preimageData.StepProof = StepProof
	return prestate, proofData, preimageData, nil
}

func FindAncestorAtDepth(game types.Game, claim types.Claim, depth types.Depth) (types.Claim, error) {
	for claim.Depth() > depth {
		parent, err := game.GetParent(claim)
		if err != nil {
			return types.Claim{}, fmt.Errorf("failed to find ancestor at depth %v: %w", depth, err)
		}
		claim = parent
	}
	return claim, nil
}
