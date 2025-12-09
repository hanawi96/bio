package repository

import "time"

type Profile struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Username    string                 `json:"username"`
	AvatarURL   *string                `json:"avatar_url"`
	Bio         *string                `json:"bio"`
	ThemeConfig map[string]interface{} `json:"theme_config"`
	CustomCSS   *string                `json:"custom_css"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Link struct {
	ID           string     `json:"id"`
	ProfileID    string     `json:"profile_id"`
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	ThumbnailURL *string    `json:"thumbnail_url"`
	LayoutType   string     `json:"layout_type"`
	Position     int        `json:"position"`
	Clicks       int        `json:"clicks"`
	IsActive     bool       `json:"is_active"`
	IsPinned     bool       `json:"is_pinned"`
	ScheduledAt  *time.Time `json:"scheduled_at"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Block struct {
	ID            string                   `json:"id"`
	ProfileID     string                   `json:"profile_id"`
	BlockType     string                   `json:"block_type"`
	Position      int                      `json:"position"`
	IsActive      bool                     `json:"is_active"`
	Content       *string                  `json:"content,omitempty"`
	TextStyle     *string                  `json:"text_style,omitempty"`
	Style         *string                  `json:"style,omitempty"`
	ImageURL      *string                  `json:"image_url,omitempty"`
	AltText       *string                  `json:"alt_text,omitempty"`
	VideoURL      *string                  `json:"video_url,omitempty"`
	SocialLinks   []map[string]interface{} `json:"social_links,omitempty"`
	DividerStyle  *string                  `json:"divider_style,omitempty"`
	Placeholder   *string                  `json:"placeholder,omitempty"`
	EmbedURL      *string                  `json:"embed_url,omitempty"`
	EmbedType     *string                  `json:"embed_type,omitempty"`
	LinkID        *string                  `json:"link_id,omitempty"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedAt     time.Time                `json:"updated_at"`
}
