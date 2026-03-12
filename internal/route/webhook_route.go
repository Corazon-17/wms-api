package route

import (
	"github.com/gofiber/fiber/v3"

	"wms-api/internal/handler"
)

func RegisterWebhookRoutes(api fiber.Router, handler *handler.WebhookHandler) {

	api.Post("/webhook/order-status", handler.OrderStatus)
	api.Post("/webhook/shipping-status", handler.ShippingStatus)
}
