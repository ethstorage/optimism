package l2blob

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	mrand "math/rand"
	"testing"
	"time"

	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-e2e/system/e2esys"
	"github.com/ethereum-optimism/optimism/op-node/node"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/testutils"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus/misc/eip4844"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethstorage/da-server/pkg/da"
	"github.com/ethstorage/da-server/pkg/da/client"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

const (
	dacPort = 37777
)

var (
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	dacUrl = fmt.Sprintf("http://127.0.0.1:%d", dacPort)
)

func TestSubmitTXWithBlobsFunctionSuccess(t *testing.T) {
	op_e2e.InitParallel(t)
	dacServer := startDACServer(t)
	defer func() {
		if err := dacServer.Stop(ctx); err != nil {
			t.Errorf("Failed to stop DAC server: %v", err)
		}
	}()

	sys, l2Client := startSystemWithDAC(t)
	t.Cleanup(sys.Close)

	var (
		toAddress = testutils.RandomAddress(mrand.New(mrand.NewSource(dacPort)))
		blobs     = make([]*eth.Blob, 3)
	)
	for i := range blobs {
		b := getRandBlob(t, int64(i))
		blobs[i] = &b
	}

	tx, err := sendTransactionWithBlobs(t, ctx, l2Client, sys.TestAccount(0), toAddress, blobs)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, l2Client, tx.Hash())
	require.NoError(t, err)

	dblobs, err := downloadBlobs(dacUrl, tx.BlobHashes())
	require.NoError(t, err)
	require.True(t, len(dblobs) == len(tx.BlobHashes()), "blobs downloaded is not equal to blob hashes")

	for i, blob := range dblobs {
		require.True(t, len(blob) == eth.BlobSize, fmt.Sprintf("invalid downloaded blob, index %d; len %d", i, len(blob)))
		require.True(t, bytes.Equal(blob, blobs[i][:]), fmt.Sprintf("blob content diff: %s vs %s",
			common.Bytes2Hex(blob[:32]), common.Bytes2Hex(blobs[i][:32])))
	}
}

func startSystemWithDAC(t *testing.T) (*e2esys.System, *ethclient.Client) {
	cfg := e2esys.DefaultSystemConfig(t)
	delete(cfg.Nodes, "verifier")
	c, ok := cfg.Nodes["sequencer"]
	require.True(t, ok, "sequencer is required")
	c.DACConfig = &node.DACConfig{URLS: []string{dacUrl}}
	c.Driver.SequencerEnabled = true
	cfg.DeployConfig.L2GenesisBlobTimeOffset = new(hexutil.Uint64)
	// Disable proposer creating fast games automatically - required games are manually created
	cfg.DisableProposer = true
	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	return sys, sys.NodeClient(e2esys.RoleSeq)
}

func sendTransactionWithBlobs(t *testing.T, ctx context.Context, l2Client *ethclient.Client, sender *ecdsa.PrivateKey,
	toAddr common.Address, blobs []*eth.Blob) (*types.Transaction, error) {
	chainID, err := l2Client.ChainID(ctx)
	require.NoError(t, err)
	gasTipCap, gasFeeCap, blobFeeCap, err := gasPriceEstimator(ctx, l2Client)
	require.NoError(t, err)
	nonce, err := l2Client.NonceAt(ctx, crypto.PubkeyToAddress(sender.PublicKey), nil)
	require.NoError(t, err)
	sidecar, blobHashes, err := txmgr.MakeSidecar(blobs)
	require.NoError(t, err)
	tx := types.MustSignNewTx(sender, types.LatestSignerForChainID(chainID), &types.BlobTx{
		ChainID:    uint256.NewInt(chainID.Uint64()),
		Nonce:      nonce,
		GasFeeCap:  uint256.NewInt(gasFeeCap.Uint64()),
		GasTipCap:  uint256.NewInt(gasTipCap.Uint64()),
		Gas:        uint64(22000),
		To:         toAddr,
		Value:      uint256.NewInt(0),
		BlobFeeCap: uint256.NewInt(blobFeeCap.Uint64()),
		BlobHashes: blobHashes,
		Sidecar:    sidecar,
	})
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = l2Client.SendTransaction(ctx, tx)
	return tx, err
}

func gasPriceEstimator(ctx context.Context, client *ethclient.Client) (*big.Int, *big.Int, *big.Int, error) {
	tip, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	head, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	if head.BaseFee == nil {
		return nil, nil, nil, errors.New("head BaseFee is nil")
	}

	var blobFee *big.Int
	if head.ExcessBlobGas != nil {
		blobFee = eip4844.CalcBlobFee(*head.ExcessBlobGas)
	}

	gasFeeCap := new(big.Int).Add(
		tip,
		new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
	)
	return tip, gasFeeCap, blobFee, nil
}

func getRandBlob(t *testing.T, seed int64) eth.Blob {
	r := mrand.New(mrand.NewSource(seed))
	bigData := eth.Data(make([]byte, eth.MaxBlobDataSize))
	_, err := r.Read(bigData)
	require.NoError(t, err)
	var b eth.Blob
	err = b.FromData(bigData)
	require.NoError(t, err)
	return b
}

func startDACServer(t *testing.T) *da.Server {
	config := da.Config{
		SequencerIP: "127.0.0.1",
		ListenAddr:  fmt.Sprintf("0.0.0.0:%d", dacPort),
		StorePath:   t.TempDir(),
	}
	server := da.NewServer(&config)
	err := server.Start(ctx)
	require.NoError(t, err)

	return server
}

func downloadBlobs(dacUrl string, blobHashes []common.Hash) (blobs []hexutil.Bytes, err error) {
	client := client.New([]string{dacUrl})
	blobs, err = client.GetBlobs(ctx, blobHashes)
	return
}
