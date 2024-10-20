package alphabet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum/go-ethereum/common"
)

func (ap *AlphabetTraceProvider) GetStepData2(ctx context.Context, pos types.Position, parentPos types.Position) ([]byte, []byte, *types.PreimageOracleData, error) {
	traceIndex := pos.TraceIndex(ap.depth)
	key := preimage.LocalIndexKey(L2ClaimBlockNumberLocalIndex).PreimageKey()
	preimageData := types.NewPreimageOracleData(key[:], ap.startingBlockNumber.Bytes(), 0)
	if traceIndex.Cmp(common.Big0) == 0 {
		return absolutePrestate, []byte{}, preimageData, nil
	}
	if traceIndex.Cmp(new(big.Int).SetUint64(ap.maxLen)) > 0 {
		return nil, nil, nil, fmt.Errorf("%w depth: %v index: %v max: %v", ErrIndexTooLarge, ap.depth, traceIndex, ap.maxLen)
	}
	initialTraceIndex := new(big.Int).Lsh(ap.startingBlockNumber, uint(ap.splitDepth))
	initialClaim := new(big.Int).Add(absolutePrestateInt, initialTraceIndex)
	newTraceIndex := new(big.Int).Add(initialTraceIndex, traceIndex)
	// newClaim is always, otherwise we would counter the first false subClaim
	newClaim := new(big.Int).Add(initialClaim, traceIndex)
	return BuildAlphabetPreimage(newTraceIndex, newClaim), []byte{}, preimageData, nil
}
