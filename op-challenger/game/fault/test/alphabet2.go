package test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/alphabet"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

type AlphabetWithOutputRootProvider struct {
	*alphabet.AlphabetTraceProvider
	depth            types.Depth
	OracleError      error
	L2BlockChallenge *types.InvalidL2BlockNumberChallenge
	pre              bool // preOutputRoot if pre is true; postOutputRoot if pre is false
}

func (a *AlphabetWithOutputRootProvider) GetStepData(ctx context.Context, i types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	preimage, _, _, err := a.AlphabetTraceProvider.GetStepData(ctx, i)
	if err != nil {
		return nil, nil, nil, err
	}
	traceIndex := i.TraceIndex(a.depth).Uint64()
	var key []byte
	if a.pre {
		// localPreimageKey(1) + STARTING_OUTPUT_ROOT(2)
		key = []byte{0x01, 0x02}
	} else {
		// localPreimageKey(1) + DISPUTED_OUTPUT_ROOT(3)
		key = []byte{0x01, 0x03}
	}
	data := types.NewPreimageOracleData(key, []byte{byte(traceIndex - 1)}, uint32(traceIndex-1))
	return preimage, []byte{byte(traceIndex - 1)}, data, nil
}

func NewAlphabetWithProofProvider2(t *testing.T, startingL2BlockNumber *big.Int, maxDepth types.Depth, rootDepth types.Depth, oracleError error) *AlphabetWithProofProvider {
	return &AlphabetWithProofProvider{
		alphabet.NewTraceProvider2(startingL2BlockNumber, maxDepth, rootDepth),
		maxDepth,
		oracleError,
		nil,
	}
}

func NewAlphabetWithPreoutputRootProvider(t *testing.T, startingL2BlockNumber *big.Int, maxDepth types.Depth, rootDepth types.Depth, pre bool, oracleError error) *AlphabetWithOutputRootProvider {
	return &AlphabetWithOutputRootProvider{
		alphabet.NewTraceProvider2(startingL2BlockNumber, maxDepth, rootDepth),
		maxDepth,
		oracleError,
		nil,
		pre,
	}
}
