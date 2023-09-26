//go:build !(wasm || wasip1)
// +build !wasm,!wasip1

package client

import (
	"errors"
	"os"

	"github.com/ethereum/go-ethereum/log"

	cldr "github.com/ethereum-optimism/optimism/op-program/client/driver"
	"github.com/ethereum-optimism/optimism/op-program/client/l1"
	"github.com/ethereum-optimism/optimism/op-program/client/l2"
)

// Main executes the client program in a detached context and exits the current process.
// The client runtime environment must be preset before calling this function.
func Main(logger log.Logger) {
	log.Info("Starting fault proof program client")

	if err := RunProgramWithDefault(logger); errors.Is(err, cldr.ErrClaimNotValid) {
		log.Error("Claim is invalid", "err", err)
		os.Exit(1)
	} else if err != nil {
		log.Error("Program failed", "err", err)
		os.Exit(2)
	} else {
		log.Info("Claim successfully verified")
		os.Exit(0)
	}
}

// RunProgramWithDefault executes the Program, while attached to an IO based pre-image oracle, to be served by a host.
func RunProgramWithDefault(logger log.Logger) error {
	pClient, hClient := NewOracleClientAndHintWriter()
	l1PreimageOracle := l1.NewPreimageOracle(pClient, hClient)
	l2PreimageOracle := l2.NewPreimageOracle(pClient, hClient)

	bootInfo := NewBootstrapClient(pClient).BootInfo()
	logger.Info("Program Bootstrapped", "bootInfo", bootInfo)
	return runDerivation(
		logger,
		bootInfo.RollupConfig,
		bootInfo.L2ChainConfig,
		bootInfo.L1Head,
		bootInfo.L2Head,
		bootInfo.L2Claim,
		bootInfo.L2ClaimBlockNumber,
		l1PreimageOracle,
		l2PreimageOracle,
	)
}
