//go:build faultdisputegamen
// +build faultdisputegamen

package types

import "github.com/ethereum/go-ethereum/common"

type ActionType string

func (a ActionType) String() string {
	return string(a)
}

const (
	ActionTypeMove                   ActionType = "move"
	ActionTypeAttack                 ActionType = "attack"
	ActionTypeStep                   ActionType = "step"
	ActionTypeChallengeL2BlockNumber ActionType = "challenge-l2-block-number"
)

type Action struct {
	Type ActionType

	// Moves and Steps
	ParentClaim Claim
	IsAttack    bool

	// Moves
	Value     common.Hash
	SubValues []common.Hash

	// Steps
	PreState   []byte
	ProofData  []byte
	OracleData *PreimageOracleData

	// Challenge L2 Block Number
	InvalidL2BlockNumberChallenge *InvalidL2BlockNumberChallenge
}
