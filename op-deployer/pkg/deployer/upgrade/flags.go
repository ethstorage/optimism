package upgrade

import (
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/bootstrap"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	"github.com/urfave/cli/v2"
)

var ASRFlags = []cli.Flag{
	deployer.L1RPCURLFlag,
	deployer.WorkdirFlag,
	deployer.PrivateKeyFlag,
	bootstrap.ArtifactsLocatorFlag,
}

var Commands = []*cli.Command{
	{
		Name:   "asr",
		Usage:  "Upgrade ASR implementation.",
		Flags:  cliapp.ProtectFlags(ASRFlags),
		Action: ASRCLI,
	},
}
