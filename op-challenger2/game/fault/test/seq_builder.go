package test

import (
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
)

type SequenceBuilder struct {
	builder   *ClaimBuilder
	lastClaim types.Claim
}

// Seq starts building a claim by following a sequence of attack and defend moves from the root
// The returned SequenceBuilder can be used to add additional moves. e.g:
// claim := Seq(true).Attack(false).Attack(true).Defend(true).Get()
func (c *ClaimBuilder) Seq(rootOpts ...ClaimOpt) *SequenceBuilder {
	claim := c.CreateRootClaim(rootOpts...)
	return &SequenceBuilder{
		builder:   c,
		lastClaim: claim,
	}
}

func (s *SequenceBuilder) Attack2(subValues []common.Hash, branch uint64, opts ...ClaimOpt) *SequenceBuilder {
	claim := s.builder.AttackClaim2(s.lastClaim, subValues, branch, opts...)
	return &SequenceBuilder{
		builder:   s.builder,
		lastClaim: claim,
	}
}

func (s *SequenceBuilder) Defend(opts ...ClaimOpt) *SequenceBuilder {
	claim := s.builder.DefendClaim(s.lastClaim, opts...)
	return &SequenceBuilder{
		builder:   s.builder,
		lastClaim: claim,
	}
}

func (s *SequenceBuilder) Get() types.Claim {
	return s.lastClaim
}
