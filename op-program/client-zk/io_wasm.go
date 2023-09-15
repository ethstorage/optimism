//go:build js || wasm || wasip1
// +build js wasm wasip1

package client

import (
	"encoding/binary"

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
	_key := key.PreimageKey()
	_, _isPublic := key.(preimage.LocalIndexKey)
	var isPublic uint32
	if _isPublic {
		isPublic = 1
	}
	size := wasm_input(isPublic)
	buf := make([]byte, size)

	ssize := size / 8
	for i := uint64(0); i < ssize; i++ {
		data := wasm_input(isPublic)
		binary.BigEndian.PutUint64(buf[i*8:], data)
	}

	if ssize*8 < size {
		data := wasm_input(isPublic)
		var sv uint64 = 56
		for i := uint64(ssize * 8); i < size; i++ {
			buf[i] = byte(data >> sv)
			sv = sv - 8
		}
	}
	// Integrity check
	// TODO: can use customized circuit to optimize
	if !_isPublic {
		hash := crypto.Keccak256Hash(buf)
		hash[0] = _key[0]
		require_bool(hash == _key)
	}
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
