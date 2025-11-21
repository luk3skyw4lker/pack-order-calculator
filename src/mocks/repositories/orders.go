package repositories

import (
	"errors"
	"sync"

	"github.com/luk3skyw4lker/order-pack-calculator/src/database/models"
)

type InMemoryOrdersRepository struct {
	mu     sync.RWMutex
	orders map[string]models.Order
}

func NewInMemoryOrdersRepository() *InMemoryOrdersRepository {
	return &InMemoryOrdersRepository{
		orders: make(map[string]models.Order),
	}
}

func (r *InMemoryOrdersRepository) GetAllOrders() ([]models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	orders := make([]models.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *InMemoryOrdersRepository) SaveOrder(order models.Order) (models.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID.String()] = order
	return order, nil
}

func (r *InMemoryOrdersRepository) FetchOrder(orderID string) (models.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, exists := r.orders[orderID]
	if !exists {
		return models.Order{}, errors.New("order not found")
	}

	return order, nil
}

// Helper method for testing - clear all orders
func (r *InMemoryOrdersRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders = make(map[string]models.Order)
}

// Helper method for testing - get count of orders
func (r *InMemoryOrdersRepository) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.orders)
}
