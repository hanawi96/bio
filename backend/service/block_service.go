package service

import "github.com/yourusername/linkbio/repository"

type BlockService struct {
	repo *repository.BlockRepository
}

func NewBlockService(repo *repository.BlockRepository) *BlockService {
	return &BlockService{repo: repo}
}

func (s *BlockService) GetBlocks(userID string) ([]repository.Block, error) {
	return s.repo.GetByUserID(userID)
}

func (s *BlockService) CreateBlock(userID string, data map[string]interface{}) (*repository.Block, error) {
	return s.repo.Create(userID, data)
}

func (s *BlockService) UpdateBlock(blockID string, data map[string]interface{}) (*repository.Block, error) {
	return s.repo.Update(blockID, data)
}

func (s *BlockService) DeleteBlock(blockID string) error {
	return s.repo.Delete(blockID)
}

func (s *BlockService) ReorderBlocks(userID string, blockIDs []string) error {
	return s.repo.Reorder(userID, blockIDs)
}

func (s *BlockService) BulkDeleteBlocks(userID string, blockIDs []string) error {
	return s.repo.BulkDelete(userID, blockIDs)
}
