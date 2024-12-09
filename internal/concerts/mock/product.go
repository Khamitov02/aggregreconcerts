package mock

import (
	"aggregconcerts/internal/concerts"
	"context"

	"github.com/google/uuid"
)

type Mock struct{}

func NewFridge() *Mock {
	return &Mock{}
}

func (m *Mock) Products(ctx context.Context) ([]concerts.Product, error) {
	products := []concerts.Product{
		{
			ID:    uuid.New().String(),
			Name:  "Test name",
			Count: 17,
		},
	}

	return products, nil
}

func (m *Mock) Place(ctx context.Context, product concerts.Product) (id string, err error) {
	return uuid.New().String(), nil
}
