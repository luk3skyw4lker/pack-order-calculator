package payload

type CreateOrder struct {
	ItemsCount int `json:"items_count" validate:"required,gt=0"`
}
