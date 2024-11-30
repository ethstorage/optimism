package disputegame

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/contracts"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const getTraceTimeout = 10 * time.Minute

type OutputHonestHelper struct {
	t            *testing.T
	require      *require.Assertions
	game         *OutputGameHelper
	contract     contracts.FaultDisputeGameContract
	correctTrace types.TraceAccessor
}

func NewOutputHonestHelper(t *testing.T, require *require.Assertions, game *OutputGameHelper, contract contracts.FaultDisputeGameContract, correctTrace types.TraceAccessor) *OutputHonestHelper {
	return &OutputHonestHelper{
		t:            t,
		require:      require,
		game:         game,
		contract:     contract,
		correctTrace: correctTrace,
	}
}

func (h *OutputHonestHelper) CounterClaim(ctx context.Context, claim *ClaimHelper, opts ...MoveOpt) *ClaimHelper {
	game, target := h.loadState(ctx, claim.Index)
	value, err := h.correctTrace.Get(ctx, game, target, target.Position)
	h.require.NoErrorf(err, "Failed to determine correct claim at position %v with g index %v", target.Position, target.Position.ToGIndex())

	if value == claim.claim {
		return h.AttackClaim(ctx, claim, 1<<h.game.Nbits-1, opts...)
	} else {
		return h.AttackClaim(ctx, claim, 0, opts...)
	}
}

func (h *OutputHonestHelper) AttackClaim(ctx context.Context, claim *ClaimHelper, attackBranch uint64, opts ...MoveOpt) *ClaimHelper {
	h.Attack2(ctx, claim.Index, attackBranch, opts...)
	return claim.WaitForCounterClaim(ctx)
}

func (h *OutputHonestHelper) DefendClaim(ctx context.Context, claim *ClaimHelper, opts ...MoveOpt) *ClaimHelper {
	h.Defend(ctx, claim.Index, opts...)
	return claim.WaitForCounterClaim(ctx)
}

func (h *OutputHonestHelper) Attack2(ctx context.Context, claimIdx int64, attackBranch uint64, opts ...MoveOpt) {
	// Ensure the claim exists
	h.game.WaitForClaimCount(ctx, claimIdx+1)

	ctx, cancel := context.WithTimeout(ctx, getTraceTimeout)
	defer cancel()

	game, claim := h.loadState(ctx, claimIdx)
	attackPos := claim.Position.MoveN(h.game.Nbits, attackBranch)
	h.t.Logf("Attacking claim %v at position %v with g index %v", claimIdx, attackPos, attackPos.ToGIndex())
	subValues := []common.Hash{}
	for i := uint64(0); i < 1<<h.game.Nbits-1; i++ {
		value, err := h.correctTrace.Get(ctx, game, claim, attackPos.MoveRightN(i))
		h.require.NoErrorf(err, "Get correct claim at position %v with g index %v", attackPos, attackPos.ToGIndex())
		subValues = append(subValues, value)
	}
	h.t.Log("Performing attack")
	h.game.Attack2(ctx, claimIdx, attackBranch, subValues, opts...)
	h.t.Log("Attack complete")
}

func (h *OutputHonestHelper) Attack(ctx context.Context, claimIdx int64, opts ...MoveOpt) {
	panic("unimplemented, use attack2 instead")
}

func (h *OutputHonestHelper) Defend(ctx context.Context, claimIdx int64, opts ...MoveOpt) {
	panic("unimplemented, use attack2 instead")
}

func (h *OutputHonestHelper) StepClaimFails(ctx context.Context, claim *ClaimHelper, isAttack bool) {
	h.StepFails(ctx, claim.Index, isAttack)
}

func (h *OutputHonestHelper) StepFails(ctx context.Context, claimIdx int64, isAttack bool) {
	// Ensure the claim exists
	h.game.WaitForClaimCount(ctx, claimIdx+1)

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	game, claim := h.loadState(ctx, claimIdx)
	pos := claim.Position
	if !isAttack {
		// If we're defending, then the step will be from the trace to the next one
		pos = pos.MoveRight()
	}
	prestate, proofData, _, err := h.correctTrace.GetStepData(ctx, game, claim, pos)
	h.require.NoError(err, "Get step data")
	h.game.StepFails(claimIdx, isAttack, prestate, proofData)
}

func (h *OutputHonestHelper) loadState(ctx context.Context, claimIdx int64) (types.Game, types.Claim) {
	claims, err := h.contract.GetAllClaimsWithSubValues(ctx)
	h.require.NoError(err, "Failed to load claims from game")
	game := types.NewGameState2(claims, h.game.MaxDepth(ctx), h.game.Nbits, h.game.SplitDepth(ctx))

	claim := game.Claims()[claimIdx]
	return game, claim
}
