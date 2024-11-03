package types

import "github.com/ethereum/go-ethereum/common"

type ActionType string

func (a ActionType) String() string {
	return string(a)
}

const (
	ActionTypeAttackV2               ActionType = "attackV2"
	ActionTypeMove                   ActionType = "move"
	ActionTypeStep                   ActionType = "step"
	ActionTypeChallengeL2BlockNumber ActionType = "challenge-l2-block-number"
)

type Action struct {
	Type ActionType

	// Moves and Steps
	ParentClaim Claim
	IsAttack    bool

	// Moves
	Value        common.Hash
	SubValues    *[]common.Hash
	AttackBranch uint64
	DAType       DAType

	// Steps
	PreState   []byte
	ProofData  []byte
	OracleData *PreimageOracleData

	// Challenge L2 Block Number
	InvalidL2BlockNumberChallenge *InvalidL2BlockNumberChallenge
}
