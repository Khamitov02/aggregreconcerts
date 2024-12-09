package memory

import (
	"aggregconcerts/internal/concerts"
	"context"
	"strings"
	"sync"
)

type Storage struct {
	mu              sync.RWMutex
	bands           []concerts.Band
	recommendations map[string][]concerts.Band // userID -> recommended bands
}

func NewStorage() *Storage {
	return &Storage{
		bands:           make([]concerts.Band, 0),
		recommendations: make(map[string][]concerts.Band),
	}
}

func (s *Storage) SaveBands(ctx context.Context, bands []concerts.Band) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bands = bands
	return nil
}

func (s *Storage) GetAllBands(ctx context.Context) ([]concerts.Band, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.bands, nil
}

func (s *Storage) SaveUserRecommendations(ctx context.Context, userID string, bands []concerts.Band) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recommendations[userID] = bands
	return nil
}

func (s *Storage) GetAllUserRecommendations(ctx context.Context) ([]concerts.UserBands, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := make([]concerts.UserBands, 0, len(s.recommendations))
	for userID, bands := range s.recommendations {
		result = append(result, concerts.UserBands{
			UserID: userID,
			Bands:  bands,
		})
	}
	return result, nil
}
