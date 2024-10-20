package solver

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
)

// AttemptStep determines what step, if any, should occur for a given leaf claim.
// An error will be returned if the claim is not at the max depth.
// Returns nil, nil if no step should be performed.
func (s *claimSolver) AttemptStep2(ctx context.Context, game types.Game, claim types.Claim, honestClaims *honestClaimTracker) (*StepData, error) {
	if claim.Depth() != s.gameDepth {
		return nil, ErrStepNonLeafNode
	}

	if counter, err := s.shouldCounter(game, claim, honestClaims); err != nil {
		return nil, fmt.Errorf("failed to determine if claim should be countered: %w", err)
	} else if !counter {
		return nil, nil
	}

	claimCorrect, err := s.agreeWithClaim(ctx, game, claim)
	if err != nil {
		return nil, err
	}

	// prestateItem
	// postStateItem

	var position types.Position
	if !claimCorrect {
		// Attack the claim by executing step index, so we need to get the pre-state of that index
		position = claim.Position
	} else {
		// Defend and use this claim as the starting point to execute the step after.
		// Thus, we need the pre-state of the next step.
		position = claim.Position.MoveRight()
	}

	preState, proofData, oracleData, err := s.trace.GetStepData(ctx, game, claim, position)
	if err != nil {
		return nil, err
	}

	return &StepData{
		LeafClaim:  claim,
		IsAttack:   !claimCorrect,
		PreState:   preState,
		ProofData:  proofData,
		OracleData: oracleData,
	}, nil
}


func findPrestateItem()