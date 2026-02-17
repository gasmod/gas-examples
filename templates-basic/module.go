package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gasmod/gas"
)

type Module struct {
	router *gas.Router
	ui     gas.UIProvider
}

type Option func(*Module)

func WithRouter(r *gas.Router) Option {
	return func(m *Module) { m.router = r }
}

func WithUIProvider(ui gas.UIProvider) Option {
	return func(m *Module) { m.ui = ui }
}

func NewModule(opts ...Option) *Module {
	m := &Module{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Module) Name() string {
	return "landing-module"
}

func (m *Module) Init() error {
	if m.ui == nil {
		return fmt.Errorf("%s: UIProvider is required", m.Name())
	}

	if err := m.router.Handle(m.Name(), http.MethodGet, "/", m.handleIndex); err != nil {
		return err
	}

	m.router.Mux().NotFound(m.handleNotFound)

	return nil
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) handleIndex(w http.ResponseWriter, _ *http.Request) {
	if err := m.ui.Render(w, "home", map[string]any{
		"Year": time.Now().Year(),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (m *Module) handleNotFound(w http.ResponseWriter, _ *http.Request) {
	if err := m.ui.RenderWithStatus(w, http.StatusNotFound, "404", map[string]any{
		"Title": "Not Found",
		"Year":  time.Now().Year(),
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
