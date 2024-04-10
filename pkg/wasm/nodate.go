package wasm

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/pkg/errors"
)

func (r *Runner) execNoDate(ctx context.Context) ([]byte, error) {
	// Instantiate the guest Wasm into the same runtime. It exports the `today`
	// function, implemented in WebAssembly.
	mod, err := r.r.Instantiate(ctx, dateWasm)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate module: %v", err)
	}

	// Call the `add` function and print the results to the console.
	today := mod.ExportedFunction("today")

	// need to free the memory that the module returned the result in.
	free := mod.ExportedFunction("free")

	// we don't need to pass anything to it, it does not depend on input.
	ptrSize, err := today.Call(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to call today with context")
	}

	// These were encoded inside the module source.
	datePtr := uint32(ptrSize[0] >> 32)
	dateSize := uint32(ptrSize[0])

	// This pointer is managed by TinyGo, but TinyGo is unaware of external usage.
	// So, we have to free it when finished
	if datePtr != 0 {
		defer func() {
			_, err = free.Call(ctx, uint64(datePtr))
			if err != nil {
				log.Panicln(err)
			}
		}()
	}

	// The pointer is a linear memory offset, which is where we write the name.
	bytes, ok := mod.Memory().Read(datePtr, dateSize)
	if !ok {
		return nil, fmt.Errorf("failed to read memory")
	}

	return bytes, nil
}
