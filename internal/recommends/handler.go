package recommends

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	router  *chi.Mux
	service Service
}

type MusicRequest struct {
	UserID string   `json:"user_id"`
	Bands  []string `json:"bands"`
}

func NewHandler(router *chi.Mux, service Service) *Handler {
	return &Handler{
		router:  router,
		service: service,
	}
}

func (h *Handler) Register() {
	h.router.Group(func(r chi.Router) {
		r.Get("/api/v1/products", h.getProducts)
		r.Post("/api/v1/putMusic", h.putMusic)
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
