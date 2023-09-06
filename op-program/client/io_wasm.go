//go:build js || wasm || wasip1
// +build js wasm wasip1

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

//go:wasmimport _gotest get_preimage_len
//go:noescape
func getPreimageLenFromOracle(keyPtr uint32) uint32

//go:wasmimport _gotest get_preimage_from_oracle
//go:noescape
func getPreimageFromOracle(keyPtr uint32, retBufPtr uint32, size uint32) uint32

//go:wasmimport _gotest hint_oracle
//go:noescape
func hintOracle(retBufPtr uint32, retBufSize uint32)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (o wasmHostIO) Get(key preimage.Key) []byte {
	h := key.PreimageKey()
	//get preimage size
	size := getPreimageLenFromOracle(uint32(uintptr(unsafe.Pointer(&h[0]))))

	buf := make([]byte, size)
	readedLen := getPreimageFromOracle(uint32(uintptr(unsafe.Pointer(&h[0]))), uint32(uintptr(unsafe.Pointer(&buf[0]))), size)
	if readedLen < size {
		getPreimageFromOracle(uint32(uintptr(unsafe.Pointer(&h[0]))), uint32(uintptr(unsafe.Pointer(&buf[readedLen]))), size-readedLen)
	}
	return buf
}

func (o wasmHostIO) Hint(v preimage.Hint) {
	hint := v.Hint()
	var hintBytes []byte
	hintBytes = binary.BigEndian.AppendUint32(hintBytes, uint32(len(hint)))
	hintBytes = append(hintBytes, []byte(hint)...)
	hintOracle(uint32(uintptr(unsafe.Pointer(&hintBytes[0]))), uint32(len(hintBytes)))
}
