package concerts

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Handler struct {
	router  *chi.Mux
	service Service
}

func NewHandler(router *chi.Mux, service Service) *Handler {
	return &Handler{
		router:  router,
		service: service,
	}
}

func (h *Handler) Register() {
	h.router.Group(func(r chi.Router) {
		r.Get("/api/v1/getRecommendations", h.getRecommendations)
		r.Post("/api/v1/putConcerts", h.putConcerts)
	})
}

func (h *Handler) getRecommendations(w http.ResponseWriter, r *http.Request) {
	recommendations, err := h.service.GetUserRecommendations(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": recommendations,
	})
}

func (h *Handler) putConcerts(w http.ResponseWriter, r *http.Request) {
	var input BandsInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.SaveBands(r.Context(), input.Bands); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// After saving bands, process recommendations for all users
	if err := h.service.ProcessAllUserBands(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
