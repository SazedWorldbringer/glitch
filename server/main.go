package main

import (
	"log"
	"os"

	"github.com/SazedWorldbringer/url-shortener/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get(":url", routes.Resolve)
	app.Post("/api/v1", routes.Shorten)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load environment file.")
	}

	app := fiber.New()

	// middleware
	app.Use(logger.New())
	app.Use(cors.New())

	setupRoutes(app)

	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	// Listen on port 4000
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
