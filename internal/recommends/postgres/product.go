package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"aggregconcerts/internal/recommends"

	"github.com/jmoiron/sqlx"
)

type Product struct {
	ID       sql.NullString
	UserID   sql.NullString
	BandName sql.NullString
}

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) LoadProducts(ctx context.Context) ([]recommends.Product, error) {
	var dbProducts []Product
	query := "SELECT ID, user_id, band_name FROM music_bands"

	err := s.db.SelectContext(ctx, &dbProducts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to select products: %w", err)
	}

	var products []recommends.Product
	for _, dbProduct := range dbProducts {
		if !dbProduct.ID.Valid || !dbProduct.UserID.Valid || !dbProduct.BandName.Valid {
			return nil, fmt.Errorf("one of the required fields is NULL")
		}

		product := recommends.Product{
			ID:       dbProduct.ID.String,
			UserID:   dbProduct.UserID.String,
			BandName: dbProduct.BandName.String,
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *Storage) SaveProduct(ctx context.Context, product recommends.Product) (id string, err error) {
	// First check if this band already exists for this user
	var exists bool
	err = s.db.QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM music_bands WHERE user_id = $1 AND band_name = $2)",
		product.UserID, product.BandName).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("failed to check band existence: %w", err)
	}
	if exists {
		return "", fmt.Errorf("band %s already exists for user %s", product.BandName, product.UserID)
	}

	// If not exists, insert new record
	query := `
		INSERT INTO music_bands (user_id, band_name)
		VALUES ($1, $2)
		RETURNING id
	`

	var newID string
	err = s.db.QueryRowContext(ctx, query, product.UserID, product.BandName).Scan(&newID)
	if err != nil {
		return "", fmt.Errorf("failed to insert band: %w", err)
	}

	return newID, nil
}
