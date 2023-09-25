package main

import (
	smoke_test_client "github.com/ethereum-optimism/optimism/op-program/client/cmd/smoke_test/client"
	oplog "github.com/ethereum-optimism/optimism/op-program/client/logger"
)

func main() {
	// Default to a machine parsable but relatively human friendly log format.
	// Don't do anything fancy to detect if color output is supported.
	logger := oplog.NewLogger()
	smoke_test_client.Main(logger)
}
