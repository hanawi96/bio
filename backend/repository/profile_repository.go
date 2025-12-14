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
	var themeJSON, customThemeJSON, headerJSON sql.NullString
	query := `
		SELECT p.id, p.user_id, u.username, p.avatar_url, p.bio, p.theme_name,
		       p.theme_config, p.custom_theme_config, p.header_config, p.social_links, p.custom_css,
		       COALESCE(p.show_share_button, true), COALESCE(p.show_subscribe_button, true), COALESCE(p.hide_branding, false),
		       p.created_at, p.updated_at
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE u.username = $1
	`
	err := r.db.QueryRow(query, username).Scan(
		&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL,
		&profile.Bio, &profile.ThemeName, &themeJSON, &customThemeJSON, &headerJSON, &profile.SocialLinks, &profile.CustomCSS,
		&profile.ShowShareButton, &profile.ShowSubscribeButton, &profile.HideBranding,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if themeJSON.Valid && themeJSON.String != "" {
		json.Unmarshal([]byte(themeJSON.String), &profile.ThemeConfig)
	}
	if customThemeJSON.Valid && customThemeJSON.String != "" {
		json.Unmarshal([]byte(customThemeJSON.String), &profile.CustomThemeConfig)
	}
	if headerJSON.Valid && headerJSON.String != "" {
		json.Unmarshal([]byte(headerJSON.String), &profile.HeaderConfig)
	}
	return &profile, nil
}

func (r *ProfileRepository) GetByUserID(userID string) (*Profile, error) {
	var profile Profile
	var themeJSON, customThemeJSON, headerJSON sql.NullString
	query := `
		SELECT p.id, p.user_id, u.username, p.avatar_url, p.bio, p.theme_name,
		       p.theme_config, p.custom_theme_config, p.header_config, p.social_links, p.custom_css,
		       COALESCE(p.show_share_button, true), COALESCE(p.show_subscribe_button, true), COALESCE(p.hide_branding, false),
		       p.created_at, p.updated_at
		FROM profiles p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Username, &profile.AvatarURL,
		&profile.Bio, &profile.ThemeName, &themeJSON, &customThemeJSON, &headerJSON, &profile.SocialLinks, &profile.CustomCSS,
		&profile.ShowShareButton, &profile.ShowSubscribeButton, &profile.HideBranding,
		&profile.CreatedAt, &profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if themeJSON.Valid && themeJSON.String != "" {
		json.Unmarshal([]byte(themeJSON.String), &profile.ThemeConfig)
	}
	if customThemeJSON.Valid && customThemeJSON.String != "" {
		json.Unmarshal([]byte(customThemeJSON.String), &profile.CustomThemeConfig)
	}
	if headerJSON.Valid && headerJSON.String != "" {
		json.Unmarshal([]byte(headerJSON.String), &profile.HeaderConfig)
	}
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
		    theme_name = COALESCE($4, theme_name),
		    theme_config = COALESCE($5, theme_config),
		    custom_theme_config = COALESCE($6, custom_theme_config),
		    header_config = COALESCE($7, header_config),
		    social_links = COALESCE($8, social_links),
		    show_share_button = COALESCE($9, show_share_button),
		    show_subscribe_button = COALESCE($10, show_subscribe_button),
		    hide_branding = COALESCE($11, hide_branding),
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	
	var themeConfig interface{}
	if tc, ok := data["theme_config"]; ok {
		if tcStr, ok := tc.(string); ok {
			themeConfig = []byte(tcStr)
		} else {
			themeJSON, _ := json.Marshal(tc)
			themeConfig = themeJSON
		}
	}
	
	var customThemeConfig interface{}
	if ctc, ok := data["custom_theme_config"]; ok {
		if ctcStr, ok := ctc.(string); ok {
			customThemeConfig = []byte(ctcStr)
		} else {
			customThemeJSON, _ := json.Marshal(ctc)
			customThemeConfig = customThemeJSON
		}
	}
	
	var headerConfig interface{}
	if hc, ok := data["header_config"]; ok {
		if hcStr, ok := hc.(string); ok {
			headerConfig = []byte(hcStr)
		} else {
			headerJSON, _ := json.Marshal(hc)
			headerConfig = headerJSON
		}
	}
	
	_, err := r.db.Exec(query, userID, data["bio"], data["avatar_url"], data["theme_name"], themeConfig, customThemeConfig, headerConfig, data["social_links"],
		data["show_share_button"], data["show_subscribe_button"], data["hide_branding"])
	if err != nil {
		return nil, err
	}
	return r.GetByUserID(userID)
}
