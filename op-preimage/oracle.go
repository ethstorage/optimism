package preimage

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// OracleClient implements the Oracle by writing the pre-image key to the given stream,
// and reading back a length-prefixed value.
type OracleClient struct {
	rw io.ReadWriter
}

func NewOracleClient(rw io.ReadWriter) *OracleClient {
	return &OracleClient{rw: rw}
}

var _ Oracle = (*OracleClient)(nil)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (o *OracleClient) Get(key Key) []byte {
	h := key.PreimageKey()
	if _, err := o.rw.Write(h[:]); err != nil {
		panic(fmt.Errorf("failed to write key %s (%T) to pre-image oracle: %w", key, key, err))
	}

	var length uint64
	if err := binary.Read(o.rw, binary.LittleEndian, &length); err != nil {
		panic(fmt.Errorf("failed to read pre-image length of key %s (%T) from pre-image oracle: %w", key, key, err))
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(o.rw, payload); err != nil {
		panic(fmt.Errorf("failed to read pre-image payload (length %d) of key %s (%T) from pre-image oracle: %w", length, key, key, err))
	}
	return payload
}

// OracleServer serves the pre-image requests of the OracleClient, implementing the same protocol as the onchain VM.
type OracleServer struct {
	rw io.ReadWriter
}

func NewOracleServer(rw io.ReadWriter) *OracleServer {
	return &OracleServer{rw: rw}
}

type PreimageGetter func(key [32]byte) ([]byte, error)

var PreimageFile *os.File

func (o *OracleServer) NextPreimageRequest(getPreimage PreimageGetter) error {
	var key [32]byte
	if _, err := io.ReadFull(o.rw, key[:]); err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("failed to read requested pre-image key: %w", err)
	}
	value, err := getPreimage(key)
	if err != nil {
		return fmt.Errorf("failed to serve pre-image %s request: %w", hex.EncodeToString(key[:]), err)
	}

	if PreimageFile != nil {
		// write length & data to preimages.bin
		binary.Write(PreimageFile, binary.LittleEndian, uint64(len(value)))
		_, err := PreimageFile.Write(value)
		if err != nil {
			return fmt.Errorf("failed to dump pre-image binary file: %w", err)
		}

		// padding some zeros to make preimages length can be divided by 8
		if len(value)%8 != 0 {
			_, err := PreimageFile.Write(make([]byte, 8-len(value)%8))
			if err != nil {
				return fmt.Errorf("failed to dump pre-image binary file: %w", err)
			}
		}
	}

	if err := binary.Write(o.rw, binary.LittleEndian, uint64(len(value))); err != nil {
		return fmt.Errorf("failed to write length-prefix %d: %w", len(value), err)
	}
	if len(value) == 0 {
		return nil
	}
	if _, err := o.rw.Write(value); err != nil {
		return fmt.Errorf("failed to write pre-image value (%d long): %w", len(value), err)
	}
	return nil
}
