package main

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type wasmHostIO struct {
}

func (o wasmHostIO) Get(key [32]byte) []byte {
	size := wasm_input(0)

	buf := make([]byte, size)

	ssize := size / 8
	for i := uint64(0); i < ssize; i++ {
		data := wasm_input(0)
		binary.BigEndian.PutUint64(buf[i*8:], data)
	}

	data := wasm_input(0)
	sv := 8*(size%8) - 8
	for i := uint64(ssize * 8); i < size; i++ {
		buf[i] = byte(data >> sv)
		sv = sv - 8
	}

	// Integrity check
	// TODO: can use customized circuit to optimize
	require_bool(crypto.Keccak256Hash(buf) == key)
	return buf
}

//go:wasmimport env wasm_input
//go:noescape
func wasm_input(isPublic uint32) uint64

func (o wasmHostIO) Hint(v []byte) {
	// do nothing
}

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

func main() {
	zkmain()
}

func zkmain() {
	hio := wasmHostIO{}
	data := hio.Get(common.HexToHash("1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8"))
	require_bool(string(data) == "hello")

	data1 := hio.Get(common.HexToHash("acaf3289d7b601cbd114fb36c4d29c85bbfd5e133f14cb355c3fd8d99367964f"))
	require_bool(string(data1) == "Hello, World!")
}
