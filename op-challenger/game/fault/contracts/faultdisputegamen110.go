package contracts

import (
	"context"
	_ "embed"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
)

//go:embed abis/FaultDisputeGameN-1.1.0.json
var faultDisputeGameNAbi110 []byte

type FaultDisputeGameContractN110 struct {
	FaultDisputeGameContractLatest
}

func (f *FaultDisputeGameContractN110) AttackTx(ctx context.Context, parent types.Claim, pivot common.Hash) (txmgr.TxCandidate, error) {
	call := f.contract.Call(methodAttack, big.NewInt(int64(parent.ContractIndex)), pivot)
	return f.txWithBond(ctx, parent.Position.Attack(), call)
}

func (f *FaultDisputeGameContractN110) DefendTx(ctx context.Context, parent types.Claim, pivot common.Hash) (txmgr.TxCandidate, error) {
	call := f.contract.Call(methodDefend, big.NewInt(int64(parent.ContractIndex)), pivot)
	return f.txWithBond(ctx, parent.Position.Defend(), call)
}
