package ports

import (
	"encoding/json"
	"net/http"
	"strings"

	_ "kursy-kriptovalyut/docs"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger"

	"kursy-kriptovalyut/internal/entities"
	"kursy-kriptovalyut/pkg/dto"
	"kursy-kriptovalyut/pkg/logger"
)

var log = logger.NewLogger()

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

// @title Chi Swagger Example
// @version 1.0
// @host localhost:8080
func (s *Server) routes() {
	s.server.Get("/rates/last", s.GetLastRates)
	s.server.Get("/rates/agg", s.GetAggregateRates)

	s.server.Mount("/swagger", httpSwagger.WrapHandler)
	// url := httpSwagger.URL("http://localhost:8080/swagger/doc.json")
	// s.server.Get("/swagger/*any", httpSwagger.WrapHandler(swaggerFiles.Handler, url))
}

// @Summary Get last rates
// @Description Get the latest rates for specified coins
// @Tags rates
// @Produce json
// @Param titles query string true "Comma-separated list of coin titles" example(BTC,ETH)
// @Success 200 {array} dto.CoinDTO
// @Failure 400 {object} dto.ErrRespDTO
// @Failure 500 {object} dto.ErrRespDTO
// @Router /rates/last [get]
func (s *Server) GetLastRates(rw http.ResponseWriter, r *http.Request) {
	log.Info("(server.GetLastRates)")
	titlesQueryParam := r.URL.Query().Get("titles")
	if titlesQueryParam == "" {
		log.Warn("(server.GetLastRates) missing 'titles' query parameter")
		respondWithJSON(rw, http.StatusBadRequest, dto.ErrRespDTO{
			StatusCode: http.StatusBadRequest,
			Msg:        "missing 'titles' query parameter",
		})
		return
	}

	coinTitles := strings.Split(strings.ToUpper(strings.TrimSpace(titlesQueryParam)), ",")

	log.Info("(server.service.GetLastRates)")
	coins, err := s.service.GetLastRates(r.Context(), coinTitles)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			respondWithJSON(rw, http.StatusNotFound, dto.ErrRespDTO{
				StatusCode: http.StatusNotFound,
				Msg:        err.Error(),
			})
			return
		}

		respondWithJSON(rw, http.StatusInternalServerError, dto.ErrRespDTO{
			StatusCode: http.StatusInternalServerError,
			Msg:        err.Error(),
		})
		return
	}

	response := make([]dto.CoinDTO, 0, len(coins))
	for _, coin := range coins {
		response = append(response, dto.CoinDTO{Title: coin.Title, Price: coin.Price})
	}

	respondWithJSON(rw, http.StatusOK, response)
}

// @Summary Get aggregated rates
// @Description Get aggregated rates for specified coins using an aggregation function
// @Tags rates
// @Produce json
// @Param titles query string true "Comma-separated list of coin titles" example(BTC,ETH)
// @Param aggFunc query string true "Aggregation function (MAX, MIN, AVG)" example(MAX)
// @Success 200 {array} dto.CoinDTO
// @Failure 400 {object} dto.ErrRespDTO
// @Failure 500 {object} dto.ErrRespDTO
// @Router /rates/agg [get]
func (s *Server) GetAggregateRates(rw http.ResponseWriter, r *http.Request) {
	log.Info("(server.GetAggregateRates)")
	titlesQueryParam := r.URL.Query().Get("titles")
	aggFuncQueryParam := r.URL.Query().Get("aggFunc")
	if titlesQueryParam == "" || aggFuncQueryParam == "" {
		log.Warn("(server.GetAggregateRates) missing 'titles' or 'aggFunc' query parameters")
		respondWithJSON(rw, http.StatusBadRequest, dto.ErrRespDTO{
			StatusCode: http.StatusBadRequest,
			Msg:        "missing 'titles' or 'aggFunc' query parameters"})
		return
	}

	coinTitles := strings.Split(strings.ToUpper(titlesQueryParam), ",")
	aggFuncName := strings.ToUpper(aggFuncQueryParam)

	coins, err := s.service.GetAggRates(r.Context(), coinTitles, aggFuncName)

	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			respondWithJSON(rw, http.StatusNotFound, dto.ErrRespDTO{
				StatusCode: http.StatusNotFound,
				Msg:        err.Error(),
			})
			return
		}

		respondWithJSON(rw, http.StatusBadRequest, dto.ErrRespDTO{
			StatusCode: http.StatusBadRequest,
			Msg:        err.Error(),
		})
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
		log.Error("(respondWithJSON) failed to encode response")
		http.Error(rw, `{"error": "internal server error"}`, http.StatusInternalServerError)
	}
}
