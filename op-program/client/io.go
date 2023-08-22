//go:build !tinygo
// +build !tinygo

package client

import (
	"os"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	oppio "github.com/ethereum-optimism/optimism/op-program/io"
)

func NewOracleClientAndHintWriter() (preimage.Oracle, preimage.Hinter) {
	preimageOracle := CreatePreimageChannel()
	preimageHinter := CreateHinterChannel()

	return preimage.NewOracleClient(preimageOracle), preimage.NewHintWriter(preimageHinter)
}

func CreateHinterChannel() oppio.FileChannel {
	r := os.NewFile(HClientRFd, "preimage-hint-read")
	w := os.NewFile(HClientWFd, "preimage-hint-write")
	return oppio.NewReadWritePair(r, w)
}

// CreatePreimageChannel returns a FileChannel for the preimage oracle in a detached context
func CreatePreimageChannel() oppio.FileChannel {
	r := os.NewFile(PClientRFd, "preimage-oracle-read")
	w := os.NewFile(PClientWFd, "preimage-oracle-write")
	return oppio.NewReadWritePair(r, w)
}
