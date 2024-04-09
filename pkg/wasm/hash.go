package wasm

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

//go:embed testdata/hash/hash.wasm
var hashWasm []byte

func (r *Runner) execHash(ctx context.Context) ([]byte, error) {
	// We need to specifically allow the module to have access to the system walltime (ie "what is the current time?")
	// modConfig := wazero.NewModuleConfig().WithSysWalltime()

	// Instantiate the guest Wasm into the same runtime. It exports the `today`
	// function, implemented in WebAssembly.
	mod, err := r.r.Instantiate(ctx, hashWasm)
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

	// Call the `hash` function and print the results to the console.
	hash := mod.ExportedFunction("hash")

	// this is needed so we can write our input into memory for wasm.
	malloc := mod.ExportedFunction("malloc")

	// need to free the memory that the module returned the result in.
	free := mod.ExportedFunction("free")

	// Let's use the argument to this main function in Wasm.
	paramSize := uint64(len(paramsParsed))

	// Instead of an arbitrary memory offset, use TinyGo's allocator. Notice
	// there is nothing string-specific in this allocation function. The same
	// function could be used to pass binary serialized data to Wasm.
	results, err := malloc.Call(ctx, paramSize)
	if err != nil {
		log.Panicln(err)
	}
	paramPtr := results[0]
	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	defer free.Call(ctx, paramPtr)

	// The pointer is a linear memory offset, which is where we write the name.
	if !mod.Memory().Write(uint32(paramPtr), []byte(paramsParsed)) {
		return nil, errors.New("failed to write input memory")
	}

	result, err := hash.Call(ctx, paramPtr, paramSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call today with context")
	}

	// These were encoded inside the module source.
	resultPtr := uint32(result[0] >> 32)
	resultSize := uint32(result[0])

	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	if resultPtr != 0 {
		defer func() {
			_, err = free.Call(ctx, uint64(resultPtr))
			if err != nil {
				log.Panicln(err)
			}
		}()
	}

	// The pointer is a linear memory offset, which is where we write the name.
	bytes, ok := mod.Memory().Read(resultPtr, resultSize)
	if !ok {
		return nil, fmt.Errorf("failed to read memory")
	}

	return bytes, nil
}
