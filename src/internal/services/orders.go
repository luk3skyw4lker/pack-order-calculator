package services

import (
	"database/sql"
	"errors"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/luk3skyw4lker/order-pack-calculator/src/database/models"
	"github.com/luk3skyw4lker/order-pack-calculator/src/payload"
)

type PackCombinationResult struct {
	Packs      map[int]int // key: number of packs, value: pack size
	TotalPacks int
}

type OrdersRepository interface {
	GetAllOrders() ([]models.Order, error)
	SaveOrder(order models.Order) (models.Order, error)
	FetchOrder(orderID string) (models.Order, error)
}

type OrdersService struct {
	ordersRepository OrdersRepository
	packSizesRepo    PackSizeRepository
}

func NewOrdersService(ordersRepository OrdersRepository, packSizesRepo PackSizeRepository) *OrdersService {
	ordersService := &OrdersService{
		ordersRepository: ordersRepository,
		packSizesRepo:    packSizesRepo,
	}

	return ordersService
}

func (s *OrdersService) GetAllOrders() ([]models.Order, error) {
	orders, err := s.ordersRepository.GetAllOrders()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return []models.Order{}, payload.ErrOrderNotFound
		}

		return []models.Order{}, err
	}

	return orders, nil
}

func (s *OrdersService) GetOrder(orderID uuid.UUID) (models.Order, error) {
	order, err := s.ordersRepository.FetchOrder(orderID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return models.Order{}, payload.ErrOrderNotFound
		}

		return models.Order{}, err
	}

	return order, nil
}

func (s *OrdersService) CreateOrder(itemsCount int) (models.Order, error) {
	packSizes, err := s.packSizesRepo.GetAllPackSizes()
	if err != nil {
		return models.Order{}, err
	}

	combination := calculatePackCombination(itemsCount, formatPackSizes(packSizes))

	order := models.Order{
		ID:         uuid.New(),
		ItemsCount: itemsCount,
		PackSetup:  formatPackSetup(combination.Packs),
	}

	return s.ordersRepository.SaveOrder(order)
}

func formatPackSizes(packSizes []models.PackSize) []int {
	sizes := make([]int, len(packSizes))
	for i, ps := range packSizes {
		sizes[i] = ps.Size
	}

	return sizes
}

// We save the pack setup as a formatted string like "2x500, 1x1000"
// which is not optimal for querying but works for demonstration purposes.
func formatPackSetup(packs map[int]int) string {
	result := ""
	for size, count := range packs {
		if result != "" {
			result += ", "
		}
		result += fmt.Sprintf("%dx%d", count, size)
	}

	return result
}

// We're using a dynamic programming with greedy optimization approach
// to find the optimal pack combination that meets or exceeds
// the itemsCount with the least number of packs and items.
func calculatePackCombination(itemsCount int, packSizes []int) PackCombinationResult {
	if itemsCount <= 0 {
		return PackCombinationResult{Packs: make(map[int]int)}
	}

	// Find the maximum reasonable target (itemsCount + largest pack - 1)
	maxTarget := itemsCount + packSizes[len(packSizes)-1]

	dp, parent := buildDPAndParent(maxTarget, packSizes)

	bestTarget := findBestTarget(dp, itemsCount, maxTarget)

	if bestTarget == -1 {
		// Fallback (shouldn't happen with these pack sizes)
		return PackCombinationResult{Packs: make(map[int]int), TotalPacks: 0}
	}

	// Backtrack to reconstruct the solution
	packs := make(map[int]int)
	current := bestTarget
	for current > 0 && parent[current] != -1 {
		packSize := parent[current]
		packs[packSize]++
		current -= packSize
	}

	return PackCombinationResult{
		Packs: packs,
	}
}

func buildDPAndParent(maxTarget int, packSizes []int) ([]int, []int) {
	// DP arrays: dp[i] = minimum packs needed to make exactly i items
	dp := make([]int, maxTarget+1)
	parent := make([]int, maxTarget+1) // Track which pack size was used

	// Initialize with impossible values
	for i := range dp {
		dp[i] = math.MaxInt32
		parent[i] = -1
	}
	dp[0] = 0

	// Fill DP table
	for i := 1; i <= maxTarget; i++ {
		for _, pack := range packSizes {
			if i >= pack && dp[i-pack] != math.MaxInt32 {
				if dp[i-pack]+1 < dp[i] {
					dp[i] = dp[i-pack] + 1
					parent[i] = pack
				}
			}
		}
	}

	return dp, parent
}

func findBestTarget(dp []int, start, end int) int {
	// Find the minimum valid target >= itemsCount
	bestTarget := -1

	for target := start; target <= end; target++ {
		if dp[target] != math.MaxInt32 {
			bestTarget = target
			break
		}
	}

	return bestTarget
}
