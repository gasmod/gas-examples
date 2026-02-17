package main

import (
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
	m.router.Handle(m.Name(), http.MethodGet, "/", m.handleIndex)
	return nil
}

func (m *Module) Close() error {
	return nil
}

func (m *Module) handleIndex(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, world!"))
}
