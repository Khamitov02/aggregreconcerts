package recommends

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Handler struct {
	router         *chi.Mux
	service        Service
	concerts       []Concert
	recommendations map[string][]Concert
	mu             sync.Mutex
}

type Concert struct {
	BandName string `json:"band_name"`
	Link     string `json:"link"`
}

type ConcertRequest struct {
	Bands []Concert `json:"bands"`
}

func NewHandler(router *chi.Mux, service Service) *Handler {
	return &Handler{
		router:         router,
		service:        service,
		recommendations: make(map[string][]Concert),
	}
}

func (h *Handler) Register() {
	h.router.Group(func(r chi.Router) {
		r.Post("/api/v1/putConcerts", h.putConcerts)
		r.Get("/api/v1/getRecommendations", h.getRecommendations)
	})
}

func (h *Handler) putConcerts(w http.ResponseWriter, r *http.Request) {
	var req ConcertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.concerts = req.Bands
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Concerts added successfully")
}

func (h *Handler) getRecommendations(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:3434/api/v1/getMusic")
	if err != nil {
		http.Error(w, "Failed to fetch user bands", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var userBands map[string][]string
	if err := json.Unmarshal(body, &userBands); err != nil {
		http.Error(w, "Invalid JSON format from getMusic", http.StatusInternalServerError)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for userID, bands := range userBands {
		for _, userBand := range bands {
			for _, concert := range h.concerts {
				if strings.Contains(concert.BandName, userBand) {
					h.recommendations[userID] = append(h.recommendations[userID], concert)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.recommendations); err != nil {
		http.Error(w, "Failed to encode recommendations", http.StatusInternalServerError)
	}
}
