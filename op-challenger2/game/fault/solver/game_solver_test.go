package solver

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	faulttest "github.com/ethereum-optimism/optimism/op-challenger2/game/fault/test"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/trace"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestCalculateNextActions_ChallengeL2BlockNumber(t *testing.T) {
	startingBlock := big.NewInt(5)
	maxDepth := types.Depth(8)
	challenge := &types.InvalidL2BlockNumberChallenge{
		Output: &eth.OutputResponse{OutputRoot: eth.Bytes32{0xbb}},
	}
	nbits := uint64(2)
	splitDepth := types.Depth(4)
	claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, startingBlock, maxDepth, nbits, splitDepth)
	traceProvider := faulttest.NewAlphabetWithProofProvider(t, startingBlock, maxDepth, nil, 0, faulttest.OracleDefaultKey)
	solver := NewGameSolver(maxDepth, trace.NewSimpleTraceAccessor(traceProvider), types.CallDataType)

	// Do not challenge when provider returns error indicating l2 block is valid
	actions, err := solver.CalculateNextActions(context.Background(), claimBuilder.GameBuilder().Game)
	require.NoError(t, err)
	require.Len(t, actions, 0)

	// Do challenge when the provider returns a challenge
	traceProvider.L2BlockChallenge = challenge
	actions, err = solver.CalculateNextActions(context.Background(), claimBuilder.GameBuilder().Game)
	require.NoError(t, err)
	require.Len(t, actions, 1)
	action := actions[0]
	require.Equal(t, types.ActionTypeChallengeL2BlockNumber, action.Type)
	require.Equal(t, challenge, action.InvalidL2BlockNumberChallenge)
}

func runStep(t *testing.T, solver *GameSolver, game types.Game, correctTraceProvider types.TraceProvider) (types.Game, []types.Action) {
	actions, err := solver.CalculateNextActions(context.Background(), game)
	require.NoError(t, err)

	postState := applyActions(game, challengerAddr, actions)

	for i, action := range actions {
		t.Logf("Move %v: Type: %v, ParentIdx: %v, Attack: %v, Value: %v, PreState: %v, ProofData: %v",
			i, action.Type, action.ParentClaim.ContractIndex, action.IsAttack, action.Value, hex.EncodeToString(action.PreState), hex.EncodeToString(action.ProofData))
		// Check that every move the solver returns meets the generic validation rules
		require.NoError(t, checkRules(game, action, correctTraceProvider), "Attempting to perform invalid action")
	}
	return postState, actions
}

func TestMultipleRoundsWithNbits1(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		actor actor
	}{
		{
			name:  "SingleRoot",
			actor: doNothingActor,
		},
		{
			name:  "LinearAttackCorrect",
			actor: correctAttackLastClaim,
		},
		{
			name:  "LinearDefendCorrect",
			actor: correctDefendLastClaim,
		},
		{
			name:  "LinearAttackIncorrect",
			actor: incorrectAttackLastClaim,
		},
		{
			name:  "LinearDefendInorrect",
			actor: incorrectDefendLastClaim,
		},
		{
			name:  "LinearDefendIncorrectDefendCorrect",
			actor: combineActors(incorrectDefendLastClaim, correctDefendLastClaim),
		},
		{
			name:  "LinearAttackIncorrectDefendCorrect",
			actor: combineActors(incorrectAttackLastClaim, correctDefendLastClaim),
		},
		{
			name:  "LinearDefendIncorrectDefendIncorrect",
			actor: combineActors(incorrectDefendLastClaim, incorrectDefendLastClaim),
		},
		{
			name:  "LinearAttackIncorrectDefendIncorrect",
			actor: combineActors(incorrectAttackLastClaim, incorrectDefendLastClaim),
		},
		{
			name:  "AttackEverythingCorrect",
			actor: attackEverythingCorrect,
		},
		{
			name:  "DefendEverythingCorrect",
			actor: defendEverythingCorrect,
		},
		{
			name:  "AttackEverythingIncorrect",
			actor: attackEverythingIncorrect,
		},
		{
			name:  "DefendEverythingIncorrect",
			actor: defendEverythingIncorrect,
		},
		{
			name:  "Exhaustive",
			actor: exhaustive,
		},
	}
	for _, test := range tests {
		test := test
		for _, rootClaimCorrect := range []bool{true, false} {
			rootClaimCorrect := rootClaimCorrect
			t.Run(fmt.Sprintf("%v-%v", test.name, rootClaimCorrect), func(t *testing.T) {
				t.Parallel()

				maxDepth := types.Depth(6)
				startingL2BlockNumber := big.NewInt(50)
				nbits := uint64(1)
				splitDepth := types.Depth(3)
				claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, startingL2BlockNumber, maxDepth, nbits, splitDepth)
				builder := claimBuilder.GameBuilder(faulttest.WithInvalidValue(!rootClaimCorrect))
				game := builder.Game

				correctTrace := claimBuilder.CorrectTraceProvider()
				solver := NewGameSolver(maxDepth, trace.NewSimpleTraceAccessor(correctTrace), types.CallDataType)

				roundNum := 0
				done := false
				for !done {
					t.Logf("------ ROUND %v ------", roundNum)
					game, _ = runStep2(t, solver, game, correctTrace)
					verifyGameRules(t, game, rootClaimCorrect)

					game, done = test.actor.Apply(t, game, correctTrace)
					roundNum++
				}
			})
		}
	}
}

func applyActions(game types.Game, claimant common.Address, actions []types.Action) types.Game {
	claims := game.Claims()
	for _, action := range actions {
		switch action.Type {
		case types.ActionTypeMove:
			newPosition := action.ParentClaim.Position.Attack()
			if !action.IsAttack {
				newPosition = action.ParentClaim.Position.Defend()
			}
			claim := types.Claim{
				ClaimData: types.ClaimData{
					Value:    action.Value,
					Bond:     big.NewInt(0),
					Position: newPosition,
				},
				Claimant:            claimant,
				ContractIndex:       len(claims),
				ParentContractIndex: action.ParentClaim.ContractIndex,
			}
			claims = append(claims, claim)
		case types.ActionTypeStep:
			counteredClaim := claims[action.ParentClaim.ContractIndex]
			counteredClaim.CounteredBy = claimant
			claims[action.ParentClaim.ContractIndex] = counteredClaim
		default:
			panic(fmt.Errorf("unknown move type: %v", action.Type))
		}
	}
	return types.NewGameState2(claims, game.MaxDepth(), game.NBits(), game.SplitDepth())
}

func TestCalculateNextActions(t *testing.T) {
	maxDepth := types.Depth(8)
	startingL2BlockNumber := big.NewInt(0)
	nbits := uint64(2)
	splitDepth := types.Depth(4)
	claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, startingL2BlockNumber, maxDepth, nbits, splitDepth)

	tests := []struct {
		name             string
		rootClaimCorrect bool
		setupGame        func(builder *faulttest.GameBuilder)
	}{
		{
			name: "AttackRootClaim",
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().ExpectAttackV2(0)
			},
		},
		{
			name:             "DoNotAttackCorrectRootClaim_AgreeWithOutputRoot",
			rootClaimCorrect: true,
			setupGame:        func(builder *faulttest.GameBuilder) {},
		},
		{
			name: "DoNotPerformDuplicateMoves",
			setupGame: func(builder *faulttest.GameBuilder) {
				builder.Seq().Attack2(nil, 0)
			},
		},
		{
			name: "RespondToAllClaimsAtDisagreeingLevel",
			setupGame: func(builder *faulttest.GameBuilder) {
				honestClaim := builder.Seq().Attack2(nil, 0)
				honestClaim.Attack2(nil, 0).ExpectAttackV2(3)
			},
		},
		{
			name: "StepAtMaxDepth",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 3). // dishonest
					Attack2(nil, 0)     // honest
				lastHonestClaim.Attack2([]common.Hash{{0x1}, {0x2}, {0x3}}, 0).ExpectStepV2(0)
			},
		},
		{
			name: "PoisonedPreState",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				honest := builder.Seq().Attack2(nil, 0)
				dishonest := honest.Attack2(values, 0)
				dishonest.ExpectAttackV2(0)
				dishonest.Attack2(values, 0)
				dishonest.ExpectAttackV2(0)
			},
		},
		{
			name:             "HonestRoot-OneLevelAttack",
			rootClaimCorrect: true,
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				honest := builder.Seq()
				honest.Attack2(nil, 0).ExpectAttackV2(3)
				honest.Attack2(values, 0).ExpectAttackV2(0)
				values = claimBuilder.GetClaimsAtPosition(builder.Game.RootClaim().Position.MoveN(nbits, 0), []uint64{1})
				honest.Attack2(values, 0).ExpectAttackV2(1)
				values = claimBuilder.GetClaimsAtPosition(builder.Game.RootClaim().Position.MoveN(nbits, 0), []uint64{2})
				honest.Attack2(values, 0).ExpectAttackV2(2)
			},
		},
		{
			name: "DishonestRoot-OneLevelAttack",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				dishonest := builder.Seq().ExpectAttackV2(0)
				dishonest.Attack2(values, 0).ExpectAttackV2(0)
				values = claimBuilder.GetClaimsAtPosition(builder.Game.RootClaim().Position.MoveN(nbits, 0), []uint64{1})
				dishonest.Attack2(values, 0).ExpectAttackV2(1)
				values = claimBuilder.GetClaimsAtPosition(builder.Game.RootClaim().Position.MoveN(nbits, 0), []uint64{2})
				dishonest.Attack2(values, 0).ExpectAttackV2(2)
			},
		},
		{
			name:             "HonestRoot-TwoLevelAttack-FirstLevelCorrect-SecondLevelCorrect",
			rootClaimCorrect: true,
			setupGame: func(builder *faulttest.GameBuilder) {
				dishonest := builder.Seq().Attack2(nil, 0)
				dishonest.Attack2(nil, 0).ExpectAttackV2(3)
				dishonest.Attack2(nil, 1).ExpectAttackV2(3)
				dishonest.Attack2(nil, 2).ExpectAttackV2(3)
				dishonest.Attack2(nil, 3)
			},
		},
		{
			name:             "HonestRoot-TwoLevelAttack-FirstLevelCorrect-SecondLevelIncorrect",
			rootClaimCorrect: true,
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				dishonest := builder.Seq().Attack2(nil, 0).ExpectAttackV2(3)
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				dishonest.Attack2(values, 0).ExpectAttackV2(0)
				dishonest.Attack2(values, 1).ExpectAttackV2(0)
				dishonest.Attack2(values, 2).ExpectAttackV2(0)
				dishonest.Attack2(values, 3).ExpectAttackV2(0)

				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{0})
				dishonest.Attack2(values, 0).ExpectAttackV2(0)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{1})
				dishonest.Attack2(values, 0).ExpectAttackV2(1)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{2})
				dishonest.Attack2(values, 0).ExpectAttackV2(2)

				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 1), []uint64{0})
				dishonest.Attack2(values, 1).ExpectAttackV2(0)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 1), []uint64{1})
				dishonest.Attack2(values, 1).ExpectAttackV2(1)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 1), []uint64{2})
				dishonest.Attack2(values, 1).ExpectAttackV2(2)

				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 2), []uint64{0})
				dishonest.Attack2(values, 2).ExpectAttackV2(0)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 2), []uint64{1})
				dishonest.Attack2(values, 2).ExpectAttackV2(1)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 2), []uint64{2})
				dishonest.Attack2(values, 2).ExpectAttackV2(2)

				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 3), []uint64{0})
				dishonest.Attack2(values, 3).ExpectAttackV2(0)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 3), []uint64{1})
				dishonest.Attack2(values, 3).ExpectAttackV2(1)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 3), []uint64{2})
				dishonest.Attack2(values, 3).ExpectAttackV2(2)
			},
		},
		{
			name:             "HonestRoot-TwoLevelAttack-FirstLevelIncorrect-SecondLevelCorrect",
			rootClaimCorrect: true,
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				honest := builder.Seq()
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				honest.Attack2(values, 0).Attack2(nil, 0)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{1})
				honest.Attack2(values, 0).ExpectAttackV2(1).Attack2(nil, 0).ExpectAttackV2(3)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{2})
				honest.Attack2(values, 0).ExpectAttackV2(2).Attack2(nil, 0).ExpectAttackV2(3)
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{3})
				honest.Attack2(values, 0).ExpectAttackV2(3).Attack2(nil, 0).ExpectAttackV2(3)
			},
		},
		{
			name: "StepAtMaxDepth-LeftBranch-StepBranch0",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 0). // dishonest
					Attack2(nil, 0)     // honest
				lastHonestClaim.Attack2([]common.Hash{{0x1}, {0x2}, {0x3}}, 0).ExpectStepV2(0)
			},
		},
		{
			name: "StepAtMaxDepth-LeftBranch-StepBranch1",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 0). // dishonest
					Attack2(nil, 0)     // honest
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{1})
				lastHonestClaim.Attack2(values, 0).ExpectStepV2(1)
			},
		},
		{
			name: "StepAtMaxDepth-LeftBranch-StepBranch2",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 0). // dishonest
					Attack2(nil, 0)     // honest
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{2})
				lastHonestClaim.Attack2(values, 0).ExpectStepV2(2)
			},
		},
		{
			name: "StepAtMaxDepth-LeftBranch-StepBranch3",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 0). // dishonest
					Attack2(nil, 0)     // honest
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{})
				lastHonestClaim.Attack2(values, 0).ExpectStepV2(3)
			},
		},
		{
			name: "StepAtMaxDepth-MiddleBranch-StepBranch0",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 1). // dishonest
					Attack2(nil, 0)     // honest
				lastHonestClaim.Attack2([]common.Hash{{0x1}, {0x2}, {0x3}}, 0).ExpectStepV2(0)
			},
		},
		{
			name: "StepAtMaxDepth-MiddleBranch-StepBranch1",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 1). // dishonest
					Attack2(nil, 0)     // honest
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{1})
				lastHonestClaim.Attack2(values, 0).ExpectStepV2(1)
			},
		},
		{
			name: "StepAtMaxDepth-MiddleBranch-StepBranch2",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 1). // dishonest
					Attack2(nil, 0)     // honest
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{2})
				lastHonestClaim.Attack2(values, 0).ExpectStepV2(2)
			},
		},
		{
			name: "StepAtMaxDepth-MiddleBranch-StepBranch3",
			setupGame: func(builder *faulttest.GameBuilder) {
				values := []common.Hash{{0x81}, {0x82}, {0x83}}
				lastHonestClaim := builder.Seq().
					Attack2(nil, 0).    // honest
					Attack2(values, 1). // dishonest
					Attack2(nil, 0)     // honest
				lastPosition := builder.Game.Claims()[len(builder.Game.Claims())-1].Position
				values = claimBuilder.GetClaimsAtPosition(lastPosition.MoveN(nbits, 0), []uint64{})
				lastHonestClaim.Attack2(values, 0).ExpectStepV2(3)
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			builder := claimBuilder.GameBuilder(faulttest.WithInvalidValue(!test.rootClaimCorrect))
			test.setupGame(builder)
			game := builder.Game
			accessor := trace.NewSimpleTraceAccessor(claimBuilder.CorrectTraceProvider())
			solver := NewGameSolver(maxDepth, accessor, types.CallDataType)
			postState, actions := runStep2(t, solver, game, claimBuilder.CorrectTraceProvider())
			preimage := getStepPreimage(t, builder, game, accessor)
			for i, action := range builder.ExpectedActions {
				if action.Type == types.ActionTypeStep {
					action.OracleData = preimage
				}
				t.Logf("Expect %v: Type: %v, Position: %v, ParentIdx: %v, Branch: %v, Value: %v, SubValues: %v, PreState: %v, ProofData: %v",
					i, action.Type, action.ParentClaim.Position, action.ParentClaim.ContractIndex, action.AttackBranch, action.Value, action.SubValues, hex.EncodeToString(action.PreState), hex.EncodeToString(action.ProofData))
				require.Containsf(t, actions, action, "Expected claim %v missing", i)
			}
			require.Len(t, actions, len(builder.ExpectedActions), "Incorrect number of actions")
			verifyGameRules(t, postState, test.rootClaimCorrect)
		})
	}
}

func runStep2(t *testing.T, solver *GameSolver, game types.Game, correctTraceProvider types.TraceProvider) (types.Game, []types.Action) {
	actions, err := solver.CalculateNextActions(context.Background(), game)
	require.NoError(t, err)

	postState := applyActions2(game, challengerAddr, actions)

	for i, action := range actions {
		t.Logf("Real %v: Type: %v, Position: %v, ParentIdx: %v, Branch: %v, Value: %v, SubValues: %v, PreState: %v, ProofData: %v",
			i, action.Type,
			action.ParentClaim.Position,
			action.ParentClaim.ContractIndex,
			action.AttackBranch, action.Value, action.SubValues, hex.EncodeToString(action.PreState), hex.EncodeToString(action.ProofData))
		// Check that every move the solver returns meets the generic validation rules
		require.NoError(t, checkRules(game, action, correctTraceProvider), "Attempting to perform invalid action")
	}
	return postState, actions
}

func applyActions2(game types.Game, claimant common.Address, actions []types.Action) types.Game {
	claims := game.Claims()
	for _, action := range actions {
		switch action.Type {
		case types.ActionTypeAttackV2:
			newPosition := action.ParentClaim.Position.MoveN(game.NBits(), action.AttackBranch)
			claim := types.Claim{
				ClaimData: types.ClaimData{
					Value:    action.Value,
					Bond:     big.NewInt(0),
					Position: newPosition,
				},
				Claimant:            claimant,
				ContractIndex:       len(claims),
				ParentContractIndex: action.ParentClaim.ContractIndex,
				SubValues:           action.SubValues,
				AttackBranch:        action.AttackBranch,
			}
			claims = append(claims, claim)
		case types.ActionTypeStep:
			counteredClaim := claims[action.ParentClaim.ContractIndex]
			counteredClaim.CounteredBy = claimant
			claims[action.ParentClaim.ContractIndex] = counteredClaim
		default:
			panic(fmt.Errorf("unknown move type: %v", action.Type))
		}
	}
	return types.NewGameState2(claims, game.MaxDepth(), game.NBits(), game.SplitDepth())
}

func getStepPreimage(t *testing.T, builder *faulttest.GameBuilder, game types.Game, accessor types.TraceAccessor) *types.PreimageOracleData {
	if len(builder.ExpectedActions) == 0 {
		return nil
	}
	claims := game.Claims()
	lastClaim := claims[len(claims)-1]
	lastExpectedAction := builder.ExpectedActions[len(builder.ExpectedActions)-1]
	if lastExpectedAction.Type != types.ActionTypeStep {
		return nil
	}
	_, _, preimage, err := accessor.GetStepData2(context.Background(), game, lastClaim, lastClaim.MoveRightN(lastExpectedAction.AttackBranch))
	require.NoError(t, err)
	return preimage
}

func TestMultipleRoundsWithNbits2(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		actor actor
	}{
		{
			name:  "SingleRoot",
			actor: doNothingActor,
		},
		{
			name:  "LinearAttackCorrect",
			actor: correctAttackLastClaim,
		},
		{
			name:  "LinearDefendCorrect",
			actor: correctDefendLastClaim,
		},
		{
			name:  "LinearAttackIncorrect",
			actor: incorrectAttackLastClaim,
		},
		{
			name:  "LinearDefendInorrect",
			actor: incorrectDefendLastClaim,
		},
		{
			name:  "LinearDefendIncorrectDefendCorrect",
			actor: combineActors(incorrectDefendLastClaim, correctDefendLastClaim),
		},
		{
			name:  "LinearAttackIncorrectDefendCorrect",
			actor: combineActors(incorrectAttackLastClaim, correctDefendLastClaim),
		},
		{
			name:  "LinearDefendIncorrectDefendIncorrect",
			actor: combineActors(incorrectDefendLastClaim, incorrectDefendLastClaim),
		},
		{
			name:  "LinearAttackIncorrectDefendIncorrect",
			actor: combineActors(incorrectAttackLastClaim, incorrectDefendLastClaim),
		},
		{
			name:  "AttackEverythingCorrect",
			actor: attackEverythingCorrect,
		},
		{
			name:  "DefendEverythingCorrect",
			actor: defendEverythingCorrect,
		},
		{
			name:  "AttackEverythingIncorrect",
			actor: attackEverythingIncorrect,
		},
		{
			name:  "DefendEverythingIncorrect",
			actor: defendEverythingIncorrect,
		},
		{
			name:  "Exhaustive",
			actor: exhaustive,
		},
	}
	for _, test := range tests {
		test := test
		for _, rootClaimCorrect := range []bool{false} {
			rootClaimCorrect := rootClaimCorrect
			t.Run(fmt.Sprintf("%v-%v", test.name, rootClaimCorrect), func(t *testing.T) {
				t.Parallel()

				maxDepth := types.Depth(10)
				startingL2BlockNumber := big.NewInt(50)
				nbits := uint64(2)
				// splitDepth can't be 4 when maxDepth=8, because we can only attack branch 0 at claim with depth of splitDepth+nbits
				splitDepth := types.Depth(4)
				claimBuilder := faulttest.NewAlphabetClaimBuilder2(t, startingL2BlockNumber, maxDepth, nbits, splitDepth)
				builder := claimBuilder.GameBuilder(faulttest.WithInvalidValue(!rootClaimCorrect))
				game := builder.Game

				correctTrace := claimBuilder.CorrectTraceProvider()
				solver := NewGameSolver(maxDepth, trace.NewSimpleTraceAccessor(correctTrace), types.CallDataType)

				roundNum := 0
				done := false
				for !done {
					t.Logf("------ ROUND %v ------", roundNum)
					game, _ = runStep2(t, solver, game, correctTrace)
					verifyGameRules(t, game, rootClaimCorrect)

					game, done = test.actor.Apply(t, game, correctTrace)
					roundNum++
				}
			})
		}
	}
}
