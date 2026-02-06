package service

import (
	"music-service/internal/repository/postgres"
)

type PreferenceService struct {
	repo *postgres.PreferenceRepository
}

func NewPreferenceService(r *postgres.PreferenceRepository) *PreferenceService {
	return &PreferenceService{repo: r}
}

func (s *PreferenceService) Save(userID string, genres map[string]int) error {
	for name, score := range genres {
		err := s.repo.UpsertUserGenreScore(userID, name, score)
		if err != nil {
			return err
		}
	}

	return nil
}
