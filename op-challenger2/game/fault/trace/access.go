package trace

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace/utils"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum/go-ethereum/common"
)

func NewSimpleTraceAccessor(trace types.TraceProvider) *Accessor {
	selector := func(_ context.Context, _ types.Game, _ types.Claim, _ types.Position) (types.TraceProvider, error) {
		return trace, nil
	}
	return NewAccessor(selector)
}

type ProviderSelector func(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (types.TraceProvider, error)

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

// Get the traceIdx's accestor claims hash and its merkel proof in subValues. traceIdx must be ref's traceIdx ±1.
// Params:
// - relativeTraceIdx: the trace index relative to splitDepth
func findAncestorProofAtDepth2(ctx context.Context, provider types.TraceProvider, game types.Game, ref types.Claim, relativeTraceIdx *big.Int) (types.DAItem, error) {
	maxTraceDepth := game.MaxDepth() - game.SplitDepth()
	// Find the ancestor claim at the (splitDepth + nbits) level.
	splitLeaf, err := FindAncestorAtDepth(game, ref, types.Depth(game.TraceRootDepth()))
	if err != nil {
		return types.DAItem{}, err
	}

	ancestor := ref
	relativeAncestorPos, err := ancestor.RelativeToAncestorAtDepth(game.SplitDepth())
	if err != nil {
		return types.DAItem{}, err
	}
	minTraceIndex := relativeAncestorPos.TraceIndex(maxTraceDepth)
	offset := new(big.Int).Lsh(big.NewInt(int64(game.MaxAttackBranch())-1), uint(maxTraceDepth-relativeAncestorPos.Depth()))
	maxTraceIndex := new(big.Int).Add(minTraceIndex, offset)
	minPrePostTraceIdx := new(big.Int).Sub(minTraceIndex, big.NewInt(1))
	maxPrePostTraceIdx := new(big.Int).Add(maxTraceIndex, big.NewInt(1))
	if relativeTraceIdx.Cmp(minPrePostTraceIdx) == -1 || relativeTraceIdx.Cmp(maxPrePostTraceIdx) == 1 {
		return types.DAItem{}, fmt.Errorf("relativeTraceIdx %v is beyond ref claims's pre/post trace range [%v, %v]", relativeTraceIdx, minPrePostTraceIdx, maxPrePostTraceIdx)
	}

	// If traceIdx equals -1, the absolute prestate is returned.
	if relativeTraceIdx.Cmp(big.NewInt(-1)) == 0 {
		absolutePresate, err := provider.AbsolutePreStateCommitment(ctx)
		if err != nil {
			return types.DAItem{}, fmt.Errorf("failed to get absolutePrestate: %w", err)
		}
		return types.DAItem{
			DataHash: absolutePresate,
		}, nil
	}

	// If traceIdx is the right most branch, the root Claim is returned.
	if new(big.Int).Add(relativeTraceIdx, big.NewInt(1)).Cmp(new(big.Int).Lsh(big.NewInt(1), uint(maxTraceDepth-types.Depth(game.NBits())))) == 0 {
		return types.DAItem{
			DataHash: splitLeaf.Value,
		}, nil
	}

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
	subValues := *ancestor.SubValues
	ancestorClaim := subValues[branch]
	merkleProof := utils.GenerateProofForSubValues(subValues, uint32(branch))

	return types.DAItem{
		DataHash: ancestorClaim,
		Proof:    merkleProof,
	}, nil
}

func NewAccessor(selector ProviderSelector) *Accessor {
	return &Accessor{selector}
}

type Accessor struct {
	selector ProviderSelector
}

func (t *Accessor) Get(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (common.Hash, error) {
	provider, err := t.selector(ctx, game, ref, pos)
	if err != nil {
		return common.Hash{}, err
	}
	return provider.Get(ctx, pos)
}

func (t *Accessor) GetStepData(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (prestate []byte, proofData []byte, preimageData *types.PreimageOracleData, err error) {
	provider, err := t.selector(ctx, game, ref, pos)
	if err != nil {
		return nil, nil, nil, err
	}
	return provider.GetStepData(ctx, pos)
}

func (t *Accessor) GetStepData2(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (prestate []byte, proofData []byte, preimageData *types.PreimageOracleData, err error) {
	// Get oracle data
	provider, err := t.selector(ctx, game, ref, pos)
	if err != nil {
		return nil, nil, nil, err
	}
	// Get output root oracle for addLocaldata
	outputRootDA, err := provider.GetLocalData(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get local data: %w", err)
	}

	prestate, proofData, preimageData, err = provider.GetStepData(ctx, pos)
	if err != nil {
		return nil, nil, nil, err
	}
	// Get stepProof for stepV2
	relativePos, err := pos.RelativeToAncestorAtDepth(game.TraceRootDepth())
	if err != nil {
		return nil, nil, nil, err
	}
	postTraceIdx := relativePos.TraceIndex(game.MaxDepth() - game.TraceRootDepth())
	preTraceIdx := new(big.Int).Sub(postTraceIdx, big.NewInt(1))
	preStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, preTraceIdx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get preStateDaItem at trace index %v: %w", preTraceIdx, err)
	}
	postStateDaItem, err := findAncestorProofAtDepth2(ctx, provider, game, ref, postTraceIdx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get postStateDaItem at trace index %v: %w", postTraceIdx, err)
	}
	stateData := types.DAData{
		PreDA:  preStateDaItem,
		PostDA: postStateDaItem,
	}

	preimageData.VMStateDA = stateData

	if preimageData.OracleKey != nil {
		keyType := preimage.KeyType(preimageData.OracleKey[0])
		if keyType == preimage.LocalKeyType {
			ident := preimageData.GetIdent()
			addlocalDataDaItem := types.DAItem{}
			if ident.Cmp(big.NewInt(types.LocalPreimageKeyStartingOutputRoot)) == 0 {
				addlocalDataDaItem = outputRootDA.PreDA
			} else if ident.Cmp(big.NewInt(types.LocalPreimageKeyDisputedOutputRoot)) == 0 {
				addlocalDataDaItem = outputRootDA.PostDA
			}
			preimageData.OutputRootDAItem = addlocalDataDaItem
		}
	}

	return prestate, proofData, preimageData, nil
}

func (t *Accessor) GetL2BlockNumberChallenge(ctx context.Context, game types.Game) (*types.InvalidL2BlockNumberChallenge, error) {
	provider, err := t.selector(ctx, game, game.Claims()[0], types.RootPosition)
	if err != nil {
		return nil, err
	}
	return provider.GetL2BlockNumberChallenge(ctx)
}

var _ types.TraceAccessor = (*Accessor)(nil)
