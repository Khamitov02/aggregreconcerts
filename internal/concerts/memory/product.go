package memory

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"musicadviser/internal/concerts"
	"sync"
)

type Storage struct {
	mu       sync.RWMutex
	id       int
	products map[string]fridge.Product
}

// init Storage
func NewStorage() *Storage {
	return &Storage{
		products: make(map[string]fridge.Product),
		id:       1,
	}
}

// get all from in-memory store
func (s *Storage) LoadProducts(ctx context.Context) ([]fridge.Product, error) {
	//read lock
	s.mu.RLock()
	defer s.mu.RUnlock()

	products := make([]fridge.Product, 0, len(s.products))
	for _, product := range s.products {
		products = append(products, product)
	}
	return products, nil
}

// SaveProduct adds a new product to the in-memory store and returns its ID
func (s *Storage) SaveProduct(ctx context.Context, product fridge.Product) (string, error) {
	//rw lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate a new unique ID
	//or  s.generateID2()
	product.ID = s.generateID()

	// Check for duplicates
	if _, exists := s.products[product.ID]; exists {
		return "", errors.New("duplicate ID")
	}

	s.products[product.ID] = product
	return product.ID, nil
}

// incremental ID
func (s *Storage) generateID() string {
	id := s.id
	s.id++
	return fmt.Sprintf("%d", id)
}

// UUID generated randomly, could be the same
func (s *Storage) generateID2() string {
	return uuid.New().String() // Generate a new UUID
}
