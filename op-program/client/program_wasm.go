//go:build wasm || wasip1
// +build wasm wasip1

package client

import (
	"errors"

	"github.com/ethereum/go-ethereum/log"

	cldr "github.com/ethereum-optimism/optimism/op-program/client/driver"
)

// Main executes the client program in a detached context and exits the current process.
// The client runtime environment must be preset before calling this function.
func Main(logger log.Logger) {
	log.Info("Starting fault proof program client")

	if err := RunProgramWithDefault(logger); errors.Is(err, cldr.ErrClaimNotValid) {
		log.Error("Claim is invalid", "err", err)
		wasm_output(1022)
		require(1)
	} else if err != nil {
		wasm_output(1023)
		log.Error("Program failed", "err", err)
		require(2)
	} else {
		wasm_output(1024)
		log.Info("Claim successfully verified")
	}
}
