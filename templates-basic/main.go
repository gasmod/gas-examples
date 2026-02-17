package main

import (
	"log"
	"net/http"

	"github.com/gasmod/gas"
	config "github.com/gasmod/gas-config"
	ui "github.com/gasmod/gas-ui"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := gas.DefaultConfig()

	cfgMod := config.New()
	if err := cfgMod.Init(); err != nil {
		log.Fatal(err)
	}

	reg := gas.NewMiddlewareRegistry()
	router := gas.NewRouter(reg)
	bus := gas.NewEventBus()

	uiMod := ui.New(
		ui.WithConfig(&ui.Config{
			TemplateDir: "templates",
			StaticDir:   "static",
			StaticPath:  "/static/",
			LayoutName:  "base",
		}),
		ui.WithRouter(router),
	)

	router.Mux().Use(middleware.Logger)
	router.Mux().Use(middleware.Recoverer)
	router.Mux().Use(middleware.Compress(5))
	router.Mux().Use(securityHeaders)
	router.Mux().Use(cacheControl)

	app := gas.NewApp(
		gas.WithConfig(cfg),
		gas.WithRouter(router),
		gas.WithMiddlewareRegistry(reg),
		gas.WithEventBus(bus),
		gas.WithModule(cfgMod),
		gas.WithModule(uiMod),
		gas.WithModule(NewModule(
			WithRouter(router),
			WithUIProvider(uiMod),
		)),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

func cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.URL.Path) > 8 && r.URL.Path[:8] == "/static/" {
			w.Header().Set("Cache-Control", "public, max-age=2592000, immutable")
		}
		next.ServeHTTP(w, r)
	})
}
