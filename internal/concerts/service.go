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

func (s *AppService) ProcessAllUserBands(ctx context.Context) error {
	// Get all available bands from memory
	availableBands, err := s.store.GetAllBands(ctx)
	if err != nil {
		return err
	}

	// Make HTTP request to get all users' music
	resp, err := http.Get("http://localhost:3434/api/v1/getMusic")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var userMusic UserMusicResponse
	if err := json.NewDecoder(resp.Body).Decode(&userMusic); err != nil {
		return err
	}

	// Process each user's bands
	for userID, userBands := range userMusic {
		var matchingBands []Band
		for _, availableBand := range availableBands {
			for _, userBand := range userBands {
				if strings.Contains(
					strings.ToLower(availableBand.BandName),
					strings.ToLower(userBand),
				) {
					matchingBands = append(matchingBands, availableBand)
					break
				}
			}
		}

		// Save recommendations for this user
		if err := s.store.SaveUserRecommendations(ctx, userID, matchingBands); err != nil {
			return err
		}
	}

	return nil
}
