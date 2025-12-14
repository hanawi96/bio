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

// CreateGroup creates a new link group
func (s *LinkService) CreateGroup(userID string, title string, layout string) (*repository.Link, error) {
	return s.linkRepo.CreateGroup(userID, title, layout)
}

// AddToGroup adds a link to an existing group
func (s *LinkService) AddToGroup(userID string, groupID string, data map[string]interface{}) (*repository.Link, error) {
	return s.linkRepo.AddToGroup(userID, groupID, data)
}

// MoveToGroup moves an existing link into a group
func (s *LinkService) MoveToGroup(userID string, linkID string, groupID string) (*repository.Link, error) {
	return s.linkRepo.MoveToGroup(userID, linkID, groupID)
}

// RemoveFromGroup removes a link from its group
func (s *LinkService) RemoveFromGroup(userID string, linkID string) (*repository.Link, error) {
	return s.linkRepo.RemoveFromGroup(userID, linkID)
}



// DuplicateGroup duplicates a group and all its children
func (s *LinkService) DuplicateGroup(userID string, groupID string) (*repository.Link, error) {
	return s.linkRepo.DuplicateGroup(userID, groupID)
}

// ReorderGroupLinks reorders links within a group
func (s *LinkService) ReorderGroupLinks(userID string, groupID string, linkIDs []string) error {
	return s.linkRepo.ReorderGroupLinks(userID, groupID, linkIDs)
}

func (s *LinkService) UpdateAllGroupStyles(userID string, styles map[string]interface{}) error {
	return s.linkRepo.UpdateAllGroupStyles(userID, styles)
}
