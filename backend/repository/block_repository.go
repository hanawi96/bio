package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type BlockRepository struct {
	db *sql.DB
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

func (r *BlockRepository) GetByUserID(userID string) ([]Block, error) {
	query := `
		SELECT b.id, b.profile_id, b.block_type, b.position, b.is_active,
		       b.content, b.text_style, b.style, b.image_url, b.alt_text, b.video_url,
		       b.social_links, b.divider_style, b.placeholder, b.embed_url, b.embed_type,
		       b.link_id, b.created_at, b.updated_at
		FROM blocks b
		JOIN profiles p ON b.profile_id = p.id
		WHERE p.user_id = $1
		ORDER BY b.position ASC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var block Block
		var socialLinksJSON []byte
		
		err := rows.Scan(
			&block.ID, &block.ProfileID, &block.BlockType, &block.Position, &block.IsActive,
			&block.Content, &block.TextStyle, &block.Style, &block.ImageURL, &block.AltText, &block.VideoURL,
			&socialLinksJSON, &block.DividerStyle, &block.Placeholder, &block.EmbedURL, &block.EmbedType,
			&block.LinkID, &block.CreatedAt, &block.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if socialLinksJSON != nil {
			json.Unmarshal(socialLinksJSON, &block.SocialLinks)
		}
		
		blocks = append(blocks, block)
	}
	return blocks, nil
}

func (r *BlockRepository) Create(userID string, data map[string]interface{}) (*Block, error) {
	var profileID string
	err := r.db.QueryRow(`SELECT id FROM profiles WHERE user_id = $1`, userID).Scan(&profileID)
	if err != nil {
		return nil, err
	}

	// Get max position
	var maxPosition int
	err = r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM blocks WHERE profile_id = $1`, profileID).Scan(&maxPosition)
	if err != nil {
		return nil, err
	}

	// Helper to get value or nil
	getVal := func(key string) interface{} {
		if val, ok := data[key]; ok && val != nil && val != "" {
			return val
		}
		return nil
	}

	// Serialize social_links if present and not empty
	var socialLinksJSON interface{}
	if socialLinks, ok := data["social_links"]; ok && socialLinks != nil {
		// Check if it's an empty array
		if arr, isArray := socialLinks.([]interface{}); isArray && len(arr) == 0 {
			socialLinksJSON = nil
		} else {
			jsonBytes, _ := json.Marshal(socialLinks)
			socialLinksJSON = jsonBytes
		}
	}

	// Get is_active value, default to true
	isActive := true
	if val, ok := data["is_active"].(bool); ok {
		isActive = val
	}

	var block Block
	query := `
		INSERT INTO blocks (profile_id, block_type, position, is_active, content, text_style, style, image_url, alt_text, video_url, social_links, divider_style, placeholder, embed_url, embed_type, link_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, profile_id, block_type, position, is_active, content, text_style, style, image_url, alt_text, video_url, social_links, divider_style, placeholder, embed_url, embed_type, link_id, created_at, updated_at
	`
	
	var socialLinksResult []byte
	err = r.db.QueryRow(
		query,
		profileID,
		data["block_type"],
		maxPosition+1,
		isActive,
		getVal("content"),
		getVal("text_style"),
		getVal("style"),
		getVal("image_url"),
		getVal("alt_text"),
		getVal("video_url"),
		socialLinksJSON,
		getVal("divider_style"),
		getVal("placeholder"),
		getVal("embed_url"),
		getVal("embed_type"),
		getVal("link_id"),
	).Scan(
		&block.ID, &block.ProfileID, &block.BlockType, &block.Position, &block.IsActive,
		&block.Content, &block.TextStyle, &block.Style, &block.ImageURL, &block.AltText, &block.VideoURL,
		&socialLinksResult, &block.DividerStyle, &block.Placeholder, &block.EmbedURL, &block.EmbedType,
		&block.LinkID, &block.CreatedAt, &block.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if socialLinksResult != nil {
		json.Unmarshal(socialLinksResult, &block.SocialLinks)
	}
	
	return &block, nil
}

func (r *BlockRepository) Update(blockID string, data map[string]interface{}) (*Block, error) {
	// Serialize social_links if present
	var socialLinksJSON interface{}
	if socialLinks, ok := data["social_links"]; ok && socialLinks != nil {
		jsonBytes, _ := json.Marshal(socialLinks)
		socialLinksJSON = jsonBytes
	}

	// Helper to get value or nil
	getVal := func(key string) interface{} {
		if val, ok := data[key]; ok {
			return val
		}
		return nil
	}

	query := `
		UPDATE blocks
		SET content = COALESCE($2, content),
		    text_style = COALESCE($3, text_style),
		    style = COALESCE($4, style),
		    image_url = COALESCE($5, image_url),
		    alt_text = COALESCE($6, alt_text),
		    video_url = COALESCE($7, video_url),
		    social_links = COALESCE($8, social_links),
		    divider_style = COALESCE($9, divider_style),
		    placeholder = COALESCE($10, placeholder),
		    embed_url = COALESCE($11, embed_url),
		    embed_type = COALESCE($12, embed_type),
		    is_active = COALESCE($13, is_active),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, profile_id, block_type, position, is_active, content, text_style, style, image_url, alt_text, video_url, social_links, divider_style, placeholder, embed_url, embed_type, link_id, created_at, updated_at
	`
	
	var block Block
	var socialLinksResult []byte
	
	err := r.db.QueryRow(
		query,
		blockID,
		getVal("content"),
		getVal("text_style"),
		getVal("style"),
		getVal("image_url"),
		getVal("alt_text"),
		getVal("video_url"),
		socialLinksJSON,
		getVal("divider_style"),
		getVal("placeholder"),
		getVal("embed_url"),
		getVal("embed_type"),
		getVal("is_active"),
	).Scan(
		&block.ID, &block.ProfileID, &block.BlockType, &block.Position, &block.IsActive,
		&block.Content, &block.TextStyle, &block.Style, &block.ImageURL, &block.AltText, &block.VideoURL,
		&socialLinksResult, &block.DividerStyle, &block.Placeholder, &block.EmbedURL, &block.EmbedType,
		&block.LinkID, &block.CreatedAt, &block.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	if socialLinksResult != nil {
		json.Unmarshal(socialLinksResult, &block.SocialLinks)
	}
	
	return &block, nil
}

func (r *BlockRepository) Delete(blockID string) error {
	_, err := r.db.Exec(`DELETE FROM blocks WHERE id = $1`, blockID)
	return err
}

func (r *BlockRepository) Reorder(userID string, blockIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, blockID := range blockIDs {
		_, err := tx.Exec(`
			UPDATE blocks 
			SET position = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 
			AND profile_id IN (SELECT id FROM profiles WHERE user_id = $3)
		`, i, blockID, userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *BlockRepository) BulkDelete(userID string, blockIDs []string) error {
	if len(blockIDs) == 0 {
		return nil
	}

	placeholders := ""
	args := []interface{}{userID}
	for i, id := range blockIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", i+2)
		args = append(args, id)
	}

	query := `DELETE FROM blocks WHERE id IN (` + placeholders + `) 
	          AND profile_id IN (SELECT id FROM profiles WHERE user_id = $1)`
	
	_, err := r.db.Exec(query, args...)
	return err
}
