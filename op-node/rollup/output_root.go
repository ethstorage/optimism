package rollup

import (
	"errors"

	"github.com/ethereum-optimism/optimism/op-node/eth"
	"github.com/ethereum/go-ethereum/crypto"
)

var NilProof = errors.New("Output root proof is nil")

// TypesOutputRootProof is an auto generated low-level Go binding around an user-defined struct.
type TypesOutputRootProof struct {
	Version                  [32]byte
	StateRoot                [32]byte
	MessagePasserStorageRoot [32]byte
	LatestBlockhash          [32]byte
}

// ComputeL2OutputRoot computes the L2 output root by hashing an output root proof.
func ComputeL2OutputRoot(proofElements *TypesOutputRootProof) (eth.Bytes32, error) {
	if proofElements == nil {
		return eth.Bytes32{}, NilProof
	}

	digest := crypto.Keccak256Hash(
		proofElements.Version[:],
		proofElements.StateRoot[:],
		proofElements.MessagePasserStorageRoot[:],
		proofElements.LatestBlockhash[:],
	)
	return eth.Bytes32(digest), nil
}

func ComputeL2OutputRootV0(block eth.BlockInfo, storageRoot [32]byte) (eth.Bytes32, error) {
	var l2OutputRootVersion eth.Bytes32 // it's zero for now
	return ComputeL2OutputRoot(&TypesOutputRootProof{
		Version:                  l2OutputRootVersion,
		StateRoot:                block.Root(),
		MessagePasserStorageRoot: storageRoot,
		LatestBlockhash:          block.Hash(),
	})
}
