# Go + Wasm

All of these examples are adapted from the wazero examples. They have all been tested prior to giving the talk.

There are three `make` commands that you might find useful:

1. `make run`: it starts the go application which has the embedded wasm modules. This will fail if the wasm files are not present, so you should probably run the other two make commands first
2. `make gengo`: generates the .wasm files from the tinygo projects of add, date, and hash in the pkg/wasm/testdata directory
3. `make genrust`: generates the .wasm file from the rust package in pkg/wasm/testdata/greet directory

## Prerequisites

You'll need the following toolchains for everything to work:
* Go itself. I'm using 1.21.0 darwin/arm64
* TinyGo - this one is needed because wasm only supports a subset of what Go can do, which tinygo also accounted for. I'm using tinygo version 0.31.2 darwin/arm64 (using go version go1.21.0 and LLVM version 17.0.1)
* Rust, cargo, and the rustup toolchain - I'm using rustup 1.27.0 (bbb9276d2 2024-03-08)
* Rust wasm32-wasi target. Install this with `rustup target add wasm32-wasi`

## Legal bit

All of the code is open source, and it's been adapted from either the wazero examples, or using echo itself.

## Resources

* https://webassembly.org/docs/security/
* https://wasi.dev/
* https://tinygo.org/docs/guides/webassembly/
* https://rustwasm.github.io/book/
* https://github.com/tetratelabs/wazero
* https://wasmtime.dev/
