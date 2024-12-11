package ports

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"kursy-kriptovalyut/internal/entities"
	"kursy-kriptovalyut/pkg/dto"
	"kursy-kriptovalyut/pkg/logger"
)

type Server struct {
	service Service
	server  *chi.Mux
}

func NewServer(service Service) (*Server, error) {
	if service == nil || service == Service(nil) {
		return nil, errors.Wrap(entities.ErrInvalidParam, "service not set")
	}

	s := &Server{
		service: service,
		server:  chi.NewRouter(),
	}

	s.routes()
	return s, nil
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	s.server.ServeHTTP(rw, r)
}

func (s *Server) routes() {
	s.server.Get("/rates/last", s.GetLastRates)
	s.server.Get("/rates/agg", s.GetAggregateRates)
}

var log = logger.GetLogger()

func (s *Server) GetLastRates(rw http.ResponseWriter, r *http.Request) {
	titlesQueryParam := r.URL.Query().Get("titles")
	if titlesQueryParam == "" {
		log.Info("1")
		respondWithJSON(rw, http.StatusBadRequest, map[string]string{"error": "missing 'titles' query parameter"})
		return
	}

	coinTitles := strings.Split(strings.ToUpper(strings.TrimSpace(titlesQueryParam)), ",")

	coins, err := s.service.GetLastRates(r.Context(), coinTitles)
	if err != nil {
		log.Info("2")
		if errors.Is(err, entities.ErrNotFound) {
			log.Info("3")
			respondWithJSON(rw, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		log.Info("4")
		respondWithJSON(rw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := make([]dto.CoinDTO, 0, len(coins))
	for _, coin := range coins {
		response = append(response, dto.CoinDTO{Title: coin.Title, Price: coin.Price})
	}

	respondWithJSON(rw, http.StatusOK, response)
}

func (s *Server) GetAggregateRates(rw http.ResponseWriter, r *http.Request) {
	titlesQueryParam := r.URL.Query().Get("titles")
	aggFuncQueryParam := r.URL.Query().Get("aggFunc")
	if titlesQueryParam == "" || aggFuncQueryParam == "" {
		log.Info("6")
		respondWithJSON(rw, http.StatusBadRequest, map[string]string{"error": "missing 'titles' or 'aggFunc' query parameters"})
		return
	}

	coinTitles := strings.Split(strings.ToUpper(titlesQueryParam), ",")
	aggFuncName := strings.ToUpper(aggFuncQueryParam)

	coins, err := s.service.GetAggRates(r.Context(), coinTitles, aggFuncName)

	if err != nil {
		log.Info("7")
		if errors.Is(err, entities.ErrNotFound) {
			log.Info("8")
			respondWithJSON(rw, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}

		log.Info("9")
		respondWithJSON(rw, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	response := make([]dto.CoinDTO, 0, len(coins))
	for _, coin := range coins {
		response = append(response, dto.CoinDTO{Title: coin.Title, Price: coin.Price})
	}

	respondWithJSON(rw, http.StatusOK, response)
}

func respondWithJSON(rw http.ResponseWriter, status int, data interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	if err := json.NewEncoder(rw).Encode(data); err != nil {
		log.Error("failed to encode response")
		http.Error(rw, `{"error": "internal server error"}`, http.StatusInternalServerError)
	}
}
