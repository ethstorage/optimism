//go:build js
// +build js

package client

import (
	"encoding/binary"
	"fmt"
	"syscall/js"
	"unsafe"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
)

func NewOracleClientAndHintWriter() (preimage.Oracle, preimage.Hinter) {
	o := wasmHostIO{}

	return o, o
}

type wasmHostIO struct {
}

func allocateBufferFunc(this js.Value, args []js.Value) interface{} {
	return int(allocateBuffer(uint32(args[0].Int())))
}

//export allocate_buffer
func allocateBuffer(size uint32) uintptr {
	// Allocate the in-Wasm memory region and returns its pointer to hosts.
	// The region is supposed to store random strings generated in hosts,
	// meaning that this is called "inside" of get_random_string.
	buf := make([]uint8, size)
	buf[0] = 3
	p := uintptr(unsafe.Pointer(&buf[0]))
	println("go uintptr:", &buf[0], p)

	return p
}

//go:wasmimport _gotest get_preimage_from_oracle
func getPreimageFromOracle(keyPtr uint32, retBufPtr uint32, retBufSize uint32)

//go:wasmimport _gotest hint_oracle
func hintOracle(retBufPtr uint32, retBufSize uint32)

func (o wasmHostIO) Get(key preimage.Key) []byte {
	var bufPtr *byte
	var bufSize uint32
	h := key.PreimageKey()
	fmt.Println("PreimageKey==========>", h)
	getPreimageFromOracle(
		uint32(uintptr(unsafe.Pointer(&h[0]))),
		uint32(uintptr(unsafe.Pointer(&bufPtr))),
		uint32(uintptr(unsafe.Pointer(&bufSize))))
	res := unsafe.Slice(bufPtr, bufSize)
	return res
}

func (o wasmHostIO) Hint(v preimage.Hint) {
	hint := v.Hint()
	var hintBytes []byte
	hintBytes = binary.BigEndian.AppendUint32(hintBytes, uint32(len(hint)))
	hintBytes = append(hintBytes, []byte(hint)...)
	hintOracle(uint32(uintptr(unsafe.Pointer(&hintBytes[0]))), uint32(len(hintBytes)))
}

func init() {
	//setup preimage oracle
	js.Global().Set("allocate_buffer", js.FuncOf(allocateBufferFunc))
}
