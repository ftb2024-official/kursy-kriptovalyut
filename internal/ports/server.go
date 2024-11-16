package ports

import "kursy-kriptovalyut/internal/cases"

type Server struct {
	service *cases.Service
}

func NewServer(service cases.Service) *Server {
	return &Server{&service}
}

func (serv *Server) GetLastRates() {}

func (serv *Server) GetAggregateRates() {}

// GetLastRates - ???
// GetAggregateRates - ???
