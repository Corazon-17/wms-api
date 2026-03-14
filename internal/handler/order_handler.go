package handler

import (
	"wms-api/internal/dto"
	"wms-api/internal/service"
	"wms-api/pkg/pagination"

	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrders(c fiber.Ctx) error {

	var q dto.OrderQuery

	if err := c.Bind().Query(&q); err != nil {
		return err
	}

	if q.Page < 1 || q.PageSize < 1 {
		return c.Status(400).JSON(fiber.Map{
			"error": "page and pageSize must be greater than 0",
		})
	}

	orders, total, err := h.service.GetOrders(c.Context(), q)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": orders,
		"meta": pagination.GeneratePaginationMeta(q.Page, q.PageSize, total),
	})
}

func (h *OrderHandler) GetOrder(c fiber.Ctx) error {

	type Params struct {
		OrderSN string `uri:"order_sn"`
	}

	var p Params

	if err := c.Bind().URI(&p); err != nil {
		return err
	}

	order, err := h.service.GetOrder(c.Context(), p.OrderSN)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "order not found"})
	}

	return c.JSON(order)
}

func (h *OrderHandler) PickOrder(c fiber.Ctx) error {

	type Params struct {
		OrderSN string `uri:"order_sn"`
	}

	var p Params
	if err := c.Bind().URI(&p); err != nil {
		return err
	}

	err := h.service.PickOrder(c.Context(), p.OrderSN)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "order picking started",
	})
}

func (h *OrderHandler) PackOrder(c fiber.Ctx) error {

	type Params struct {
		OrderSN string `uri:"order_sn"`
	}

	var p Params
	if err := c.Bind().URI(&p); err != nil {
		return err
	}

	err := h.service.PackOrder(c.Context(), p.OrderSN)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "order packed",
	})
}

func (h *OrderHandler) ShipOrder(c fiber.Ctx) error {

	type Params struct {
		OrderSN string `uri:"order_sn"`
	}

	var p Params
	if err := c.Bind().URI(&p); err != nil {
		return err
	}

	order, err := h.service.ShipOrder(c.Context(), p.OrderSN, "JNE")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "order shipped",
		"data":    order,
	})
}

func (h *OrderHandler) GetWMSStatuses(c fiber.Ctx) error {

	result, err := h.service.GetWMSStatuses(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *OrderHandler) GetMarketplaceStatuses(c fiber.Ctx) error {

	result, err := h.service.GetMarketplaceStatuses(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *OrderHandler) GetShippingStatuses(c fiber.Ctx) error {

	result, err := h.service.GetShippingStatuses(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}

func (h *OrderHandler) GetOrderSummary(c fiber.Ctx) error {

	result, err := h.service.GetOrderSummary(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
