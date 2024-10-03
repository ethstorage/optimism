package chainconfig

import (
	"embed"
	"fmt"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/ethereum/go-ethereum/params"
)

//go:embed configs
var customChainConfigFS embed.FS

var OPSepoliaChainConfig, OPMainnetChainConfig *params.ChainConfig

func init() {
	mustLoadConfig := func(chainID uint64) *params.ChainConfig {
		cfg, err := params.LoadOPStackChainConfig(chainID)
		if err != nil {
			panic(err)
		}
		return cfg
	}
	OPSepoliaChainConfig = mustLoadConfig(11155420)
	OPMainnetChainConfig = mustLoadConfig(10)
	superchainEntry, err := customChainConfigFS.ReadDir("configs")
	if err == nil && len(superchainEntry) == 1 {
		superchain.LoadSuperchain(customChainConfigFS, superchainEntry[0])
	}

}

var L2ChainConfigsByChainID = map[uint64]*params.ChainConfig{
	11155420: OPSepoliaChainConfig,
	10:       OPMainnetChainConfig,
}

func RollupConfigByChainID(chainID uint64) (*rollup.Config, error) {
	config, err := rollup.LoadOPStackRollupConfig(chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rollup config for chain ID %d: %w", chainID, err)
	}
	return config, nil
}

func ChainConfigByChainID(chainID uint64) (*params.ChainConfig, error) {
	return params.LoadOPStackChainConfig(chainID)
}
