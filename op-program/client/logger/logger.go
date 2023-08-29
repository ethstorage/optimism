//go:build !js
// +build !js

package logger

import (
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/log"
)

func NewLogger() log.Logger {
	return oplog.NewLogger(oplog.CLIConfig{
		Level:  "info",
		Format: "logfmt",
		Color:  false,
	})
}
