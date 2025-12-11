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
		SELECT b.id, b.profile_id, b.parent_id, b.is_group, b.group_title, b.group_layout,
		       b.grid_columns, b.grid_aspect_ratio,
		       b.block_type, b.position, b.is_active,
		       b.content, b.text_style, b.style, b.image_url, b.alt_text, b.video_url,
		       b.social_links, b.divider_style, b.placeholder, b.embed_url, b.embed_type,
		       b.created_at, b.updated_at
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

	var allBlocks []Block
	blockMap := make(map[string]*Block)

	for rows.Next() {
		var block Block
		var socialLinksJSON []byte

		err := rows.Scan(
			&block.ID, &block.ProfileID, &block.ParentID, &block.IsGroup, &block.GroupTitle, &block.GroupLayout,
			&block.GridColumns, &block.GridAspectRatio,
			&block.BlockType, &block.Position, &block.IsActive,
			&block.Content, &block.TextStyle, &block.Style, &block.ImageURL, &block.AltText, &block.VideoURL,
			&socialLinksJSON, &block.DividerStyle, &block.Placeholder, &block.EmbedURL, &block.EmbedType,
			&block.CreatedAt, &block.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if socialLinksJSON != nil {
			json.Unmarshal(socialLinksJSON, &block.SocialLinks)
		}

		allBlocks = append(allBlocks, block)
		blockMap[block.ID] = &allBlocks[len(allBlocks)-1]
	}

	// Build tree structure: attach children to parents
	var rootBlocks []Block
	for i := range allBlocks {
		block := &allBlocks[i]
		if block.ParentID != nil && *block.ParentID != "" {
			// This is a child block
			if parent, exists := blockMap[*block.ParentID]; exists {
				parent.Children = append(parent.Children, *block)
			}
		} else {
			// This is a root block
			rootBlocks = append(rootBlocks, *block)
		}
	}

	// Update the parent blocks with their children in the result
	for i := range rootBlocks {
		if rootBlocks[i].IsGroup {
			if updated, exists := blockMap[rootBlocks[i].ID]; exists {
				rootBlocks[i].Children = updated.Children
			}
		}
	}

	return rootBlocks, nil
}

func (r *BlockRepository) Create(userID string, data map[string]interface{}) (*Block, error) {
	var profileID string
	err := r.db.QueryRow(`SELECT id FROM profiles WHERE user_id = $1`, userID).Scan(&profileID)
	if err != nil {
		return nil, err
	}

	// Get max position from both blocks and links
	var maxBlockPos, maxLinkPos int
	r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM blocks WHERE profile_id = $1`, profileID).Scan(&maxBlockPos)
	r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM links WHERE profile_id = $1`, profileID).Scan(&maxLinkPos)

	maxPosition := maxBlockPos
	if maxLinkPos > maxPosition {
		maxPosition = maxLinkPos
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

	// Get is_group and set defaults if group
	isGroup := false
	if val, ok := data["is_group"].(bool); ok {
		isGroup = val
	}

	// Helper for group-related fields - returns interface for nullable fields
	getGroupVal := func(key string, defaultVal interface{}) interface{} {
		if !isGroup {
			return nil
		}
		if val, ok := data[key]; ok && val != nil && val != "" {
			return val
		}
		return defaultVal
	}

	var block Block
	query := `
		INSERT INTO blocks (profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio, block_type, position, is_active, content, text_style, style, image_url, alt_text, video_url, social_links, divider_style, placeholder, embed_url, embed_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio, block_type, position, is_active, content, text_style, style, image_url, alt_text, video_url, social_links, divider_style, placeholder, embed_url, embed_type, created_at, updated_at
	`

	var socialLinksResult []byte
	err = r.db.QueryRow(
		query,
		profileID,
		getVal("parent_id"),
		isGroup,
		getVal("group_title"),
		getGroupVal("group_layout", "list"),
		getGroupVal("grid_columns", 2),
		getGroupVal("grid_aspect_ratio", "3:2"),
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
	).Scan(
		&block.ID, &block.ProfileID, &block.ParentID, &block.IsGroup, &block.GroupTitle, &block.GroupLayout,
		&block.GridColumns, &block.GridAspectRatio,
		&block.BlockType, &block.Position, &block.IsActive,
		&block.Content, &block.TextStyle, &block.Style, &block.ImageURL, &block.AltText, &block.VideoURL,
		&socialLinksResult, &block.DividerStyle, &block.Placeholder, &block.EmbedURL, &block.EmbedType,
		&block.CreatedAt, &block.UpdatedAt,
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
		SET parent_id = COALESCE($2, parent_id),
		    is_group = COALESCE($3, is_group),
		    group_title = COALESCE($4, group_title),
		    group_layout = COALESCE($5, group_layout),
		    grid_columns = COALESCE($6, grid_columns),
		    grid_aspect_ratio = COALESCE($7, grid_aspect_ratio),
		    content = COALESCE($8, content),
		    text_style = COALESCE($9, text_style),
		    style = COALESCE($10, style),
		    image_url = COALESCE($11, image_url),
		    alt_text = COALESCE($12, alt_text),
		    video_url = COALESCE($13, video_url),
		    social_links = COALESCE($14, social_links),
		    divider_style = COALESCE($15, divider_style),
		    placeholder = COALESCE($16, placeholder),
		    embed_url = COALESCE($17, embed_url),
		    embed_type = COALESCE($18, embed_type),
		    is_active = COALESCE($19, is_active),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio, block_type, position, is_active, content, text_style, style, image_url, alt_text, video_url, social_links, divider_style, placeholder, embed_url, embed_type, created_at, updated_at
	`

	var block Block
	var socialLinksResult []byte

	err := r.db.QueryRow(
		query,
		blockID,
		getVal("parent_id"),
		getVal("is_group"),
		getVal("group_title"),
		getVal("group_layout"),
		getVal("grid_columns"),
		getVal("grid_aspect_ratio"),
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
		&block.ID, &block.ProfileID, &block.ParentID, &block.IsGroup, &block.GroupTitle, &block.GroupLayout,
		&block.GridColumns, &block.GridAspectRatio,
		&block.BlockType, &block.Position, &block.IsActive,
		&block.Content, &block.TextStyle, &block.Style, &block.ImageURL, &block.AltText, &block.VideoURL,
		&socialLinksResult, &block.DividerStyle, &block.Placeholder, &block.EmbedURL, &block.EmbedType,
		&block.CreatedAt, &block.UpdatedAt,
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

	// Get profile ID
	var profileID string
	err = tx.QueryRow(`SELECT id FROM profiles WHERE user_id = $1`, userID).Scan(&profileID)
	if err != nil {
		return err
	}

	// Update block positions
	position := 0
	for _, blockID := range blockIDs {
		_, err := tx.Exec(`
			UPDATE blocks 
			SET position = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND profile_id = $3
		`, position, blockID, profileID)
		if err != nil {
			return err
		}
		position++
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

// ReorderGroupBlocks reorders blocks within a group
func (r *BlockRepository) ReorderGroupBlocks(userID string, groupID string, blockIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify group belongs to user
	var profileID string
	err = tx.QueryRow(`
		SELECT b.profile_id 
		FROM blocks b
		JOIN profiles p ON b.profile_id = p.id
		WHERE b.id = $1 AND p.user_id = $2 AND b.is_group = true
	`, groupID, userID).Scan(&profileID)
	if err != nil {
		return err
	}

	// Update position for each block
	for i, blockID := range blockIDs {
		_, err := tx.Exec(`
			UPDATE blocks 
			SET position = $1 
			WHERE id = $2 AND parent_id = $3
		`, i, blockID, groupID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
