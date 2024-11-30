package solver

import (
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/test"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
)

type actor interface {
	Apply(t *testing.T, game types.Game, correctTrace types.TraceProvider) (types.Game, bool)
}

type actorFn func(t *testing.T, game types.Game, correctTrace types.TraceProvider) (types.Game, bool)

func (a actorFn) Apply(t *testing.T, game types.Game, correctTrace types.TraceProvider) (types.Game, bool) {
	return a(t, game, correctTrace)
}

type builderFn func(builder *test.GameBuilder) bool

func (a builderFn) Apply(t *testing.T, game types.Game, correctTrace types.TraceProvider) (types.Game, bool) {
	builder := test.NewGameBuilderFromGame(t, correctTrace, game)
	done := a(builder)
	return builder.Game, done
}

func combineActors(actors ...actor) actor {
	return actorFn(func(t *testing.T, game types.Game, correctTrace types.TraceProvider) (types.Game, bool) {
		done := true
		for _, actor := range actors {
			newGame, actorDone := actor.Apply(t, game, correctTrace)
			game = newGame
			done = done && actorDone
		}
		return game, done
	})
}

var doNothingActor builderFn = func(builder *test.GameBuilder) bool {
	return true
}

var correctAttackLastClaim = respondLastClaim(func(seq *test.GameBuilderSeq) {
	seq.Attack2(nil, 0)
})

var correctDefendLastClaim = respondLastClaim(func(seq *test.GameBuilderSeq) {
	if seq.IsRoot() {
		// Must attack the root
		seq.Attack2(nil, 0)
	} else {
		seq.Attack2(nil, seq.MaxAttackBranch())
	}
})

var incorrectAttackLastClaim = respondLastClaim(func(seq *test.GameBuilderSeq) {
	incorrectSubValues := []common.Hash{}
	for i := uint64(0); i < seq.MaxAttackBranch(); i++ {
		incorrectSubValues = append(incorrectSubValues, common.Hash{0xaa})
	}
	seq.Attack2(incorrectSubValues, 0)
})

var incorrectDefendLastClaim = respondLastClaim(func(seq *test.GameBuilderSeq) {
	incorrectSubValues := []common.Hash{}
	for i := uint64(0); i < seq.MaxAttackBranch(); i++ {
		incorrectSubValues = append(incorrectSubValues, common.Hash{0xdd})
	}
	if seq.IsRoot() {
		// Must attack the root
		seq.Attack2(incorrectSubValues, 0)
	} else {
		seq.Attack2(incorrectSubValues, seq.MaxAttackBranch())
	}
})

var attackEverythingCorrect = respondAllClaims(func(seq *test.GameBuilderSeq) {
	seq.Attack2(nil, 0)
})

var defendEverythingCorrect = respondAllClaims(func(seq *test.GameBuilderSeq) {
	if seq.IsRoot() {
		// Must attack root
		seq.Attack2(nil, 0)
	} else {
		seq.Attack2(nil, seq.MaxAttackBranch())
	}
})

var attackEverythingIncorrect = respondAllClaims(func(seq *test.GameBuilderSeq) {
	incorrectSubValues := []common.Hash{}
	for i := uint64(0); i < seq.MaxAttackBranch(); i++ {
		incorrectSubValues = append(incorrectSubValues, common.Hash{0xaa})
	}
	seq.Attack2(incorrectSubValues, 0)
})

var defendEverythingIncorrect = respondAllClaims(func(seq *test.GameBuilderSeq) {
	incorrectSubValues := []common.Hash{}
	for i := uint64(0); i < seq.MaxAttackBranch(); i++ {
		incorrectSubValues = append(incorrectSubValues, common.Hash{0xbb})
	}
	if seq.IsRoot() {
		// Must attack root
		seq.Attack2(incorrectSubValues, 0)
	} else {
		seq.Attack2(incorrectSubValues, seq.MaxAttackBranch())
	}
})

var exhaustive = respondAllClaims(func(seq *test.GameBuilderSeq) {
	seq.Attack2(nil, 0)
	incorrectSubValues := []common.Hash{}
	for i := uint64(0); i < seq.MaxAttackBranch(); i++ {
		incorrectSubValues = append(incorrectSubValues, common.Hash{0xaa})
		if seq.IsSplitDepth() {
			// at splitDepth, there is only one subValue
			break
		}
	}
	seq.Attack2(incorrectSubValues, 0)
	if !seq.IsRoot() && !seq.IsTraceRoot() {
		seq.Attack2(nil, seq.MaxAttackBranch())
		for i := uint64(0); i < seq.MaxAttackBranch(); i++ {
			incorrectSubValues[i] = common.Hash{0xdd}
			if seq.IsSplitDepth() {
				break
			}
		}
		seq.Attack2(incorrectSubValues, seq.MaxAttackBranch())
	}
})

func respondLastClaim(respond func(seq *test.GameBuilderSeq)) builderFn {
	return func(builder *test.GameBuilder) bool {
		seq := seqFromLastClaim(builder)
		if seq.IsMaxDepth() {
			// Can't counter the leaf claim
			return true
		}
		respond(seq)
		return false
	}
}

func respondAllClaims(respond func(seq *test.GameBuilderSeq)) builderFn {
	return func(builder *test.GameBuilder) bool {
		startingCount := len(builder.Game.Claims())
		for _, claim := range builder.Game.Claims() {
			if claim.Depth() == builder.Game.MaxDepth() {
				continue
			}
			respond(builder.SeqFrom(claim))
		}
		finalCount := len(builder.Game.Claims())
		return finalCount == startingCount
	}
}

func seqFromLastClaim(builder *test.GameBuilder) *test.GameBuilderSeq {
	claims := builder.Game.Claims()
	claim := claims[len(claims)-1]
	return builder.SeqFrom(claim)
}
