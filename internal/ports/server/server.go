package server

import (
	"context"
	"cryptocurrency/internal/entity"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Service interface {
	GetLatestPrices(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMinPrices(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetMaxPrices(ctx context.Context, titles []string) ([]entity.Coin, error)
	GetPriceChangePercent(ctx context.Context, titles []string) ([]entity.Coin, error)
}
type Server struct {
	service Service
	router  *chi.Mux
	server  *http.Server
}

func NewServer(service Service, address string) (*Server, error) {
	if service == nil {
		return nil, errors.Wrap(entity.ErrInvalidParams, "http server: service is nil")
	}
	if address == "" {
		return nil, errors.Wrap(entity.ErrInvalidParams, "http server: address is empty")
	}
	router := chi.NewRouter()
	server := &Server{
		service: service,
		router:  router,
	}
	server.routes()
	server.server = &http.Server{
		Addr:    address,
		Handler: router,
	}
	return server, nil
}

func (s *Server) routes() {
	s.router.Get("/coins/prices", s.GetLatestPrices)
	s.router.Get("/coins/prices/min", s.GetMinPrices)
	s.router.Get("/coins/prices/max", s.GetMaxPrices)
	s.router.Get("/coins/prices/change-percent", s.GetPriceChangePercent)
}

func (s *Server) GetLatestPrices(w http.ResponseWriter, r *http.Request) {
	titles := r.URL.Query()["title"]
	if len(titles) == 0 {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	coins, err := s.service.GetLatestPrices(r.Context(), titles)
	if err != nil {
		Error(w, err)
		return
	}
	response := make([]CoinDTO, len(coins))
	for i, coin := range coins {
		response[i] = NewCoinDTO(coin)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

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
	response := make([]CoinDTO, len(coins))
	for i, coin := range coins {
		response[i] = NewCoinDTO(coin)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

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
	response := make([]CoinDTO, len(coins))
	for i, coin := range coins {
		response[i] = NewCoinDTO(coin)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

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
	response := make([]CoinDTO, len(coins))
	for i, coin := range coins {
		response[i] = NewCoinDTO(coin)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func Error(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entity.ErrInvalidParams):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, entity.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) Run() error {
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "http server: ListenAndServe")
	}
	return nil
}
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "http server: Shutdown")
	}
	return nil
}
