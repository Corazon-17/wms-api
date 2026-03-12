package middleware

import (
	"wms-api/internal/config"

	"github.com/gofiber/fiber/v3"
)

func APIKeyAuth(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		key := c.Get("x-api-key")

		if key != cfg.InternalAPIKey {
			return c.Status(401).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		return c.Next()
	}
}
