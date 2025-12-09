package service

import (
	"database/sql"
	"log"
	"time"

	"github.com/yourusername/linkbio/repository"
)

type SchedulerService struct {
	db       *sql.DB
	linkRepo *repository.LinkRepository
	ticker   *time.Ticker
	done     chan bool
}

func NewSchedulerService(db *sql.DB) *SchedulerService {
	return &SchedulerService{
		db:       db,
		linkRepo: repository.NewLinkRepository(db),
		done:     make(chan bool),
	}
}

// Start begins the scheduler that runs every minute
func (s *SchedulerService) Start() {
	log.Println("📅 Scheduler service started")
	s.ticker = time.NewTicker(1 * time.Minute)

	go func() {
		// Run immediately on start
		s.processScheduledLinks()

		for {
			select {
			case <-s.ticker.C:
				s.processScheduledLinks()
			case <-s.done:
				log.Println("📅 Scheduler service stopped")
				return
			}
		}
	}()
}

// Stop stops the scheduler
func (s *SchedulerService) Stop() {
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
}

// processScheduledLinks checks and updates link statuses based on schedule
func (s *SchedulerService) processScheduledLinks() {
	now := time.Now()

	// Activate links that should be published
	activatedCount, err := s.activateScheduledLinks(now)
	if err != nil {
		log.Printf("❌ Error activating scheduled links: %v", err)
	} else if activatedCount > 0 {
		log.Printf("✅ Activated %d scheduled link(s)", activatedCount)
	}

	// Deactivate links that have expired
	deactivatedCount, err := s.deactivateExpiredLinks(now)
	if err != nil {
		log.Printf("❌ Error deactivating expired links: %v", err)
	} else if deactivatedCount > 0 {
		log.Printf("⏰ Deactivated %d expired link(s)", deactivatedCount)
	}
}

// activateScheduledLinks activates links whose scheduled_at time has passed
func (s *SchedulerService) activateScheduledLinks(now time.Time) (int, error) {
	query := `
		UPDATE links
		SET is_active = true, updated_at = NOW()
		WHERE scheduled_at IS NOT NULL
		  AND scheduled_at <= $1
		  AND is_active = false
	`

	result, err := s.db.Exec(query, now)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// deactivateExpiredLinks deactivates links whose expires_at time has passed
func (s *SchedulerService) deactivateExpiredLinks(now time.Time) (int, error) {
	query := `
		UPDATE links
		SET is_active = false, updated_at = NOW()
		WHERE expires_at IS NOT NULL
		  AND expires_at <= $1
		  AND is_active = true
	`

	result, err := s.db.Exec(query, now)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetUpcomingSchedules returns links with upcoming schedules
func (s *SchedulerService) GetUpcomingSchedules(profileID string, limit int) ([]repository.Link, error) {
	query := `
		SELECT id, profile_id, title, url, thumbnail_url, layout_type, position, clicks,
		       is_active, is_pinned, scheduled_at, expires_at, created_at, updated_at
		FROM links
		WHERE profile_id = $1
		  AND (scheduled_at IS NOT NULL OR expires_at IS NOT NULL)
		  AND (scheduled_at > NOW() OR expires_at > NOW())
		ORDER BY COALESCE(scheduled_at, expires_at) ASC
		LIMIT $2
	`

	rows, err := s.db.Query(query, profileID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []repository.Link
	for rows.Next() {
		var link repository.Link
		err := rows.Scan(
			&link.ID, &link.ProfileID, &link.Title, &link.URL, &link.ThumbnailURL,
			&link.LayoutType, &link.Position, &link.Clicks, &link.IsActive, &link.IsPinned,
			&link.ScheduledAt, &link.ExpiresAt, &link.CreatedAt, &link.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	return links, nil
}
