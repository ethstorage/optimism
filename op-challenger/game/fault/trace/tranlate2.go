package trace

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

// fix it with real one
func (p *TranslatingProvider) GetStepData2(ctx context.Context, pos types.Position, parentPos types.Position) (prestate []byte, proofData []byte, preimageData *types.PreimageOracleData, err error) {
	relativePos, err := pos.RelativeToAncestorAtDepth(p.rootDepth)
	if err != nil {
		return nil, nil, nil, err
	}
	return p.provider.GetStepData2(ctx, relativePos, parentPos)
}
