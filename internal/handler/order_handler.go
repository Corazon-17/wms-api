package handler

import (
	"wms-api/internal/dto"
	"wms-api/internal/service"

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

	orders, total, err := h.service.GetOrders(c.Context(), q)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": orders,
		"meta": fiber.Map{
			"total": total,
			"page":  q.Page,
			"limit": q.Limit,
		},
	})
}

func (h *OrderHandler) GetOrder(c fiber.Ctx) error {

	type Request struct {
		OrderSN string `params:"order_sn"`
	}

	var req Request

	if err := c.Bind().Query(&req); err != nil {
		return err
	}

	order, err := h.service.GetOrder(c.Context(), req.OrderSN)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "order not found"})
	}

	return c.JSON(order)
}

func (h *OrderHandler) PickOrder(c fiber.Ctx) error {

	type Params struct {
		OrderSN string `params:"order_sn"`
	}

	var p Params
	if err := c.Bind().Query(&p); err != nil {
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
		OrderSN string `params:"order_sn"`
	}

	var p Params
	if err := c.Bind().Query(&p); err != nil {
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

	type Request struct {
		OrderSN   string `params:"order_sn"`
		ChannelId string `params:"channel_id"`
	}

	var req Request
	c.Bind().Query(&req)

	order, err := h.service.ShipOrder(c.Context(), req.OrderSN, req.ChannelId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(order)
}
