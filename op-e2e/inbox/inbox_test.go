package inbox

import (
	"context"
	"math/big"
	"testing"
	"time"

	batcherFlags "github.com/ethereum-optimism/optimism/op-batcher/flags"
	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"
	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-e2e/system/e2esys"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

var (
	ctx, _ = context.WithTimeout(context.Background(), 20*time.Second)
	cost   = big.NewInt(1500000000000000)
)

func TestBatchInboxFunctionSuccess(t *testing.T) {
	op_e2e.InitParallel(t)

	sys, l1Client := startSystemWithBatchInboxContract(t)
	t.Cleanup(sys.Close)

	// Wait for batch submitted and check event
	requireEventualBatcherTx(t, &sys.Cfg, l1Client, 8*time.Second)
}

func startSystemWithBatchInboxContract(t *testing.T) (*e2esys.System, *ethclient.Client) {
	cfg := e2esys.DefaultSystemConfig(t)
	cfg.DataAvailabilityType = batcherFlags.BlobsType
	cfg.BatcherTargetNumFrames = eth.MaxBlobsPerBlobTx
	cfg.DeployConfig.UseInboxContract = true
	c, ok := cfg.Nodes["sequencer"]
	require.True(t, ok, "sequencer is required")
	c.Driver.SequencerEnabled = true

	sys, err := cfg.Start(t, e2esys.StartOption{
		Key: "afterL1Start",
		Action: func(cfg *e2esys.SystemConfig, s *e2esys.System) {
			l1Client := s.NodeClient(e2esys.RoleL1)
			// Deploy mock storage contract
			mockStorageAddr := deployContract(t, cfg, l1Client, bindings.MockEthStorageMetaData, cost)
			// Deploy BatchInbox.sol contract
			batchInboxAddr := deployContract(t, cfg, l1Client, bindings.BatchInboxMetaData, mockStorageAddr)
			t.Logf("mock storage %s, batchInbox %s", mockStorageAddr.Hex(), batchInboxAddr.Hex())
			// Set BatchInboxAddress
			cfg.DeployConfig.BatchInboxAddress = batchInboxAddr
			// Deposit token
			transferNativeTokenToBatchInboxAddress(t, cfg, l1Client, new(big.Int).Mul(cost, big.NewInt(1000)))
		},
	})
	require.Nil(t, err, "Error starting up system")
	return sys, sys.NodeClient(e2esys.RoleL1)
}

func requireEventualBatcherTx(t *testing.T, cfg *e2esys.SystemConfig, l1Client *ethclient.Client, timeout time.Duration) {
	require.Eventually(t, func() bool {
		b, err := l1Client.BlockByNumber(ctx, nil)
		require.NoError(t, err)
		for _, tx := range b.Transactions() {
			if tx.To() == nil || tx.To().Cmp(cfg.DeployConfig.BatchInboxAddress) != 0 {
				continue
			}
			receipt, err := l1Client.TransactionReceipt(ctx, tx.Hash())
			require.NoError(t, err)
			if len(receipt.Logs) == 0 {
				continue
			}
			balanceBefore, err := l1Client.BalanceAt(ctx, receipt.Logs[0].Address, new(big.Int).Add(receipt.BlockNumber, big.NewInt(-1)))
			require.NoError(t, err)
			balanceAfter, err := l1Client.BalanceAt(ctx, receipt.Logs[0].Address, receipt.BlockNumber)
			require.NoError(t, err)
			require.True(t, balanceAfter.Uint64()-balanceBefore.Uint64() == cost.Uint64()*uint64(len(receipt.Logs)), "Cost is mismatch")
			return true
		}
		return false
	}, timeout, time.Second, "expected batcher tx type didn't arrive")
}

func deployContract(t *testing.T, cfg *e2esys.SystemConfig, client *ethclient.Client, meta *bind.MetaData,
	params ...interface{}) common.Address {
	ethPrivKey := cfg.Secrets.Batcher
	fromAddr := cfg.Secrets.Addresses().Batcher

	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	require.NoError(t, err)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(ethPrivKey, cfg.L1ChainIDBig())
	require.NoError(t, err)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	parsed, err := meta.GetAbi()
	require.NoError(t, err)

	address, tx, _, err := bind.DeployContract(auth, *parsed, common.FromHex(meta.Bin), client, params...)
	require.NoError(t, err)

	_, err = wait.ForReceiptOK(ctx, client, tx.Hash())
	require.NoError(t, err)
	return address
}

func transferNativeTokenToBatchInboxAddress(t *testing.T, cfg *e2esys.SystemConfig, client *ethclient.Client, amount *big.Int) {
	ethPrivKey := cfg.Secrets.Batcher
	fromAddr := cfg.Secrets.Addresses().Batcher

	gasTipCap, err := client.SuggestGasTipCap(ctx)
	require.NoError(t, err)
	head, err := client.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	gasFeeCap := new(big.Int).Add(
		gasTipCap,
		new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
	)

	nonce, err := client.NonceAt(ctx, fromAddr, nil)
	require.NoError(t, err)
	tx := types.MustSignNewTx(ethPrivKey, types.LatestSignerForChainID(cfg.L1ChainIDBig()), &types.DynamicFeeTx{
		ChainID:   cfg.L1ChainIDBig(),
		Nonce:     nonce,
		To:        &cfg.DeployConfig.BatchInboxAddress,
		Value:     amount,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       1000000,
	})
	err = client.SendTransaction(ctx, tx)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(context.Background(), client, tx.Hash())
	require.NoError(t, err)

	balance, err := client.BalanceAt(ctx, cfg.DeployConfig.BatchInboxAddress, nil)
	require.NoError(t, err)
	require.True(t, balance.Uint64() == amount.Uint64(), "balance no match")
}
