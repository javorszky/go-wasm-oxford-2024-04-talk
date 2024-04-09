package wasm

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

//go:embed testdata/greet/target/wasm32-wasi/release/greet.wasm
var greetWasm []byte

func (r *Runner) execGreet(ctx context.Context) ([]byte, error) {
	// Instantiate the guest Wasm into the same runtime. It exports the `today`
	// function, implemented in WebAssembly.
	mod, err := r.r.Instantiate(ctx, greetWasm)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %v", err)
	}

	params := ctx.Value(ContextKey)
	if params == nil {
		return nil, fmt.Errorf("failed to instantiate module due to empty context")
	}

	paramsParsed, ok := params.(string)
	if !ok {
		return nil, fmt.Errorf("failed to instantiate module due to bad context type")
	}

	// Call the `greeting` to return the greeting from the wasm module.
	greeting := mod.ExportedFunction("greeting")

	// this is needed so we can write our input into memory for wasm.
	allocate := mod.ExportedFunction("allocate")

	// need to free the memory that the module returned the result in.
	deallocate := mod.ExportedFunction("deallocate")

	// Let's use the argument to this main function in Wasm.
	paramSize := uint64(len(paramsParsed))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := allocate.Call(ctx, paramSize)
	if err != nil {
		log.Panicln(err)
	}
	paramPtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer deallocate.Call(ctx, paramPtr, paramSize)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(uint32(paramPtr), []byte(paramsParsed)) {
		return nil, errors.New("failed to write input memory")
	}

	ptrSize, err := greeting.Call(ctx, paramPtr, paramSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call today with context")

	}
	greetingPtr := uint32(ptrSize[0] >> 32)
	greetingSize := uint32(ptrSize[0])

	// This pointer was allocated by Rust, but owned by Go, So, we have to
	// deallocate it when finished
	defer deallocate.Call(ctx, uint64(greetingPtr), uint64(greetingSize))

	// The pointer is a linear memory offset, which is where we write the name.
	bytes, ok := mod.Memory().Read(greetingPtr, greetingSize)
	if !ok {
		return nil, fmt.Errorf("failed to read memory")
	}

	ret := string(bytes)

	// The bottom return produces a corrupted memory state :)
	return []byte(ret), nil
	// return bytes, nil
}
