package concerts

import (
	"context"
)

type Band struct {
	BandName string `json:"band_name"`
	Link     string `json:"link"`
}

type UserBands struct {
	UserID string
	Bands  []Band
}

type BandsInput struct {
	Bands []Band `json:"bands"`
}

type Service interface {
	SaveBands(ctx context.Context, bands []Band) error
	GetUserRecommendations(ctx context.Context) ([]UserBands, error)
	ProcessUserBands(ctx context.Context, userID string, userBands []string) error
}

type Store interface {
	SaveBands(ctx context.Context, bands []Band) error
	GetAllBands(ctx context.Context) ([]Band, error)
	SaveUserRecommendations(ctx context.Context, userID string, bands []Band) error
	GetAllUserRecommendations(ctx context.Context) ([]UserBands, error)
}
