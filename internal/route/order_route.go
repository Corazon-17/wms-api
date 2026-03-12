package route

import (
	"github.com/gofiber/fiber/v3"

	"wms-api/internal/config"
	"wms-api/internal/handler"
	"wms-api/internal/middleware"
)

func RegisterOrderRoutes(api fiber.Router, handler *handler.OrderHandler, cfg *config.Config) {

	protected := api.Group("", middleware.APIKeyAuth(cfg))

	protected.Get("/orders", handler.GetOrders)
	protected.Get("/orders/:order_sn", handler.GetOrder)
	protected.Post("/orders/:order_sn/pick", handler.PickOrder)
	protected.Post("/orders/:order_sn/pack", handler.PackOrder)
	protected.Post("/orders/:order_sn/ship", handler.ShipOrder)
}
