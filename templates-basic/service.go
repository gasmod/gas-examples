package main

import (
	"net/http"
	"time"

	"github.com/gasmod/gas"
)

type Service struct {
	router *gas.Router
	ui     gas.UIProvider
}

func NewService(router *gas.Router, ui gas.UIProvider) *Service {
	return &Service{router: router, ui: ui}
}

func (m *Service) Name() string {
	return "landing-service"
}

func (m *Service) Init() error {
	m.router.Handle(m.Name(), http.MethodGet, "/", m.handleIndex)
	m.router.NotFound(m.Name(), m.handleNotFound)

	return nil
}

func (m *Service) Close() error {
	return nil
}

func (m *Service) handleIndex(w http.ResponseWriter, _ *http.Request) {
	_ = m.ui.Render(w, "home", map[string]any{
		"Year": time.Now().Year(),
	})
}

func (m *Service) handleNotFound(w http.ResponseWriter, _ *http.Request) {
	_ = m.ui.RenderWithStatus(w, http.StatusNotFound, "404", map[string]any{
		"Title": "Not Found",
		"Year":  time.Now().Year(),
	})
}
