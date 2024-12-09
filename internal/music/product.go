package music

import (
	"context"
)

type Product struct {
	ID        string
	UserID    string
	BandName  string
}

type Concert struct {
	BandName string
	Link     string
}

type Service interface {
	Products(ctx context.Context) ([]Product, error)
	Place(ctx context.Context, product Product) (id string, err error)
}

type Store interface {
	LoadProducts(ctx context.Context) ([]Product, error)
	SaveProduct(ctx context.Context, product Product) (id string, err error)
}
