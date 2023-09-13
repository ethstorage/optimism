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
		binary.LittleEndian.PutUint64(buf[i*8:], data)
	}

	data := wasm_input(0)
	for i := uint64(ssize * 8); i < size; i++ {
		buf[i] = byte(data)
		data = data >> 8
	}

	// Integrity check
	// TODO: can use customized circuit to optimize
	if crypto.Keccak256Hash(buf) != key {
		panic("incorrect preimage data")
	}
	return buf
}

//go:wasmimport env wasm_input
//go:noescape
func wasm_input(isPublic uint32) uint64

func (o wasmHostIO) Hint(v []byte) {
	// do nothing
}

func main() {
	hio := wasmHostIO{}
	data := hio.Get(common.HexToHash("1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8"))
	if string(data) != "hello" {
		panic("preimage wrong")
	}

	data1 := hio.Get(common.HexToHash("acaf3289d7b601cbd114fb36c4d29c85bbfd5e133f14cb355c3fd8d99367964f"))
	if string(data1) != "Hello, World!" {
		panic("preimage wrong")
	}
}
