package repository

import (
	"database/sql"
	"encoding/json"
)

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetByUsername(username string) (*Profile, error) {
	var profile Profile
	var themeJSON []byte
	query := `
		SELECT p.id, p.user_id, u.username, p.avatar_url, p.bio, 
		       p.theme_config, p.custom_css, p.created_at, p.updated_at
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE u.username = $1
	`
	err := r.db.QueryRow(query, username).Scan(
		&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL,
		&profile.Bio, &themeJSON, &profile.CustomCSS, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(themeJSON, &profile.ThemeConfig)
	return &profile, nil
}

func (r *ProfileRepository) GetByUserID(userID string) (*Profile, error) {
	var profile Profile
	var themeJSON []byte
	query := `
		SELECT p.id, p.user_id, u.username, p.avatar_url, p.bio,
		       p.theme_config, p.custom_css, p.created_at, p.updated_at
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL,
		&profile.Bio, &themeJSON, &profile.CustomCSS, &profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(themeJSON, &profile.ThemeConfig)
	return &profile, nil
}

func (r *ProfileRepository) Create(userID string) (*Profile, error) {
	query := `
		INSERT INTO profiles (user_id, theme_config)
		VALUES ($1, '{}'::jsonb)
		RETURNING id
	`
	var profileID string
	err := r.db.QueryRow(query, userID).Scan(&profileID)
	if err != nil {
		return nil, err
	}
	return r.GetByUserID(userID)
}

func (r *ProfileRepository) Update(userID string, data map[string]interface{}) (*Profile, error) {
	query := `
		UPDATE profiles
		SET bio = COALESCE($2, bio),
		    avatar_url = COALESCE($3, avatar_url),
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID, data["bio"], data["avatar_url"])
	if err != nil {
		return nil, err
	}
	return r.GetByUserID(userID)
}
