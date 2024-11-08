package test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace/alphabet"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	oppreimage "github.com/ethereum-optimism/optimism/op-preimage"
)

type TestKeyType int

const (
	OraclePreKey     TestKeyType = iota
	OraclPostKey     TestKeyType = iota
	OracleDefaultKey TestKeyType = iota
)

func NewAlphabetWithProofProvider(t *testing.T, startingL2BlockNumber *big.Int, maxDepth types.Depth, oracleError error, rootDepth types.Depth, oracleKey TestKeyType) *AlphabetWithProofProvider {
	return &AlphabetWithProofProvider{
		alphabet.NewTraceProvider(startingL2BlockNumber, maxDepth),
		maxDepth,
		oracleError,
		nil,
		oracleKey,
	}
}

func NewAlphabetClaimBuilder2(t *testing.T, startingL2BlockNumber *big.Int, maxDepth types.Depth, nbits uint64, splitDepth types.Depth) *ClaimBuilder {
	alphabetProvider := NewAlphabetWithProofProvider(t, startingL2BlockNumber, maxDepth, nil, splitDepth+types.Depth(nbits), OracleDefaultKey)
	return NewClaimBuilder2(t, maxDepth, nbits, splitDepth, alphabetProvider)
}

type AlphabetWithProofProvider struct {
	*alphabet.AlphabetTraceProvider
	depth            types.Depth
	OracleError      error
	L2BlockChallenge *types.InvalidL2BlockNumberChallenge
	pre              TestKeyType
}

func (a *AlphabetWithProofProvider) GetStepData(ctx context.Context, i types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	preimage, _, _, err := a.AlphabetTraceProvider.GetStepData(ctx, i)
	if err != nil {
		return nil, nil, nil, err
	}
	traceIndex := i.TraceIndex(a.depth).Uint64()
	var key []byte
	switch a.pre {
	case OraclePreKey:
		// localPreimageKey(1) + STARTING_OUTPUT_ROOT(2)
		key = []byte{byte(oppreimage.LocalKeyType), types.LocalPreimageKeyStartingOutputRoot}
	case OraclPostKey:
		// localPreimageKey(1) + DISPUTED_OUTPUT_ROOT(3)
		key = []byte{byte(oppreimage.LocalKeyType), types.LocalPreimageKeyDisputedOutputRoot}
	case OracleDefaultKey:
		key = []byte{byte(traceIndex)}
	}
	data := types.NewPreimageOracleData(key, []byte{byte(traceIndex - 1)}, uint32(traceIndex-1))
	return preimage, []byte{byte(traceIndex - 1)}, data, nil
}

func (ap *AlphabetWithProofProvider) GetStepData2(ctx context.Context, pos types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	return nil, nil, nil, fmt.Errorf("alphabetWithProofProvider GetStepData2 is not supported, use GetStepData instead")
}

func (c *AlphabetWithProofProvider) GetL2BlockNumberChallenge(_ context.Context) (*types.InvalidL2BlockNumberChallenge, error) {
	if c.L2BlockChallenge != nil {
		return c.L2BlockChallenge, nil
	} else {
		return nil, types.ErrL2BlockNumberValid
	}
}
