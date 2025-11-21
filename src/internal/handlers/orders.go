package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
	"github.com/luk3skyw4lker/order-pack-calculator/src/database/models"
	"github.com/luk3skyw4lker/order-pack-calculator/src/payload"
	"github.com/luk3skyw4lker/order-pack-calculator/src/utils"
)

type OrderService interface {
	CreateOrder(itemsCount int) (models.Order, error)
	GetOrder(orderID uuid.UUID) (models.Order, error)
	GetAllOrders() ([]models.Order, error)
}

type OrdersHandler struct {
	orderService OrderService
}

func NewOrdersHandler(orderService OrderService) *OrdersHandler {
	return &OrdersHandler{
		orderService: orderService,
	}
}

// CreateOrder godoc
//
//	@Summary		Create an order
//	@Description	Create an order with the specified number of items
//	@Tags			Orders
//	@Accept			json
//	@Produces		json
//	@Param			order	body		payload.CreateOrder	true	"the order to be created"
//	@Success		201			{object}	models.Order
//	@Failure		400			{object}	payload.ErrorResponse
//	@Failure		500			{object}	payload.ErrorResponse
//	@Router			/orders [post]
func (h *OrdersHandler) CreateOrder(ctx fiber.Ctx) error {
	input, err := utils.UnmarshalRequest[payload.CreateOrder](ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(payload.ErrorResponse{Message: "badly formed request"})
	}

	order, err := h.orderService.CreateOrder(input.ItemsCount)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(payload.ErrorResponse{Message: "failed to create order"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(order)
}

// GetOrder godoc
//
//	@Summary		Get an order by ID
//	@Description	Retrieve the details of an order using its ID
//	@Tags			Orders
//	@Accept			json
//	@Produces		json
//	@Param			order_id	path		string	true	"The ID of the order to retrieve"
//	@Success		200			{object}	models.Order
//	@Failure		400			{object}	payload.ErrorResponse
//	@Failure		404			{object}	payload.ErrorResponse
//	@Failure		500			{object}	payload.ErrorResponse
//	@Router			/orders/{order_id} [get]
func (h *OrdersHandler) GetOrder(ctx fiber.Ctx) error {
	orderIDStr := ctx.Params("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		log.Error("invalid order ID:", err)

		return ctx.Status(fiber.StatusBadRequest).JSON(payload.ErrorResponse{Message: "invalid order ID"})
	}

	order, err := h.orderService.GetOrder(orderID)
	if err != nil {
		if errors.Is(err, payload.ErrOrderNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(payload.ErrorResponse{Message: "order not found"})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(payload.ErrorResponse{Message: "failed to retrieve order"})
	}

	return ctx.Status(fiber.StatusOK).JSON(order)
}

// GetAllOrders godoc
//
//	@Summary		Get all orders
//	@Description	Retrieve a list of all orders
//	@Tags			Orders
//	@Accept			json
//	@Produces		json
//	@Success		200	{array}		models.Order
//	@Failure		500	{object}	payload.ErrorResponse
//	@Router			/orders [get]
func (h *OrdersHandler) GetAllOrders(ctx fiber.Ctx) error {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		if errors.Is(err, payload.ErrOrderNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(payload.ErrorResponse{Message: "no orders found"})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(payload.ErrorResponse{Message: "failed to retrieve orders"})
	}

	return ctx.Status(fiber.StatusOK).JSON(orders)
}
