//go:build js || wasm || wasip1
// +build js wasm wasip1
package types

import (
	"encoding/binary"
	"io"
)

//go:wasmimport env keccak_new
//go:noescape
func keccak_new(uint64)

//go:wasmimport env keccak_push
//go:noescape
func keccak_push(uint64)

//go:wasmimport env keccak_finalize
//go:noescape
func keccak_finalize() uint64

func NewKeccak256Helper() *Keccak256Helper {
	return &Keccak256Helper{
		data: make([]byte, 0),
	}
}

type Keccak256Helper struct {
	data []byte
}

func (b *Keccak256Helper) Write(p []byte) (n int, err error) {
	b.data = append(b.data, p...)
	return len(p), nil
}

func (b *Keccak256Helper) WriteTo(w io.Writer) (err error) {
	w.Write(b.data)
	return nil
}

func (b *Keccak256Helper) Hash() (hash [32]byte) {
	size := uint64(len(b.data))
	padding := size % 136
	if padding != 0 {
		padding = 136 - padding
	} else {
		padding = 136
	}
	data := make([]byte, size+padding)
	copy(data, b.data)
	totalLen := len(data)
	if padding == 1 {
		data[totalLen-1] = 0x81
	} else {
		data[size] = 0x01
		data[totalLen-1] = 0x80
	}
	round := totalLen / 136
	var hash_0, hash_1, hash_2, hash_3 uint64
	keccak_new(1)
	for i := 0; i < round; i++ {
		for j := 0; j < 17; j++ {
			start := i*136 + j*8
			val := binary.LittleEndian.Uint64(data[start : start+8])
			keccak_push(val)
		}
		hash_0 = keccak_finalize()
		hash_1 = keccak_finalize()
		hash_2 = keccak_finalize()
		hash_3 = keccak_finalize()
		keccak_new(0)
	}
	binary.LittleEndian.PutUint64(hash[:], hash_0)
	binary.LittleEndian.PutUint64(hash[8:], hash_1)
	binary.LittleEndian.PutUint64(hash[16:], hash_2)
	binary.LittleEndian.PutUint64(hash[24:], hash_3)
	return hash
}

/*
// for test

func keccak256check(input []byte, output []byte) {
	result := Keccak256Hash(input)
	for i := 0; i < len(result); i++ {
		if result[i] != output[i] {
			require(1)
			require(0)
		}
	}
}

func main() {
	input := make([]byte, 0)
	emtpy_output := []byte{
		197, 210, 70, 1, 134, 247, 35, 60, 146, 126, 125, 178, 220, 199, 3, 192, 229, 0, 182, 83,
		202, 130, 39, 59, 123, 250, 216, 4, 93, 133, 164, 112,
	}
	keccak256check(input, emtpy_output)

	input = []byte{102, 111, 111, 98, 97, 114, 97, 97}
	short_output := []byte{
		172, 132, 33, 155, 248, 181, 178, 245, 199, 105, 157, 164, 188, 53, 193, 25, 7, 35, 159,
		188, 30, 123, 91, 143, 30, 100, 188, 128, 172, 248, 137, 202,
	}
	keccak256check(input, short_output)
}
*/
