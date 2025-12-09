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
		SELECT l.id, l.profile_id, l.title, l.url, l.thumbnail_url, l.layout_type, l.position, l.clicks,
		       l.is_active, l.is_pinned, l.scheduled_at, l.expires_at, l.created_at, l.updated_at
		FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE p.user_id = $1
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
		query += " ORDER BY l.is_pinned DESC, l.position ASC"
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
			&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.ThumbnailURL, &link.LayoutType, &link.Position,
			&link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
			&link.CreatedAt, &link.UpdatedAt,
		)
		if err != nil {
			return nil, err
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
		INSERT INTO links (profile_id, title, url, position)
		VALUES ($1, $2, $3, $4)
		RETURNING id, profile_id, title, url, thumbnail_url, layout_type, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = r.db.QueryRow(query, profileID, data["title"], data["url"], maxPosition+1).Scan(
		&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.ThumbnailURL, &link.LayoutType, &link.Position,
		&link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
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
		    layout_type = COALESCE($5, layout_type),
		    is_active = COALESCE($6, is_active),
		    scheduled_at = CASE WHEN $7::text IS NOT NULL THEN $7::timestamp ELSE scheduled_at END,
		    expires_at = CASE WHEN $8::text IS NOT NULL THEN $8::timestamp ELSE expires_at END,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, profile_id, title, url, thumbnail_url, layout_type, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	var link Link
	err := r.db.QueryRow(query, linkID, data["title"], data["url"], data["thumbnail_url"], data["layout_type"], data["is_active"], data["scheduled_at"], data["expires_at"]).Scan(
		&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.ThumbnailURL, &link.LayoutType, &link.Position,
		&link.Clicks, &link.IsActive, &link.IsPinned, &link.ScheduledAt, &link.ExpiresAt,
		&link.CreatedAt, &link.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *LinkRepository) Delete(linkID string) error {
	_, err := r.db.Exec(`DELETE FROM links WHERE id = $1`, linkID)
	return err
}

func (r *LinkRepository) Reorder(userID string, linkIDs []string) error {
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

	// Update link positions
	for i, linkID := range linkIDs {
		_, err := tx.Exec(`
			UPDATE links 
			SET position = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 AND profile_id = $3
		`, i, linkID, profileID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
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
		SELECT l.id, l.profile_id, l.title, l.url, l.thumbnail_url, l.layout_type, l.position, l.clicks,
		       l.is_active, l.is_pinned, l.scheduled_at, l.expires_at, l.created_at, l.updated_at
		FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE l.id = $1 AND p.user_id = $2
	`
	err := r.db.QueryRow(query, linkID, userID).Scan(
		&original.ID, &original.ProfileID, &original.Title, &original.URL, &original.ThumbnailURL,
		&original.LayoutType, &original.Position, &original.Clicks, &original.IsActive, &original.IsPinned,
		&original.ScheduledAt, &original.ExpiresAt, &original.CreatedAt, &original.UpdatedAt,
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
		INSERT INTO links (profile_id, title, url, thumbnail_url, layout_type, position, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, false)
		RETURNING id, profile_id, title, url, thumbnail_url, layout_type, position, clicks, is_active, is_pinned,
		          scheduled_at, expires_at, created_at, updated_at
	`
	newTitle := original.Title + " (Copy)"
	err = r.db.QueryRow(
		insertQuery,
		original.ProfileID,
		newTitle,
		original.URL,
		original.ThumbnailURL,
		original.LayoutType,
		maxPosition+1,
	).Scan(
		&duplicate.ID, &duplicate.ProfileID, &duplicate.Title, &duplicate.URL, &duplicate.ThumbnailURL,
		&duplicate.LayoutType, &duplicate.Position, &duplicate.Clicks, &duplicate.IsActive, &duplicate.IsPinned,
		&duplicate.ScheduledAt, &duplicate.ExpiresAt, &duplicate.CreatedAt, &duplicate.UpdatedAt,
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

	// Get profile ID and current pin status
	var profileID string
	var currentPinned bool
	err = tx.QueryRow(`
		SELECT p.id, COALESCE(l.is_pinned, false)
		FROM profiles p
		LEFT JOIN links l ON l.id = $2 AND l.profile_id = p.id
		WHERE p.user_id = $1
	`, userID, linkID).Scan(&profileID, &currentPinned)
	if err != nil {
		return nil, err
	}

	var link Link
	if currentPinned {
		// Unpin this link
		query := `
			UPDATE links 
			SET is_pinned = false, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND profile_id = $2
			RETURNING id, profile_id, title, url, thumbnail_url, layout_type, position, clicks, is_active, is_pinned,
			          scheduled_at, expires_at, created_at, updated_at
		`
		err = tx.QueryRow(query, linkID, profileID).Scan(
			&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.ThumbnailURL,
			&link.LayoutType, &link.Position, &link.Clicks, &link.IsActive, &link.IsPinned,
			&link.ScheduledAt, &link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
		)
	} else {
		// Unpin all other links first
		_, err = tx.Exec(`UPDATE links SET is_pinned = false WHERE profile_id = $1`, profileID)
		if err != nil {
			return nil, err
		}

		// Pin this link
		query := `
			UPDATE links 
			SET is_pinned = true, position = 0, updated_at = CURRENT_TIMESTAMP
			WHERE id = $1 AND profile_id = $2
			RETURNING id, profile_id, title, url, thumbnail_url, layout_type, position, clicks, is_active, is_pinned,
			          scheduled_at, expires_at, created_at, updated_at
		`
		err = tx.QueryRow(query, linkID, profileID).Scan(
			&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.ThumbnailURL,
			&link.LayoutType, &link.Position, &link.Clicks, &link.IsActive, &link.IsPinned,
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
