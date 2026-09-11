package public

import (
	"context"
	"cryptocurrency/internal/entity"
	"cryptocurrency/pkg/dto"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	service Service
	router  *chi.Mux
	server  *http.Server
}

func NewServer(service Service, address string) (*Server, error) {
	if service == nil || service == Service(nil) {
		return nil, errors.Wrap(entity.ErrInvalidParams, "http: service is nil")
	}
	router := chi.NewRouter()
	server := &Server{
		service: service,
		router:  router,
	}
	server.routes()
	server.server = &http.Server{
		Handler: router,
		Addr:    address,
	}
	return server, nil
}

func (s *Server) routes() {
	s.router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	s.router.Route("/api/v1/coins", func(r chi.Router) {
		r.Get("/get/latest", s.GetLatestPrices)
		r.Get("/get/min", s.GetMinPrices)
		r.Get("/get/max", s.GetMaxPrices)
		r.Get("/get/change-percent", s.GetPriceChangePercent)
	})
}

// GetLatestPrices
// @Summary Get latest prices
// @Description Returns the latest prices for requested cryptocurrencies
// @Tags coins
// @Produce json
// @Param title query []string true "Coin titles"
// @Success 200 {array} dto.CoinDTO
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /coins/get/latest [get]
func (s *Server) GetLatestPrices(rw http.ResponseWriter, req *http.Request) {
	titles := req.URL.Query()["title"]
	if len(titles) == 0 {
		http.Error(rw, "title is required", http.StatusBadRequest)
		return
	}
	coins, err := s.service.GetLatestPrices(req.Context(), titles)
	if err != nil {
		Error(rw, err)
		return
	}
	response := make(dto.CoinsDTO, len(coins))
	for i, coin := range coins {
		response[i] = dto.CoinDTO(coin)
	}
	rw.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

// GetMinPrices
// @Summary Get minimum prices
// @Description Returns minimum cryptocurrency prices for the last 24 hours
// @Tags coins
// @Produce json
// @Param title query []string true "Coin titles"
// @Success 200 {array} dto.CoinDTO
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /coins/get/min [get]
func (s *Server) GetMinPrices(w http.ResponseWriter, r *http.Request) {
	titles := r.URL.Query()["title"]
	if len(titles) == 0 {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	coins, err := s.service.GetMinPrices(r.Context(), titles)
	if err != nil {
		Error(w, err)
		return
	}
	response := make(dto.CoinsDTO, len(coins))
	for i, coin := range coins {
		response[i] = dto.CoinDTO(coin)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetMaxPrices
// @Summary Get maximum prices
// @Description Returns maximum cryptocurrency prices for the last 24 hours
// @Tags coins
// @Produce json
// @Param title query []string true "Coin titles"
// @Success 200 {array} dto.CoinDTO
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /coins/get/max [get]
func (s *Server) GetMaxPrices(w http.ResponseWriter, r *http.Request) {
	titles := r.URL.Query()["title"]
	if len(titles) == 0 {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	coins, err := s.service.GetMaxPrices(r.Context(), titles)
	if err != nil {
		Error(w, err)
		return
	}
	response := make(dto.CoinsDTO, len(coins))
	for i, coin := range coins {
		response[i] = dto.CoinDTO(coin)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPriceChangePercent
// @Summary Get cryptocurrency price change percent
// @Description Returns cryptocurrency price change percentage for the last hour
// @Tags coins
// @Produce json
// @Param title query []string true "Coin titles"
// @Success 200 {array} dto.CoinDTO
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /coins/get/change-percent [get]
func (s *Server) GetPriceChangePercent(w http.ResponseWriter, r *http.Request) {
	titles := r.URL.Query()["title"]
	if len(titles) == 0 {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	coins, err := s.service.GetPriceChangePercent(r.Context(), titles)
	if err != nil {
		Error(w, err)
		return
	}
	response := make(dto.CoinsDTO, len(coins))
	for i, coin := range coins {
		response[i] = dto.CoinDTO(coin)
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Error(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrInvalidParams):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, entity.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "http http: ListenAndServe")
	}
	return nil
}
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "http http: Shutdown")
	}
	return nil
}
