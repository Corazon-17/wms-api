package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

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

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{cfg.CorsAllowOrigins},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-api-key"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	db := database.NewDB(cfg)

	orderRepo := repository.NewOrderRepository(db)

	marketplaceClient := marketplace.NewClient(cfg)

	orderService := service.NewOrderService(orderRepo, marketplaceClient)

	go func() {
		err := orderService.SyncOrders(context.Background())
		if err != nil {
			log.Println("initial order sync failed:", err)
		}
	}()

	authHandler := handler.NewAuthHandler(cfg)
	orderHandler := handler.NewOrderHandler(orderService)
	webhookHandler := handler.NewWebhookHandler(orderService)

	api := app.Group("api")

	route.RegisterAuthRoutes(api, authHandler)
	route.RegisterOrderRoutes(api, orderHandler, cfg)
	route.RegisterWebhookRoutes(api, webhookHandler)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
