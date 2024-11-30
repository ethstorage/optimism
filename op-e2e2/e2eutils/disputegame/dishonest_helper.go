package disputegame

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
)

type DishonestHelper struct {
	*OutputGameHelper
	*OutputHonestHelper
	defender bool
}

func newDishonestHelper(g *OutputGameHelper, correctTrace *OutputHonestHelper, defender bool) *DishonestHelper {
	return &DishonestHelper{g, correctTrace, defender}
}

// ExhaustDishonestClaims makes all possible significant moves (mod honest challenger's) in a game.
// It is very inefficient and should NOT be used on games with large depths
func (d *DishonestHelper) ExhaustDishonestClaims(ctx context.Context, rootClaim *ClaimHelper) {
	depth := d.MaxDepth(ctx)
	splitDepth := d.SplitDepth(ctx)

	move := func(claimIndex int64, claimData types.Claim) {
		// dishonest level, valid attack
		// dishonest level, invalid attack
		// dishonest level, valid defense
		// dishonest level, invalid defense
		// honest level, invalid attack
		// honest level, invalid defense

		if claimData.Depth() == depth {
			return
		}

		d.LogGameData(ctx)
		d.OutputGameHelper.T.Logf("Dishonest moves against claimIndex %d", claimIndex)
		agreeWithLevel := d.defender == (claimData.Depth()%(types.Depth(2*d.Nbits)) == 0)
		maxAttackBranch := 1<<d.Nbits - 1
		if !agreeWithLevel {
			d.OutputHonestHelper.Attack2(ctx, claimIndex, 0, WithIgnoreDuplicates())
			if claimIndex != 0 && claimData.Depth() != splitDepth+types.Depth(d.Nbits) {
				d.OutputHonestHelper.Attack2(ctx, claimIndex, uint64(maxAttackBranch), WithIgnoreDuplicates())
			}
		}
		incorresctSubValues := []common.Hash{}
		for i := 0; i < maxAttackBranch; i++ {
			incorresctSubValues = append(incorresctSubValues, common.Hash{byte(claimIndex)})
		}
		d.OutputGameHelper.Attack2(ctx, claimIndex, 0, incorresctSubValues, WithIgnoreDuplicates())
		if claimIndex != 0 && claimData.Depth() != splitDepth+types.Depth(d.Nbits) {
			d.OutputGameHelper.Attack2(ctx, claimIndex, uint64(maxAttackBranch), incorresctSubValues, WithIgnoreDuplicates())
		}
	}

	numClaimsSeen := rootClaim.Index
	for {
		// Use a short timeout since we don't know the challenger will respond,
		// and this is only designed for the alphabet game where the response should be fast.
		newCount, err := d.waitForNewClaim(ctx, numClaimsSeen, 30*time.Second)
		if errors.Is(err, context.DeadlineExceeded) {
			// we assume that the honest challenger has stopped responding
			// There's nothing to respond to.
			break
		}
		d.OutputGameHelper.Require.NoError(err)

		for ; numClaimsSeen < newCount; numClaimsSeen++ {
			claimData := d.getClaim(ctx, numClaimsSeen)
			move(numClaimsSeen, claimData)
		}
	}
}
