//go:build js || wasm || wasip1
// +build js wasm wasip1

package client

import (
	"encoding/binary"
	"fmt"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum/go-ethereum/crypto"
)

func NewOracleClientAndHintWriter() (preimage.Oracle, preimage.Hinter) {
	o := wasmHostIO{}

	return o, o
}

type wasmHostIO struct {
}

func (o wasmHostIO) Get(key preimage.Key) []byte {
	// _key := key.PreimageKey()
	size := wasm_input(0)
	fmt.Println("go size:", size)

	buf := make([]byte, size)

	ssize := size / 8
	for i := uint64(0); i < ssize; i++ {
		data := wasm_input(0)
		binary.BigEndian.PutUint64(buf[i*8:], data)
	}

	if ssize*8 < size {
		data := wasm_input(0)
		sv := 56
		for i := uint64(ssize * 8); i < size; i++ {
			buf[i] = byte(data >> sv)
			sv = sv - 8
		}
	}
	// Integrity check
	// TODO: can use customized circuit to optimize
	fmt.Printf("buf:%02x\n", buf[size-8:size])
	fmt.Printf("crypto.Keccak256Hash:%02x\n", crypto.Keccak256Hash(buf))
	// require_bool(crypto.Keccak256Hash(buf) == _key)
	return buf
}

func (o wasmHostIO) Hint(v preimage.Hint) {
	//do nothing
	return
}

//go:wasmimport env wasm_input
//go:noescape
func wasm_input(isPublic uint32) uint64

//go:wasmimport env require
//go:noescape
func require(uint32)

func require_bool(cond bool) {
	if cond {
		require(1)
	} else {
		require(0)
	}
}
