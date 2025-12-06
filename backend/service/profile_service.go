package service

import (
	"github.com/yourusername/linkbio/repository"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
	userRepo    *repository.UserRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository, userRepo *repository.UserRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

func (s *ProfileService) GetByUsername(username string) (*repository.Profile, error) {
	return s.profileRepo.GetByUsername(username)
}

func (s *ProfileService) GetByUserID(userID string) (*repository.Profile, error) {
	profile, err := s.profileRepo.GetByUserID(userID)
	if err != nil {
		// If profile doesn't exist, create it automatically
		return s.profileRepo.Create(userID)
	}
	return profile, nil
}

func (s *ProfileService) Create(userID string) (*repository.Profile, error) {
	return s.profileRepo.Create(userID)
}

func (s *ProfileService) Update(userID string, data map[string]interface{}) (*repository.Profile, error) {
	return s.profileRepo.Update(userID, data)
}
