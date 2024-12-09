package recommends

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
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
	log.Printf("Received PUT request to /api/v1/putConcerts")
	
	var req ConcertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	
	log.Printf("Received concerts data: %+v", req.Bands)

	h.mu.Lock()
	h.concerts = req.Bands
	h.mu.Unlock()

	log.Printf("Successfully stored %d concerts", len(req.Bands))
	
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Concerts added successfully")
}

func (h *Handler) getRecommendations(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received GET request to /api/v1/getRecommendations")
	
	resp, err := http.Get("http://localhost:3434/api/v1/getMusic")
	if err != nil {
		log.Printf("Error fetching user bands: %v", err)
		http.Error(w, "Failed to fetch user bands", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	log.Printf("Received user bands data from external service: %s", string(body))

	var userBands map[string][]string
	if err := json.Unmarshal(body, &userBands); err != nil {
		log.Printf("Error unmarshaling user bands: %v", err)
		http.Error(w, "Invalid JSON format from getMusic", http.StatusInternalServerError)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Clear previous recommendations
	h.recommendations = make(map[string][]Concert)

	for userID, bands := range userBands {
		log.Printf("Processing recommendations for user %s with bands: %v", userID, bands)
		for _, userBand := range bands {
			for _, concert := range h.concerts {
				if strings.Contains(concert.BandName, userBand) {
					log.Printf("Found match: User band '%s' matches concert '%s'", userBand, concert.BandName)
					h.recommendations[userID] = append(h.recommendations[userID], concert)
				}
			}
		}
	}

	log.Printf("Generated recommendations: %+v", h.recommendations)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(h.recommendations); err != nil {
		log.Printf("Error encoding recommendations: %v", err)
		http.Error(w, "Failed to encode recommendations", http.StatusInternalServerError)
		return
	}
	
	log.Printf("Successfully sent recommendations response")
}
