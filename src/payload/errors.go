package payload

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
)

type ErrorResponse struct {
	Message string `json:"message"`
}
