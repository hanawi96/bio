package repository

import (
	"database/sql"
	"fmt"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) GetByUserID(userID string) ([]Link, error) {
	return r.GetByUserIDWithFilters(userID, "", "", "", "")
}

func (r *LinkRepository) GetByUserIDWithFilters(userID, search, status, layoutType, sortBy string) ([]Link, error) {
	query := `
		SELECT l.id, l.profile_id, l.parent_id, l.is_group, l.group_title, l.group_layout, l.grid_columns, l.grid_aspect_ratio,
		       l.title, l.url, l.description, l.thumbnail_url, l.image_shape, l.layout_type, 
		       l.image_placement, l.text_alignment, l.text_size, l.show_outline, l.show_shadow, l.shadow_x, l.shadow_y, l.shadow_blur, l.show_description, l.show_text,
		       l.has_card_background, l.card_background_color, l.card_background_opacity, l.card_border_radius, l.card_text_color,
		       l.has_card_border, l.card_border_color, l.card_border_style, l.card_border_width, l.style,
		       l.position, l.clicks, l.is_active, l.is_pinned, l.scheduled_at, l.expires_at, 
		       l.created_at, l.updated_at
		FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE p.user_id = $1 AND l.parent_id IS NULL
	`
	args := []interface{}{userID}
	argCount := 1

	// Search filter
	if search != "" {
		argCount++
		searchPattern := "%" + search + "%"
		query += fmt.Sprintf(" AND (LOWER(l.title) LIKE LOWER($%d) OR LOWER(l.url) LIKE LOWER($%d))", argCount, argCount)
		args = append(args, searchPattern)
		fmt.Printf("🔍 Search filter added: pattern='%s'\n", searchPattern)
	}

	// Status filter
	if status == "active" {
		query += " AND l.is_active = true"
	} else if status == "inactive" {
		query += " AND l.is_active = false"
	}

	// Layout type filter
	if layoutType != "" {
		argCount++
		query += fmt.Sprintf(" AND l.layout_type = $%d", argCount)
		args = append(args, layoutType)
	}

	// Sorting
	switch sortBy {
	case "clicks":
		query += " ORDER BY l.clicks DESC, l.position ASC"
	case "created":
		query += " ORDER BY l.created_at DESC"
	case "updated":
		query += " ORDER BY l.updated_at DESC"
	case "title":
		query += " ORDER BY LOWER(l.title) ASC"
	default:
		// Default: sort by position only (pinned status is handled on frontend for public view)
		query += " ORDER BY l.position ASC"
	}

	fmt.Printf("📝 Final SQL Query: %s\n", query)
	fmt.Printf("📦 Args: %v\n", args)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		fmt.Printf("❌ Query error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(
			&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout, &link.GridColumns, &link.GridAspectRatio,
			&link.Title, &link.URL, &link.Description, &link.ThumbnailURL, &link.ImageShape, &link.LayoutType,
			&link.ImagePlacement, &link.TextAlignment, &link.TextSize, &link.ShowOutline, &link.ShowShadow, &link.ShadowX, &link.ShadowY, &link.ShadowBlur, &link.ShowDescription, &link.ShowText,
			&link.HasCardBackground, &link.CardBackgroundColor, &link.CardBackgroundOpacity, &link.CardBorderRadius, &link.CardTextColor,
			&link.HasCardBorder, &link.CardBorderColor, &link.CardBorderStyle, &link.CardBorderWidth, &link.Style,
			&link.Position, &link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
			&link.CreatedAt, &link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// If it's a group, fetch children
		if link.IsGroup {
			children, err := r.GetChildrenByParentID(link.ID)
			if err == nil {
				link.Children = children
			}
		}

		links = append(links, link)
	}
	return links, nil
}

func (r *LinkRepository) Create(userID string, data map[string]interface{}) (*Link, error) {
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

	var link Link
	query := `
		INSERT INTO links (profile_id, title, url, position, image_placement, text_alignment, text_size, show_outline, show_shadow, show_description, is_group)
		VALUES ($1, $2, $3, $4, 'left', 'left', 'M', false, false, true, false)
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout,
		          title, url, thumbnail_url, layout_type, image_placement, text_alignment, 
		          text_size, show_outline, show_shadow, show_description, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = r.db.QueryRow(query, profileID, data["title"], data["url"], maxPosition+1).Scan(
		&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout,
		&link.Title, &link.URL, &link.ThumbnailURL, &link.LayoutType,
		&link.ImagePlacement, &link.TextAlignment, &link.TextSize, &link.ShowOutline, &link.ShowShadow, &link.ShowDescription,
		&link.Position, &link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
		&link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) Update(linkID string, data map[string]interface{}) (*Link, error) {
	query := `
		UPDATE links
		SET title = COALESCE($2, title),
		    url = COALESCE($3, url),
		    thumbnail_url = COALESCE($4, thumbnail_url),
		    image_shape = COALESCE($5, image_shape),
		    layout_type = COALESCE($6, layout_type),
		    image_placement = COALESCE($7, image_placement),
		    text_alignment = COALESCE($8, text_alignment),
		    text_size = COALESCE($9, text_size),
		    show_outline = COALESCE($10, show_outline),
		    show_shadow = COALESCE($11, show_shadow),
		    shadow_x = COALESCE($12, shadow_x),
		    shadow_y = COALESCE($13, shadow_y),
		    shadow_blur = COALESCE($14, shadow_blur),
		    show_description = COALESCE($15, show_description),
		    show_text = COALESCE($16, show_text),
		    is_active = COALESCE($17, is_active),
		    scheduled_at = CASE WHEN $18::text IS NOT NULL THEN $18::timestamp ELSE scheduled_at END,
		    expires_at = CASE WHEN $19::text IS NOT NULL THEN $19::timestamp ELSE expires_at END,
		    group_title = COALESCE($20, group_title),
		    group_layout = COALESCE($21, group_layout),
		    has_card_background = COALESCE($22, has_card_background),
		    card_background_color = COALESCE($23, card_background_color),
		    card_background_opacity = COALESCE($24, card_background_opacity),
		    card_border_radius = COALESCE($25, card_border_radius),
		    card_text_color = COALESCE($26, card_text_color),
		    has_card_border = COALESCE($27, has_card_border),
		    card_border_color = COALESCE($28, card_border_color),
		    card_border_style = COALESCE($29, card_border_style),
		    card_border_width = COALESCE($30, card_border_width),
		    style = COALESCE($31, style),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio,
		          title, url, thumbnail_url, image_shape, layout_type, image_placement, text_alignment,
		          text_size, show_outline, show_shadow, shadow_x, shadow_y, shadow_blur, show_description, show_text,
		          has_card_background, card_background_color, card_background_opacity, card_border_radius, card_text_color,
		          has_card_border, card_border_color, card_border_style, card_border_width, style,
		          position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at, description
	`
	var link Link
	err := r.db.QueryRow(query, linkID, data["title"], data["url"], data["thumbnail_url"], data["image_shape"], data["layout_type"],
		data["image_placement"], data["text_alignment"], data["text_size"], data["show_outline"], data["show_shadow"], 
		data["shadow_x"], data["shadow_y"], data["shadow_blur"], data["show_description"],
		data["show_text"], data["is_active"], data["scheduled_at"], data["expires_at"], data["group_title"], data["group_layout"],
		data["has_card_background"], data["card_background_color"], data["card_background_opacity"], data["card_border_radius"], data["card_text_color"],
		data["has_card_border"], data["card_border_color"], data["card_border_style"], data["card_border_width"], data["style"]).Scan(
		&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout, &link.GridColumns, &link.GridAspectRatio,
		&link.Title, &link.URL, &link.ThumbnailURL, &link.ImageShape, &link.LayoutType,
		&link.ImagePlacement, &link.TextAlignment, &link.TextSize, &link.ShowOutline, &link.ShowShadow, &link.ShadowX, &link.ShadowY, &link.ShadowBlur, &link.ShowDescription,
		&link.ShowText, &link.HasCardBackground, &link.CardBackgroundColor, &link.CardBackgroundOpacity, &link.CardBorderRadius, &link.CardTextColor,
		&link.HasCardBorder, &link.CardBorderColor, &link.CardBorderStyle, &link.CardBorderWidth, &link.Style,
		&link.Position, &link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
		&link.CreatedAt, &link.UpdatedAt, &link.Description,
	)
	if err != nil {
		return nil, err
	}

	// If it's a group, fetch children
	if link.IsGroup {
		children, err := r.GetChildrenByParentID(link.ID)
		if err == nil {
			link.Children = children
		}
	}

	return &link, nil
}

func (r *LinkRepository) Delete(linkID string) error {
	_, err := r.db.Exec(`DELETE FROM links WHERE id = $1`, linkID)
	return err
}

// ReorderWithBlocks updates positions for both links and blocks in unified order
func (r *LinkRepository) ReorderWithBlocks(userID string, items []map[string]interface{}) error {
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

	// Update positions for all items
	for i, item := range items {
		itemType := item["type"].(string)
		itemID := item["id"].(string)

		if itemType == "link" {
			_, err := tx.Exec(`
				UPDATE links 
				SET position = $1, updated_at = CURRENT_TIMESTAMP
				WHERE id = $2 AND profile_id = $3
			`, i, itemID, profileID)
			if err != nil {
				return err
			}
		} else if itemType == "block" {
			_, err := tx.Exec(`
				UPDATE blocks 
				SET position = $1, updated_at = CURRENT_TIMESTAMP
				WHERE id = $2 AND profile_id = $3
			`, i, itemID, profileID)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *LinkRepository) Duplicate(userID string, linkID string) (*Link, error) {
	// Get original link
	var original Link
	query := `
		SELECT l.id, l.profile_id, l.parent_id, l.is_group, l.group_title, l.group_layout, l.grid_columns, l.grid_aspect_ratio,
		       l.title, l.url, l.description, l.thumbnail_url, l.image_shape, l.layout_type, 
		       l.image_placement, l.text_alignment, l.text_size, l.show_outline, l.show_shadow, l.show_description, l.show_text,
		       l.position, l.clicks, l.is_active, l.is_pinned, l.scheduled_at, l.expires_at, 
		       l.created_at, l.updated_at
		FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1 AND p.user_id = $2
	`
	err := r.db.QueryRow(query, linkID, userID).Scan(
		&original.ID, &original.ProfileID, &original.ParentID, &original.IsGroup, &original.GroupTitle, &original.GroupLayout, &original.GridColumns, &original.GridAspectRatio,
		&original.Title, &original.URL, &original.Description, &original.ThumbnailURL, &original.ImageShape,
		&original.LayoutType, &original.ImagePlacement, &original.TextAlignment, &original.TextSize,
		&original.ShowOutline, &original.ShowShadow, &original.ShowDescription, &original.ShowText, &original.Position, &original.Clicks,
		&original.IsActive, &original.IsPinned, &original.ScheduledAt, &original.ExpiresAt,
		&original.CreatedAt, &original.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get max position from both blocks and links
	var maxBlockPos, maxLinkPos int
	r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM blocks WHERE profile_id = $1`, original.ProfileID).Scan(&maxBlockPos)
	r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM links WHERE profile_id = $1`, original.ProfileID).Scan(&maxLinkPos)

	maxPosition := maxBlockPos
	if maxLinkPos > maxPosition {
		maxPosition = maxLinkPos
	}

	// Create duplicate with "(Copy)" suffix
	var duplicate Link
	insertQuery := `
		INSERT INTO links (profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio,
		                   title, url, description, thumbnail_url, image_shape, layout_type, image_placement, 
		                   text_alignment, text_size, show_outline, show_shadow, show_description, show_text, position, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, true)
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio,
		          title, url, description, thumbnail_url, image_shape, layout_type, image_placement, text_alignment,
		          text_size, show_outline, show_shadow, show_description, show_text, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	newTitle := original.Title
	if original.IsGroup && original.GroupTitle != nil {
		newTitle = *original.GroupTitle + " (Copy)"
	} else if !original.IsGroup {
		newTitle = original.Title + " (Copy)"
	}

	err = r.db.QueryRow(
		insertQuery,
		original.ProfileID,
		original.ParentID,
		original.IsGroup,
		&newTitle,
		original.GroupLayout,
		original.GridColumns,
		original.GridAspectRatio,
		newTitle,
		original.URL,
		original.Description,
		original.ThumbnailURL,
		original.ImageShape,
		original.LayoutType,
		original.ImagePlacement,
		original.TextAlignment,
		original.TextSize,
		original.ShowOutline,
		original.ShowShadow,
		original.ShowDescription,
		original.ShowText,
		maxPosition+1,
	).Scan(
		&duplicate.ID, &duplicate.ProfileID, &duplicate.ParentID, &duplicate.IsGroup, &duplicate.GroupTitle,
		&duplicate.GroupLayout, &duplicate.GridColumns, &duplicate.GridAspectRatio,
		&duplicate.Title, &duplicate.URL, &duplicate.Description, &duplicate.ThumbnailURL, &duplicate.ImageShape,
		&duplicate.LayoutType, &duplicate.ImagePlacement, &duplicate.TextAlignment, &duplicate.TextSize,
		&duplicate.ShowOutline, &duplicate.ShowShadow, &duplicate.ShowDescription, &duplicate.ShowText, &duplicate.Position, &duplicate.Clicks,
		&duplicate.IsActive, &duplicate.IsPinned, &duplicate.ScheduledAt, &duplicate.ExpiresAt,
		&duplicate.CreatedAt, &duplicate.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &duplicate, nil
}

func (r *LinkRepository) TogglePin(userID string, linkID string) (*Link, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Get profile ID, parent_id and current pin status
	var profileID string
	var parentID *string
	var currentPinned bool
	err = tx.QueryRow(`
		SELECT p.id, l.parent_id, COALESCE(l.is_pinned, false)
		FROM profiles p
		JOIN links l ON l.id = $2 AND l.profile_id = p.id
		WHERE p.user_id = $1
	`, userID, linkID).Scan(&profileID, &parentID, &currentPinned)
	if err != nil {
		return nil, err
	}

	var link Link
	if currentPinned {
		// UNPIN: Just toggle flag, keep position unchanged
		query := `
			UPDATE links 
			SET is_pinned = false, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND profile_id = $2
			RETURNING id, profile_id, parent_id, is_group, group_title, group_layout,
			          title, url, thumbnail_url, layout_type, image_placement, text_alignment,
			          text_size, show_outline, show_shadow, position, clicks, is_active, is_pinned,
			          scheduled_at, expires_at, created_at, updated_at
		`
		err = tx.QueryRow(query, linkID, profileID).Scan(
			&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout,
			&link.Title, &link.URL, &link.ThumbnailURL,
			&link.LayoutType, &link.ImagePlacement, &link.TextAlignment, &link.TextSize,
			&link.ShowOutline, &link.ShowShadow, &link.Position, &link.Clicks, &link.IsActive, &link.IsPinned,
			&link.ScheduledAt, &link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
		)
	} else {
		// PIN: Unpin all other links in the same group only, then pin this link
		// Position remains unchanged - pinned items will be sorted first in queries
		if parentID != nil {
			// Unpin other links in the same group
			_, err = tx.Exec(`UPDATE links SET is_pinned = false WHERE parent_id = $1 AND id != $2`, parentID, linkID)
		} else {
			// Unpin other top-level links
			_, err = tx.Exec(`UPDATE links SET is_pinned = false WHERE profile_id = $1 AND parent_id IS NULL AND id != $2`, profileID, linkID)
		}
		if err != nil {
			return nil, err
		}

		// Pin this link, keep position unchanged
		query := `
			UPDATE links 
			SET is_pinned = true, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND profile_id = $2
			RETURNING id, profile_id, parent_id, is_group, group_title, group_layout,
			          title, url, thumbnail_url, layout_type, image_placement, text_alignment,
			          text_size, show_outline, show_shadow, position, clicks, is_active, is_pinned,
			          scheduled_at, expires_at, created_at, updated_at
		`
		err = tx.QueryRow(query, linkID, profileID).Scan(
			&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout,
			&link.Title, &link.URL, &link.ThumbnailURL,
			&link.LayoutType, &link.ImagePlacement, &link.TextAlignment, &link.TextSize,
			&link.ShowOutline, &link.ShowShadow, &link.Position, &link.Clicks, &link.IsActive, &link.IsPinned,
			&link.ScheduledAt, &link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &link, nil
}

func (r *LinkRepository) BulkAction(userID string, linkIDs []string, action string) error {
	if len(linkIDs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Build placeholders for IN clause
	placeholders := ""
	args := []interface{}{userID}
	for i, id := range linkIDs {
		if i > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", i+2)
		args = append(args, id)
	}

	var query string
	switch action {
	case "delete":
		query = `DELETE FROM links WHERE id IN (` + placeholders + `) 
		         AND profile_id IN (SELECT id FROM profiles WHERE user_id = $1)`
	case "activate":
		query = `UPDATE links SET is_active = true, updated_at = CURRENT_TIMESTAMP 
		         WHERE id IN (` + placeholders + `) 
		         AND profile_id IN (SELECT id FROM profiles WHERE user_id = $1)`
	case "deactivate":
		query = `UPDATE links SET is_active = false, updated_at = CURRENT_TIMESTAMP 
		         WHERE id IN (` + placeholders + `) 
		         AND profile_id IN (SELECT id FROM profiles WHERE user_id = $1)`
	default:
		return nil
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetChildrenByParentID retrieves all child links belonging to a parent group
func (r *LinkRepository) GetChildrenByParentID(parentID string) ([]Link, error) {
	query := `
		SELECT id, profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio,
		       title, url, description, thumbnail_url, image_shape, layout_type, 
		       image_placement, text_alignment, text_size, show_outline, show_shadow, shadow_x, shadow_y, shadow_blur, show_description, show_text,
		       has_card_background, card_background_color, card_background_opacity, card_border_radius, card_text_color,
		       has_card_border, card_border_color, card_border_style, card_border_width, style,
		       position, clicks, is_active, is_pinned, scheduled_at, expires_at, 
		       created_at, updated_at
		FROM links
		WHERE parent_id = $1
		ORDER BY position ASC
	`

	rows, err := r.db.Query(query, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var children []Link
	for rows.Next() {
		var child Link
		err := rows.Scan(
			&child.ID, &child.ProfileID, &child.ParentID, &child.IsGroup, &child.GroupTitle, &child.GroupLayout, &child.GridColumns, &child.GridAspectRatio,
			&child.Title, &child.URL, &child.Description, &child.ThumbnailURL, &child.ImageShape, &child.LayoutType,
			&child.ImagePlacement, &child.TextAlignment, &child.TextSize, &child.ShowOutline, &child.ShowShadow, &child.ShadowX, &child.ShadowY, &child.ShadowBlur, &child.ShowDescription, &child.ShowText,
			&child.HasCardBackground, &child.CardBackgroundColor, &child.CardBackgroundOpacity, &child.CardBorderRadius, &child.CardTextColor,
			&child.HasCardBorder, &child.CardBorderColor, &child.CardBorderStyle, &child.CardBorderWidth, &child.Style,
			&child.Position, &child.Clicks, &child.IsActive, &child.IsPinned, &child.ScheduledAt, &child.ExpiresAt,
			&child.CreatedAt, &child.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

// CreateGroup creates a new link group
func (r *LinkRepository) CreateGroup(userID string, title string, layout string) (*Link, error) {
	var profileID string
	err := r.db.QueryRow(`SELECT id FROM profiles WHERE user_id = $1`, userID).Scan(&profileID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("profile not found for user. Please create a profile first")
		}
		return nil, err
	}

	// Get max position from top-level items only
	var maxBlockPos, maxLinkPos int
	r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM blocks WHERE profile_id = $1 AND parent_id IS NULL`, profileID).Scan(&maxBlockPos)
	r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM links WHERE profile_id = $1 AND parent_id IS NULL`, profileID).Scan(&maxLinkPos)

	maxPosition := maxBlockPos
	if maxLinkPos > maxPosition {
		maxPosition = maxLinkPos
	}

	var group Link
	query := `
		INSERT INTO links (profile_id, is_group, group_title, group_layout, title, url, position, is_active)
		VALUES ($1, true, $2, $3, $2, '#', $4, true)
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout,
		          title, url, thumbnail_url, layout_type, image_placement, text_alignment,
		          text_size, show_outline, show_shadow, show_description, show_text, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = r.db.QueryRow(query, profileID, title, layout, maxPosition+1).Scan(
		&group.ID, &group.ProfileID, &group.ParentID, &group.IsGroup, &group.GroupTitle, &group.GroupLayout,
		&group.Title, &group.URL, &group.ThumbnailURL, &group.LayoutType,
		&group.ImagePlacement, &group.TextAlignment, &group.TextSize, &group.ShowOutline, &group.ShowShadow, &group.ShowDescription, &group.ShowText,
		&group.Position, &group.Clicks, &group.IsActive, &group.IsPinned, &group.ScheduledAt, &group.ExpiresAt,
		&group.CreatedAt, &group.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	group.Children = []Link{} // Initialize empty children array

	return &group, nil
}

// AddToGroup adds a link to an existing group
func (r *LinkRepository) AddToGroup(userID string, groupID string, data map[string]interface{}) (*Link, error) {
	// Verify group exists, is a group, and belongs to user
	var profileID string
	var isGroup bool
	err := r.db.QueryRow(`
		SELECT l.profile_id, l.is_group FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1 AND p.user_id = $2
	`, groupID, userID).Scan(&profileID, &isGroup)
	if err != nil {
		return nil, fmt.Errorf("group not found")
	}
	if !isGroup {
		return nil, fmt.Errorf("target is not a group")
	}

	// Get max position within group
	var maxPos int
	r.db.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM links WHERE parent_id = $1`, groupID).Scan(&maxPos)

	var link Link
	query := `
		INSERT INTO links (profile_id, parent_id, title, url, description, position, image_placement, text_alignment, text_size, show_outline, show_shadow, show_description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, 'left', 'left', 'M', false, false, true, true)
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout,
		          title, url, description, thumbnail_url, layout_type, image_placement, text_alignment,
		          text_size, show_outline, show_shadow, show_description, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = r.db.QueryRow(query, profileID, groupID, data["title"], data["url"], data["description"], maxPos+1).Scan(
		&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout,
		&link.Title, &link.URL, &link.Description, &link.ThumbnailURL, &link.LayoutType,
		&link.ImagePlacement, &link.TextAlignment, &link.TextSize, &link.ShowOutline, &link.ShowShadow, &link.ShowDescription,
		&link.Position, &link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
		&link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// MoveToGroup moves an existing link into a group
func (r *LinkRepository) MoveToGroup(userID string, linkID string, groupID string) (*Link, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Verify link exists and belongs to user
	var profileID string
	err = tx.QueryRow(`
		SELECT l.profile_id FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1 AND p.user_id = $2 AND l.is_group = false
	`, linkID, userID).Scan(&profileID)
	if err != nil {
		return nil, fmt.Errorf("link not found")
	}

	// Verify group exists and belongs to same profile
	var groupProfileID string
	var isGroup bool
	err = tx.QueryRow(`SELECT profile_id, is_group FROM links WHERE id = $1`, groupID).Scan(&groupProfileID, &isGroup)
	if err != nil {
		return nil, fmt.Errorf("group not found")
	}
	if !isGroup {
		return nil, fmt.Errorf("target is not a group")
	}
	if groupProfileID != profileID {
		return nil, fmt.Errorf("group belongs to different profile")
	}

	// Get max position within group
	var maxPos int
	tx.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM links WHERE parent_id = $1`, groupID).Scan(&maxPos)

	// Move link to group
	var link Link
	query := `
		UPDATE links
		SET parent_id = $2, position = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout,
		          title, url, thumbnail_url, layout_type, image_placement, text_alignment,
		          text_size, show_outline, show_shadow, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = tx.QueryRow(query, linkID, groupID, maxPos+1).Scan(
		&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout,
		&link.Title, &link.URL, &link.ThumbnailURL, &link.LayoutType,
		&link.ImagePlacement, &link.TextAlignment, &link.TextSize, &link.ShowOutline, &link.ShowShadow,
		&link.Position, &link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
		&link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &link, nil
}

// RemoveFromGroup removes a link from its group (makes it top-level)
func (r *LinkRepository) RemoveFromGroup(userID string, linkID string) (*Link, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Verify link exists, belongs to user, and is in a group
	var profileID string
	var parentID *string
	err = tx.QueryRow(`
		SELECT l.profile_id, l.parent_id FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1 AND p.user_id = $2
	`, linkID, userID).Scan(&profileID, &parentID)
	if err != nil {
		return nil, fmt.Errorf("link not found")
	}
	if parentID == nil {
		return nil, fmt.Errorf("link is not in a group")
	}

	// Get max position from top-level items
	var maxBlockPos, maxLinkPos int
	tx.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM blocks WHERE profile_id = $1 AND parent_id IS NULL`, profileID).Scan(&maxBlockPos)
	tx.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM links WHERE profile_id = $1 AND parent_id IS NULL`, profileID).Scan(&maxLinkPos)

	maxPosition := maxBlockPos
	if maxLinkPos > maxPosition {
		maxPosition = maxLinkPos
	}

	// Remove from group
	var link Link
	query := `
		UPDATE links
		SET parent_id = NULL, position = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout,
		          title, url, thumbnail_url, layout_type, image_placement, text_alignment,
		          text_size, show_outline, show_shadow, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = tx.QueryRow(query, linkID, maxPosition+1).Scan(
		&link.ID, &link.ProfileID, &link.ParentID, &link.IsGroup, &link.GroupTitle, &link.GroupLayout,
		&link.Title, &link.URL, &link.ThumbnailURL, &link.LayoutType,
		&link.ImagePlacement, &link.TextAlignment, &link.TextSize, &link.ShowOutline, &link.ShowShadow,
		&link.Position, &link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
		&link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &link, nil
}

// DuplicateGroup duplicates a group and all its children
func (r *LinkRepository) DuplicateGroup(userID string, groupID string) (*Link, error) {
	fmt.Printf("🔄 DuplicateGroup START: userID=%s, groupID=%s\n", userID, groupID)

	tx, err := r.db.Begin()
	if err != nil {
		fmt.Printf("❌ Failed to begin transaction: %v\n", err)
		return nil, err
	}
	defer tx.Rollback()
	fmt.Printf("✅ Transaction started\n")

	// Get original group
	var originalGroup Link
	query := `
		SELECT l.id, l.profile_id, l.parent_id, l.is_group, l.group_title, l.group_layout, l.grid_columns, l.grid_aspect_ratio,
		       l.title, l.url, l.description, l.thumbnail_url, l.image_shape, l.layout_type, 
		       l.image_placement, l.text_alignment, l.text_size, l.show_outline, l.show_shadow, l.show_description, l.show_text,
		       l.position, l.clicks, l.is_active, l.is_pinned, l.scheduled_at, l.expires_at, 
		       l.created_at, l.updated_at
		FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1 AND p.user_id = $2 AND l.is_group = true
	`
	fmt.Printf("📝 Fetching original group...\n")
	err = tx.QueryRow(query, groupID, userID).Scan(
		&originalGroup.ID, &originalGroup.ProfileID, &originalGroup.ParentID, &originalGroup.IsGroup,
		&originalGroup.GroupTitle, &originalGroup.GroupLayout, &originalGroup.GridColumns, &originalGroup.GridAspectRatio, &originalGroup.Title, &originalGroup.URL,
		&originalGroup.Description, &originalGroup.ThumbnailURL, &originalGroup.ImageShape, &originalGroup.LayoutType, &originalGroup.ImagePlacement,
		&originalGroup.TextAlignment, &originalGroup.TextSize, &originalGroup.ShowOutline,
		&originalGroup.ShowShadow, &originalGroup.ShowDescription, &originalGroup.ShowText, &originalGroup.Position, &originalGroup.Clicks, &originalGroup.IsActive,
		&originalGroup.IsPinned, &originalGroup.ScheduledAt, &originalGroup.ExpiresAt,
		&originalGroup.CreatedAt, &originalGroup.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("❌ Failed to fetch original group: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ Original group fetched: ID=%s, Title=%v, Layout=%v\n", originalGroup.ID, originalGroup.GroupTitle, originalGroup.GroupLayout)

	// Get max position from links only (simpler and safer)
	var maxPosition int
	err = tx.QueryRow(`SELECT COALESCE(MAX(position), -1) FROM links WHERE profile_id = $1 AND parent_id IS NULL`, originalGroup.ProfileID).Scan(&maxPosition)
	if err != nil {
		fmt.Printf("❌ Failed to get max position: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ Max position: %d\n", maxPosition)

	// Create duplicate group
	var newGroup Link
	newTitle := *originalGroup.GroupTitle + " (Copy)"
	fmt.Printf("📝 Creating duplicate group with title: %s\n", newTitle)
	insertQuery := `
		INSERT INTO links (profile_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio, image_shape,
		                   text_alignment, text_size, show_outline, show_shadow, show_description, show_text, title, url, position, is_active)
		VALUES ($1, true, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $2, '#', $13, true)
		RETURNING id, profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio,
		          title, url, description, thumbnail_url, image_shape, layout_type, image_placement, text_alignment,
		          text_size, show_outline, show_shadow, show_description, show_text, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = tx.QueryRow(insertQuery, originalGroup.ProfileID, newTitle, originalGroup.GroupLayout,
		originalGroup.GridColumns, originalGroup.GridAspectRatio, originalGroup.ImageShape,
		originalGroup.TextAlignment, originalGroup.TextSize, originalGroup.ShowOutline, originalGroup.ShowShadow, originalGroup.ShowDescription,
		originalGroup.ShowText, maxPosition+1).Scan(
		&newGroup.ID, &newGroup.ProfileID, &newGroup.ParentID, &newGroup.IsGroup, &newGroup.GroupTitle,
		&newGroup.GroupLayout, &newGroup.GridColumns, &newGroup.GridAspectRatio, &newGroup.Title, &newGroup.URL,
		&newGroup.Description, &newGroup.ThumbnailURL, &newGroup.ImageShape, &newGroup.LayoutType, &newGroup.ImagePlacement,
		&newGroup.TextAlignment, &newGroup.TextSize, &newGroup.ShowOutline, &newGroup.ShowShadow, &newGroup.ShowDescription,
		&newGroup.ShowText, &newGroup.Position, &newGroup.Clicks, &newGroup.IsActive, &newGroup.IsPinned,
		&newGroup.ScheduledAt, &newGroup.ExpiresAt, &newGroup.CreatedAt, &newGroup.UpdatedAt,
	)
	if err != nil {
		fmt.Printf("❌ Failed to create duplicate group: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ Duplicate group created: ID=%s\n", newGroup.ID)

	// Get all children of original group - LOAD ALL INTO MEMORY FIRST
	fmt.Printf("📝 Fetching children of original group...\n")
	childrenQuery := `
		SELECT id, title, url, description, thumbnail_url, layout_type, image_placement, text_alignment, 
		       text_size, show_outline, show_shadow, show_description, position, is_active
		FROM links
		WHERE parent_id = $1
		ORDER BY position ASC
	`
	rows, err := tx.Query(childrenQuery, groupID)
	if err != nil {
		fmt.Printf("❌ Failed to query children: %v\n", err)
		return nil, err
	}

	// Load all children into slice first
	type ChildData struct {
		ID              string
		Title           string
		URL             string
		Description     *string
		ThumbnailURL    *string
		LayoutType      string
		ImagePlacement  string
		TextAlignment   string
		TextSize        string
		ShowOutline     bool
		ShowShadow      bool
		ShowDescription bool
		Position        int
		IsActive        bool
	}
	var childrenData []ChildData

	for rows.Next() {
		var child ChildData
		err = rows.Scan(&child.ID, &child.Title, &child.URL, &child.Description, &child.ThumbnailURL, &child.LayoutType,
			&child.ImagePlacement, &child.TextAlignment, &child.TextSize, &child.ShowOutline,
			&child.ShowShadow, &child.ShowDescription, &child.Position, &child.IsActive)
		if err != nil {
			rows.Close()
			fmt.Printf("❌ Failed to scan child: %v\n", err)
			return nil, err
		}
		childrenData = append(childrenData, child)
	}
	rows.Close() // CLOSE ROWS BEFORE INSERT
	fmt.Printf("✅ Loaded %d children from original group\n", len(childrenData))

	// Now insert duplicates
	for i, child := range childrenData {
		fmt.Printf("  📌 Duplicating child #%d: ID=%s, Title=%s\n", i+1, child.ID, child.Title)

		_, err = tx.Exec(`
			INSERT INTO links (profile_id, parent_id, title, url, description, thumbnail_url, layout_type, 
			                   image_placement, text_alignment, text_size, show_outline, show_shadow, show_description,
			                   position, is_active, is_pinned)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, false)
		`, originalGroup.ProfileID, newGroup.ID, child.Title, child.URL, child.Description, child.ThumbnailURL,
			child.LayoutType, child.ImagePlacement, child.TextAlignment, child.TextSize,
			child.ShowOutline, child.ShowShadow, child.ShowDescription, child.Position, child.IsActive)
		if err != nil {
			fmt.Printf("❌ Failed to insert duplicate child #%d: %v\n", i+1, err)
			return nil, err
		}
		fmt.Printf("  ✅ Child #%d duplicated\n", i+1)
	}
	fmt.Printf("✅ Total children duplicated: %d\n", len(childrenData))

	// Load children for the new group BEFORE commit
	fmt.Printf("📝 Loading children for new group...\n")
	childrenQuery2 := `
		SELECT id, profile_id, parent_id, is_group, group_title, group_layout, grid_columns, grid_aspect_ratio,
		       title, url, description, thumbnail_url, image_shape, layout_type, 
		       image_placement, text_alignment, text_size, show_outline, show_shadow, show_description,
		       position, clicks, is_active, is_pinned, scheduled_at, expires_at, 
		       created_at, updated_at
		FROM links
		WHERE parent_id = $1
		ORDER BY is_pinned DESC, position ASC
	`

	rows2, err := tx.Query(childrenQuery2, newGroup.ID)
	if err != nil {
		fmt.Printf("❌ Failed to load children for new group: %v\n", err)
		return nil, err
	}
	defer rows2.Close()

	var children []Link
	loadedCount := 0
	for rows2.Next() {
		loadedCount++
		var child Link
		err := rows2.Scan(
			&child.ID, &child.ProfileID, &child.ParentID, &child.IsGroup, &child.GroupTitle, &child.GroupLayout,
			&child.GridColumns, &child.GridAspectRatio, &child.Title, &child.URL, &child.Description, &child.ThumbnailURL, &child.ImageShape,
			&child.LayoutType, &child.ImagePlacement, &child.TextSize, &child.TextAlignment, &child.ShowOutline, &child.ShowShadow, &child.ShowDescription,
			&child.Position, &child.Clicks, &child.IsActive, &child.IsPinned, &child.ScheduledAt, &child.ExpiresAt,
			&child.CreatedAt, &child.UpdatedAt,
		)
		if err != nil {
			fmt.Printf("❌ Failed to scan loaded child #%d: %v\n", loadedCount, err)
			return nil, err
		}
		children = append(children, child)
	}
	newGroup.Children = children
	fmt.Printf("✅ Loaded %d children for new group\n", loadedCount)

	fmt.Printf("📝 Committing transaction...\n")
	if err = tx.Commit(); err != nil {
		fmt.Printf("❌ Failed to commit transaction: %v\n", err)
		return nil, err
	}
	fmt.Printf("✅ Transaction committed successfully\n")
	fmt.Printf("🎉 DuplicateGroup COMPLETE: New group ID=%s with %d children\n", newGroup.ID, len(newGroup.Children))

	return &newGroup, nil
}

// ReorderGroupLinks reorders links within a group
func (r *LinkRepository) ReorderGroupLinks(userID string, groupID string, linkIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Verify group belongs to user
	var profileID string
	err = tx.QueryRow(`
		SELECT l.profile_id 
		FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1 AND p.user_id = $2 AND l.is_group = true
	`, groupID, userID).Scan(&profileID)
	if err != nil {
		return err
	}

	// Update position for each link
	for i, linkID := range linkIDs {
		_, err := tx.Exec(`
			UPDATE links 
			SET position = $1 
			WHERE id = $2 AND parent_id = $3
		`, i, linkID, groupID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}


// UpdateAllGroupsCardStyles updates card styles for all link groups of a user
func (r *LinkRepository) UpdateAllGroupsCardStyles(userID string, cardStyles map[string]interface{}) error {
	query := `
		UPDATE links l
		SET 
			card_background_color = COALESCE($2, card_background_color),
			card_background_opacity = COALESCE($3, card_background_opacity),
			card_text_color = COALESCE($4, card_text_color),
			card_border_radius = COALESCE($5, card_border_radius),
			show_shadow = COALESCE($6, show_shadow),
			shadow_x = COALESCE($7, shadow_x),
			shadow_y = COALESCE($8, shadow_y),
			shadow_blur = COALESCE($9, shadow_blur),
			has_card_border = COALESCE($10, has_card_border),
			card_border_color = COALESCE($11, card_border_color),
			card_border_width = COALESCE($12, card_border_width),
			has_card_background = COALESCE($13, has_card_background),
			updated_at = CURRENT_TIMESTAMP
		FROM profiles p
		WHERE l.profile_id = p.id 
		  AND p.user_id = $1 
		  AND l.is_group = true
	`

	_, err := r.db.Exec(query,
		userID,
		cardStyles["card_background_color"],
		cardStyles["card_background_opacity"],
		cardStyles["card_text_color"],
		cardStyles["card_border_radius"],
		cardStyles["show_shadow"],
		cardStyles["shadow_x"],
		cardStyles["shadow_y"],
		cardStyles["shadow_blur"],
		cardStyles["has_card_border"],
		cardStyles["card_border_color"],
		cardStyles["card_border_width"],
		cardStyles["has_card_background"],
	)

	return err
}
