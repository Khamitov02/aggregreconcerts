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

type MusicRequest struct {
	UserID string   `json:"user_id"`
	Bands  []string `json:"bands"`
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
		r.Get("/api/v1/products", h.getProducts)
		r.Post("/api/v1/putMusic", h.putMusic)
		r.Post("/api/v1/putConcerts", h.putConcerts)
		r.Get("/api/v1/getRecommendations", h.getRecommendations)
	})
}

func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler: GetProducts accessed - UserAgent: %s, RemoteAddr: %s", r.UserAgent(), r.RemoteAddr)
	// validate r
	data, err := h.service.Products(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "", data)
}

func (h *Handler) putMusic(w http.ResponseWriter, r *http.Request) {
	log.Printf("Handler: PutMusic accessed - UserAgent: %s, RemoteAddr: %s", r.UserAgent(), r.RemoteAddr)
	var req MusicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Process each band in the request
	for _, bandName := range req.Bands {
		product := Product{
			UserID:   req.UserID,
			BandName: bandName,
		}
		
		_, err := h.service.Place(r.Context(), product)
		if err != nil {
			// If the band already exists, continue to the next one
			if strings.Contains(err.Error(), "already exists") {
				continue
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Music bands added successfully")
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
