package api

import (
	"kursy-kriptovalyut/internal/usecases"
)

type Server struct {
	service *usecases.Service
}

func NewServer(service usecases.Service) *Server {
	return &Server{&service}
}

func (serv *Server) GetLastRates() {}

func (serv *Server) GetAggregateRates() {}

// GetLastRates - ???
// GetAggregateRates - ???
