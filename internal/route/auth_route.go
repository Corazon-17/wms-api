package route

import (
	"github.com/gofiber/fiber/v3"

	"wms-api/internal/handler"
)

func RegisterAuthRoutes(api fiber.Router, handler *handler.AuthHandler) {

	api.Post("/auth/login", handler.Login)
}
