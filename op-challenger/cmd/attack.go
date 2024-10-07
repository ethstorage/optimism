// +build faultdisputegamen

package main

import (
	"context"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-challenger/flags"
	"github.com/ethereum-optimism/optimism/op-challenger/game/fault/contracts"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/urfave/cli/v2"
)

var (
	AttackBranchFlag = &cli.BoolFlag{
		Name:    "attack-branch",
		Usage:   "Specify the branch that needs to be attacked.",
		EnvVars: opservice.PrefixEnvVar(flags.EnvVarPrefix, "ATTACK_BRANCH"),
	}
	DaTypeFlag = &cli.StringFlag{
		Name:    "da-type",
		Usage:   "Type of DA (either 4844:1 or calldata:0).",
		EnvVars: opservice.PrefixEnvVar(flags.EnvVarPrefix, "DA_TYPE"),
	}
	ClaimsFlag = &cli.StringFlag{
		Name:    "claims",
		Usage:   "The claims.",
		EnvVars: opservice.PrefixEnvVar(flags.EnvVarPrefix, "CLAIMS"),
	}
)

func Attack(ctx *cli.Context) error {
	attackBranch := ctx.Uint64(AttackBranchFlag.Name)
	daType := ctx.Uint64(DaTypeFlag.Name)
	parentIndex := ctx.Uint64(ParentIndexFlag.Name)
	claims := []byte(ctx.String(ClaimsFlag.Name))

	contract, txMgr, err := NewContractWithTxMgr[contracts.FaultDisputeGameContract](ctx, GameAddressFlag.Name, contracts.NewFaultDisputeGameContract)
	if err != nil {
		return fmt.Errorf("failed to create dispute game bindings: %w", err)
	}

	parentClaim, err := contract.GetClaim(ctx.Context, parentIndex)
	if err != nil {
		return fmt.Errorf("failed to get parent claim: %w", err)
	}
	var tx txmgr.TxCandidate
	tx, err = contract.AttackV2Tx(ctx.Context, parentClaim, attackBranch, daType, claims)
	if err != nil {
		return fmt.Errorf("failed to create attack tx: %w", err)
	}
	rct, err := txMgr.Send(context.Background(), tx)
	if err != nil {
		return fmt.Errorf("failed to send tx: %w", err)
	}
	fmt.Printf("Sent tx with status: %v, hash: %s\n", rct.Status, rct.TxHash.String())

	return nil
}

func attackFlags() []cli.Flag {
	cliFlags := []cli.Flag{
		flags.L1EthRpcFlag,
		GameAddressFlag,
		AttackBranchFlag,
		DaTypeFlag,
		ParentIndexFlag,
		ClaimsFlag,
	}
	cliFlags = append(cliFlags, txmgr.CLIFlagsWithDefaults(flags.EnvVarPrefix, txmgr.DefaultChallengerFlagValues)...)
	cliFlags = append(cliFlags, oplog.CLIFlags(flags.EnvVarPrefix)...)
	return cliFlags
}

var AttackCommand = &cli.Command{
	Name:        "attack",
	Usage:       "Creates and sends a attack transaction to the dispute game",
	Description: "Creates and sends a attack transaction to the dispute game",
	Action:      Attack,
	Flags:       attackFlags(),
}
