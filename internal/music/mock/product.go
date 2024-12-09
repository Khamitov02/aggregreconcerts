package mock

import (
	"context"
	"musicadviser/internal/music"

	"github.com/google/uuid"
)

type Mock struct{}

func NewFridge() *Mock {
	return &Mock{}
}

func (m *Mock) Products(ctx context.Context) ([]music.Product, error) {
	products := []music.Product{
		{
			ID:    uuid.New().String(),
			Name:  "Test name",
			Count: 17,
		},
	}

	return products, nil
}

func (m *Mock) Place(ctx context.Context, product music.Product) (id string, err error) {
	return uuid.New().String(), nil
}
