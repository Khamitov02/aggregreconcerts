package recommends_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"musicadviser/internal/recommends"
	"musicadviser/internal/recommends/memory"
)

func TestAppService_Products(t *testing.T) {
	// Setup
	store := memory.NewStore()
	service := recommends.NewAppService(store)
	ctx := context.Background()

	// Test data
	testProducts := []recommends.Product{
		{
			ID:       "1",
			UserID:   "user1",
			BandName: "Test Band 1",
		},
		{
			ID:       "2",
			UserID:   "user2",
			BandName: "Test Band 2",
		},
	}

	// Store test data
	for _, p := range testProducts {
		_, err := store.SaveProduct(ctx, p)
		assert.NoError(t, err)
	}

	// Test retrieval
	products, err := service.Products(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, len(testProducts))

	// Verify contents
	for i, p := range products {
		assert.Equal(t, testProducts[i].BandName, p.BandName)
		assert.Equal(t, testProducts[i].UserID, p.UserID)
	}
}

func TestAppService_Place(t *testing.T) {
	// Setup
	store := memory.NewStore()
	service := recommends.NewAppService(store)
	ctx := context.Background()

	// Test data
	product := recommends.Product{
		UserID:   "user1",
		BandName: "Test Band",
	}

	// Test placing product
	id, err := service.Place(ctx, product)
	assert.NoError(t, err)
	assert.NotEmpty(t, id)

	// Verify product was stored
	products, err := service.Products(ctx)
	assert.NoError(t, err)
	assert.Len(t, products, 1)
	assert.Equal(t, product.BandName, products[0].BandName)
	assert.Equal(t, product.UserID, products[0].UserID)
}
