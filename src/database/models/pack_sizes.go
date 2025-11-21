package models

import "github.com/google/uuid"

type PackSize struct {
	ID   uuid.UUID `json:"id"`
	Size int       `json:"size"`
}
