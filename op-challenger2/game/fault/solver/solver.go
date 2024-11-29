package solver

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/contracts"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrStepNonLeafNode = errors.New("cannot step on non-leaf claims")
)

// claimSolver uses a [TraceProvider] to determine the moves to make in a dispute game.
type claimSolver struct {
	trace     types.TraceAccessor
	gameDepth types.Depth
	daType    types.DAType
}

// newClaimSolver creates a new [claimSolver] using the provided [TraceProvider].
func newClaimSolver(gameDepth types.Depth, trace types.TraceAccessor, daType types.DAType) *claimSolver {
	return &claimSolver{
		trace,
		gameDepth,
		daType,
	}
}

func (s *claimSolver) shouldCounter(game types.Game, claim types.Claim, honestClaims *honestClaimTracker) (bool, error) {
	// Do not counter honest claims
	if honestClaims.IsHonest(claim) {
		return false, nil
	}

	if claim.IsRoot() {
		// Always counter the root claim if it is not honest
		return true, nil
	}

	parent, err := game.GetParent(claim)
	if err != nil {
		return false, fmt.Errorf("no parent for claim %v: %w", claim.ContractIndex, err)
	}

	// Counter all claims that are countering an honest claim
	if honestClaims.IsHonest(parent) {
		return true, nil
	}

	counter, hasCounter := honestClaims.HonestCounter(parent)
	// Do not respond to any claim countering a claim the honest actor ignored
	if !hasCounter {
		return false, nil
	}

	// Do not counter sibling to an honest claim that are right of the honest claim.
	honestIdx := counter.TraceIndex(game.MaxDepth())
	claimIdx := claim.TraceIndex(game.MaxDepth())
	return claimIdx.Cmp(honestIdx) <= 0, nil
}

// NextMove returns the next move to make given the current state of the game.
func (s *claimSolver) NextMove(ctx context.Context, claim types.Claim, game types.Game, honestClaims *honestClaimTracker, branch uint64) (*types.Claim, error) {
	if claim.Depth() == s.gameDepth {
		return nil, types.ErrGameDepthReached
	}

	if counter, err := s.shouldCounter(game, claim, honestClaims); err != nil {
		return nil, fmt.Errorf("failed to determine if claim should be countered: %w", err)
	} else if !counter {
		return nil, nil
	}

	agree, err := s.agreeWithClaimV2(ctx, game, claim, branch)
	if err != nil {
		return nil, err
	}
	if agree {
		if claim.Depth() == game.TraceRootDepth() && branch == 0 {
			// Only useful in alphabet game, because the alphabet game has a constant status byte, and is not safe from someone being dishonest in
			// output bisection and then posting a correct execution trace bisection root claim.
			// when root claim of output bisection is dishonest,
			// root claim of execution trace bisection is made by the dishonest actor but is honest
			// we should counter it.
			return s.attackV2(ctx, game, claim, branch)
		}
		if branch < game.MaxAttackBranch()-1 {
			return nil, nil
		}
		return s.attackV2(ctx, game, claim, branch+1)
	} else {
		return s.attackV2(ctx, game, claim, branch)
	}
}

type StepData struct {
	LeafClaim    types.Claim
	AttackBranch uint64
	PreState     []byte
	ProofData    []byte
	OracleData   *types.PreimageOracleData
}

// AttemptStep determines what step, if any, should occur for a given leaf claim.
// An error will be returned if the claim is not at the max depth.
// Returns nil, nil if no step should be performed.
func (s *claimSolver) AttemptStep(ctx context.Context, game types.Game, claim types.Claim, honestClaims *honestClaimTracker, branch uint64) (*StepData, error) {
	if claim.Depth() != s.gameDepth {
		return nil, ErrStepNonLeafNode
	}

	if counter, err := s.shouldCounter(game, claim, honestClaims); err != nil {
		return nil, fmt.Errorf("failed to determine if claim should be countered: %w", err)
	} else if !counter {
		return nil, nil
	}

	claimCorrect, err := s.agreeWithClaimV2(ctx, game, claim, branch)
	if err != nil {
		return nil, err
	}

	var position types.Position
	attackBranch := branch
	if !claimCorrect {
		// Attack the claim by executing step index, so we need to get the pre-state of that index
		position = claim.Position.MoveRightN(branch)
	} else {
		if branch == game.MaxAttackBranch()-1 {
			// If we are at the max attack branch, we need to step on the next branch
			position = claim.Position.MoveRightN(branch + 1)
			attackBranch = branch + 1
		} else {
			return nil, nil
		}
	}
	preState, proofData, oracleData, err := s.trace.GetStepData2(ctx, game, claim, position)
	if err != nil {
		return nil, err
	}
	return &StepData{
		LeafClaim:    claim,
		AttackBranch: attackBranch,
		PreState:     preState,
		ProofData:    proofData,
		OracleData:   oracleData,
	}, nil
}

// agreeWithClaim returns true if the claim is correct according to the internal [TraceProvider].
func (s *claimSolver) agreeWithClaimV2(ctx context.Context, game types.Game, claim types.Claim, branch uint64) (bool, error) {
	if branch >= uint64(len(*claim.SubValues)) {
		return true, fmt.Errorf("branch must be less than maxAttachBranch")
	}
	ourValue, err := s.trace.Get(ctx, game, claim, claim.Position.MoveRightN(branch))
	return bytes.Equal(ourValue[:], (*claim.SubValues)[branch][:]), err
}

func (s *claimSolver) attackV2(ctx context.Context, game types.Game, claim types.Claim, branch uint64) (*types.Claim, error) {
	var err error
	var value common.Hash
	var values []common.Hash
	maxAttackBranch := game.MaxAttackBranch()
	position := claim.MoveN(game.NBits(), branch)
	for i := uint64(0); i < maxAttackBranch; i++ {
		tmpPosition := position.MoveRightN(i)
		if tmpPosition.Depth() == (game.SplitDepth()+types.Depth(game.NBits())) && i != 0 {
			value = common.Hash{}
		} else {
			value, err = s.trace.Get(ctx, game, claim, tmpPosition)
			if err != nil {
				return nil, fmt.Errorf("attack claim: %w", err)
			}
		}
		values = append(values, value)
	}
	hash := contracts.SubValuesHash(values)
	return &types.Claim{
		ClaimData:           types.ClaimData{Value: hash, Position: position},
		ParentContractIndex: claim.ContractIndex,
		SubValues:           &values,
		AttackBranch:        branch,
	}, nil
}
