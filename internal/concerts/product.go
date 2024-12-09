package concerts

import (
	"context"
)

type UserMusicResponse map[string][]string // userID -> []bandName

type Band struct {
	BandName string `json:"band_name"`
	Link     string `json:"link"`
}

type UserBands struct {
	UserID string `json:"user_id"`
	Bands  []Band `json:"bands"`
}

type BandsInput struct {
	Bands []Band `json:"bands"`
}

type Service interface {
	SaveBands(ctx context.Context, bands []Band) error
	GetUserRecommendations(ctx context.Context) ([]UserBands, error)
	ProcessAllUserBands(ctx context.Context) error
}

type Store interface {
	SaveBands(ctx context.Context, bands []Band) error
	GetAllBands(ctx context.Context) ([]Band, error)
	SaveUserRecommendations(ctx context.Context, userID string, bands []Band) error
	GetAllUserRecommendations(ctx context.Context) ([]UserBands, error)
}
