package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/javorszky/go-wasm-talk/pkg/wasm"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
)

func Home() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func Add(r *wasm.Runner) echo.HandlerFunc {
	return func(c echo.Context) error {
		x, y := c.Param("x"), c.Param("y")

		xu, err := strconv.ParseUint(x, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse x")
		}

		yu, err := strconv.ParseUint(y, 10, 64)
		if err != nil {
			return errors.Wrap(err, "failed to parse y")
		}

		wasmCtx := context.WithValue(context.Background(), wasm.ContextKey, []uint64{xu, yu})

		result, err := r.Exec(wasmCtx, wasm.Add)
		if err != nil {
			return errors.Wrap(err, "error while executing wasm module")
		}

		return c.String(http.StatusOK, string(result))
	}
}

func Date(r *wasm.Runner) echo.HandlerFunc {
	return func(c echo.Context) error {
		wasmCtx := context.Background()

		result, err := r.Exec(wasmCtx, wasm.Date)
		if err != nil {
			return errors.Wrap(err, "error while executing wasm module")
		}

		return c.String(http.StatusOK, string(result))
	}
}

func NoDate(r *wasm.Runner) echo.HandlerFunc {
	return func(c echo.Context) error {
		wasmCtx := context.Background()

		result, err := r.Exec(wasmCtx, wasm.NoDate)
		if err != nil {
			return errors.Wrap(err, "error while executing wasm module")
		}

		return c.String(http.StatusOK, string(result))
	}
}

func Hash(r *wasm.Runner) echo.HandlerFunc {
	return func(c echo.Context) error {
		wasmCtx := context.WithValue(context.Background(), wasm.ContextKey, c.Param("what"))

		result, err := r.Exec(wasmCtx, wasm.Hash)
		if err != nil {
			return errors.Wrap(err, "error while executing wasm module")
		}

		return c.String(http.StatusOK, string(result))
	}
}

func Greet(r *wasm.Runner) echo.HandlerFunc {
	return func(c echo.Context) error {
		wasmCtx := context.WithValue(context.Background(), wasm.ContextKey, c.Param("who"))

		result, err := r.Exec(wasmCtx, wasm.Greet)
		if err != nil {
			return errors.Wrap(err, "error while executing wasm module")
		}

		return c.String(http.StatusOK, string(result))
	}
}
