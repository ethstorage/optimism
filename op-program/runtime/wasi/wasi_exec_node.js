import { readFile } from 'node:fs/promises';
import { WASI } from 'wasi';
import { argv, env } from 'node:process';
import fs from "fs"
import path from 'node:path';

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
    await readFile( new URL(path.join(process.cwd(),process.argv[2]), import.meta.url))
  );

  let instance
  let wasm_input_counter = 0;
  let wasm_input_state = 0;
  let wasm_inputs = ["68656c6c6f", "48656c6c6f2c20576f726c6421"]; // hex of input data


  const hostio = {
    "env":
    {
      "wasm_input": (isPublic) => {
        // a simple case to return "hello" and "Hello, World!"
        if (wasm_input_state == 0) {
          wasm_input_state = 1;
          console.log(BigInt(wasm_inputs[wasm_input_counter].length / 2))
          return BigInt(wasm_inputs[wasm_input_counter].length / 2);
        }

        let start = (wasm_input_state - 1) * 16;
        let end = wasm_input_state * 16;
        let data;
        if (end >= wasm_inputs[wasm_input_counter].length) {
          end = wasm_inputs[wasm_input_counter].length;
          data = BigInt("0x" + wasm_inputs[wasm_input_counter].substring(start, end));
          // move to next data
          wasm_input_state = 0;
          wasm_input_counter = wasm_input_counter + 1;
        } else {
          wasm_input_state = wasm_input_state + 1;
          data = BigInt("0x" + wasm_inputs[wasm_input_counter].substring(start, end));
        }

        if (max_mem < instance.exports.memory.buffer.byteLength) {
          max_mem = instance.exports.memory.buffer.byteLength
        }

        console.log(data);
        return data;
      },

      "wasm_exit": (code) => {
        console.log("\nmaximum memory usage==========================>",max_mem);
        process.exit(code);
      },

      "require": (cond) => {
        if (cond == 0) {
          console.log("require is not satisfied, which is a false assertion in the wasm code. Please check the logic of your image or input.");
          process.exit(1);
        }
      }
    },
    "_gotest": //func get_preimage_len
    {
      get_preimage_len: (keyPtr) => {
        let key = new Uint8Array(instance.exports.memory.buffer,keyPtr, 32)
        // console.log("key is:", key.toString())

        //read preimage from file descriptor
        let PClientRFd = 5
        let PClientWFd = 6
        fs.writeSync(PClientWFd, Buffer.from(key))

        //write to go-wasm
        let lenBuf = Buffer.alloc(8)
        fs.readSync(PClientRFd,lenBuf,0,8)
        // console.log("lenBuf====>",lenBuf)
        let len = parseInt(lenBuf.toString("hex"),16)
        // console.log("len js:", len)
        return len
      },

      //func getKeyFromOracle() []byte
      get_preimage_from_oracle: (keyPtr,offset,len) => {
        let mem = new DataView(instance.exports.memory.buffer)
        let PClientRFd = 5
        let key = new Uint8Array(instance.exports.memory.buffer,keyPtr, 32)
        // console.log("key is:", key.toString())

        let data = Buffer.alloc(len)
        let readed_len = fs.readSync(PClientRFd,data)
        // console.log("read length",readed_len)
        // console.log("read data:",  data.subarray(0,32))

        //send data back to go-wasm
        for(let i=0; i< readed_len; i++){
          mem.setUint8(offset,data[i],true)
          offset = offset + 1
        }
        return readed_len

      },

      "hint_oracle": (retBufPtr, retBufSize) => {
        //load hintstr
        let hintArr = new Uint8Array(instance.exports.memory.buffer,retBufPtr, retBufSize)
        let HClientWFd = 4
        fs.writeSync(HClientWFd, Buffer.from(hintArr))
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

