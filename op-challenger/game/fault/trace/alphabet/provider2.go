package alphabet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

// NewTraceProvider returns a new [AlphabetProvider].
func NewTraceProvider2(startingBlockNumber *big.Int, depth types.Depth, rootDepth types.Depth) *AlphabetTraceProvider {
	return &AlphabetTraceProvider{
		startingBlockNumber: startingBlockNumber,
		depth:               depth,
		maxLen:              1 << depth,
		rootDepth:           rootDepth,
	}
}

func (ap *AlphabetTraceProvider) GetStepData2(ctx context.Context, pos types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	return nil, nil, nil, fmt.Errorf("alphabet GetStepData2 is not supported, use GetStepData instead")
}
