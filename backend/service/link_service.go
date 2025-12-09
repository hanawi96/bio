package service

import (
	"github.com/yourusername/linkbio/repository"
)

type LinkService struct {
	linkRepo *repository.LinkRepository
}

func NewLinkService(linkRepo *repository.LinkRepository) *LinkService {
	return &LinkService{linkRepo: linkRepo}
}

func (s *LinkService) GetByUserID(userID string) ([]repository.Link, error) {
	return s.linkRepo.GetByUserID(userID)
}

func (s *LinkService) GetByUserIDWithFilters(userID, search, status, layoutType, sortBy string) ([]repository.Link, error) {
	return s.linkRepo.GetByUserIDWithFilters(userID, search, status, layoutType, sortBy)
}

func (s *LinkService) Create(userID string, data map[string]interface{}) (*repository.Link, error) {
	return s.linkRepo.Create(userID, data)
}

func (s *LinkService) Update(linkID string, data map[string]interface{}) (*repository.Link, error) {
	return s.linkRepo.Update(linkID, data)
}

func (s *LinkService) Delete(linkID string) error {
	return s.linkRepo.Delete(linkID)
}

func (s *LinkService) ReorderLinks(userID string, linkIDs []string) error {
	return s.linkRepo.Reorder(userID, linkIDs)
}

func (s *LinkService) Duplicate(userID string, linkID string) (*repository.Link, error) {
	return s.linkRepo.Duplicate(userID, linkID)
}

func (s *LinkService) BulkAction(userID string, linkIDs []string, action string) error {
	return s.linkRepo.BulkAction(userID, linkIDs, action)
}

func (s *LinkService) TogglePin(userID string, linkID string) (*repository.Link, error) {
	return s.linkRepo.TogglePin(userID, linkID)
}

func (s *LinkService) ReorderWithBlocks(userID string, items []map[string]interface{}) error {
	return s.linkRepo.ReorderWithBlocks(userID, items)
}
