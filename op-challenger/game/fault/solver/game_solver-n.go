//go:build faultdisputegamen
// +build faultdisputegamen

package solver

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
)

type GameSolver struct {
	claimSolver *claimSolver
}

func NewGameSolver(gameDepth types.Depth, trace types.TraceAccessor) *GameSolver {
	return &GameSolver{
		claimSolver: newClaimSolver(gameDepth, trace),
	}
}

func (s *GameSolver) AgreeWithRootClaim(ctx context.Context, game types.Game) (bool, error) {
	return s.claimSolver.agreeWithClaim(ctx, game, game.Claims()[0])
}

func (s *GameSolver) CalculateNextActions(ctx context.Context, game types.Game) ([]types.Action, error) {
	agreeWithRootClaim, err := s.AgreeWithRootClaim(ctx, game)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if root claim is correct: %w", err)
	}

	// Challenging the L2 block number will only work if we have the same output root as the claim
	// Otherwise our output root preimage won't match. We can just proceed and invalidate the output root by disputing claims instead.
	if agreeWithRootClaim {
		if challenge, err := s.claimSolver.trace.GetL2BlockNumberChallenge(ctx, game); errors.Is(err, types.ErrL2BlockNumberValid) {
			// We agree with the L2 block number, proceed to processing claims
		} else if err != nil {
			// Failed to check L2 block validity
			return nil, fmt.Errorf("failed to determine L2 block validity: %w", err)
		} else {
			return []types.Action{
				{
					Type:                          types.ActionTypeChallengeL2BlockNumber,
					InvalidL2BlockNumberChallenge: challenge,
				},
			}, nil
		}
	}

	var actions []types.Action
	agreedClaims := newHonestClaimTracker()
	if agreeWithRootClaim {
		agreedClaims.AddHonestClaim(types.Claim{}, game.Claims()[0])
	}
	for _, claim := range game.Claims() {
		var action *types.Action
		subValues, err := s.getClaimRealValues(ctx, game, claim)
		if err != nil {
			return nil, fmt.Errorf("failed to get real values for claim %v: %w", claim.ContractIndex, err)
		}
		claimV2 := types.ClaimV2{
			Claim:     claim,
			SubValues: subValues,
		}
		for branch, _ := range subValues {
			if claim.Depth() == game.MaxDepth() {
				// TODO: implement
				action, err = s.calculateStep(ctx, game, claim, agreedClaims)
			} else {
				action, err = s.calculateMove(ctx, game, claimV2, agreedClaims, uint64(branch))
			}
			if err != nil {
				// Unable to continue iterating claims safely because we may not have tracked the required honest moves
				// for this claim which affects the response to later claims.
				// Any actions we've already identified are still safe to apply.
				return actions, fmt.Errorf("failed to determine response to claim %v: %w", claim.ContractIndex, err)
			}
			if action == nil {
				continue
			}
			actions = append(actions, *action)
			break
		}
	}
	return actions, nil
}

func (s *GameSolver) calculateStep(ctx context.Context, game types.Game, claim types.Claim, agreedClaims *honestClaimTracker) (*types.Action, error) {
	if claim.CounteredBy != (common.Address{}) {
		return nil, nil
	}
	step, err := s.claimSolver.AttemptStep(ctx, game, claim, agreedClaims)
	if err != nil {
		return nil, err
	}
	if step == nil {
		return nil, nil
	}
	return &types.Action{
		Type:        types.ActionTypeStep,
		ParentClaim: step.LeafClaim,
		IsAttack:    step.IsAttack,
		PreState:    step.PreState,
		ProofData:   step.ProofData,
		OracleData:  step.OracleData,
	}, nil
}

func (s *GameSolver) calculateMove(ctx context.Context, game types.Game, claimV2 types.ClaimV2, honestClaims *honestClaimTracker, branch uint64) (*types.Action, error) {
	move, err := s.claimSolver.NextMove(ctx, claimV2, game, honestClaims, branch)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate next move for claim index %v: %w", claimV2.Claim.ContractIndex, err)
	}
	if move == nil {
		return nil, nil
	}
	honestClaims.AddHonestClaim(claimV2.Claim, move.Claim)
	if game.IsDuplicate(move.Claim) {
		return nil, nil
	}
	return &types.Action{
		Type:        types.ActionTypeAttack,
		IsAttack:    !game.DefendsParent(move.Claim),
		ParentClaim: game.Claims()[move.Claim.ParentContractIndex],
		Value:       move.Claim.Value,
		SubValues:   move.SubValues,
	}, nil
}

func (s *GameSolver) getClaimRealValues(ctx context.Context, game types.Game, claim types.Claim) ([]common.Hash, error) {
	values := []common.Hash{claim.Value}
	// TODO: implement
	return values, nil
}
