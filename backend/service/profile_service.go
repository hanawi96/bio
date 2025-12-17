package service

import (
	"github.com/yourusername/linkbio/repository"
)

type ProfileService struct {
	profileRepo *repository.ProfileRepository
	userRepo    *repository.UserRepository
	linkRepo    *repository.LinkRepository
	blockRepo   *repository.BlockRepository
}

func NewProfileService(profileRepo *repository.ProfileRepository, userRepo *repository.UserRepository, linkRepo *repository.LinkRepository, blockRepo *repository.BlockRepository) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		userRepo:    userRepo,
		linkRepo:    linkRepo,
		blockRepo:   blockRepo,
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

func (s *ProfileService) GetPublicProfileWithLinks(username string) (map[string]interface{}, error) {
	profile, err := s.profileRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	// Get user to fetch links and blocks
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	links, err := s.linkRepo.GetByUserID(user.ID)
	if err != nil {
		links = []repository.Link{}
	}

	blocks, err := s.blockRepo.GetByUserID(user.ID)
	if err != nil {
		blocks = []repository.Block{}
	}

	return map[string]interface{}{
		"profile": profile,
		"links":   links,
		"blocks":  blocks,
	}, nil
}

// ApplyTheme applies theme preset to profile and all groups
func (s *ProfileService) ApplyTheme(userID string, themeName string, themeConfig map[string]interface{}, cardStyles map[string]interface{}, textStyles string, headerConfig map[string]interface{}) (map[string]interface{}, error) {
	// 1. Update profile theme_config, theme_name and header_config
	updateData := map[string]interface{}{
		"theme_name":   themeName,
		"theme_config": themeConfig,
	}
	if headerConfig != nil {
		updateData["header_config"] = headerConfig
	}
	
	profile, err := s.profileRepo.Update(userID, updateData)
	if err != nil {
		return nil, err
	}

	// 2. Update all link groups with card styles
	err = s.linkRepo.UpdateAllGroupsCardStyles(userID, cardStyles)
	if err != nil {
		return nil, err
	}

	// 3. Update all text groups with text styles
	err = s.blockRepo.UpdateAllGroupsStyle(userID, textStyles)
	if err != nil {
		return nil, err
	}

	// 4. Fetch updated data
	links, err := s.linkRepo.GetByUserID(userID)
	if err != nil {
		links = []repository.Link{}
	}

	blocks, err := s.blockRepo.GetByUserID(userID)
	if err != nil {
		blocks = []repository.Block{}
	}

	return map[string]interface{}{
		"profile": profile,
		"links":   links,
		"blocks":  blocks,
	}, nil
}
