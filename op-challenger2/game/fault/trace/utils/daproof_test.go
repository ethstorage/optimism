package utils

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestUnpaddingMerkelProof(t *testing.T) {
	oneleave := []common.Hash{{0xaa}}
	twoLeaves := []common.Hash{{0xaa}, {0xbb}}
	fourLeaves := []common.Hash{{0xaa}, {0xbb}, {0xcc}, {0xdd}}
	fiveLeaves := []common.Hash{{0xaa}, {0xbb}, {0xcc}, {0xdd}, {0xee}}
	tests := []struct {
		data        []common.Hash
		index       uint32
		expectProof []byte
	}{
		// 1 leave
		{
			data:        oneleave,
			index:       0,
			expectProof: nil,
		},
		// 2 leaves
		{
			data:        twoLeaves,
			index:       0,
			expectProof: twoLeaves[1][:],
		},
		{
			data:        twoLeaves,
			index:       1,
			expectProof: twoLeaves[0][:],
		},
		// 4 leaves
		{
			data:        fourLeaves,
			index:       0,
			expectProof: append(fourLeaves[1][:], crypto.Keccak256(fourLeaves[2][:], fourLeaves[3][:])...),
		},
		{
			data:        fourLeaves,
			index:       3,
			expectProof: append(fourLeaves[2][:], crypto.Keccak256(fourLeaves[0][:], fourLeaves[1][:])...),
		},
		// 5 leaves
		{
			data:        fiveLeaves,
			index:       0,
			expectProof: append(append(fiveLeaves[1][:], crypto.Keccak256(fiveLeaves[2][:], fiveLeaves[3][:])...), fiveLeaves[4][:]...),
		},
		{
			data:        fiveLeaves,
			index:       1,
			expectProof: append(append(fiveLeaves[0][:], crypto.Keccak256(fiveLeaves[2][:], fiveLeaves[3][:])...), fiveLeaves[4][:]...),
		},
		{
			data:        fiveLeaves,
			index:       2,
			expectProof: append(append(fiveLeaves[3][:], crypto.Keccak256(fiveLeaves[0][:], fiveLeaves[1][:])...), fiveLeaves[4][:]...),
		},
		{
			data:        fiveLeaves,
			index:       3,
			expectProof: append(append(fiveLeaves[2][:], crypto.Keccak256(fiveLeaves[0][:], fiveLeaves[1][:])...), fiveLeaves[4][:]...),
		},
		{
			data:        fiveLeaves,
			index:       4,
			expectProof: crypto.Keccak256(crypto.Keccak256(fiveLeaves[0][:], fiveLeaves[1][:]), crypto.Keccak256(fiveLeaves[2][:], fiveLeaves[3][:])),
		},
	}

	for _, tCase := range tests {
		// Build the Merkle tree without padding
		tree := NewMerkleTree(tCase.data)
		proof := tree.GenerateProof(tCase.index)
		require.Equal(t, tCase.expectProof, proof)
	}
}
