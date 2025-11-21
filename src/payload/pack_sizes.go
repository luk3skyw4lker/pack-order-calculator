package payload

import "github.com/google/uuid"

type CreatePackSize struct {
	Size int `json:"size" validate:"required,gt=0"`
}

type UpdatePackSize struct {
	ID   uuid.UUID `json:"id" validate:"required"`
	Size int       `json:"size" validate:"required,gt=0"`
}
