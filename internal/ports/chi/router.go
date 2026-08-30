package chi

import (
	"cryptocurrency/internal/cases"
	"cryptocurrency/internal/entity"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Router struct {
	service *cases.CoinService
}

func NewRouter(service *cases.CoinService) (*Router, error) {
	if service == nil {
		return nil, errors.Wrap(entity.ErrInvalidParams, "router: service is nil")
	}
	return &Router{
		service: service,
	}, nil
}

func (r *Router) Routes() chi.Router {
	router := chi.NewRouter()

	router.Get("/coins/latest", r.GetLatestPrices)
	router.Get("/coins/min", r.GetMinPrices)
	router.Get("/coins/max", r.GetMaxPrices)
	router.Get("/coins/percent", r.GetChangePercent)
	return router
}
func (r *Router) GetLatestPrices(w http.ResponseWriter, req *http.Request) {
	titles := req.URL.Query()["title"]

	coins, err := r.service.GetLatestPrices(req.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(coins); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (r *Router) GetMinPrices(w http.ResponseWriter, req *http.Request) {
	titles := req.URL.Query()["title"]

	coins, err := r.service.GetMinPrices(req.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(coins); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *Router) GetMaxPrices(w http.ResponseWriter, req *http.Request) {
	titles := req.URL.Query()["title"]

	coins, err := r.service.GetMaxPrices(req.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(coins); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *Router) GetChangePercent(w http.ResponseWriter, req *http.Request) {
	titles := req.URL.Query()["title"]

	coins, err := r.service.GetPriceChangePercent(req.Context(), titles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(coins); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
