package models

import "github.com/google/uuid"

type Order struct {
	ID         uuid.UUID `json:"id"`
	ItemsCount int       `json:"items_count"`
	PackSetup  string    `json:"pack_setup"`
}
