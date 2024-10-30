package trace

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

func (p *TranslatingProvider) GetStepData2(ctx context.Context, pos types.Position) (prestate []byte, proofData []byte, preimageData *types.PreimageOracleData, err error) {
	return nil, nil, nil, fmt.Errorf("translate GetStepData2 is not supported, use GetStepData instead")
}
