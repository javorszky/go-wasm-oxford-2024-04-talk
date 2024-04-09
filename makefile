.PHONY: run gengo genrust

run:
	go run main.go

gengo:
	tinygo build -o pkg/wasm/testdata/add/add.wasm -target=wasi pkg/wasm/testdata/add/add.go
	tinygo build -o pkg/wasm/testdata/date/date.wasm -target=wasi pkg/wasm/testdata/date/date.go
	tinygo build -o pkg/wasm/testdata/hash/hash.wasm -target=wasi pkg/wasm/testdata/hash/hash.go

genrust:
	cargo build --manifest-path pkg/wasm/testdata/greet/Cargo.toml --release --target wasm32-wasi
