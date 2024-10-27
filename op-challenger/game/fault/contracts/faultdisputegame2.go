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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

func (f *FaultDisputeGameContractLatest) GetSubClaims(ctx context.Context, block rpcblock.Block, aggClaim *types.Claim) ([]common.Hash, error) {
	defer f.metrics.StartContractRequest("GetAllSubClaims")()

	filter, err := bindings.NewFaultDisputeGameFilterer(f.contract.Addr(), f.multiCaller)
	if err != nil {
		return nil, err
	}

	parentIndex := [...]*big.Int{big.NewInt(int64(aggClaim.ParentContractIndex))}
	claim := [...][32]byte{aggClaim.ClaimData.ValueBytes()}
	claimant := [...]common.Address{aggClaim.Claimant}
	var end *uint64
	if block.ArgValue() == rpcblock.Latest {
		end = nil
	} else {
		blockNumber, ok := block.ArgValue().(rpc.BlockNumber)
		if !ok {
			return nil, fmt.Errorf("block number is not lastest or int64")
		}
		blockNumberU64 := uint64(blockNumber)
		end = &blockNumberU64
	}
	moveIter, err := filter.FilterMove(&bind.FilterOpts{End: end, Context: ctx}, parentIndex[:], claim[:], claimant[:])
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
	result, err := f.multiCaller.SingleCall(ctx, block, txCall)
	if err != nil {
		return nil, fmt.Errorf("failed to load claim calldata: %w", err)
	}

	txn := result.GetTx()

	var subClaims []common.Hash

	if len(txn.BlobHashes()) > 0 {
		// todo: fetch Blobs and unpack it into subClaims
		return nil, fmt.Errorf("blob tx hasn't been supported")
	} else {
		inputMap, err := txCall.UnpackCallData(&txn)
		if err != nil {
			return nil, fmt.Errorf("failed to unpack tx resp: %w", err)
		}
		// todo: replace claim with nary-subclaims
		claim := *abi.ConvertType(inputMap[subClaimField], new([32]byte)).(*[32]byte)
		subClaims = append(subClaims, claim)
	}

	return subClaims, nil
}
