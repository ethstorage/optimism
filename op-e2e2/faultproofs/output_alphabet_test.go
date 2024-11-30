package faultproofs

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger2/config"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/types"
	op_e2e "github.com/ethereum-optimism/optimism/op-e2e2"
	"github.com/ethereum-optimism/optimism/op-e2e2/e2eutils/challenger"
	"github.com/ethereum-optimism/optimism/op-e2e2/e2eutils/disputegame"
	"github.com/ethereum-optimism/optimism/op-e2e2/e2eutils/wait"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

const (
	daType = config.DACalldata
)

func TestOutputAlphabetGame_ChallengerWins(t *testing.T) {
	op_e2e.InitParallel(t)
	ctx := context.Background()
	sys, l1Client := StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)

	disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys, daType)
	game := disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 3, common.Hash{0xff})
	correctTrace := game.CreateHonestActor(ctx, "sequencer")
	game.LogGameData(ctx)

	opts := challenger.WithPrivKey(sys.Cfg.Secrets.Alice)
	game.StartChallenger(ctx, "sequencer", "Challenger", opts, challenger.WithDaType(daType))
	game.LogGameData(ctx)

	// Challenger should post an output root to counter claims down to the leaf level of the top game
	claim := game.RootClaim(ctx)
	incorrectValues := setIncorrectValues(common.Hash{0xaa}, 1<<game.Nbits-1)
	for claim.IsOutputRoot(ctx) && !claim.IsOutputRootLeaf(ctx) {
		if claim.AgreesWithOutputRoot() {
			// If the latest claim agrees with the output root, expect the honest challenger to counter it
			claim = claim.WaitForCounterClaim(ctx)
			game.LogGameData(ctx)
			claim.RequireCorrectOutputRoot(ctx)
		} else {
			// Otherwise we should counter
			claim = claim.Attack2(ctx, 0, incorrectValues)
			game.LogGameData(ctx)
		}
	}

	// Wait for the challenger to post the first claim in the cannon trace
	claim = claim.WaitForCounterClaim(ctx)
	game.LogGameData(ctx)
	// Attack the root of the alphabet trace subgame
	claim = correctTrace.AttackClaim(ctx, claim, 0)
	for !claim.IsMaxDepth(ctx) {
		if claim.AgreesWithOutputRoot() {
			// If the latest claim supports the output root, wait for the honest challenger to respond
			claim = claim.WaitForCounterClaim(ctx)
			game.LogGameData(ctx)
		} else {
			// Otherwise we need to counter the honest claim
			claim = correctTrace.AttackClaim(ctx, claim, 0)
			game.LogGameData(ctx)
		}
	}
	// Challenger should be able to call step and counter the leaf claim.
	claim.WaitForCountered(ctx)
	game.LogGameData(ctx)

	sys.TimeTravelClock.AdvanceTime(game.MaxClockDuration(ctx))
	require.NoError(t, wait.ForNextBlock(ctx, l1Client))
	game.WaitForGameStatus(ctx, types.GameStatusChallengerWon)
	game.LogGameData(ctx)
}

func TestOutputAlphabetGame_ReclaimBond(t *testing.T) {
	op_e2e.InitParallel(t)
	ctx := context.Background()
	sys, l1Client := StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)

	disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys, daType)
	game := disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 3, common.Hash{0xff})
	game.LogGameData(ctx)

	// The dispute game should have a zero balance
	balance := game.WethBalance(ctx, game.Addr)
	require.Zero(t, balance.Uint64())

	alice := sys.Cfg.Secrets.Addresses().Alice

	// Grab the root claim
	claim := game.RootClaim(ctx)
	opts := challenger.WithPrivKey(sys.Cfg.Secrets.Alice)
	game.StartChallenger(ctx, "sequencer", "Challenger", opts)
	game.LogGameData(ctx)

	// Perform a few moves
	claim = claim.WaitForCounterClaim(ctx)
	game.LogGameData(ctx)
	incorrectValues := setIncorrectValues(common.Hash{0xaa}, 1<<game.Nbits-1)
	claim = claim.Attack2(ctx, 0, incorrectValues)
	claim = claim.WaitForCounterClaim(ctx)
	game.LogGameData(ctx)
	claim = claim.Attack2(ctx, 0, incorrectValues)
	game.LogGameData(ctx)
	// when nbits=2, splitDepth is 4, maxDepth is 8, the following line should be commentted,
	// because the maxDepth is reached and we don't need to waitForCounterClaim.
	_ = claim.WaitForCounterClaim(ctx)

	// Expect posted claims so the game balance is non-zero
	balance = game.WethBalance(ctx, game.Addr)
	require.Truef(t, balance.Cmp(big.NewInt(0)) > 0, "Expected game balance to be above zero")

	sys.TimeTravelClock.AdvanceTime(game.MaxClockDuration(ctx))
	require.NoError(t, wait.ForNextBlock(ctx, l1Client))
	game.WaitForGameStatus(ctx, types.GameStatusChallengerWon)
	game.LogGameData(ctx)

	// Expect Alice's credit to be non-zero
	// But it can't be claimed right now since there is a delay on the weth unlock
	require.Truef(t, game.AvailableCredit(ctx, alice).Cmp(big.NewInt(0)) > 0, "Expected alice credit to be above zero")

	// The actor should have no credit available because all its bonds were paid to Alice.
	actorCredit := game.AvailableCredit(ctx, disputegame.TestAddress)
	require.True(t, actorCredit.Cmp(big.NewInt(0)) == 0, "Expected alice available credit to be zero")

	// Advance the time past the weth unlock delay
	sys.TimeTravelClock.AdvanceTime(game.CreditUnlockDuration(ctx))
	require.NoError(t, wait.ForNextBlock(ctx, l1Client))

	// Wait for alice to have no available credit
	// aka, wait for the challenger to claim its credit
	game.WaitForNoAvailableCredit(ctx, alice)

	// The dispute game delayed weth balance should be zero since it's all claimed
	require.True(t, game.WethBalance(ctx, game.Addr).Cmp(big.NewInt(0)) == 0)
}

func TestOutputAlphabetGame_ValidOutputRoot(t *testing.T) {
	op_e2e.InitParallel(t)
	ctx := context.Background()
	sys, l1Client := StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)

	disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys, daType)
	game := disputeGameFactory.StartOutputAlphabetGameWithCorrectRoot(ctx, "sequencer", 2)
	correctTrace := game.CreateHonestActor(ctx, "sequencer")
	game.LogGameData(ctx)
	claim := game.DisputeLastBlock(ctx)
	// Invalid root claim of the alphabet game
	claim = claim.Attack2(ctx, 0, setIncorrectValues(common.Hash{0x01}, 1<<game.Nbits-1))

	opts := challenger.WithPrivKey(sys.Cfg.Secrets.Alice)
	game.StartChallenger(ctx, "sequencer", "Challenger", opts)

	claim = claim.WaitForCounterClaim(ctx)
	game.LogGameData(ctx)
	for !claim.IsMaxDepth(ctx) {
		// Dishonest actor always attacks with the correct trace
		claim = correctTrace.AttackClaim(ctx, claim, 0)
		claim = claim.WaitForCounterClaim(ctx)
		game.LogGameData(ctx)
	}

	sys.TimeTravelClock.AdvanceTime(game.MaxClockDuration(ctx))
	require.NoError(t, wait.ForNextBlock(ctx, l1Client))
	game.WaitForGameStatus(ctx, types.GameStatusDefenderWon)
}

func TestChallengerCompleteExhaustiveDisputeGame(t *testing.T) {
	op_e2e.InitParallel(t)

	testCase := func(t *testing.T, isRootCorrect bool) {
		ctx := context.Background()
		sys, l1Client := StartFaultDisputeSystem(t)
		t.Cleanup(sys.Close)

		disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys, daType)
		var game *disputegame.OutputAlphabetGameHelper
		if isRootCorrect {
			game = disputeGameFactory.StartOutputAlphabetGameWithCorrectRoot(ctx, "sequencer", 1)
		} else {
			game = disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 1, common.Hash{0xaa, 0xbb, 0xcc})
		}
		claim := game.DisputeLastBlock(ctx)

		game.LogGameData(ctx)

		// Start honest challenger
		game.StartChallenger(ctx, "sequencer", "Challenger",
			challenger.WithAlphabet(),
			challenger.WithPrivKey(sys.Cfg.Secrets.Alice),
			// Ensures the challenger responds to all claims before test timeout
			challenger.WithPollInterval(time.Millisecond*400),
		)

		if isRootCorrect {
			// Attack the correct output root with an invalid alphabet trace
			claim = claim.Attack2(ctx, 0, setIncorrectValues(common.Hash{0x01}, 1<<game.Nbits-1))
		} else {
			// Wait for the challenger to counter the invalid output root
			claim = claim.WaitForCounterClaim(ctx)
		}

		// Start dishonest challenger
		dishonestHelper := game.CreateDishonestHelper(ctx, "sequencer", !isRootCorrect)
		dishonestHelper.ExhaustDishonestClaims(ctx, claim)

		// Wait until we've reached max depth before checking for inactivity
		game.WaitForClaimAtDepth(ctx, game.MaxDepth(ctx))

		// Wait for 4 blocks of no challenger responses. The challenger may still be stepping on invalid claims at max depth
		game.WaitForInactivity(ctx, 4, false)

		gameDuration := game.MaxClockDuration(ctx)
		sys.TimeTravelClock.AdvanceTime(gameDuration)
		require.NoError(t, wait.ForNextBlock(ctx, l1Client))

		expectedStatus := types.GameStatusChallengerWon
		if isRootCorrect {
			expectedStatus = types.GameStatusDefenderWon
		}
		game.WaitForGameStatus(ctx, expectedStatus)
		game.LogGameData(ctx)
	}

	t.Run("RootCorrect", func(t *testing.T) {
		op_e2e.InitParallel(t)
		testCase(t, true)
	})
	t.Run("RootIncorrect", func(t *testing.T) {
		op_e2e.InitParallel(t)
		testCase(t, false)
	})
}

func TestOutputAlphabetGame_FreeloaderEarnsNothing(t *testing.T) {
	op_e2e.InitParallel(t)
	ctx := context.Background()
	sys, l1Client := StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)

	freeloaderOpts, err := bind.NewKeyedTransactorWithChainID(sys.Cfg.Secrets.Mallory, sys.Cfg.L1ChainIDBig())
	require.Nil(t, err)

	disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys, daType)
	game := disputeGameFactory.StartOutputAlphabetGameWithCorrectRoot(ctx, "sequencer", 2)
	correctTrace := game.CreateHonestActor(ctx, "sequencer")
	game.LogGameData(ctx)
	claim := game.DisputeLastBlock(ctx)
	// Invalid root claim of the alphabet game
	maxAttackBranch := 1<<game.Nbits - 1
	claim = claim.Attack2(ctx, 0, setIncorrectValues(common.Hash{0x01}, maxAttackBranch))

	// Chronology of claims:
	// dishonest root claim:
	//   - honest counter
	//     - dishonest
	//       - freeloader
	//       - honest
	//       - freeloader
	// The freeloader must be positioned leftmost (gindex positioning) or at the same position as honest claims.

	// honest counter
	claim = correctTrace.AttackClaim(ctx, claim, 0)

	var freeloaders []*disputegame.ClaimHelper

	// dishonest
	dishonest := correctTrace.AttackClaim(ctx, claim, 0)

	maxAttakBranch := uint64(1<<game.Nbits - 1)
	freeloaders = append(freeloaders, correctTrace.AttackClaim(ctx, dishonest, 0, disputegame.WithTransactOpts(freeloaderOpts)))
	freeloaders = append(freeloaders, dishonest.Attack2(ctx, 0, setIncorrectValues(common.Hash{0x02}, maxAttackBranch), disputegame.WithTransactOpts(freeloaderOpts)))
	freeloaders = append(freeloaders, dishonest.Attack2(ctx, maxAttakBranch, setIncorrectValues(common.Hash{0x03}, maxAttackBranch), disputegame.WithTransactOpts(freeloaderOpts)))

	// Ensure freeloaders respond before the honest challenger
	game.StartChallenger(ctx, "sequencer", "Challenger", challenger.WithPrivKey(sys.Cfg.Secrets.Alice))
	dishonest.WaitForCounterClaim(ctx, freeloaders...)

	// Freeloaders after the honest challenger
	freeloaders = append(freeloaders, dishonest.Attack2(ctx, 0, setIncorrectValues(common.Hash{0x04}, maxAttackBranch), disputegame.WithTransactOpts(freeloaderOpts)))
	freeloaders = append(freeloaders, dishonest.Attack2(ctx, maxAttakBranch, setIncorrectValues(common.Hash{0x05}, maxAttackBranch), disputegame.WithTransactOpts(freeloaderOpts)))

	for _, freeloader := range freeloaders {
		if freeloader.IsMaxDepth(ctx) {
			freeloader.WaitForCountered(ctx)
		} else {
			freeloader.WaitForCounterClaim(ctx)
		}
	}

	game.LogGameData(ctx)
	sys.TimeTravelClock.AdvanceTime(game.MaxClockDuration(ctx))
	require.NoError(t, wait.ForNextBlock(ctx, l1Client))
	game.WaitForGameStatus(ctx, types.GameStatusDefenderWon)

	game.LogGameData(ctx)
	amt := game.Credit(ctx, freeloaderOpts.From)
	require.Truef(t, amt.BitLen() == 0, "freeloaders should not be rewarded. Credit: %v", amt)
}

func TestHighestActedL1BlockMetric(t *testing.T) {
	op_e2e.InitParallel(t)
	ctx := context.Background()
	sys, l1Client := StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)

	disputeGameFactory := disputegame.NewFactoryHelper(t, ctx, sys, daType)
	honestChallenger := disputeGameFactory.StartChallenger(ctx, "Honest", challenger.WithAlphabet(), challenger.WithPrivKey(sys.Cfg.Secrets.Alice))

	game1 := disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 1, common.Hash{0xaa})
	sys.AdvanceTime(game1.MaxClockDuration(ctx))
	require.NoError(t, wait.ForNextBlock(ctx, l1Client))

	game1.WaitForGameStatus(ctx, types.GameStatusDefenderWon)

	disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 2, common.Hash{0xaa})
	disputeGameFactory.StartOutputAlphabetGame(ctx, "sequencer", 3, common.Hash{0xaa})

	honestChallenger.WaitL1HeadActedOn(ctx, l1Client)

	require.NoError(t, wait.ForNextBlock(ctx, l1Client))
	honestChallenger.WaitL1HeadActedOn(ctx, l1Client)
}

func setIncorrectValues(v common.Hash, nelements int) (subValues []common.Hash) {
	for i := 0; i < nelements; i++ {
		subValues = append(subValues, v)
	}
	return subValues
}
