package opcm

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
)

type UpgradeAnchorStateRegistryInput struct {
	Release                  string
	StandardVersionsToml     string
	VmAddress                common.Address
	GameKind                 string
	GameType                 uint32
	AbsolutePrestate         common.Hash
	MaxGameDepth             uint64
	SplitDepth               uint64
	ClockExtension           uint64
	MaxClockDuration         uint64
	DelayedWethProxy         common.Address
	AnchorStateRegistryProxy common.Address
	L2ChainId                uint64
	Proposer                 common.Address
	Challenger               common.Address
}

func (input *UpgradeAnchorStateRegistryInput) InputSet() bool {
	return true
}

type UpgradeAnchorStateRegistryOutput struct {
	DisputeGameImpl common.Address
}

func (output *UpgradeAnchorStateRegistryOutput) CheckOutput(input common.Address) error {
	return nil
}

func UpgradeAnchorStateRegistry(
	host *script.Host,
	input UpgradeAnchorStateRegistryInput,
) (UpgradeAnchorStateRegistryOutput, error) {
	return RunBasicScript[UpgradeAnchorStateRegistryInput, UpgradeAnchorStateRegistryOutput](host, input, "UpgradeAnchorStateRegistry.s.sol", "UpgradeAnchorStateRegistry")
}
