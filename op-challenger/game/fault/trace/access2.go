package trace

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

func (t *Accessor) GetStepData2(ctx context.Context, game types.Game, ref types.Claim, pos types.Position) (prestate []byte, proofData []byte, preimageData *types.PreimageOracleData, err error) {
	provider, err := t.selector(ctx, game, ref, pos)
	if err != nil {
		return nil, nil, nil, err
	}
	parentPos, err := game.GetParent(ref)
	if err != nil {
		return nil, nil, nil, err
	}
	return provider.GetStepData2(ctx, pos, parentPos.Position)
}
