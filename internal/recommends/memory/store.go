package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"musicadviser/internal/recommends"
)

type Store struct {
	mu       sync.RWMutex
	products map[string]recommends.Product
}

func NewStore() *Store {
	return &Store{
		products: make(map[string]recommends.Product),
	}
}

func (s *Store) LoadProducts(_ context.Context) ([]recommends.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]recommends.Product, 0, len(s.products))
	for _, p := range s.products {
		products = append(products, p)
	}
	return products, nil
}

func (s *Store) SaveProduct(_ context.Context, product recommends.Product) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if product.ID == "" {
		product.ID = uuid.New().String()
	}
	s.products[product.ID] = product
	return product.ID, nil
} 