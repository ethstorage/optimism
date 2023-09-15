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
	fmt.Printf("PreimageKey==========>%02x\n", h)
	if _, err := o.rw.Write(h[:]); err != nil {
		panic(fmt.Errorf("failed to write key %s (%T) to pre-image oracle: %w", key, key, err))
	}

	var length uint64
	if err := binary.Read(o.rw, binary.BigEndian, &length); err != nil {
		panic(fmt.Errorf("failed to read pre-image length of key %s (%T) from pre-image oracle: %w", key, key, err))
	}
	fmt.Println("PreimageSize==========>", length)
	payload := make([]byte, length)
	if _, err := io.ReadFull(o.rw, payload); err != nil {
		panic(fmt.Errorf("failed to read pre-image payload (length %d) of key %s (%T) from pre-image oracle: %w", length, key, key, err))
	}
	// os.Exit(3)
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

var Preimages = map[string]string{}

var PreimageFile, _ = os.Create("./bin/preimages.bin")

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

	//add preimage k,v to the Preimage map for dumping it to json file
	bytes := make([]byte, 32)
	copy(bytes[:], key[:])
	Preimages[hex.EncodeToString(bytes)] = hex.EncodeToString(value)

	//write length & data to preimages.bin
	binary.Write(PreimageFile, binary.BigEndian, uint64(len(value)))
	for i := 0; i < len(value); i++ {
		var ii = value[i]
		err := binary.Write(PreimageFile, binary.LittleEndian, ii)
		if err != nil {
			return fmt.Errorf("failed to dump pre-image binary file: %w", err)
		}
	}
	//padding some zeros to make preimages length can be divided by 8
	if len(value)%8 != 0 {
		for i := 0; i < 8-len(value)%8; i++ {
			err := binary.Write(PreimageFile, binary.LittleEndian, byte(0))
			if err != nil {
				return fmt.Errorf("failed to dump pre-image binary file: %w", err)
			}
		}
	}

	fmt.Printf("go buf is:===>%02x\n", value[len(value)-8:])

	if err := binary.Write(o.rw, binary.BigEndian, uint64(len(value))); err != nil {
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
