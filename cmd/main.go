package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"wms-api/internal/config"
	"wms-api/internal/database"
	"wms-api/internal/handler"
	"wms-api/internal/provider/marketplace"
	"wms-api/internal/repository"
	"wms-api/internal/route"
	"wms-api/internal/service"
)

func main() {

	cfg := config.Load()

	app := fiber.New()

	db := database.NewDB(cfg)

	orderRepo := repository.NewOrderRepository(db)

	marketplaceClient := marketplace.NewClient(cfg)

	orderService := service.NewOrderService(orderRepo, marketplaceClient)

	authHandler := handler.NewAuthHandler(cfg)
	orderHandler := handler.NewOrderHandler(orderService)
	webhookHandler := handler.NewWebhookHandler(orderService)

	api := app.Group("api")

	route.RegisterAuthRoutes(api, authHandler)
	route.RegisterOrderRoutes(api, orderHandler, cfg)
	route.RegisterWebhookRoutes(api, webhookHandler)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
