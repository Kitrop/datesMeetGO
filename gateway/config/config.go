package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartServer() *fiber.App{
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*",
		MaxAge: 43200,
	}))

	return app
}