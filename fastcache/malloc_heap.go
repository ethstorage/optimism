//go:build appengine || windows || js || tinygo || wasip1
// +build appengine windows js tinygo wasip1

package fastcache

func getChunk() []byte {
	return make([]byte, chunkSize)
}

func putChunk(chunk []byte) {
	// No-op.
}
