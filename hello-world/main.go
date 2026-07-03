// The simplest possible Gas application. No services, no config — just a
// single inline DI-aware handler registered directly on the app router.
package main

import (
	"log"
	"net/http"

	"github.com/gasmod/gas"
)

func main() {
	// NewApp creates a Router and EventBus internally.
	// CSRF protection is enabled by default.
	app := gas.NewApp()

	// Routes can be registered directly on the router without a service.
	// The first argument is the owning service name (empty here since there
	// is no service). The handler is a DI-aware function: its first parameter
	// must be gas.Context and it must return error.
	app.Router().Handle("", http.MethodGet, "/", func(ctx gas.Context) error {
		return ctx.Text(http.StatusOK, "Hello, World!")
	})

	// Run starts the full lifecycle: init services → run migrations →
	// execute ready hooks → start HTTP server. On SIGINT/SIGTERM it
	// gracefully shuts down.
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
