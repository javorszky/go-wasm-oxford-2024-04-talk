package main

import (
	"context"

	"github.com/javorszky/go-wasm-talk/pkg/handlers"
	"github.com/javorszky/go-wasm-talk/pkg/wasm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Create a new runner for our webassembly shenanigans
	r := wasm.NewRunner(context.Background())
	defer r.Stop()

	// Handle http requests.
	e := echo.New()

	e.Use(handlers.PanicRecovery())
	e.Use(middleware.Logger())

	e.GET("/", handlers.Home())
	e.GET("/add/:x/:y", handlers.Add(r))
	e.GET("/date", handlers.Date(r))
	e.GET("/nodate", handlers.NoDate(r))

	e.GET("/hash/:what", handlers.Hash(r))
	e.GET("/greet/:who", handlers.Greet(r))

	e.Logger.Fatal(e.Start(":1323"))
}
