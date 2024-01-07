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

func (o wasmHostIO) OldGet(key preimage.Key) []byte {
	_key := key.PreimageKey()
	_, _isPublic := key.(preimage.LocalIndexKey)

	size := wasm_input(0)
	buf := make([]byte, size)

	ssize := size / 8
	for i := uint64(0); i < ssize; i++ {
		data := wasm_input(0)
		binary.BigEndian.PutUint64(buf[i*8:], data)
	}

	if ssize*8 < size {
		data := wasm_input(0)
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

func (o wasmHostIO) Uint8Get(key preimage.Key) []byte {
	_key := key.PreimageKey()
	_, _isPublic := key.(preimage.LocalIndexKey)

	size := wasm_input(0)
	padding := size % 136
	if padding != 0 {
		padding = 136 - padding
	} else {
		padding = 136
	}

	buf := make([]byte, size+padding)

	ssize := size / 8
	for i := uint64(0); i < ssize; i++ {
		data := wasm_input(0)
		binary.BigEndian.PutUint64(buf[i*8:], data)
	}

	if ssize*8 < size {
		data := wasm_input(0)
		var sv uint64 = 56
		for i := uint64(ssize * 8); i < size; i++ {
			buf[i] = byte(data >> sv)
			sv = sv - 8
		}
	}
	// Integrity check
	// TODO: can use customized circuit to optimize
	if !_isPublic {
		// hash := crypto.Keccak256Hash(buf)
		hash := Keccak256Hash(buf, size, padding)
		hash[0] = _key[0]
		require_bool(hash == _key)
	}
	return buf[:size]
}

func (o wasmHostIO) Get(key preimage.Key) []byte {
	_key := key.PreimageKey()
	_, _isPublic := key.(preimage.LocalIndexKey)

	size := wasm_input(0)
	padding := size % 136
	if padding != 0 {
		padding = 136 - padding
	} else {
		padding = 136
	}

	totalLen := size + padding
	totalUint64Len := totalLen / 8
	buf := make([]byte, totalLen)
	bufUint64 := make([]uint64, totalUint64Len)

	ssize := size / 8
	for i := uint64(0); i < ssize; i++ {
		bufUint64[i] = wasm_input(0)
		binary.LittleEndian.PutUint64(buf[i*8:], bufUint64[i])
	}

	if ssize*8 < size {
		data := wasm_input(0)
		var sv uint64 = 56
		for i := uint64(ssize * 8); i < size; i++ {
			buf[i] = byte(data >> sv)
			sv = sv - 8
		}
	}
	if padding == 1 {
		buf[totalLen-1] = 0x81
	} else {
		buf[size] = 0x01
		buf[totalLen-1] = 0x80
	}
	for i := ssize; i < totalUint64Len; i++ {
		bufUint64[i] = binary.LittleEndian.Uint64(buf[i*8:])
	}
	// Integrity check
	// TODO: can use customized circuit to optimize
	if !_isPublic {
		// hash := crypto.Keccak256Hash(buf)
		hash := Keccak256HashUint64(bufUint64)
		hash[0] = _key[0]
		require_bool(hash == _key)
	}
	return buf[:size]
}

func (o wasmHostIO) Hint(v preimage.Hint) {
	// do nothing
	return
}

//go:wasmimport env wasm_input
//go:noescape
func wasm_input(isPublic uint32) uint64

//go:wasmimport env wasm_output
//go:noescape
func wasm_output(value uint64)

//go:wasmimport env require
//go:noescape
func require(uint32)

//go:wasmimport env wasm_dbg
//go:noescape
func wasm_dbg(uint64)

func require_bool(cond bool) {
	if cond {
		require(1)
	} else {
		require(0)
	}
}
