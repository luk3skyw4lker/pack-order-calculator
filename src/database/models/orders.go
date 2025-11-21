package models

import "github.com/google/uuid"

type Order struct {
	ID         uuid.UUID
	ItemsCount int
	PackCount  int
}
