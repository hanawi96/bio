package repository

import (
	"database/sql"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) *LinkRepository {
	return &LinkRepository{db: db}
}

func (r *LinkRepository) GetByUserID(userID string) ([]Link, error) {
	query := `
		SELECT l.id, l.profile_id, l.title, l.url, l.position, l.clicks,
		       l.is_active, l.scheduled_at, l.expires_at, l.created_at, l.updated_at
		FROM links l
		JOIN profiles p ON l.profile_id = p.id
		WHERE p.user_id = $1
		ORDER BY l.position ASC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []Link
	for rows.Next() {
		var link Link
		err := rows.Scan(
			&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.Position,
			&link.Clicks, &link.IsActive, &link.ScheduledAt, &link.ExpiresAt,
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

	var link Link
	query := `
		INSERT INTO links (profile_id, title, url, position)
		VALUES ($1, $2, $3, COALESCE($4, 0))
		RETURNING id, profile_id, title, url, position, clicks, is_active,
		          scheduled_at, expires_at, created_at, updated_at
	`
	err = r.db.QueryRow(query, profileID, data["title"], data["url"], data["position"]).Scan(
		&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.Position,
		&link.Clicks, &link.IsActive, &link.ScheduledAt, &link.ExpiresAt,
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
		    is_active = COALESCE($4, is_active),
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, profile_id, title, url, position, clicks, is_active,
		          scheduled_at, expires_at, created_at, updated_at
	`
	var link Link
	err := r.db.QueryRow(query, linkID, data["title"], data["url"], data["is_active"]).Scan(
		&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.Position,
		&link.Clicks, &link.IsActive, &link.ScheduledAt, &link.ExpiresAt,
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

	for i, linkID := range linkIDs {
		_, err := tx.Exec(`
			UPDATE links 
			SET position = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id = $2 
			AND profile_id IN (SELECT id FROM profiles WHERE user_id = $3)
		`, i, linkID, userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
