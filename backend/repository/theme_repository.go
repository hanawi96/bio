package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type UserTheme struct {
	ID             string                 `json:"id"`
	UserID         string                 `json:"user_id"`
	Name           string                 `json:"name"`
	Slug           *string                `json:"slug"`
	Description    *string                `json:"description"`
	Config         map[string]interface{} `json:"config"`
	ThumbnailURL   *string                `json:"thumbnail_url"`
	IsPublic       bool                   `json:"is_public"`
	DownloadsCount int                    `json:"downloads_count"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

type ThemeRepository struct {
	db *sql.DB
}

func NewThemeRepository(db *sql.DB) *ThemeRepository {
	return &ThemeRepository{db: db}
}

// GetByID retrieves a theme by ID
func (r *ThemeRepository) GetByID(themeID string) (*UserTheme, error) {
	var theme UserTheme
	var configJSON sql.NullString

	query := `
		SELECT id, user_id, name, slug, description, config, thumbnail_url,
		       is_public, downloads_count, created_at, updated_at
		FROM user_themes
		WHERE id = $1
	`

	err := r.db.QueryRow(query, themeID).Scan(
		&theme.ID, &theme.UserID, &theme.Name, &theme.Slug, &theme.Description,
		&configJSON, &theme.ThumbnailURL, &theme.IsPublic, &theme.DownloadsCount,
		&theme.CreatedAt, &theme.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if configJSON.Valid && configJSON.String != "" {
		if err := json.Unmarshal([]byte(configJSON.String), &theme.Config); err != nil {
			return nil, fmt.Errorf("failed to parse theme config: %w", err)
		}
	}

	return &theme, nil
}

// GetByUserID retrieves all themes for a user
func (r *ThemeRepository) GetByUserID(userID string) ([]UserTheme, error) {
	query := `
		SELECT id, user_id, name, slug, description, config, thumbnail_url,
		       is_public, downloads_count, created_at, updated_at
		FROM user_themes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var themes []UserTheme
	for rows.Next() {
		var theme UserTheme
		var configJSON sql.NullString

		err := rows.Scan(
			&theme.ID, &theme.UserID, &theme.Name, &theme.Slug, &theme.Description,
			&configJSON, &theme.ThumbnailURL, &theme.IsPublic, &theme.DownloadsCount,
			&theme.CreatedAt, &theme.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if configJSON.Valid && configJSON.String != "" {
			if err := json.Unmarshal([]byte(configJSON.String), &theme.Config); err != nil {
				return nil, fmt.Errorf("failed to parse theme config: %w", err)
			}
		}

		themes = append(themes, theme)
	}

	return themes, rows.Err()
}

// GetPublicThemes retrieves all public themes (for marketplace)
func (r *ThemeRepository) GetPublicThemes(limit, offset int) ([]UserTheme, error) {
	query := `
		SELECT id, user_id, name, slug, description, config, thumbnail_url,
		       is_public, downloads_count, created_at, updated_at
		FROM user_themes
		WHERE is_public = true
		ORDER BY downloads_count DESC, created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var themes []UserTheme
	for rows.Next() {
		var theme UserTheme
		var configJSON sql.NullString

		err := rows.Scan(
			&theme.ID, &theme.UserID, &theme.Name, &theme.Slug, &theme.Description,
			&configJSON, &theme.ThumbnailURL, &theme.IsPublic, &theme.DownloadsCount,
			&theme.CreatedAt, &theme.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if configJSON.Valid && configJSON.String != "" {
			if err := json.Unmarshal([]byte(configJSON.String), &theme.Config); err != nil {
				return nil, fmt.Errorf("failed to parse theme config: %w", err)
			}
		}

		themes = append(themes, theme)
	}

	return themes, rows.Err()
}

// GetBySlug retrieves a public theme by slug
func (r *ThemeRepository) GetBySlug(slug string) (*UserTheme, error) {
	var theme UserTheme
	var configJSON sql.NullString

	query := `
		SELECT id, user_id, name, slug, description, config, thumbnail_url,
		       is_public, downloads_count, created_at, updated_at
		FROM user_themes
		WHERE slug = $1 AND is_public = true
	`

	err := r.db.QueryRow(query, slug).Scan(
		&theme.ID, &theme.UserID, &theme.Name, &theme.Slug, &theme.Description,
		&configJSON, &theme.ThumbnailURL, &theme.IsPublic, &theme.DownloadsCount,
		&theme.CreatedAt, &theme.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if configJSON.Valid && configJSON.String != "" {
		if err := json.Unmarshal([]byte(configJSON.String), &theme.Config); err != nil {
			return nil, fmt.Errorf("failed to parse theme config: %w", err)
		}
	}

	return &theme, nil
}

// Create creates a new theme
func (r *ThemeRepository) Create(userID, name string, description *string, config map[string]interface{}) (*UserTheme, error) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config: %w", err)
	}

	query := `
		INSERT INTO user_themes (user_id, name, description, config)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var themeID string
	err = r.db.QueryRow(query, userID, name, description, configJSON).Scan(&themeID)
	if err != nil {
		return nil, err
	}

	return r.GetByID(themeID)
}

// Update updates a theme
func (r *ThemeRepository) Update(themeID, userID string, data map[string]interface{}) (*UserTheme, error) {
	// First verify ownership
	var ownerID string
	err := r.db.QueryRow("SELECT user_id FROM user_themes WHERE id = $1", themeID).Scan(&ownerID)
	if err != nil {
		return nil, err
	}
	if ownerID != userID {
		return nil, fmt.Errorf("unauthorized: theme does not belong to user")
	}

	// Build dynamic update query
	query := "UPDATE user_themes SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}
	argCount := 1

	if name, ok := data["name"].(string); ok && name != "" {
		query += fmt.Sprintf(", name = $%d", argCount)
		args = append(args, name)
		argCount++
	}

	if description, ok := data["description"].(string); ok {
		query += fmt.Sprintf(", description = $%d", argCount)
		args = append(args, description)
		argCount++
	}

	if config, ok := data["config"].(map[string]interface{}); ok {
		configJSON, err := json.Marshal(config)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		query += fmt.Sprintf(", config = $%d", argCount)
		args = append(args, configJSON)
		argCount++
	}

	if thumbnailURL, ok := data["thumbnail_url"].(string); ok {
		query += fmt.Sprintf(", thumbnail_url = $%d", argCount)
		args = append(args, thumbnailURL)
		argCount++
	}

	if isPublic, ok := data["is_public"].(bool); ok {
		query += fmt.Sprintf(", is_public = $%d", argCount)
		args = append(args, isPublic)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, themeID)

	_, err = r.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return r.GetByID(themeID)
}

// Delete deletes a theme
func (r *ThemeRepository) Delete(themeID, userID string) error {
	// Verify ownership
	var ownerID string
	err := r.db.QueryRow("SELECT user_id FROM user_themes WHERE id = $1", themeID).Scan(&ownerID)
	if err != nil {
		return err
	}
	if ownerID != userID {
		return fmt.Errorf("unauthorized: theme does not belong to user")
	}

	query := "DELETE FROM user_themes WHERE id = $1"
	_, err = r.db.Exec(query, themeID)
	return err
}

// IncrementDownloads increments the downloads count for a theme
func (r *ThemeRepository) IncrementDownloads(themeID string) error {
	query := "UPDATE user_themes SET downloads_count = downloads_count + 1 WHERE id = $1"
	_, err := r.db.Exec(query, themeID)
	return err
}

// CheckNameExists checks if a theme name already exists for a user
func (r *ThemeRepository) CheckNameExists(userID, name string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM user_themes WHERE user_id = $1 AND name = $2)"
	err := r.db.QueryRow(query, userID, name).Scan(&exists)
	return exists, err
}
