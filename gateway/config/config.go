package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartServer() *fiber.App{
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:  "*", // Разрешенные домены
		AllowCredentials: true,
		AllowMethods:  "GET, POST, DELETE",
		MaxAge: 43200,
	}))

	return app
}