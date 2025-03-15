package handler

import (
	"gateway/models"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func Handlers(app *fiber.App) {
	apiUserService := app.Group("/users")

	apiUserService.Post("/create", func(c *fiber.Ctx) error {
		var input models.UserCreateInput
		// Парсинг тела запроса в структуру
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Неверный формат запроса",
			})
		}
		// Валидация входных данных
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка валидации: " + err.Error(),
			})
		}

		// Здесь можно добавить логику создания пользователя, например, запись в БД

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Пользователь успешно создан",
			"user":    input,
		})
	})

	apiUserService.Post("/login", func(c *fiber.Ctx) error {
		var input models.UserLoginInput
		// Парсинг тела запроса в структуру
		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Неверный формат запроса",
			})
		}
		// Валидация входных данных
		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Ошибка валидации: " + err.Error(),
			})
		}

		// Здесь можно добавить логику аутентификации пользователя

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Пользователь успешно авторизован",
		})
	})
}
