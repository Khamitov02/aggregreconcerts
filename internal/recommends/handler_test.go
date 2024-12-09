package recommends_test

import (
	"encoding/json"
	"musicadviser/internal/recommends"
	"musicadviser/internal/recommends/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-cmp/cmp"
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
