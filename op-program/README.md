# op-program-wasm

## build wasm

```
cd op-program
op-program-client-mips-wasm
```

## problems with op-wasm and solutions
### mock syscall
- `github.com/gofrs/flock`:
    - modify `go-ethereum\core\rawdb\freezer.go`,`go-ethereum\core\rawdb\chain_freezer.go`
    - add `go-ethereum\core\rawdb\fileutil_mock.go`
    - modify `op-geth/node/node_mock.go`
    - add `op-geth/node/fileutil_mock.go`
- `/go-ethereum/ethdb/leveldb`:
    - `replace github.com/syndtr/goleveldb => ./goleveldb`
    - add `goleveldb/leveldb/storage/file_storage_wasm.go`
- `go.uber.org/fx`
    - `replace go.uber.org/fx => ./fx`
    - add `fx/app_wasm.go`
- `github.com/libp2p/go-libp2p`:
    - `replace github.com/libp2p/go-libp2p => ./go-libp2p`
    - add `go-libp2p/p2p/transport/websocket/websocket_wasm.go`
- `/go-ethereum/trie` & `VictoriaMetrics/fastcache`: use [arb's fastcache](https://github.com/OffchainLabs/fastcache) `replace github.com/VictoriaMetrics/fastcache => ./fastcache`
    - Goetherum related commit:
    - [First commit](https://github.com/OffchainLabs/go-ethereum/commits?after=1319d385dc35f0a3be7166ec4a63ce83de89c376+244&author=PlasmaPower)
    - [Another offchainlabs](https://github.com/OffchainLabs/go-ethereum/commits?author=Tristan-Wilson&before=1319d385dc35f0a3be7166ec4a63ce83de89c376+70)

### how to remove syscall/host-io
- [1th discussion between ethstorage & HO](https://drive.google.com/file/d/15XvyltLqXBuF26LD8hRYo7CsuLyajcUv/view?usp=drivesdk)
- [current explorations](https://literate-wolfsbane-bf0.notion.site/Wasm-Compile-Explore-6f92d2a5d23f4dc18232d1b15266f422?pvs=4)

## Arb reference
- [code structure](https://docs.google.com/presentation/d/1I2eZ9qp6tkyuV2w-z-ffRaVwpOTMUloBX3n_tfvtdWo/edit#slide=id.p)
    - Challenge manager:
    ![image](https://literate-wolfsbane-bf0.notion.site/image/https%3A%2F%2Fs3-us-west-2.amazonaws.com%2Fsecure.notion-static.com%2F75a34284-0857-485c-ab91-401b67161d50%2FUntitled.png?table=block&id=caa17fae-dd87-4cd1-8bd5-5b2346637457&spaceId=1c6bced1-30e7-47d5-9ac7-860911f88272&width=2000&userId=&cache=v2)

- How to build:
    1. follow [arbitrum build nitro tutorial](https://docs.arbitrum.io/node-running/how-tos/build-nitro-locally) to install requirements (Don't `make`, all you'll have to clean the builded file)
    2. `make build build-replay-env`

- Compile replay.wasm separately: `make build-wasm-bin`

- Test a wavm(Prover) can [load wasm file](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/arbitrator/prover/src/main.rs#L185) compiled from `arbitrator/prover/test-cases/go/main.go`, [parse wasm](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/arbitrator/prover/src/machine.rs#L868C19-L868C24)to op-codes, [replace the imports with Rust functions](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/arbitrator/prover/src/machine.rs#L301), [transfrom wasm op-code to wavm op-code](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/arbitrator/prover/src/machine.rs#L341), [run the wavm Op-codes](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/arbitrator/prover/src/main.rs#L226) using its [customized VM](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/arbitrator/prover/src/machine.rs#L1370):
    `make contracts/test/prover/proofs/go.json`

- Test Prover can load Rust compiled wasm file which **requires third-party wasm libraries**, and run it like the above test case:
    `make contracts/test/prover/proofs/rust-host-io.json`
- [serialize wavm to binary `machine.wavm.br`](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/arbitrator/prover/src/main.rs#L200) for valiator to generate proof:
    `make target/machines/latest/machine.wavm.br`

- Validator use `machine.wavm.br` to generate proof:
    - [Validator spawn a arbitrator thread](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/valnode/valnode.go#L126)
    - [load `machine.wavm.br`](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/server_arb/validator_spawner.go#L64)
    - call `machine.wavm.br`'s c_binded function to execute or generate proof, for example:
        - [execute](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/server_arb/validator_spawner.go#L122)
        - [call `machine.wavm.br`'s c_binded function](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/server_arb/machine.go#L124)

- Validator use JIT to run replay.wasm (compiled by `cmd/replay/main.go`):
    - [Validator spawn a JIT thread](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/valnode/valnode.go#L103)
    - configuaration: [replay.wasm](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/server_jit/machine_loader.go#L20), [JIT binary](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/server_jit/machine_loader.go#L31)
    - [load config](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/server_jit/spawner.go#L48)
    - [execute](https://github.com/OffchainLabs/nitro/blob/6a2078a8b91a74d826fff37c0f3d2ecdcd6c4bff/validator/server_jit/jit_machine.go#L34)




# op-program

Implements a fault proof program that runs through the rollup state-transition to verify an L2 output from L1 inputs.
This verifiable output can then resolve a disputed output on L1.

The program is designed such that it can be run in a deterministic way such that two invocations with the same input
data wil result in not only the same output, but the same program execution trace. This allows it to be run in an
on-chain VM as part of the dispute resolution process.


## Run [cannon](https://github.com/ethstorage/optimism/tree/develop/cannon)

Replace the `{APIKEY}` to yours in the following commands.

```
git clone https://github.com/ethstorage/optimism.git
# Build op-program server-mode and MIPS-client binaries.
cd op-program
make op-program # build

# Switch back to cannon, and build the CLI
cd ../cannon
make cannon

# Transform MIPS op-program client binary into first VM state.
# This outputs state.json (VM state) and meta.json (for debug symbols).
./bin/cannon load-elf --path=../op-program/bin/op-program-client.elf

# Run cannon emulator (with example inputs)
# Note that the server-mode op-program command is passed into cannon (after the --),
# it runs as sub-process to provide the pre-image data.
#
# Note:
#  - The L2 RPC is an archive L2 node on OP goerli.
#  - The L1 RPC is a non-archive RPC, also change `--l1.rpckind` to reflect the correct L1 RPC type.
./bin/cannon run \
    --pprof.cpu \
    --info-at '%10000000' \
    --proof-at never \
    --input ./state.json \
    -- \
    ../op-program/bin/op-program \
    --l2 https://optimism-goerli.infura.io/v3/{APIKEY} \
    --l1 https://goerli.infura.io/v3/{APIKEY} \
    --l1.trustrpc \
    --l1.rpckind debug_geth \
    --log.format terminal \
    --l2.head 0xedc79de4d616a9100fdd42192224580daee81ea3d6303de8089d48a6c1bf4816 \
    --network goerli \
    --l1.head 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab \
    --l2.claim 0x530658ab1b1b3ff4829731fc8d5955f0e6b8410db2cd65b572067ba58df1f2b9 \
    --l2.blocknumber 8813570 \
    --datadir /tmp/fpp-database \
    --server

# Add --proof-at '=12345' (or pick other pattern, see --help)
# to pick a step to build a proof for (e.g. exact step, every N steps, etc.)

# Also see `./bin/cannon run --help` for more options

```

## Compiling

To build op-program, from within the `op-program` directory run:

```shell
make op-program
```

This resulting executable will be in `./bin/op-program`

## Testing

To run op-program unit tests, from within the `op-program` directory run:

```shell
make test
```

## Lint

To run the linter, from within the `op-program` directory run:
```shell
make lint
```

This requires having `golangci-lint` installed.

## Running

From within the `op-program` directory, options can be reviewed with:

```shell
./bin/op-program --help
```
