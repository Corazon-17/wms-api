package handler

import (
	"wms-api/internal/service"

	"github.com/gofiber/fiber/v3"
)

type WebhookHandler struct {
	service *service.OrderService
}

func NewWebhookHandler(s *service.OrderService) *WebhookHandler {
	return &WebhookHandler{service: s}
}

func (h *WebhookHandler) OrderStatus(c fiber.Ctx) error {

	type Payload struct {
		OrderSN string `json:"orderSN"`
		Status  string `json:"status"`
	}

	var p Payload

	if err := c.Bind().Body(&p); err != nil {
		return err
	}

	err := h.service.UpdateMarketplaceStatus(
		c.Context(),
		p.OrderSN,
		p.Status,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "order status updated",
	})
}

func (h *WebhookHandler) ShippingStatus(c fiber.Ctx) error {

	type Payload struct {
		OrderSN       string `json:"orderSN"`
		ShippingState string `json:"shipping_state"`
	}

	var p Payload

	if err := c.Bind().Body(&p); err != nil {
		return err
	}

	err := h.service.UpdateShippingStatus(
		c.Context(),
		p.OrderSN,
		p.ShippingState,
	)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "shipping status updated",
	})
}
