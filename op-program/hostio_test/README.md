# Steps to Run

0. Compile go following [link](../README.md).  Make sure activate the compiled go
1. `cd op-program`
2. `make hostio-test`, which will generate `hostio.wasi` in bin folder
3. `node runtime/wasi/wasi_exec_node.js bin/hostio-test.wasi`.  You may see a stack overflow error, which means that the code runs correctly (but we lack `proc_exit` support at the moment.) 
