package services

import "github.com/luk3skyw4lker/order-pack-calculator/src/database/models"

type PackSizeRepository interface {
	GetAllPackSizes() ([]models.PackSize, error)
	CreatePackSize(packSize models.PackSize) (models.PackSize, error)
	UpdatePackSize(packSize models.PackSize) (models.PackSize, error)
}

type PackSizesService struct {
	repo PackSizeRepository
}

func NewPackSizesService(repo PackSizeRepository) *PackSizesService {
	return &PackSizesService{
		repo: repo,
	}
}

func (s *PackSizesService) GetAllPackSizes() ([]models.PackSize, error) {
	packSizes, err := s.repo.GetAllPackSizes()
	if err != nil {
		return nil, err
	}

	return packSizes, nil
}

func (s *PackSizesService) CreatePackSize(packSize models.PackSize) (models.PackSize, error) {
	createdPackSize, err := s.repo.CreatePackSize(packSize)
	if err != nil {
		return models.PackSize{}, err
	}

	return createdPackSize, nil
}

func (s *PackSizesService) UpdatePackSize(packSize models.PackSize) (models.PackSize, error) {
	updatedPackSize, err := s.repo.UpdatePackSize(packSize)
	if err != nil {
		return models.PackSize{}, err
	}

	return updatedPackSize, nil
}
