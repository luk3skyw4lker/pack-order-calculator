package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/luk3skyw4lker/order-pack-calculator/src/config"
	"github.com/luk3skyw4lker/order-pack-calculator/src/database"
	_ "github.com/luk3skyw4lker/order-pack-calculator/src/docs"
	"github.com/luk3skyw4lker/order-pack-calculator/src/internal/handlers"
	"github.com/luk3skyw4lker/order-pack-calculator/src/internal/repositories"
	"github.com/luk3skyw4lker/order-pack-calculator/src/internal/services"
	"github.com/luk3skyw4lker/order-pack-calculator/src/utils"
)

// @title Orders Calculation API
// @version 1.0
// @description This is a server for calculating order pack breakdowns.
// @termsOfService http://swagger.io/terms/
// @contact.name Lucas Lemos
// @contact.email lucashenriqueblemos@gmail.com
// @license.name MIT
// @license.url https://mit-license.org/
// @host orders-calculation.luk3skyw4lker.com
// @BasePath /
func main() {
	var cfg config.Config
	if err := config.LoadConfig(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := fiber.New()

	utils.InitDocs(app)

	database, err := database.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	ordersRepo := repositories.NewOrdersRepository(database)
	ordersService := services.NewOrdersService(ordersRepo)
	ordersHandler := handlers.NewOrdersHandler(ordersService)

	setupRoutes(app, ordersHandler)

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(app *fiber.App, ordersHandler *handlers.OrdersHandler) {
	app.Post("/orders", ordersHandler.CreateOrder)
	app.Get("/orders/:order_id", ordersHandler.GetOrder)
	app.Get("/orders", ordersHandler.GetAllOrders)
}
