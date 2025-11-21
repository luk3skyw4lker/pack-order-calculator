package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/luk3skyw4lker/order-pack-calculator/src/database/models"
	"github.com/luk3skyw4lker/order-pack-calculator/src/payload"
	"github.com/luk3skyw4lker/order-pack-calculator/src/utils"
)

type PackSizesService interface {
	GetAllPackSizes() ([]models.PackSize, error)
	CreatePackSize(packSize models.PackSize) (models.PackSize, error)
	UpdatePackSize(packSize models.PackSize) (models.PackSize, error)
}

type PackSizesHandler struct {
	service PackSizesService
}

func NewPackSizesHandler(service PackSizesService) *PackSizesHandler {
	return &PackSizesHandler{
		service: service,
	}
}

// CreatePackSize godoc
//
//	@Summary		Create a new pack size
//	@Description	Add a new pack size to the system
//	@Tags			PackSizes
//	@Accept			json
//	@Produce		json
//	@Param			packSize	body		payload.CreatePackSize	true	"The pack size to create"
//	@Success		201			{object}	models.PackSize
//	@Failure		400			{object}	payload.ErrorResponse
//	@Failure		500			{object}	payload.ErrorResponse
//	@Router			/pack-sizes [post]
func (h *PackSizesHandler) CreatePackSize(ctx fiber.Ctx) error {
	input, err := utils.UnmarshalRequest[payload.CreatePackSize](ctx)
	if err != nil {
		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(payload.ErrorResponse{Message: "invalid request body"})
	}

	createdPackSize, err := h.service.CreatePackSize(models.PackSize{
		ID:   uuid.New(),
		Size: input.Size,
	})
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(payload.ErrorResponse{Message: "failed to create pack size"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(createdPackSize)
}

// UpdatePackSize godoc
//
//	@Summary		Update an existing pack size
//	@Description	Update the size of an existing pack size
//	@Tags			PackSizes
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"The ID of the pack size to update"
//	@Param			packSize	body		payload.UpdatePackSize	true	"The updated pack size data"
//	@Success		200			{object}	models.PackSize
//	@Failure		400			{object}	payload.ErrorResponse
//	@Failure		500			{object}	payload.ErrorResponse
//	@Router			/pack-sizes/{pack_size_id} [put]
func (h *PackSizesHandler) UpdatePackSize(ctx fiber.Ctx) error {
	input, err := utils.UnmarshalRequest[payload.UpdatePackSize](ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(payload.ErrorResponse{Message: "invalid request body"})
	}

	packSizeID, err := uuid.Parse(ctx.Params("pack_size_id"))
	if err != nil {
		log.Error("invalid pack size ID:", err)

		return ctx.
			Status(fiber.StatusBadRequest).
			JSON(payload.ErrorResponse{Message: "invalid pack size ID"})
	}

	updatedPackSize, err := h.service.UpdatePackSize(models.PackSize{
		ID:   packSizeID,
		Size: input.Size,
	})
	if err != nil {
		log.Error("failed to update pack size:", err)

		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(payload.ErrorResponse{Message: "failed to update pack size"})
	}

	return ctx.Status(fiber.StatusOK).JSON(updatedPackSize)
}

// GetAllPackSizes godoc
//
//	@Summary		Get all pack sizes
//	@Description	Retrieve a list of all available pack sizes
//	@Tags			PackSizes
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.PackSize
//	@Failure		500	{object}	payload.ErrorResponse
//	@Router			/pack-sizes [get]
func (h *PackSizesHandler) GetAllPackSizes(ctx fiber.Ctx) error {
	packSizes, err := h.service.GetAllPackSizes()
	if err != nil {
		return ctx.
			Status(fiber.StatusInternalServerError).
			JSON(payload.ErrorResponse{Message: "failed to retrieve pack sizes"})
	}

	return ctx.Status(fiber.StatusOK).JSON(packSizes)
}
