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
