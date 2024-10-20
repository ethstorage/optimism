package asterisc

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

func (p *AsteriscTraceProvider) GetStepData2(ctx context.Context, pos types.Position, parentPos types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	traceIndex := pos.TraceIndex(p.gameDepth)
	if !traceIndex.IsUint64() {
		return nil, nil, nil, errors.New("trace index out of bounds")
	}
	proof, err := p.loadProof(ctx, traceIndex.Uint64())
	if err != nil {
		return nil, nil, nil, err
	}
	value := ([]byte)(proof.StateData)
	if len(value) == 0 {
		return nil, nil, nil, errors.New("proof missing state data")
	}
	data := ([]byte)(proof.ProofData)
	if data == nil {
		return nil, nil, nil, errors.New("proof missing proof data")
	}
	oracleData, err := p.preimageLoader.LoadPreimage(proof)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to load preimage: %w", err)
	}
	return value, data, oracleData, nil
}
