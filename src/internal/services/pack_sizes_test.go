package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/luk3skyw4lker/order-pack-calculator/src/database/models"
	"github.com/luk3skyw4lker/order-pack-calculator/src/mocks/repositories"
)

func TestPackSizesService_GetAllPackSizes(t *testing.T) {
	repo := repositories.NewInMemoryPackSizesRepository()
	service := NewPackSizesService(repo)

	// Initially should be empty
	packSizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(packSizes) != 0 {
		t.Errorf("expected 0 pack sizes, got %d", len(packSizes))
	}

	// Add some pack sizes
	sizes := []int{250, 500, 1000}
	for _, size := range sizes {
		_, err := repo.CreatePackSize(models.PackSize{
			ID:   uuid.New(),
			Size: size,
		})
		if err != nil {
			t.Fatalf("failed to create pack size: %v", err)
		}
	}

	// Should now have 3 pack sizes
	packSizes, err = service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(packSizes) != 3 {
		t.Errorf("expected 3 pack sizes, got %d", len(packSizes))
	}

	// Verify all sizes are present
	sizeMap := make(map[int]bool)
	for _, ps := range packSizes {
		sizeMap[ps.Size] = true
	}
	for _, expectedSize := range sizes {
		if !sizeMap[expectedSize] {
			t.Errorf("expected to find pack size %d, but it was not present", expectedSize)
		}
	}
}

func TestPackSizesService_CreatePackSize(t *testing.T) {
	testCases := []struct {
		name     string
		packSize models.PackSize
	}{
		{
			name: "Create pack size 250",
			packSize: models.PackSize{
				ID:   uuid.New(),
				Size: 250,
			},
		},
		{
			name: "Create pack size 500",
			packSize: models.PackSize{
				ID:   uuid.New(),
				Size: 500,
			},
		},
		{
			name: "Create pack size 1000",
			packSize: models.PackSize{
				ID:   uuid.New(),
				Size: 1000,
			},
		},
		{
			name: "Create pack size with large value",
			packSize: models.PackSize{
				ID:   uuid.New(),
				Size: 10000,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := repositories.NewInMemoryPackSizesRepository()
			service := NewPackSizesService(repo)

			createdPackSize, err := service.CreatePackSize(tc.packSize)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify the returned pack size
			if createdPackSize.ID != tc.packSize.ID {
				t.Errorf("expected ID %v, got %v", tc.packSize.ID, createdPackSize.ID)
			}
			if createdPackSize.Size != tc.packSize.Size {
				t.Errorf("expected size %d, got %d", tc.packSize.Size, createdPackSize.Size)
			}

			// Verify it was actually saved in the repository
			allPackSizes, err := repo.GetAllPackSizes()
			if err != nil {
				t.Fatalf("failed to get all pack sizes: %v", err)
			}

			found := false
			for _, ps := range allPackSizes {
				if ps.ID == tc.packSize.ID && ps.Size == tc.packSize.Size {
					found = true
					break
				}
			}
			if !found {
				t.Error("created pack size was not found in repository")
			}
		})
	}
}

func TestPackSizesService_UpdatePackSize(t *testing.T) {
	repo := repositories.NewInMemoryPackSizesRepository()
	service := NewPackSizesService(repo)

	// Create initial pack size
	initialPackSize := models.PackSize{
		ID:   uuid.New(),
		Size: 250,
	}
	_, err := repo.CreatePackSize(initialPackSize)
	if err != nil {
		t.Fatalf("failed to create initial pack size: %v", err)
	}

	// Update the pack size
	updatedPackSize := models.PackSize{
		ID:   initialPackSize.ID,
		Size: 500,
	}
	result, err := service.UpdatePackSize(updatedPackSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the returned pack size
	if result.ID != updatedPackSize.ID {
		t.Errorf("expected ID %v, got %v", updatedPackSize.ID, result.ID)
	}
	if result.Size != updatedPackSize.Size {
		t.Errorf("expected size %d, got %d", updatedPackSize.Size, result.Size)
	}

	// Verify it was actually updated in the repository
	allPackSizes, err := repo.GetAllPackSizes()
	if err != nil {
		t.Fatalf("failed to get all pack sizes: %v", err)
	}

	found := false
	for _, ps := range allPackSizes {
		if ps.ID == updatedPackSize.ID {
			found = true
			if ps.Size != updatedPackSize.Size {
				t.Errorf("expected updated size %d, got %d", updatedPackSize.Size, ps.Size)
			}
			break
		}
	}
	if !found {
		t.Error("updated pack size was not found in repository")
	}
}

func TestPackSizesService_CreateMultiplePackSizes(t *testing.T) {
	repo := repositories.NewInMemoryPackSizesRepository()
	service := NewPackSizesService(repo)

	sizes := []int{250, 500, 1000, 2000, 5000}
	createdIDs := make([]uuid.UUID, 0, len(sizes))

	// Create multiple pack sizes
	for _, size := range sizes {
		packSize := models.PackSize{
			ID:   uuid.New(),
			Size: size,
		}
		created, err := service.CreatePackSize(packSize)
		if err != nil {
			t.Fatalf("failed to create pack size %d: %v", size, err)
		}
		createdIDs = append(createdIDs, created.ID)
	}

	// Verify all were created
	allPackSizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(allPackSizes) != len(sizes) {
		t.Errorf("expected %d pack sizes, got %d", len(sizes), len(allPackSizes))
	}

	// Verify all sizes and IDs are present
	foundSizes := make(map[int]bool)
	foundIDs := make(map[uuid.UUID]bool)
	for _, ps := range allPackSizes {
		foundSizes[ps.Size] = true
		foundIDs[ps.ID] = true
	}

	for _, size := range sizes {
		if !foundSizes[size] {
			t.Errorf("expected to find pack size %d, but it was not present", size)
		}
	}

	for _, id := range createdIDs {
		if !foundIDs[id] {
			t.Errorf("expected to find pack size with ID %v, but it was not present", id)
		}
	}
}

func TestPackSizesService_UpdateNonExistentPackSize(t *testing.T) {
	repo := repositories.NewInMemoryPackSizesRepository()
	service := NewPackSizesService(repo)

	// Try to update a pack size that doesn't exist
	nonExistentPackSize := models.PackSize{
		ID:   uuid.New(),
		Size: 999,
	}

	// This should still work as the in-memory repository doesn't validate existence
	// It will just create a new entry
	result, err := service.UpdatePackSize(nonExistentPackSize)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != nonExistentPackSize.ID {
		t.Errorf("expected ID %v, got %v", nonExistentPackSize.ID, result.ID)
	}

	// Verify it exists now
	allPackSizes, err := service.GetAllPackSizes()
	if err != nil {
		t.Fatalf("failed to get all pack sizes: %v", err)
	}

	found := false
	for _, ps := range allPackSizes {
		if ps.ID == nonExistentPackSize.ID {
			found = true
			break
		}
	}
	if !found {
		t.Error("pack size should exist after update")
	}
}
