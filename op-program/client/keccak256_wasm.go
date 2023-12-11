//go:build js || wasm || wasip1
// +build js wasm wasip1

package client

import "encoding/binary"

//go:wasmimport env keccak_new
//go:noescape
func keccak_new(uint64)

//go:wasmimport env keccak_push
//go:noescape
func keccak_push(uint64)

//go:wasmimport env keccak_finalize
//go:noescape
func keccak_finalize() uint64

func Keccak256Hash(data ...[]byte) (output [32]byte) {
	dataBytes := make([]byte, 0)
	for _, value := range data {
		dataBytes = append(dataBytes, value...)
	}
	require_bool(len(dataBytes)%8 == 0)
	input := ByteSliceToUint64Slice(dataBytes)
	keccak_new(0)
	for _, value := range input {
		keccak_push(value)
	}
	result := make([]uint64, 0)
	for i := 0; i < 4; i++ {
		result = append(result, keccak_finalize())
	}
	resultBytes := Uint64SliceToByteSlice(result)
	for i := 0; i < 32; i++ {
		output[i] = resultBytes[i]
	}
	return output
}

func Uint64SliceToByteSlice(uint64Slice []uint64) []byte {
	byteSlice := make([]byte, len(uint64Slice)*8)
	for i, val := range uint64Slice {
		binary.LittleEndian.PutUint64(byteSlice[i*8:], val)
	}
	return byteSlice
}

func ByteSliceToUint64Slice(byteSlice []byte) []uint64 {
	require_bool(len(byteSlice)%8 == 0)
	uint64Slice := make([]uint64, len(byteSlice)/8)
	for i := 0; i < len(byteSlice); i += 8 {
		val := binary.LittleEndian.Uint64(byteSlice[i:])
		uint64Slice[i/8] = val
	}
	return uint64Slice
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
