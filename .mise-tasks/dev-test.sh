#!/usr/bin/env bash

#MISE description="Developers' local tests"
#MISE alias="dt"

set -e
SECONDS=0

error_handler() {
    echo "Execution time: ${SECONDS} seconds"
    exit 1
}

trap 'error_handler' ERR

# Environment tests

forge --version

for var in SEPOLIA_RPC_URL MAINNET_RPC_URL; do
    if [ -z "${!var}" ]; then
        echo "Error: $var is not set."
        exit 1
    fi
done

STATUS=$(kurtosis engine status)
if echo "$STATUS" | grep -q "1.4.3"; then
    echo "Kurtosis engine is running."
else
    echo "The Kurtosis engine is not running, or there is a version mismatch."
    exit 1
fi

# Runs semgrep tests on the entire monorepo

just semgrep
just semgrep-test

# Solidity

cd packages/contracts-bedrock
just lint-check
just pre-pr
just test

# Go

cd ../..
make lint-go
make build-go

cd op-program && make op-program-client && cd ..
cd cannon && make elf && cd ..
cd op-e2e && make pre-test && cd ..

make devnet-allocs

export ENABLE_KURTOSIS=true
export OP_E2E_CANNON_ENABLED="false"
export OP_E2E_SKIP_SLOW_TEST=true
export OP_E2E_USE_HTTP=true
export ENABLE_ANVIL=true

# Note: not all packages are tested.
# For example the test `TestFinalization` in `op-alt-da` package fails even in upstream.
packages=(
    op-batcher
    op-chain-ops
    op-node
    op-proposer
    op-challenger
    op-dispute-mon
    op-conductor
    op-program
    op-service
    op-supervisor
    op-deployer
    op-e2e/system
    op-e2e/e2eutils
    op-e2e/opgeth
    op-e2e/interop
    op-e2e/actions
    op-e2e/faultproofs
    op-e2e/l2blob
    op-e2e/inbox
    op-e2e/sgt
    packages/contracts-bedrock/scripts/checks
)
formatted_packages=""
for package in "${packages[@]}"; do
    formatted_packages="$formatted_packages ./$package/..."
done

gotestsum --no-summary=skipped,output \
   --packages="$formatted_packages" \
   --format=short-verbose \
   --rerun-fails=2

echo "Execution time: $((SECONDS / 60)) minute(s) and $((SECONDS % 60)) second(s)"