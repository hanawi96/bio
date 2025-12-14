package repository

import "time"

type Profile struct {
	ID                    string                 `json:"id"`
	UserID                string                 `json:"user_id"`
	Username              string                 `json:"username"`
	AvatarURL             *string                `json:"avatar_url"`
	Bio                   *string                `json:"bio"`
	ThemeName             *string                `json:"theme_name"`
	ThemeConfig           map[string]interface{} `json:"theme_config"`
	CustomThemeConfig     map[string]interface{} `json:"custom_theme_config"`
	HeaderConfig          map[string]interface{} `json:"header_config"`
	SocialLinks           *string                `json:"social_links"`
	CustomCSS             *string                `json:"custom_css"`
	ShowShareButton       bool                   `json:"show_share_button"`
	ShowSubscribeButton   bool                   `json:"show_subscribe_button"`
	HideBranding          bool                   `json:"hide_branding"`
	CreatedAt             time.Time              `json:"created_at"`
	UpdatedAt             time.Time              `json:"updated_at"`
}

type Link struct {
	ID                    string     `json:"id"`
	ProfileID             string     `json:"profile_id"`
	ParentID              *string    `json:"parent_id"`
	IsGroup               bool       `json:"is_group"`
	GroupTitle            *string    `json:"group_title"`
	GroupLayout           string     `json:"group_layout"`
	GridColumns           int        `json:"grid_columns"`
	GridAspectRatio       string     `json:"grid_aspect_ratio"`
	Title                 string     `json:"title"`
	URL                   string     `json:"url"`
	Description           *string    `json:"description"`
	ThumbnailURL          *string    `json:"thumbnail_url"`
	ImageShape            string     `json:"image_shape"`
	LayoutType            string     `json:"layout_type"`
	ImagePlacement        string     `json:"image_placement"`
	TextAlignment         string     `json:"text_alignment"`
	TextSize              string     `json:"text_size"`
	ShowOutline           bool       `json:"show_outline"`
	ShowShadow            bool       `json:"show_shadow"`
	ShadowX               int        `json:"shadow_x"`
	ShadowY               int        `json:"shadow_y"`
	ShadowBlur            int        `json:"shadow_blur"`
	ShowDescription       bool       `json:"show_description"`
	ShowText              bool       `json:"show_text"`
	HasCardBackground     bool       `json:"has_card_background"`
	CardBackgroundColor   string     `json:"card_background_color"`
	CardBackgroundOpacity int        `json:"card_background_opacity"`
	CardBorderRadius      int        `json:"card_border_radius"`
	CardTextColor         string     `json:"card_text_color"`
	HasCardBorder         bool       `json:"has_card_border"`
	CardBorderColor       string     `json:"card_border_color"`
	CardBorderStyle       string     `json:"card_border_style"`
	CardBorderWidth       int        `json:"card_border_width"`
	Style                 *string    `json:"style"`
	Position              int        `json:"position"`
	Clicks                int        `json:"clicks"`
	IsActive              bool       `json:"is_active"`
	IsPinned              bool       `json:"is_pinned"`
	ScheduledAt           *time.Time `json:"scheduled_at"`
	ExpiresAt             *time.Time `json:"expires_at"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	Children              []Link     `json:"children,omitempty"`
}

type Block struct {
	ID              string                   `json:"id"`
	ProfileID       string                   `json:"profile_id"`
	ParentID        *string                  `json:"parent_id"`
	IsGroup         bool                     `json:"is_group"`
	GroupTitle      *string                  `json:"group_title"`
	GroupLayout     *string                  `json:"group_layout"`
	GridColumns     *int                     `json:"grid_columns"`
	GridAspectRatio *string                  `json:"grid_aspect_ratio"`
	BlockType       string                   `json:"block_type"`
	Position        int                      `json:"position"`
	IsActive        bool                     `json:"is_active"`
	Content         *string                  `json:"content,omitempty"`
	TextStyle       *string                  `json:"text_style,omitempty"`
	Style           *string                  `json:"style,omitempty"`
	ImageURL        *string                  `json:"image_url,omitempty"`
	AltText         *string                  `json:"alt_text,omitempty"`
	VideoURL        *string                  `json:"video_url,omitempty"`
	SocialLinks     []map[string]interface{} `json:"social_links,omitempty"`
	DividerStyle    *string                  `json:"divider_style,omitempty"`
	Placeholder     *string                  `json:"placeholder,omitempty"`
	EmbedURL        *string                  `json:"embed_url,omitempty"`
	EmbedType       *string                  `json:"embed_type,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
	Children        []Block                  `json:"children,omitempty"`
}
