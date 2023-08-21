package main

import (
	"github.com/ethereum-optimism/optimism/op-program/client"
	oplog "github.com/ethereum-optimism/optimism/op-program/client/logger"
)

func main() {
	// Default to a machine parsable but relatively human friendly log format.
	// Don't do anything fancy to detect if color output is supported.
	logger := oplog.NewLogger()
	client.Main(logger)
}
