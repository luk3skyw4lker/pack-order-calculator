package main

import (
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/luk3skyw4lker/order-pack-calculator/src/config"
	_ "github.com/luk3skyw4lker/order-pack-calculator/src/docs"
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

	if err := app.Listen(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
