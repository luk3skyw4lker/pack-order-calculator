package repositories

import "github.com/luk3skyw4lker/order-pack-calculator/src/database/models"

type Database interface {
	QueryWithScan(query string, dest interface{}, args ...any) error
	Query(query string) error
}

type OrdersRepository struct {
	db Database
}

func NewOrdersRepository(db Database) *OrdersRepository {
	return &OrdersRepository{
		db: db,
	}
}

func (r *OrdersRepository) queryWithScan(query string, args ...any) (models.Order, error) {
	var dest models.Order
	if err := r.db.QueryWithScan(query, &dest, args...); err != nil {
		return models.Order{}, err
	}

	return dest, nil
}

func (r *OrdersRepository) GetAllOrders() ([]models.Order, error) {
	query := "SELECT * FROM orders"

	var dest []models.Order
	if err := r.db.QueryWithScan(query, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *OrdersRepository) SaveOrder(order models.Order) (models.Order, error) {
	query := "INSERT INTO orders (id, items_count, pack_setup) VALUES ($1, $2, $3) RETURNING *"

	return r.queryWithScan(query, order.ID, order.ItemsCount, order.PackSetup)
}

func (r *OrdersRepository) FetchOrder(orderID string) (models.Order, error) {
	query := "SELECT * FROM orders WHERE id = $1"

	return r.queryWithScan(query, orderID)
}
