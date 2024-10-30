package contracts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching"
	"github.com/ethereum-optimism/optimism/op-service/sources/batching/rpcblock"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func (f *FaultDisputeGameContractLatest) GetAllClaimsWithSubValues(ctx context.Context, block rpcblock.Block) ([]types.Claim, error) {
	defer f.metrics.StartContractRequest("GetAllClaims")()
	results, err := batching.ReadArray(ctx, f.multiCaller, block, f.contract.Call(methodClaimCount), func(i *big.Int) *batching.ContractCall {
		return f.contract.Call(methodClaim, i)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load claims: %w", err)
	}

	var claims []types.Claim
	for idx, result := range results {
		claim := f.decodeClaim(result, idx)
		subValues, err := f.GetSubValues(ctx, block, &claim)
		if err != nil {
			return nil, fmt.Errorf("failed to load subValues for claim %s: %w", claim, err)
		}
		claim.SetSubValues(&subValues)
		claims = append(claims, claim)
	}
	return claims, nil
}

func (f *FaultDisputeGameContractLatest) GetSubValues(ctx context.Context, block rpcblock.Block, aggClaim *types.Claim) ([]common.Hash, error) {
	defer f.metrics.StartContractRequest("GetAllSubValues")()

	filter, err := bindings.NewFaultDisputeGameFilterer(f.contract.Addr(), f.multiCaller)
	if err != nil {
		return nil, err
	}

	parentIndex := [...]*big.Int{big.NewInt(int64(aggClaim.ParentContractIndex))}
	claim := [...][32]byte{aggClaim.ClaimData.ValueBytes()}
	claimant := [...]common.Address{aggClaim.Claimant}
	moveIter, err := filter.FilterMove(nil, parentIndex[:], claim[:], claimant[:])
	if err != nil {
		return nil, fmt.Errorf("failed to filter move event log: %w", err)
	}
	ok := moveIter.Next()
	if !ok {
		return nil, fmt.Errorf("failed to get move event log: %w", moveIter.Error())
	}
	txHash := moveIter.Event.Raw.TxHash

	// todo: replace hardcoded method name
	txCall := batching.NewTxGetByHash(f.contract.Abi(), txHash, "move")
	result, err := f.multiCaller.SingleCall(ctx, rpcblock.Latest, txCall)
	if err != nil {
		return nil, fmt.Errorf("failed to load claim calldata: %w", err)
	}

	txn, err := txCall.DecodeToTx(result)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tx: %w", err)
	}

	var subValues []common.Hash

	if len(txn.BlobHashes()) > 0 {
		// todo: fetch Blobs and unpack it into subValues
		return nil, fmt.Errorf("blob tx hasn't been supported")
	} else {
		inputMap, err := txCall.UnpackCallData(txn)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack tx resp: %w", err)
		}
		// todo: replace claim with nary-subValues
		claim := *abi.ConvertType(inputMap[subClaimField], new([32]byte)).(*[32]byte)
		subValues = append(subValues, claim)
	}

	return subValues, nil
}
