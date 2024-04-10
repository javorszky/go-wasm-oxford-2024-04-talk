package wasm

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type Mod int

const (
	Add Mod = iota
	Date
	NoDate
	Hash
	Greet
)

type Runner struct {
	ctx context.Context
	r   wazero.Runtime
}

// CTXKey is a private type alias so it can't be overwritten from outside of the package.
type ctxKey string

// ContextKey is used to store input data from an http request that will be passed onto the wasm modules.
const ContextKey ctxKey = "wazero.ctx.key"

func NewRunner(ctx context.Context) *Runner {
	r := wazero.NewRuntime(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	runner := Runner{
		ctx: ctx,
		r:   r,
	}

	return &runner
}

// Exec is the entry point into the actual individual methods. In production this would be abstracted away so every
// wasm module would have a "take input from memory - write output to memory" pattern with the same exported function
// names that users would move inside of.
func (r *Runner) Exec(ctx context.Context, what Mod) ([]byte, error) {
	switch what {
	case Add:
		return r.execAdd(ctx)
	case Date:
		return r.execDate(ctx)
	case NoDate:
		return r.execNoDate(ctx)
	case Hash:
		return r.execHash(ctx)
	case Greet:
		return r.execGreet(ctx)
	default:
		return nil, fmt.Errorf("tried to call a module that doesn't exist")
	}
}

// Stop is needed so we can defer call this from main.
func (r *Runner) Stop() {
	_ = r.r.Close(r.ctx)
}
