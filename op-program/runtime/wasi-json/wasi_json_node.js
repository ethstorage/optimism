import { readFile } from 'node:fs/promises';
import { WASI } from 'wasi';
import { argv, env } from 'node:process';
import fs from "fs"
import path from 'node:path';

(async function () {
  // parse preimages json file
  let preimages_data = fs.readFileSync("./bin/preimages.json")
  let preimages_json = JSON.parse(preimages_data)
  let preimages = {}
  for(let key in preimages_json) {
    let data_buf = Buffer.from(preimages_json[key],'hex')
    preimages[key] = {
      len: data_buf.length,
      buf: data_buf
    }
  }

  const wasi = new WASI({
    version: 'preview1',
    args: argv,
    env,
    returnOnExit: true
    // preopens: {
    //   '/sandbox': '/root/now/wasm-runtime',
    // },
  });
  const wasm = await WebAssembly.compile(
    await readFile( new URL(path.join(process.cwd(),process.argv[2]), import.meta.url))
  );

  console.log("start ")

  let instance

  const hostio = {
    "_gotest": //func get_preimage_len
    {
      get_preimage_len: (keyPtr) => {
        let key = new Uint8Array(instance.exports.memory.buffer,keyPtr, 32)
        key = Buffer.from(key).toString("hex")
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

  const wasi_imports = wasi.getImportObject()
  const previous_clock_time_get = wasi_imports.wasi_snapshot_preview1.clock_time_get
  wasi_imports.wasi_snapshot_preview1.clock_time_get =  (clockId, precision, time) => {
    let res = previous_clock_time_get(clockId, precision, time)
    if (max_mem < instance.exports.memory.buffer.byteLength) {
      max_mem = instance.exports.memory.buffer.byteLength
    }
    return res
  }

  instance = await WebAssembly.instantiate(wasm, {...wasi_imports,...hostio});
  wasi.start(instance);

  process.on("exit", ()=> {
    console.log("\nmaximum memory usage==========================>",max_mem)
  })

})()

