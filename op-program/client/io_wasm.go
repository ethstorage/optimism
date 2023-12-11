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

	size := wasm_input(0)
	ssize := size / 8
	uint64Size := (size + 7) / 8
	buf := make([]byte, size)
	bufU64 := make([]uint64, uint64Size)
	for i := uint64(0); i < ssize; i++ {
		bufU64[i] = wasm_input(0)
		binary.BigEndian.PutUint64(buf[i*8:], bufU64[i])
	}

	if ssize*8 < size {
		bufU64[uint64Size-1] = wasm_input(0)
		var sv uint64 = 56
		for i := uint64(ssize * 8); i < size; i++ {
			buf[i] = byte(bufU64[uint64Size-1] >> sv)
			sv = sv - 8
		}
	}
	// Integrity check
	// TODO: can use customized circuit to optimize
	if !_isPublic {
		hash := Keccak256HashInputU64(bufU64)
		hash[0] = _key[0]
		if _key[1] != 79 {
			print_log_flag()
			for _, val := range buf {
				wasm_log(uint64(val))
			}
			print_log_flag()
			for _, val := range hash {
				wasm_log(uint64(val))
			}
			print_log_flag()
			for _, val := range _key {
				wasm_log(uint64(val))
			}
		}

		require_bool(hash == _key)
	}
	return buf
}

func (o wasmHostIO) GetBak(key preimage.Key) []byte {
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

//go:wasmimport env wasm_log
//go:noescape
func wasm_log(value uint64)

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

func print_log_flag() {
	wasm_log(1)
	wasm_log(1)
	wasm_log(1)
	wasm_log(0)
	wasm_log(0)
	wasm_log(0)
}
