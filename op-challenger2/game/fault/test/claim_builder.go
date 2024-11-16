package test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/contracts"
	"github.com/ethereum-optimism/optimism/op-challenger2/game/fault/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var DefaultClaimant = common.Address{0xba, 0xdb, 0xad, 0xba, 0xdb, 0xad}

type claimCfg struct {
	value          common.Hash
	invalidValue   bool
	claimant       common.Address
	parentIdx      int
	clockTimestamp time.Time
	clockDuration  time.Duration
	branch         uint64
	subValues      *[]common.Hash
}

func newClaimCfg(opts ...ClaimOpt) *claimCfg {
	cfg := &claimCfg{
		clockTimestamp: time.Unix(math.MaxInt64-1, 0),
	}
	for _, opt := range opts {
		opt.Apply(cfg)
	}
	return cfg
}

type ClaimOpt interface {
	Apply(cfg *claimCfg)
}

type claimOptFn func(cfg *claimCfg)

func (c claimOptFn) Apply(cfg *claimCfg) {
	c(cfg)
}

func WithValue(value common.Hash) ClaimOpt {
	return claimOptFn(func(cfg *claimCfg) {
		cfg.value = value
	})
}

func WithInvalidValue(invalid bool) ClaimOpt {
	return claimOptFn(func(cfg *claimCfg) {
		cfg.invalidValue = invalid
	})
}

func WithClaimant(claimant common.Address) ClaimOpt {
	return claimOptFn(func(cfg *claimCfg) {
		cfg.claimant = claimant
	})
}

func WithParent(claim types.Claim) ClaimOpt {
	return claimOptFn(func(cfg *claimCfg) {
		cfg.parentIdx = claim.ContractIndex
	})
}

func WithClock(timestamp time.Time, duration time.Duration) ClaimOpt {
	return claimOptFn(func(cfg *claimCfg) {
		cfg.clockTimestamp = timestamp
		cfg.clockDuration = duration
	})
}

func WithBranch(branch uint64) ClaimOpt {
	return claimOptFn(func(cfg *claimCfg) {
		cfg.branch = branch
	})
}

func WithSubValues(subValues *[]common.Hash) ClaimOpt {
	return claimOptFn(func(cfg *claimCfg) {
		cfg.subValues = subValues
	})
}

// ClaimBuilder is a test utility to enable creating claims in a wide range of situations
type ClaimBuilder struct {
	require    *require.Assertions
	maxDepth   types.Depth
	nbits      uint64
	splitDepth types.Depth
	correct    types.TraceProvider
}

// NewClaimBuilder creates a new [ClaimBuilder].
func NewClaimBuilder2(t *testing.T, maxDepth types.Depth, nbits uint64, splitDepth types.Depth, provider types.TraceProvider) *ClaimBuilder {
	return &ClaimBuilder{
		require:    require.New(t),
		maxDepth:   maxDepth,
		nbits:      nbits,
		splitDepth: splitDepth,
		correct:    provider,
	}
}

// CorrectTraceProvider returns a types.TraceProvider that provides the canonical trace.
func (c *ClaimBuilder) CorrectTraceProvider() types.TraceProvider {
	return c.correct
}

// CorrectClaimAtPosition returns the canonical claim at a specified position
func (c *ClaimBuilder) CorrectClaimAtPosition(pos types.Position) common.Hash {
	value, err := c.correct.Get(context.Background(), pos)
	c.require.NoError(err)
	return value
}

// CorrectPreState returns the pre-state (not hashed) required to execute the valid step at the specified trace index
func (c *ClaimBuilder) CorrectPreState(idx *big.Int) []byte {
	pos := types.NewPosition(c.maxDepth, idx)
	preimage, _, _, err := c.correct.GetStepData(context.Background(), pos)
	c.require.NoError(err)
	return preimage
}

// CorrectProofData returns the proof-data required to execute the valid step at the specified trace index
func (c *ClaimBuilder) CorrectProofData(idx *big.Int) []byte {
	pos := types.NewPosition(c.maxDepth, idx)
	_, proof, _, err := c.correct.GetStepData(context.Background(), pos)
	c.require.NoError(err)
	return proof
}

func (c *ClaimBuilder) CorrectOracleData(idx *big.Int) *types.PreimageOracleData {
	pos := types.NewPosition(c.maxDepth, idx)
	_, _, data, err := c.correct.GetStepData(context.Background(), pos)
	c.require.NoError(err)
	return data
}

func (c *ClaimBuilder) incorrectClaim(pos types.Position) common.Hash {
	return common.BigToHash(pos.TraceIndex(c.maxDepth))
}

func (c *ClaimBuilder) claim(pos types.Position, opts ...ClaimOpt) types.Claim {
	cfg := newClaimCfg(opts...)
	claim := types.Claim{
		ClaimData: types.ClaimData{
			Position: pos,
		},
		Claimant: DefaultClaimant,
		Clock: types.Clock{
			Duration:  cfg.clockDuration,
			Timestamp: cfg.clockTimestamp,
		},
		AttackBranch: cfg.branch,
	}
	if cfg.subValues != nil {
		claim.SubValues = cfg.subValues
		claim.Value = contracts.SubValuesHash(*claim.SubValues)
	} else {
		values := []common.Hash{}
		if pos.ToGIndex().Cmp(big.NewInt(1)) == 0 {
			if cfg.invalidValue {
				values = append(values, c.incorrectClaim(pos))
			} else {
				values = append(values, c.CorrectClaimAtPosition(pos))
			}
			for i := uint64(0); i < c.MaxAttackBranch()-1; i++ {
				values = append(values, common.Hash{})
			}
		} else {
			for i := uint64(0); i < c.MaxAttackBranch(); i++ {
				pos := pos.MoveRightN(i)
				if cfg.invalidValue {
					values = append(values, c.incorrectClaim(pos))
				} else {
					values = append(values, c.CorrectClaimAtPosition(pos))
				}
			}
		}
		claim.SubValues = &values
		claim.Value = contracts.SubValuesHash(*claim.SubValues)
	}
	claim.ParentContractIndex = cfg.parentIdx
	fmt.Printf("claim.Value: %v\n", claim.Value)
	return claim
}

func (c *ClaimBuilder) CreateRootClaim(opts ...ClaimOpt) types.Claim {
	pos := types.NewPositionFromGIndex(big.NewInt(1))
	return c.claim(pos, opts...)
}

func (c *ClaimBuilder) CreateLeafClaim(traceIndex *big.Int, opts ...ClaimOpt) types.Claim {
	pos := types.NewPosition(c.maxDepth, traceIndex)
	return c.claim(pos, opts...)
}

func (c *ClaimBuilder) AttackClaim2(claim types.Claim, subValues []common.Hash, branch uint64, opts ...ClaimOpt) types.Claim {
	pos := claim.Position.MoveN(c.nbits, branch)
	if subValues == nil {
		// aggCLaim will be auto generated in this function
		return c.claim(pos, append([]ClaimOpt{WithParent(claim), WithBranch(branch)}, opts...)...)
	} else {
		aggClaim := contracts.SubValuesHash(subValues)
		return c.claim(pos, append([]ClaimOpt{WithParent(claim), WithValue(aggClaim), WithBranch(branch), WithSubValues(&subValues)}, opts...)...)
	}
}

func (c *ClaimBuilder) DefendClaim(claim types.Claim, opts ...ClaimOpt) types.Claim {
	pos := claim.Position.Defend()
	return c.claim(pos, append([]ClaimOpt{WithParent(claim)}, opts...)...)
}

func (c *ClaimBuilder) NBits() uint64 {
	return c.nbits
}

func (c *ClaimBuilder) MaxAttackBranch() uint64 {
	return 1<<c.nbits - 1
}

func (c *ClaimBuilder) SplitDepth() types.Depth {
	return c.splitDepth
}

func (c *ClaimBuilder) TraceRootDepth() types.Depth {
	return c.splitDepth + types.Depth(c.nbits)
}
