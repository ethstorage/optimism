package test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestGetClaimsHash(t *testing.T) {
	claims := []common.Hash{{0x01}, {0x02}, {0x03}}
	rootHash := GetClaimsHash(claims)
	expectedHash := crypto.Keccak256Hash(claims[0][:], claims[1][:])
	expectedHash = crypto.Keccak256Hash(expectedHash[:], claims[2][:])
	require.Equal(t, expectedHash, rootHash)
}
