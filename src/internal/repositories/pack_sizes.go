package repositories

import "github.com/luk3skyw4lker/order-pack-calculator/src/database/models"

type PackSizesRepository struct {
	db Database
}

func NewPackSizesRepository(db Database) *PackSizesRepository {
	return &PackSizesRepository{
		db: db,
	}
}

func (r *PackSizesRepository) CreatePackSize(packSize models.PackSize) (models.PackSize, error) {
	query := "INSERT INTO pack_sizes (id, size) VALUES ($1, $2) RETURNING *"

	var dest models.PackSize
	if err := r.db.QueryWithScan(query, &dest, packSize.ID, packSize.Size); err != nil {
		return models.PackSize{}, err
	}

	return dest, nil
}

func (r *PackSizesRepository) GetAllPackSizes() ([]models.PackSize, error) {
	query := "SELECT * FROM pack_sizes"

	var dest []models.PackSize
	if err := r.db.QueryWithScan(query, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (r *PackSizesRepository) UpdatePackSize(packSize models.PackSize) (models.PackSize, error) {
	query := "UPDATE pack_sizes SET size = $1 WHERE id = $2 RETURNING *"

	var dest models.PackSize
	if err := r.db.QueryWithScan(query, &dest, packSize.Size, packSize.ID); err != nil {
		return models.PackSize{}, err
	}

	return dest, nil
}
