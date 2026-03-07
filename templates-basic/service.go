package main

import (
	"net/http"
	"time"

	"github.com/gasmod/gas"
)

type Module struct {
	router *gas.Router
}

func NewModule(router *gas.Router) *Module {
	return &Module{router: router}
}

func (m *Module) Name() string {
	return "landing-service"
}

func (m *Module) Init() error {
	m.router.Handle(m.Name(), http.MethodGet, "/", m.handleIndex)
	m.router.NotFound(m.Name(), m.handleNotFound)

	return nil
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) handleIndex(ctx gas.Context, ui gas.UIProvider) error {
	return ui.Render(ctx.ResponseWriter(), "home", map[string]any{
		"Year": time.Now().Year(),
	})
}

func (m *Module) handleNotFound(ctx gas.Context, ui gas.UIProvider) error {
	return ui.RenderWithStatus(ctx.ResponseWriter(), http.StatusNotFound, "404", map[string]any{
		"Title": "Not Found",
		"Year":  time.Now().Year(),
	})
}
