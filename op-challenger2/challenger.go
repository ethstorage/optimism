package op_challenger2

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-challenger2/metrics"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-challenger2/config"
	"github.com/ethereum-optimism/optimism/op-challenger2/game"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
)

// Main is the programmatic entry-point for running op-challenger2 with a given configuration.
func Main(ctx context.Context, logger log.Logger, cfg *config.Config, m metrics.Metricer) (cliapp.Lifecycle, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	return game.NewService(ctx, logger, cfg, m)
}
