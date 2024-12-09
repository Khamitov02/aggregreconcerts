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
		r.Get("/api/v1/getMusic/{user_id}", h.getUserMusic)
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

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getUserMusic(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	
	// Make HTTP request to external service
	resp, err := http.Get("http://user-service/api/v1/getMusic/" + userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userBands []string
	if err := json.NewDecoder(resp.Body).Decode(&userBands); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process user bands and save recommendations
	if err := h.service.ProcessUserBands(r.Context(), userID, userBands); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
