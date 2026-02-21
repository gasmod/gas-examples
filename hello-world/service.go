package main

import (
	"net/http"

	"github.com/gasmod/gas"
)

type Service struct {
	router *gas.Router
}

func NewService(router *gas.Router) *Service {
	return &Service{router: router}
}

func (s *Service) Name() string {
	return "hello-world-service"
}

func (s *Service) Init() error {
	s.router.Handle(s.Name(), http.MethodGet, "/", s.handleIndex)
	return nil
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) handleIndex(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("Hello, world!")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
