#!/usr/bin/env bash
set -e

#MISE description="Developers' local tests"
#MISE alias="dt"

forge --version

cd packages/contracts-bedrock

# pwd

# just pre-pr

# just test

# cd ../..

# make lint-go

# make build

# # Build other essential components
# cd op-program && make op-program-client && cd ..
# cd cannon && make elf && cd ..
# cd op-e2e && make pre-test && cd ..

# export ENABLE_KURTOSIS=true
# export OP_E2E_CANNON_ENABLED="false"
# export OP_E2E_SKIP_SLOW_TEST=true
# export OP_E2E_USE_HTTP=true
# export ENABLE_ANVIL=true

# # Set your own RPC URLs below.
# export SEPOLIA_RPC_URL=<YOUR_SEPOLIA_RPC_URL>
# export MAINNET_RPC_URL=<YOUR_MAINNET_RPC_URL>


#             packages=(
#                 op-batcher
#                 op-chain-ops
#                 op-node
#                 op-proposer
#                 op-challenger
#                 op-dispute-mon
#                 op-conductor
#                 op-program
#                 op-service
#                 op-supervisor
#                 op-deployer
#                 op-e2e/system
#                 op-e2e/e2eutils
#                 op-e2e/opgeth
#                 op-e2e/interop
#                 op-e2e/actions
#                 op-e2e/faultproofs
#                 packages/contracts-bedrock/scripts/checks
#             )
#             formatted_packages=""
#             for package in "${packages[@]}"; do
#               formatted_packages="$formatted_packages ./$package/..."
#             done

# gotestsum --no-summary=skipped,output \
#    --packages="$formatted_packages" \
#    --format=short-verbose \
#    --rerun-fails=2
#    --format=testname \
#    --rerun-fails=2


# cd op-e2e
# make test-actions
# make test-ws