package cannon

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

func (p *CannonTraceProvider) GetStepData2(ctx context.Context, pos types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	return nil, nil, nil, fmt.Errorf("cannon GetStepData2 is not supported, use GetStepData instead")
}
