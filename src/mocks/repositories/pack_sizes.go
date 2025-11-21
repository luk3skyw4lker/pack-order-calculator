package repositories

import (
	"sync"

	"github.com/luk3skyw4lker/order-pack-calculator/src/database/models"
)

type InMemoryPackSizesRepository struct {
	mu        sync.RWMutex
	packSizes map[string]models.PackSize
}

func NewInMemoryPackSizesRepository() *InMemoryPackSizesRepository {
	return &InMemoryPackSizesRepository{
		packSizes: make(map[string]models.PackSize),
	}
}

func (r *InMemoryPackSizesRepository) CreatePackSize(packSize models.PackSize) (models.PackSize, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.packSizes[packSize.ID.String()] = packSize
	return packSize, nil
}

func (r *InMemoryPackSizesRepository) GetAllPackSizes() ([]models.PackSize, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	packSizes := make([]models.PackSize, 0, len(r.packSizes))
	for _, packSize := range r.packSizes {
		packSizes = append(packSizes, packSize)
	}

	return packSizes, nil
}

func (r *InMemoryPackSizesRepository) UpdatePackSize(packSize models.PackSize) (models.PackSize, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.packSizes[packSize.ID.String()] = packSize
	return packSize, nil
}

// Helper method for testing - clear all pack sizes
func (r *InMemoryPackSizesRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.packSizes = make(map[string]models.PackSize)
}
