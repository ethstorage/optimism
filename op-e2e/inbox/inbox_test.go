package inbox

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	batcherFlags "github.com/ethereum-optimism/optimism/op-batcher/flags"
	op_e2e "github.com/ethereum-optimism/optimism/op-e2e"
	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/transactions"
	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils/wait"
	"github.com/ethereum-optimism/optimism/op-e2e/system/e2esys"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

// MockStorageMetaData contains all meta data concerning the L1Block contract.
var MockStorageMetaData = &bind.MetaData{
	ABI: "[{\"inputs\": [{\"internalType\": \"uint256\",\"name\": \"_cost\",\"type\": \"uint256\"}],\"stateMutability\": \"nonpayable\",\"type\": \"constructor\"},{\"anonymous\": false,\"inputs\": [{\"indexed\": true,\"internalType\": \"uint256\",\"name\": \"kvIdx\",\"type\": \"uint256\"},{\"indexed\": true,\"internalType\": \"uint256\",\"name\": \"kvSize\",\"type\": \"uint256\"},{\"indexed\": true,\"internalType\": \"bytes32\",\"name\": \"dataHash\",\"type\": \"bytes32\"}],\"name\": \"PutBlob\",\"type\": \"event\"},{\"inputs\": [],\"name\": \"kvEntryCount\",\"outputs\": [{\"internalType\": \"uint256\",\"name\": \"\",\"type\": \"uint256\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"bytes32\",\"name\": \"_key\",\"type\": \"bytes32\"},{\"internalType\": \"uint256\",\"name\": \"_blobIdx\",\"type\": \"uint256\"},{\"internalType\": \"uint256\",\"name\": \"_length\",\"type\": \"uint256\"}],\"name\": \"putBlob\",\"outputs\": [],\"stateMutability\": \"payable\",\"type\": \"function\"},{\"inputs\": [],\"name\": \"upfrontPayment\",\"outputs\": [{\"internalType\": \"uint256\",\"name\": \"\",\"type\": \"uint256\"}],\"stateMutability\": \"view\",\"type\": \"function\"}]",
	Bin: "60a0604052348015600e575f5ffd5b506040516104fc3803806104fc8339818101604052810190602e9190606d565b8060808181525050506093565b5f5ffd5b5f819050919050565b604f81603f565b81146058575f5ffd5b50565b5f815190506067816048565b92915050565b5f60208284031215607f57607e603b565b5b5f608a84828501605b565b91505092915050565b6080516104526100aa5f395f60aa01526104525ff3fe608060405260043610610033575f3560e01c80631ccbc6da146100375780634581a92014610061578063638ba9e91461007d575b5f5ffd5b348015610042575f5ffd5b5061004b6100a7565b60405161005891906101c6565b60405180910390f35b61007b60048036038101906100769190610240565b6100ce565b005b348015610088575f5ffd5b506100916101a9565b60405161009e91906101c6565b60405180910390f35b5f7f0000000000000000000000000000000000000000000000000000000000000000905090565b5f824990505f5f1b8103610117576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161010e90610310565b60405180910390fd5b5f5f1b840361015b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101529061039e565b60405180910390fd5b5f5f54905060015f5461016e91906103e9565b5f819055508183827f8b7a21215282409938287ae262331bfe6411d35d3d46aa7e505ef02000524ac260405160405180910390a45050505050565b5f5481565b5f819050919050565b6101c0816101ae565b82525050565b5f6020820190506101d95f8301846101b7565b92915050565b5f5ffd5b5f819050919050565b6101f5816101e3565b81146101ff575f5ffd5b50565b5f81359050610210816101ec565b92915050565b61021f816101ae565b8114610229575f5ffd5b50565b5f8135905061023a81610216565b92915050565b5f5f5f60608486031215610257576102566101df565b5b5f61026486828701610202565b93505060206102758682870161022c565b92505060406102868682870161022c565b9150509250925092565b5f82825260208201905092915050565b7f45746853746f72616765436f6e74726163743a206661696c656420746f2067655f8201527f7420626c6f622068617368000000000000000000000000000000000000000000602082015250565b5f6102fa602b83610290565b9150610305826102a0565b604082019050919050565b5f6020820190508181035f830152610327816102ee565b9050919050565b7f45746853746f72616765436f6e74726163743a206661696c656420746f2067655f8201527f7420626c6f62206b657900000000000000000000000000000000000000000000602082015250565b5f610388602a83610290565b91506103938261032e565b604082019050919050565b5f6020820190508181035f8301526103b58161037c565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6103f3826101ae565b91506103fe836101ae565b9250828201905080821115610416576104156103bc565b5b9291505056fea2646970667358221220843189a0c55b11adcf6318b93e4bb7cddf32d681f005f2ceb5eec50aa5cf3c4d64736f6c634300081c0033",
}

// BatchInboxMetaData contains all meta data concerning the L1Block contract.
var BatchInboxMetaData = &bind.MetaData{
	ABI: "[{\"inputs\": [{\"internalType\": \"address\",\"name\": \"_esStorageContract\",\"type\": \"address\"}],\"stateMutability\": \"nonpayable\",\"type\": \"constructor\"},{\"inputs\": [],\"name\": \"BalanceNotEnough\",\"type\": \"error\"},{\"stateMutability\": \"payable\",\"type\": \"fallback\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"\",\"type\": \"address\"}],\"name\": \"balances\",\"outputs\": [{\"internalType\": \"uint256\",\"name\": \"\",\"type\": \"uint256\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"_to\",\"type\": \"address\"}],\"name\": \"deposit\",\"outputs\": [],\"stateMutability\": \"payable\",\"type\": \"function\"},{\"inputs\": [],\"name\": \"esStorageContract\",\"outputs\": [{\"internalType\": \"contract StorageContract\",\"name\": \"\",\"type\": \"address\"}],\"stateMutability\": \"view\",\"type\": \"function\"},{\"inputs\": [{\"internalType\": \"address\",\"name\": \"_to\",\"type\": \"address\"},{\"internalType\": \"uint256\",\"name\": \"_amount\",\"type\": \"uint256\"}],\"name\": \"withdraw\",\"outputs\": [],\"stateMutability\": \"nonpayable\",\"type\": \"function\"},{\"stateMutability\": \"payable\",\"type\": \"receive\"}]",
	Bin: "60a060405234801561000f575f5ffd5b50604051610984380380610984833981810160405281019061003191906100c9565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff1681525050506100f4565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100988261006f565b9050919050565b6100a88161008e565b81146100b2575f5ffd5b50565b5f815190506100c38161009f565b92915050565b5f602082840312156100de576100dd61006b565b5b5f6100eb848285016100b5565b91505092915050565b60805161086a61011a5f395f81816101990152818161022901526103b3015261086a5ff3fe608060405260043610610042575f3560e01c806322eb767d1461006f57806327e235e314610099578063f340fa01146100d5578063f3fef3a3146100f15761005b565b3661005b576100513334610119565b610059610178565b005b6100653334610119565b61006d610178565b005b34801561007a575f5ffd5b506100836103b1565b604051610090919061057f565b60405180910390f35b3480156100a4575f5ffd5b506100bf60048036038101906100ba91906105d7565b6103d5565b6040516100cc919061061a565b60405180910390f35b6100ef60048036038101906100ea91906105d7565b6103e9565b005b3480156100fc575f5ffd5b506101176004803603810190610112919061065d565b6103f6565b005b5f81031561017457805f5f8473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f82825461016c91906106c8565b925050819055505b5050565b5f5f90505f5f90505f5b824990505f5f1b8103156102c5575f8203610227577f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16631ccbc6da6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610200573d5f5f3e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610224919061070f565b91505b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16634581a920838386620200006040518563ffffffff1660e01b81526004016102889392919061078b565b5f604051808303818588803b15801561029f575f5ffd5b505af11580156102b1573d5f5f3e3d5ffd5b505050505082806001019350506001610182575b5f83036102d4575050506103af565b5f83836102e191906107c0565b90505f5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205490508181101561035d576040517f9882883500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b81816103699190610801565b5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f208190555050505050505b565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f602052805f5260405f205f915090505481565b6103f38134610119565b50565b5f5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f2054905081811015610470576040517f9882883500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b818161047c9190610801565b5f5f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f20819055508273ffffffffffffffffffffffffffffffffffffffff166108fc8390811502906040515f60405180830381858888f193505050501580156104ff573d5f5f3e3d5ffd5b50505050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f61054761054261053d84610505565b610524565b610505565b9050919050565b5f6105588261052d565b9050919050565b5f6105698261054e565b9050919050565b6105798161055f565b82525050565b5f6020820190506105925f830184610570565b92915050565b5f5ffd5b5f6105a682610505565b9050919050565b6105b68161059c565b81146105c0575f5ffd5b50565b5f813590506105d1816105ad565b92915050565b5f602082840312156105ec576105eb610598565b5b5f6105f9848285016105c3565b91505092915050565b5f819050919050565b61061481610602565b82525050565b5f60208201905061062d5f83018461060b565b92915050565b61063c81610602565b8114610646575f5ffd5b50565b5f8135905061065781610633565b92915050565b5f5f6040838503121561067357610672610598565b5b5f610680858286016105c3565b925050602061069185828601610649565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6106d282610602565b91506106dd83610602565b92508282019050808211156106f5576106f461069b565b5b92915050565b5f8151905061070981610633565b92915050565b5f6020828403121561072457610723610598565b5b5f610731848285016106fb565b91505092915050565b5f819050919050565b61074c8161073a565b82525050565b5f819050919050565b5f61077561077061076b84610752565b610524565b610602565b9050919050565b6107858161075b565b82525050565b5f60608201905061079e5f830186610743565b6107ab602083018561060b565b6107b8604083018461077c565b949350505050565b5f6107ca82610602565b91506107d583610602565b92508282026107e381610602565b915082820484148315176107fa576107f961069b565b5b5092915050565b5f61080b82610602565b915061081683610602565b925082820390508181111561082e5761082d61069b565b5b9291505056fea2646970667358221220d3a380efe4d22309c483b9af73fa3cb244a03a404fa85e10382e6a18e9a5d89764736f6c634300081c0033",
}

var (
	ctx, _          = context.WithTimeout(context.Background(), 10*time.Second)
	cost            = big.NewInt(1500000000000000)
	depositVal      = new(big.Int).Mul(cost, big.NewInt(10))
	mockStorageAddr = common.Address{}
)

func TestBatchInboxFunctionSuccess(t *testing.T) {
	op_e2e.InitParallel(t)

	sys, l1Client := startSystemWithBatchInboxContract(t)
	t.Cleanup(sys.Close)

	sendTxs(t, &sys.Cfg, l1Client)

	// Wait for batch submitted and check event
	requireEventualBatcherTx(t, &sys.Cfg, l1Client, 8*time.Second)
}

func startSystemWithBatchInboxContract(t *testing.T) (*e2esys.System, *ethclient.Client) {
	cfg := e2esys.DefaultSystemConfig(t)
	cfg.DataAvailabilityType = batcherFlags.BlobsType
	cfg.BatcherTargetNumFrames = eth.MaxBlobsPerBlobTx
	c, ok := cfg.Nodes["sequencer"]
	require.True(t, ok, "sequencer is required")
	c.SafeDBPath = t.TempDir()
	c.Driver.SequencerEnabled = true

	sys, err := cfg.Start(t, e2esys.StartOption{
		Key: "afterL1Start",
		Action: func(cfg *e2esys.SystemConfig, s *e2esys.System) {
			l1Client := s.NodeClient(e2esys.RoleL1)
			// Deploy mock storage contract
			mockStorageAddr = deployContract(t, cfg, l1Client, MockStorageMetaData, cost)
			// Deploy BatchInbox.sol contract
			batchInboxAddr := deployContract(t, cfg, l1Client, BatchInboxMetaData, mockStorageAddr)
			// Deposit token
			transferNativeTokenToBatchInboxAddress(t, cfg, l1Client, depositVal)
			// Set BatchInboxAddress
			cfg.DeployConfig.BatchInboxAddress = batchInboxAddr
		},
	})
	require.Nil(t, err, "Error starting up system")
	return sys, sys.NodeClient(e2esys.RoleL1)
}

func requireEventualBatcherTx(t *testing.T, cfg *e2esys.SystemConfig, l1Client *ethclient.Client, timeout time.Duration) {
	var foundOtherTxType bool
	require.Eventually(t, func() bool {
		b, err := l1Client.BlockByNumber(ctx, nil)
		require.NoError(t, err)
		for _, tx := range b.Transactions() {
			if tx.To().Cmp(cfg.DeployConfig.BatchInboxAddress) != 0 {
				continue
			}
			receipt, err := l1Client.TransactionReceipt(ctx, tx.Hash())
			require.NoError(t, err)
			require.True(t, len(receipt.Logs) > 0, "Storage event missing")
			balanceBefore, err := l1Client.BalanceAt(ctx, receipt.Logs[0].Address, new(big.Int).Add(receipt.BlockNumber, big.NewInt(-1)))
			require.NoError(t, err)
			balanceAfter, err := l1Client.BalanceAt(ctx, receipt.Logs[0].Address, receipt.BlockNumber)
			require.NoError(t, err)
			require.True(t, balanceAfter.Uint64()-balanceBefore.Uint64() == cost.Uint64()*uint64(len(receipt.Logs)), "Cost is mismatch")
		}
		return false
	}, timeout, time.Second, "expected batcher tx type didn't arrive")
	require.False(t, foundOtherTxType, "unexpected batcher tx type found")
}

func sendTxs(t *testing.T, cfg *e2esys.SystemConfig, l1Client *ethclient.Client) []*types.Transaction {
	ethPrivKey := cfg.Secrets.Alice
	fromAddr := cfg.Secrets.Addresses().Alice

	// Send deposit transactions in a loop to drive up L1 base fee
	depAmount := big.NewInt(1_000_000_000_000)
	const numDeps = 3
	txs := make([]*types.Transaction, 0, numDeps)
	t.Logf("Sending %d deposits...", numDeps)
	for i := int64(0); i < numDeps; i++ {
		opts, err := bind.NewKeyedTransactorWithChainID(ethPrivKey, cfg.L1ChainIDBig())
		require.NoError(t, err)
		opts.Value = depAmount
		opts.Nonce = big.NewInt(i)
		depositContract, err := bindings.NewOptimismPortal(cfg.L1Deployments.OptimismPortalProxy, l1Client)
		require.NoError(t, err)

		tx, err := transactions.PadGasEstimate(opts, 2, func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return depositContract.DepositTransaction(opts, fromAddr, depAmount, 1_000_000, false, nil)
		})
		require.NoErrorf(t, err, "failed to send deposit tx[%d]", i)
		t.Logf("Deposit submitted[%d]: tx hash: %v", i, tx.Hash())
		txs = append(txs, tx)
	}
	require.Len(t, txs, numDeps)
	return txs
}

func deployContract(t *testing.T, cfg *e2esys.SystemConfig, client *ethclient.Client, meta *bind.MetaData,
	params ...interface{}) common.Address {
	ethPrivKey := cfg.Secrets.Alice
	fromAddr := cfg.Secrets.Addresses().Alice

	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	require.NoError(t, err)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(ethPrivKey, cfg.L1ChainIDBig())
	require.NoError(t, err)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // 部署时无需转账
	auth.GasLimit = uint64(3000000) // 设置足够的Gas限制
	auth.GasPrice = gasPrice

	bytecode, err := hex.DecodeString(meta.Bin)
	require.NoError(t, err)
	parsed, err := meta.GetAbi()
	require.NoError(t, err)

	address, tx, _, err := bind.DeployContract(auth, *parsed, bytecode, client, params...)
	require.NoError(t, err)

	_, err = wait.ForReceiptOK(ctx, client, tx.Hash())
	require.NoError(t, err)
	return address
}

func transferNativeTokenToBatchInboxAddress(t *testing.T, cfg *e2esys.SystemConfig, client *ethclient.Client, amount *big.Int) {
	ethPrivKey := cfg.Secrets.Alice
	fromAddr := cfg.Secrets.Addresses().Alice

	chainID, err := client.ChainID(ctx)
	require.NoError(t, err)
	gasFeeCap := big.NewInt(200)
	gasTipCap := big.NewInt(10)

	nonce, err := client.NonceAt(ctx, fromAddr, nil)
	require.NoError(t, err)
	tx := types.MustSignNewTx(ethPrivKey, types.LatestSignerForChainID(chainID), &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &cfg.DeployConfig.BatchInboxAddress,
		Value:     amount,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Gas:       3000000,
	})
	err = client.SendTransaction(ctx, tx)
	require.NoError(t, err)
	_, err = wait.ForReceiptOK(ctx, client, tx.Hash())
	require.NoError(t, err)
}
