package sources

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/zhiqiangxu/multicall"
)

type ESClient struct {
	logger         log.Logger
	client         *ethclient.Client
	abi            abi.ABI
	batchInboxAddr common.Address
}

func NewESClient(nodeAddr string, batchInboxAddr common.Address, logger log.Logger) (*ESClient, error) {
	client, err := ethclient.Dial(nodeAddr)
	if err != nil {
		return nil, err
	}

	abi, err := abi.JSON(strings.NewReader(abiDecentralizedKV))
	if err != nil {
		return nil, err
	}
	return &ESClient{logger: logger, client: client, abi: abi, batchInboxAddr: batchInboxAddr}, nil
}

func (s *ESClient) GetBlobs(blobHashes []eth.IndexedBlobHash) ([]*eth.BlobSidecar, error) {

	invokes := make([]multicall.Invoke, 0, len(blobHashes))
	for _, blobHash := range blobHashes {
		invokes = append(invokes, multicall.Invoke{
			Contract: common.HexToAddress("0x804C520d3c084C805E37A35E90057Ac32831F96f"), // TODO make it configurable
			Name:     "get",
			Args:     []interface{}{blobHash.Hash, uint8(0), big.NewInt(0), big.NewInt(blobSize)},
		})
	}

	output := make([][]byte, len(invokes))
	err := multicall.DoFrom(context.Background(), s.client, &s.abi, invokes, output, s.batchInboxAddr)
	if err != nil {
		return nil, err
	}

	blobSideCars := make([]*eth.BlobSidecar, 0, len(output))
	for i, blobBytes := range output {
		if len(blobBytes) != blobSize {
			return nil, fmt.Errorf("invalid blob size:%d, index:%d", len(blobBytes), blobHashes[i].Index)
		}
		var blob eth.Blob
		copy(blob[:], blobBytes)
		blobSideCar, err := makeBlobSideCar(blob, blobHashes[i].Index)
		if err != nil {
			return nil, err
		}
		blobSideCars = append(blobSideCars, blobSideCar)
	}

	return blobSideCars, nil
}

func makeBlobSideCar(blob eth.Blob, blobIndex uint64) (*eth.BlobSidecar, error) {

	rawBlob := *blob.KZGBlob()
	commitment, err := kzg4844.BlobToCommitment(rawBlob)
	if err != nil {
		return nil, fmt.Errorf("cannot compute KZG commitment of blob %d in tx candidate: %w", blobIndex, err)
	}

	proof, err := kzg4844.ComputeBlobProof(rawBlob, commitment)
	if err != nil {
		return nil, fmt.Errorf("cannot compute KZG proof for fast commitment verification of blob %d in tx candidate: %w", blobIndex, err)
	}

	return &eth.BlobSidecar{Blob: blob, Index: eth.Uint64String(blobIndex), KZGCommitment: eth.Bytes48(commitment), KZGProof: eth.Bytes48(proof)}, nil
}

const (
	abiDecentralizedKV = "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"kvIdx\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"lastKvIdx\",\"type\":\"uint256\"}],\"name\":\"Remove\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_maxKvSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_storageCost\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_dcfFactor\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"__init_KV\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dcfFactor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"exist\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"enumDecentralizedKV.DecodeType\",\"name\":\"decodeType\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"off\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"len\",\"type\":\"uint256\"}],\"name\":\"get\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"kvIndices\",\"type\":\"uint256[]\"}],\"name\":\"getKvMetas\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"hash\",\"outputs\":[{\"internalType\":\"bytes24\",\"name\":\"\",\"type\":\"bytes24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastKvIdx\",\"outputs\":[{\"internalType\":\"uint40\",\"name\":\"\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxKvSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"remove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"removeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"}],\"name\":\"size\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"storageCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"upfrontPayment\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	blobSize           = 32 * 4096
)
