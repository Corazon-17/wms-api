package handler

import (
	"wms-api/internal/config"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

func (h *AuthHandler) Login(c fiber.Ctx) error {

	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest

	if err := c.Bind().Body(&req); err != nil {
		return err
	}

	if req.Email != h.cfg.WMSEmail ||
		req.Password != h.cfg.WMSPassword {

		return c.Status(401).JSON(fiber.Map{
			"error": "invalid credentials",
		})
	}

	return c.JSON(fiber.Map{
		"token": h.cfg.InternalAPIKey,
	})
}
