package test

import (
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/trace/alphabet"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

func NewAlphabetWithProofProvider2(t *testing.T, startingL2BlockNumber *big.Int, maxDepth types.Depth, rootDepth types.Depth, oracleError error) *AlphabetWithProofProvider {
	return &AlphabetWithProofProvider{
		alphabet.NewTraceProvider2(startingL2BlockNumber, maxDepth, rootDepth),
		maxDepth,
		oracleError,
		nil,
	}
}
