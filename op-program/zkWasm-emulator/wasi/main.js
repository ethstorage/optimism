import { readFile } from 'node:fs/promises';
import { WASI } from 'wasi';
import { argv, env } from 'node:process';
import fs from "fs"

(async function () {
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
    await readFile(new URL(process.argv[2], import.meta.url)),
  );

  console.log("start ")

  let instance
  let cur = 0
  let preimages = fs.readFileSync("./bin/preimages.bin")
  const hostio = {
    env: {
      wasm_input: (ispulic) => {
				let data = preimages.readBigInt64BE(cur)
        console.log("data:",data)
				cur += 8
				return data
      },
      require: (cond) => {
        if (cond == 0) {
          console.log("require is not satisfied, which is a false assertion in the wasm code. Please check the logic of your image or input.");
          process.exit(1);
        }
      },
      wasm_output:(value) => {
        console.log("successfully run:", value)
      }
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

