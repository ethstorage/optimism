# op-program-wasm

## build js-wasm and replay
> require go version>=1.21.0 (you can use gvm to change go version)

### build js-wasm
```
cd op-program
make op-program-client-wasm
```
### replay js-wasm
```
alias replay="./bin/op-program --l2 http://65.108.75.40:8645     --l1 http://65.108.75.40:8745     --l1.trustrpc     --l1.rpckind debug_geth     --log.format terminal     --l2.head 0xedc79de4d616a9100fdd42192224580daee81ea3d6303de8089d48a6c1bf4816     --network goerli     --l1.head 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab     --l2.claim 0x530658ab1b1b3ff4829731fc8d5955f0e6b8410db2cd65b572067ba58df1f2b9     --l2.blocknumber 8813570     --datadir /tmp/fpp-database     --exec \"node ./runtime/js/wasm_exec_node.js ./bin/op-program-client.wasm\""

replay
```

## build wasi and replay
> require go version>=1.21.0 (you can use gvm to change go version)

### build wasi
```
cd op-program
make op-program-client-wasi
```
### replay wasi
```
alias replay="./bin/op-program --l2 http://65.108.75.40:8645     --l1 http://65.108.75.40:8745     --l1.trustrpc     --l1.rpckind debug_geth     --log.format terminal     --l2.head 0xedc79de4d616a9100fdd42192224580daee81ea3d6303de8089d48a6c1bf4816     --network goerli     --l1.head 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab     --l2.claim 0x530658ab1b1b3ff4829731fc8d5955f0e6b8410db2cd65b572067ba58df1f2b9     --l2.blocknumber 8813570     --datadir /tmp/fpp-database     --exec \"node ./runtime/wasi/wasi_exec_node.js ./bin/op-program-client.wasi\""

replay
```


## build wasi and replay without op-host program
> require go version>=1.21.0 (you can use gvm to change go version)

### build wasi
```
cd op-program
make op-program-client-wasi
```
### replay wasi with op-host program for dumping preimages json file(./bin/preimages.json)
```
alias replay="./bin/op-program --l2 http://65.108.75.40:8645     --l1 http://65.108.75.40:8745     --l1.trustrpc     --l1.rpckind debug_geth     --log.format terminal     --l2.head 0xedc79de4d616a9100fdd42192224580daee81ea3d6303de8089d48a6c1bf4816     --network goerli     --l1.head 0x204f815790ca3bb43526ad60ebcc64784ec809bdc3550e82b54a0172f981efab     --l2.claim 0x530658ab1b1b3ff4829731fc8d5955f0e6b8410db2cd65b572067ba58df1f2b9     --l2.blocknumber 8813570     --datadir /tmp/fpp-database     --exec \"node ./runtime/wasi/wasi_exec_node.js ./bin/op-program-client.wasi\""
replay
```

### replay **without op-host program**
```
node ./runtime/wasi-json/wasi_json_node.js  ./bin/op-program-client.wasi
```

## problems with op-wasm and solutions
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
- `/go-ethereum/trie`:
    - `replace github.com/ethereum/go-ethereum => ./op-geth`
    - ref arb's go-etherum, `replace github.com/ethereum/go-ethereum => ./go-ethereum`, [related pr](https://github.com/OffchainLabs/go-ethereum/pull/205)
- `/go-ethereum/trie` & `VictoriaMetrics/fastcache`: use [arb's fastcache](https://github.com/OffchainLabs/fastcache) `replace github.com/VictoriaMetrics/fastcache => ./fastcache`

## Arb reference
- [replay code](https://github.com/OffchainLabs/nitro/blob/master/cmd/replay/main.go)
- Build:
    1. follow [arbitrum build nitro tutorial](https://docs.arbitrum.io/node-running/how-tos/build-nitro-locally) to install requirements
    2. `make build-wasm-bin`
- Goetherum related commit:
    - [First commit](https://github.com/OffchainLabs/go-ethereum/commits?after=1319d385dc35f0a3be7166ec4a63ce83de89c376+244&author=PlasmaPower)
    - [Another offchainlabs](https://github.com/OffchainLabs/go-ethereum/commits?author=Tristan-Wilson&before=1319d385dc35f0a3be7166ec4a63ce83de89c376+70)


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
