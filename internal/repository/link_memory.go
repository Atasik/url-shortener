package repository

import (
	"link-shortener/internal/domain"
	"sync"
)

type LinkMemoryRepo struct {
	lastID      int64
	mu          sync.RWMutex
	shortToLong map[int64]domain.Link
	longToShort map[string]int64
}

func NewLinkMemoryRepo() *LinkMemoryRepo {
	return &LinkMemoryRepo{
		shortToLong: map[int64]domain.Link{},
		longToShort: map[string]int64{},
	}
}

func (repo *LinkMemoryRepo) AddOriginalURL(link domain.Link) (int64, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	id, ok := repo.longToShort[link.OriginalURL]
	if !ok {
		repo.lastID++
		repo.shortToLong[repo.lastID] = link
		repo.longToShort[link.OriginalURL] = repo.lastID
		return repo.lastID, nil
	}
	return id, nil
}

func (repo *LinkMemoryRepo) GetOrginalURL(link domain.Link) (string, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	res, ok := repo.shortToLong[link.ID]
	if !ok {
		return "", domain.ErrNotFound
	}
	return res.OriginalURL, nil
}
