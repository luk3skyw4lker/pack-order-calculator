package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/luk3skyw4lker/order-pack-calculator/src/database/models"
	"github.com/luk3skyw4lker/order-pack-calculator/src/mocks/repositories"
)

var defaultPackSizes = []int{250, 500, 1000, 2000, 5000}

func TestOrdersService_CreateOrder(t *testing.T) {
	testCases := []struct {
		name          string
		itemsCount    int
		expectedPacks map[int]int // pack size -> count
	}{
		{
			name:       "Single item",
			itemsCount: 1,
			expectedPacks: map[int]int{
				250: 1,
			},
		},
		{
			name:       "Exact match 250",
			itemsCount: 250,
			expectedPacks: map[int]int{
				250: 1,
			},
		},
		{
			name:       "251 items - should use 1x500",
			itemsCount: 251,
			expectedPacks: map[int]int{
				500: 1,
			},
		},
		{
			name:       "501 items - should use 1x500 + 1x250",
			itemsCount: 501,
			expectedPacks: map[int]int{
				500: 1,
				250: 1,
			},
		},
		{
			name:       "Large order 12001",
			itemsCount: 12001,
			expectedPacks: map[int]int{
				5000: 2,
				2000: 1,
				250:  1,
			},
		},
		{
			name:       "Exact match 5000",
			itemsCount: 5000,
			expectedPacks: map[int]int{
				5000: 1,
			},
		},
		{
			name:       "751 items - tricky case",
			itemsCount: 751,
			expectedPacks: map[int]int{
				1000: 1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := repositories.NewInMemoryOrdersRepository()
			packSizesRepo := setupPackSizesRepositoryWithDefaults()
			defer packSizesRepo.Clear()

			service := NewOrdersService(repo, packSizesRepo)

			order, err := service.CreateOrder(tc.itemsCount)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify order was created
			if order.ID == (uuid.UUID{}) {
				t.Error("order ID should not be empty")
			}

			if order.ItemsCount != tc.itemsCount {
				t.Errorf("expected ItemsCount %d, got %d", tc.itemsCount, order.ItemsCount)
			}

			// Verify order was saved in repository
			savedOrder, err := repo.FetchOrder(order.ID.String())
			if err != nil {
				t.Fatalf("order should be saved in repository: %v", err)
			}

			if savedOrder.ID != order.ID {
				t.Errorf("saved order ID mismatch: expected %v, got %v", order.ID, savedOrder.ID)
			}

			// Verify pack combination by calculating it again
			combination := calculatePackCombination(tc.itemsCount, defaultPackSizes)
			if len(combination.Packs) != len(tc.expectedPacks) {
				t.Errorf("expected %d different pack sizes, got %d", len(tc.expectedPacks), len(combination.Packs))
			}

			for packSize, expectedCount := range tc.expectedPacks {
				if gotCount, exists := combination.Packs[packSize]; !exists {
					t.Errorf("expected pack size %d not found in combination", packSize)
				} else if gotCount != expectedCount {
					t.Errorf("pack size %d: expected %d packs, got %d", packSize, expectedCount, gotCount)
				}
			}

			// Verify total items >= order size
			totalItems := 0
			for packSize, count := range combination.Packs {
				totalItems += packSize * count
			}
			if totalItems < tc.itemsCount {
				t.Errorf("total items %d is less than order size %d", totalItems, tc.itemsCount)
			}
		})
	}
}

func TestOrdersService_GetAllOrders(t *testing.T) {
	repo := repositories.NewInMemoryOrdersRepository()
	packSizesRepo := setupPackSizesRepositoryWithDefaults()
	defer packSizesRepo.Clear()

	service := NewOrdersService(repo, packSizesRepo)

	// Initially should be empty
	orders, err := service.GetAllOrders()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 0 {
		t.Errorf("expected 0 orders, got %d", len(orders))
	}

	// Create some orders
	orderCounts := []int{250, 500, 1000}
	for _, count := range orderCounts {
		_, err := service.CreateOrder(count)
		if err != nil {
			t.Fatalf("failed to create order: %v", err)
		}
	}

	// Should now have 3 orders
	orders, err = service.GetAllOrders()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 3 {
		t.Errorf("expected 3 orders, got %d", len(orders))
	}
}

func TestOrdersService_GetOrder(t *testing.T) {
	repo := repositories.NewInMemoryOrdersRepository()
	packSizesRepo := setupPackSizesRepositoryWithDefaults()
	defer packSizesRepo.Clear()
	service := NewOrdersService(repo, packSizesRepo)

	// Create an order
	createdOrder, err := service.CreateOrder(500)
	if err != nil {
		t.Fatalf("failed to create order: %v", err)
	}

	// Fetch the order
	fetchedOrder, err := service.GetOrder(createdOrder.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetchedOrder.ID != createdOrder.ID {
		t.Errorf("expected order ID %v, got %v", createdOrder.ID, fetchedOrder.ID)
	}

	if fetchedOrder.ItemsCount != createdOrder.ItemsCount {
		t.Errorf("expected items count %d, got %d", createdOrder.ItemsCount, fetchedOrder.ItemsCount)
	}

	// Try to fetch non-existent order
	nonExistentID := uuid.New()
	_, err = service.GetOrder(nonExistentID)
	if err == nil {
		t.Error("expected error when fetching non-existent order")
	}
}

func TestCalculatePackCombination(t *testing.T) {
	testCases := []struct {
		name            string
		itemsCount      int
		packSizes       []int
		expectedPacks   map[int]int
		expectedTotal   int
		expectedMinimum bool // whether this should be the minimum items solution
	}{
		{
			name:            "Zero items",
			itemsCount:      0,
			packSizes:       defaultPackSizes,
			expectedPacks:   map[int]int{},
			expectedTotal:   0,
			expectedMinimum: true,
		},
		{
			name:            "Negative items",
			itemsCount:      -1,
			packSizes:       defaultPackSizes,
			expectedPacks:   map[int]int{},
			expectedTotal:   0,
			expectedMinimum: true,
		},
		{
			name:            "1 item needs 1x250",
			itemsCount:      1,
			packSizes:       defaultPackSizes,
			expectedPacks:   map[int]int{250: 1},
			expectedTotal:   250,
			expectedMinimum: true,
		},
		{
			name:            "251 needs 1x500 not 2x250",
			itemsCount:      251,
			packSizes:       defaultPackSizes,
			expectedPacks:   map[int]int{500: 1},
			expectedTotal:   500,
			expectedMinimum: true,
		},
		// odd test case found in the pdf file
		{
			name:            "500000 items with custom pack sizes 23, 31, 53",
			itemsCount:      500000,
			packSizes:       []int{23, 31, 53},
			expectedPacks:   map[int]int{23: 2, 31: 7, 53: 9429},
			expectedTotal:   500000,
			expectedMinimum: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculatePackCombination(tc.itemsCount, tc.packSizes)

			// Check pack counts match
			if len(result.Packs) != len(tc.expectedPacks) {
				t.Errorf("expected %d pack types, got %d", len(tc.expectedPacks), len(result.Packs))
			}

			for packSize, expectedCount := range tc.expectedPacks {
				if gotCount, exists := result.Packs[packSize]; !exists {
					t.Errorf("expected pack size %d not found", packSize)
				} else if gotCount != expectedCount {
					t.Errorf("pack size %d: expected %d, got %d", packSize, expectedCount, gotCount)
				}
			}

			// Calculate total items
			totalItems := 0
			for packSize, count := range result.Packs {
				totalItems += packSize * count
			}

			if totalItems != tc.expectedTotal {
				t.Errorf("expected total items %d, got %d", tc.expectedTotal, totalItems)
			}

			// Verify we meet the order requirement
			if tc.itemsCount > 0 && totalItems < tc.itemsCount {
				t.Errorf("total items %d is less than required %d", totalItems, tc.itemsCount)
			}
		})
	}
}

func setupPackSizesRepositoryWithDefaults() *repositories.InMemoryPackSizesRepository {
	repo := repositories.NewInMemoryPackSizesRepository()
	defaultPackSizes := []int{250, 500, 1000, 2000, 5000}
	for _, size := range defaultPackSizes {
		_, _ = repo.CreatePackSize(models.PackSize{
			ID:   uuid.New(),
			Size: size,
		})
	}

	return repo
}
