package blobs

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

	"github.com/ethereum-optimism/optimism/op-e2e"
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
	ctx       = context.Background()
	dacUrl    = fmt.Sprintf("http://127.0.0.1:%d", dacPort)
	toAddress = testutils.RandomAddress(mrand.New(mrand.NewSource(dacPort)))
	blobs     = make([]*eth.Blob, 3)
)

func TestFunctionSuccess(t *testing.T) {
	op_e2e.InitParallel(t)
	dacServer := StartDACServer(t)
	defer dacServer.Stop(ctx)

	sys, l2Client := StartSystemWithDAC(t)
	t.Cleanup(sys.Close)

	for i := range blobs {
		b := GetRandBlob(t, int64(i))
		blobs[i] = &b
	}

	tx, err := sendTransactionWithBlobs(t, ctx, l2Client, sys.TestAccount(0), toAddress, blobs)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, l2Client, tx.Hash())

	// TODO: wait for DAC server bug fix check in to refactor here
	for i, blobHash := range tx.BlobHashes() {
		dblobs, err := downloadBlobs(dacUrl, []common.Hash{blobHash})
		require.NoError(t, err)
		require.True(t, len(dblobs) == 1, "we should get one blob with blob hash")
		blob := dblobs[0]
		if len(blob) != eth.BlobSize {
			t.Error("Invalid downloaded blob len", "blob hash", blobHash, "blob len", len(blob))
		}
		if bytes.Compare(blob, blobs[i][:]) != 0 {
			t.Error("blob content", blob[:32], blobs[i][:32])
		}
	}
}

func StartSystemWithDAC(t *testing.T) (*e2esys.System, *ethclient.Client) {
	cfg := e2esys.DefaultSystemConfig(t)
	delete(cfg.Nodes, "verifier")
	if c, ok := cfg.Nodes["sequencer"]; ok {
		c.SafeDBPath = t.TempDir()
		c.DACConfig = &node.DACConfig{URLS: []string{dacUrl}}
		c.Driver.SequencerEnabled = true
	}
	cfg.DeployConfig.SequencerWindowSize = 30
	cfg.DeployConfig.FinalizationPeriodSeconds = 2
	cfg.SupportL1TimeTravel = true
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
	gasTipCap, gasFeeCap, blobFeeCap, err := GasPriceEstimator(ctx, l2Client)
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

func GasPriceEstimator(ctx context.Context, client *ethclient.Client) (*big.Int, *big.Int, *big.Int, error) {
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

func GetRandBlob(t *testing.T, seed int64) eth.Blob {
	r := mrand.New(mrand.NewSource(seed))
	bigData := eth.Data(make([]byte, eth.MaxBlobDataSize))
	for i := range bigData {
		bigData[i] = byte(r.Intn(256))
	}
	var b eth.Blob
	err := b.FromData(bigData)
	require.NoError(t, err)
	return b
}

func StartDACServer(t *testing.T) *da.Server {
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
