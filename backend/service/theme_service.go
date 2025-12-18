package service

import (
	"fmt"

	"github.com/yourusername/linkbio/repository"
)

type ThemeService struct {
	themeRepo *repository.ThemeRepository
}

func NewThemeService(themeRepo *repository.ThemeRepository) *ThemeService {
	return &ThemeService{
		themeRepo: themeRepo,
	}
}

// GetMyThemes retrieves all themes for a user
func (s *ThemeService) GetMyThemes(userID string) ([]repository.UserTheme, error) {
	return s.themeRepo.GetByUserID(userID)
}

// GetThemeByID retrieves a theme by ID (with ownership check)
func (s *ThemeService) GetThemeByID(themeID, userID string) (*repository.UserTheme, error) {
	theme, err := s.themeRepo.GetByID(themeID)
	if err != nil {
		return nil, err
	}

	// Check ownership or public access
	if theme.UserID != userID && !theme.IsPublic {
		return nil, fmt.Errorf("unauthorized: theme is private")
	}

	return theme, nil
}

// CreateTheme creates a new theme
func (s *ThemeService) CreateTheme(userID, name string, description *string, config map[string]interface{}) (*repository.UserTheme, error) {
	// Validate theme name
	if name == "" {
		return nil, fmt.Errorf("theme name is required")
	}

	// Check if name already exists for this user
	exists, err := s.themeRepo.CheckNameExists(userID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("theme with name '%s' already exists", name)
	}

	// Validate config
	if config == nil || len(config) == 0 {
		return nil, fmt.Errorf("theme config is required")
	}

	return s.themeRepo.Create(userID, name, description, config)
}

// UpdateTheme updates a theme
func (s *ThemeService) UpdateTheme(themeID, userID string, data map[string]interface{}) (*repository.UserTheme, error) {
	// If updating name, check for duplicates
	if name, ok := data["name"].(string); ok && name != "" {
		exists, err := s.themeRepo.CheckNameExists(userID, name)
		if err != nil {
			return nil, err
		}
		if exists {
			// Check if it's the same theme
			theme, err := s.themeRepo.GetByID(themeID)
			if err != nil {
				return nil, err
			}
			if theme.Name != name {
				return nil, fmt.Errorf("theme with name '%s' already exists", name)
			}
		}
	}

	return s.themeRepo.Update(themeID, userID, data)
}

// DeleteTheme deletes a theme
func (s *ThemeService) DeleteTheme(themeID, userID string) error {
	return s.themeRepo.Delete(themeID, userID)
}

// PublishTheme makes a theme public
func (s *ThemeService) PublishTheme(themeID, userID string) (*repository.UserTheme, error) {
	data := map[string]interface{}{
		"is_public": true,
	}
	return s.themeRepo.Update(themeID, userID, data)
}

// UnpublishTheme makes a theme private
func (s *ThemeService) UnpublishTheme(themeID, userID string) (*repository.UserTheme, error) {
	data := map[string]interface{}{
		"is_public": false,
	}
	return s.themeRepo.Update(themeID, userID, data)
}

// GetPublicThemes retrieves public themes for marketplace
func (s *ThemeService) GetPublicThemes(limit, offset int) ([]repository.UserTheme, error) {
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}
	if offset < 0 {
		offset = 0
	}
	return s.themeRepo.GetPublicThemes(limit, offset)
}

// GetThemeBySlug retrieves a public theme by slug
func (s *ThemeService) GetThemeBySlug(slug string) (*repository.UserTheme, error) {
	return s.themeRepo.GetBySlug(slug)
}

// ImportTheme imports a theme from config (creates a copy)
func (s *ThemeService) ImportTheme(userID, name string, description *string, config map[string]interface{}, sourceThemeID *string) (*repository.UserTheme, error) {
	// Validate
	if name == "" {
		return nil, fmt.Errorf("theme name is required")
	}
	if config == nil || len(config) == 0 {
		return nil, fmt.Errorf("theme config is required")
	}

	// Check for duplicate name
	exists, err := s.themeRepo.CheckNameExists(userID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("theme with name '%s' already exists", name)
	}

	// If importing from a public theme, increment downloads
	if sourceThemeID != nil && *sourceThemeID != "" {
		if err := s.themeRepo.IncrementDownloads(*sourceThemeID); err != nil {
			// Log error but don't fail the import
			fmt.Printf("Warning: failed to increment downloads for theme %s: %v\n", *sourceThemeID, err)
		}
	}

	return s.themeRepo.Create(userID, name, description, config)
}

// ExportTheme exports a theme config (just returns the theme)
func (s *ThemeService) ExportTheme(themeID, userID string) (*repository.UserTheme, error) {
	return s.GetThemeByID(themeID, userID)
}
