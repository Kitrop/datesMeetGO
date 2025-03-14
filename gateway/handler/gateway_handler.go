package handler

import "github.com/gofiber/fiber/v2"



func Handlers(app *fiber.App) {
	apiUserService := app.Group("/users")

	apiUserService.Post("/create", func(c *fiber.Ctx) error {
		return c.SendString("POST request")
	})
	apiUserService.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("POST request")
	})
}