//go:build tinygo
// +build tinygo

package client

import preimage "github.com/ethereum-optimism/optimism/op-preimage"

func NewOracleClientAndHintWriter() (preimage.Oracle, preimage.Hinter) {
	o := wasmHostIO{}

	return o, o
}

type wasmHostIO struct {
}

func (o wasmHostIO) Get(key preimage.Key) []byte {
	return getKeyFromOracle(key.PreimageKey())
}

func (o wasmHostIO) Hint(hint preimage.Hint) {
	hintToHinter(hint.Hint())
}

//export getKeyFromOracle
func getKeyFromOracle([32]byte) []byte

//export hintToHinter
func hintToHinter(string)
