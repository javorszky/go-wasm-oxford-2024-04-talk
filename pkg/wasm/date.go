package wasm

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/tetratelabs/wazero"
)

//go:embed testdata/date/date.wasm
var dateWasm []byte

func (r *Runner) execDate(ctx context.Context) ([]byte, error) {
	// We need to specifically allow the module to have access to the system walltime (ie "what is the current time?")
	// If we don't, the returned time will always be 2022-01-01 00:00:00
	modConfig := wazero.NewModuleConfig().WithSysWalltime()

	// Instantiate the guest Wasm into the same runtime. It exports the `today`
	// function, implemented in WebAssembly.
	mod, err := r.r.InstantiateWithConfig(ctx, dateWasm, modConfig)
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
