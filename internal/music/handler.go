package music

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	router       *chi.Mux
	service      Service
	concerts     []Concert
	recommendations map[string][]Concert
}

type ConcertRequest struct {
	Bands []Concert `json:"bands"`
}

func NewHandler(router *chi.Mux, service Service) *Handler {
	return &Handler{
		router:  router,
		service: service,
		concerts: []Concert{},
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

	h.concerts = req.Bands
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
		http.Error(w, "Failed to parse user bands", http.StatusInternalServerError)
		return
	}

	h.recommendations = make(map[string][]Concert)
	for userID, bands := range userBands {
		for _, band := range bands {
			for _, concert := range h.concerts {
				if strings.Contains(concert.BandName, band) {
					h.recommendations[userID] = append(h.recommendations[userID], concert)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.recommendations)
}
