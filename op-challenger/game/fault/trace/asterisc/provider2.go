package asterisc

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

func (p *AsteriscTraceProvider) GetStepData2(ctx context.Context, pos types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	return nil, nil, nil, fmt.Errorf("asterisc GetStepData2 is not supported, use GetStepData instead")
}
