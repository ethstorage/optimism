import { readFile } from 'node:fs/promises';
import { WASI } from 'wasi';
import { argv, env } from 'node:process';
import fs from "fs"
import path from 'node:path';

(async function () {
  // parse preimages json file
  let preimages_data = fs.readFileSync("./bin/preimages.json")
  console.log("parseing......",preimages_data.length)
  let preimages_json = JSON.parse(preimages_data)
  let preimages = {}
  for(let key in preimages_json) {
    let data_buf = Buffer.from(preimages_json[key],'hex')
    preimages[key] = {
      len: data_buf.length,
      buf: data_buf
    }
  }

  console.log("parsed......")


  const wasi = new WASI({
    version: 'preview1',
    args: [],
    env:{},
    returnOnExit: true
    // preopens: {
    //   '/sandbox': '/root/now/wasm-runtime',
    // },
  });
  const wasm = await WebAssembly.compile(
    await readFile( new URL(path.join(process.cwd(),process.argv[2]), import.meta.url))
  );

  console.log("starting replay using ./bin/preimages.json...")

  let instance

  const hostio = {
    "_gotest": //func get_preimage_len
    {
      get_preimage_len: (keyPtr) => {
        let key = new Uint8Array(instance.exports.memory.buffer,keyPtr, 32)
        key = Buffer.from(key).toString("hex")
        console.log("key is===>", key)
        //read preimage len from json file
        return preimages[key].len
      },

      //func getKeyFromOracle() []byte
      get_preimage_from_oracle: (keyPtr,offset,len) => {
        let mem = new DataView(instance.exports.memory.buffer)

        let key = new Uint8Array(instance.exports.memory.buffer,keyPtr, 32)
        key = Buffer.from(key).toString("hex")
        let data = preimages[key].buf
        //send data back to wasm
        for(let i=0; i< len; i++){
          mem.setUint8(offset,data[i],true)
          offset = offset + 1
        }
        return len
      },

      "hint_oracle": (retBufPtr, retBufSize) => {
        //do nothing, cause we have load all preimages
        return
      },
    }
  }
  let max_mem = 0

  let CPUTIME_START = BigInt(0)

  const WASI_CLOCK_MONOTONIC = 0;
  const WASI_CLOCK_PROCESS_CPUTIME_ID = 1;
  const WASI_CLOCK_REALTIME = 2;
  const WASI_CLOCK_THREAD_CPUTIME_ID = 3;

  const now = (clockId) => {
    switch (clockId) {
      case WASI_CLOCK_MONOTONIC:
        return CPUTIME_START;
      case WASI_CLOCK_REALTIME:
        return undefined
      case WASI_CLOCK_PROCESS_CPUTIME_ID:
      case WASI_CLOCK_THREAD_CPUTIME_ID:
        CPUTIME_START += 1n
        return CPUTIME_START;
      default:
        return null;
    }
  };

  const WASI_EINVAL = 28;
  const WASI_ESUCCESS = 0;
  const wasi_imports = {
    wasi_snapshot_preview1: {
      ...wasi.getImportObject().wasi_snapshot_preview1,
      ...{
        sched_yield: ()=>{},
        proc_exit: (rval)=>{
          console.log("\nmaximum memory usage==========================>",max_mem)
          process.exit(rval);
        },
        args_get: (_argv, argvBuf) => {
          let coffset = _argv;
          let offset = argvBuf;
          let view = new DataView(instance.exports.memory.buffer);
          argv.forEach((a) => {
            view.setUint32(coffset, offset, true);
            coffset += 4;
            offset += Buffer.from(instance.exports.memory.buffer).write(`${a}\0`, offset);
          });
          console.log("argv================>", argv)
          return WASI_ESUCCESS;
        },
        args_sizes_get: (argc, argvBufSize)=>{
          let view = new DataView(instance.exports.memory.buffer);
          view.setUint32(argc, argv.length, true);
          const size = argv.reduce((acc, a) => acc + Buffer.byteLength(a) + 1, 0);
          view.setUint32(argvBufSize, size, true);
          return WASI_ESUCCESS;
        },
        clock_time_get: (clockId, precision, time) => {
          const n = now(clockId);
          if (n === null) {
            return WASI_EINVAL;
          }
          let view = new DataView(instance.exports.memory.buffer);
          view.setBigUint64(time, n, true);

          //for debug memory ony
          if (max_mem < instance.exports.memory.buffer.byteLength) {
            max_mem = instance.exports.memory.buffer.byteLength
          }
        },
        environ_get: (environ, environBuf) => {
          let coffset = environ;
          let offset = environBuf;
          const cache = Buffer.from(instance.exports.memory.buffer);
          let view = new DataView(instance.exports.memory.buffer);
          Object.entries(env)
            .forEach(([key, value]) => {
              view.setUint32(coffset, offset, true);
              coffset += 4;
              offset += cache.write(`${key}=${value}\0`, offset);
            });
          return WASI_ESUCCESS;
        },
        environ_sizes_get: (environCount, environBufSize) => {
          let view = new DataView(instance.exports.memory.buffer);
          view.setUint32(environCount, 0, true);
          view.setUint32(environBufSize, 0, true);
        },

        fd_write: (fd, iovs, iovsLen, offset, nwritten) => {
          let view = new DataView(instance.exports.memory.buffer);
          view.setUint32(nwritten, 0, true);
        },
        random_get: ()=>{},
        poll_oneoff: ()=>{},
        fd_close: (fd) => {
          return
        },
        fd_fdstat_get: ()=>{},
        fd_fdstat_set_flags: ()=>{},
        fd_prestat_get: (fd, bufPtr)=>{
          return 8
        },
        fd_prestat_dir_name: ()=>{}
      }
    }
  }

  instance = await WebAssembly.instantiate(wasm, {...wasi_imports,...hostio});
  wasi.start(instance);

})()

