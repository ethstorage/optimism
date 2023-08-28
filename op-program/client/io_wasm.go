//go:build tinygo
// +build tinygo

package client

import (
	"encoding/binary"
	"unsafe"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
)

func NewOracleClientAndHintWriter() (preimage.Oracle, preimage.Hinter) {
	o := wasmHostIO{}

	return o, o
}

type wasmHostIO struct {
}

//export allocate_buffer
func allocateBuffer(size uint32) *uint8 {
	// Allocate the in-Wasm memory region and returns its pointer to hosts.
	// The region is supposed to store random strings generated in hosts,
	// meaning that this is called "inside" of get_random_string.
	buf := make([]uint8, size)
	return &buf[0]
}

func (o wasmHostIO) Get(key preimage.Key) []byte {
	var bufPtr *byte
	var bufSize uint32
	h := key.PreimageKey()
	getPreimageFromOracle(h, &bufPtr, &bufSize)
	res := unsafe.Slice(bufPtr, bufSize)
	return res
}

func (o wasmHostIO) Hint(v preimage.Hint) {
	hint := v.Hint()
	var hintBytes []byte
	hintBytes = binary.BigEndian.AppendUint32(hintBytes, uint32(len(hint)))
	hintBytes = append(hintBytes, []byte(hint)...)
	hintOracle(&hintBytes[0], uint32(len(hintBytes)))
}

//export get_preimage_from_oracle
func getPreimageFromOracle(key [32]byte, retBufPtr **byte, retBufSize *uint32)

//export hint_oracle
func hintOracle(retBufPtr *byte, retBufSize uint32)
