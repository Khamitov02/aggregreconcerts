package concerts

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type AppService struct {
	store Store
}

func NewAppService(s Store) *AppService {
	return &AppService{
		store: s,
	}
}

func (s *AppService) SaveBands(ctx context.Context, bands []Band) error {
	return s.store.SaveBands(ctx, bands)
}

func (s *AppService) GetUserRecommendations(ctx context.Context) ([]UserBands, error) {
	return s.store.GetAllUserRecommendations(ctx)
}

func (s *AppService) ProcessUserBands(ctx context.Context, userID string, userBands []string) error {
	// Get all available bands from memory
	availableBands, err := s.store.GetAllBands(ctx)
	if err != nil {
		return err
	}

	// Find matching bands
	var matchingBands []Band
	for _, availableBand := range availableBands {
		for _, userBand := range userBands {
			if strings.Contains(strings.ToLower(availableBand.BandName), strings.ToLower(userBand)) {
				matchingBands = append(matchingBands, availableBand)
				break
			}
		}
	}

	// Save recommendations for this user
	return s.store.SaveUserRecommendations(ctx, userID, matchingBands)
}
