//go:build js
// +build js

package client

import (
	"encoding/binary"
	"fmt"
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
func getPreimageLenFromOracle(keyPtr uint32) uint32

//go:wasmimport _gotest get_preimage_from_oracle
func getPreimageFromOracle(keyPtr uint32, retBufPtr uint32, size uint32)

//go:wasmimport _gotest hint_oracle
func hintOracle(retBufPtr uint32, retBufSize uint32)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (o wasmHostIO) Get(key preimage.Key) []byte {
	h := key.PreimageKey()
	fmt.Printf("PreimageKey==========>%02x\n", h)
	//get preimage size
	size := getPreimageLenFromOracle(uint32(uintptr(unsafe.Pointer(&h[0]))))
	fmt.Println("PreimageSize==========>", size)

	buf := make([]byte, size)
	getPreimageFromOracle(uint32(uintptr(unsafe.Pointer(&h[0]))), uint32(uintptr(unsafe.Pointer(&buf[0]))), size)

	trunc := min(32, int(size))
	fmt.Printf("PreimageBytes==========>%02x\n", buf[0:trunc])
	return buf
}

func (o wasmHostIO) Hint(v preimage.Hint) {
	hint := v.Hint()
	var hintBytes []byte
	hintBytes = binary.BigEndian.AppendUint32(hintBytes, uint32(len(hint)))
	hintBytes = append(hintBytes, []byte(hint)...)
	hintOracle(uint32(uintptr(unsafe.Pointer(&hintBytes[0]))), uint32(len(hintBytes)))
}
