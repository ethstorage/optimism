package sgt

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum-optimism/optimism/op-e2e/bindings"
	"github.com/ethereum-optimism/optimism/op-e2e/system/e2esys"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

type SgtHelper struct {
	T           *testing.T
	Require     *require.Assertions
	L2Client    *ethclient.Client
	SysCfg      e2esys.SystemConfig
	SgtAddr     common.Address
	SgtContract *bindings.SoulGasToken
	ChainID     *big.Int
	sys         *e2esys.System
}

func NewSgtHelper(t *testing.T, ctx context.Context, sys *e2esys.System) *SgtHelper {
	require := require.New(t)
	// use sequencer's L2 client
	client := sys.NodeClient(e2esys.RoleSeq)
	chainID, err := client.ChainID(ctx)
	require.NoError(err)

	sgtAddr := predeploys.SoulGasTokenAddr
	sgtContract, err := bindings.NewSoulGasToken(sgtAddr, client)
	require.NoError(err)

	return &SgtHelper{
		T:           t,
		Require:     require,
		L2Client:    client,
		SysCfg:      sys.Cfg,
		SgtAddr:     sgtAddr,
		SgtContract: sgtContract,
		ChainID:     chainID,
		sys:         sys,
	}
}

func (s *SgtHelper) GetTestAccount(idx int) *ecdsa.PrivateKey {
	return s.sys.TestAccount(idx)
}
