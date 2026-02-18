package main

import (
	"fmt"
	"net/http"

	"github.com/gasmod/gas"
)

type Module struct {
	router *gas.Router
}

type Option func(*Module)

func WithRouter(r *gas.Router) Option {
	return func(m *Module) {
		m.router = r
	}
}

func NewModule(opts ...Option) *Module {
	m := &Module{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *Module) Name() string {
	return "hello-world-module"
}

func (m *Module) Init() error {
	if m.router == nil {
		return fmt.Errorf("gas: router not set")
	}
	if err := m.router.Handle(m.Name(), http.MethodGet, "/", m.handleIndex); err != nil {
		return err
	}
	return nil
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) handleIndex(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("Hello, world!")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
