package sgt

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-e2e/system/e2esys"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

type SgtHelper struct {
	T           *testing.T
	L2Client    *ethclient.Client
	SysCfg      e2esys.SystemConfig
	SgtContract *bindings.SoulGasToken
	ChainID     *big.Int
	sys         *e2esys.System
}

func NewSgtHelper(t *testing.T, ctx context.Context, sys *e2esys.System) *SgtHelper {
	// use sequencer's L2 client
	client := sys.NodeClient(e2esys.RoleSeq)
	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)

	sgtAddr := predeploys.SoulGasTokenAddr
	sgtContract, err := bindings.NewSoulGasToken(sgtAddr, client)
	require.NoError(t, err)

	return &SgtHelper{
		T:           t,
		L2Client:    client,
		SysCfg:      sys.Cfg,
		SgtContract: sgtContract,
		ChainID:     chainID,
		sys:         sys,
	}
}

func (s *SgtHelper) GetTestAccount(idx int) *ecdsa.PrivateKey {
	return s.sys.TestAccount(idx)
}

func (s *SgtHelper) depositSgtAndNativeFromGenesisAccountToAccount(t *testing.T, ctx context.Context, toAddr common.Address, sgtValue *big.Int, l2Value *big.Int) {
	privKey := s.GetTestAccount(0) // Genesis Account with lots of native balances
	// deposit some sgt and native tokens first
	txOpts, err := bind.NewKeyedTransactorWithChainID(privKey, s.ChainID)
	require.NoError(t, err)
	txOpts.Value = sgtValue
	sgtTx, err := s.SgtContract.BatchDepositForAll(txOpts, []common.Address{toAddr}, sgtValue)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, s.L2Client, sgtTx.Hash())
	require.NoError(t, err)
	nativeTx, err := s.transferNativeToken(t, ctx, privKey, toAddr, l2Value)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, s.L2Client, nativeTx.Hash())
	require.NoError(t, err)
}

func (s *SgtHelper) transferNativeToken(t *testing.T, ctx context.Context, sender *ecdsa.PrivateKey, toAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	chainID, err := s.L2Client.ChainID(ctx)
	require.NoError(t, err)
	gasFeeCap := big.NewInt(200)
	gasTipCap := big.NewInt(10)

	nonce, err := s.L2Client.NonceAt(ctx, crypto.PubkeyToAddress(sender.PublicKey), nil)
	require.NoError(t, err)
	tx := types.MustSignNewTx(sender, types.LatestSignerForChainID(chainID), &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &toAddr,
		Value:     amount,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       21000,
	})
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = s.L2Client.SendTransaction(ctx, tx)
	return tx, err
}
