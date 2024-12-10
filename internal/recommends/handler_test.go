package recommends_test

import (
	"bytes"
	"encoding/json"
	"musicadviser/internal/recommends"
	"musicadviser/internal/recommends/mock"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestHandler_getProducts(t *testing.T) {
	service := mock.NewFridge()
	router := chi.NewRouter()

	h := recommends.NewHandler(router, service)

	h.Register()

	req, err := http.NewRequest(http.MethodGet, "/api/v1/products", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	t.Run("status", func(t *testing.T) {
		if rr.Code != http.StatusOK {
			t.Errorf("handler return wrong status code: want %d, got: %s", http.StatusOK, rr.Code)
		}
	})

	t.Run("body", func(t *testing.T) {
		var got recommends.Product
		err := json.NewDecoder(rr.Body).Decode(&got)
		if err != nil {
			t.Fatal(err)
		}

		want := recommends.Product{
			// заполнить данными из мока
		}

		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("GET /api/v1/products mismatch: (-want +got)\n%s", diff)
		}
	})
}

func TestHandler_putConcerts(t *testing.T) {
	// Setup
	router := chi.NewRouter()
	service := mock.NewService()
	handler := recommends.NewHandler(router, service)
	handler.Register()

	tests := []struct {
		name       string
		input      recommends.ConcertRequest
		wantStatus int
	}{
		{
			name: "valid request",
			input: recommends.ConcertRequest{
				Bands: []recommends.Concert{
					{BandName: "Test Band", Link: "http://example.com"},
				},
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "empty request",
			input:      recommends.ConcertRequest{},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/putConcerts", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
		})
	}
}

func TestHandler_getRecommendations(t *testing.T) {
	// Setup mock server for external music service
	mockMusicServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userBands := map[string][]string{
			"user1": {"Band1", "Band2"},
			"user2": {"Band3"},
		}
		json.NewEncoder(w).Encode(userBands)
	}))
	defer mockMusicServer.Close()

	// Setup main handler
	router := chi.NewRouter()
	service := mock.NewService()
	handler := recommends.NewHandler(router, service)
	handler.Register()

	// Add some test concerts
	concertsReq := recommends.ConcertRequest{
		Bands: []recommends.Concert{
			{BandName: "Band1 Live", Link: "http://example.com/band1"},
			{BandName: "Band2 Concert", Link: "http://example.com/band2"},
		},
	}

	// First, put some concerts
	body, _ := json.Marshal(concertsReq)
	putReq := httptest.NewRequest(http.MethodPost, "/api/v1/putConcerts", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, putReq)

	// Now test recommendations
	req := httptest.NewRequest(http.MethodGet, "/api/v1/getRecommendations", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var recommendations map[string][]recommends.Concert
	err := json.NewDecoder(rr.Body).Decode(&recommendations)
	assert.NoError(t, err)

	// Verify recommendations structure
	assert.Contains(t, recommendations, "user1")
	assert.Contains(t, recommendations, "user2")
}
