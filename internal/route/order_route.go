package route

import (
	"github.com/gofiber/fiber/v3"

	"wms-api/internal/config"
	"wms-api/internal/handler"
	"wms-api/internal/middleware"
)

func RegisterOrderRoutes(api fiber.Router, handler *handler.OrderHandler, cfg *config.Config) {

	protected := api.Group("orders", middleware.APIKeyAuth(cfg))

	protected.Get("wms-statuses", handler.GetWMSStatuses)
	protected.Get("marketplace-statuses", handler.GetMarketplaceStatuses)
	protected.Get("shipping-statuses", handler.GetShippingStatuses)

	protected.Get("", handler.GetOrders)

	protected.Get(":order_sn", handler.GetOrder)
	protected.Post(":order_sn/pick", handler.PickOrder)
	protected.Post(":order_sn/pack", handler.PackOrder)
	protected.Post(":order_sn/ship", handler.ShipOrder)
}
