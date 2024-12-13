package sgt

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"math/rand"
	"testing"

	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-e2e/faultproofs"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum-optimism/optimism/op-service/testutils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

var (
	seqVault                 = predeploys.SequencerFeeVaultAddr
	baseVault                = predeploys.BaseFeeVaultAddr
	l1Vault                  = predeploys.L1FeeVault
	dummyAddr                = common.Address{0xff, 0xff}
	errorInsufficientBalance = errors.New("insufficientBalance")
)

func TestSGTDepositFunctionSuccess(t *testing.T) {
	op_e2e.InitParallel(t)
	sys, _ := faultproofs.StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)
	ctx := context.Background()

	sgt := NewSgtHelper(t, ctx, sys)
	depositSgtValue := big.NewInt(10000)
	_, _ = setUpTestAccount(t, ctx, 0, sgt, depositSgtValue, big.NewInt(0))
}

// Diverse test scenarios to verify that the SoulGasToken(sgt) is utilized for gas payment firstly,
// unless there is insufficient sgt balance, in which case the native balance will be used instead.
func TestSGTAsGasPayment(t *testing.T) {
	op_e2e.InitParallel(t)
	sys, _ := faultproofs.StartFaultDisputeSystem(t)
	t.Cleanup(sys.Close)
	ctx := context.Background()

	sgt := NewSgtHelper(t, ctx, sys)
	// 1. setup a test account and deposit specified amount of sgt tokens (`depositSgtValue``) and native tokens (`depositL2Value``) into it.
	// 2. execute a token transfer tx to `dummyAddr` and validate that the gas payment behavior using sgt is as anticipated.
	tests := []struct {
		name            string
		depositSgtValue *big.Int
		depositL2Value  *big.Int
		txValue         *big.Int
		expectedErr     error
		action          func(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper)
	}{
		{
			name:            "NativaGasPaymentWithoutSGTSuccess",
			depositSgtValue: big.NewInt(0),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(0),
			action:          nativaGasPaymentWithoutSGTSuccess,
		},
		{
			name:            "FullSGTGasPaymentWithoutNativeBalanceSuccess",
			depositSgtValue: big.NewInt(10000000000000),
			depositL2Value:  big.NewInt(0),
			txValue:         big.NewInt(0),
			action:          fullSGTGasPaymentWithoutNativeBalanceSuccess,
		},
		{
			name:            "FullSGTGasPaymentWithNativeBalanceSuccess",
			depositSgtValue: big.NewInt(10000000000000),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(0),
			action:          fullSGTGasPaymentWithNativeBalanceSuccess,
		},
		{
			name:            "PartialSGTGasPaymentSuccess",
			depositSgtValue: big.NewInt(1000),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(0),
			action:          partialSGTGasPaymentSuccess,
		},
		{
			name:            "FullSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess",
			depositSgtValue: big.NewInt(10000000000000),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(10000),
			action:          fullSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess,
		},
		{
			name:            "PartialSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess",
			depositSgtValue: big.NewInt(1000),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(10000),
			action:          partialSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess,
		},
		{
			name:            "FullSGTInsufficientGasPaymentFail",
			depositSgtValue: big.NewInt(10000),
			depositL2Value:  big.NewInt(0),
			txValue:         big.NewInt(0),
			expectedErr:     errorInsufficientBalance,
			action:          fullSGTInsufficientGasPaymentFail,
		},
		{
			name:            "FullNativeInsufficientGasPaymentFail",
			depositSgtValue: big.NewInt(0),
			depositL2Value:  big.NewInt(10000),
			txValue:         big.NewInt(0),
			expectedErr:     errorInsufficientBalance,
			action:          fullNativeInsufficientGasPaymentFail,
		},
		{
			name:            "PartialSGTInsufficientGasPaymentFail",
			depositSgtValue: big.NewInt(10000),
			depositL2Value:  big.NewInt(10000),
			txValue:         big.NewInt(0),
			expectedErr:     errorInsufficientBalance,
			action:          partialSGTInsufficientGasPaymentFail,
		},
		{
			name:            "FullSGTGasPaymentAndNonZeroTxValueWithInsufficientNativeBalanceFail",
			depositSgtValue: big.NewInt(10000000000000),
			depositL2Value:  big.NewInt(10000),
			txValue:         big.NewInt(10001),
			expectedErr:     errorInsufficientBalance,
			action:          fullSGTGasPaymentAndNonZeroTxValueWithInsufficientNativeBalanceFail,
		},
		{
			name:            "PartialSGTGasPaymentAndNonZeroTxValueWithInsufficientNativeBalanceFail",
			depositSgtValue: big.NewInt(10000),
			depositL2Value:  big.NewInt(10000000000000),
			txValue:         big.NewInt(10000000000000 - 10000),
			expectedErr:     errorInsufficientBalance,
			action:          partialSGTGasPaymentAndNonZeroTxValueWithInsufficientNativeBalanceFail,
		},
	}

	for index, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			tCase.action(t, ctx, int64(index), tCase.depositSgtValue, tCase.depositL2Value, tCase.txValue, sgt)
		})
	}
}

func setUpTestAccount(t *testing.T, ctx context.Context, index int64, sgt *SgtHelper, depositSgtValue *big.Int, depositL2Value *big.Int) (*ecdsa.PrivateKey, common.Address) {
	opts := &bind.CallOpts{Context: ctx}
	rng := rand.New(rand.NewSource(index))
	testPrivKey := testutils.InsecureRandomKey(rng)
	testAddr := crypto.PubkeyToAddress(testPrivKey.PublicKey)

	// check it's a fresh account
	sgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, int64(0), sgtBalance.Int64())
	l2Balance, err := sgt.L2Client.BalanceAt(ctx, testAddr, nil)
	require.NoError(t, err)
	require.Equal(t, int64(0), l2Balance.Int64())

	// deposit initial sgt and native(L2) balance to the test account
	sgt.depositSgtAndNativeFromGenesisAccountToAccount(t, ctx, testAddr, depositSgtValue, depositL2Value)
	// ensure that sgt and native balance of testAccount are correctly initialized
	preSgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, depositSgtValue.Cmp(preSgtBalance), 0)
	preL2Balance, err := sgt.L2Client.BalanceAt(ctx, testAddr, nil)
	require.NoError(t, err)
	require.Equal(t, depositL2Value.Cmp(preL2Balance), 0)
	return testPrivKey, testAddr
}

// balance invariant check: preTotalBalance = postTotalBalance + gasCost + txValue
func invariantBalanceCheck(t *testing.T, ctx context.Context, sgt *SgtHelper, addr common.Address, gasCost *big.Int, txValue *big.Int, preSgtBalance *big.Int, preL2Balance *big.Int, postSgtBalance *big.Int) {
	postL2Balance, err := sgt.L2Client.BalanceAt(ctx, addr, nil)
	require.NoError(t, err)
	preBalance := preSgtBalance.Add(preSgtBalance, preL2Balance)
	postBalance := postSgtBalance.Add(postSgtBalance, gasCost)
	postBalance = postBalance.Add(postBalance, txValue)
	postBalance = postBalance.Add(postBalance, postL2Balance)
	require.Equal(t, 0, preBalance.Cmp(postBalance))
}

func nativaGasPaymentWithoutSGTSuccess(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that sgt balance is 0
	require.Equal(t, common.Big0.Cmp(sgtValue), 0)
	// ensure that txValue is 0
	require.Equal(t, common.Big0.Cmp(txValue), 0)
	testAccount, testAddr := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.NoError(t, err)
	receipt, err := wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)
	gasCost := calcGasFee(receipt)

	// post sgt balance check: it should be 0
	opts := &bind.CallOpts{Context: ctx}
	postSgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, common.Big0.Cmp(postSgtBalance), 0)
	// balance invariant check
	invariantBalanceCheck(t, ctx, sgt, testAddr, gasCost, txValue, sgtValue, l2Value, postSgtBalance)
}

func fullSGTGasPaymentWithoutNativeBalanceSuccess(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that native balance is 0
	require.Equal(t, common.Big0.Cmp(l2Value), 0)
	// ensure that txValue is 0
	require.Equal(t, common.Big0.Cmp(txValue), 0)
	testAccount, testAddr := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.NoError(t, err)
	receipt, err := wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)
	gasCost := calcGasFee(receipt)

	// post sgt balance check: sgt should be used as gas first
	opts := &bind.CallOpts{Context: ctx}
	postSgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, new(big.Int).Add(postSgtBalance, gasCost).Cmp(sgtValue), 0)
	// balance invariant check
	invariantBalanceCheck(t, ctx, sgt, testAddr, gasCost, txValue, sgtValue, l2Value, postSgtBalance)
}

func fullSGTGasPaymentWithNativeBalanceSuccess(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that native balance is positive
	require.Equal(t, common.Big0.Cmp(l2Value), -1)
	// ensure that txValue is 0
	require.Equal(t, common.Big0.Cmp(txValue), 0)
	testAccount, testAddr := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.NoError(t, err)
	receipt, err := wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)
	gasCost := calcGasFee(receipt)

	// post sgt balance check: sgt should be used as gas first
	opts := &bind.CallOpts{Context: ctx}
	postSgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, new(big.Int).Add(postSgtBalance, gasCost).Cmp(sgtValue), 0)
	// balance invariant check
	invariantBalanceCheck(t, ctx, sgt, testAddr, gasCost, txValue, sgtValue, l2Value, postSgtBalance)
}

func partialSGTGasPaymentSuccess(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that native balance is positive
	require.Equal(t, common.Big0.Cmp(l2Value), -1)
	// ensure that txValue is 0
	require.Equal(t, common.Big0.Cmp(txValue), 0)
	testAccount, testAddr := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.NoError(t, err)
	receipt, err := wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)
	gasCost := calcGasFee(receipt)

	// post sgt balance check: sgt should be used as gas first and should be spent all
	opts := &bind.CallOpts{Context: ctx}
	postSgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, common.Big0.Cmp(postSgtBalance), 0)
	// balance invariant check
	invariantBalanceCheck(t, ctx, sgt, testAddr, gasCost, txValue, sgtValue, l2Value, postSgtBalance)
}

func fullSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that native balance is positive
	require.Equal(t, common.Big0.Cmp(l2Value), -1)
	// ensure that txValue is positive
	require.Equal(t, common.Big0.Cmp(txValue), -1)
	testAccount, testAddr := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.NoError(t, err)
	receipt, err := wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)
	gasCost := calcGasFee(receipt)

	// post sgt balance check: sgt should be used as gas first
	opts := &bind.CallOpts{Context: ctx}
	postSgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, new(big.Int).Add(postSgtBalance, gasCost).Cmp(sgtValue), 0)
	// balance invariant check
	invariantBalanceCheck(t, ctx, sgt, testAddr, gasCost, txValue, sgtValue, l2Value, postSgtBalance)
}

func partialSGTGasPaymentAndNonZeroTxValueWithSufficientNativeBalanceSuccess(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that native balance is positive
	require.Equal(t, common.Big0.Cmp(l2Value), -1)
	// ensure that txValue is positive
	require.Equal(t, common.Big0.Cmp(txValue), -1)
	testAccount, testAddr := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.NoError(t, err)
	receipt, err := wait.ForReceiptOK(ctx, sgt.L2Client, tx.Hash())
	require.NoError(t, err)
	gasCost := calcGasFee(receipt)

	// post sgt balance check: sgt should be used as gas first and should be spent all
	opts := &bind.CallOpts{Context: ctx}
	postSgtBalance, err := sgt.SgtContract.BalanceOf(opts, testAddr)
	require.NoError(t, err)
	require.Equal(t, common.Big0.Cmp(postSgtBalance), 0)
	// balance invariant check
	invariantBalanceCheck(t, ctx, sgt, testAddr, gasCost, txValue, sgtValue, l2Value, postSgtBalance)
}

func fullSGTInsufficientGasPaymentFail(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that native balance is 0
	require.Equal(t, common.Big0.Cmp(l2Value), 0)
	// ensure that txValue is 0
	require.Equal(t, common.Big0.Cmp(txValue), 0)
	testAccount, _ := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	_, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.Error(t, err)
}

func fullNativeInsufficientGasPaymentFail(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that sgt balance is 0
	require.Equal(t, common.Big0.Cmp(sgtValue), 0)
	// ensure that txValue is 0
	require.Equal(t, common.Big0.Cmp(txValue), 0)
	testAccount, _ := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	_, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.Error(t, err)
}

func partialSGTInsufficientGasPaymentFail(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that txValue is 0
	require.Equal(t, common.Big0.Cmp(txValue), 0)
	testAccount, _ := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	_, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	require.Error(t, err)
}

func fullSGTGasPaymentAndNonZeroTxValueWithInsufficientNativeBalanceFail(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	// ensure that native balance is less than txValue
	require.Equal(t, l2Value.Cmp(txValue), -1)
	testAccount, _ := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	// ensure sgt balance is large enough to cover the gas cost
	require.Equal(t, sgtValue.Cmp(tx.GasCost()), 1)
	require.Error(t, err)
}

func partialSGTGasPaymentAndNonZeroTxValueWithInsufficientNativeBalanceFail(t *testing.T, ctx context.Context, index int64, sgtValue *big.Int, l2Value *big.Int, txValue *big.Int, sgt *SgtHelper) {
	testAccount, _ := setUpTestAccount(t, ctx, index, sgt, sgtValue, l2Value)

	// make a simple tx with the testAccount: transfer txValue from testAccount to dummyAddr
	tx, err := sgt.transferNativeToken(t, ctx, testAccount, dummyAddr, txValue)
	// ensure native balance is large enough to cover the gas cost
	require.Equal(t, l2Value.Cmp(tx.GasCost()), 1)
	require.Error(t, err)
}

func calcGasFee(receipt *types.Receipt) *big.Int {
	// OPStackTxFee = L2ExecutionGasFee + L1DataFee
	fees := new(big.Int).Mul(receipt.EffectiveGasPrice, big.NewInt(int64(receipt.GasUsed)))
	fees = fees.Add(fees, receipt.L1Fee)
	return fees
}
