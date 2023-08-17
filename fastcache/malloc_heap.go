//go:build appengine || windows || js || tinygo
// +build appengine windows js tinygo

package fastcache

func getChunk() []byte {
	return make([]byte, chunkSize)
}

func putChunk(chunk []byte) {
	// No-op.
}
