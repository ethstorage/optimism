# op-program-zkwasm

## build op-program-host
> Notice: require go version<=1.20.7 (you can use gvm to change go version)
```
cd op-program
make op-program-host
```

## build op-program-client and replay (Generate witness by replaying op-program-client)
> Notice: require go version<=1.20.7

```
make op-program-client

# specify the `--preimage {file path}` flag to change preimages file location
alias replay="./bin/op-program --l2 http://65.108.75.40:8645     --l1 http://65.108.75.40:8745     --l1.trustrpc     --l1.rpckind debug_geth     --log.format terminal     --l2.head 0xedc79de4d616a9100fdd42192224580daee81ea3d6303de8089d48a6c1bf4816     --network goerli     --l1.head 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab     --l2.claim 0x530658ab1b1b3ff4829731fc8d5955f0e6b8410db2cd65b572067ba58df1f2b9     --l2.blocknumber 8813570     --datadir /tmp/fpp-database --preimage ./bin/preimages.bin    --exec ./bin/op-program-client"

replay
```

## Build op-program-client-wasi for zkWasm image

### build customized zkwasm-go
```
go clean -cache
git clone -b wasi https://github.com/ethstorage/go
cd go/src
./all.bash
```

### build op-program-client-wasi with zkwasm-go
```
# make sure your go path is the above zkwasm-go(you can change the relevant go path in makefile)
make op-program-client-wasi
```

### check witness file with Node.js zkwasm emulator
> Notice: require node.js version>=20.5.1
```
node ./zkWasm-emulator/wasi/wasi_exec_node.js ./bin/op-program-client.wasi ./bin/preimages.bin
```
> Notice: it will print `wasm_output:1024` if correct.

## zkWasm emulator

### build zkWasm
```
git clone -b dev https://github.com/ethstorage/zkWasm
cd zkWasm
git submodule update --init
cargo build --release
```

### dry run
```
{/target/release/delphinus-cli path} -k 22 --function zkmain --output ./output --wasm {op-program-client-preimage.wasi path} dry-run --preimages "{preimages.bin path}"
```
> Notice: it will print `wasm_output:1024` if correct. Dry run costs about 22 hours.


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