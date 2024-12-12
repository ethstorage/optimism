package sgt

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"math/rand"
	"testing"
	"time"

	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-e2e/faultproofs"
	"github.com/ethereum-optimism/optimism/op-e2e/system/helpers"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum-optimism/optimism/op-service/testutils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

var (
	vaultAddr                = predeploys.SequencerFeeVaultAddr
	dummyAddr                = common.Address{0xff, 0xff}
	errorInsufficientBalance = errors.New("insufficientBalance")
)

type BalanceCheck struct {
	sgtBalance       *big.Int
	gasCost          *big.Int
	l2Balance        *big.Int
	vaultBalanceDiff *big.Int
}

func TestDepositSGTBalance(t *testing.T) {
	op_e2e.InitParallel(t)
	sys, _ := faultproofs.StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)
	ctx := context.Background()

	sgt := NewSgtHelper(t, ctx, sys)

	tests := []struct {
		name              string
		depositSgtValue   *big.Int
		depositL2Value    *big.Int
		txValue           *big.Int
		expectedL2Balance *big.Int
		expectedErr       error
		action            func(ctx context.Context, t *testing.T, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) (*BalanceCheck, error)
	}{
		{
			name:              "SGTDepositSuccess",
			depositSgtValue:   big.NewInt(10000000000000),
			expectedL2Balance: big.NewInt(0),
			action:            sgtDeposit,
		},
		{
			name:            "NativaGasPaymentWithoutSGTSuccess",
			depositSgtValue: big.NewInt(0),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(0),
			action:          sgtTxSuccess,
		},
		{
			name:              "SGTTxWithoutNativeBalanceSuccess",
			depositSgtValue:   big.NewInt(10000000000000),
			depositL2Value:    big.NewInt(0),
			txValue:           big.NewInt(0),
			expectedL2Balance: big.NewInt(0),
			action:            sgtTxSuccess,
		},
		{
			name:              "FullSGTGasPaymentWithNativeBalanceSuccess",
			depositSgtValue:   big.NewInt(10000000000000),
			depositL2Value:    big.NewInt(10000000000000),
			txValue:           big.NewInt(0),
			expectedL2Balance: big.NewInt(10000000000000),
			action:            sgtTxSuccess,
		},
		{
			name:            "PartialSGTGasPaymentSuccess",
			depositSgtValue: big.NewInt(1000),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(0),
			action:          sgtTxSuccess,
		},
		{
			name:              "FullSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess",
			depositSgtValue:   big.NewInt(10000000000000),
			depositL2Value:    big.NewInt(10000000000000),
			txValue:           big.NewInt(10000),
			expectedL2Balance: big.NewInt(10000000000000 - 10000),
			action:            sgtTxSuccess,
		},
		{
			name:            "PartialSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess",
			depositSgtValue: big.NewInt(1000),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(10000),
			action:          sgtTxSuccess,
		},
		{
			name:            "InsufficientGasPaymentFail",
			depositSgtValue: big.NewInt(10000),
			depositL2Value:  big.NewInt(10000),
			txValue:         common.Big0,
			expectedErr:     errorInsufficientBalance,
			action:          sgtTxFail,
		},
		{
			name:            "PartialSGTGasPaymentAndNonZeroTxValueWithInsufficientNativeBalanceFail",
			depositSgtValue: big.NewInt(10000),
			depositL2Value:  big.NewInt(10000),
			txValue:         big.NewInt(10000000000000),
			expectedErr:     errorInsufficientBalance,
			action:          sgtTxFail,
		},
	}

	for index, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			b, err := tCase.action(ctx, t, int64(index), tCase.depositSgtValue, tCase.depositL2Value, tCase.txValue, sgt)
			if tCase.depositSgtValue.Cmp(b.gasCost) >= 0 {
				require.EqualValues(t, tCase.expectedErr, err)
				require.Equal(t, new(big.Int).Sub(tCase.depositSgtValue, b.gasCost).Cmp(b.sgtBalance), 0)
				require.Equal(t, tCase.expectedL2Balance.Cmp(b.l2Balance), 0)
			} else {
				balance := tCase.depositSgtValue.Add(tCase.depositSgtValue, tCase.depositL2Value)
				if balance.Cmp(b.gasCost) >= 0 {
					require.EqualValues(t, tCase.expectedErr, err)
					require.Equal(t, int64(0), b.sgtBalance.Int64())

					require.Equal(t, balance.Sub(balance, b.gasCost).Cmp(b.l2Balance.Add(b.l2Balance, tCase.txValue)), 0)
				} else {
					require.EqualValues(t, tCase.expectedErr, err)
				}
			}
			if b.vaultBalanceDiff != nil {
				// todo: how to check vault's balance diff
			}
		})
	}

}

func sgtDeposit(ctx context.Context, t *testing.T, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) (*BalanceCheck, error) {
	opts := &bind.CallOpts{Context: ctx}
	rng := rand.New(rand.NewSource(index))
	testPrivKey := testutils.InsecureRandomKey(rng)
	addr := crypto.PubkeyToAddress(testPrivKey.PublicKey)
	sgtBalance, err := sgt.SgtContract.BalanceOf(opts, addr)
	require.NoError(t, err)
	require.Equal(t, sgtBalance.Cmp(common.Big0), 0)

	privKey := sgt.GetTestAccount(0)
	txOpts, err := bind.NewKeyedTransactorWithChainID(privKey, sgt.ChainID)
	txOpts.Value = sgtValue
	tx, err := sgt.SgtContract.BatchDepositForAll(txOpts, []common.Address{addr}, sgtValue)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)

	sgtBalance, err = sgt.SgtContract.BalanceOf(opts, addr)
	require.NoError(t, err)
	l2Balance, err := sgt.L2Client.BalanceAt(ctx, addr, nil)
	require.NoError(t, err)
	if l2Balance == nil {
		l2Balance = common.Big0
	}

	sgtGas := common.Big0
	var vaultBalanceDiff *big.Int
	return &BalanceCheck{
		sgtBalance,
		sgtGas,
		l2Balance,
		vaultBalanceDiff,
	}, nil
}

func sgtTxSuccess(ctx context.Context, t *testing.T, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) (*BalanceCheck, error) {
	opts := &bind.CallOpts{Context: ctx}
	rng := rand.New(rand.NewSource(index))
	testPrivKey := testutils.InsecureRandomKey(rng)
	addr := crypto.PubkeyToAddress(testPrivKey.PublicKey)

	// check it's a fresh account
	sgtBalance, err := sgt.SgtContract.BalanceOf(opts, addr)
	require.NoError(t, err)
	require.Equal(t, int64(0), sgtBalance.Int64())
	l2Balance, err := sgt.L2Client.BalanceAt(ctx, addr, nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), l2Balance.Int64())
	// deposit sgt and send native balance first
	privKey := sgt.GetTestAccount(0)
	txOpts, err := bind.NewKeyedTransactorWithChainID(privKey, sgt.ChainID)
	require.NoError(t, err)
	txOpts.Value = sgtValue
	tx, err := sgt.SgtContract.BatchDepositForAll(txOpts, []common.Address{addr}, sgtValue)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)
	transferNativeToken(t, ctx, sgt, privKey, addr, l2Value)

	vaultBalanceBefore, err := sgt.L2Client.BalanceAt(ctx, vaultAddr, nil)
	// send tx with sgt as gas
	sgtGas := transferNativeToken(t, ctx, sgt, testPrivKey, dummyAddr, txValue)

	vaultBalanceAfter, err := sgt.L2Client.BalanceAt(ctx, vaultAddr, nil)
	vaultBalanceDiff := vaultBalanceAfter.Sub(vaultBalanceAfter, vaultBalanceBefore)
	sgtBalance, err = sgt.SgtContract.BalanceOf(opts, addr)
	require.NoError(t, err)

	l2Balance, err = sgt.L2Client.BalanceAt(ctx, addr, nil)
	require.NoError(t, err)

	if l2Balance == nil {
		l2Balance = big.NewInt(0)
	}

	return &BalanceCheck{
		sgtBalance,
		sgtGas,
		l2Balance,
		vaultBalanceDiff,
	}, nil
}

func sgtTxFail(ctx context.Context, t *testing.T, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) (*BalanceCheck, error) {
	opts := &bind.CallOpts{Context: ctx}
	rng := rand.New(rand.NewSource(index))
	testPrivKey := testutils.InsecureRandomKey(rng)
	addr := crypto.PubkeyToAddress(testPrivKey.PublicKey)

	// check it's a fresh account
	sgtBalance, err := sgt.SgtContract.BalanceOf(opts, addr)
	require.NoError(t, err)
	require.Equal(t, int64(0), sgtBalance.Int64())
	l2Balance, err := sgt.L2Client.BalanceAt(ctx, addr, nil)
	require.NoError(t, err)
	require.Equal(t, common.Big0.Cmp(l2Balance), 0)
	// deposit sgt and send native balance first
	privKey := sgt.GetTestAccount(0)
	txOpts, err := bind.NewKeyedTransactorWithChainID(privKey, sgt.ChainID)
	txOpts.Value = sgtValue
	depositTx, err := sgt.SgtContract.BatchDepositForAll(txOpts, []common.Address{addr}, sgtValue)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, sgt.L2Client, depositTx.Hash())
	require.NoError(t, err)
	transferNativeToken(t, ctx, sgt, privKey, addr, l2Value)

	// send tx with sgt as gas
	chainID, err := sgt.L2Client.ChainID(ctx)
	require.NoError(t, err)
	gasFeeCap := big.NewInt(200)
	gasTipCap := big.NewInt(10)
	tx := types.MustSignNewTx(testPrivKey, types.LatestSignerForChainID(chainID), &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     0, // Already have deposit
		To:        &dummyAddr,
		Value:     txValue,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       21000,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = sgt.L2Client.SendTransaction(ctx, tx)
	if err != nil {
		err = errorInsufficientBalance
	}
	// tx.GasCost() doesn't include the L1 fee, but it's large enough for this test case
	return &BalanceCheck{gasCost: tx.GasCost()}, err
}

func transferNativeToken(t *testing.T, ctx context.Context, sys *SgtHelper, sender *ecdsa.PrivateKey, toAddr common.Address, amount *big.Int) *big.Int {
	l2Seq := sys.L2Client
	gasTip := big.NewInt(10)
	receipt := helpers.SendL2Tx(t, sys.SysCfg, l2Seq, sender, func(opts *helpers.TxOpts) {
		opts.ToAddr = &toAddr
		opts.Value = amount
		opts.GasTipCap = gasTip
		opts.Gas = 21000
		opts.GasFeeCap = big.NewInt(200)
		nonce, err := l2Seq.NonceAt(ctx, crypto.PubkeyToAddress(sender.PublicKey), nil)
		require.NoError(t, err)
		opts.Nonce = nonce
	})

	require.Equal(t, receipt.Status, types.ReceiptStatusSuccessful)
	return calcGasFee(receipt)
}

func calcGasFee(receipt *types.Receipt) *big.Int {
	// OPStackTxFee = L2ExecutionGasFee + L1DataFee
	fees := new(big.Int).Mul(receipt.EffectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
	fees = fees.Add(fees, receipt.L1Fee)
	return fees
}
